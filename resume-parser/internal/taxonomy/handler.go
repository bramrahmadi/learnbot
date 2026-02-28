// Package taxonomy – handler.go provides the HTTP API for skill taxonomy
// operations.
//
// Endpoints:
//
//	POST /api/v1/skills/extract    – extract skills from free-form text
//	POST /api/v1/skills/normalize  – normalize raw skill strings to taxonomy
//	GET  /api/v1/skills/lookup     – look up a skill by canonical ID
//	GET  /api/v1/skills/search     – search the taxonomy
package taxonomy

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handler holds the HTTP handler dependencies for the taxonomy API.
type Handler struct {
	taxonomy  *Taxonomy
	extractor *Extractor
	logger    *log.Logger
}

// NewHandler creates a new taxonomy Handler.
func NewHandler(logger *log.Logger) *Handler {
	t := New()
	return &Handler{
		taxonomy:  t,
		extractor: NewExtractor(t),
		logger:    logger,
	}
}

// RegisterRoutes registers the taxonomy routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/skills/extract", h.withMiddleware(h.ExtractHandler))
	mux.HandleFunc("/api/v1/skills/normalize", h.withMiddleware(h.NormalizeHandler))
	mux.HandleFunc("/api/v1/skills/lookup", h.withMiddleware(h.LookupHandler))
	mux.HandleFunc("/api/v1/skills/search", h.withMiddleware(h.SearchHandler))
}

// withMiddleware wraps a handler with logging and panic recovery.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC in taxonomy: %v", rec)
				h.writeError(w, http.StatusInternalServerError,
					"an unexpected error occurred")
			}
		}()
		h.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// ExtractHandler handles POST /api/v1/skills/extract
//
// Request body (JSON):
//
//	{
//	  "text": "We are looking for a Go developer with experience in Kubernetes...",
//	  "include_unknown": false
//	}
//
// Response body (JSON):
//
//	{
//	  "success": true,
//	  "data": {
//	    "skills": [...],
//	    "technical_skills": [...],
//	    "soft_skills": [...],
//	    "domain_skills": [...],
//	    "unknown_skills": [...]
//	  }
//	}
func (h *Handler) ExtractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	result := h.extractor.Extract(req.Text, req.IncludeUnknown)
	h.writeJSON(w, http.StatusOK, ExtractResponse{
		Success: true,
		Data:    &result,
	})
}

// NormalizeHandler handles POST /api/v1/skills/normalize
//
// Request body (JSON):
//
//	{
//	  "skills": ["golang", "k8s", "react.js", "ML"]
//	}
//
// Response body (JSON):
//
//	{
//	  "success": true,
//	  "data": [
//	    {"input": "golang", "canonical_id": "go", "canonical_name": "Go", ...},
//	    ...
//	  ]
//	}
func (h *Handler) NormalizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var req NormalizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if len(req.Skills) == 0 {
		h.writeError(w, http.StatusBadRequest, "skills array must not be empty")
		return
	}

	results := h.taxonomy.NormalizeMany(req.Skills)
	h.writeJSON(w, http.StatusOK, NormalizeResponse{
		Success: true,
		Data:    results,
	})
}

// LookupHandler handles GET /api/v1/skills/lookup?id=<canonical-id>
//
// Query parameters:
//
//	id – the canonical skill ID (e.g. "go", "python", "kubernetes")
//
// Response body (JSON):
//
//	{
//	  "success": true,
//	  "data": { "id": "go", "canonical_name": "Go", ... }
//	}
func (h *Handler) LookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "query parameter 'id' is required")
		return
	}

	node := h.taxonomy.Lookup(id)
	if node == nil {
		h.writeJSON(w, http.StatusNotFound, LookupResponse{
			Success: false,
			Error:   "skill not found: " + id,
		})
		return
	}

	h.writeJSON(w, http.StatusOK, LookupResponse{
		Success: true,
		Data:    node,
	})
}

// SearchHandler handles GET /api/v1/skills/search?q=<query>[&domain=...][&category=...][&limit=...]
//
// Query parameters:
//
//	q        – search term (required)
//	domain   – filter by domain (optional)
//	category – filter by category (optional)
//	limit    – max results (optional, default 20)
//
// Response body (JSON):
//
//	{
//	  "success": true,
//	  "data": [...],
//	  "total": 5
//	}
func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	domain := Domain(strings.TrimSpace(r.URL.Query().Get("domain")))
	category := Category(strings.TrimSpace(r.URL.Query().Get("category")))

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	results := h.taxonomy.Search(q, domain, category, limit)
	h.writeJSON(w, http.StatusOK, SearchResponse{
		Success: true,
		Data:    results,
		Total:   len(results),
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
	h.writeJSON(w, status, ExtractResponse{
		Success: false,
		Error:   message,
	})
}
