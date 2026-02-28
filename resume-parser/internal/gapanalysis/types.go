// Package gapanalysis implements the skill gap analysis engine for the
// LearnBot platform. It compares a candidate's current skills against job
// requirements and produces a prioritized list of skill gaps with actionable
// recommendations.
//
// Gap categories (from LearnBot specification):
//   - Critical: must-have requirements the candidate is missing
//   - Important: preferred requirements the candidate is missing
//   - NiceToHave: optional skills the candidate is missing
//
// Gaps are ranked by:
//   - Importance to job acceptance
//   - Estimated time to acquire (learning hours)
//   - Transferability to other roles
package gapanalysis

import "github.com/learnbot/resume-parser/internal/scorer"

// ─────────────────────────────────────────────────────────────────────────────
// Gap category and priority types
// ─────────────────────────────────────────────────────────────────────────────

// GapCategory classifies a skill gap by its importance to the job.
type GapCategory string

const (
	// GapCategoryCritical represents must-have skills that are missing.
	// These are required skills listed in the job's RequiredSkills.
	GapCategoryCritical GapCategory = "critical"

	// GapCategoryImportant represents preferred skills that are missing.
	// These are skills listed in the job's PreferredSkills.
	GapCategoryImportant GapCategory = "important"

	// GapCategoryNiceToHave represents optional skills that are missing.
	// These are skills inferred from the job context but not explicitly required.
	GapCategoryNiceToHave GapCategory = "nice_to_have"
)

// DifficultyLevel classifies how hard a skill is to acquire.
type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
	DifficultyExpert       DifficultyLevel = "expert"
)

// ─────────────────────────────────────────────────────────────────────────────
// Core gap types
// ─────────────────────────────────────────────────────────────────────────────

// SkillGap represents a single missing skill with its analysis metadata.
type SkillGap struct {
	// SkillName is the name of the missing skill.
	SkillName string `json:"skill_name"`

	// Category classifies the gap as critical, important, or nice_to_have.
	Category GapCategory `json:"category"`

	// PriorityScore is a composite score [0.0, 1.0] ranking the gap by:
	//   - Importance to job acceptance (weight: 0.50)
	//   - Transferability to other roles (weight: 0.30)
	//   - Inverse of time to acquire (weight: 0.20)
	// Higher score = higher priority to address.
	PriorityScore float64 `json:"priority_score"`

	// ImportanceScore is the raw importance weight [0.0, 1.0] based on
	// how critical the skill is to the job (critical=1.0, important=0.6,
	// nice_to_have=0.3).
	ImportanceScore float64 `json:"importance_score"`

	// EstimatedLearningHours is the estimated number of hours to acquire
	// this skill to a job-ready level, given the candidate's current level.
	EstimatedLearningHours int `json:"estimated_learning_hours"`

	// TransferabilityScore is a score [0.0, 1.0] indicating how broadly
	// useful this skill is across different roles and industries.
	// Higher = more transferable (e.g. Python > niche framework).
	TransferabilityScore float64 `json:"transferability_score"`

	// CurrentLevel is the candidate's current proficiency in this skill,
	// if they have partial knowledge. Empty if they have no knowledge.
	CurrentLevel string `json:"current_level,omitempty"`

	// TargetLevel is the proficiency level required by the job.
	TargetLevel string `json:"target_level"`

	// SemanticSimilarityScore is the similarity [0.0, 1.0] between this
	// missing skill and the candidate's closest existing skill. A high
	// score means the candidate has a related skill that reduces the gap.
	SemanticSimilarityScore float64 `json:"semantic_similarity_score"`

	// ClosestExistingSkill is the name of the candidate's skill that is
	// most semantically similar to this gap skill.
	ClosestExistingSkill string `json:"closest_existing_skill,omitempty"`

	// Recommendations contains actionable steps to address this gap.
	Recommendations []Recommendation `json:"recommendations"`

	// Difficulty is the estimated difficulty level to acquire this skill.
	Difficulty DifficultyLevel `json:"difficulty"`
}

