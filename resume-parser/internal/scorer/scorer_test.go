package scorer

import (
	"math"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// approxEqual returns true if a and b differ by less than epsilon.
func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// mustContain fails the test if want is not in got.
func mustContain(t *testing.T, got []string, want string) {
	t.Helper()
	for _, s := range got {
		if s == want {
			return
		}
	}
	t.Errorf("expected %q to be in %v", want, got)
}

// mustNotContain fails the test if unwanted is in got.
func mustNotContain(t *testing.T, got []string, unwanted string) {
	t.Helper()
	for _, s := range got {
		if s == unwanted {
			t.Errorf("expected %q NOT to be in %v", unwanted, got)
			return
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Calculate (overall)
// ─────────────────────────────────────────────────────────────────────────────

func TestCalculate_PerfectMatch(t *testing.T) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "PostgreSQL", Proficiency: "advanced"},
			{Name: "Docker", Proficiency: "advanced"},
		},
		YearsOfExperience: 5,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Senior Software Engineer", Industry: "software", DurationMonths: 36},
			{Title: "Software Engineer", Industry: "software", DurationMonths: 24},
		},
		Education: []EducationEntry{
			{DegreeLevel: "bachelor", FieldOfStudy: "Computer Science"},
		},
		LocationCity:    "San Francisco",
		LocationCountry: "US",
		RemotePreference: "any",
	}

	job := JobRequirements{
		Title:               "Senior Software Engineer",
		RequiredSkills:      []string{"Go", "PostgreSQL", "Docker"},
		MinYearsExperience:  3,
		RequiredDegreeLevel: "bachelor",
		PreferredFields:     []string{"Computer Science"},
		LocationType:        "remote",
		Industry:            "software",
		ExperienceLevel:     "senior",
	}

	result := Calculate(profile, job)

	if result.OverallScore < 80 {
		t.Errorf("expected overall score >= 80 for perfect match, got %.2f", result.OverallScore)
	}
	// Skill score is capped at 0.80 when no preferred skills are specified
	// (required skills contribute 80% of the skill score).
	if result.SkillMatchScore < 0.75 {
		t.Errorf("expected skill match >= 0.75, got %.2f", result.SkillMatchScore)
	}
	if result.ExperienceMatchScore < 0.8 {
		t.Errorf("expected experience match >= 0.8, got %.2f", result.ExperienceMatchScore)
	}
	if result.EducationMatchScore < 0.9 {
		t.Errorf("expected education match >= 0.9, got %.2f", result.EducationMatchScore)
	}
	if result.LocationFitScore < 0.9 {
		t.Errorf("expected location fit >= 0.9, got %.2f", result.LocationFitScore)
	}
	if result.IndustryRelevanceScore < 0.9 {
		t.Errorf("expected industry relevance >= 0.9, got %.2f", result.IndustryRelevanceScore)
	}
}

func TestCalculate_NoMatch(t *testing.T) {
	profile := CandidateProfile{
		Skills:            []CandidateSkill{{Name: "Excel", Proficiency: "beginner"}},
		YearsOfExperience: 0,
		Education:         []EducationEntry{{DegreeLevel: "high_school"}},
		LocationCountry:   "Brazil",
		WillingToRelocate: false,
	}

	job := JobRequirements{
		Title:               "Principal Machine Learning Engineer",
		RequiredSkills:      []string{"Python", "TensorFlow", "PyTorch", "Kubernetes"},
		MinYearsExperience:  10,
		RequiredDegreeLevel: "doctorate",
		LocationCity:        "Seattle",
		LocationCountry:     "US",
		LocationType:        "on_site",
		Industry:            "artificial intelligence",
		ExperienceLevel:     "executive",
	}

	result := Calculate(profile, job)

	if result.OverallScore > 30 {
		t.Errorf("expected overall score <= 30 for no match, got %.2f", result.OverallScore)
	}
}

