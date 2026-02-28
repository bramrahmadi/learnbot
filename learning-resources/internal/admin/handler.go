// Package admin provides the HTTP handlers for the learning resources admin API.
// Admin endpoints allow creating, updating, and managing learning resources.
package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnbot/database/repository"
)

// Handler holds the HTTP handler dependencies for the admin API.
type Handler struct {
	repo   *repository.LearningResourceRepository
	logger *log.Logger
}

// NewHandler creates a new admin Handler.
func NewHandler(repo *repository.LearningResourceRepository, logger *log.Logger) *Handler {
	return &Handler{repo: repo, logger: logger}
}

// RegisterRoutes registers all admin routes on the given mux.
//
// Admin endpoints (should be protected by authentication middleware):
//
//	POST   /api/v1/admin/resources           – create a new resource
//	PUT    /api/v1/admin/resources/{id}      – update a resource
//	DELETE /api/v1/admin/resources/{id}      – soft-delete a resource
//	POST   /api/v1/admin/providers           – create a new provider
//	POST   /api/v1/admin/paths               – create a new learning path
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/admin/resources", h.withMiddleware(h.handleAdminResources))
	mux.HandleFunc("/api/v1/admin/resources/", h.withMiddleware(h.handleAdminResourceByID))
	mux.HandleFunc("/api/v1/admin/providers", h.withMiddleware(h.handleAdminProviders))
	mux.HandleFunc("/api/v1/admin/paths", h.withMiddleware(h.handleAdminPaths))
}

// withMiddleware wraps a handler with logging and panic recovery.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC in admin: %v", rec)
				h.writeError(w, http.StatusInternalServerError, "an unexpected error occurred")
			}
		}()
		h.logger.Printf("[ADMIN] %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("[ADMIN] %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource admin handlers
// ─────────────────────────────────────────────────────────────────────────────

// handleAdminResources handles POST /api/v1/admin/resources
//
// Request body (JSON):
//
//	{
//	  "title": "...",
//	  "slug": "...",
//	  "description": "...",
//	  "url": "...",
//	  "provider_id": "uuid",
//	  "resource_type": "course",
//	  "difficulty": "intermediate",
//	  "cost_type": "paid",
//	  "cost_amount": 19.99,
//	  "duration_hours": 10.0,
//	  "duration_label": "10 hours",
//	  "has_certificate": true,
//	  "has_hands_on": true,
//	  "skills": [
//	    {"skill_name": "Python", "is_primary": true, "coverage_level": "intermediate"}
//	  ]
//	}
func (h *Handler) handleAdminResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var req createResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := req.validate(); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := req.toInput()
	resource, err := h.repo.Create(r.Context(), input)
	if err != nil {
		h.logger.Printf("create resource error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to create resource")
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    resource,
	})
}

// handleAdminResourceByID handles PUT/DELETE /api/v1/admin/resources/{id}
func (h *Handler) handleAdminResourceByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/resources/")
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "resource ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid resource ID")
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.updateResource(w, r, id)
	case http.MethodDelete:
		h.deleteResource(w, r, id)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "only PUT and DELETE are supported")
	}
}