// Recommendation is a single actionable step to address a skill gap.
type Recommendation struct {
	// Title is a short description of the recommended action.
	Title string `json:"title"`

	// Description provides more detail about the recommendation.
	Description string `json:"description"`

	// ResourceType classifies the resource: "course", "certification",
	// "project", "book", "documentation", "practice".
	ResourceType string `json:"resource_type"`

	// EstimatedHours is the estimated time investment for this resource.
	EstimatedHours int `json:"estimated_hours"`

	// Priority is the order in which this recommendation should be pursued
	// (1 = highest priority).
	Priority int `json:"priority"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Analysis result types
// ─────────────────────────────────────────────────────────────────────────────

// GapAnalysisResult is the complete output of the gap analysis engine.
type GapAnalysisResult struct {
	// CriticalGaps lists must-have skills the candidate is missing.
	// Sorted by PriorityScore descending.
	CriticalGaps []SkillGap `json:"critical_gaps"`

	// ImportantGaps lists preferred skills the candidate is missing.
	// Sorted by PriorityScore descending.
	ImportantGaps []SkillGap `json:"important_gaps"`

	// NiceToHaveGaps lists optional skills the candidate is missing.
	// Sorted by PriorityScore descending.
	NiceToHaveGaps []SkillGap `json:"nice_to_have_gaps"`

	// TotalGaps is the total number of identified gaps across all categories.
	TotalGaps int `json:"total_gaps"`

	// CriticalGapCount is the number of critical gaps.
	CriticalGapCount int `json:"critical_gap_count"`

	// ImportantGapCount is the number of important gaps.
	ImportantGapCount int `json:"important_gap_count"`

	// NiceToHaveGapCount is the number of nice-to-have gaps.
	NiceToHaveGapCount int `json:"nice_to_have_gap_count"`

	// TotalEstimatedLearningHours is the sum of estimated learning hours
	// across all gaps.
	TotalEstimatedLearningHours int `json:"total_estimated_learning_hours"`

	// ReadinessScore is a score [0.0, 100.0] indicating how ready the
	// candidate is for the role. Derived from the inverse of gap severity.
	ReadinessScore float64 `json:"readiness_score"`

	// TopPriorityGaps is a ranked list of the top 5 gaps to address first,
	// selected across all categories by PriorityScore.
	TopPriorityGaps []SkillGap `json:"top_priority_gaps"`

	// MatchedSkills lists skills the candidate already has that match
	// the job requirements.
	MatchedSkills []string `json:"matched_skills"`

	// VisualData provides a JSON-friendly representation for frontend
	// visualization of the gap analysis.
	VisualData GapVisualData `json:"visual_data"`
}

// GapVisualData provides structured data for frontend visualization.
type GapVisualData struct {
	// RadarChart provides data for a radar/spider chart showing skill coverage.
	RadarChart RadarChartData `json:"radar_chart"`

	// GapsByCategory provides gap counts and hours by category for bar charts.
	GapsByCategory []CategorySummary `json:"gaps_by_category"`

	// LearningTimeline provides a suggested learning order with cumulative hours.
	LearningTimeline []TimelineEntry `json:"learning_timeline"`
}

// RadarChartData provides data for a radar chart visualization.
type RadarChartData struct {
	// Labels are the skill category names.
	Labels []string `json:"labels"`

	// CandidateScores are the candidate's scores per category [0.0, 1.0].
	CandidateScores []float64 `json:"candidate_scores"`

	// RequiredScores are the required scores per category [0.0, 1.0].
	RequiredScores []float64 `json:"required_scores"`
}

// CategorySummary summarizes gaps for a single category.
type CategorySummary struct {
	// Category is the gap category name.
	Category string `json:"category"`

	// GapCount is the number of gaps in this category.
	GapCount int `json:"gap_count"`

	// TotalLearningHours is the total estimated learning hours for this category.
	TotalLearningHours int `json:"total_learning_hours"`

	// AveragePriority is the average priority score for gaps in this category.
	AveragePriority float64 `json:"average_priority"`
}

// TimelineEntry represents a single step in the suggested learning timeline.
type TimelineEntry struct {
	// Order is the suggested learning order (1 = first).
	Order int `json:"order"`

	// SkillName is the skill to learn.
	SkillName string `json:"skill_name"`

	// Category is the gap category.
	Category string `json:"category"`

	// EstimatedHours is the estimated hours for this skill.
	EstimatedHours int `json:"estimated_hours"`

	// CumulativeHours is the running total of hours up to and including this step.
	CumulativeHours int `json:"cumulative_hours"`

	// Rationale explains why this skill is prioritized at this position.
	Rationale string `json:"rationale"`
}

// ─────────────────────────────────────────────────────────────────────────────
// API request/response types
// ─────────────────────────────────────────────────────────────────────────────

// GapAnalysisRequest is the input to the gap analysis API endpoint.
type GapAnalysisRequest struct {
	// Profile is the candidate's professional profile.
	Profile scorer.CandidateProfile `json:"profile"`

	// Job is the job requirements to analyze gaps against.
	Job scorer.JobRequirements `json:"job"`
}

// GapAnalysisResponse is the output of the gap analysis API endpoint.
type GapAnalysisResponse struct {
	// Success indicates whether the analysis succeeded.
	Success bool `json:"success"`

	// Data contains the gap analysis result when Success is true.
	Data *GapAnalysisResult `json:"data,omitempty"`

	// Error contains an error message when Success is false.
	Error string `json:"error,omitempty"`
}