func TestCalculate_ScoreInRange(t *testing.T) {
	// Score must always be in [0, 100].
	profiles := []CandidateProfile{
		{}, // empty profile
		{
			Skills:            []CandidateSkill{{Name: "Go", Proficiency: "expert"}},
			YearsOfExperience: 100, // extreme over-qualification
		},
	}
	jobs := []JobRequirements{
		{}, // empty job
		{
			RequiredSkills:     []string{"Python"},
			MinYearsExperience: 1,
		},
	}

	for _, p := range profiles {
		for _, j := range jobs {
			result := Calculate(p, j)
			if result.OverallScore < 0 || result.OverallScore > 100 {
				t.Errorf("score out of range [0,100]: %.2f", result.OverallScore)
			}
		}
	}
}

func TestCalculate_WeightedFormula(t *testing.T) {
	// Verify the weighted formula is applied correctly by constructing a case
	// where we can predict the exact component scores.
	profile := CandidateProfile{
		// All required skills matched at expert level → skill score = 1.0
		Skills: []CandidateSkill{
			{Name: "Python", Proficiency: "expert"},
		},
		// Exactly meets minimum experience → experience score ≈ 1.0
		YearsOfExperience: 3,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Data Scientist", Industry: "finance", DurationMonths: 36},
		},
		// Meets degree requirement → education score = 1.0
		Education: []EducationEntry{
			{DegreeLevel: "master", FieldOfStudy: "Data Science"},
		},
		// Remote preference matches remote job → location score = 1.0
		RemotePreference: "remote",
		// Industry matches → industry score = 1.0
	}

	job := JobRequirements{
		Title:               "Data Scientist",
		RequiredSkills:      []string{"Python"},
		MinYearsExperience:  3,
		RequiredDegreeLevel: "master",
		PreferredFields:     []string{"Data Science"},
		LocationType:        "remote",
		Industry:            "finance",
		ExperienceLevel:     "mid",
	}

	result := Calculate(profile, job)

	// Skill score is 0.80 (required only, no preferred skills → 80% of max).
	// Other components are 1.0.
	// Expected: (0.80*0.35 + 1.0*0.25 + 1.0*0.15 + 1.0*0.10 + 1.0*0.15) * 100 = 93.0
	if result.OverallScore < 90 {
		t.Errorf("expected score >= 90, got %.2f (breakdown: skill=%.2f exp=%.2f edu=%.2f loc=%.2f ind=%.2f)",
			result.OverallScore,
			result.SkillMatchScore,
			result.ExperienceMatchScore,
			result.EducationMatchScore,
			result.LocationFitScore,
			result.IndustryRelevanceScore,
		)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Skill match
// ─────────────────────────────────────────────────────────────────────────────

func TestScoreSkillMatch_AllRequired(t *testing.T) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Docker", Proficiency: "advanced"},
			{Name: "Kubernetes", Proficiency: "intermediate"},
		},
	}
	job := JobRequirements{
		RequiredSkills: []string{"Go", "Docker", "Kubernetes"},
	}

	result := Calculate(profile, job)

	// All 3 required skills matched. Score = weighted_avg_proficiency * 0.80
	// (expert=1.0, advanced=0.9, intermediate=0.75) → avg=0.883 → *0.80 = 0.706
	if result.SkillMatchScore < 0.65 {
		t.Errorf("expected skill match >= 0.65, got %.2f", result.SkillMatchScore)
	}
	if len(result.MissingRequiredSkills) != 0 {
		t.Errorf("expected no missing skills, got %v", result.MissingRequiredSkills)
	}
	if len(result.MatchedRequiredSkills) != 3 {
		t.Errorf("expected 3 matched skills, got %d", len(result.MatchedRequiredSkills))
	}
}

