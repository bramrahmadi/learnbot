// Package handler – analysis.go implements gap analysis and training recommendation endpoints.
package handler

import (
	"net/http"
	"strings"

	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/api-gateway/internal/types"
	"github.com/learnbot/resume-parser/pkg/analysis"
	"github.com/learnbot/resume-parser/pkg/recommend"
	"github.com/learnbot/resume-parser/pkg/scoring"
)

// AnalysisHandler handles gap analysis and training recommendation endpoints.
type AnalysisHandler struct {
	gapAnalyzer *analysis.Analyzer
	recEngine   *recommend.Engine
}

// NewAnalysisHandler creates a new AnalysisHandler.
func NewAnalysisHandler() *AnalysisHandler {
	return &AnalysisHandler{
		gapAnalyzer: analysis.New(),
		recEngine:   recommend.New(),
	}
}

// RegisterRoutes registers analysis routes on the mux.
//
//	POST /api/analysis/gaps           – analyze skill gaps for a target job
//	GET  /api/training/recommendations – get personalized training plan
func (h *AnalysisHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/api/analysis/gaps",
		authMiddleware(http.HandlerFunc(h.GapAnalysis)))
	mux.Handle("/api/training/recommendations",
		authMiddleware(http.HandlerFunc(h.TrainingRecommendations)))
}

// GapAnalysis handles POST /api/analysis/gaps.
//
// Analyzes the skill gaps between the current user's profile and a target job.
//
// Request body:
//
//	{
//	  "job_id": "job-001",
//	  "job": {
//	    "title": "Senior Go Engineer",
//	    "required_skills": ["Go", "PostgreSQL", "Docker"],
//	    "preferred_skills": ["Kubernetes", "AWS"],
//	    "min_years_experience": 5,
//	    "experience_level": "senior"
//	  }
//	}
//
// Response includes critical gaps, important gaps, readiness score, and visual data.
func (h *AnalysisHandler) GapAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	userID := middleware.GetUserID(r)

	var req types.GapAnalysisRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	// Resolve job requirements.
	jobReqs, err := h.resolveJobRequirements(req.JobID, req.Job)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JOB",
			"job_id or job details are required")
		return
	}

	// Build candidate profile.
	profile := buildCandidateProfile(userID)

	// Run gap analysis.
	result := h.gapAnalyzer.Analyze(profile, jobReqs)

	WriteSuccess(w, http.StatusOK, result)
}

// TrainingRecommendations handles GET/POST /api/training/recommendations.
//
// Generates a personalized learning plan for the current user.
//
// Query parameters (GET):
//   - job_id: target job ID (optional)
//
// Or POST with body:
//
//	{
//	  "job_id": "job-001",
//	  "job": {...},
//	  "preferences": {
//	    "prefer_free": false,
//	    "weekly_hours_available": 10,
//	    "prefer_hands_on": true
//	  }
//	}
func (h *AnalysisHandler) TrainingRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req types.TrainingRecommendationRequest

	switch r.Method {
	case http.MethodGet:
		// Build request from query params.
		req.JobID = r.URL.Query().Get("job_id")
	case http.MethodPost:
		if !DecodeJSON(w, r, &req) {
			return
		}
	default:
		WriteMethodNotAllowed(w)
		return
	}

	// Resolve job requirements.
	jobReqs, err := h.resolveJobRequirements(req.JobID, req.Job)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JOB",
			"job_id or job details are required")
		return
	}

	// Build candidate profile.
	profile := buildCandidateProfile(userID)

	// Build preferences.
	prefs := recommend.UserPreferences{
		PreferFree:             req.Preferences.PreferFree,
		MaxBudgetUSD:           req.Preferences.MaxBudgetUSD,
		WeeklyHoursAvailable:   req.Preferences.WeeklyHoursAvailable,
		PreferHandsOn:          req.Preferences.PreferHandsOn,
		PreferCertificates:     req.Preferences.PreferCertificates,
		TargetDate:             req.Preferences.TargetDate,
		PreferredResourceTypes: req.Preferences.PreferredResourceTypes,
		ExcludedProviders:      req.Preferences.ExcludedProviders,
	}

	// Generate learning plan.
	plan := h.recEngine.Generate(profile, jobReqs, prefs)

	WriteSuccess(w, http.StatusOK, plan)
}

