package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/learnbot/job-aggregator/internal/httpclient"
	"github.com/learnbot/job-aggregator/internal/model"
	"golang.org/x/net/html"
)

// LinkedInScraper scrapes job postings from LinkedIn Jobs.
// Note: LinkedIn's official API requires partnership approval.
// This scraper uses the public job search endpoint which is accessible
// without authentication for basic job listings.
// Always respect LinkedIn's robots.txt and Terms of Service.
type LinkedInScraper struct {
	*BaseScraper
	baseURL string
}

// NewLinkedInScraper creates a new LinkedIn scraper.
func NewLinkedInScraper(logger *log.Logger) (*LinkedInScraper, error) {
	cfg := httpclient.DefaultConfig()
	cfg.RequestsPerMinute = 5 // LinkedIn is strict about rate limiting
	cfg.MaxRetries = 3
	cfg.RetryDelay = 5 * time.Second
	cfg.UserAgent = "Mozilla/5.0 (compatible; LearnBot/1.0; +https://learnbot.io)"

	base, err := NewBaseScraper(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &LinkedInScraper{
		BaseScraper: base,
		baseURL:     "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
	}, nil
}

func (s *LinkedInScraper) Source() model.JobSource { return model.SourceLinkedIn }
func (s *LinkedInScraper) Name() string            { return "LinkedIn Jobs" }

// Scrape fetches job listings from LinkedIn's public job search.
func (s *LinkedInScraper) Scrape(ctx context.Context, params model.SearchParams, jobs chan<- *model.ScrapedJob) error {
	if params.PageSize <= 0 {
		params.PageSize = 25
	}

	s.Logger.Printf("[linkedin] starting scrape: query=%q location=%q", params.Query, params.Location)

	page := 0
	maxPages := 10 // LinkedIn limits public access to ~250 results

	for page < maxPages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		pageJobs, hasMore, err := s.scrapePage(ctx, params, page*params.PageSize)
		if err != nil {
			s.Logger.Printf("[linkedin] page %d error: %v", page, err)
			if page == 0 {
				return fmt.Errorf("linkedin scrape failed: %w", err)
			}
			break // Stop on error after first page
		}

		for _, job := range pageJobs {
			select {
			case jobs <- job:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		s.Logger.Printf("[linkedin] page %d: found %d jobs", page, len(pageJobs))

		if !hasMore || len(pageJobs) == 0 {
			break
		}
		page++

		// Polite delay between pages
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return nil
}

// scrapePage fetches a single page of LinkedIn job results.
func (s *LinkedInScraper) scrapePage(ctx context.Context, params model.SearchParams, start int) ([]*model.ScrapedJob, bool, error) {
	// Build search URL
	q := url.Values{}
	q.Set("keywords", params.Query)
	if params.Location != "" {
		q.Set("location", params.Location)
	}
	if params.Remote {
		q.Set("f_WT", "2") // LinkedIn's remote filter
	}
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("count", fmt.Sprintf("%d", params.PageSize))

	searchURL := s.baseURL + "?" + q.Encode()

	// Check robots.txt
	if !s.Robots.IsAllowed(ctx, searchURL) {
		return nil, false, fmt.Errorf("robots.txt disallows scraping %s", searchURL)
	}

	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml",
		"Accept-Language": "en-US,en;q=0.9",
		"Referer":         "https://www.linkedin.com/jobs/search/",
	}

	body, err := s.Client.GetBody(ctx, searchURL, headers)
	if err != nil {
		return nil, false, fmt.Errorf("fetch page: %w", err)
	}

	jobs, err := parseLinkedInHTML(body)
	if err != nil {
		return nil, false, fmt.Errorf("parse HTML: %w", err)
	}

	hasMore := len(jobs) >= params.PageSize
	return jobs, hasMore, nil
}

// parseLinkedInHTML parses LinkedIn job listing HTML and extracts job data.
func parseLinkedInHTML(body string) ([]*model.ScrapedJob, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	var jobs []*model.ScrapedJob
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// LinkedIn job cards have class "base-card" or "job-search-card"
			if hasClass(n, "base-card") || hasClass(n, "job-search-card") {
				if job := extractLinkedInJobCard(n); job != nil {
					jobs = append(jobs, job)
				}
				return // Don't recurse into job cards
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return jobs, nil
}

// extractLinkedInJobCard extracts job data from a LinkedIn job card HTML node.
func extractLinkedInJobCard(n *html.Node) *model.ScrapedJob {
	job := &model.ScrapedJob{
		Source:         model.SourceLinkedIn,
		SalaryCurrency: "USD",
	}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch {
			case node.Data == "a" && hasClass(node, "base-card__full-link"):
				job.ApplicationURL = getAttr(node, "href")
				// Extract external ID from URL
				if parts := strings.Split(job.ApplicationURL, "-"); len(parts) > 0 {
					job.ExternalID = parts[len(parts)-1]
					// Clean up query params
					if idx := strings.Index(job.ExternalID, "?"); idx >= 0 {
						job.ExternalID = job.ExternalID[:idx]
					}
				}

			case hasClass(node, "base-search-card__title"):
				job.Title = CleanText(extractText(node))

			case hasClass(node, "base-search-card__subtitle"):
				job.CompanyName = CleanText(extractText(node))

			case hasClass(node, "job-search-card__location"):
				job.LocationRaw = CleanText(extractText(node))
				parseLocation(job.LocationRaw, job)

			case hasClass(node, "job-search-card__listdate"):
				dateStr := getAttr(node, "datetime")
				if dateStr == "" {
					dateStr = CleanText(extractText(node))
				}
				if t := ParseRelativeDate(dateStr); t != nil {
					job.PostedAt = t
				}

			case hasClass(node, "job-search-card__salary-info"):
				job.SalaryRaw = CleanText(extractText(node))
				job.SalaryMin, job.SalaryMax, job.SalaryCurrency = ParseSalary(job.SalaryRaw)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	if job.Title == "" || job.CompanyName == "" {
		return nil
	}

	// Infer fields from available data
	job.LocationType = ExtractLocationType(job.LocationRaw, "")
	job.EmploymentType = model.EmploymentFullTime
	job.ExperienceLevel = ExtractExperienceLevel(job.Title, "")

	return job
}

// LinkedInJobDetail holds detailed job information from a job detail page.
type LinkedInJobDetail struct {
	Description string
	Skills      []string
}

// FetchJobDetail fetches the full job description from a LinkedIn job page.
func (s *LinkedInScraper) FetchJobDetail(ctx context.Context, jobURL string) (*LinkedInJobDetail, error) {
	// Use the JSON API endpoint for job details
	// Extract job ID from URL
	parts := strings.Split(jobURL, "-")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid job URL: %s", jobURL)
	}
	jobID := parts[len(parts)-1]
	if idx := strings.Index(jobID, "?"); idx >= 0 {
		jobID = jobID[:idx]
	}

	apiURL := fmt.Sprintf("https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/%s", jobID)

	body, err := s.Client.GetBody(ctx, apiURL, map[string]string{
		"Accept": "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("fetch job detail: %w", err)
	}

	// Try JSON first
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jsonData); err == nil {
		detail := &LinkedInJobDetail{}
		if desc, ok := jsonData["description"].(map[string]interface{}); ok {
			if text, ok := desc["text"].(string); ok {
				detail.Description = CleanText(text)
			}
		}
		return detail, nil
	}

	// Fall back to HTML parsing
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse detail HTML: %w", err)
	}

	detail := &LinkedInJobDetail{}
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if hasClass(n, "show-more-less-html__markup") || hasClass(n, "description__text") {
				detail.Description = CleanText(extractText(n))
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return detail, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// HTML utility functions
// ─────────────────────────────────────────────────────────────────────────────

// hasClass returns true if the HTML node has the given CSS class.
func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, c := range strings.Fields(attr.Val) {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

// getAttr returns the value of an HTML attribute.
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// extractText recursively extracts all text content from an HTML node.
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractText(c))
	}
	return sb.String()
}

// parseLocation splits a raw location string into city, state, country.
func parseLocation(raw string, job *model.ScrapedJob) {
	parts := strings.Split(raw, ",")
	switch len(parts) {
	case 1:
		job.LocationCity = strings.TrimSpace(parts[0])
	case 2:
		job.LocationCity = strings.TrimSpace(parts[0])
		job.LocationState = strings.TrimSpace(parts[1])
	case 3:
		job.LocationCity = strings.TrimSpace(parts[0])
		job.LocationState = strings.TrimSpace(parts[1])
		job.LocationCountry = strings.TrimSpace(parts[2])
	}
}
