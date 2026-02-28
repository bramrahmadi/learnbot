package scorer

import (
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Performance benchmarks for acceptance likelihood scoring
// ─────────────────────────────────────────────────────────────────────────────

// BenchmarkCalculate_EmptyInputs benchmarks scoring with empty inputs.
func BenchmarkCalculate_EmptyInputs(b *testing.B) {
	profile := CandidateProfile{}
	job := JobRequirements{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Calculate(profile, job)
	}
}

// BenchmarkCalculate_SingleSkill benchmarks scoring with a single skill.
func BenchmarkCalculate_SingleSkill(b *testing.B) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	job := JobRequirements{
		RequiredSkills: []string{"Go"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Calculate(profile, job)
	}
}

// BenchmarkCalculate_ManySkills benchmarks scoring with many skills.
func BenchmarkCalculate_ManySkills(b *testing.B) {
	skills := []CandidateSkill{
		{Name: "Go", Proficiency: "expert"},
		{Name: "Python", Proficiency: "advanced"},
		{Name: "Java", Proficiency: "intermediate"},
		{Name: "JavaScript", Proficiency: "advanced"},
		{Name: "TypeScript", Proficiency: "advanced"},
		{Name: "Rust", Proficiency: "beginner"},
		{Name: "PostgreSQL", Proficiency: "advanced"},
		{Name: "MySQL", Proficiency: "intermediate"},
		{Name: "MongoDB", Proficiency: "intermediate"},
		{Name: "Redis", Proficiency: "advanced"},
		{Name: "Docker", Proficiency: "expert"},
		{Name: "Kubernetes", Proficiency: "advanced"},
		{Name: "AWS", Proficiency: "advanced"},
		{Name: "Terraform", Proficiency: "intermediate"},
		{Name: "React", Proficiency: "intermediate"},
		{Name: "Node.js", Proficiency: "advanced"},
		{Name: "TensorFlow", Proficiency: "beginner"},
		{Name: "Git", Proficiency: "expert"},
		{Name: "Linux", Proficiency: "advanced"},
		{Name: "Kafka", Proficiency: "intermediate"},
	}

	profile := CandidateProfile{
		Skills:            skills,
		YearsOfExperience: 8,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Senior Engineer", Industry: "software", DurationMonths: 48},
			{Title: "Engineer", Industry: "fintech", DurationMonths: 36},
		},
		Education: []EducationEntry{
			{DegreeLevel: "master", FieldOfStudy: "Computer Science"},
		},
		LocationCity:     "San Francisco",
		LocationCountry:  "US",
		RemotePreference: "hybrid",
	}

	job := JobRequirements{
		Title:               "Staff Engineer",
		RequiredSkills:      []string{"Go", "PostgreSQL", "Docker", "Kubernetes", "AWS"},
		PreferredSkills:     []string{"Python", "Terraform", "Kafka", "Redis"},
		MinYearsExperience:  6,
		RequiredDegreeLevel: "bachelor",
		LocationType:        "hybrid",
		Industry:            "software",
		ExperienceLevel:     "senior",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Calculate(profile, job)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Additional unit tests for acceptance likelihood calculation
// ─────────────────────────────────────────────────────────────────────────────

// TestCalculate_ScoreInRange verifies the overall score is always in [0, 100].
func TestCalculate_ScoreInRange(t *testing.T) {
	testCases := []struct {
		name    string
		profile CandidateProfile
		job     JobRequirements
	}{
		{
			name:    "empty inputs",
			profile: CandidateProfile{},
			job:     JobRequirements{},
		},
		{
			name: "perfect match",
			profile: CandidateProfile{
				Skills:            []CandidateSkill{{Name: "Go", Proficiency: "expert"}},
				YearsOfExperience: 10,
				Education:         []EducationEntry{{DegreeLevel: "master"}},
				RemotePreference:  "remote",
			},
			job: JobRequirements{
				RequiredSkills:      []string{"Go"},
				MinYearsExperience:  5,
				RequiredDegreeLevel: "bachelor",
				LocationType:        "remote",
			},
		},
		{
			name: "no match",
			profile: CandidateProfile{
				Skills: []CandidateSkill{{Name: "Excel", Proficiency: "beginner"}},
			},
			job: JobRequirements{
				RequiredSkills:     []string{"Go", "Kubernetes", "TensorFlow"},
				MinYearsExperience: 10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Calculate(tc.profile, tc.job)
			if result.OverallScore < 0 || result.OverallScore > 100 {
				t.Errorf("overall score out of [0, 100]: %.2f", result.OverallScore)
			}
			if result.SkillMatchScore < 0 || result.SkillMatchScore > 1 {
				t.Errorf("skill match score out of [0, 1]: %.2f", result.SkillMatchScore)
			}
			if result.ExperienceMatchScore < 0 || result.ExperienceMatchScore > 1 {
				t.Errorf("experience match score out of [0, 1]: %.2f", result.ExperienceMatchScore)
			}
			if result.EducationMatchScore < 0 || result.EducationMatchScore > 1 {
				t.Errorf("education match score out of [0, 1]: %.2f", result.EducationMatchScore)
			}
			if result.LocationFitScore < 0 || result.LocationFitScore > 1 {
				t.Errorf("location fit score out of [0, 1]: %.2f", result.LocationFitScore)
			}
			if result.IndustryRelevanceScore < 0 || result.IndustryRelevanceScore > 1 {
				t.Errorf("industry relevance score out of [0, 1]: %.2f", result.IndustryRelevanceScore)
			}
		})
	}
}

// TestCalculate_SkillMatchImproves verifies that having more required skills improves score.
func TestCalculate_SkillMatchImproves(t *testing.T) {
	job := JobRequirements{
		RequiredSkills: []string{"Go", "PostgreSQL", "Docker"},
	}

	profileNoSkills := CandidateProfile{}
	profileSomeSkills := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	profileAllSkills := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "PostgreSQL", Proficiency: "advanced"},
			{Name: "Docker", Proficiency: "advanced"},
		},
	}

	scoreNone := Calculate(profileNoSkills, job).OverallScore
	scoreSome := Calculate(profileSomeSkills, job).OverallScore
	scoreAll := Calculate(profileAllSkills, job).OverallScore

	if scoreSome <= scoreNone {
		t.Errorf("having some skills should improve score: none=%.2f, some=%.2f", scoreNone, scoreSome)
	}
	if scoreAll <= scoreSome {
		t.Errorf("having all skills should improve score: some=%.2f, all=%.2f", scoreSome, scoreAll)
	}
}

