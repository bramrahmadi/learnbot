// Package gapanalysis provides the HTTP handler for the skill gap analysis
// API endpoint.
package gapanalysis

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Handler holds the HTTP handler dependencies for the gap analysis API.
type Handler struct {
	analyzer *Analyzer
	logger   *log.Logger
}

// NewHandler creates a new gap analysis Handler.
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		analyzer: New(),
		logger:   logger,
	}
}

// RegisterRoutes registers the gap analysis routes on the given mux.
//
//	POST /api/v1/gap-analysis  – perform skill gap analysis
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/gap-analysis", h.withMiddleware(h.GapAnalysisHandler))
}

// withMiddleware wraps a handler with logging and panic recovery.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC in gap-analysis: %v", rec)
				h.writeError(w, http.StatusInternalServerError,
					"an unexpected error occurred")
			}
		}()
		h.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// GapAnalysisHandler handles POST /api/v1/gap-analysis.
//
// Request body (JSON):
//
//	{
//	  "profile": { ... CandidateProfile ... },
//	  "job":     { ... JobRequirements  ... }
//	}
//
// Response body (JSON):
//
//	{
//	  "success": true,
//	  "data": { ... GapAnalysisResult ... }
//	}
//
// The response includes:
//   - critical_gaps: must-have skills the candidate is missing
//   - important_gaps: preferred skills the candidate is missing
//   - nice_to_have_gaps: optional skills the candidate is missing
//   - top_priority_gaps: top 5 gaps ranked by priority score
//   - readiness_score: overall readiness percentage [0, 100]
//   - total_estimated_learning_hours: total hours to close all gaps
//   - visual_data: JSON-friendly data for frontend visualization
//
// Example curl:
//
//	curl -X POST http://localhost:8080/api/v1/gap-analysis \
//	  -H "Content-Type: application/json" \
//	  -d '{"profile":{...},"job":{...}}'
func (h *Handler) GapAnalysisHandler(w http.ResponseWriter, r *http.Request) {
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

	var req GapAnalysisRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest,
			"invalid request body: "+err.Error())
		return
	}

	result := h.analyzer.Analyze(req.Profile, req.Job)

	h.writeJSON(w, http.StatusOK, GapAnalysisResponse{
		Success: true,
		Data:    &result,
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
	h.writeJSON(w, status, GapAnalysisResponse{
		Success: false,
		Error:   message,
	})
}
