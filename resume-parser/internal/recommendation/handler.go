// Package recommendation provides the HTTP handler for the training
// recommendation API endpoint.
package recommendation

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Handler holds the HTTP handler dependencies for the recommendation API.
type Handler struct {
	engine *Engine
	logger *log.Logger
}

// NewHandler creates a new recommendation Handler.
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		engine: New(),
		logger: logger,
	}
}

// RegisterRoutes registers the recommendation routes on the given mux.
//
//	POST /api/v1/recommendations  – generate a personalized learning plan
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/recommendations", h.withMiddleware(h.RecommendationHandler))
}

// withMiddleware wraps a handler with logging and panic recovery.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC in recommendation: %v", rec)
				h.writeError(w, http.StatusInternalServerError,
					"an unexpected error occurred")
			}
		}()
		h.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// RecommendationHandler handles POST /api/v1/recommendations.
//
// Request body (JSON):
//
//	{
//	  "profile": { ... CandidateProfile ... },
//	  "job":     { ... JobRequirements  ... },
//	  "preferences": {
//	    "prefer_free": false,
//	    "max_budget_usd": 100,
//	    "weekly_hours_available": 10,
//	    "prefer_hands_on": true,
//	    "prefer_certificates": false,
//	    "target_date": "2025-06-01",
//	    "preferred_resource_types": ["course", "documentation"],
//	    "excluded_providers": []
//	  }
//	}
//
// Response body (JSON):
//
//	{
//	  "success": true,
//	  "data": {
//	    "job_title": "...",
//	    "readiness_score": 45.0,
//	    "total_gaps": 3,
//	    "total_estimated_hours": 120.0,
//	    "phases": [...],
//	    "timeline": {...},
//	    "matched_skills": [...],
//	    "summary": {...}
//	  }
//	}
//
// Example curl:
//
//	curl -X POST http://localhost:8080/api/v1/recommendations \
//	  -H "Content-Type: application/json" \
//	  -d '{"profile":{...},"job":{...},"preferences":{"weekly_hours_available":10}}'
func (h *Handler) RecommendationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed,
			"only POST is supported")
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "" &&
		len(ct) >= 16 && ct[:16] != "application/json" {
		h.writeError(w, http.StatusUnsupportedMediaType,
			"Content-Type must be application/json")
		return
	}

	var req RecommendationRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest,
			"invalid request body: "+err.Error())
		return
	}

	plan := h.engine.Generate(req.Profile, req.Job, req.Preferences)

	h.writeJSON(w, http.StatusOK, RecommendationResponse{
		Success: true,
		Data:    &plan,
	})
}

// writeJSON serialises v as JSON and writes it to the response.
func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Printf("failed to encode JSON response: %v", err)
	}
}

// writeError writes a structured error response.
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, RecommendationResponse{
		Success: false,
		Error:   message,
	})
}
