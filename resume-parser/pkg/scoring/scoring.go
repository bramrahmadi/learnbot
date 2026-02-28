// Package scoring provides a public API for acceptance likelihood scoring.
// This package wraps the internal scorer package for use by external modules.
package scoring

import (
	"github.com/learnbot/resume-parser/internal/scorer"
)

// Re-export types from scorer for external use.

// CandidateProfile describes the user's professional profile.
type CandidateProfile = scorer.CandidateProfile

// CandidateSkill represents a single skill with optional proficiency metadata.
type CandidateSkill = scorer.CandidateSkill

// WorkHistoryEntry represents a single job in the candidate's work history.
type WorkHistoryEntry = scorer.WorkHistoryEntry

// EducationEntry represents a single educational qualification.
type EducationEntry = scorer.EducationEntry

// JobRequirements describes the requirements extracted from a job posting.
type JobRequirements = scorer.JobRequirements

// ScoreBreakdown holds the individual component scores and the final result.
type ScoreBreakdown = scorer.ScoreBreakdown

// Calculate computes the acceptance likelihood score for a candidate against a job.
func Calculate(profile CandidateProfile, job JobRequirements) ScoreBreakdown {
	return scorer.Calculate(profile, job)
}
