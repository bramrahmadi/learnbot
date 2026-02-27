// Package scheduler provides the job scraping scheduler with concurrent
// worker pool processing.
package scheduler

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/learnbot/job-aggregator/internal/model"
	"github.com/learnbot/job-aggregator/internal/scraper"
	"github.com/learnbot/job-aggregator/internal/storage"
)

// Config holds scheduler configuration.
type Config struct {
	// Number of concurrent workers for processing scraped jobs
	WorkerCount int
	// Buffer size for the jobs channel
	JobChannelBuffer int
	// Default search queries to run
	DefaultQueries []SearchQuery
	// How long before a job is considered stale and marked expired
	JobStaleDuration time.Duration
}

// SearchQuery defines a search to run on each scrape cycle.
type SearchQuery struct {
	Query    string
	Location string
	Remote   bool
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		WorkerCount:      5,
		JobChannelBuffer: 100,
		JobStaleDuration: 7 * 24 * time.Hour, // 7 days
		DefaultQueries: []SearchQuery{
			{Query: "software engineer", Location: "United States", Remote: true},
			{Query: "backend developer", Location: "United States", Remote: true},
			{Query: "golang developer", Remote: true},
			{Query: "python developer", Remote: true},
			{Query: "data scientist", Remote: true},
		},
	}
}

// Scheduler orchestrates scraping runs with concurrent job processing.
type Scheduler struct {
	repo     *storage.JobRepository
	scrapers []scraper.Scraper
	config   Config
	logger   *log.Logger
	mu       sync.Mutex
	running  bool
}

// New creates a new Scheduler.
func New(db *sql.DB, scrapers []scraper.Scraper, cfg Config, logger *log.Logger) *Scheduler {
	return &Scheduler{
		repo:     storage.NewJobRepository(db),
		scrapers: scrapers,
		config:   cfg,
		logger:   logger,
	}
}

// RunOnce executes a single scraping cycle for all configured scrapers.
// It uses a worker pool to process scraped jobs concurrently.
func (s *Scheduler) RunOnce(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		s.logger.Println("[scheduler] already running, skipping")
		return nil
	}
	s.running = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	s.logger.Printf("[scheduler] starting scrape cycle with %d scrapers", len(s.scrapers))
	start := time.Now()

	var wg sync.WaitGroup
	for _, sc := range s.scrapers {
		wg.Add(1)
		go func(sc scraper.Scraper) {
			defer wg.Done()
			s.runScraper(ctx, sc)
		}(sc)
	}

	wg.Wait()
	s.logger.Printf("[scheduler] scrape cycle completed in %v", time.Since(start))
	return nil
}

// runScraper runs a single scraper for all configured search queries.
func (s *Scheduler) runScraper(ctx context.Context, sc scraper.Scraper) {
	for _, query := range s.config.DefaultQueries {
		select {
		case <-ctx.Done():
			return
		default:
		}

		params := model.SearchParams{
			Query:    query.Query,
			Location: query.Location,
			Remote:   query.Remote,
			PageSize: 25,
		}

		s.runScraperQuery(ctx, sc, params)
	}
}

