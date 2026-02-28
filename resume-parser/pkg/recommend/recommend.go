// Package recommend provides a public API for training recommendations.
// This package wraps the internal recommendation package for use by external modules.
package recommend

import (
	"github.com/learnbot/resume-parser/internal/recommendation"
	"github.com/learnbot/resume-parser/internal/scorer"
)

// Re-export types from recommendation for external use.

// UserPreferences captures the user's learning preferences.
type UserPreferences = recommendation.UserPreferences

// ResourceEntry represents a single learning resource.
type ResourceEntry = recommendation.ResourceEntry

// LearningPlan is the complete personalized learning plan output.
type LearningPlan = recommendation.LearningPlan

// Engine is the training recommendation engine.
type Engine struct {
	inner *recommendation.Engine
}

// New creates a new recommendation Engine with the built-in resource catalog.
func New() *Engine {
	return &Engine{inner: recommendation.New()}
}

// Generate produces a personalized learning plan.
func (e *Engine) Generate(
	profile scorer.CandidateProfile,
	job scorer.JobRequirements,
	prefs UserPreferences,
) LearningPlan {
	return e.inner.Generate(profile, job, prefs)
}

// GetCatalog returns the built-in resource catalog.
func GetCatalog() []ResourceEntry {
	return recommendation.GetCatalog()
}
