package scraper

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/learnbot/job-aggregator/internal/httpclient"
	"github.com/learnbot/job-aggregator/internal/model"
	"golang.org/x/net/html"
)

// IndeedScraper scrapes job postings from Indeed.
// Note: Indeed's official API (Indeed Publisher API) requires approval.
// This scraper uses the public job search which is accessible without auth.
// Always respect Indeed's robots.txt and Terms of Service.
type IndeedScraper struct {
	*BaseScraper
	baseURL string
}

// NewIndeedScraper creates a new Indeed scraper.
func NewIndeedScraper(logger *log.Logger) (*IndeedScraper, error) {
	cfg := httpclient.DefaultConfig()
	cfg.RequestsPerMinute = 6
	cfg.MaxRetries = 3
	cfg.RetryDelay = 4 * time.Second
	cfg.UserAgent = "Mozilla/5.0 (compatible; LearnBot/1.0; +https://learnbot.io)"

	base, err := NewBaseScraper(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &IndeedScraper{
		BaseScraper: base,
		baseURL:     "https://www.indeed.com/jobs",
	}, nil
}

func (s *IndeedScraper) Source() model.JobSource { return model.SourceIndeed }
func (s *IndeedScraper) Name() string            { return "Indeed" }

// Scrape fetches job listings from Indeed.
func (s *IndeedScraper) Scrape(ctx context.Context, params model.SearchParams, jobs chan<- *model.ScrapedJob) error {
	if params.PageSize <= 0 {
		params.PageSize = 15 // Indeed shows 15 results per page
	}

	s.Logger.Printf("[indeed] starting scrape: query=%q location=%q", params.Query, params.Location)

	page := 0
	maxPages := 20

	for page < maxPages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		pageJobs, hasMore, err := s.scrapePage(ctx, params, page*params.PageSize)
		if err != nil {
			s.Logger.Printf("[indeed] page %d error: %v", page, err)
			if page == 0 {
				return fmt.Errorf("indeed scrape failed: %w", err)
			}
			break
		}

		for _, job := range pageJobs {
			select {
			case jobs <- job:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		s.Logger.Printf("[indeed] page %d: found %d jobs", page, len(pageJobs))

		if !hasMore || len(pageJobs) == 0 {
			break
		}
		page++

		// Polite delay
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}

	return nil
}

// scrapePage fetches a single page of Indeed job results.
func (s *IndeedScraper) scrapePage(ctx context.Context, params model.SearchParams, start int) ([]*model.ScrapedJob, bool, error) {
	q := url.Values{}
	q.Set("q", params.Query)
	if params.Location != "" {
		q.Set("l", params.Location)
	}
	if params.Remote {
		q.Set("remotejob", "032b3046-06a3-4876-8dfd-474eb5e7ed11") // Indeed's remote filter
	}
	q.Set("start", fmt.Sprintf("%d", start))

	searchURL := s.baseURL + "?" + q.Encode()

	if !s.Robots.IsAllowed(ctx, searchURL) {
		return nil, false, fmt.Errorf("robots.txt disallows scraping %s", searchURL)
	}

	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml",
		"Accept-Language": "en-US,en;q=0.9",
		"Referer":         "https://www.indeed.com/",
	}

	body, err := s.Client.GetBody(ctx, searchURL, headers)
	if err != nil {
		return nil, false, fmt.Errorf("fetch page: %w", err)
	}

	jobs, err := parseIndeedHTML(body)
	if err != nil {
		return nil, false, fmt.Errorf("parse HTML: %w", err)
	}

	hasMore := len(jobs) >= params.PageSize
	return jobs, hasMore, nil
}

// parseIndeedHTML parses Indeed job listing HTML.
func parseIndeedHTML(body string) ([]*model.ScrapedJob, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	var jobs []*model.ScrapedJob
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Indeed job cards have data-jk attribute (job key)
			if n.Data == "div" || n.Data == "li" {
				if jk := getAttr(n, "data-jk"); jk != "" {
					if job := extractIndeedJobCard(n, jk); job != nil {
						jobs = append(jobs, job)
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return jobs, nil
}

// extractIndeedJobCard extracts job data from an Indeed job card.
func extractIndeedJobCard(n *html.Node, jobKey string) *model.ScrapedJob {
	job := &model.ScrapedJob{
		Source:         model.SourceIndeed,
		ExternalID:     jobKey,
		ApplicationURL: fmt.Sprintf("https://www.indeed.com/viewjob?jk=%s", jobKey),
		SalaryCurrency: "USD",
	}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch {
			case node.Data == "h2" && hasClass(node, "jobTitle"):
				job.Title = CleanText(extractText(node))

			case node.Data == "span" && hasClass(node, "companyName"):
				job.CompanyName = CleanText(extractText(node))

			case node.Data == "div" && hasClass(node, "companyLocation"):
				job.LocationRaw = CleanText(extractText(node))
				parseLocation(job.LocationRaw, job)

			case node.Data == "div" && hasClass(node, "salary-snippet-container"):
				job.SalaryRaw = CleanText(extractText(node))
				job.SalaryMin, job.SalaryMax, job.SalaryCurrency = ParseSalary(job.SalaryRaw)

			case node.Data == "div" && hasClass(node, "job-snippet"):
				job.Description = CleanText(extractText(node))

			case node.Data == "span" && hasClass(node, "date"):
				dateStr := CleanText(extractText(node))
				if t := ParseRelativeDate(dateStr); t != nil {
					job.PostedAt = t
				}

			case node.Data == "div" && hasClass(node, "metadata"):
				// Extract employment type and other metadata
				text := strings.ToLower(CleanText(extractText(node)))
				if strings.Contains(text, "full-time") {
					job.EmploymentType = model.EmploymentFullTime
				} else if strings.Contains(text, "part-time") {
					job.EmploymentType = model.EmploymentPartTime
				} else if strings.Contains(text, "contract") {
					job.EmploymentType = model.EmploymentContract
				}
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

	// Infer fields
	job.LocationType = ExtractLocationType(job.LocationRaw, job.Description)
	if job.EmploymentType == "" {
		job.EmploymentType = ExtractEmploymentType(job.Title, job.Description)
	}
	job.ExperienceLevel = ExtractExperienceLevel(job.Title, job.Description)

	return job
}
