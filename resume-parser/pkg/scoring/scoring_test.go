package scoring_test

import (
	"testing"

	"github.com/learnbot/resume-parser/pkg/scoring"
)

func TestCalculate_EmptyInputs(t *testing.T) {
	profile := scoring.CandidateProfile{}
	job := scoring.JobRequirements{}

	result := scoring.Calculate(profile, job)

	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("overall score out of range: %.2f", result.OverallScore)
	}
}

func TestCalculate_PerfectSkillMatch(t *testing.T) {
	profile := scoring.CandidateProfile{
		Skills: []scoring.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "PostgreSQL", Proficiency: "advanced"},
		},
	}
	job := scoring.JobRequirements{
		RequiredSkills: []string{"Go", "PostgreSQL"},
	}

	result := scoring.Calculate(profile, job)

	if result.SkillMatchScore <= 0 {
		t.Errorf("expected positive skill match score, got %.2f", result.SkillMatchScore)
	}
	if len(result.MatchedRequiredSkills) != 2 {
		t.Errorf("expected 2 matched skills, got %d", len(result.MatchedRequiredSkills))
	}
}

func TestCalculate_ScoreInRange(t *testing.T) {
	profile := scoring.CandidateProfile{
		Skills: []scoring.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
		YearsOfExperience: 5,
	}
	job := scoring.JobRequirements{
		RequiredSkills:     []string{"Go"},
		MinYearsExperience: 3,
	}

	result := scoring.Calculate(profile, job)

	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("overall score out of [0,100]: %.2f", result.OverallScore)
	}
	if result.SkillMatchScore < 0 || result.SkillMatchScore > 1 {
		t.Errorf("skill match score out of [0,1]: %.2f", result.SkillMatchScore)
	}
	if result.ExperienceMatchScore < 0 || result.ExperienceMatchScore > 1 {
		t.Errorf("experience match score out of [0,1]: %.2f", result.ExperienceMatchScore)
	}
}
