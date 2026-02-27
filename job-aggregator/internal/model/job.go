// Package model defines the core data types for the job aggregation system.
package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// JobSource identifies the origin of a job posting.
type JobSource string

const (
	SourceLinkedIn          JobSource = "linkedin"
	SourceIndeed            JobSource = "indeed"
	SourceCompanyCareerPage JobSource = "company_career_page"
	SourceGlassdoor         JobSource = "glassdoor"
	SourceOther             JobSource = "other"
)

// JobStatus represents the current state of a job posting.
type JobStatus string

const (
	StatusActive  JobStatus = "active"
	StatusExpired JobStatus = "expired"
	StatusFilled  JobStatus = "filled"
	StatusUnknown JobStatus = "unknown"
)

// ExperienceLevel represents the required experience level.
type ExperienceLevel string

const (
	LevelInternship ExperienceLevel = "internship"
	LevelEntry      ExperienceLevel = "entry"
	LevelMid        ExperienceLevel = "mid"
	LevelSenior     ExperienceLevel = "senior"
	LevelLead       ExperienceLevel = "lead"
	LevelExecutive  ExperienceLevel = "executive"
	LevelUnknown    ExperienceLevel = "unknown"
)

// EmploymentType represents the type of employment.
type EmploymentType string

const (
	EmploymentFullTime   EmploymentType = "full_time"
	EmploymentPartTime   EmploymentType = "part_time"
	EmploymentContract   EmploymentType = "contract"
	EmploymentTemporary  EmploymentType = "temporary"
	EmploymentInternship EmploymentType = "internship"
	EmploymentVolunteer  EmploymentType = "volunteer"
	EmploymentOther      EmploymentType = "other"
)

// WorkLocationType represents the work location arrangement.
type WorkLocationType string

const (
	LocationOnSite  WorkLocationType = "on_site"
	LocationRemote  WorkLocationType = "remote"
	LocationHybrid  WorkLocationType = "hybrid"
	LocationUnknown WorkLocationType = "unknown"
)

// ScrapeStatus represents the status of a scraping run.
type ScrapeStatus string

const (
	ScrapeStatusPending     ScrapeStatus = "pending"
	ScrapeStatusRunning     ScrapeStatus = "running"
	ScrapeStatusCompleted   ScrapeStatus = "completed"
	ScrapeStatusFailed      ScrapeStatus = "failed"
	ScrapeStatusRateLimited ScrapeStatus = "rate_limited"
)

// Company represents a deduplicated company record.
type Company struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	Name           string         `db:"name" json:"name"`
	NormalizedName string         `db:"normalized_name" json:"normalized_name"`
	Industry       sql.NullString `db:"industry" json:"industry,omitempty"`
	WebsiteURL     sql.NullString `db:"website_url" json:"website_url,omitempty"`
	LinkedInURL    sql.NullString `db:"linkedin_url" json:"linkedin_url,omitempty"`
	LogoURL        sql.NullString `db:"logo_url" json:"logo_url,omitempty"`
	SizeRange      sql.NullString `db:"size_range" json:"size_range,omitempty"`
	Headquarters   sql.NullString `db:"headquarters" json:"headquarters,omitempty"`
	Description    sql.NullString `db:"description" json:"description,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// Job represents a single job posting.
type Job struct {
	ID              uuid.UUID        `db:"id" json:"id"`
	DedupHash       string           `db:"dedup_hash" json:"dedup_hash"`
	Source          JobSource        `db:"source" json:"source"`
	ExternalID      sql.NullString   `db:"external_id" json:"external_id,omitempty"`
	CompanyID       *uuid.UUID       `db:"company_id" json:"company_id,omitempty"`
	CompanyName     string           `db:"company_name" json:"company_name"`
	Title           string           `db:"title" json:"title"`
	Description     sql.NullString   `db:"description" json:"description,omitempty"`
	DescriptionHTML sql.NullString   `db:"description_html" json:"description_html,omitempty"`
	Industry        sql.NullString   `db:"industry" json:"industry,omitempty"`
	LocationCity    sql.NullString   `db:"location_city" json:"location_city,omitempty"`
	LocationState   sql.NullString   `db:"location_state" json:"location_state,omitempty"`
	LocationCountry sql.NullString   `db:"location_country" json:"location_country,omitempty"`
	LocationRaw     sql.NullString   `db:"location_raw" json:"location_raw,omitempty"`
	LocationType    WorkLocationType `db:"location_type" json:"location_type"`
	EmploymentType  EmploymentType   `db:"employment_type" json:"employment_type"`
	ExperienceLevel ExperienceLevel  `db:"experience_level" json:"experience_level"`
	RequiredSkills  pq.StringArray   `db:"required_skills" json:"required_skills,omitempty"`
	PreferredSkills pq.StringArray   `db:"preferred_skills" json:"preferred_skills,omitempty"`
	SalaryMin       sql.NullInt32    `db:"salary_min" json:"salary_min,omitempty"`
	SalaryMax       sql.NullInt32    `db:"salary_max" json:"salary_max,omitempty"`
	SalaryCurrency  string           `db:"salary_currency" json:"salary_currency"`
	SalaryRaw       sql.NullString   `db:"salary_raw" json:"salary_raw,omitempty"`
	ApplicationURL  string           `db:"application_url" json:"application_url"`
	CompanyURL      sql.NullString   `db:"company_url" json:"company_url,omitempty"`
	PostedAt        sql.NullTime     `db:"posted_at" json:"posted_at,omitempty"`
	ExpiresAt       sql.NullTime     `db:"expires_at" json:"expires_at,omitempty"`
	ScrapedAt       time.Time        `db:"scraped_at" json:"scraped_at"`
	LastSeenAt      time.Time        `db:"last_seen_at" json:"last_seen_at"`
	Status          JobStatus        `db:"status" json:"status"`
	IsFeatured      bool             `db:"is_featured" json:"is_featured"`
	RawData         []byte           `db:"raw_data" json:"raw_data,omitempty"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at" json:"updated_at"`
}

