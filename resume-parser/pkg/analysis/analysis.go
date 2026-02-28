// Package analysis provides a public API for skill gap analysis.
// This package wraps the internal gapanalysis package for use by external modules.
package analysis

import (
	"github.com/learnbot/resume-parser/internal/gapanalysis"
	"github.com/learnbot/resume-parser/internal/scorer"
)

// Re-export types from gapanalysis for external use.

// GapCategory classifies a skill gap by its importance.
type GapCategory = gapanalysis.GapCategory

const (
	GapCategoryCritical    = gapanalysis.GapCategoryCritical
	GapCategoryImportant   = gapanalysis.GapCategoryImportant
	GapCategoryNiceToHave  = gapanalysis.GapCategoryNiceToHave
)

// SkillGap represents a single missing skill with analysis metadata.
type SkillGap = gapanalysis.SkillGap

// GapAnalysisResult is the complete output of the gap analysis engine.
type GapAnalysisResult = gapanalysis.GapAnalysisResult

// Analyzer performs skill gap analysis.
type Analyzer struct {
	inner *gapanalysis.Analyzer
}

// New creates a new Analyzer.
func New() *Analyzer {
	return &Analyzer{inner: gapanalysis.New()}
}

// Analyze computes the skill gap analysis for a candidate against a job.
func (a *Analyzer) Analyze(profile scorer.CandidateProfile, job scorer.JobRequirements) GapAnalysisResult {
	return a.inner.Analyze(profile, job)
}
