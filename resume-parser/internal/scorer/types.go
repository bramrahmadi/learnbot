// Package scorer implements the acceptance likelihood scoring algorithm
// for matching user profiles against job postings.
//
// The algorithm uses a weighted formula from the LearnBot specification:
//
//	Score = (skill_match * 0.35) + (experience_match * 0.25) +
//	        (education_match * 0.15) + (location_fit * 0.10) +
//	        (industry_relevance * 0.15)
//
// Each component returns a value in [0.0, 1.0], and the final score
// is expressed as a percentage in [0.0, 100.0].
package scorer

// Weights defines the contribution of each scoring component.
// They must sum to 1.0.
const (
	WeightSkillMatch        = 0.35
	WeightExperienceMatch   = 0.25
	WeightEducationMatch    = 0.15
	WeightLocationFit       = 0.10
	WeightIndustryRelevance = 0.15
)

// JobRequirements describes the requirements extracted from a job posting.
type JobRequirements struct {
	// Title is the job title (e.g. "Senior Software Engineer").
	Title string `json:"title"`

	// RequiredSkills is the list of must-have skills.
	RequiredSkills []string `json:"required_skills"`

	// PreferredSkills is the list of nice-to-have skills.
	PreferredSkills []string `json:"preferred_skills,omitempty"`

	// MinYearsExperience is the minimum years of experience required.
	MinYearsExperience float64 `json:"min_years_experience"`

	// MaxYearsExperience is the maximum years of experience (0 = no upper limit).
	MaxYearsExperience float64 `json:"max_years_experience,omitempty"`

	// RequiredDegreeLevel is the minimum degree level required.
	// Accepted values: "high_school", "associate", "bachelor", "master",
	// "doctorate", "professional", "certificate", "diploma", "other", "".
	// Empty string means no degree requirement.
	RequiredDegreeLevel string `json:"required_degree_level,omitempty"`

	// PreferredFields lists preferred fields of study (e.g. "Computer Science").
	PreferredFields []string `json:"preferred_fields,omitempty"`

	// LocationCity is the job's city.
	LocationCity string `json:"location_city,omitempty"`

	// LocationCountry is the job's country.
	LocationCountry string `json:"location_country,omitempty"`

	// LocationType is the work arrangement: "remote", "hybrid", "on_site".
	LocationType string `json:"location_type"`

	// Industry is the job's industry (e.g. "Software", "Finance").
	Industry string `json:"industry,omitempty"`

	// RelatedIndustries lists industries considered equivalent or transferable.
	RelatedIndustries []string `json:"related_industries,omitempty"`

	// ExperienceLevel is the seniority level: "internship", "entry", "mid",
	// "senior", "lead", "executive".
	ExperienceLevel string `json:"experience_level,omitempty"`
}

// CandidateProfile describes the user's professional profile used for scoring.
type CandidateProfile struct {
	// Skills is the list of skills the candidate possesses.
	Skills []CandidateSkill `json:"skills"`

	// YearsOfExperience is the total years of professional experience.
	YearsOfExperience float64 `json:"years_of_experience"`

	// WorkHistory is the list of past and current positions.
	WorkHistory []WorkHistoryEntry `json:"work_history"`

	// Education is the list of educational qualifications.
	Education []EducationEntry `json:"education"`

	// LocationCity is the candidate's current city.
	LocationCity string `json:"location_city,omitempty"`

	// LocationCountry is the candidate's current country.
	LocationCountry string `json:"location_country,omitempty"`

	// WillingToRelocate indicates whether the candidate will relocate.
	WillingToRelocate bool `json:"willing_to_relocate"`

	// RemotePreference indicates the candidate's preferred work arrangement.
	// Accepted values: "remote", "hybrid", "on_site", "any".
	RemotePreference string `json:"remote_preference,omitempty"`
}

// CandidateSkill represents a single skill with optional proficiency metadata.
type CandidateSkill struct {
	// Name is the skill name (e.g. "Go", "Python").
	Name string `json:"name"`

	// Proficiency is the self-assessed level: "beginner", "intermediate",
	// "advanced", "expert".
	Proficiency string `json:"proficiency,omitempty"`

	// YearsOfExperience is the number of years using this skill.
	YearsOfExperience float64 `json:"years_of_experience,omitempty"`
}

// WorkHistoryEntry represents a single job in the candidate's work history.
type WorkHistoryEntry struct {
	// Title is the job title held.
	Title string `json:"title"`

	// Industry is the industry of the employer.
	Industry string `json:"industry,omitempty"`

	// DurationMonths is the length of the role in months.
	DurationMonths int `json:"duration_months"`

	// IsCurrent indicates whether this is the candidate's current role.
	IsCurrent bool `json:"is_current"`
}

// EducationEntry represents a single educational qualification.
type EducationEntry struct {
	// DegreeLevel is the level of the degree (e.g. "bachelor", "master").
	DegreeLevel string `json:"degree_level"`

	// FieldOfStudy is the field or major (e.g. "Computer Science").
	FieldOfStudy string `json:"field_of_study,omitempty"`
}

// ScoreBreakdown holds the individual component scores and the final result.
type ScoreBreakdown struct {
	// OverallScore is the final weighted score expressed as a percentage [0, 100].
	OverallScore float64 `json:"overall_score"`

	// SkillMatchScore is the skill match component score [0, 1].
	SkillMatchScore float64 `json:"skill_match_score"`

	// ExperienceMatchScore is the experience match component score [0, 1].
	ExperienceMatchScore float64 `json:"experience_match_score"`

	// EducationMatchScore is the education match component score [0, 1].
	EducationMatchScore float64 `json:"education_match_score"`

	// LocationFitScore is the location fit component score [0, 1].
	LocationFitScore float64 `json:"location_fit_score"`

	// IndustryRelevanceScore is the industry relevance component score [0, 1].
	IndustryRelevanceScore float64 `json:"industry_relevance_score"`

	// MatchedRequiredSkills lists required skills the candidate has.
	MatchedRequiredSkills []string `json:"matched_required_skills"`

	// MissingRequiredSkills lists required skills the candidate lacks.
	MissingRequiredSkills []string `json:"missing_required_skills"`

	// MatchedPreferredSkills lists preferred skills the candidate has.
	MatchedPreferredSkills []string `json:"matched_preferred_skills,omitempty"`
}

// ScoreRequest is the input to the scoring API endpoint.
type ScoreRequest struct {
	// Profile is the candidate's professional profile.
	Profile CandidateProfile `json:"profile"`

	// Job is the job requirements to score against.
	Job JobRequirements `json:"job"`
}

// ScoreResponse is the output of the scoring API endpoint.
type ScoreResponse struct {
	// Success indicates whether the scoring succeeded.
	Success bool `json:"success"`

	// Data contains the score breakdown when Success is true.
	Data *ScoreBreakdown `json:"data,omitempty"`

	// Error contains an error message when Success is false.
	Error string `json:"error,omitempty"`
}
