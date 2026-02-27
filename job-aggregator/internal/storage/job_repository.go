// Package storage provides database access for the job aggregation system.
package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnbot/job-aggregator/internal/model"
	"github.com/lib/pq"
)

// JobRepository provides CRUD operations for jobs and related entities.
type JobRepository struct {
	db *sql.DB
}

// NewJobRepository creates a new JobRepository.
func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// Job storage
// ─────────────────────────────────────────────────────────────────────────────

// UpsertJob inserts a new job or updates it if the dedup_hash already exists.
// Returns (job, isNew, error).
func (r *JobRepository) UpsertJob(ctx context.Context, scraped *model.ScrapedJob) (*model.Job, bool, error) {
	hash := ComputeDedupHash(scraped)

	rawData, _ := json.Marshal(scraped.RawData)

	job := &model.Job{}
	var isNew bool

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO jobs (
			dedup_hash, source, external_id, company_name, title,
			description, description_html, industry,
			location_city, location_state, location_country, location_raw,
			location_type, employment_type, experience_level,
			required_skills, preferred_skills,
			salary_min, salary_max, salary_currency, salary_raw,
			application_url, company_url, posted_at, expires_at, raw_data
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15,
			$16, $17,
			$18, $19, $20, $21,
			$22, $23, $24, $25, $26
		)
		ON CONFLICT (dedup_hash) DO UPDATE SET
			last_seen_at     = NOW(),
			status           = 'active',
			description      = EXCLUDED.description,
			description_html = EXCLUDED.description_html,
			required_skills  = EXCLUDED.required_skills,
			preferred_skills = EXCLUDED.preferred_skills,
			salary_min       = EXCLUDED.salary_min,
			salary_max       = EXCLUDED.salary_max,
			salary_raw       = EXCLUDED.salary_raw,
			expires_at       = EXCLUDED.expires_at,
			raw_data         = EXCLUDED.raw_data,
			updated_at       = NOW()
		RETURNING id, dedup_hash, source, external_id, company_name, title,
		          description, location_city, location_state, location_country,
		          location_raw, location_type, employment_type, experience_level,
		          required_skills, preferred_skills, salary_min, salary_max,
		          salary_currency, salary_raw, application_url, company_url,
		          posted_at, expires_at, scraped_at, last_seen_at, status,
		          is_featured, created_at, updated_at,
		          (xmax = 0) AS is_new`,
		hash, scraped.Source, nullString(scraped.ExternalID), scraped.CompanyName, scraped.Title,
		nullString(scraped.Description), nullString(scraped.DescriptionHTML), nullString(scraped.Industry),
		nullString(scraped.LocationCity), nullString(scraped.LocationState),
		nullString(scraped.LocationCountry), nullString(scraped.LocationRaw),
		scraped.LocationType, scraped.EmploymentType, scraped.ExperienceLevel,
		pq.Array(scraped.RequiredSkills), pq.Array(scraped.PreferredSkills),
		scraped.SalaryMin, scraped.SalaryMax, nullString(scraped.SalaryCurrency), nullString(scraped.SalaryRaw),
		scraped.ApplicationURL, nullString(scraped.CompanyURL),
		scraped.PostedAt, scraped.ExpiresAt, rawData,
	).Scan(
		&job.ID, &job.DedupHash, &job.Source, &job.ExternalID,
		&job.CompanyName, &job.Title, &job.Description,
		&job.LocationCity, &job.LocationState, &job.LocationCountry,
		&job.LocationRaw, &job.LocationType, &job.EmploymentType, &job.ExperienceLevel,
		&job.RequiredSkills, &job.PreferredSkills,
		&job.SalaryMin, &job.SalaryMax, &job.SalaryCurrency, &job.SalaryRaw,
		&job.ApplicationURL, &job.CompanyURL,
		&job.PostedAt, &job.ExpiresAt, &job.ScrapedAt, &job.LastSeenAt,
		&job.Status, &job.IsFeatured, &job.CreatedAt, &job.UpdatedAt,
		&isNew,
	)
	if err != nil {
		return nil, false, fmt.Errorf("upsert job: %w", err)
	}

	return job, isNew, nil
}

// GetJobByID retrieves a job by its UUID.
func (r *JobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*model.Job, error) {
	job := &model.Job{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, dedup_hash, source, external_id, company_name, title,
		       description, description_html, industry,
		       location_city, location_state, location_country, location_raw,
		       location_type, employment_type, experience_level,
		       required_skills, preferred_skills,
		       salary_min, salary_max, salary_currency, salary_raw,
		       application_url, company_url, posted_at, expires_at,
		       scraped_at, last_seen_at, status, is_featured, created_at, updated_at
		FROM jobs WHERE id = $1`, id,
	).Scan(
		&job.ID, &job.DedupHash, &job.Source, &job.ExternalID,
		&job.CompanyName, &job.Title, &job.Description, &job.DescriptionHTML,
		&job.Industry, &job.LocationCity, &job.LocationState, &job.LocationCountry,
		&job.LocationRaw, &job.LocationType, &job.EmploymentType, &job.ExperienceLevel,
		&job.RequiredSkills, &job.PreferredSkills,
		&job.SalaryMin, &job.SalaryMax, &job.SalaryCurrency, &job.SalaryRaw,
		&job.ApplicationURL, &job.CompanyURL, &job.PostedAt, &job.ExpiresAt,
		&job.ScrapedAt, &job.LastSeenAt, &job.Status, &job.IsFeatured,
		&job.CreatedAt, &job.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get job by id: %w", err)
	}
	return job, nil
}