func TestScoreSkillMatch_NoneRequired(t *testing.T) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{{Name: "Excel"}},
	}
	job := JobRequirements{
		RequiredSkills: []string{"Python", "TensorFlow", "PyTorch"},
	}

	result := Calculate(profile, job)

	if result.SkillMatchScore > 0.3 {
		t.Errorf("expected skill match <= 0.3, got %.2f", result.SkillMatchScore)
	}
	if len(result.MissingRequiredSkills) != 3 {
		t.Errorf("expected 3 missing skills, got %d: %v", len(result.MissingRequiredSkills), result.MissingRequiredSkills)
	}
}

func TestScoreSkillMatch_NoRequiredSkillsInJob(t *testing.T) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{{Name: "Go"}},
	}
	job := JobRequirements{} // no required skills

	result := Calculate(profile, job)

	if result.SkillMatchScore != 1.0 {
		t.Errorf("expected skill match = 1.0 when no required skills, got %.2f", result.SkillMatchScore)
	}
}

func TestScoreSkillMatch_AliasMatching(t *testing.T) {
	tests := []struct {
		candidateSkill string
		requiredSkill  string
		wantMatch      bool
	}{
		{"golang", "Go", true},
		{"Golang", "go", true},
		{"nodejs", "node.js", true},
		{"postgres", "postgresql", true},
		{"k8s", "kubernetes", true},
		{"js", "javascript", true},
		{"reactjs", "react", true},
		{"Python", "python", true},
		{"Excel", "Python", false},
	}

	for _, tt := range tests {
		t.Run(tt.candidateSkill+"_vs_"+tt.requiredSkill, func(t *testing.T) {
			profile := CandidateProfile{
				Skills: []CandidateSkill{{Name: tt.candidateSkill, Proficiency: "expert"}},
			}
			job := JobRequirements{
				RequiredSkills: []string{tt.requiredSkill},
			}
			result := Calculate(profile, job)
			if tt.wantMatch && len(result.MatchedRequiredSkills) == 0 {
				t.Errorf("expected %q to match %q", tt.candidateSkill, tt.requiredSkill)
			}
			if !tt.wantMatch && len(result.MatchedRequiredSkills) > 0 {
				t.Errorf("expected %q NOT to match %q", tt.candidateSkill, tt.requiredSkill)
			}
		})
	}
}

func TestScoreSkillMatch_PreferredSkillsBonus(t *testing.T) {
	// Profile has all required skills but none of the preferred.
	profileNoPreferred := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	// Profile has all required AND all preferred skills.
	profileWithPreferred := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Docker", Proficiency: "advanced"},
			{Name: "Kubernetes", Proficiency: "advanced"},
		},
	}

	job := JobRequirements{
		RequiredSkills:  []string{"Go"},
		PreferredSkills: []string{"Docker", "Kubernetes"},
	}

	resultNo := Calculate(profileNoPreferred, job)
	resultWith := Calculate(profileWithPreferred, job)

	if resultWith.SkillMatchScore <= resultNo.SkillMatchScore {
		t.Errorf("expected preferred skills to increase score: with=%.2f, without=%.2f",
			resultWith.SkillMatchScore, resultNo.SkillMatchScore)
	}
	if len(resultWith.MatchedPreferredSkills) != 2 {
		t.Errorf("expected 2 matched preferred skills, got %d", len(resultWith.MatchedPreferredSkills))
	}
}

func TestScoreSkillMatch_ProficiencyImpact(t *testing.T) {
	// Expert proficiency should yield a higher score than beginner.
	profileExpert := CandidateProfile{
		Skills: []CandidateSkill{{Name: "Python", Proficiency: "expert"}},
	}
	profileBeginner := CandidateProfile{
		Skills: []CandidateSkill{{Name: "Python", Proficiency: "beginner"}},
	}
	job := JobRequirements{
		RequiredSkills: []string{"Python"},
	}

	resultExpert := Calculate(profileExpert, job)
	resultBeginner := Calculate(profileBeginner, job)

	if resultExpert.SkillMatchScore <= resultBeginner.SkillMatchScore {
		t.Errorf("expert proficiency should score higher than beginner: expert=%.2f, beginner=%.2f",
			resultExpert.SkillMatchScore, resultBeginner.SkillMatchScore)
	}
}

