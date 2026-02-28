// Package api provides the HTTP handlers for the learning resources API.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnbot/database/repository"
)

// Handler holds the HTTP handler dependencies for the learning resources API.
type Handler struct {
	repo   *repository.LearningResourceRepository
	logger *log.Logger
}

// NewHandler creates a new learning resources Handler.
func NewHandler(repo *repository.LearningResourceRepository, logger *log.Logger) *Handler {
	return &Handler{repo: repo, logger: logger}
}

// RegisterRoutes registers all learning resource routes on the given mux.
//
// Public endpoints:
//
//	GET  /api/v1/resources              – list/search resources
//	GET  /api/v1/resources/{slug}       – get resource by slug
//	GET  /api/v1/resources/featured     – get featured resources
//	GET  /api/v1/resources/by-skill     – get resources for a skill
//	GET  /api/v1/paths                  – list learning paths
//	GET  /api/v1/paths/{slug}           – get learning path by slug
//	GET  /api/v1/providers              – list resource providers
//
// User endpoints (require user context):
//
//	GET  /api/v1/users/{id}/progress    – get user's resource progress
//	POST /api/v1/users/{id}/progress    – update user's resource progress
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/resources", h.withMiddleware(h.handleResources))
	mux.HandleFunc("/api/v1/resources/featured", h.withMiddleware(h.handleFeaturedResources))
	mux.HandleFunc("/api/v1/resources/by-skill", h.withMiddleware(h.handleResourcesBySkill))
	mux.HandleFunc("/api/v1/resources/", h.withMiddleware(h.handleResourceBySlug))
	mux.HandleFunc("/api/v1/paths", h.withMiddleware(h.handlePaths))
	mux.HandleFunc("/api/v1/paths/", h.withMiddleware(h.handlePathBySlug))
	mux.HandleFunc("/api/v1/providers", h.withMiddleware(h.handleProviders))
	mux.HandleFunc("/api/v1/users/", h.withMiddleware(h.handleUserProgress))
}

// withMiddleware wraps a handler with logging and panic recovery.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC in learning-resources: %v", rec)
				h.writeError(w, http.StatusInternalServerError, "an unexpected error occurred")
			}
		}()
		h.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource handlers
// ─────────────────────────────────────────────────────────────────────────────

// handleResources handles GET /api/v1/resources
//
// Query parameters:
//   - skill: filter by skill name
//   - type: filter by resource type (course, certification, etc.)
//   - difficulty: filter by difficulty level
//   - cost_type: filter by cost model
//   - free: "true" to show only free resources
//   - has_certificate: "true" to show only resources with certificates
//   - has_hands_on: "true" to show only resources with hands-on exercises
//   - min_rating: minimum rating (0.0-5.0)
//   - q: full-text search query
//   - limit: max results (default 20, max 100)
//   - offset: pagination offset
func (h *Handler) handleResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	q := r.URL.Query()
	filter := repository.ResourceQueryFilter{
		SkillName:      q.Get("skill"),
		ResourceType:   repository.ResourceType(q.Get("type")),
		Difficulty:     repository.ResourceDifficulty(q.Get("difficulty")),
		CostType:       repository.ResourceCostType(q.Get("cost_type")),
		IsFree:         q.Get("free") == "true",
		HasCertificate: q.Get("has_certificate") == "true",
		HasHandsOn:     q.Get("has_hands_on") == "true",
		SearchQuery:    q.Get("q"),
	}

	if minRating := q.Get("min_rating"); minRating != "" {
		if v, err := strconv.ParseFloat(minRating, 64); err == nil {
			filter.MinRating = v
		}
	}
	if limit := q.Get("limit"); limit != "" {
		if v, err := strconv.Atoi(limit); err == nil {
			filter.Limit = v
		}
	}
	if offset := q.Get("offset"); offset != "" {
		if v, err := strconv.Atoi(offset); err == nil {
			filter.Offset = v
		}
	}

	resources, total, err := h.repo.List(r.Context(), filter)
	if err != nil {
		h.logger.Printf("list resources error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list resources")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resources,
		"total":   total,
		"limit":   filter.Limit,
		"offset":  filter.Offset,
	})
}

// handleFeaturedResources handles GET /api/v1/resources/featured
func (h *Handler) handleFeaturedResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	resources, err := h.repo.GetFeatured(r.Context(), limit)
	if err != nil {
		h.logger.Printf("get featured resources error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get featured resources")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resources,
	})
}