// SearchJobs queries jobs with filters and pagination.
func (r *JobRepository) SearchJobs(ctx context.Context, filter model.JobFilter) ([]model.Job, int, error) {
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	where := []string{"1=1"}
	args := []interface{}{}
	idx := 1

	// Status filter (default to active)
	status := filter.Status
	if status == "" {
		status = model.StatusActive
	}
	where = append(where, fmt.Sprintf("status = $%d", idx))
	args = append(args, status)
	idx++

	// Source filter
	if len(filter.Sources) > 0 {
		where = append(where, fmt.Sprintf("source = ANY($%d)", idx))
		args = append(args, pq.Array(filter.Sources))
		idx++
	}

	// Location type filter
	if len(filter.LocationTypes) > 0 {
		where = append(where, fmt.Sprintf("location_type = ANY($%d)", idx))
		args = append(args, pq.Array(filter.LocationTypes))
		idx++
	}

	// Experience level filter
	if len(filter.ExperienceLevels) > 0 {
		where = append(where, fmt.Sprintf("experience_level = ANY($%d)", idx))
		args = append(args, pq.Array(filter.ExperienceLevels))
		idx++
	}

	// Skills filter (job must have at least one of the required skills)
	if len(filter.Skills) > 0 {
		where = append(where, fmt.Sprintf("required_skills && $%d", idx))
		args = append(args, pq.Array(filter.Skills))
		idx++
	}

	// Salary filter
	if filter.SalaryMin != nil {
		where = append(where, fmt.Sprintf("(salary_max IS NULL OR salary_max >= $%d)", idx))
		args = append(args, *filter.SalaryMin)
		idx++
	}

	// Posted after filter
	if filter.PostedAfter != nil {
		where = append(where, fmt.Sprintf("(posted_at IS NULL OR posted_at >= $%d)", idx))
		args = append(args, *filter.PostedAfter)
		idx++
	}

	// Full-text title search
	if filter.TitleSearch != "" {
		where = append(where, fmt.Sprintf(
			"to_tsvector('english', title) @@ plainto_tsquery('english', $%d)", idx))
		args = append(args, filter.TitleSearch)
		idx++
	}

	// Company name filter
	if filter.CompanyName != "" {
		where = append(where, fmt.Sprintf("company_name ILIKE $%d", idx))
		args = append(args, "%"+filter.CompanyName+"%")
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count query
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count jobs: %w", err)
	}

	// Data query with pagination
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)
	dataQuery := fmt.Sprintf(`
		SELECT id, dedup_hash, source, external_id, company_name, title,
		       description, location_city, location_state, location_country,
		       location_raw, location_type, employment_type, experience_level,
		       required_skills, preferred_skills,
		       salary_min, salary_max, salary_currency, salary_raw,
		       application_url, company_url, posted_at, expires_at,
		       scraped_at, last_seen_at, status, is_featured, created_at, updated_at
		FROM jobs
		WHERE %s
		ORDER BY posted_at DESC NULLS LAST, scraped_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, idx, idx+1)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("search jobs: %w", err)
	}
	defer rows.Close()

	var jobs []model.Job
	for rows.Next() {
		var j model.Job
		if err := rows.Scan(
			&j.ID, &j.DedupHash, &j.Source, &j.ExternalID,
			&j.CompanyName, &j.Title, &j.Description,
			&j.LocationCity, &j.LocationState, &j.LocationCountry,
			&j.LocationRaw, &j.LocationType, &j.EmploymentType, &j.ExperienceLevel,
			&j.RequiredSkills, &j.PreferredSkills,
			&j.SalaryMin, &j.SalaryMax, &j.SalaryCurrency, &j.SalaryRaw,
			&j.ApplicationURL, &j.CompanyURL, &j.PostedAt, &j.ExpiresAt,
			&j.ScrapedAt, &j.LastSeenAt, &j.Status, &j.IsFeatured,
			&j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, j)
	}
	return jobs, total, rows.Err()
}

// MarkExpiredJobs marks jobs not seen since the cutoff as expired.
func (r *JobRepository) MarkExpiredJobs(ctx context.Context, source model.JobSource, cutoff time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE jobs SET status = 'expired', updated_at = NOW()
		WHERE source = $1 AND last_seen_at < $2 AND status = 'active'`,
		source, cutoff,
	)
	if err != nil {
		return 0, fmt.Errorf("mark expired jobs: %w", err)
	}
	return result.RowsAffected()
}

