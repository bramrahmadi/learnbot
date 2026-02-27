// Command server starts the job aggregation service.
package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/learnbot/job-aggregator/internal/admin"
	"github.com/learnbot/job-aggregator/internal/scheduler"
	"github.com/learnbot/job-aggregator/internal/scraper"
	"github.com/learnbot/job-aggregator/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8081", "HTTP server address")
	dbURL := flag.String("db", getEnv("DATABASE_URL", "postgres://localhost/learnbot?sslmode=disable"), "PostgreSQL connection URL")
	runNow := flag.Bool("run-now", false, "Run scrapers immediately on startup")
	flag.Parse()

	logger := log.New(os.Stdout, "[job-aggregator] ", log.LstdFlags|log.Lshortfile)

	// Connect to database
	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Printf("warning: database not available: %v (continuing without DB)", err)
	}

	// Initialize scrapers
	linkedInScraper, err := scraper.NewLinkedInScraper(logger)
	if err != nil {
		logger.Fatalf("failed to create LinkedIn scraper: %v", err)
	}

	indeedScraper, err := scraper.NewIndeedScraper(logger)
	if err != nil {
		logger.Fatalf("failed to create Indeed scraper: %v", err)
	}

	scrapers := []scraper.Scraper{linkedInScraper, indeedScraper}

	// Load career page scrapers from database
	repo := storage.NewJobRepository(db)
	careerPages, err := repo.GetCareerPages(context.Background())
	if err != nil {
		logger.Printf("warning: failed to load career pages: %v", err)
	}
	for _, page := range careerPages {
		cpScraper, err := scraper.NewCareerPageScraper(page, logger)
		if err != nil {
			logger.Printf("warning: failed to create career page scraper for %s: %v", page.CompanyName, err)
			continue
		}
		scrapers = append(scrapers, cpScraper)
	}

	logger.Printf("initialized %d scrapers", len(scrapers))

	// Initialize scheduler
	schedConfig := scheduler.DefaultConfig()
	sched := scheduler.New(db, scrapers, schedConfig, logger)

	// Set up HTTP server
	mux := http.NewServeMux()
	adminHandler := admin.NewHandler(repo, sched, logger)
	adminHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start daily schedule
	sched.StartDailySchedule(ctx)

	// Optionally run immediately
	if *runNow {
		logger.Println("running scrapers immediately (--run-now flag)")
		sched.RunNow(ctx)
	}

	go func() {
		logger.Printf("starting server on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	logger.Println("shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("forced shutdown: %v", err)
	}

	logger.Println("server stopped")
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
