// Package admin provides the HTTP admin dashboard for monitoring scraping status.
package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/learnbot/job-aggregator/internal/model"
	"github.com/learnbot/job-aggregator/internal/scheduler"
	"github.com/learnbot/job-aggregator/internal/storage"
)

// Handler provides HTTP endpoints for the admin dashboard.
type Handler struct {
	repo      *storage.JobRepository
	scheduler *scheduler.Scheduler
	logger    *log.Logger
}

// NewHandler creates a new admin Handler.
func NewHandler(repo *storage.JobRepository, sched *scheduler.Scheduler, logger *log.Logger) *Handler {
	return &Handler{
		repo:      repo,
		scheduler: sched,
		logger:    logger,
	}
}

// RegisterRoutes registers all admin routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Dashboard stats
	mux.HandleFunc("/admin/stats", h.GetStats)
	// Scrape runs
	mux.HandleFunc("/admin/runs", h.GetRecentRuns)
	// Trigger manual scrape
	mux.HandleFunc("/admin/scrape/trigger", h.TriggerScrape)
	// Job management
	mux.HandleFunc("/admin/jobs", h.SearchJobs)
	mux.HandleFunc("/admin/jobs/", h.GetJob)
	// Career pages
	mux.HandleFunc("/admin/career-pages", h.GetCareerPages)
	// Health
	mux.HandleFunc("/admin/health", h.Health)
}

// GetStats returns aggregated scraping statistics.
// GET /admin/stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats, err := h.repo.GetAdminStats(r.Context())
	if err != nil {
		h.logger.Printf("[admin] GetStats error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get stats")
		return
	}

	// Add scheduler status
	response := map[string]interface{}{
		"stats":      stats,
		"is_running": h.scheduler.IsRunning(),
		"timestamp":  time.Now().UTC(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// GetRecentRuns returns the most recent scraping runs.
// GET /admin/runs?limit=20
func (h *Handler) GetRecentRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	runs, err := h.repo.GetRecentScrapeRuns(r.Context(), limit)
	if err != nil {
		h.logger.Printf("[admin] GetRecentRuns error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get runs")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"runs":  runs,
		"count": len(runs),
	})
}

// TriggerScrape triggers an immediate scraping run.
// POST /admin/scrape/trigger
func (h *Handler) TriggerScrape(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if h.scheduler.IsRunning() {
		h.writeJSON(w, http.StatusConflict, map[string]interface{}{
			"message": "scraper is already running",
			"running": true,
		})
		return
	}

	h.scheduler.RunNow(r.Context())

	h.writeJSON(w, http.StatusAccepted, map[string]interface{}{
		"message": "scraping run triggered",
		"running": true,
	})
}

// SearchJobs searches for jobs with filters.
// GET /admin/jobs?q=engineer&location=remote&page=1&page_size=20
func (h *Handler) SearchJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	q := r.URL.Query()

	filter := model.JobFilter{
		TitleSearch: q.Get("q"),
		CompanyName: q.Get("company"),
		Status:      model.StatusActive,
	}

	// Parse page
	if p := q.Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			filter.Page = n
		}
	}
	if ps := q.Get("page_size"); ps != "" {
		if n, err := strconv.Atoi(ps); err == nil {
			filter.PageSize = n
		}
	}

	// Parse location type
	if lt := q.Get("location_type"); lt != "" {
		filter.LocationTypes = []model.WorkLocationType{model.WorkLocationType(lt)}
	}

	// Parse experience level
	if el := q.Get("experience"); el != "" {
		filter.ExperienceLevels = []model.ExperienceLevel{model.ExperienceLevel(el)}
	}

	// Parse status
	if s := q.Get("status"); s != "" {
		filter.Status = model.JobStatus(s)
	}

	// Parse posted_after
	if pa := q.Get("posted_after"); pa != "" {
		if t, err := time.Parse("2006-01-02", pa); err == nil {
			filter.PostedAfter = &t
		}
	}

	jobs, total, err := h.repo.SearchJobs(r.Context(), filter)
	if err != nil {
		h.logger.Printf("[admin] SearchJobs error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to search jobs")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"jobs":      jobs,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

// GetJob retrieves a single job by ID.
// GET /admin/jobs/{id}
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract ID from path
	path := r.URL.Path
	idStr := path[len("/admin/jobs/"):]
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "job ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid job ID format")
		return
	}

	job, err := h.repo.GetJobByID(r.Context(), id)
	if err == storage.ErrNotFound {
		h.writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if err != nil {
		h.logger.Printf("[admin] GetJob error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get job")
		return
	}

	h.writeJSON(w, http.StatusOK, job)
}

// GetCareerPages returns all configured company career pages.
// GET /admin/career-pages
func (h *Handler) GetCareerPages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	pages, err := h.repo.GetCareerPages(r.Context())
	if err != nil {
		h.logger.Printf("[admin] GetCareerPages error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get career pages")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"career_pages": pages,
		"count":        len(pages),
	})
}

// Health returns the health status of the service.
// GET /admin/health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":     "ok",
		"version":    "1.0.0",
		"is_running": h.scheduler.IsRunning(),
		"time":       time.Now().UTC(),
	})
}

// writeJSON serializes v as JSON and writes it to the response.
func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Printf("[admin] JSON encode error: %v", err)
	}
}

// writeError writes a JSON error response.
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
