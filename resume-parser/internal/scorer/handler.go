// Package scorer provides the HTTP handler for the acceptance likelihood
// scoring API endpoint.
package scorer

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Handler holds the HTTP handler dependencies for the scoring API.
type Handler struct {
	logger *log.Logger
}

// NewHandler creates a new scoring Handler.
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{logger: logger}
}

// RegisterRoutes registers the scoring routes on the given mux.
//
//	POST /api/v1/score  – calculate acceptance likelihood score
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/score", h.withMiddleware(h.ScoreHandler))
}

// withMiddleware wraps a handler with logging and panic recovery.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC in scorer: %v", rec)
				h.writeError(w, http.StatusInternalServerError,
					"an unexpected error occurred")
			}
		}()
		h.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// ScoreHandler handles POST /api/v1/score.
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
//	  "data": { ... ScoreBreakdown ... }
//	}
//
// Example curl:
//
//	curl -X POST http://localhost:8080/api/v1/score \
//	  -H "Content-Type: application/json" \
//	  -d '{"profile":{...},"job":{...}}'
func (h *Handler) ScoreHandler(w http.ResponseWriter, r *http.Request) {
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

	var req ScoreRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest,
			"invalid request body: "+err.Error())
		return
	}

	breakdown := Calculate(req.Profile, req.Job)

	h.writeJSON(w, http.StatusOK, ScoreResponse{
		Success: true,
		Data:    &breakdown,
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
	h.writeJSON(w, status, ScoreResponse{
		Success: false,
		Error:   message,
	})
}