// resolveJobRequirements resolves job requirements from a job ID or inline job details.
func (h *AnalysisHandler) resolveJobRequirements(jobID string, inline *types.JobRequirementsInput) (scoring.JobRequirements, error) {
	// Try job ID first.
	if jobID != "" {
		for _, job := range sampleJobs {
			if job.ID == jobID {
				return jobToRequirements(job), nil
			}
		}
	}

	// Fall back to inline job details.
	if inline != nil && len(inline.RequiredSkills) > 0 {
		return scoring.JobRequirements{
			Title:               inline.Title,
			RequiredSkills:      inline.RequiredSkills,
			PreferredSkills:     inline.PreferredSkills,
			MinYearsExperience:  inline.MinYearsExperience,
			RequiredDegreeLevel: inline.RequiredDegreeLevel,
			LocationType:        inline.LocationType,
			Industry:            inline.Industry,
			ExperienceLevel:     inline.ExperienceLevel,
		}, nil
	}

	return scoring.JobRequirements{}, errNoJobSpecified
}

// errNoJobSpecified is returned when no job is specified.
var errNoJobSpecified = &apiError{message: "no job specified"}

type apiError struct{ message string }

func (e *apiError) Error() string { return e.message }

// ─────────────────────────────────────────────────────────────────────────────
// ResourcesHandler
// ─────────────────────────────────────────────────────────────────────────────

// ResourcesHandler handles learning resource search endpoints.
type ResourcesHandler struct {
	catalog []recommend.ResourceEntry
}

// NewResourcesHandler creates a new ResourcesHandler.
func NewResourcesHandler() *ResourcesHandler {
	return &ResourcesHandler{
		catalog: recommend.GetCatalog(),
	}
}

// RegisterRoutes registers resource routes on the mux.
//
//	GET /api/resources/search – search learning resources
func (h *ResourcesHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/resources/search", h.Search)
}

// Search handles GET /api/resources/search.
//
// Query parameters:
//   - skill: filter by skill name
//   - type: filter by resource type
//   - difficulty: filter by difficulty
//   - free: "true" for free resources only
//   - has_certificate: "true" for resources with certificates
//   - has_hands_on: "true" for hands-on resources
//   - min_rating: minimum rating (0.0-5.0)
//   - q: full-text search query
//   - limit: max results (default 20)
//   - offset: pagination offset
func (h *ResourcesHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteMethodNotAllowed(w)
		return
	}

	q := r.URL.Query()
	skill := q.Get("skill")
	resourceType := q.Get("type")
	difficulty := q.Get("difficulty")
	free := q.Get("free") == "true"
	hasCert := q.Get("has_certificate") == "true"
	hasHandsOn := q.Get("has_hands_on") == "true"
	searchQuery := q.Get("q")
	limit := queryParamInt(r, "limit", 20)
	offset := queryParamInt(r, "offset", 0)

	if limit > 100 {
		limit = 100
	}

	// Filter catalog.
	var results []recommend.ResourceEntry
	for _, res := range h.catalog {
		if !matchesResourceFilter(res, skill, resourceType, difficulty, free, hasCert, hasHandsOn, searchQuery) {
			continue
		}
		results = append(results, res)
	}

	total := len(results)

	// Apply pagination.
	if offset >= total {
		results = nil
	} else {
		end := offset + limit
		if end > total {
			end = total
		}
		results = results[offset:end]
	}

	WriteSuccessWithMeta(w, http.StatusOK, results, &types.ResponseMeta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// matchesResourceFilter returns true if a resource matches the search filters.
func matchesResourceFilter(
	res recommend.ResourceEntry,
	skill, resourceType, difficulty string,
	free, hasCert, hasHandsOn bool,
	query string,
) bool {
	if skill != "" {
		skillNorm := strings.ToLower(strings.TrimSpace(skill))
		found := false
		for _, s := range res.Skills {
			if strings.ToLower(s) == skillNorm {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if resourceType != "" && !strings.EqualFold(res.ResourceType, resourceType) {
		return false
	}

	if difficulty != "" && !strings.EqualFold(res.Difficulty, difficulty) &&
		res.Difficulty != "all_levels" {
		return false
	}

	if free && res.CostType != "free" && res.CostType != "free_audit" {
		return false
	}

	if hasCert && !res.HasCertificate {
		return false
	}

	if hasHandsOn && !res.HasHandsOn {
		return false
	}

	if query != "" {
		qLower := strings.ToLower(query)
		if !strings.Contains(strings.ToLower(res.Title), qLower) &&
			!strings.Contains(strings.ToLower(res.Description), qLower) &&
			!strings.Contains(strings.ToLower(res.Provider), qLower) {
			return false
		}
	}

	return true
}
