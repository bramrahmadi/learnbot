package recommend_test

import (
	"testing"

	"github.com/learnbot/resume-parser/internal/scorer"
	"github.com/learnbot/resume-parser/pkg/recommend"
)

func TestNew(t *testing.T) {
	e := recommend.New()
	if e == nil {
		t.Fatal("expected non-nil Engine")
	}
}

func TestGetCatalog(t *testing.T) {
	catalog := recommend.GetCatalog()
	if len(catalog) == 0 {
		t.Error("expected non-empty resource catalog")
	}
}

func TestGenerate_EmptyInputs(t *testing.T) {
	e := recommend.New()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{}
	prefs := recommend.UserPreferences{
		WeeklyHoursAvailable: 10,
	}

	plan := e.Generate(profile, job, prefs)

	// With no gaps, the plan should be valid but may have no phases
	if plan.TotalEstimatedHours < 0 {
		t.Errorf("expected non-negative total hours, got %.2f", plan.TotalEstimatedHours)
	}
}

func TestGenerate_WithSkillGap(t *testing.T) {
	e := recommend.New()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go", "Python"},
	}
	prefs := recommend.UserPreferences{
		WeeklyHoursAvailable: 10,
		PreferFree:           true,
	}

	plan := e.Generate(profile, job, prefs)

	// Should have at least one phase for the Python gap
	if len(plan.Phases) == 0 {
		t.Error("expected at least one learning phase for skill gap")
	}
}
