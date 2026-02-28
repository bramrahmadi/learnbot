// Command server starts the resume parser HTTP API server.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/learnbot/resume-parser/internal/api"
	"github.com/learnbot/resume-parser/internal/parser"
	"github.com/learnbot/resume-parser/internal/scorer"
	"github.com/learnbot/resume-parser/internal/taxonomy"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	logger := log.New(os.Stdout, "[resume-parser] ", log.LstdFlags|log.Lshortfile)

	resumeParser := parser.NewResumeParser()
	handler := api.NewHandler(resumeParser, logger)
	scorerHandler := scorer.NewHandler(logger)
	taxonomyHandler := taxonomy.NewHandler(logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	scorerHandler.RegisterRoutes(mux)
	taxonomyHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("forced shutdown: %v", err)
	}

	logger.Println("server stopped")
}