// ─────────────────────────────────────────────────────────────────────────────
// Scrape run logging
// ─────────────────────────────────────────────────────────────────────────────

// CreateScrapeRun creates a new scrape run log entry.
func (r *JobRepository) CreateScrapeRun(ctx context.Context, source model.JobSource, query, location string) (*model.ScrapeRun, error) {
	run := &model.ScrapeRun{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO scrape_runs (source, search_query, search_location, status, started_at)
		VALUES ($1, $2, $3, 'running', NOW())
		RETURNING id, source, search_query, search_location, status,
		          jobs_found, jobs_new, jobs_updated, jobs_failed, pages_scraped,
		          error_message, started_at, completed_at, duration_ms, created_at`,
		source, query, location,
	).Scan(
		&run.ID, &run.Source, &run.SearchQuery, &run.SearchLocation, &run.Status,
		&run.JobsFound, &run.JobsNew, &run.JobsUpdated, &run.JobsFailed, &run.PagesScraped,
		&run.ErrorMessage, &run.StartedAt, &run.CompletedAt, &run.DurationMs, &run.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create scrape run: %w", err)
	}
	return run, nil
}

// UpdateScrapeRun updates a scrape run with final statistics.
func (r *JobRepository) UpdateScrapeRun(ctx context.Context, runID uuid.UUID, status model.ScrapeStatus, stats model.ScrapeRun) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE scrape_runs SET
			status        = $2,
			jobs_found    = $3,
			jobs_new      = $4,
			jobs_updated  = $5,
			jobs_failed   = $6,
			pages_scraped = $7,
			error_message = $8,
			completed_at  = NOW()
		WHERE id = $1`,
		runID, status,
		stats.JobsFound, stats.JobsNew, stats.JobsUpdated,
		stats.JobsFailed, stats.PagesScraped, stats.ErrorMessage,
	)
	return err
}