// ScrapeRun represents a single scraping run log entry.
type ScrapeRun struct {
	ID             uuid.UUID    `db:"id" json:"id"`
	Source         JobSource    `db:"source" json:"source"`
	SearchQuery    string       `db:"search_query" json:"search_query"`
	SearchLocation string       `db:"search_location" json:"search_location"`
	Status         ScrapeStatus `db:"status" json:"status"`
	JobsFound      int          `db:"jobs_found" json:"jobs_found"`
	JobsNew        int          `db:"jobs_new" json:"jobs_new"`
	JobsUpdated    int          `db:"jobs_updated" json:"jobs_updated"`
	JobsFailed     int          `db:"jobs_failed" json:"jobs_failed"`
	PagesScraped   int          `db:"pages_scraped" json:"pages_scraped"`
	ErrorMessage   string       `db:"error_message" json:"error_message,omitempty"`
	StartedAt      sql.NullTime `db:"started_at" json:"started_at,omitempty"`
	CompletedAt    sql.NullTime `db:"completed_at" json:"completed_at,omitempty"`
	DurationMs     sql.NullInt32 `db:"duration_ms" json:"duration_ms,omitempty"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
}

// ScrapeConfig holds configuration for a scraper source.
type ScrapeConfig struct {
	ID                 uuid.UUID      `db:"id" json:"id"`
	Source             JobSource      `db:"source" json:"source"`
	Name               string         `db:"name" json:"name"`
	IsEnabled          bool           `db:"is_enabled" json:"is_enabled"`
	RequestsPerMinute  int            `db:"requests_per_minute" json:"requests_per_minute"`
	MaxRetries         int            `db:"max_retries" json:"max_retries"`
	RetryDelayMs       int            `db:"retry_delay_ms" json:"retry_delay_ms"`
	CronSchedule       string         `db:"cron_schedule" json:"cron_schedule"`
	LastRunAt          sql.NullTime   `db:"last_run_at" json:"last_run_at,omitempty"`
	NextRunAt          sql.NullTime   `db:"next_run_at" json:"next_run_at,omitempty"`
	Config             []byte         `db:"config" json:"config"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at" json:"updated_at"`
}

// CompanyCareerPage holds configuration for a company career page scraper.
type CompanyCareerPage struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	CompanyID     *uuid.UUID `db:"company_id" json:"company_id,omitempty"`
	CompanyName   string     `db:"company_name" json:"company_name"`
	CareerPageURL string     `db:"career_page_url" json:"career_page_url"`
	Selectors     []byte     `db:"selectors" json:"selectors"`
	IsEnabled     bool       `db:"is_enabled" json:"is_enabled"`
	LastScrapedAt sql.NullTime `db:"last_scraped_at" json:"last_scraped_at,omitempty"`
	JobsFound     int        `db:"jobs_found" json:"jobs_found"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

// ScrapedJob is the intermediate representation returned by scrapers
// before being stored in the database.
type ScrapedJob struct {
	Source          JobSource
	ExternalID      string
	CompanyName     string
	Title           string
	Description     string
	DescriptionHTML string
	Industry        string
	LocationRaw     string
	LocationCity    string
	LocationState   string
	LocationCountry string
	LocationType    WorkLocationType
	EmploymentType  EmploymentType
	ExperienceLevel ExperienceLevel
	RequiredSkills  []string
	PreferredSkills []string
	SalaryMin       *int
	SalaryMax       *int
	SalaryCurrency  string
	SalaryRaw       string
	ApplicationURL  string
	CompanyURL      string
	PostedAt        *time.Time
	ExpiresAt       *time.Time
	RawData         map[string]interface{}
}

// SearchParams holds parameters for job search queries.
type SearchParams struct {
	Query          string
	Location       string
	Remote         bool
	ExperienceLevel ExperienceLevel
	EmploymentType  EmploymentType
	SalaryMin      *int
	SalaryMax      *int
	Skills         []string
	Page           int
	PageSize       int
}

// JobFilter holds filter criteria for querying stored jobs.
type JobFilter struct {
	Sources         []JobSource
	Status          JobStatus
	LocationTypes   []WorkLocationType
	ExperienceLevels []ExperienceLevel
	EmploymentTypes []EmploymentType
	Skills          []string
	SalaryMin       *int
	SalaryMax       *int
	PostedAfter     *time.Time
	CompanyName     string
	TitleSearch     string
	Page            int
	PageSize        int
}

// AdminStats holds aggregated statistics for the admin dashboard.
type AdminStats struct {
	TotalJobs       int                    `json:"total_jobs"`
	ActiveJobs      int                    `json:"active_jobs"`
	JobsBySource    map[string]int         `json:"jobs_by_source"`
	RecentRuns      []ScrapeRun            `json:"recent_runs"`
	SourceStats     []SourceStat           `json:"source_stats"`
}

// SourceStat holds per-source statistics.
type SourceStat struct {
	Source          string  `json:"source"`
	TotalRuns       int     `json:"total_runs"`
	SuccessfulRuns  int     `json:"successful_runs"`
	FailedRuns      int     `json:"failed_runs"`
	TotalJobsFound  int     `json:"total_jobs_found"`
	TotalJobsNew    int     `json:"total_jobs_new"`
	AvgDurationMs   float64 `json:"avg_duration_ms"`
	LastSuccessfulRun *time.Time `json:"last_successful_run,omitempty"`
}
