# Job Aggregator

A high-performance Go job aggregation service that collects job postings from multiple sources (LinkedIn, Indeed, and company career pages) with rate limiting, retry logic, deduplication, and a daily scheduler. Part of the [LearnBot AI Career Development Platform](../LearnBot.md).

## Features

- **Multi-source scraping**: LinkedIn Jobs, Indeed, and configurable company career pages
- **Rate limiting**: Per-source configurable requests/minute with `golang.org/x/time/rate`
- **Retry logic**: Exponential backoff with configurable max retries
- **Deduplication**: SHA-256 hash-based deduplication prevents duplicate job entries
- **Concurrent processing**: Worker pool for parallel job storage
- **Daily scheduler**: Automatic daily scraping at 2am UTC
- **robots.txt compliance**: Checks robots.txt before scraping any URL
- **Admin dashboard**: HTTP endpoints for monitoring scraping status
- **Structured logging**: Per-scraper logging with timestamps

## Architecture

```
job-aggregator/
├── cmd/server/          # HTTP server entry point
├── internal/
│   ├── model/           # Data types (Job, Company, ScrapeRun, etc.)
│   ├── storage/         # PostgreSQL repository with deduplication
│   ├── httpclient/      # Rate-limited HTTP client with retry + robots.txt
│   ├── scraper/
│   │   ├── scraper.go       # Scraper interface + text extraction utilities
│   │   ├── linkedin.go      # LinkedIn Jobs scraper
│   │   ├── indeed.go        # Indeed scraper
│   │   └── career_page.go   # Configurable company career page scraper
│   ├── scheduler/       # Concurrent worker pool + daily schedule
│   └── admin/           # Admin dashboard HTTP handlers
└── migrations/
    └── 001_create_jobs_schema.sql
```

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 15+

### Setup

```bash
# Run database migration
psql -d learnbot -f migrations/001_create_jobs_schema.sql

# Build and run
cd job-aggregator
go build -o job-aggregator ./cmd/server
./job-aggregator -addr :8081 -db "postgres://localhost/learnbot?sslmode=disable"

# Run scrapers immediately on startup
./job-aggregator --run-now
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://localhost/learnbot?sslmode=disable` | PostgreSQL connection URL |

## Admin Dashboard API

### `GET /admin/health`
Health check endpoint.

```json
{
  "status": "ok",
  "version": "1.0.0",
  "is_running": false,
  "time": "2024-01-15T10:30:00Z"
}
```

### `GET /admin/stats`
Aggregated scraping statistics.

```json
{
  "stats": {
    "total_jobs": 15420,
    "active_jobs": 12350,
    "jobs_by_source": {
      "linkedin": 8200,
      "indeed": 3500,
      "company_career_page": 650
    },
    "source_stats": [...]
  },
  "is_running": false,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### `GET /admin/runs?limit=20`
Recent scraping runs.

### `POST /admin/scrape/trigger`
Trigger an immediate scraping run.

```json
{ "message": "scraping run triggered", "running": true }
```

### `GET /admin/jobs?q=engineer&location_type=remote&page=1&page_size=20`
Search jobs with filters.

**Query parameters:**
| Parameter | Description |
|-----------|-------------|
| `q` | Full-text search on job title |
| `company` | Filter by company name |
| `location_type` | `on_site`, `remote`, `hybrid` |
| `experience` | `entry`, `mid`, `senior`, `lead`, `executive` |
| `status` | `active` (default), `expired`, `filled` |
| `posted_after` | ISO date (e.g., `2024-01-01`) |
| `page` | Page number (default: 1) |
| `page_size` | Results per page (default: 20) |

### `GET /admin/jobs/{id}`
Get a single job by UUID.

### `GET /admin/career-pages`
List all configured company career pages.

---

## Scraper Details

### LinkedIn Jobs Scraper

Uses LinkedIn's public job search API (`/jobs-guest/jobs/api/seeMoreJobPostings/search`).
- Rate limit: 5 requests/minute
- Max pages: 10 (250 results)
- Respects robots.txt

### Indeed Scraper

Uses Indeed's public job search (`/jobs`).
- Rate limit: 6 requests/minute
- Max pages: 20 (300 results)
- Respects robots.txt

### Company Career Page Scraper

Configurable via the `company_career_pages` database table.
Supports two modes:
1. **JSON API mode**: Set `api_endpoint` in the selectors JSON
2. **HTML mode**: Set CSS selectors for job container, title, location, etc.

**Example selectors JSON:**
```json
{
  "job_container": ".job-listing",
  "title": "h2.job-title",
  "location": ".job-location",
  "apply_url": "a.apply-link",
  "posted_date": ".posted-date",
  "next_page_selector": "a.next-page"
}
```

**Example JSON API selectors:**
```json
{
  "api_endpoint": "https://company.com/api/jobs",
  "api_jobs_field": "positions"
}
```

---

## Deduplication

Jobs are deduplicated using a SHA-256 hash:
- **With external ID**: `hash(source:external_id)`
- **Without external ID**: `hash(source:title:company:location)`

On conflict, the job's `last_seen_at` and description are updated.

---

## Scheduler

The scheduler runs all scrapers concurrently with a configurable worker pool:

```go
schedConfig := scheduler.DefaultConfig()
// schedConfig.WorkerCount = 5
// schedConfig.DefaultQueries = []SearchQuery{...}
// schedConfig.JobStaleDuration = 7 * 24 * time.Hour
```

Jobs not seen within `JobStaleDuration` are automatically marked as `expired`.

---

## Running Tests

```bash
go test ./...
go test -cover ./...
```

| Package | Coverage |
|---------|----------|
| `internal/httpclient` | ~85% |
| `internal/scraper` | ~80% |
| `internal/storage` | ~75% |
