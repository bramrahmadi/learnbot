// Command server starts the learning resources HTTP API server.
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

	"github.com/learnbot/database/repository"
	"github.com/learnbot/learning-resources/internal/admin"
	"github.com/learnbot/learning-resources/internal/api"
)

func main() {
	addr := flag.String("addr", ":8081", "HTTP server address")
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	flag.Parse()

	logger := log.New(os.Stdout, "[learning-resources] ", log.LstdFlags|log.Lshortfile)

	if *dsn == "" {
		logger.Fatal("DATABASE_URL environment variable or -dsn flag is required")
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Configure connection pool.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connectivity.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	logger.Println("connected to database")

	repo := repository.NewLearningResourceRepository(db)
	apiHandler := api.NewHandler(repo, logger)
	adminHandler := admin.NewHandler(repo, logger)

	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)
	adminHandler.RegisterRoutes(mux)

	// Health check endpoint.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"learning-resources"}`))
	})

	srv := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("starting server on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	logger.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("forced shutdown: %v", err)
	}

	logger.Println("server stopped")
}