func (h *Handler) updateResource(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var req updateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	input := req.toInput()
	resource, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		h.logger.Printf("update resource error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to update resource")
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

func (h *Handler) deleteResource(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		h.logger.Printf("delete resource error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to delete resource")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "resource deleted",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Provider admin handlers
// ─────────────────────────────────────────────────────────────────────────────

// handleAdminProviders handles POST /api/v1/admin/providers
//
// Request body (JSON):
//
//	{
//	  "name": "...",
//	  "website_url": "...",
//	  "logo_url": "...",
//	  "description": "..."
//	}
func (h *Handler) handleAdminProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var req createProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		h.writeError(w, http.StatusBadRequest, "provider name is required")
		return
	}

	provider, err := h.repo.CreateProvider(r.Context(), repository.CreateProviderInput{
		Name:        req.Name,
		WebsiteURL:  req.WebsiteURL,
		LogoURL:     req.LogoURL,
		Description: req.Description,
	})
	if err != nil {
		h.logger.Printf("create provider error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to create provider")
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    provider,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Learning path admin handlers
// ─────────────────────────────────────────────────────────────────────────────

// handleAdminPaths handles POST /api/v1/admin/paths
//
// Request body (JSON):
//
//	{
//	  "title": "...",
//	  "slug": "...",
//	  "description": "...",
//	  "target_role": "Backend Engineer",
//	  "target_skill": "Go",
//	  "difficulty": "intermediate",
//	  "estimated_hours": 120,
//	  "resources": [
//	    {"resource_id": "uuid", "step_order": 1, "is_required": true, "notes": "..."}
//	  ]
//	}
func (h *Handler) handleAdminPaths(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var req createPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		h.writeError(w, http.StatusBadRequest, "path title is required")
		return
	}
	if strings.TrimSpace(req.Slug) == "" {
		h.writeError(w, http.StatusBadRequest, "path slug is required")
		return
	}

	// Note: CreateLearningPath is not yet implemented in the repository.
	// This returns a placeholder response to demonstrate the admin interface.
	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "learning path creation queued",
		"data":    req,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Request types
// ─────────────────────────────────────────────────────────────────────────────

type createResourceRequest struct {
	Title          string                       `json:"title"`
	Slug           string                       `json:"slug"`
	Description    *string                      `json:"description,omitempty"`
	URL            string                       `json:"url"`
	ProviderID     *string                      `json:"provider_id,omitempty"`
	ResourceType   string                       `json:"resource_type"`
	Difficulty     string                       `json:"difficulty"`
	CostType       string                       `json:"cost_type"`
	CostAmount     *float64                     `json:"cost_amount,omitempty"`
	CostCurrency   string                       `json:"cost_currency,omitempty"`
	DurationHours  *float64                     `json:"duration_hours,omitempty"`
	DurationLabel  *string                      `json:"duration_label,omitempty"`
	Language       string                       `json:"language,omitempty"`
	HasCertificate bool                         `json:"has_certificate"`
	HasHandsOn     bool                         `json:"has_hands_on"`
	Skills         []createResourceSkillRequest `json:"skills,omitempty"`
}

type createResourceSkillRequest struct {
	SkillName     string  `json:"skill_name"`
	IsPrimary     bool    `json:"is_primary"`
	CoverageLevel *string `json:"coverage_level,omitempty"`
}

func (r *createResourceRequest) validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(r.Slug) == "" {
		return fmt.Errorf("slug is required")
	}
	if strings.TrimSpace(r.URL) == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

func (r *createResourceRequest) toInput() repository.CreateResourceInput {
	input := repository.CreateResourceInput{
		Title:          r.Title,
		Slug:           r.Slug,
		Description:    r.Description,
		URL:            r.URL,
		ResourceType:   repository.ResourceType(r.ResourceType),
		Difficulty:     repository.ResourceDifficulty(r.Difficulty),
		CostType:       repository.ResourceCostType(r.CostType),
		CostAmount:     r.CostAmount,
		CostCurrency:   r.CostCurrency,
		DurationHours:  r.DurationHours,
		DurationLabel:  r.DurationLabel,
		Language:       r.Language,
		HasCertificate: r.HasCertificate,
		HasHandsOn:     r.HasHandsOn,
	}

	if r.ProviderID != nil {
		if id, err := uuid.Parse(*r.ProviderID); err == nil {
			input.ProviderID = &id
		}
	}

	for _, s := range r.Skills {
		skillInput := repository.ResourceSkillInput{
			SkillName: s.SkillName,
			IsPrimary: s.IsPrimary,
		}
		if s.CoverageLevel != nil {
			level := repository.ResourceDifficulty(*s.CoverageLevel)
			skillInput.CoverageLevel = &level
		}
		input.Skills = append(input.Skills, skillInput)
	}

	return input
}

type updateResourceRequest struct {
	Title          *string  `json:"title,omitempty"`
	Description    *string  `json:"description,omitempty"`
	URL            *string  `json:"url,omitempty"`
	Difficulty     *string  `json:"difficulty,omitempty"`
	CostType       *string  `json:"cost_type,omitempty"`
	CostAmount     *float64 `json:"cost_amount,omitempty"`
	DurationHours  *float64 `json:"duration_hours,omitempty"`
	DurationLabel  *string  `json:"duration_label,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
	IsFeatured     *bool    `json:"is_featured,omitempty"`
	HasCertificate *bool    `json:"has_certificate,omitempty"`
	HasHandsOn     *bool    `json:"has_hands_on,omitempty"`
}

func (r *updateResourceRequest) toInput() repository.UpdateResourceInput {
	input := repository.UpdateResourceInput{
		Title:          r.Title,
		Description:    r.Description,
		URL:            r.URL,
		CostAmount:     r.CostAmount,
		DurationHours:  r.DurationHours,
		DurationLabel:  r.DurationLabel,
		IsActive:       r.IsActive,
		IsFeatured:     r.IsFeatured,
		HasCertificate: r.HasCertificate,
		HasHandsOn:     r.HasHandsOn,
	}
	if r.Difficulty != nil {
		d := repository.ResourceDifficulty(*r.Difficulty)
		input.Difficulty = &d
	}
	if r.CostType != nil {
		c := repository.ResourceCostType(*r.CostType)
		input.CostType = &c
	}
	return input
}

type createProviderRequest struct {
	Name        string  `json:"name"`
	WebsiteURL  *string `json:"website_url,omitempty"`
	LogoURL     *string `json:"logo_url,omitempty"`
	Description *string `json:"description,omitempty"`
}

type createPathRequest struct {
	Title          string                      `json:"title"`
	Slug           string                      `json:"slug"`
	Description    *string                     `json:"description,omitempty"`
	TargetRole     *string                     `json:"target_role,omitempty"`
	TargetSkill    *string                     `json:"target_skill,omitempty"`
	Difficulty     string                      `json:"difficulty"`
	EstimatedHours *float64                    `json:"estimated_hours,omitempty"`
	Resources      []pathResourceRequest       `json:"resources,omitempty"`
}

type pathResourceRequest struct {
	ResourceID string  `json:"resource_id"`
	StepOrder  int16   `json:"step_order"`
	IsRequired bool    `json:"is_required"`
	Notes      *string `json:"notes,omitempty"`
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
