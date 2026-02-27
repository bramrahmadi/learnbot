-- Migration 001: Create job aggregation schema

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- Enum types
-- ─────────────────────────────────────────────────────────────────────────────

CREATE TYPE job_source AS ENUM (
    'linkedin',
    'indeed',
    'company_career_page',
    'glassdoor',
    'other'
);

CREATE TYPE job_status AS ENUM (
    'active',       -- currently open
    'expired',      -- past closing date
    'filled',       -- position filled
    'unknown'       -- status not determinable
);

CREATE TYPE experience_level AS ENUM (
    'internship',
    'entry',        -- 0-2 years
    'mid',          -- 2-5 years
    'senior',       -- 5-10 years
    'lead',         -- 8+ years, team lead
    'executive',    -- director, VP, C-suite
    'unknown'
);

CREATE TYPE employment_type AS ENUM (
    'full_time',
    'part_time',
    'contract',
    'temporary',
    'internship',
    'volunteer',
    'other'
);

CREATE TYPE work_location_type AS ENUM (
    'on_site',
    'remote',
    'hybrid',
    'unknown'
);

CREATE TYPE scrape_status AS ENUM (
    'pending',
    'running',
    'completed',
    'failed',
    'rate_limited'
);

-- ─────────────────────────────────────────────────────────────────────────────
-- companies: Deduplicated company records
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE companies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    industry        TEXT,
    website_url     TEXT,
    linkedin_url    TEXT,
    logo_url        TEXT,
    size_range      TEXT,                          -- '1-10', '11-50', '51-200', etc.
    headquarters    TEXT,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT companies_normalized_unique UNIQUE (normalized_name)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- jobs: Core job postings table
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE jobs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Deduplication key: hash of (source, external_id) or (title, company, location)
    dedup_hash          TEXT NOT NULL,
    source              job_source NOT NULL,
    external_id         TEXT,                      -- source's own job ID
    company_id          UUID REFERENCES companies(id),
    company_name        TEXT NOT NULL,             -- denormalized for fast reads
    title               TEXT NOT NULL,
    description         TEXT,
    description_html    TEXT,                      -- raw HTML if available
    industry            TEXT,
    location_city       TEXT,
    location_state      TEXT,
    location_country    TEXT,
    location_raw        TEXT,                      -- original location string
    location_type       work_location_type NOT NULL DEFAULT 'unknown',
    employment_type     employment_type NOT NULL DEFAULT 'full_time',
    experience_level    experience_level NOT NULL DEFAULT 'unknown',
    -- Skills extracted from description
    required_skills     TEXT[],
    preferred_skills    TEXT[],
    -- Salary
    salary_min          INTEGER,                   -- annual, in USD
    salary_max          INTEGER,
    salary_currency     TEXT DEFAULT 'USD',
    salary_raw          TEXT,                      -- original salary string
    -- URLs
    application_url     TEXT NOT NULL,
    company_url         TEXT,
    -- Dates
    posted_at           TIMESTAMPTZ,
    expires_at          TIMESTAMPTZ,
    scraped_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Status
    status              job_status NOT NULL DEFAULT 'active',
    is_featured         BOOLEAN NOT NULL DEFAULT FALSE,
    -- Metadata
    raw_data            JSONB,                     -- full scraped payload
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT jobs_dedup_unique UNIQUE (dedup_hash),
    CONSTRAINT jobs_salary_range CHECK (
        salary_min IS NULL OR salary_max IS NULL OR salary_min <= salary_max
    ),
    CONSTRAINT jobs_salary_positive CHECK (
        (salary_min IS NULL OR salary_min >= 0) AND
        (salary_max IS NULL OR salary_max >= 0)
    ),
    CONSTRAINT jobs_title_not_empty CHECK (LENGTH(TRIM(title)) > 0),
    CONSTRAINT jobs_company_not_empty CHECK (LENGTH(TRIM(company_name)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- scrape_runs: Log of each scraping run
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE scrape_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source          job_source NOT NULL,
    search_query    TEXT,                          -- search terms used
    search_location TEXT,                          -- location filter
    status          scrape_status NOT NULL DEFAULT 'pending',
    jobs_found      INTEGER NOT NULL DEFAULT 0,
    jobs_new        INTEGER NOT NULL DEFAULT 0,
    jobs_updated    INTEGER NOT NULL DEFAULT 0,
    jobs_failed     INTEGER NOT NULL DEFAULT 0,
    pages_scraped   INTEGER NOT NULL DEFAULT 0,
    error_message   TEXT,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    duration_ms     INTEGER GENERATED ALWAYS AS (
        CASE
            WHEN completed_at IS NOT NULL AND started_at IS NOT NULL
            THEN EXTRACT(EPOCH FROM (completed_at - started_at))::INTEGER * 1000
            ELSE NULL
        END
    ) STORED,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ─────────────────────────────────────────────────────────────────────────────
-- scrape_configs: Configuration for each scraper source
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE scrape_configs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source              job_source NOT NULL,
    name                TEXT NOT NULL,
    is_enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    -- Rate limiting
    requests_per_minute INTEGER NOT NULL DEFAULT 10,
    max_retries         INTEGER NOT NULL DEFAULT 3,
    retry_delay_ms      INTEGER NOT NULL DEFAULT 1000,
    -- Scheduling
    cron_schedule       TEXT NOT NULL DEFAULT '0 2 * * *',  -- daily at 2am
    last_run_at         TIMESTAMPTZ,
    next_run_at         TIMESTAMPTZ,
    -- Source-specific config (API keys, URLs, etc.)
    config              JSONB NOT NULL DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT scrape_configs_source_unique UNIQUE (source, name)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- company_career_pages: Configurable list of company career pages to scrape
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE company_career_pages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID REFERENCES companies(id),
    company_name    TEXT NOT NULL,
    career_page_url TEXT NOT NULL,
    -- CSS selectors for extracting job data
    selectors       JSONB NOT NULL DEFAULT '{}',
    is_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    last_scraped_at TIMESTAMPTZ,
    jobs_found      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT career_pages_url_unique UNIQUE (career_page_url)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- Indexes (optimized for job search and matching)
-- ─────────────────────────────────────────────────────────────────────────────

-- companies
CREATE INDEX idx_companies_normalized ON companies(normalized_name);
CREATE INDEX idx_companies_industry ON companies(industry);

-- jobs - primary search patterns
CREATE INDEX idx_jobs_status ON jobs(status) WHERE status = 'active';
CREATE INDEX idx_jobs_source ON jobs(source);
CREATE INDEX idx_jobs_company_id ON jobs(company_id);
CREATE INDEX idx_jobs_company_name ON jobs(company_name);
CREATE INDEX idx_jobs_location ON jobs(location_country, location_state, location_city);
CREATE INDEX idx_jobs_location_type ON jobs(location_type);
CREATE INDEX idx_jobs_employment_type ON jobs(employment_type);
CREATE INDEX idx_jobs_experience_level ON jobs(experience_level);
CREATE INDEX idx_jobs_posted_at ON jobs(posted_at DESC NULLS LAST);
CREATE INDEX idx_jobs_scraped_at ON jobs(scraped_at DESC);
CREATE INDEX idx_jobs_last_seen ON jobs(last_seen_at DESC);
CREATE INDEX idx_jobs_salary ON jobs(salary_min, salary_max) WHERE salary_min IS NOT NULL;
-- Full-text search on title and description
CREATE INDEX idx_jobs_title_fts ON jobs USING GIN(to_tsvector('english', title));
CREATE INDEX idx_jobs_description_fts ON jobs USING GIN(to_tsvector('english', COALESCE(description, '')));
-- Skills arrays
CREATE INDEX idx_jobs_required_skills ON jobs USING GIN(required_skills);
CREATE INDEX idx_jobs_preferred_skills ON jobs USING GIN(preferred_skills);
-- Composite for common filter combinations
CREATE INDEX idx_jobs_active_location ON jobs(status, location_type, experience_level)
    WHERE status = 'active';

-- scrape_runs
CREATE INDEX idx_scrape_runs_source ON scrape_runs(source);
CREATE INDEX idx_scrape_runs_status ON scrape_runs(status);
CREATE INDEX idx_scrape_runs_created_at ON scrape_runs(created_at DESC);

-- ─────────────────────────────────────────────────────────────────────────────
-- Triggers
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER jobs_updated_at
    BEFORE UPDATE ON jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER scrape_configs_updated_at
    BEFORE UPDATE ON scrape_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER career_pages_updated_at
    BEFORE UPDATE ON company_career_pages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ─────────────────────────────────────────────────────────────────────────────
-- Views
-- ─────────────────────────────────────────────────────────────────────────────

-- Active jobs with company info
CREATE VIEW v_active_jobs AS
SELECT
    j.id,
    j.source,
    j.title,
    j.company_name,
    c.industry,
    c.size_range AS company_size,
    j.location_city,
    j.location_state,
    j.location_country,
    j.location_type,
    j.employment_type,
    j.experience_level,
    j.required_skills,
    j.preferred_skills,
    j.salary_min,
    j.salary_max,
    j.salary_currency,
    j.application_url,
    j.posted_at,
    j.scraped_at
FROM jobs j
LEFT JOIN companies c ON c.id = j.company_id
WHERE j.status = 'active'
ORDER BY j.posted_at DESC NULLS LAST;

-- Scraping statistics per source
CREATE VIEW v_scrape_stats AS
SELECT
    source,
    COUNT(*) AS total_runs,
    COUNT(*) FILTER (WHERE status = 'completed') AS successful_runs,
    COUNT(*) FILTER (WHERE status = 'failed') AS failed_runs,
    SUM(jobs_found) AS total_jobs_found,
    SUM(jobs_new) AS total_jobs_new,
    AVG(duration_ms) FILTER (WHERE duration_ms IS NOT NULL) AS avg_duration_ms,
    MAX(completed_at) AS last_successful_run
FROM scrape_runs
GROUP BY source;

-- Job count by source and status
CREATE VIEW v_job_counts AS
SELECT
    source,
    status,
    COUNT(*) AS count,
    MIN(posted_at) AS oldest_job,
    MAX(posted_at) AS newest_job
FROM jobs
GROUP BY source, status;

COMMIT;