// TestCalculate_ExperienceMatchImproves verifies that more experience improves score.
func TestCalculate_ExperienceMatchImproves(t *testing.T) {
	job := JobRequirements{
		MinYearsExperience: 5,
	}

	profileNoExp := CandidateProfile{YearsOfExperience: 0}
	profileSomeExp := CandidateProfile{YearsOfExperience: 3}
	profileEnoughExp := CandidateProfile{YearsOfExperience: 7}

	scoreNone := Calculate(profileNoExp, job).ExperienceMatchScore
	scoreSome := Calculate(profileSomeExp, job).ExperienceMatchScore
	scoreEnough := Calculate(profileEnoughExp, job).ExperienceMatchScore

	if scoreSome <= scoreNone {
		t.Errorf("more experience should improve score: none=%.2f, some=%.2f", scoreNone, scoreSome)
	}
	if scoreEnough <= scoreSome {
		t.Errorf("enough experience should improve score: some=%.2f, enough=%.2f", scoreSome, scoreEnough)
	}
}

// TestCalculate_MatchedSkillsTracked verifies matched and missing skills are tracked.
func TestCalculate_MatchedSkillsTracked(t *testing.T) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	job := JobRequirements{
		RequiredSkills: []string{"Go", "Python", "Rust"},
	}

	result := Calculate(profile, job)

	mustContain(t, result.MatchedRequiredSkills, "Go")
	mustContain(t, result.MatchedRequiredSkills, "Python")
	mustContain(t, result.MissingRequiredSkills, "Rust")
	mustNotContain(t, result.MissingRequiredSkills, "Go")
	mustNotContain(t, result.MissingRequiredSkills, "Python")
}