// GetRecentScrapeRuns returns the most recent scrape runs.
func (r *JobRepository) GetRecentScrapeRuns(ctx context.Context, limit int) ([]model.ScrapeRun, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, source, search_query, search_location, status,
		       jobs_found, jobs_new, jobs_updated, jobs_failed, pages_scraped,
		       error_message, started_at, completed_at, duration_ms, created_at
		FROM scrape_runs
		ORDER BY created_at DESC
		LIMIT $1`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get recent scrape runs: %w", err)
	}
	defer rows.Close()

	var runs []model.ScrapeRun
	for rows.Next() {
		var r model.ScrapeRun
		if err := rows.Scan(
			&r.ID, &r.Source, &r.SearchQuery, &r.SearchLocation, &r.Status,
			&r.JobsFound, &r.JobsNew, &r.JobsUpdated, &r.JobsFailed, &r.PagesScraped,
			&r.ErrorMessage, &r.StartedAt, &r.CompletedAt, &r.DurationMs, &r.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan scrape run: %w", err)
		}
		runs = append(runs, r)
	}
	return runs, rows.Err()
}

// GetAdminStats returns aggregated statistics for the admin dashboard.
func (r *JobRepository) GetAdminStats(ctx context.Context) (*model.AdminStats, error) {
	stats := &model.AdminStats{
		JobsBySource: make(map[string]int),
	}

	// Total and active job counts
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs`).Scan(&stats.TotalJobs)
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE status = 'active'`).Scan(&stats.ActiveJobs)

	// Jobs by source
	rows, err := r.db.QueryContext(ctx, `
		SELECT source, COUNT(*) FROM jobs WHERE status = 'active' GROUP BY source`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var src string
			var count int
			rows.Scan(&src, &count)
			stats.JobsBySource[src] = count
		}
	}

	// Recent runs
	stats.RecentRuns, _ = r.GetRecentScrapeRuns(ctx, 10)

	// Source stats from view
	srcRows, err := r.db.QueryContext(ctx, `
		SELECT source, total_runs, successful_runs, failed_runs,
		       total_jobs_found, total_jobs_new, avg_duration_ms, last_successful_run
		FROM v_scrape_stats ORDER BY source`)
	if err == nil {
		defer srcRows.Close()
		for srcRows.Next() {
			var s model.SourceStat
			srcRows.Scan(
				&s.Source, &s.TotalRuns, &s.SuccessfulRuns, &s.FailedRuns,
				&s.TotalJobsFound, &s.TotalJobsNew, &s.AvgDurationMs, &s.LastSuccessfulRun,
			)
			stats.SourceStats = append(stats.SourceStats, s)
		}
	}

	return stats, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Company management
// ─────────────────────────────────────────────────────────────────────────────

// UpsertCompany creates or updates a company record.
func (r *JobRepository) UpsertCompany(ctx context.Context, name string) (*model.Company, error) {
	normalized := normalizeCompanyName(name)
	company := &model.Company{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO companies (name, normalized_name)
		VALUES ($1, $2)
		ON CONFLICT (normalized_name) DO UPDATE SET
			name = EXCLUDED.name,
			updated_at = NOW()
		RETURNING id, name, normalized_name, industry, website_url, linkedin_url,
		          logo_url, size_range, headquarters, description, created_at, updated_at`,
		name, normalized,
	).Scan(
		&company.ID, &company.Name, &company.NormalizedName,
		&company.Industry, &company.WebsiteURL, &company.LinkedInURL,
		&company.LogoURL, &company.SizeRange, &company.Headquarters,
		&company.Description, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert company: %w", err)
	}
	return company, nil
}

// GetCareerPages returns all enabled company career pages.
func (r *JobRepository) GetCareerPages(ctx context.Context) ([]model.CompanyCareerPage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, company_id, company_name, career_page_url, selectors,
		       is_enabled, last_scraped_at, jobs_found, created_at, updated_at
		FROM company_career_pages
		WHERE is_enabled = TRUE
		ORDER BY company_name`)
	if err != nil {
		return nil, fmt.Errorf("get career pages: %w", err)
	}
	defer rows.Close()

	var pages []model.CompanyCareerPage
	for rows.Next() {
		var p model.CompanyCareerPage
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.CompanyName, &p.CareerPageURL,
			&p.Selectors, &p.IsEnabled, &p.LastScrapedAt,
			&p.JobsFound, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan career page: %w", err)
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// ComputeDedupHash generates a deduplication hash for a scraped job.
func ComputeDedupHash(job *model.ScrapedJob) string {
	// Use external ID if available (most reliable)
	if job.ExternalID != "" {
		key := fmt.Sprintf("%s:%s", job.Source, job.ExternalID)
		h := sha256.Sum256([]byte(key))
		return fmt.Sprintf("%x", h[:16])
	}
	// Fall back to title + company + location hash
	key := fmt.Sprintf("%s:%s:%s:%s",
		strings.ToLower(string(job.Source)),
		strings.ToLower(strings.TrimSpace(job.Title)),
		strings.ToLower(strings.TrimSpace(job.CompanyName)),
		strings.ToLower(strings.TrimSpace(job.LocationRaw)),
	)
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h[:16])
}

// normalizeCompanyName returns a lowercase, trimmed company name for deduplication.
func normalizeCompanyName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	// Remove common suffixes
	suffixes := []string{", inc.", " inc.", ", llc", " llc", ", ltd.", " ltd.", ", corp.", " corp."}
	for _, s := range suffixes {
		name = strings.TrimSuffix(name, s)
	}
	return strings.TrimSpace(name)
}

// nullString returns a sql.NullString from a string (empty = null).
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// ErrNotFound is returned when a record is not found.
var ErrNotFound = fmt.Errorf("record not found")
