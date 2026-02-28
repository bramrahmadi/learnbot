package analysis_test

import (
	"testing"

	"github.com/learnbot/resume-parser/internal/scorer"
	"github.com/learnbot/resume-parser/pkg/analysis"
)

func TestNew(t *testing.T) {
	a := analysis.New()
	if a == nil {
		t.Fatal("expected non-nil Analyzer")
	}
}

func TestAnalyze_EmptyInputs(t *testing.T) {
	a := analysis.New()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{}

	result := a.Analyze(profile, job)

	if result.TotalGaps != 0 {
		t.Errorf("expected 0 gaps for empty inputs, got %d", result.TotalGaps)
	}
	if result.ReadinessScore < 0 || result.ReadinessScore > 100 {
		t.Errorf("readiness score out of range: %.2f", result.ReadinessScore)
	}
}

func TestAnalyze_WithSkillGap(t *testing.T) {
	a := analysis.New()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go", "Python"},
	}

	result := a.Analyze(profile, job)

	if result.CriticalGapCount != 1 {
		t.Errorf("expected 1 critical gap (Python), got %d", result.CriticalGapCount)
	}
}

func TestGapCategoryConstants(t *testing.T) {
	if analysis.GapCategoryCritical == "" {
		t.Error("GapCategoryCritical should not be empty")
	}
	if analysis.GapCategoryImportant == "" {
		t.Error("GapCategoryImportant should not be empty")
	}
	if analysis.GapCategoryNiceToHave == "" {
		t.Error("GapCategoryNiceToHave should not be empty")
	}
}