// TestCalculate_PreferredSkillsBonus verifies preferred skills add a bonus.
func TestCalculate_PreferredSkillsBonus(t *testing.T) {
	job := JobRequirements{
		RequiredSkills:  []string{"Go"},
		PreferredSkills: []string{"Python", "Docker"},
	}

	profileWithRequired := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	profileWithAll := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Python", Proficiency: "advanced"},
			{Name: "Docker", Proficiency: "advanced"},
		},
	}

	scoreRequired := Calculate(profileWithRequired, job).OverallScore
	scoreAll := Calculate(profileWithAll, job).OverallScore

	if scoreAll <= scoreRequired {
		t.Errorf("having preferred skills should improve score: required=%.2f, all=%.2f",
			scoreRequired, scoreAll)
	}
}

// TestCalculate_LocationFit_Remote verifies remote preference matching.
func TestCalculate_LocationFit_Remote(t *testing.T) {
	jobRemote := JobRequirements{LocationType: "remote"}
	jobOnsite := JobRequirements{LocationType: "onsite"}

	profileRemote := CandidateProfile{RemotePreference: "remote"}

	scoreRemote := Calculate(profileRemote, jobRemote).LocationFitScore
	scoreOnsite := Calculate(profileRemote, jobOnsite).LocationFitScore

	if scoreRemote <= scoreOnsite {
		t.Errorf("remote preference should score higher for remote job: remote=%.2f, onsite=%.2f",
			scoreRemote, scoreOnsite)
	}
}

// TestCalculate_EducationMatch verifies education level matching.
func TestCalculate_EducationMatch(t *testing.T) {
	jobRequiresBachelor := JobRequirements{RequiredDegreeLevel: "bachelor"}

	profileNoDegree := CandidateProfile{}
	profileBachelor := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "bachelor"}},
	}
	profileMaster := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "master"}},
	}

	scoreNone := Calculate(profileNoDegree, jobRequiresBachelor).EducationMatchScore
	scoreBachelor := Calculate(profileBachelor, jobRequiresBachelor).EducationMatchScore
	scoreMaster := Calculate(profileMaster, jobRequiresBachelor).EducationMatchScore

	if scoreBachelor <= scoreNone {
		t.Errorf("having required degree should improve score: none=%.2f, bachelor=%.2f",
			scoreNone, scoreBachelor)
	}
	if scoreMaster < scoreBachelor {
		t.Errorf("higher degree should not decrease score: bachelor=%.2f, master=%.2f",
			scoreBachelor, scoreMaster)
	}
}

// TestCalculate_IndustryRelevance verifies industry matching.
func TestCalculate_IndustryRelevance(t *testing.T) {
	jobSoftware := JobRequirements{Industry: "software"}

	profileSoftware := CandidateProfile{
		WorkHistory: []WorkHistoryEntry{
			{Industry: "software", DurationMonths: 36},
		},
	}
	profileFinance := CandidateProfile{
		WorkHistory: []WorkHistoryEntry{
			{Industry: "finance", DurationMonths: 36},
		},
	}

	scoreSoftware := Calculate(profileSoftware, jobSoftware).IndustryRelevanceScore
	scoreFinance := Calculate(profileFinance, jobSoftware).IndustryRelevanceScore

	if scoreSoftware <= scoreFinance {
		t.Errorf("matching industry should score higher: software=%.2f, finance=%.2f",
			scoreSoftware, scoreFinance)
	}
}