// handleResourcesBySkill handles GET /api/v1/resources/by-skill?skill=Python
func (h *Handler) handleResourcesBySkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	skill := r.URL.Query().Get("skill")
	if skill == "" {
		h.writeError(w, http.StatusBadRequest, "skill parameter is required")
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	resources, err := h.repo.GetBySkill(r.Context(), skill, limit)
	if err != nil {
		h.logger.Printf("get resources by skill error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get resources for skill")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"skill":   skill,
		"data":    resources,
	})
}

// handleResourceBySlug handles GET /api/v1/resources/{slug}
func (h *Handler) handleResourceBySlug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	// Extract slug from path: /api/v1/resources/{slug}
	slug := strings.TrimPrefix(r.URL.Path, "/api/v1/resources/")
	if slug == "" {
		h.writeError(w, http.StatusBadRequest, "resource slug is required")
		return
	}

	resource, err := h.repo.GetBySlug(r.Context(), slug)
	if err != nil {
		h.logger.Printf("get resource by slug error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get resource")
		return
	}
	if resource == nil {
		h.writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resource,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Learning path handlers
// ─────────────────────────────────────────────────────────────────────────────

// handlePaths handles GET /api/v1/paths
//
// Query parameters:
//   - skill: filter by target skill
//   - role: filter by target role
func (h *Handler) handlePaths(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	skill := r.URL.Query().Get("skill")
	role := r.URL.Query().Get("role")

	paths, err := h.repo.ListPaths(r.Context(), skill, role)
	if err != nil {
		h.logger.Printf("list paths error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list learning paths")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    paths,
	})
}

// handlePathBySlug handles GET /api/v1/paths/{slug}
func (h *Handler) handlePathBySlug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/v1/paths/")
	if slug == "" {
		h.writeError(w, http.StatusBadRequest, "path slug is required")
		return
	}

	path, err := h.repo.GetPathBySlug(r.Context(), slug)
	if err != nil {
		h.logger.Printf("get path by slug error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get learning path")
		return
	}
	if path == nil {
		h.writeError(w, http.StatusNotFound, "learning path not found")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    path,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Provider handlers
// ─────────────────────────────────────────────────────────────────────────────

// handleProviders handles GET /api/v1/providers
func (h *Handler) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	providers, err := h.repo.ListProviders(r.Context())
	if err != nil {
		h.logger.Printf("list providers error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list providers")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    providers,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// User progress handlers
// ─────────────────────────────────────────────────────────────────────────────

// handleUserProgress handles GET/POST /api/v1/users/{user_id}/progress
func (h *Handler) handleUserProgress(w http.ResponseWriter, r *http.Request) {
	// Parse path: /api/v1/users/{user_id}/progress[/{resource_id}]
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	parts := strings.SplitN(path, "/", 3)

	if len(parts) < 2 || parts[1] != "progress" {
		h.writeError(w, http.StatusNotFound, "not found")
		return
	}

	userID, err := uuid.Parse(parts[0])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUserProgress(w, r, userID, parts)
	case http.MethodPost:
		h.updateUserProgress(w, r, userID, parts)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "only GET and POST are supported")
	}
}

func (h *Handler) getUserProgress(w http.ResponseWriter, r *http.Request, userID uuid.UUID, parts []string) {
	// GET /api/v1/users/{user_id}/progress[?status=in_progress]
	status := repository.UserResourceStatus(r.URL.Query().Get("status"))

	progress, err := h.repo.ListUserProgress(r.Context(), userID, status)
	if err != nil {
		h.logger.Printf("list user progress error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get user progress")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    progress,
	})
}

func (h *Handler) updateUserProgress(w http.ResponseWriter, r *http.Request, userID uuid.UUID, parts []string) {
	// POST /api/v1/users/{user_id}/progress/{resource_id}
	if len(parts) < 3 {
		h.writeError(w, http.StatusBadRequest, "resource ID is required")
		return
	}

	resourceID, err := uuid.Parse(parts[2])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid resource ID")
		return
	}

	var input repository.UpsertProgressInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	progress, err := h.repo.UpsertUserProgress(r.Context(), userID, resourceID, input)
	if err != nil {
		h.logger.Printf("upsert user progress error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to update progress")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    progress,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Response helpers
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Printf("failed to encode JSON response: %v", err)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