func TestScoreSkillMatch_EmptyProfile(t *testing.T) {
	profile := CandidateProfile{}
	job := JobRequirements{
		RequiredSkills: []string{"Go", "Python"},
	}

	result := Calculate(profile, job)

	if result.SkillMatchScore != 0 {
		t.Errorf("expected skill match = 0 for empty profile, got %.2f", result.SkillMatchScore)
	}
	if len(result.MissingRequiredSkills) != 2 {
		t.Errorf("expected 2 missing skills, got %d", len(result.MissingRequiredSkills))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Experience match
// ─────────────────────────────────────────────────────────────────────────────

func TestScoreExperienceMatch_MeetsMinimum(t *testing.T) {
	profile := CandidateProfile{YearsOfExperience: 5}
	job := JobRequirements{MinYearsExperience: 3, ExperienceLevel: "mid"}

	result := Calculate(profile, job)

	if result.ExperienceMatchScore < 0.7 {
		t.Errorf("expected experience match >= 0.7, got %.2f", result.ExperienceMatchScore)
	}
}

func TestScoreExperienceMatch_BelowMinimum(t *testing.T) {
	profile := CandidateProfile{YearsOfExperience: 1}
	job := JobRequirements{MinYearsExperience: 5}

	result := Calculate(profile, job)

	if result.ExperienceMatchScore > 0.5 {
		t.Errorf("expected experience match <= 0.5 for under-qualified, got %.2f", result.ExperienceMatchScore)
	}
}

func TestScoreExperienceMatch_ZeroExperience(t *testing.T) {
	profile := CandidateProfile{YearsOfExperience: 0}
	job := JobRequirements{MinYearsExperience: 5}

	result := Calculate(profile, job)

	if result.ExperienceMatchScore > 0.4 {
		t.Errorf("expected low experience match for 0 years vs 5 required, got %.2f", result.ExperienceMatchScore)
	}
}

func TestScoreExperienceMatch_InferFromWorkHistory(t *testing.T) {
	// YearsOfExperience not set; should be inferred from work history.
	profile := CandidateProfile{
		YearsOfExperience: 0,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Engineer", DurationMonths: 24},
			{Title: "Senior Engineer", DurationMonths: 24},
		},
	}
	job := JobRequirements{MinYearsExperience: 3}

	result := Calculate(profile, job)

	// 48 months = 4 years, which exceeds the 3-year minimum.
	if result.ExperienceMatchScore < 0.7 {
		t.Errorf("expected experience match >= 0.7 when inferred from history, got %.2f", result.ExperienceMatchScore)
	}
}

func TestScoreExperienceMatch_ExperienceLevelFallback(t *testing.T) {
	// No MinYearsExperience set; should fall back to ExperienceLevel.
	profile := CandidateProfile{YearsOfExperience: 4}
	job := JobRequirements{
		MinYearsExperience: 0,
		ExperienceLevel:    "mid", // maps to ~3 years
	}

	result := Calculate(profile, job)

	if result.ExperienceMatchScore < 0.7 {
		t.Errorf("expected experience match >= 0.7 using level fallback, got %.2f", result.ExperienceMatchScore)
	}
}

func TestScoreExperienceMatch_TitleSimilarity(t *testing.T) {
	profileMatch := CandidateProfile{
		YearsOfExperience: 5,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Senior Software Engineer", DurationMonths: 36},
		},
	}
	profileNoMatch := CandidateProfile{
		YearsOfExperience: 5,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Marketing Manager", DurationMonths: 36},
		},
	}
	job := JobRequirements{
		Title:              "Senior Software Engineer",
		MinYearsExperience: 5,
	}

	resultMatch := Calculate(profileMatch, job)
	resultNoMatch := Calculate(profileNoMatch, job)

	if resultMatch.ExperienceMatchScore <= resultNoMatch.ExperienceMatchScore {
		t.Errorf("title match should increase experience score: match=%.2f, no_match=%.2f",
			resultMatch.ExperienceMatchScore, resultNoMatch.ExperienceMatchScore)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Education match
// ─────────────────────────────────────────────────────────────────────────────

func TestScoreEducationMatch_NoRequirement(t *testing.T) {
	profile := CandidateProfile{} // no education
	job := JobRequirements{}      // no degree requirement

	result := Calculate(profile, job)

	if result.EducationMatchScore != 1.0 {
		t.Errorf("expected education match = 1.0 when no requirement, got %.2f", result.EducationMatchScore)
	}
}

func TestScoreEducationMatch_MeetsRequirement(t *testing.T) {
	profile := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "bachelor", FieldOfStudy: "Computer Science"}},
	}
	job := JobRequirements{
		RequiredDegreeLevel: "bachelor",
		PreferredFields:     []string{"Computer Science"},
	}

	result := Calculate(profile, job)

	if result.EducationMatchScore < 0.9 {
		t.Errorf("expected education match >= 0.9, got %.2f", result.EducationMatchScore)
	}
}

