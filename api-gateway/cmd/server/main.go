// Command server starts the LearnBot API gateway.
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

	"github.com/learnbot/api-gateway/internal/handler"
	"github.com/learnbot/api-gateway/internal/middleware"
)

func main() {
	addr := flag.String("addr", ":8090", "HTTP server address")
	jwtSecret := flag.String("jwt-secret", os.Getenv("JWT_SECRET"), "JWT signing secret")
	flag.Parse()

	logger := log.New(os.Stdout, "[api-gateway] ", log.LstdFlags|log.Lshortfile)

	// JWT configuration.
	jwtCfg := middleware.DefaultJWTConfig(*jwtSecret)

	// Rate limiter: 10 requests/second, burst of 30.
	rateLimiter := middleware.NewRateLimiter(10, 30)

	// Create handlers.
	authHandler := handler.NewAuthHandler(jwtCfg)
	profileHandler := handler.NewProfileHandler(jwtCfg)
	resumeHandler := handler.NewResumeHandler()
	jobsHandler := handler.NewJobsHandler()
	analysisHandler := handler.NewAnalysisHandler()
	resourcesHandler := handler.NewResourcesHandler()

	// Auth middleware factory.
	authMiddleware := middleware.RequireAuth(jwtCfg)

	// Build mux.
	mux := http.NewServeMux()

	// Register routes.
	authHandler.RegisterRoutes(mux)
	profileHandler.RegisterRoutes(mux, authMiddleware)
	resumeHandler.RegisterRoutes(mux, authMiddleware)
	jobsHandler.RegisterRoutes(mux, authMiddleware)
	analysisHandler.RegisterRoutes(mux, authMiddleware)
	resourcesHandler.RegisterRoutes(mux)

	// Health check.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"api-gateway","version":"1.0.0"}`))
	})

	// Apply global middleware chain.
	globalChain := middleware.Chain(
		middleware.Recovery(logger),
		middleware.Logger(logger),
		middleware.CORS([]string{"*"}),
		rateLimiter.Middleware,
	)

	srv := &http.Server{
		Addr:         *addr,
		Handler:      globalChain(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("starting API gateway on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	logger.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("forced shutdown: %v", err)
	}
	logger.Println("server stopped")
}
