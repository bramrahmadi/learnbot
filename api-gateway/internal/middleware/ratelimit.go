// Package middleware – ratelimit.go implements a token-bucket rate limiter.
package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a per-IP token bucket rate limiter.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64 // tokens per second
	capacity float64 // max tokens
	cleanup  time.Duration
}

// tokenBucket holds the state for a single IP's rate limit.
type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

// NewRateLimiter creates a new RateLimiter.
//
// Parameters:
//   - requestsPerSecond: the sustained request rate allowed per IP
//   - burst: the maximum burst size (capacity of the bucket)
func NewRateLimiter(requestsPerSecond, burst float64) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     requestsPerSecond,
		capacity: burst,
		cleanup:  5 * time.Minute,
	}
	go rl.cleanupLoop()
	return rl
}

// Allow returns true if the request from the given IP is allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.buckets[ip]
	if !exists {
		rl.buckets[ip] = &tokenBucket{
			tokens:   rl.capacity - 1, // consume one token immediately
			lastSeen: now,
		}
		return true
	}

	// Refill tokens based on elapsed time.
	elapsed := now.Sub(bucket.lastSeen).Seconds()
	bucket.tokens += elapsed * rl.rate
	if bucket.tokens > rl.capacity {
		bucket.tokens = rl.capacity
	}
	bucket.lastSeen = now

	if bucket.tokens < 1 {
		return false
	}
	bucket.tokens--
	return true
}

// Middleware returns an HTTP middleware that applies rate limiting.
// Returns 429 Too Many Requests when the limit is exceeded.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractClientIP(r)
		if !rl.Allow(ip) {
			w.Header().Set("Retry-After", "1")
			writeJSONError(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED",
				"too many requests, please slow down")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// cleanupLoop periodically removes stale buckets to prevent memory leaks.
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-rl.cleanup)
		for ip, bucket := range rl.buckets {
			if bucket.lastSeen.Before(cutoff) {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// extractClientIP extracts the client IP from the request.
// Respects X-Forwarded-For and X-Real-IP headers for proxied requests.
func extractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain.
		parts := splitAndTrim(xff, ",")
		if len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr (strip port).
	addr := r.RemoteAddr
	if idx := lastIndex(addr, ':'); idx >= 0 {
		return addr[:idx]
	}
	return addr
}

// splitAndTrim splits a string by sep and trims whitespace from each part.
func splitAndTrim(s, sep string) []string {
	parts := make([]string, 0)
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			part := trimSpace(s[start:i])
			parts = append(parts, part)
			start = i + len(sep)
		}
	}
	parts = append(parts, trimSpace(s[start:]))
	return parts
}

// trimSpace removes leading and trailing whitespace.
func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// lastIndex returns the index of the last occurrence of c in s.
func lastIndex(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
