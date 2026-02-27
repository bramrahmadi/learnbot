package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/learnbot/job-aggregator/internal/httpclient"
	"github.com/learnbot/job-aggregator/internal/model"
	"golang.org/x/net/html"
)

// CareerPageSelectors defines CSS-like selectors for extracting job data
// from a company's career page. Each field is a list of selectors to try
// in order (first match wins).
type CareerPageSelectors struct {
	// Container selector for each job listing
	JobContainer string `json:"job_container"`
	// Field selectors within each job container
	Title          string `json:"title"`
	Location       string `json:"location"`
	Department     string `json:"department"`
	EmploymentType string `json:"employment_type"`
	Description    string `json:"description"`
	ApplyURL       string `json:"apply_url"`
	PostedDate     string `json:"posted_date"`
	// Pagination
	NextPageSelector string `json:"next_page_selector"`
	// JSON API mode (if the career page uses a JSON API)
	APIEndpoint  string `json:"api_endpoint,omitempty"`
	APIJobsField string `json:"api_jobs_field,omitempty"`
}

// CareerPageScraper scrapes job postings from company career pages.
// It uses configurable CSS selectors stored in the database.
type CareerPageScraper struct {
	*BaseScraper
	page model.CompanyCareerPage
}

// NewCareerPageScraper creates a new career page scraper for a specific company.
func NewCareerPageScraper(page model.CompanyCareerPage, logger *log.Logger) (*CareerPageScraper, error) {
	cfg := httpclient.DefaultConfig()
	cfg.RequestsPerMinute = 8
	cfg.MaxRetries = 3
	cfg.RetryDelay = 3 * time.Second

	base, err := NewBaseScraper(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &CareerPageScraper{
		BaseScraper: base,
		page:        page,
	}, nil
}

func (s *CareerPageScraper) Source() model.JobSource { return model.SourceCompanyCareerPage }
func (s *CareerPageScraper) Name() string {
	return fmt.Sprintf("Career Page: %s", s.page.CompanyName)
}

// Scrape fetches job listings from the company's career page.
func (s *CareerPageScraper) Scrape(ctx context.Context, params model.SearchParams, jobs chan<- *model.ScrapedJob) error {
	var selectors CareerPageSelectors
	if err := json.Unmarshal(s.page.Selectors, &selectors); err != nil {
		return fmt.Errorf("invalid selectors config: %w", err)
	}

	s.Logger.Printf("[career_page] scraping %s (%s)", s.page.CompanyName, s.page.CareerPageURL)

	// Check robots.txt
	if !s.Robots.IsAllowed(ctx, s.page.CareerPageURL) {
		return fmt.Errorf("robots.txt disallows scraping %s", s.page.CareerPageURL)
	}

	// Use JSON API if configured
	if selectors.APIEndpoint != "" {
		return s.scrapeAPI(ctx, selectors, params, jobs)
	}

	// Otherwise scrape HTML
	return s.scrapeHTML(ctx, selectors, params, jobs)
}

// scrapeAPI fetches jobs from a JSON API endpoint.
func (s *CareerPageScraper) scrapeAPI(ctx context.Context, selectors CareerPageSelectors, params model.SearchParams, jobs chan<- *model.ScrapedJob) error {
	body, err := s.Client.GetBody(ctx, selectors.APIEndpoint, map[string]string{
		"Accept": "application/json",
	})
	if err != nil {
		return fmt.Errorf("fetch API: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return fmt.Errorf("parse API response: %w", err)
	}

	// Navigate to jobs field
	jobsField := selectors.APIJobsField
	if jobsField == "" {
		jobsField = "jobs"
	}

	jobsData, ok := data[jobsField]
	if !ok {
		return fmt.Errorf("jobs field %q not found in API response", jobsField)
	}

	jobsList, ok := jobsData.([]interface{})
	if !ok {
		return fmt.Errorf("jobs field is not an array")
	}

	for _, item := range jobsList {
		jobMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		job := s.extractAPIJob(jobMap)
		if job == nil {
			continue
		}

		// Filter by query if provided
		if params.Query != "" {
			if !strings.Contains(strings.ToLower(job.Title), strings.ToLower(params.Query)) &&
				!strings.Contains(strings.ToLower(job.Description), strings.ToLower(params.Query)) {
				continue
			}
		}

		select {
		case jobs <- job:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// extractAPIJob converts a JSON job object to a ScrapedJob.
func (s *CareerPageScraper) extractAPIJob(data map[string]interface{}) *model.ScrapedJob {
	job := &model.ScrapedJob{
		Source:         model.SourceCompanyCareerPage,
		CompanyName:    s.page.CompanyName,
		SalaryCurrency: "USD",
	}

	// Common field names across different career page APIs
	titleFields := []string{"title", "job_title", "name", "position"}
	locationFields := []string{"location", "city", "office", "work_location"}
	descFields := []string{"description", "job_description", "summary", "content"}
	urlFields := []string{"url", "apply_url", "link", "job_url", "absolute_url"}
	idFields := []string{"id", "job_id", "req_id", "requisition_id"}

	for _, f := range titleFields {
		if v, ok := data[f].(string); ok && v != "" {
			job.Title = v
			break
		}
	}
	for _, f := range locationFields {
		if v, ok := data[f].(string); ok && v != "" {
			job.LocationRaw = v
			parseLocation(v, job)
			break
		}
	}
	for _, f := range descFields {
		if v, ok := data[f].(string); ok && v != "" {
			job.Description = CleanText(v)
			break
		}
	}
	for _, f := range urlFields {
		if v, ok := data[f].(string); ok && v != "" {
			job.ApplicationURL = v
			break
		}
	}
	for _, f := range idFields {
		if v, ok := data[f]; ok {
			job.ExternalID = fmt.Sprintf("%v", v)
			break
		}
	}

	if job.Title == "" {
		return nil
	}

	if job.ApplicationURL == "" {
		job.ApplicationURL = s.page.CareerPageURL
	}

	job.LocationType = ExtractLocationType(job.LocationRaw, job.Description)
	job.EmploymentType = ExtractEmploymentType(job.Title, job.Description)
	job.ExperienceLevel = ExtractExperienceLevel(job.Title, job.Description)
	job.RequiredSkills, job.PreferredSkills = ExtractSkillsFromText(job.Description)

	return job
}

// scrapeHTML fetches and parses HTML from the career page.
func (s *CareerPageScraper) scrapeHTML(ctx context.Context, selectors CareerPageSelectors, params model.SearchParams, jobs chan<- *model.ScrapedJob) error {
	currentURL := s.page.CareerPageURL
	maxPages := 10

	for page := 0; page < maxPages; page++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		body, err := s.Client.GetBody(ctx, currentURL, nil)
		if err != nil {
			return fmt.Errorf("fetch page %d: %w", page, err)
		}

		pageJobs, nextURL, err := s.parseCareerPageHTML(body, selectors, currentURL)
		if err != nil {
			s.Logger.Printf("[career_page] parse error on page %d: %v", page, err)
			break
		}

		for _, job := range pageJobs {
			if params.Query != "" {
				if !strings.Contains(strings.ToLower(job.Title), strings.ToLower(params.Query)) {
					continue
				}
			}
			select {
			case jobs <- job:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		s.Logger.Printf("[career_page] %s page %d: found %d jobs", s.page.CompanyName, page, len(pageJobs))

		if nextURL == "" || nextURL == currentURL {
			break
		}
		currentURL = nextURL

		// Polite delay
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return nil
}

// parseCareerPageHTML parses HTML using the configured selectors.
func (s *CareerPageScraper) parseCareerPageHTML(body string, selectors CareerPageSelectors, baseURL string) ([]*model.ScrapedJob, string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, "", fmt.Errorf("parse HTML: %w", err)
	}

	var jobs []*model.ScrapedJob
	var nextURL string

	// Find job containers
	containers := findBySelector(doc, selectors.JobContainer)
	for _, container := range containers {
		job := &model.ScrapedJob{
			Source:         model.SourceCompanyCareerPage,
			CompanyName:    s.page.CompanyName,
			SalaryCurrency: "USD",
		}

		// Extract fields using selectors
		if nodes := findBySelector(container, selectors.Title); len(nodes) > 0 {
			job.Title = CleanText(extractText(nodes[0]))
		}
		if nodes := findBySelector(container, selectors.Location); len(nodes) > 0 {
			job.LocationRaw = CleanText(extractText(nodes[0]))
			parseLocation(job.LocationRaw, job)
		}
		if nodes := findBySelector(container, selectors.Description); len(nodes) > 0 {
			job.Description = CleanText(extractText(nodes[0]))
		}
		if nodes := findBySelector(container, selectors.ApplyURL); len(nodes) > 0 {
			href := getAttr(nodes[0], "href")
			if href != "" {
				job.ApplicationURL = resolveURL(baseURL, href)
			}
		}
		if nodes := findBySelector(container, selectors.PostedDate); len(nodes) > 0 {
			dateStr := CleanText(extractText(nodes[0]))
			if t := ParseRelativeDate(dateStr); t != nil {
				job.PostedAt = t
			}
		}

		if job.Title == "" {
			continue
		}
		if job.ApplicationURL == "" {
			job.ApplicationURL = baseURL
		}

		job.LocationType = ExtractLocationType(job.LocationRaw, job.Description)
		job.EmploymentType = ExtractEmploymentType(job.Title, job.Description)
		job.ExperienceLevel = ExtractExperienceLevel(job.Title, job.Description)
		job.RequiredSkills, job.PreferredSkills = ExtractSkillsFromText(job.Description)

		jobs = append(jobs, job)
	}

	// Find next page URL
	if selectors.NextPageSelector != "" {
		if nodes := findBySelector(doc, selectors.NextPageSelector); len(nodes) > 0 {
			href := getAttr(nodes[0], "href")
			if href != "" {
				nextURL = resolveURL(baseURL, href)
			}
		}
	}

	return jobs, nextURL, nil
}

// findBySelector finds HTML nodes matching a simple CSS-like selector.
// Supports: element, .class, #id, element.class
func findBySelector(root *html.Node, selector string) []*html.Node {
	if selector == "" {
		return nil
	}

	var matches []*html.Node
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && matchesSelector(n, selector) {
			matches = append(matches, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(root)
	return matches
}

// matchesSelector checks if a node matches a simple selector.
func matchesSelector(n *html.Node, selector string) bool {
	if selector == "" {
		return false
	}

	// Handle .class selector
	if strings.HasPrefix(selector, ".") {
		return hasClass(n, selector[1:])
	}

	// Handle #id selector
	if strings.HasPrefix(selector, "#") {
		return getAttr(n, "id") == selector[1:]
	}

	// Handle element.class selector
	if idx := strings.Index(selector, "."); idx > 0 {
		elem := selector[:idx]
		class := selector[idx+1:]
		return n.Data == elem && hasClass(n, class)
	}

	// Handle element selector
	return n.Data == selector
}

// resolveURL resolves a relative URL against a base URL.
func resolveURL(base, href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	if strings.HasPrefix(href, "//") {
		// Protocol-relative URL
		if strings.HasPrefix(base, "https") {
			return "https:" + href
		}
		return "http:" + href
	}
	if strings.HasPrefix(href, "/") {
		// Absolute path
		parts := strings.SplitN(base, "/", 4)
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2] + href
		}
	}
	// Relative path
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		return base[:idx+1] + href
	}
	return href
}