func TestScoreEducationMatch_ExceedsRequirement(t *testing.T) {
	profile := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "doctorate", FieldOfStudy: "Computer Science"}},
	}
	job := JobRequirements{
		RequiredDegreeLevel: "bachelor",
	}

	result := Calculate(profile, job)

	if result.EducationMatchScore < 0.9 {
		t.Errorf("expected education match >= 0.9 when exceeding requirement, got %.2f", result.EducationMatchScore)
	}
}

func TestScoreEducationMatch_BelowRequirement(t *testing.T) {
	profile := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "high_school"}},
	}
	job := JobRequirements{
		RequiredDegreeLevel: "master",
	}

	result := Calculate(profile, job)

	if result.EducationMatchScore > 0.5 {
		t.Errorf("expected education match <= 0.5 for high_school vs master, got %.2f", result.EducationMatchScore)
	}
}

func TestScoreEducationMatch_NoEducationInfo(t *testing.T) {
	profile := CandidateProfile{} // no education entries
	job := JobRequirements{
		RequiredDegreeLevel: "bachelor",
	}

	result := Calculate(profile, job)

	// Should get partial credit (0.3) since we don't know.
	if result.EducationMatchScore > 0.5 {
		t.Errorf("expected low education match for missing info, got %.2f", result.EducationMatchScore)
	}
}

