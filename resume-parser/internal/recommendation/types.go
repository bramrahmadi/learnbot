// Package recommendation implements the training recommendation engine for
// the LearnBot platform. It takes a candidate's skill gap analysis and
// produces a personalized, structured learning plan.
//
// The MVP uses a rule-based approach:
//  1. For each skill gap, query the built-in resource catalog
//  2. Filter resources by user preferences (free/paid, time commitment)
//  3. Rank resources by relevance, quality, and popularity
//  4. Generate a phased learning timeline with milestones
//
// The output is a structured learning plan with:
//   - Phases (Critical → Important → Nice-to-have)
//   - Resources per skill with estimated completion time
//   - Daily/weekly study schedule
//   - Milestones and checkpoints
//   - Alternative resources for flexibility
package recommendation

import "github.com/learnbot/resume-parser/internal/scorer"

// ─────────────────────────────────────────────────────────────────────────────
// User preferences
// ─────────────────────────────────────────────────────────────────────────────

// UserPreferences captures the user's learning preferences for filtering
// and personalizing recommendations.
type UserPreferences struct {
	// PreferFree indicates the user prefers free resources over paid ones.
	PreferFree bool `json:"prefer_free"`

	// MaxBudgetUSD is the maximum budget for paid resources (0 = no limit).
	MaxBudgetUSD float64 `json:"max_budget_usd,omitempty"`

	// WeeklyHoursAvailable is the number of hours per week the user can dedicate
	// to learning (default: 10).
	WeeklyHoursAvailable float64 `json:"weekly_hours_available"`

	// PreferHandsOn indicates the user prefers hands-on/project-based learning.
	PreferHandsOn bool `json:"prefer_hands_on"`

	// PreferCertificates indicates the user prefers resources that offer
	// completion certificates.
	PreferCertificates bool `json:"prefer_certificates"`

	// TargetDate is the optional target date for completing the learning plan
	// (ISO 8601 date string, e.g. "2025-06-01"). Empty = no deadline.
	TargetDate string `json:"target_date,omitempty"`

	// PreferredResourceTypes lists preferred resource types in priority order.
	// Empty = no preference (all types considered).
	// Valid values: "course", "certification", "documentation", "video",
	// "book", "practice", "article", "project"
	PreferredResourceTypes []string `json:"preferred_resource_types,omitempty"`

	// ExcludedProviders lists provider names to exclude from recommendations.
	ExcludedProviders []string `json:"excluded_providers,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource catalog types
// ─────────────────────────────────────────────────────────────────────────────

// ResourceEntry represents a single learning resource in the built-in catalog.
type ResourceEntry struct {
	// ID is a unique identifier for the resource.
	ID string `json:"id"`

	// Title is the resource title.
	Title string `json:"title"`

	// Description is a short description of the resource.
	Description string `json:"description"`

	// URL is the resource URL.
	URL string `json:"url"`

	// Provider is the platform or organization offering the resource.
	Provider string `json:"provider"`

	// ResourceType classifies the resource.
	// Values: "course", "certification", "documentation", "video",
	// "book", "practice", "article", "project"
	ResourceType string `json:"resource_type"`

	// Difficulty is the resource difficulty level.
	// Values: "beginner", "intermediate", "advanced", "expert", "all_levels"
	Difficulty string `json:"difficulty"`

	// CostType is the cost model.
	// Values: "free", "freemium", "paid", "subscription", "free_audit"
	CostType string `json:"cost_type"`

	// CostUSD is the approximate cost in USD (0 for free resources).
	CostUSD float64 `json:"cost_usd,omitempty"`

	// DurationHours is the estimated completion time in hours.
	DurationHours float64 `json:"duration_hours"`

	// DurationLabel is a human-readable duration (e.g. "8 weeks", "40 hours").
	DurationLabel string `json:"duration_label,omitempty"`

	// Skills lists the skills this resource covers (normalized lowercase).
	Skills []string `json:"skills"`

	// PrimarySkill is the main skill taught by this resource.
	PrimarySkill string `json:"primary_skill"`

	// Rating is the average user rating [0.0, 5.0].
	Rating float64 `json:"rating,omitempty"`

	// RatingCount is the number of ratings.
	RatingCount int `json:"rating_count,omitempty"`

	// HasCertificate indicates whether the resource offers a certificate.
	HasCertificate bool `json:"has_certificate"`

	// HasHandsOn indicates whether the resource includes hands-on exercises.
	HasHandsOn bool `json:"has_hands_on"`

	// IsVerified indicates the resource has been curated/verified.
	IsVerified bool `json:"is_verified"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Recommendation output types
// ─────────────────────────────────────────────────────────────────────────────

// RecommendedResource is a resource recommended for a specific skill gap,
// enriched with relevance scoring and context.
type RecommendedResource struct {
	// Resource is the underlying resource entry.
	Resource ResourceEntry `json:"resource"`

	// RelevanceScore is a composite score [0.0, 1.0] indicating how well
	// this resource matches the user's needs for this skill gap.
	// Factors: skill match, difficulty fit, user preferences, quality.
	RelevanceScore float64 `json:"relevance_score"`

	// IsAlternative indicates this is an alternative resource (not the
	// primary recommendation for this skill).
	IsAlternative bool `json:"is_alternative"`

	// RecommendationReason explains why this resource was recommended.
	RecommendationReason string `json:"recommendation_reason"`

	// EstimatedCompletionHours is the estimated hours to complete this
	// resource given the user's current level.
	EstimatedCompletionHours float64 `json:"estimated_completion_hours"`
}

// SkillRecommendation groups resources recommended for a single skill gap.
type SkillRecommendation struct {
	// SkillName is the name of the skill gap being addressed.
	SkillName string `json:"skill_name"`

	// GapCategory is the gap category: "critical", "important", "nice_to_have".
	GapCategory string `json:"gap_category"`

	// PriorityScore is the gap's priority score [0.0, 1.0].
	PriorityScore float64 `json:"priority_score"`

	// PrimaryResource is the top-ranked resource for this skill.
	PrimaryResource *RecommendedResource `json:"primary_resource,omitempty"`

	// AlternativeResources lists up to 2 alternative resources.
	AlternativeResources []RecommendedResource `json:"alternative_resources,omitempty"`

	// EstimatedHoursToJobReady is the estimated hours to reach job-ready
	// proficiency in this skill.
	EstimatedHoursToJobReady int `json:"estimated_hours_to_job_ready"`

	// CurrentLevel is the user's current proficiency (empty if none).
	CurrentLevel string `json:"current_level,omitempty"`

	// TargetLevel is the required proficiency level.
	TargetLevel string `json:"target_level"`
}

// LearningPhase represents a phase in the learning plan.
type LearningPhase struct {
	// PhaseNumber is the phase order (1 = first).
	PhaseNumber int `json:"phase_number"`

	// PhaseName is the phase name (e.g. "Critical Skills", "Preferred Skills").
	PhaseName string `json:"phase_name"`

	// PhaseDescription explains the focus of this phase.
	PhaseDescription string `json:"phase_description"`

	// Skills lists the skill recommendations in this phase, ordered by priority.
	Skills []SkillRecommendation `json:"skills"`

	// TotalHours is the total estimated hours for this phase.
	TotalHours float64 `json:"total_hours"`

	// EstimatedWeeks is the estimated weeks to complete this phase
	// given the user's weekly hours available.
	EstimatedWeeks float64 `json:"estimated_weeks"`

	// Milestone is the achievement unlocked upon completing this phase.
	Milestone string `json:"milestone"`
}

// WeeklySchedule represents a suggested weekly study schedule.
type WeeklySchedule struct {
	// WeekNumber is the week number in the plan (1-based).
	WeekNumber int `json:"week_number"`

	// PhaseNumber is the phase this week belongs to.
	PhaseNumber int `json:"phase_number"`

	// SkillFocus is the primary skill to focus on this week.
	SkillFocus string `json:"skill_focus"`

	// ResourceTitle is the resource to work on this week.
	ResourceTitle string `json:"resource_title"`

	// HoursPlanned is the planned study hours for this week.
	HoursPlanned float64 `json:"hours_planned"`

	// CumulativeHours is the running total of hours up to and including this week.
	CumulativeHours float64 `json:"cumulative_hours"`

	// Activities lists suggested activities for this week.
	Activities []string `json:"activities"`

	// IsCheckpoint indicates this week has a milestone checkpoint.
	IsCheckpoint bool `json:"is_checkpoint"`

	// CheckpointDescription describes the checkpoint goal.
	CheckpointDescription string `json:"checkpoint_description,omitempty"`
}

// LearningTimeline provides a week-by-week study schedule.
type LearningTimeline struct {
	// TotalWeeks is the total number of weeks in the plan.
	TotalWeeks int `json:"total_weeks"`

	// TotalHours is the total estimated study hours.
	TotalHours float64 `json:"total_hours"`

	// WeeklyHours is the planned weekly study hours.
	WeeklyHours float64 `json:"weekly_hours"`

	// Weeks is the week-by-week schedule.
	Weeks []WeeklySchedule `json:"weeks"`

	// TargetCompletionDate is the estimated completion date (ISO 8601).
	// Empty if no target date was specified.
	TargetCompletionDate string `json:"target_completion_date,omitempty"`
}

// LearningPlan is the complete personalized learning plan output.
type LearningPlan struct {
	// JobTitle is the target job title.
	JobTitle string `json:"job_title"`

	// ReadinessScore is the current readiness score [0, 100].
	ReadinessScore float64 `json:"readiness_score"`

	// TotalGaps is the total number of skill gaps identified.
	TotalGaps int `json:"total_gaps"`

	// TotalEstimatedHours is the total estimated learning hours.
	TotalEstimatedHours float64 `json:"total_estimated_hours"`

	// Phases is the ordered list of learning phases.
	// Phase 1: Critical gaps, Phase 2: Important gaps, Phase 3: Nice-to-have.
	Phases []LearningPhase `json:"phases"`

	// Timeline provides the week-by-week schedule.
	Timeline LearningTimeline `json:"timeline"`

	// MatchedSkills lists skills the candidate already has.
	MatchedSkills []string `json:"matched_skills"`

	// Summary is a human-readable summary of the learning plan.
	Summary LearningPlanSummary `json:"summary"`
}

// LearningPlanSummary provides a high-level overview of the learning plan.
type LearningPlanSummary struct {
	// Headline is a one-line summary (e.g. "3 critical gaps to close in 8 weeks").
	Headline string `json:"headline"`

	// CriticalGapCount is the number of critical gaps.
	CriticalGapCount int `json:"critical_gap_count"`

	// ImportantGapCount is the number of important gaps.
	ImportantGapCount int `json:"important_gap_count"`

	// FreeResourceCount is the number of free resources recommended.
	FreeResourceCount int `json:"free_resource_count"`

	// PaidResourceCount is the number of paid resources recommended.
	PaidResourceCount int `json:"paid_resource_count"`

	// EstimatedTotalCostUSD is the estimated total cost of paid resources.
	EstimatedTotalCostUSD float64 `json:"estimated_total_cost_usd"`

	// TopSkillsToLearn lists the top 3 skills to focus on first.
	TopSkillsToLearn []string `json:"top_skills_to_learn"`

	// QuickWins lists skills that can be learned quickly (< 20 hours).
	QuickWins []string `json:"quick_wins"`
}

// ─────────────────────────────────────────────────────────────────────────────
// API request/response types
// ─────────────────────────────────────────────────────────────────────────────

// RecommendationRequest is the input to the recommendation API endpoint.
type RecommendationRequest struct {
	// Profile is the candidate's professional profile.
	Profile scorer.CandidateProfile `json:"profile"`

	// Job is the job requirements to generate recommendations for.
	Job scorer.JobRequirements `json:"job"`

	// Preferences are the user's learning preferences.
	Preferences UserPreferences `json:"preferences"`
}

// RecommendationResponse is the output of the recommendation API endpoint.
type RecommendationResponse struct {
	// Success indicates whether the recommendation succeeded.
	Success bool `json:"success"`

	// Data contains the learning plan when Success is true.
	Data *LearningPlan `json:"data,omitempty"`

	// Error contains an error message when Success is false.
	Error string `json:"error,omitempty"`
}