// runScraperQuery runs a single scraper for a single search query.
func (s *Scheduler) runScraperQuery(ctx context.Context, sc scraper.Scraper, params model.SearchParams) {
	// Create scrape run log
	run, err := s.repo.CreateScrapeRun(ctx, sc.Source(), params.Query, params.Location)
	if err != nil {
		s.logger.Printf("[scheduler] failed to create scrape run: %v", err)
		return
	}

	s.logger.Printf("[scheduler] %s: scraping query=%q location=%q",
		sc.Name(), params.Query, params.Location)

	// Create jobs channel and start worker pool
	jobsCh := make(chan *model.ScrapedJob, s.config.JobChannelBuffer)
	stats := &scrapeStats{}

	// Start worker pool
	workerCount := s.config.WorkerCount
	if workerCount <= 0 {
		workerCount = 3
	}

	var workerWg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			s.processJobs(ctx, jobsCh, stats)
		}()
	}

	// Run scraper (sends to jobsCh)
	scrapeErr := sc.Scrape(ctx, params, jobsCh)
	close(jobsCh)

	// Wait for all workers to finish
	workerWg.Wait()

	// Update scrape run with final stats
	finalStatus := model.ScrapeStatusCompleted
	errMsg := ""
	if scrapeErr != nil {
		finalStatus = model.ScrapeStatusFailed
		errMsg = scrapeErr.Error()
		s.logger.Printf("[scheduler] %s scrape error: %v", sc.Name(), scrapeErr)
	}

	stats.mu.Lock()
	finalRun := model.ScrapeRun{
		JobsFound:    stats.found,
		JobsNew:      stats.newJobs,
		JobsUpdated:  stats.updated,
		JobsFailed:   stats.failed,
		PagesScraped: stats.pages,
		ErrorMessage: errMsg,
	}
	stats.mu.Unlock()

	if err := s.repo.UpdateScrapeRun(ctx, run.ID, finalStatus, finalRun); err != nil {
		s.logger.Printf("[scheduler] failed to update scrape run: %v", err)
	}

	// Mark stale jobs as expired
	cutoff := time.Now().Add(-s.config.JobStaleDuration)
	expired, err := s.repo.MarkExpiredJobs(ctx, sc.Source(), cutoff)
	if err != nil {
		s.logger.Printf("[scheduler] failed to mark expired jobs: %v", err)
	} else if expired > 0 {
		s.logger.Printf("[scheduler] marked %d jobs as expired for %s", expired, sc.Source())
	}

	s.logger.Printf("[scheduler] %s: found=%d new=%d updated=%d failed=%d",
		sc.Name(), finalRun.JobsFound, finalRun.JobsNew, finalRun.JobsUpdated, finalRun.JobsFailed)
}

// processJobs is a worker that reads from the jobs channel and stores them.
func (s *Scheduler) processJobs(ctx context.Context, jobs <-chan *model.ScrapedJob, stats *scrapeStats) {
	for job := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
		}

		stats.mu.Lock()
		stats.found++
		stats.mu.Unlock()

		_, isNew, err := s.repo.UpsertJob(ctx, job)
		if err != nil {
			s.logger.Printf("[scheduler] failed to upsert job %q at %s: %v",
				job.Title, job.CompanyName, err)
			stats.mu.Lock()
			stats.failed++
			stats.mu.Unlock()
			continue
		}

		stats.mu.Lock()
		if isNew {
			stats.newJobs++
		} else {
			stats.updated++
		}
		stats.mu.Unlock()
	}
}

// scrapeStats holds thread-safe counters for a scraping run.
type scrapeStats struct {
	mu      sync.Mutex
	found   int
	newJobs int
	updated int
	failed  int
	pages   int
}

// StartDailySchedule starts a background goroutine that runs the scraper
// on the configured schedule (default: daily at 2am UTC).
func (s *Scheduler) StartDailySchedule(ctx context.Context) {
	go func() {
		s.logger.Println("[scheduler] daily schedule started")

		for {
			// Calculate next run time (2am UTC)
			now := time.Now().UTC()
			next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, time.UTC)
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}

			waitDuration := time.Until(next)
			s.logger.Printf("[scheduler] next run at %v (in %v)", next, waitDuration)

			select {
			case <-ctx.Done():
				s.logger.Println("[scheduler] daily schedule stopped")
				return
			case <-time.After(waitDuration):
				if err := s.RunOnce(ctx); err != nil {
					s.logger.Printf("[scheduler] daily run error: %v", err)
				}
			}
		}
	}()
}

// RunNow triggers an immediate scraping run (for manual/admin use).
func (s *Scheduler) RunNow(ctx context.Context) {
	go func() {
		if err := s.RunOnce(ctx); err != nil {
			s.logger.Printf("[scheduler] manual run error: %v", err)
		}
	}()
}

// IsRunning returns true if a scraping cycle is currently in progress.
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