func TestScoreEducationMatch_FieldBonus(t *testing.T) {
	// Use associate degree (one level below bachelor) so degree score = 0.6,
	// leaving room for the field bonus to push the score higher.
	profileMatch := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "associate", FieldOfStudy: "Computer Science"}},
	}
	profileNoMatch := CandidateProfile{
		Education: []EducationEntry{{DegreeLevel: "associate", FieldOfStudy: "History"}},
	}
	job := JobRequirements{
		RequiredDegreeLevel: "bachelor",
		PreferredFields:     []string{"Computer Science", "Software Engineering"},
	}

	resultMatch := Calculate(profileMatch, job)
	resultNoMatch := Calculate(profileNoMatch, job)

	if resultMatch.EducationMatchScore <= resultNoMatch.EducationMatchScore {
		t.Errorf("field match should increase education score: match=%.2f, no_match=%.2f",
			resultMatch.EducationMatchScore, resultNoMatch.EducationMatchScore)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Location fit
// ─────────────────────────────────────────────────────────────────────────────

func TestScoreLocationFit_RemoteJobRemoteCandidate(t *testing.T) {
	profile := CandidateProfile{RemotePreference: "remote"}
	job := JobRequirements{LocationType: "remote"}

	result := Calculate(profile, job)

	if result.LocationFitScore != 1.0 {
		t.Errorf("expected location fit = 1.0 for remote/remote, got %.2f", result.LocationFitScore)
	}
}

func TestScoreLocationFit_RemoteJobOnSiteCandidate(t *testing.T) {
	profile := CandidateProfile{RemotePreference: "on_site"}
	job := JobRequirements{LocationType: "remote"}

	result := Calculate(profile, job)

	if result.LocationFitScore >= 1.0 {
		t.Errorf("expected location fit < 1.0 for remote job / on-site candidate, got %.2f", result.LocationFitScore)
	}
}

func TestScoreLocationFit_SameCity(t *testing.T) {
	profile := CandidateProfile{
		LocationCity:    "New York",
		LocationCountry: "US",
	}
	job := JobRequirements{
		LocationCity:    "New York",
		LocationCountry: "US",
		LocationType:    "on_site",
	}

	result := Calculate(profile, job)

	if result.LocationFitScore != 1.0 {
		t.Errorf("expected location fit = 1.0 for same city, got %.2f", result.LocationFitScore)
	}
}

func TestScoreLocationFit_SameCountryDifferentCity(t *testing.T) {
	profile := CandidateProfile{
		LocationCity:    "Los Angeles",
		LocationCountry: "US",
	}
	job := JobRequirements{
		LocationCity:    "New York",
		LocationCountry: "US",
		LocationType:    "on_site",
	}

	result := Calculate(profile, job)

	if result.LocationFitScore < 0.5 || result.LocationFitScore >= 1.0 {
		t.Errorf("expected location fit in (0.5, 1.0) for same country, got %.2f", result.LocationFitScore)
	}
}

func TestScoreLocationFit_DifferentCountryWillingToRelocate(t *testing.T) {
	profile := CandidateProfile{
		LocationCountry:   "Brazil",
		WillingToRelocate: true,
	}
	job := JobRequirements{
		LocationCountry: "US",
		LocationType:    "on_site",
	}

	result := Calculate(profile, job)

	if result.LocationFitScore < 0.5 {
		t.Errorf("expected location fit >= 0.5 for willing to relocate, got %.2f", result.LocationFitScore)
	}
}

func TestScoreLocationFit_DifferentCountryNotWillingToRelocate(t *testing.T) {
	profile := CandidateProfile{
		LocationCountry:   "Brazil",
		WillingToRelocate: false,
	}
	job := JobRequirements{
		LocationCountry: "US",
		LocationType:    "on_site",
	}

	result := Calculate(profile, job)

	if result.LocationFitScore > 0.4 {
		t.Errorf("expected low location fit for different country / no relocation, got %.2f", result.LocationFitScore)
	}
}

func TestScoreLocationFit_HybridJob(t *testing.T) {
	tests := []struct {
		pref      string
		wantScore float64
	}{
		{"hybrid", 1.0},
		{"any", 1.0},
		{"", 1.0},
		{"remote", 0.6},
		{"on_site", 0.8},
	}

	for _, tt := range tests {
		t.Run("pref_"+tt.pref, func(t *testing.T) {
			profile := CandidateProfile{RemotePreference: tt.pref}
			job := JobRequirements{LocationType: "hybrid"}
			result := Calculate(profile, job)
			if !approxEqual(result.LocationFitScore, tt.wantScore, 0.01) {
				t.Errorf("hybrid job with pref=%q: expected %.2f, got %.2f",
					tt.pref, tt.wantScore, result.LocationFitScore)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Industry relevance
// ─────────────────────────────────────────────────────────────────────────────

func TestScoreIndustryRelevance_NoIndustryRequirement(t *testing.T) {
	profile := CandidateProfile{}
	job := JobRequirements{} // no industry

	result := Calculate(profile, job)

	if result.IndustryRelevanceScore != 1.0 {
		t.Errorf("expected industry relevance = 1.0 when no requirement, got %.2f", result.IndustryRelevanceScore)
	}
}

func TestScoreIndustryRelevance_DirectMatch(t *testing.T) {
	profile := CandidateProfile{
		WorkHistory: []WorkHistoryEntry{
			{Title: "Engineer", Industry: "Software", DurationMonths: 36},
		},
	}
	job := JobRequirements{Industry: "software"}

	result := Calculate(profile, job)

	if result.IndustryRelevanceScore != 1.0 {
		t.Errorf("expected industry relevance = 1.0 for direct match, got %.2f", result.IndustryRelevanceScore)
	}
}

func TestScoreIndustryRelevance_RelatedIndustry(t *testing.T) {
	profile := CandidateProfile{
		WorkHistory: []WorkHistoryEntry{
			{Title: "Engineer", Industry: "fintech", DurationMonths: 36},
		},
	}
	job := JobRequirements{
		Industry:          "finance",
		RelatedIndustries: []string{"fintech", "banking"},
	}

	result := Calculate(profile, job)

	if result.IndustryRelevanceScore < 0.6 {
		t.Errorf("expected industry relevance >= 0.6 for related industry, got %.2f", result.IndustryRelevanceScore)
	}
}

func TestScoreIndustryRelevance_NoMatch(t *testing.T) {
	profile := CandidateProfile{
		WorkHistory: []WorkHistoryEntry{
			{Title: "Chef", Industry: "hospitality", DurationMonths: 60},
		},
	}
	job := JobRequirements{Industry: "aerospace"}

	result := Calculate(profile, job)

	if result.IndustryRelevanceScore > 0.3 {
		t.Errorf("expected low industry relevance for no match, got %.2f", result.IndustryRelevanceScore)
	}
}

func TestScoreIndustryRelevance_NoWorkHistory(t *testing.T) {
	profile := CandidateProfile{} // no work history
	job := JobRequirements{Industry: "software"}

	result := Calculate(profile, job)

	// Should return neutral score (0.5) when no work history.
	if result.IndustryRelevanceScore != 0.5 {
		t.Errorf("expected industry relevance = 0.5 for no work history, got %.2f", result.IndustryRelevanceScore)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Edge cases
// ─────────────────────────────────────────────────────────────────────────────

func TestCalculate_EmptyProfileAndJob(t *testing.T) {
	result := Calculate(CandidateProfile{}, JobRequirements{})

	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("score out of range for empty inputs: %.2f", result.OverallScore)
	}
}

func TestCalculate_DuplicateSkillsInProfile(t *testing.T) {
	// Duplicate skills should not inflate the score.
	profile := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Go", Proficiency: "beginner"}, // duplicate
			{Name: "go", Proficiency: "advanced"},  // same, different case
		},
	}
	job := JobRequirements{
		RequiredSkills: []string{"Go"},
	}

	result := Calculate(profile, job)

	if len(result.MatchedRequiredSkills) != 1 {
		t.Errorf("expected 1 matched skill (deduped), got %d", len(result.MatchedRequiredSkills))
	}
}

func TestCalculate_CaseInsensitiveSkillMatching(t *testing.T) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{{Name: "PYTHON", Proficiency: "expert"}},
	}
	job := JobRequirements{
		RequiredSkills: []string{"python"},
	}

	result := Calculate(profile, job)

	if len(result.MatchedRequiredSkills) == 0 {
		t.Error("expected case-insensitive skill match")
	}
}

func TestCalculate_OverQualifiedExperience(t *testing.T) {
	// Significantly over-qualified should still get a good (but not perfect) score.
	profile := CandidateProfile{YearsOfExperience: 30}
	job := JobRequirements{
		MinYearsExperience: 2,
		MaxYearsExperience: 5,
	}

	result := Calculate(profile, job)

	// Should still be a reasonable score, not penalised too harshly.
	if result.ExperienceMatchScore < 0.5 {
		t.Errorf("expected experience match >= 0.5 for over-qualified, got %.2f", result.ExperienceMatchScore)
	}
}

func TestCalculate_BreakdownSumApprox(t *testing.T) {
	// Verify the weighted sum matches the overall score.
	profile := CandidateProfile{
		Skills:            []CandidateSkill{{Name: "Go", Proficiency: "advanced"}},
		YearsOfExperience: 4,
		WorkHistory:       []WorkHistoryEntry{{Title: "Engineer", Industry: "tech", DurationMonths: 48}},
		Education:         []EducationEntry{{DegreeLevel: "bachelor", FieldOfStudy: "CS"}},
		LocationCountry:   "US",
		RemotePreference:  "remote",
	}
	job := JobRequirements{
		Title:               "Software Engineer",
		RequiredSkills:      []string{"Go"},
		MinYearsExperience:  3,
		RequiredDegreeLevel: "bachelor",
		LocationType:        "remote",
		Industry:            "tech",
	}

	result := Calculate(profile, job)

	expected := (result.SkillMatchScore*WeightSkillMatch +
		result.ExperienceMatchScore*WeightExperienceMatch +
		result.EducationMatchScore*WeightEducationMatch +
		result.LocationFitScore*WeightLocationFit +
		result.IndustryRelevanceScore*WeightIndustryRelevance) * 100

	if !approxEqual(result.OverallScore, roundTo2(expected), 0.1) {
		t.Errorf("overall score %.2f does not match weighted sum %.2f", result.OverallScore, expected)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

func TestNormalizeSkillName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Go", "go"},
		{"  Python  ", "python"},
		{"Node.js", "node.js"},
		{"C++", "c++"},
		{"", ""},
	}
	for _, tt := range tests {
		got := normalizeSkillName(tt.input)
		if got != tt.want {
			t.Errorf("normalizeSkillName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSkillsAreAliases(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"golang", "go", true},
		{"go", "golang", true},
		{"js", "javascript", true},
		{"k8s", "kubernetes", true},
		{"python", "python", true},
		{"python", "java", false},
		{"spring boot", "spring", true},
		{"ab", "a", false}, // too short
	}
	for _, tt := range tests {
		got := skillsAreAliases(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("skillsAreAliases(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestComputeYearsScore(t *testing.T) {
	tests := []struct {
		candidate float64
		targetMin float64
		targetMax float64
		wantMin   float64
		wantMax   float64
	}{
		{5, 3, 0, 1.0, 1.0},   // meets minimum, no max
		{1, 5, 0, 0.1, 0.3},   // under-qualified
		{0, 5, 0, 0.1, 0.15},  // zero experience
		{20, 5, 8, 0.7, 0.9},  // over-qualified (> 2× max)
		{0, 0, 0, 1.0, 1.0},   // no requirement
	}
	for _, tt := range tests {
		got := computeYearsScore(tt.candidate, tt.targetMin, tt.targetMax)
		if got < tt.wantMin || got > tt.wantMax {
			t.Errorf("computeYearsScore(%.1f, %.1f, %.1f) = %.2f, want [%.2f, %.2f]",
				tt.candidate, tt.targetMin, tt.targetMax, got, tt.wantMin, tt.wantMax)
		}
	}
}

func TestWordOverlapScore(t *testing.T) {
	tests := []struct {
		a, b []string
		want float64
	}{
		{[]string{"senior", "software", "engineer"}, []string{"senior", "software", "engineer"}, 1.0},
		{[]string{"software", "engineer"}, []string{"marketing", "manager"}, 0.0},
		{[]string{}, []string{"a"}, 0.0},
		{[]string{"a"}, []string{}, 0.0},
	}
	for _, tt := range tests {
		got := wordOverlapScore(tt.a, tt.b)
		if !approxEqual(got, tt.want, 0.01) {
			t.Errorf("wordOverlapScore(%v, %v) = %.2f, want %.2f", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestRoundTo2(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{1.234, 1.23},
		{1.235, 1.24},
		{0.0, 0.0},
		{100.0, 100.0},
	}
	for _, tt := range tests {
		got := roundTo2(tt.input)
		if !approxEqual(got, tt.want, 0.001) {
			t.Errorf("roundTo2(%.3f) = %.3f, want %.3f", tt.input, got, tt.want)
		}
	}
}
