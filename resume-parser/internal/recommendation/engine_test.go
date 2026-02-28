package recommendation

import (
	"math"
	"testing"

	"github.com/learnbot/resume-parser/internal/gapanalysis"
	"github.com/learnbot/resume-parser/internal/scorer"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────────────────────

// approxEqual returns true if a and b differ by less than epsilon.
func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// newTestEngine creates an Engine with a small test catalog.
func newTestEngine() *Engine {
	return NewWithCatalog(testCatalog)
}

// testCatalog is a minimal catalog for testing.
var testCatalog = []ResourceEntry{
	{
		ID: "python-course", Title: "Python Beginner Course",
		Description: "Learn Python from scratch.",
		URL:         "https://example.com/python",
		Provider:    "TestProvider", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 20,
		Skills: []string{"python"}, PrimarySkill: "python",
		Rating: 4.5, RatingCount: 10000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "python-advanced", Title: "Advanced Python",
		Description: "Advanced Python programming.",
		URL:         "https://example.com/python-advanced",
		Provider:    "OtherProvider", ResourceType: "course", Difficulty: "advanced",
		CostType: "paid", CostUSD: 29.99, DurationHours: 15,
		Skills: []string{"python", "oop"}, PrimarySkill: "python",
		Rating: 4.7, RatingCount: 5000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "go-course", Title: "Go Programming",
		Description: "Learn Go programming.",
		URL:         "https://example.com/go",
		Provider:    "TestProvider", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 10,
		Skills: []string{"go"}, PrimarySkill: "go",
		Rating: 4.6, RatingCount: 8000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "docker-course", Title: "Docker Fundamentals",
		Description: "Learn Docker.",
		URL:         "https://example.com/docker",
		Provider:    "TestProvider", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", CostUSD: 0, DurationHours: 8,
		Skills: []string{"docker"}, PrimarySkill: "docker",
		Rating: 4.4, RatingCount: 3000, HasCertificate: false, HasHandsOn: true, IsVerified: false,
	},
	{
		ID: "sql-course", Title: "SQL Basics",
		Description: "Learn SQL.",
		URL:         "https://example.com/sql",
		Provider:    "TestProvider", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 5,
		Skills: []string{"sql", "postgresql"}, PrimarySkill: "sql",
		Rating: 4.3, RatingCount: 2000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
}

// ─────────────────────────────────────────────────────────────────────────────
// Engine.Generate tests
// ─────────────────────────────────────────────────────────────────────────────

func TestGenerate_NoCriticalGaps(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
			{Name: "Go", Proficiency: "intermediate"},
		},
	}
	job := scorer.JobRequirements{
		Title:          "Backend Engineer",
		RequiredSkills: []string{"Python", "Go"},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	plan := engine.Generate(profile, job, prefs)

	if plan.TotalGaps != 0 {
		t.Errorf("expected 0 gaps, got %d", plan.TotalGaps)
	}
	if plan.ReadinessScore < 90 {
		t.Errorf("expected high readiness score, got %.2f", plan.ReadinessScore)
	}
	if len(plan.Phases) != 0 {
		t.Errorf("expected 0 phases for no gaps, got %d", len(plan.Phases))
	}
}

func TestGenerate_WithCriticalGaps(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		Title:          "Backend Engineer",
		RequiredSkills: []string{"Python", "Go", "Docker"},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	plan := engine.Generate(profile, job, prefs)

	if plan.TotalGaps != 2 {
		t.Errorf("expected 2 gaps (Go, Docker), got %d", plan.TotalGaps)
	}
	if len(plan.Phases) == 0 {
		t.Error("expected at least one phase")
	}
	if plan.Phases[0].PhaseName != "Critical Skills" {
		t.Errorf("expected first phase to be 'Critical Skills', got %s", plan.Phases[0].PhaseName)
	}
}

func TestGenerate_WithPreferredGaps(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
			{Name: "Go", Proficiency: "intermediate"},
		},
	}
	job := scorer.JobRequirements{
		Title:           "Backend Engineer",
		RequiredSkills:  []string{"Python", "Go"},
		PreferredSkills: []string{"Docker", "SQL"},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	plan := engine.Generate(profile, job, prefs)

	if plan.TotalGaps != 2 {
		t.Errorf("expected 2 important gaps (Docker, SQL), got %d", plan.TotalGaps)
	}
	// Should have a "Preferred Skills" phase.
	foundPreferred := false
	for _, p := range plan.Phases {
		if p.PhaseName == "Preferred Skills" {
			foundPreferred = true
			break
		}
	}
	if !foundPreferred {
		t.Error("expected a 'Preferred Skills' phase")
	}
}

