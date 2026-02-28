// Package handler – profile.go implements user profile endpoints.
package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/api-gateway/internal/types"
)

// ─────────────────────────────────────────────────────────────────────────────
// In-memory profile store (MVP)
// ─────────────────────────────────────────────────────────────────────────────

// profileRecord stores a user's profile.
type profileRecord struct {
	UserID            string
	Headline          string
	Summary           string
	LocationCity      string
	LocationCountry   string
	LinkedInURL       string
	GitHubURL         string
	WebsiteURL        string
	YearsOfExperience float64
	IsOpenToWork      bool
	Skills            []skillRecord
	UpdatedAt         time.Time
}

// skillRecord stores a single skill.
type skillRecord struct {
	Name              string
	Proficiency       string
	YearsOfExperience float64
	IsPrimary         bool
}

// profileStore is a thread-safe in-memory profile store.
type profileStore struct {
	mu       sync.RWMutex
	profiles map[string]*profileRecord // keyed by userID
}

var globalProfileStore = &profileStore{
	profiles: make(map[string]*profileRecord),
}

func (s *profileStore) get(userID string) *profileRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.profiles[userID]
	if !ok {
		return &profileRecord{UserID: userID}
	}
	return p
}

func (s *profileStore) update(userID string, fn func(*profileRecord)) *profileRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.profiles[userID]
	if !ok {
		p = &profileRecord{UserID: userID}
	}
	fn(p)
	p.UpdatedAt = time.Now()
	s.profiles[userID] = p
	return p
}

// ─────────────────────────────────────────────────────────────────────────────
// ProfileHandler
// ─────────────────────────────────────────────────────────────────────────────

// ProfileHandler handles user profile endpoints.
type ProfileHandler struct {
	jwtCfg middleware.JWTConfig
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(jwtCfg middleware.JWTConfig) *ProfileHandler {
	return &ProfileHandler{jwtCfg: jwtCfg}
}

// RegisterRoutes registers profile routes on the mux.
//
//	GET  /api/users/profile    – get current user's profile
//	PUT  /api/users/profile    – update current user's profile
//	GET  /api/profile/skills   – get current user's skills
//	PUT  /api/profile/skills   – update current user's skills
func (h *ProfileHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/api/users/profile",
		authMiddleware(http.HandlerFunc(h.handleProfile)))
	mux.Handle("/api/profile/skills",
		authMiddleware(http.HandlerFunc(h.handleSkills)))
}

// handleProfile handles GET/PUT /api/users/profile.
func (h *ProfileHandler) handleProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProfile(w, r, userID)
	case http.MethodPut:
		h.updateProfile(w, r, userID)
	default:
		WriteMethodNotAllowed(w)
	}
}

// getProfile handles GET /api/users/profile.
func (h *ProfileHandler) getProfile(w http.ResponseWriter, r *http.Request, userID string) {
	user, exists := globalUserStore.findByID(userID)
	if !exists {
		WriteNotFound(w, "user")
		return
	}

	profile := globalProfileStore.get(userID)

	WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"user_id":             userID,
		"email":               user.Email,
		"full_name":           user.FullName,
		"headline":            profile.Headline,
		"summary":             profile.Summary,
		"location_city":       profile.LocationCity,
		"location_country":    profile.LocationCountry,
		"linkedin_url":        profile.LinkedInURL,
		"github_url":          profile.GitHubURL,
		"website_url":         profile.WebsiteURL,
		"years_of_experience": profile.YearsOfExperience,
		"is_open_to_work":     profile.IsOpenToWork,
		"skills":              profile.Skills,
		"updated_at":          profile.UpdatedAt,
	})
}

// updateProfile handles PUT /api/users/profile.
func (h *ProfileHandler) updateProfile(w http.ResponseWriter, r *http.Request, userID string) {
	var req types.ProfileUpdateRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	profile := globalProfileStore.update(userID, func(p *profileRecord) {
		if req.Headline != nil {
			p.Headline = *req.Headline
		}
		if req.Summary != nil {
			p.Summary = *req.Summary
		}
		if req.LocationCity != nil {
			p.LocationCity = *req.LocationCity
		}
		if req.LocationCountry != nil {
			p.LocationCountry = *req.LocationCountry
		}
		if req.LinkedInURL != nil {
			p.LinkedInURL = *req.LinkedInURL
		}
		if req.GitHubURL != nil {
			p.GitHubURL = *req.GitHubURL
		}
		if req.WebsiteURL != nil {
			p.WebsiteURL = *req.WebsiteURL
		}
		if req.YearsOfExperience != nil {
			p.YearsOfExperience = *req.YearsOfExperience
		}
		if req.IsOpenToWork != nil {
			p.IsOpenToWork = *req.IsOpenToWork
		}
	})

	WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"user_id":             userID,
		"headline":            profile.Headline,
		"summary":             profile.Summary,
		"location_city":       profile.LocationCity,
		"location_country":    profile.LocationCountry,
		"linkedin_url":        profile.LinkedInURL,
		"github_url":          profile.GitHubURL,
		"website_url":         profile.WebsiteURL,
		"years_of_experience": profile.YearsOfExperience,
		"is_open_to_work":     profile.IsOpenToWork,
		"updated_at":          profile.UpdatedAt,
	})
}

// handleSkills handles GET/PUT /api/profile/skills.
func (h *ProfileHandler) handleSkills(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getSkills(w, r, userID)
	case http.MethodPut:
		h.updateSkills(w, r, userID)
	default:
		WriteMethodNotAllowed(w)
	}
}

// getSkills handles GET /api/profile/skills.
func (h *ProfileHandler) getSkills(w http.ResponseWriter, r *http.Request, userID string) {
	profile := globalProfileStore.get(userID)
	WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"skills":  profile.Skills,
		"count":   len(profile.Skills),
	})
}

// updateSkills handles PUT /api/profile/skills.
func (h *ProfileHandler) updateSkills(w http.ResponseWriter, r *http.Request, userID string) {
	var req types.SkillUpdateRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	// Validate skills.
	var v Validator
	for i, skill := range req.Skills {
		if trimSpaceStr(skill.Name) == "" {
			v.errors = append(v.errors, types.FieldError{
				Field:   "skills[" + itoa(i) + "].name",
				Message: "skill name is required",
			})
		}
		validProficiencies := map[string]bool{
			"beginner": true, "intermediate": true, "advanced": true, "expert": true,
		}
		if skill.Proficiency != "" && !validProficiencies[skill.Proficiency] {
			v.errors = append(v.errors, types.FieldError{
				Field:   "skills[" + itoa(i) + "].proficiency",
				Message: "must be one of: beginner, intermediate, advanced, expert",
			})
		}
	}
	if v.WriteIfInvalid(w) {
		return
	}

	profile := globalProfileStore.update(userID, func(p *profileRecord) {
		p.Skills = make([]skillRecord, len(req.Skills))
		for i, s := range req.Skills {
			yoe := 0.0
			if s.YearsOfExperience != nil {
				yoe = *s.YearsOfExperience
			}
			p.Skills[i] = skillRecord{
				Name:              s.Name,
				Proficiency:       s.Proficiency,
				YearsOfExperience: yoe,
				IsPrimary:         s.IsPrimary,
			}
		}
	})

	WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"skills":  profile.Skills,
		"count":   len(profile.Skills),
	})
}

// itoa converts an int to a string.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
