// Package types defines shared request/response types for the API gateway.
package types

import "time"

// ─────────────────────────────────────────────────────────────────────────────
// Standard API response envelope
// ─────────────────────────────────────────────────────────────────────────────

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	// Success indicates whether the request succeeded.
	Success bool `json:"success"`

	// Data contains the response payload when Success is true.
	Data interface{} `json:"data,omitempty"`

	// Error contains error details when Success is false.
	Error *APIError `json:"error,omitempty"`

	// Meta contains optional metadata (pagination, etc.).
	Meta *ResponseMeta `json:"meta,omitempty"`
}

// APIError represents a structured error response.
type APIError struct {
	// Code is a machine-readable error code (e.g. "INVALID_INPUT").
	Code string `json:"code"`

	// Message is a human-readable error description.
	Message string `json:"message"`

	// Details contains field-level validation errors.
	Details []FieldError `json:"details,omitempty"`
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	// Field is the field name that failed validation.
	Field string `json:"field"`

	// Message describes the validation failure.
	Message string `json:"message"`
}

// ResponseMeta contains pagination and other metadata.
type ResponseMeta struct {
	// Total is the total number of items (for paginated responses).
	Total int `json:"total,omitempty"`

	// Limit is the page size.
	Limit int `json:"limit,omitempty"`

	// Offset is the pagination offset.
	Offset int `json:"offset,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Authentication types
// ─────────────────────────────────────────────────────────────────────────────

// RegisterRequest is the input for user registration.
type RegisterRequest struct {
	// Email is the user's email address.
	Email string `json:"email"`

	// Password is the user's password (min 8 characters).
	Password string `json:"password"`

	// FullName is the user's display name.
	FullName string `json:"full_name"`
}

// LoginRequest is the input for user login.
type LoginRequest struct {
	// Email is the user's email address.
	Email string `json:"email"`

	// Password is the user's password.
	Password string `json:"password"`
}

// AuthResponse is returned on successful authentication.
type AuthResponse struct {
	// Token is the JWT access token.
	Token string `json:"token"`

	// ExpiresAt is the token expiration time.
	ExpiresAt time.Time `json:"expires_at"`

	// User contains basic user information.
	User UserInfo `json:"user"`
}

// UserInfo contains basic user information included in auth responses.
type UserInfo struct {
	// ID is the user's unique identifier.
	ID string `json:"id"`

	// Email is the user's email address.
	Email string `json:"email"`

	// FullName is the user's display name.
	FullName string `json:"full_name"`
}

// JWTClaims represents the claims stored in a JWT token.
type JWTClaims struct {
	// UserID is the authenticated user's ID.
	UserID string `json:"user_id"`

	// Email is the authenticated user's email.
	Email string `json:"email"`

	// IsAdmin indicates whether the user has admin privileges.
	IsAdmin bool `json:"is_admin"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Profile types
// ─────────────────────────────────────────────────────────────────────────────

// ProfileUpdateRequest is the input for updating a user profile.
type ProfileUpdateRequest struct {
	Headline          *string  `json:"headline,omitempty"`
	Summary           *string  `json:"summary,omitempty"`
	LocationCity      *string  `json:"location_city,omitempty"`
	LocationCountry   *string  `json:"location_country,omitempty"`
	LinkedInURL       *string  `json:"linkedin_url,omitempty"`
	GitHubURL         *string  `json:"github_url,omitempty"`
	WebsiteURL        *string  `json:"website_url,omitempty"`
	YearsOfExperience *float64 `json:"years_of_experience,omitempty"`
	IsOpenToWork      *bool    `json:"is_open_to_work,omitempty"`
}

// SkillUpdateRequest is the input for updating user skills.
type SkillUpdateRequest struct {
	// Skills is the list of skills to upsert.
	Skills []SkillInput `json:"skills"`
}

// SkillInput represents a single skill to add or update.
type SkillInput struct {
	Name              string   `json:"name"`
	Proficiency       string   `json:"proficiency"`
	YearsOfExperience *float64 `json:"years_of_experience,omitempty"`
	IsPrimary         bool     `json:"is_primary"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Job matching types
// ─────────────────────────────────────────────────────────────────────────────

// JobSearchRequest is the input for job search.
type JobSearchRequest struct {
	// Query is the free-text search query.
	Query string `json:"query,omitempty"`

	// Skills filters jobs requiring specific skills.
	Skills []string `json:"skills,omitempty"`

	// LocationType filters by work arrangement: "remote", "hybrid", "on_site".
	LocationType string `json:"location_type,omitempty"`

	// ExperienceLevel filters by seniority: "entry", "mid", "senior", etc.
	ExperienceLevel string `json:"experience_level,omitempty"`

	// Industry filters by industry.
	Industry string `json:"industry,omitempty"`

	// Limit is the maximum number of results (default 20).
	Limit int `json:"limit,omitempty"`

	// Offset is the pagination offset.
	Offset int `json:"offset,omitempty"`
}

// JobSummary is a lightweight job representation for list responses.
type JobSummary struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Company         string   `json:"company"`
	LocationType    string   `json:"location_type"`
	ExperienceLevel string   `json:"experience_level"`
	RequiredSkills  []string `json:"required_skills"`
	PostedAt        string   `json:"posted_at,omitempty"`
	MatchScore      *float64 `json:"match_score,omitempty"`
}

// JobDetail is the full job representation.
type JobDetail struct {
	JobSummary
	Description     string   `json:"description"`
	PreferredSkills []string `json:"preferred_skills,omitempty"`
	MinExperience   float64  `json:"min_years_experience"`
	Industry        string   `json:"industry,omitempty"`
	LocationCity    string   `json:"location_city,omitempty"`
	LocationCountry string   `json:"location_country,omitempty"`
	SalaryMin       *int     `json:"salary_min,omitempty"`
	SalaryMax       *int     `json:"salary_max,omitempty"`
	SalaryCurrency  string   `json:"salary_currency,omitempty"`
	ApplyURL        string   `json:"apply_url,omitempty"`
}

// JobMatchResponse is the response for job match scoring.
type JobMatchResponse struct {
	JobID          string  `json:"job_id"`
	OverallScore   float64 `json:"overall_score"`
	SkillMatch     float64 `json:"skill_match"`
	ExperienceMatch float64 `json:"experience_match"`
	EducationMatch float64 `json:"education_match"`
	LocationFit    float64 `json:"location_fit"`
	IndustryMatch  float64 `json:"industry_match"`
	MatchedSkills  []string `json:"matched_skills"`
	MissingSkills  []string `json:"missing_skills"`
	Recommendation string  `json:"recommendation"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Gap analysis types
// ─────────────────────────────────────────────────────────────────────────────

// GapAnalysisRequest is the input for gap analysis.
type GapAnalysisRequest struct {
	// JobID is the ID of the target job (optional if job details provided).
	JobID string `json:"job_id,omitempty"`

	// Job contains inline job requirements (used when JobID is not provided).
	Job *JobRequirementsInput `json:"job,omitempty"`
}

// JobRequirementsInput allows inline job requirements for gap analysis.
type JobRequirementsInput struct {
	Title               string   `json:"title"`
	RequiredSkills      []string `json:"required_skills"`
	PreferredSkills     []string `json:"preferred_skills,omitempty"`
	MinYearsExperience  float64  `json:"min_years_experience"`
	RequiredDegreeLevel string   `json:"required_degree_level,omitempty"`
	LocationType        string   `json:"location_type,omitempty"`
	Industry            string   `json:"industry,omitempty"`
	ExperienceLevel     string   `json:"experience_level,omitempty"`
}

// TrainingRecommendationRequest is the input for training recommendations.
type TrainingRecommendationRequest struct {
	// JobID is the target job ID (optional).
	JobID string `json:"job_id,omitempty"`

	// Job contains inline job requirements.
	Job *JobRequirementsInput `json:"job,omitempty"`

	// Preferences are the user's learning preferences.
	Preferences LearningPreferencesInput `json:"preferences"`
}

// LearningPreferencesInput captures user learning preferences.
type LearningPreferencesInput struct {
	PreferFree             bool     `json:"prefer_free"`
	MaxBudgetUSD           float64  `json:"max_budget_usd,omitempty"`
	WeeklyHoursAvailable   float64  `json:"weekly_hours_available"`
	PreferHandsOn          bool     `json:"prefer_hands_on"`
	PreferCertificates     bool     `json:"prefer_certificates"`
	TargetDate             string   `json:"target_date,omitempty"`
	PreferredResourceTypes []string `json:"preferred_resource_types,omitempty"`
	ExcludedProviders      []string `json:"excluded_providers,omitempty"`
}

// ResourceSearchRequest is the input for resource search.
type ResourceSearchRequest struct {
	Skill          string `json:"skill,omitempty"`
	ResourceType   string `json:"resource_type,omitempty"`
	Difficulty     string `json:"difficulty,omitempty"`
	Free           bool   `json:"free,omitempty"`
	HasCertificate bool   `json:"has_certificate,omitempty"`
	HasHandsOn     bool   `json:"has_hands_on,omitempty"`
	MinRating      float64 `json:"min_rating,omitempty"`
	Query          string `json:"q,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
}
