// Package scraper defines the Scraper interface and common utilities
// for all job source scrapers.
package scraper

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/learnbot/job-aggregator/internal/httpclient"
	"github.com/learnbot/job-aggregator/internal/model"
)

// Scraper is the interface that all job source scrapers must implement.
type Scraper interface {
	// Source returns the job source identifier.
	Source() model.JobSource

	// Scrape fetches job postings matching the given search parameters.
	// It sends scraped jobs to the jobs channel and returns when done.
	Scrape(ctx context.Context, params model.SearchParams, jobs chan<- *model.ScrapedJob) error

	// Name returns a human-readable name for the scraper.
	Name() string
}

// BaseScraper provides common functionality for all scrapers.
type BaseScraper struct {
	Client *httpclient.Client
	Logger *log.Logger
	Robots *httpclient.RobotsChecker
}

// NewBaseScraper creates a new BaseScraper with the given configuration.
func NewBaseScraper(cfg httpclient.Config, logger *log.Logger) (*BaseScraper, error) {
	client, err := httpclient.New(cfg, logger)
	if err != nil {
		return nil, err
	}
	robots := httpclient.NewRobotsChecker(client, cfg.UserAgent)
	return &BaseScraper{
		Client: client,
		Logger: logger,
		Robots: robots,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Text extraction utilities
// ─────────────────────────────────────────────────────────────────────────────

var (
	// Salary patterns
	salaryRangeRe = regexp.MustCompile(
		`(?i)\$?\s*(\d{1,3}(?:,\d{3})*(?:\.\d+)?)\s*[kK]?\s*[-–—to]+\s*\$?\s*(\d{1,3}(?:,\d{3})*(?:\.\d+)?)\s*[kK]?`)
	salarySingleRe = regexp.MustCompile(
		`(?i)\$\s*(\d{1,3}(?:,\d{3})*(?:\.\d+)?)\s*[kK]?`)

	// Experience level keywords
	seniorKeywords  = []string{"senior", "sr.", "sr ", "lead", "principal", "staff", "architect"}
	midKeywords     = []string{"mid", "mid-level", "intermediate", "ii", "iii"}
	entryKeywords   = []string{"junior", "jr.", "jr ", "entry", "associate", "graduate", "new grad"}
	internKeywords  = []string{"intern", "internship", "co-op", "coop"}
	execKeywords    = []string{"director", "vp", "vice president", "head of", "chief", "cto", "ceo", "cfo"}

	// Remote keywords
	remoteKeywords = []string{"remote", "work from home", "wfh", "distributed", "anywhere"}
	hybridKeywords = []string{"hybrid", "flexible", "partially remote"}

	// Employment type keywords
	contractKeywords   = []string{"contract", "contractor", "freelance", "consulting", "1099"}
	partTimeKeywords   = []string{"part-time", "part time", "parttime"}
	internshipKeywords = []string{"intern", "internship", "co-op"}
)

// ExtractExperienceLevel infers the experience level from a job title and description.
func ExtractExperienceLevel(title, description string) model.ExperienceLevel {
	descSnippet := description
	if len(descSnippet) > 500 {
		descSnippet = descSnippet[:500]
	}
	lower := strings.ToLower(title + " " + descSnippet)

	for _, kw := range internKeywords {
		if strings.Contains(lower, kw) {
			return model.LevelInternship
		}
	}
	for _, kw := range execKeywords {
		if strings.Contains(lower, kw) {
			return model.LevelExecutive
		}
	}
	for _, kw := range seniorKeywords {
		if strings.Contains(lower, kw) {
			return model.LevelSenior
		}
	}
	for _, kw := range midKeywords {
		if strings.Contains(lower, kw) {
			return model.LevelMid
		}
	}
	for _, kw := range entryKeywords {
		if strings.Contains(lower, kw) {
			return model.LevelEntry
		}
	}
	return model.LevelUnknown
}

// ExtractLocationType infers the work location type from location and description.
func ExtractLocationType(location, description string) model.WorkLocationType {
	descSnippet := description
	if len(descSnippet) > 300 {
		descSnippet = descSnippet[:300]
	}
	lower := strings.ToLower(location + " " + descSnippet)

	for _, kw := range hybridKeywords {
		if strings.Contains(lower, kw) {
			return model.LocationHybrid
		}
	}
	for _, kw := range remoteKeywords {
		if strings.Contains(lower, kw) {
			return model.LocationRemote
		}
	}
	if location != "" {
		return model.LocationOnSite
	}
	return model.LocationUnknown
}

// ExtractEmploymentType infers the employment type from title and description.
func ExtractEmploymentType(title, description string) model.EmploymentType {
	descSnippet := description
	if len(descSnippet) > 300 {
		descSnippet = descSnippet[:300]
	}
	lower := strings.ToLower(title + " " + descSnippet)

	for _, kw := range internshipKeywords {
		if strings.Contains(lower, kw) {
			return model.EmploymentInternship
		}
	}
	for _, kw := range contractKeywords {
		if strings.Contains(lower, kw) {
			return model.EmploymentContract
		}
	}
	for _, kw := range partTimeKeywords {
		if strings.Contains(lower, kw) {
			return model.EmploymentPartTime
		}
	}
	return model.EmploymentFullTime
}

// ParseSalary extracts salary range from a raw salary string.
// Returns (min, max, currency, raw).
func ParseSalary(raw string) (minVal, maxVal *int, currency string) {
	if raw == "" {
		return nil, nil, "USD"
	}

	// Detect currency
	currency = "USD"
	if strings.Contains(raw, "£") || strings.Contains(raw, "GBP") {
		currency = "GBP"
	} else if strings.Contains(raw, "€") || strings.Contains(raw, "EUR") {
		currency = "EUR"
	} else if strings.Contains(raw, "CAD") {
		currency = "CAD"
	}

	// Try range pattern first
	if m := salaryRangeRe.FindStringSubmatch(raw); len(m) >= 3 {
		mn := parseSalaryValue(m[1], raw)
		mx := parseSalaryValue(m[2], raw)
		return &mn, &mx, currency
	}

	// Try single value
	if m := salarySingleRe.FindStringSubmatch(raw); len(m) >= 2 {
		val := parseSalaryValue(m[1], raw)
		return &val, nil, currency
	}

	return nil, nil, currency
}

// parseSalaryValue converts a salary string to an integer annual value.
func parseSalaryValue(s, context string) int {
	s = strings.ReplaceAll(s, ",", "")
	var val float64
	fmt.Sscanf(s, "%f", &val)

	// Check for 'k' suffix in context
	if strings.Contains(strings.ToLower(context), "k") && val < 1000 {
		val *= 1000
	}

	// If value looks like hourly (< 200), convert to annual
	if val > 0 && val < 200 {
		val *= 2080 // 40 hours/week * 52 weeks
	}

	return int(val)
}

// ExtractSkillsFromText extracts skill keywords from job description text.
func ExtractSkillsFromText(text string) (required, preferred []string) {
	techSkills := []string{
		"go", "golang", "python", "java", "javascript", "typescript", "rust",
		"c++", "c#", "ruby", "php", "swift", "kotlin", "scala",
		"react", "angular", "vue", "node.js", "django", "flask", "spring",
		"postgresql", "mysql", "mongodb", "redis", "elasticsearch",
		"aws", "azure", "gcp", "docker", "kubernetes", "terraform",
		"git", "ci/cd", "agile", "scrum", "rest", "graphql", "grpc",
		"machine learning", "deep learning", "nlp", "data science",
		"sql", "nosql", "microservices", "kafka", "rabbitmq",
	}

	lower := strings.ToLower(text)
	seen := map[string]bool{}

	requiredSection := extractSection(text, []string{"required", "must have", "requirements"})
	preferredSection := extractSection(text, []string{"preferred", "nice to have", "bonus", "plus"})

	for _, skill := range techSkills {
		if seen[skill] {
			continue
		}
		if strings.Contains(lower, skill) {
			seen[skill] = true
			if requiredSection != "" && strings.Contains(strings.ToLower(requiredSection), skill) {
				required = append(required, skill)
			} else if preferredSection != "" && strings.Contains(strings.ToLower(preferredSection), skill) {
				preferred = append(preferred, skill)
			} else {
				required = append(required, skill)
			}
		}
	}

	return required, preferred
}

// extractSection finds text after a section header keyword.
func extractSection(text string, headers []string) string {
	lower := strings.ToLower(text)
	for _, h := range headers {
		idx := strings.Index(lower, h)
		if idx >= 0 {
			end := idx + 500
			if end > len(text) {
				end = len(text)
			}
			return text[idx:end]
		}
	}
	return ""
}

// CleanText removes HTML tags and normalizes whitespace.
func CleanText(s string) string {
	htmlTagRe := regexp.MustCompile(`<[^>]+>`)
	s = htmlTagRe.ReplaceAllString(s, " ")

	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")

	var sb strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				sb.WriteRune(' ')
				prevSpace = true
			}
		} else {
			prevSpace = false
			sb.WriteRune(r)
		}
	}
	return strings.TrimSpace(sb.String())
}

// ParseRelativeDate converts relative date strings to absolute times.
func ParseRelativeDate(s string) *time.Time {
	s = strings.ToLower(strings.TrimSpace(s))
	now := time.Now()

	switch {
	case s == "today" || s == "just now" || strings.Contains(s, "hour"):
		t := now
		return &t
	case strings.Contains(s, "yesterday") || s == "1 day ago":
		t := now.AddDate(0, 0, -1)
		return &t
	case strings.Contains(s, "day"):
		var days int
		fmt.Sscanf(s, "%d", &days)
		if days > 0 {
			t := now.AddDate(0, 0, -days)
			return &t
		}
	case strings.Contains(s, "week"):
		var weeks int
		fmt.Sscanf(s, "%d", &weeks)
		if weeks == 0 {
			weeks = 1
		}
		t := now.AddDate(0, 0, -weeks*7)
		return &t
	case strings.Contains(s, "month"):
		var months int
		fmt.Sscanf(s, "%d", &months)
		if months == 0 {
			months = 1
		}
		t := now.AddDate(0, -months, 0)
		return &t
	}
	return nil
}