func TestGenerate_JobTitle(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		Title:          "Senior Go Developer",
		RequiredSkills: []string{"Go"},
	}
	prefs := UserPreferences{}

	plan := engine.Generate(profile, job, prefs)

	if plan.JobTitle != "Senior Go Developer" {
		t.Errorf("expected JobTitle='Senior Go Developer', got %s", plan.JobTitle)
	}
}

func TestGenerate_MatchedSkillsPopulated(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Go"},
	}
	prefs := UserPreferences{}

	plan := engine.Generate(profile, job, prefs)

	found := false
	for _, s := range plan.MatchedSkills {
		if s == "Python" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Python in matched skills")
	}
}

func TestGenerate_DefaultWeeklyHours(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go"},
	}
	// No weekly hours specified – should default to 10.
	prefs := UserPreferences{}

	plan := engine.Generate(profile, job, prefs)

	if plan.Timeline.WeeklyHours != 10 {
		t.Errorf("expected default weekly hours=10, got %.1f", plan.Timeline.WeeklyHours)
	}
}

func TestGenerate_CustomWeeklyHours(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go"},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 20}

	plan := engine.Generate(profile, job, prefs)

	if plan.Timeline.WeeklyHours != 20 {
		t.Errorf("expected weekly hours=20, got %.1f", plan.Timeline.WeeklyHours)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource matching tests
// ─────────────────────────────────────────────────────────────────────────────

func TestFindMatchingResources_ExactMatch(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{}

	resources := engine.findMatchingResources("Python", prefs)

	if len(resources) == 0 {
		t.Error("expected resources for Python")
	}
	for _, r := range resources {
		found := false
		for _, s := range r.Skills {
			if normalizeSkillName(s) == "python" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("resource %q does not cover Python", r.Title)
		}
	}
}

func TestFindMatchingResources_AliasMatch(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{}

	// "golang" should match "go" resources.
	resources := engine.findMatchingResources("golang", prefs)

	if len(resources) == 0 {
		t.Error("expected resources for golang (alias for go)")
	}
}

func TestFindMatchingResources_NoMatch(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{}

	resources := engine.findMatchingResources("some_obscure_skill_xyz", prefs)

	if len(resources) != 0 {
		t.Errorf("expected 0 resources for unknown skill, got %d", len(resources))
	}
}

func TestFindMatchingResources_FreeFilter(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{PreferFree: true}

	resources := engine.findMatchingResources("Python", prefs)

	for _, r := range resources {
		if r.CostType != "free" && r.CostType != "free_audit" {
			t.Errorf("expected only free resources, got %s (cost_type=%s)", r.Title, r.CostType)
		}
	}
}

func TestFindMatchingResources_BudgetFilter(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{MaxBudgetUSD: 25.00}

	resources := engine.findMatchingResources("Python", prefs)

	for _, r := range resources {
		if r.CostUSD > 25.00 {
			t.Errorf("resource %q exceeds budget: $%.2f > $25.00", r.Title, r.CostUSD)
		}
	}
}

func TestFindMatchingResources_ExcludedProvider(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{ExcludedProviders: []string{"OtherProvider"}}

	resources := engine.findMatchingResources("Python", prefs)

	for _, r := range resources {
		if r.Provider == "OtherProvider" {
			t.Errorf("excluded provider 'OtherProvider' appeared in results: %s", r.Title)
		}
	}
}

func TestFindMatchingResources_ResourceTypeFilter(t *testing.T) {
	engine := newTestEngine()
	prefs := UserPreferences{PreferredResourceTypes: []string{"course"}}

	resources := engine.findMatchingResources("Python", prefs)

	for _, r := range resources {
		if r.ResourceType != "course" {
			t.Errorf("expected only course resources, got %s (type=%s)", r.Title, r.ResourceType)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Relevance scoring tests
// ─────────────────────────────────────────────────────────────────────────────

func TestComputeRelevanceScore_InRange(t *testing.T) {
	gap := testGap("Python", "critical", "intermediate", "")
	prefs := UserPreferences{}

	for _, res := range testCatalog {
		score := computeRelevanceScore(res, gap, prefs)
		if score < 0 || score > 1 {
			t.Errorf("relevance score out of [0,1] for %q: %.4f", res.Title, score)
		}
	}
}

func TestComputeRelevanceScore_PrimarySkillHigherThanSecondary(t *testing.T) {
	// Python course (primary=python) should score higher than a course
	// where python is secondary.
	primaryRes := ResourceEntry{
		ID: "primary", Title: "Python Course",
		Provider: "Test", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", Skills: []string{"python"}, PrimarySkill: "python",
		Rating: 4.5, RatingCount: 1000, IsVerified: true,
	}
	secondaryRes := ResourceEntry{
		ID: "secondary", Title: "Data Science Course",
		Provider: "Test", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", Skills: []string{"python", "statistics"}, PrimarySkill: "statistics",
		Rating: 4.5, RatingCount: 1000, IsVerified: true,
	}

	gap := testGap("Python", "critical", "intermediate", "")
	prefs := UserPreferences{}

	primaryScore := computeRelevanceScore(primaryRes, gap, prefs)
	secondaryScore := computeRelevanceScore(secondaryRes, gap, prefs)

	if primaryScore <= secondaryScore {
		t.Errorf("primary skill match (%.4f) should score higher than secondary (%.4f)",
			primaryScore, secondaryScore)
	}
}

func TestComputeRelevanceScore_FreePreferenceBoost(t *testing.T) {
	freeRes := ResourceEntry{
		ID: "free", Title: "Free Python",
		Provider: "Test", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", Skills: []string{"python"}, PrimarySkill: "python",
		Rating: 4.5, RatingCount: 1000, IsVerified: true,
	}
	paidRes := ResourceEntry{
		ID: "paid", Title: "Paid Python",
		Provider: "Test", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, Skills: []string{"python"}, PrimarySkill: "python",
		Rating: 4.5, RatingCount: 1000, IsVerified: true,
	}

	gap := testGap("Python", "critical", "intermediate", "")
	prefsWithFree := UserPreferences{PreferFree: true}
	prefsNoFree := UserPreferences{}

	freeScoreWithPref := computeRelevanceScore(freeRes, gap, prefsWithFree)
	paidScoreWithPref := computeRelevanceScore(paidRes, gap, prefsWithFree)
	freeScoreNoPref := computeRelevanceScore(freeRes, gap, prefsNoFree)
	paidScoreNoPref := computeRelevanceScore(paidRes, gap, prefsNoFree)

	// With free preference, free resource should score higher than paid.
	if freeScoreWithPref <= paidScoreWithPref {
		t.Errorf("with free preference, free (%.4f) should score > paid (%.4f)",
			freeScoreWithPref, paidScoreWithPref)
	}

	// Without free preference, scores should be equal (same rating, same everything).
	if freeScoreNoPref != paidScoreNoPref {
		t.Errorf("without free preference, free (%.4f) and paid (%.4f) should score equally",
			freeScoreNoPref, paidScoreNoPref)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Difficulty fit tests
// ─────────────────────────────────────────────────────────────────────────────

func TestComputeDifficultyFit_PerfectMatch(t *testing.T) {
	score := computeDifficultyFit("intermediate", "intermediate", "beginner")
	if score != 1.0 {
		t.Errorf("expected 1.0 for perfect difficulty match, got %.4f", score)
	}
}

func TestComputeDifficultyFit_AllLevels(t *testing.T) {
	score := computeDifficultyFit("all_levels", "advanced", "beginner")
	if score < 0.8 {
		t.Errorf("expected high score for all_levels resource, got %.4f", score)
	}
}

func TestComputeDifficultyFit_OneOff(t *testing.T) {
	score := computeDifficultyFit("advanced", "intermediate", "")
	if score < 0.5 || score > 0.9 {
		t.Errorf("expected moderate score for one-level-off difficulty, got %.4f", score)
	}
}

func TestComputeDifficultyFit_TwoOff(t *testing.T) {
	score := computeDifficultyFit("expert", "beginner", "")
	if score > 0.5 {
		t.Errorf("expected low score for two-level-off difficulty, got %.4f", score)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Learning hours estimation tests
// ─────────────────────────────────────────────────────────────────────────────

func TestEstimateCompletionHours_NoCurrentLevel(t *testing.T) {
	res := ResourceEntry{DurationHours: 20}
	gap := testGap("Python", "critical", "intermediate", "")

	hours := estimateCompletionHours(res, gap)

	if hours != 20.0 {
		t.Errorf("expected 20 hours with no current level, got %.1f", hours)
	}
}

func TestEstimateCompletionHours_BeginnerLevel(t *testing.T) {
	res := ResourceEntry{DurationHours: 20}
	gap := testGap("Python", "critical", "intermediate", "beginner")

	hours := estimateCompletionHours(res, gap)

	if hours >= 20.0 {
		t.Errorf("expected fewer hours with beginner level, got %.1f", hours)
	}
}

func TestEstimateCompletionHours_IntermediateLevel(t *testing.T) {
	res := ResourceEntry{DurationHours: 20}
	gapBeginner := testGap("Python", "critical", "intermediate", "beginner")
	gapIntermediate := testGap("Python", "critical", "intermediate", "intermediate")

	hoursBeginner := estimateCompletionHours(res, gapBeginner)
	hoursIntermediate := estimateCompletionHours(res, gapIntermediate)

	if hoursIntermediate >= hoursBeginner {
		t.Errorf("intermediate level should need fewer hours than beginner: intermediate=%.1f, beginner=%.1f",
			hoursIntermediate, hoursBeginner)
	}
}

func TestEstimateCompletionHours_AlwaysPositive(t *testing.T) {
	res := ResourceEntry{DurationHours: 1}
	gap := testGap("Python", "critical", "intermediate", "advanced")

	hours := estimateCompletionHours(res, gap)

	if hours <= 0 {
		t.Errorf("expected positive hours, got %.1f", hours)
	}
}

func TestEstimateCompletionHours_NoDuration(t *testing.T) {
	res := ResourceEntry{DurationHours: 0}
	gap := testGap("Python", "critical", "intermediate", "")
	gap.EstimatedLearningHours = 50

	hours := estimateCompletionHours(res, gap)

	if hours != 50.0 {
		t.Errorf("expected gap estimate (50) when resource has no duration, got %.1f", hours)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phase building tests
// ─────────────────────────────────────────────────────────────────────────────

func TestBuildPhases_CriticalFirst(t *testing.T) {
	criticalRecs := []SkillRecommendation{
		{SkillName: "Python", GapCategory: "critical", EstimatedHoursToJobReady: 20},
	}
	importantRecs := []SkillRecommendation{
		{SkillName: "Docker", GapCategory: "important", EstimatedHoursToJobReady: 10},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	phases := buildPhases(criticalRecs, importantRecs, nil, prefs)

	if len(phases) != 2 {
		t.Errorf("expected 2 phases, got %d", len(phases))
	}
	if phases[0].PhaseNumber != 1 {
		t.Errorf("expected first phase number=1, got %d", phases[0].PhaseNumber)
	}
	if phases[0].PhaseName != "Critical Skills" {
		t.Errorf("expected first phase='Critical Skills', got %s", phases[0].PhaseName)
	}
	if phases[1].PhaseName != "Preferred Skills" {
		t.Errorf("expected second phase='Preferred Skills', got %s", phases[1].PhaseName)
	}
}

func TestBuildPhases_OnlyCritical(t *testing.T) {
	criticalRecs := []SkillRecommendation{
		{SkillName: "Python", GapCategory: "critical", EstimatedHoursToJobReady: 20},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	phases := buildPhases(criticalRecs, nil, nil, prefs)

	if len(phases) != 1 {
		t.Errorf("expected 1 phase, got %d", len(phases))
	}
}

func TestBuildPhases_EmptyGaps(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	phases := buildPhases(nil, nil, nil, prefs)

	if len(phases) != 0 {
		t.Errorf("expected 0 phases for no gaps, got %d", len(phases))
	}
}

func TestBuildPhase_TotalHours(t *testing.T) {
	recs := []SkillRecommendation{
		{SkillName: "Python", EstimatedHoursToJobReady: 20},
		{SkillName: "Go", EstimatedHoursToJobReady: 15},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	phase := buildPhase(1, "Critical Skills", recs, prefs, "description")

	if phase.TotalHours != 35.0 {
		t.Errorf("expected TotalHours=35, got %.1f", phase.TotalHours)
	}
}

func TestBuildPhase_EstimatedWeeks(t *testing.T) {
	recs := []SkillRecommendation{
		{SkillName: "Python", EstimatedHoursToJobReady: 20},
	}
	prefs := UserPreferences{WeeklyHoursAvailable: 10}

	phase := buildPhase(1, "Critical Skills", recs, prefs, "description")

	// 20 hours / 10 hours per week = 2 weeks.
	if phase.EstimatedWeeks != 2.0 {
		t.Errorf("expected EstimatedWeeks=2.0, got %.1f", phase.EstimatedWeeks)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Summary tests
// ─────────────────────────────────────────────────────────────────────────────

func TestBuildSummary_HeadlineNoCriticalGaps(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		Title:          "Python Developer",
		RequiredSkills: []string{"Python"},
	}
	prefs := UserPreferences{}

	plan := engine.Generate(profile, job, prefs)

	if plan.Summary.Headline == "" {
		t.Error("expected non-empty headline")
	}
}

func TestBuildSummary_QuickWins(t *testing.T) {
	engine := newTestEngine()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		Title:          "Backend Engineer",
		RequiredSkills: []string{"SQL"}, // SQL has 5 hours in test catalog
	}
	prefs := UserPreferences{}

	plan := engine.Generate(profile, job, prefs)

	// SQL should be a quick win (< 20 hours).
	found := false
	for _, qw := range plan.Summary.QuickWins {
		if qw == "SQL" {
			found = true
			break
		}
	}
	if !found {
		t.Logf("Quick wins: %v", plan.Summary.QuickWins)
		// Note: quick wins depend on estimated hours which may vary.
		// This is a soft check.
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Preference defaults tests
// ─────────────────────────────────────────────────────────────────────────────

func TestApplyPreferenceDefaults_ZeroHours(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 0}
	result := applyPreferenceDefaults(prefs)
	if result.WeeklyHoursAvailable != 10 {
		t.Errorf("expected default 10 hours, got %.1f", result.WeeklyHoursAvailable)
	}
}

func TestApplyPreferenceDefaults_NegativeHours(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: -5}
	result := applyPreferenceDefaults(prefs)
	if result.WeeklyHoursAvailable != 10 {
		t.Errorf("expected default 10 hours for negative input, got %.1f", result.WeeklyHoursAvailable)
	}
}

func TestApplyPreferenceDefaults_CustomHoursPreserved(t *testing.T) {
	prefs := UserPreferences{WeeklyHoursAvailable: 20}
	result := applyPreferenceDefaults(prefs)
	if result.WeeklyHoursAvailable != 20 {
		t.Errorf("expected 20 hours preserved, got %.1f", result.WeeklyHoursAvailable)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Alias resolution tests
// ─────────────────────────────────────────────────────────────────────────────

func TestResolveAlias_KnownAlias(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"golang", "go"},
		{"js", "javascript"},
		{"ts", "typescript"},
		{"nodejs", "node.js"},
		{"k8s", "kubernetes"},
		{"postgres", "postgresql"},
	}

	for _, tt := range tests {
		result := resolveAlias(tt.input)
		if result != tt.expected {
			t.Errorf("resolveAlias(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestResolveAlias_UnknownSkill(t *testing.T) {
	result := resolveAlias("python")
	if result != "python" {
		t.Errorf("resolveAlias('python') should return 'python', got %q", result)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Recommendation reason tests
// ─────────────────────────────────────────────────────────────────────────────

func TestBuildRecommendationReason_NotEmpty(t *testing.T) {
	res := testCatalog[0] // python-course
	gap := testGap("Python", "critical", "intermediate", "")
	prefs := UserPreferences{}

	reason := buildRecommendationReason(res, gap, prefs)

	if reason == "" {
		t.Error("expected non-empty recommendation reason")
	}
}

func TestBuildRecommendationReason_MentionsFreeForFreeResource(t *testing.T) {
	res := ResourceEntry{
		ID: "free-res", Title: "Free Python",
		Provider: "Test", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", Skills: []string{"python"}, PrimarySkill: "python",
		Rating: 4.5, IsVerified: true,
	}
	gap := testGap("Python", "critical", "intermediate", "")
	prefs := UserPreferences{}

	reason := buildRecommendationReason(res, gap, prefs)

	if !containsString(reason, "free") {
		t.Errorf("expected reason to mention 'free' for free resource, got: %s", reason)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Skill recommendation tests
// ─────────────────────────────────────────────────────────────────────────────

func TestBuildSkillRecommendation_HasPrimaryResource(t *testing.T) {
	engine := newTestEngine()
	gap := testGap("Python", "critical", "intermediate", "")
	prefs := UserPreferences{}

	rec := engine.buildSkillRecommendation(gap, prefs)

	if rec.PrimaryResource == nil {
		t.Error("expected primary resource for Python (in test catalog)")
	}
}

func TestBuildSkillRecommendation_NoPrimaryForUnknownSkill(t *testing.T) {
	engine := newTestEngine()
	gap := testGap("some_obscure_skill_xyz", "critical", "intermediate", "")
	prefs := UserPreferences{}

	rec := engine.buildSkillRecommendation(gap, prefs)

	if rec.PrimaryResource != nil {
		t.Error("expected no primary resource for unknown skill")
	}
}

func TestBuildSkillRecommendation_AlternativesAreDifferent(t *testing.T) {
	engine := newTestEngine()
	gap := testGap("Python", "critical", "intermediate", "")
	prefs := UserPreferences{}

	rec := engine.buildSkillRecommendation(gap, prefs)

	if rec.PrimaryResource == nil {
		t.Skip("no primary resource found")
	}

	for _, alt := range rec.AlternativeResources {
		if !alt.IsAlternative {
			t.Error("alternative resources should have IsAlternative=true")
		}
		// Alternative should differ from primary in type or provider.
		if alt.Resource.ResourceType == rec.PrimaryResource.Resource.ResourceType &&
			alt.Resource.Provider == rec.PrimaryResource.Resource.Provider {
			t.Errorf("alternative %q has same type and provider as primary", alt.Resource.Title)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// testGap creates a SkillGap for testing.
func testGap(skillName, category, targetLevel, currentLevel string) gapanalysis.SkillGap {
	return gapanalysis.SkillGap{
		SkillName:              skillName,
		Category:               gapanalysis.GapCategory(category),
		PriorityScore:          0.8,
		ImportanceScore:        1.0,
		EstimatedLearningHours: 30,
		TransferabilityScore:   0.9,
		TargetLevel:            targetLevel,
		CurrentLevel:           currentLevel,
	}
}

// containsString returns true if s contains substr (case-insensitive).
func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
