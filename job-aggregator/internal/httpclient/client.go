// Package httpclient provides a rate-limited HTTP client with retry logic
// for web scraping operations.
package httpclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// Config holds configuration for the HTTP client.
type Config struct {
	// Rate limiting
	RequestsPerMinute int
	// Retry settings
	MaxRetries    int
	RetryDelay    time.Duration
	RetryMaxDelay time.Duration
	// Timeouts
	RequestTimeout time.Duration
	// User agent
	UserAgent string
	// Optional proxy URL
	ProxyURL string
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		RequestsPerMinute: 10,
		MaxRetries:        3,
		RetryDelay:        2 * time.Second,
		RetryMaxDelay:     30 * time.Second,
		RequestTimeout:    30 * time.Second,
		UserAgent:         "LearnBot-JobAggregator/1.0 (https://learnbot.io; jobs@learnbot.io)",
	}
}

// Client is a rate-limited HTTP client with retry logic.
type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	config     Config
	logger     *log.Logger
}

// New creates a new Client with the given configuration.
func New(cfg Config, logger *log.Logger) (*Client, error) {
	if cfg.RequestsPerMinute <= 0 {
		cfg.RequestsPerMinute = 10
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 2 * time.Second
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 30 * time.Second
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	// Configure proxy if provided
	if cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   cfg.RequestTimeout,
	}

	// Rate limiter: tokens per second = requests per minute / 60
	tokensPerSecond := float64(cfg.RequestsPerMinute) / 60.0
	limiter := rate.NewLimiter(rate.Limit(tokensPerSecond), cfg.RequestsPerMinute)

	if logger == nil {
		logger = log.Default()
	}

	return &Client{
		httpClient: httpClient,
		limiter:    limiter,
		config:     cfg,
		logger:     logger,
	}, nil
}

// Get performs a GET request with rate limiting and retry logic.
func (c *Client) Get(ctx context.Context, rawURL string, headers map[string]string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, rawURL, nil, headers)
}

// Post performs a POST request with rate limiting and retry logic.
func (c *Client) Post(ctx context.Context, rawURL string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, rawURL, body, headers)
}

// GetBody performs a GET request and returns the response body as a string.
func (c *Client) GetBody(ctx context.Context, rawURL string, headers map[string]string) (string, error) {
	resp, err := c.Get(ctx, rawURL, headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	return string(body), nil
}

// do executes an HTTP request with rate limiting and exponential backoff retry.
func (c *Client) do(ctx context.Context, method, rawURL string, body io.Reader, headers map[string]string) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := c.backoffDelay(attempt)
			c.logger.Printf("retry %d/%d for %s %s (delay: %v): %v",
				attempt, c.config.MaxRetries, method, rawURL, delay, lastErr)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Wait for rate limiter
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}

		resp, err := c.executeRequest(ctx, method, rawURL, body, headers)
		if err != nil {
			lastErr = err
			if !isRetryableError(err) {
				return nil, err
			}
			continue
		}

		// Check for retryable HTTP status codes
		if isRetryableStatus(resp.StatusCode) {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			if resp.StatusCode == http.StatusTooManyRequests {
				// Respect Retry-After header if present
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					c.logger.Printf("rate limited by server, Retry-After: %s", retryAfter)
				}
			}
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded for %s %s: %w", method, rawURL, lastErr)
}

// executeRequest performs a single HTTP request.
func (c *Client) executeRequest(ctx context.Context, method, rawURL string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set default headers
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	// Override with caller-provided headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}

// backoffDelay calculates exponential backoff delay for retry attempts.
func (c *Client) backoffDelay(attempt int) time.Duration {
	delay := float64(c.config.RetryDelay) * math.Pow(2, float64(attempt-1))
	maxDelay := float64(c.config.RetryMaxDelay)
	if delay > maxDelay {
		delay = maxDelay
	}
	return time.Duration(delay)
}

// isRetryableError returns true if the error warrants a retry.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"EOF",
		"timeout",
		"temporary failure",
		"no such host",
	}
	for _, e := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), e) {
			return true
		}
	}
	return false
}

// isRetryableStatus returns true if the HTTP status code warrants a retry.
func isRetryableStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests,      // 429
		http.StatusInternalServerError,   // 500
		http.StatusBadGateway,            // 502
		http.StatusServiceUnavailable,    // 503
		http.StatusGatewayTimeout:        // 504
		return true
	}
	return false
}

// RobotsChecker checks robots.txt compliance.
type RobotsChecker struct {
	client    *Client
	cache     map[string]string // domain -> robots.txt content
	userAgent string
}

// NewRobotsChecker creates a new RobotsChecker.
func NewRobotsChecker(client *Client, userAgent string) *RobotsChecker {
	return &RobotsChecker{
		client:    client,
		cache:     make(map[string]string),
		userAgent: userAgent,
	}
}

// IsAllowed checks if the given URL is allowed by robots.txt.
// Returns true if allowed or if robots.txt cannot be fetched.
func (rc *RobotsChecker) IsAllowed(ctx context.Context, rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return true // allow if URL is unparseable
	}

	domain := parsed.Scheme + "://" + parsed.Host
	robotsURL := domain + "/robots.txt"

	content, ok := rc.cache[domain]
	if !ok {
		body, err := rc.client.GetBody(ctx, robotsURL, nil)
		if err != nil {
			rc.cache[domain] = "" // cache empty to avoid repeated failures
			return true           // allow if robots.txt not accessible
		}
		rc.cache[domain] = body
		content = body
	}

	if content == "" {
		return true
	}

	return isPathAllowed(content, rc.userAgent, parsed.Path)
}

// isPathAllowed parses robots.txt content and checks if the path is allowed.
// It collects all disallow rules that apply to the given user agent (including wildcard).
func isPathAllowed(robotsTxt, userAgent, path string) bool {
	lines := strings.Split(robotsTxt, "\n")

	type block struct {
		agents    []string
		disallows []string
	}

	var blocks []block
	var current *block

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "user-agent:") {
			agent := strings.TrimSpace(line[len("user-agent:"):])
			if current == nil || len(current.disallows) > 0 {
				// Start a new block
				if current != nil {
					blocks = append(blocks, *current)
				}
				current = &block{}
			}
			current.agents = append(current.agents, agent)
		} else if strings.HasPrefix(lower, "disallow:") {
			if current != nil {
				disallowedPath := strings.TrimSpace(line[len("disallow:"):])
				if disallowedPath != "" {
					current.disallows = append(current.disallows, disallowedPath)
				}
			}
		} else if line == "" && current != nil && len(current.disallows) > 0 {
			blocks = append(blocks, *current)
			current = nil
		}
	}
	if current != nil {
		blocks = append(blocks, *current)
	}

	for _, b := range blocks {
		applies := false
		for _, agent := range b.agents {
			if agent == "*" || strings.EqualFold(agent, userAgent) {
				applies = true
				break
			}
		}
		if !applies {
			continue
		}
		for _, d := range b.disallows {
			if strings.HasPrefix(path, d) {
				return false
			}
		}
	}
	return true
}
