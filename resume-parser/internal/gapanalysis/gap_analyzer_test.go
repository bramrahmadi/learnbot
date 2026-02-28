package gapanalysis

import (
	"math"
	"testing"

	"github.com/learnbot/resume-parser/internal/scorer"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────────────────────

// approxEqual returns true if a and b differ by less than epsilon.
func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// containsGap returns true if a gap with the given skill name exists in the list.
func containsGap(gaps []SkillGap, skillName string) bool {
	for _, g := range gaps {
		if g.SkillName == skillName {
			return true
		}
	}
	return false
}

// findGap returns the gap with the given skill name, or a zero value.
func findGap(gaps []SkillGap, skillName string) (SkillGap, bool) {
	for _, g := range gaps {
		if g.SkillName == skillName {
			return g, true
		}
	}
	return SkillGap{}, false
}

// containsSkill returns true if the skill name is in the list.
func containsSkill(skills []string, name string) bool {
	for _, s := range skills {
		if s == name {
			return true
		}
	}
	return false
}

// newAnalyzer creates a fresh Analyzer for testing.
func newAnalyzer() *Analyzer {
	return New()
}

// ─────────────────────────────────────────────────────────────────────────────
// Analyze – overall behavior
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_NoCriticalGaps(t *testing.T) {
	// Candidate has all required skills.
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "PostgreSQL", Proficiency: "advanced"},
			{Name: "Docker", Proficiency: "intermediate"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go", "PostgreSQL", "Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) != 0 {
		t.Errorf("expected 0 critical gaps, got %d: %v", len(result.CriticalGaps), result.CriticalGaps)
	}
	if result.CriticalGapCount != 0 {
		t.Errorf("expected CriticalGapCount=0, got %d", result.CriticalGapCount)
	}
}

func TestAnalyze_AllSkillsMissing(t *testing.T) {
	// Candidate has no relevant skills.
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Excel", Proficiency: "beginner"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "TensorFlow", "Kubernetes"},
		PreferredSkills: []string{"Docker", "AWS"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) != 3 {
		t.Errorf("expected 3 critical gaps, got %d", len(result.CriticalGaps))
	}
	if len(result.ImportantGaps) != 2 {
		t.Errorf("expected 2 important gaps, got %d", len(result.ImportantGaps))
	}
	if result.TotalGaps != 5 {
		t.Errorf("expected TotalGaps=5, got %d", result.TotalGaps)
	}
}

func TestAnalyze_PartialSkillMatch(t *testing.T) {
	// Candidate has some but not all required skills.
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "TensorFlow", "PyTorch"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) != 2 {
		t.Errorf("expected 2 critical gaps (TensorFlow, PyTorch), got %d", len(result.CriticalGaps))
	}
	if !containsGap(result.CriticalGaps, "TensorFlow") {
		t.Error("expected TensorFlow in critical gaps")
	}
	if !containsGap(result.CriticalGaps, "PyTorch") {
		t.Error("expected PyTorch in critical gaps")
	}
	if containsGap(result.CriticalGaps, "Python") {
		t.Error("Python should NOT be in critical gaps (candidate has it)")
	}
}

func TestAnalyze_MatchedSkillsPopulated(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Docker", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Go", "Python"},
		PreferredSkills: []string{"Docker", "Kubernetes"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if !containsSkill(result.MatchedSkills, "Go") {
		t.Error("expected Go in matched skills")
	}
	if !containsSkill(result.MatchedSkills, "Docker") {
		t.Error("expected Docker in matched skills")
	}
	if containsSkill(result.MatchedSkills, "Python") {
		t.Error("Python should NOT be in matched skills (candidate lacks it)")
	}
}

func TestAnalyze_EmptyProfile(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go", "Python"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) != 2 {
		t.Errorf("expected 2 critical gaps for empty profile, got %d", len(result.CriticalGaps))
	}
	if result.ReadinessScore > 80 {
		t.Errorf("expected low readiness score for empty profile, got %.2f", result.ReadinessScore)
	}
}

func TestAnalyze_EmptyJob(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) != 0 {
		t.Errorf("expected 0 critical gaps for empty job, got %d", len(result.CriticalGaps))
	}
	if result.ReadinessScore < 90 {
		t.Errorf("expected high readiness score for empty job, got %.2f", result.ReadinessScore)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Gap categories
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_GapCategories(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{},
	}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	// Python should be critical.
	pythonGap, found := findGap(result.CriticalGaps, "Python")
	if !found {
		t.Fatal("expected Python in critical gaps")
	}
	if pythonGap.Category != GapCategoryCritical {
		t.Errorf("expected Python gap category=critical, got %s", pythonGap.Category)
	}

	// Docker should be important.
	dockerGap, found := findGap(result.ImportantGaps, "Docker")
	if !found {
		t.Fatal("expected Docker in important gaps")
	}
	if dockerGap.Category != GapCategoryImportant {
		t.Errorf("expected Docker gap category=important, got %s", dockerGap.Category)
	}
}

func TestAnalyze_CriticalGapsHaveHigherImportanceThanImportant(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) == 0 || len(result.ImportantGaps) == 0 {
		t.Skip("need at least one gap in each category")
	}

	criticalImportance := result.CriticalGaps[0].ImportanceScore
	importantImportance := result.ImportantGaps[0].ImportanceScore

	if criticalImportance <= importantImportance {
		t.Errorf("critical gap importance (%.2f) should be > important gap importance (%.2f)",
			criticalImportance, importantImportance)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Priority scoring
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_PriorityScoreInRange(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go", "Rust", "TensorFlow"},
		PreferredSkills: []string{"Docker", "Kubernetes"},
	}

	result := newAnalyzer().Analyze(profile, job)

	allGaps := append(result.CriticalGaps, result.ImportantGaps...)
	for _, g := range allGaps {
		if g.PriorityScore < 0 || g.PriorityScore > 1 {
			t.Errorf("gap %q has priority score out of [0,1]: %.4f", g.SkillName, g.PriorityScore)
		}
	}
}

func TestAnalyze_CriticalGapsPrioritizedOverImportant(t *testing.T) {
	// A critical gap should generally have a higher priority score than
	// an important gap (since importance weight is 0.50).
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) == 0 || len(result.ImportantGaps) == 0 {
		t.Skip("need gaps in both categories")
	}

	criticalPriority := result.CriticalGaps[0].PriorityScore
	importantPriority := result.ImportantGaps[0].PriorityScore

	if criticalPriority <= importantPriority {
		t.Errorf("critical gap priority (%.4f) should be > important gap priority (%.4f)",
			criticalPriority, importantPriority)
	}
}

func TestAnalyze_GapsSortedByPriorityDescending(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Go", "Rust", "TensorFlow", "Kubernetes"},
	}

	result := newAnalyzer().Analyze(profile, job)

	for i := 1; i < len(result.CriticalGaps); i++ {
		if result.CriticalGaps[i].PriorityScore > result.CriticalGaps[i-1].PriorityScore {
			t.Errorf("critical gaps not sorted by priority: [%d]=%.4f > [%d]=%.4f",
				i, result.CriticalGaps[i].PriorityScore,
				i-1, result.CriticalGaps[i-1].PriorityScore)
		}
	}
}

func TestAnalyze_TopPriorityGapsMaxFive(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Go", "Rust", "TensorFlow", "Kubernetes", "Docker", "AWS"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.TopPriorityGaps) > maxTopPriorityGaps {
		t.Errorf("expected at most %d top priority gaps, got %d",
			maxTopPriorityGaps, len(result.TopPriorityGaps))
	}
}

func TestAnalyze_TopPriorityGapsAreSortedByPriority(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go", "Rust"},
		PreferredSkills: []string{"Docker", "Kubernetes"},
	}

	result := newAnalyzer().Analyze(profile, job)

	for i := 1; i < len(result.TopPriorityGaps); i++ {
		if result.TopPriorityGaps[i].PriorityScore > result.TopPriorityGaps[i-1].PriorityScore {
			t.Errorf("top priority gaps not sorted: [%d]=%.4f > [%d]=%.4f",
				i, result.TopPriorityGaps[i].PriorityScore,
				i-1, result.TopPriorityGaps[i-1].PriorityScore)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Semantic similarity
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_SemanticSimilarityReducesLearningHours(t *testing.T) {
	// Candidate with TensorFlow should need fewer hours to learn PyTorch
	// than a candidate with no ML skills.
	profileWithRelated := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "TensorFlow", Proficiency: "advanced"},
		},
	}
	profileWithoutRelated := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Excel", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"PyTorch"},
	}

	resultWith := newAnalyzer().Analyze(profileWithRelated, job)
	resultWithout := newAnalyzer().Analyze(profileWithoutRelated, job)

	if len(resultWith.CriticalGaps) == 0 || len(resultWithout.CriticalGaps) == 0 {
		t.Skip("need PyTorch gap in both results")
	}

	hoursWithRelated := resultWith.CriticalGaps[0].EstimatedLearningHours
	hoursWithoutRelated := resultWithout.CriticalGaps[0].EstimatedLearningHours

	if hoursWithRelated >= hoursWithoutRelated {
		t.Errorf("expected fewer hours with related skill: with=%d, without=%d",
			hoursWithRelated, hoursWithoutRelated)
	}
}

func TestAnalyze_SemanticSimilarityScoreInRange(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "TensorFlow", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"PyTorch", "Kubernetes"},
	}

	result := newAnalyzer().Analyze(profile, job)

	for _, g := range result.CriticalGaps {
		if g.SemanticSimilarityScore < 0 || g.SemanticSimilarityScore > 1 {
			t.Errorf("gap %q has similarity score out of [0,1]: %.4f",
				g.SkillName, g.SemanticSimilarityScore)
		}
	}
}

func TestAnalyze_ClosestExistingSkillPopulated(t *testing.T) {
	// Candidate has TensorFlow; PyTorch is a related skill.
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "TensorFlow", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"PyTorch"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) == 0 {
		t.Skip("no critical gaps found")
	}

	gap := result.CriticalGaps[0]
	if gap.ClosestExistingSkill == "" {
		t.Error("expected ClosestExistingSkill to be populated when candidate has related skills")
	}
	if gap.SemanticSimilarityScore == 0 {
		t.Error("expected non-zero SemanticSimilarityScore when candidate has related skills")
	}
}

func TestAnalyze_AliasMatchPreventsGap(t *testing.T) {
	// "golang" should match "Go" requirement.
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "golang", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) != 0 {
		t.Errorf("expected 0 critical gaps (golang aliases Go), got %d: %v",
			len(result.CriticalGaps), result.CriticalGaps)
	}
}

func TestAnalyze_AliasMatchVariants(t *testing.T) {
	tests := []struct {
		candidateSkill string
		requiredSkill  string
		expectGap      bool
	}{
		{"golang", "Go", false},
		{"nodejs", "node.js", false},
		{"postgres", "postgresql", false},
		{"k8s", "kubernetes", false},
		{"js", "javascript", false},
		{"reactjs", "react", false},
		{"Excel", "Python", true}, // no alias relationship
	}

	for _, tt := range tests {
		t.Run(tt.candidateSkill+"_vs_"+tt.requiredSkill, func(t *testing.T) {
			profile := scorer.CandidateProfile{
				Skills: []scorer.CandidateSkill{
					{Name: tt.candidateSkill, Proficiency: "expert"},
				},
			}
			job := scorer.JobRequirements{
				RequiredSkills: []string{tt.requiredSkill},
			}
			result := newAnalyzer().Analyze(profile, job)
			hasGap := len(result.CriticalGaps) > 0
			if hasGap != tt.expectGap {
				t.Errorf("expected gap=%v for %q vs %q, got gap=%v",
					tt.expectGap, tt.candidateSkill, tt.requiredSkill, hasGap)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Learning hours estimation
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_LearningHoursPositive(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Go", "Rust", "TensorFlow"},
	}

	result := newAnalyzer().Analyze(profile, job)

	for _, g := range result.CriticalGaps {
		if g.EstimatedLearningHours <= 0 {
			t.Errorf("gap %q has non-positive learning hours: %d",
				g.SkillName, g.EstimatedLearningHours)
		}
	}
}

func TestAnalyze_TotalLearningHoursIsSum(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	expectedTotal := 0
	for _, g := range result.CriticalGaps {
		expectedTotal += g.EstimatedLearningHours
	}
	for _, g := range result.ImportantGaps {
		expectedTotal += g.EstimatedLearningHours
	}

	if result.TotalEstimatedLearningHours != expectedTotal {
		t.Errorf("expected TotalEstimatedLearningHours=%d, got %d",
			expectedTotal, result.TotalEstimatedLearningHours)
	}
}

func TestAdjustLearningHours_HighSimilarityReducesHours(t *testing.T) {
	base := 100
	hoursNoSim := adjustLearningHours(base, 0.0, "")
	hoursHighSim := adjustLearningHours(base, 0.9, "")

	if hoursHighSim >= hoursNoSim {
		t.Errorf("high similarity should reduce hours: highSim=%d, noSim=%d",
			hoursHighSim, hoursNoSim)
	}
}

func TestAdjustLearningHours_CurrentLevelReducesHours(t *testing.T) {
	base := 100
	hoursNone := adjustLearningHours(base, 0.0, "")
	hoursBeginner := adjustLearningHours(base, 0.0, "beginner")
	hoursIntermediate := adjustLearningHours(base, 0.0, "intermediate")
	hoursAdvanced := adjustLearningHours(base, 0.0, "advanced")

	if hoursBeginner >= hoursNone {
		t.Errorf("beginner level should reduce hours: beginner=%d, none=%d",
			hoursBeginner, hoursNone)
	}
	if hoursIntermediate >= hoursBeginner {
		t.Errorf("intermediate should have fewer hours than beginner: intermediate=%d, beginner=%d",
			hoursIntermediate, hoursBeginner)
	}
	if hoursAdvanced >= hoursIntermediate {
		t.Errorf("advanced should have fewer hours than intermediate: advanced=%d, intermediate=%d",
			hoursAdvanced, hoursIntermediate)
	}
}

func TestAdjustLearningHours_AlwaysPositive(t *testing.T) {
	cases := []struct {
		base     int
		sim      float64
		level    string
	}{
		{1, 1.0, "advanced"},
		{10, 0.9, "expert"},
		{0, 0.0, ""},
	}
	for _, c := range cases {
		h := adjustLearningHours(c.base, c.sim, c.level)
		if h <= 0 {
			t.Errorf("adjustLearningHours(%d, %.1f, %q) = %d, want > 0",
				c.base, c.sim, c.level, h)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Readiness score
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_ReadinessScoreInRange(t *testing.T) {
	cases := []struct {
		name    string
		profile scorer.CandidateProfile
		job     scorer.JobRequirements
	}{
		{
			name:    "empty profile and job",
			profile: scorer.CandidateProfile{},
			job:     scorer.JobRequirements{},
		},
		{
			name: "perfect match",
			profile: scorer.CandidateProfile{
				Skills: []scorer.CandidateSkill{
					{Name: "Go", Proficiency: "expert"},
					{Name: "Python", Proficiency: "advanced"},
				},
			},
			job: scorer.JobRequirements{
				RequiredSkills: []string{"Go", "Python"},
			},
		},
		{
			name:    "all skills missing",
			profile: scorer.CandidateProfile{},
			job: scorer.JobRequirements{
				RequiredSkills: []string{"Python", "TensorFlow", "Kubernetes", "AWS", "Docker"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := newAnalyzer().Analyze(tc.profile, tc.job)
			if result.ReadinessScore < 0 || result.ReadinessScore > 100 {
				t.Errorf("readiness score out of [0,100]: %.2f", result.ReadinessScore)
			}
		})
	}
}

func TestAnalyze_ReadinessScoreHigherWithFewerGaps(t *testing.T) {
	// Profile with all required skills should have higher readiness than
	// a profile with no skills.
	profileFull := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	profileEmpty := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go", "Python"},
	}

	resultFull := newAnalyzer().Analyze(profileFull, job)
	resultEmpty := newAnalyzer().Analyze(profileEmpty, job)

	if resultFull.ReadinessScore <= resultEmpty.ReadinessScore {
		t.Errorf("full profile should have higher readiness: full=%.2f, empty=%.2f",
			resultFull.ReadinessScore, resultEmpty.ReadinessScore)
	}
}

func TestAnalyze_PerfectMatchReadinessScore100(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Go"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if result.ReadinessScore != 100.0 {
		t.Errorf("expected readiness score=100 for perfect match, got %.2f", result.ReadinessScore)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Recommendations
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_RecommendationsPopulated(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) == 0 {
		t.Skip("no critical gaps")
	}

	gap := result.CriticalGaps[0]
	if len(gap.Recommendations) == 0 {
		t.Error("expected at least one recommendation for a critical gap")
	}
}

func TestAnalyze_RecommendationsPrioritiesArePositive(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Go", "Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	for _, g := range result.CriticalGaps {
		for _, rec := range g.Recommendations {
			if rec.Priority <= 0 {
				t.Errorf("gap %q has recommendation with non-positive priority: %d",
					g.SkillName, rec.Priority)
			}
			if rec.EstimatedHours <= 0 {
				t.Errorf("gap %q has recommendation with non-positive hours: %d",
					g.SkillName, rec.EstimatedHours)
			}
			if rec.Title == "" {
				t.Errorf("gap %q has recommendation with empty title", g.SkillName)
			}
		}
	}
}

func TestAnalyze_RelatedSkillAddsLeverageRecommendation(t *testing.T) {
	// Candidate with TensorFlow should get a "leverage existing skills"
	// recommendation for PyTorch (since they are related).
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "TensorFlow", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"PyTorch"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.CriticalGaps) == 0 {
		t.Skip("no critical gaps")
	}

	gap := result.CriticalGaps[0]
	foundLeverage := false
	for _, rec := range gap.Recommendations {
		if rec.ResourceType == "documentation" && rec.Priority == 1 {
			foundLeverage = true
			break
		}
	}
	if !foundLeverage {
		t.Error("expected a 'leverage existing skills' recommendation (documentation, priority=1) for PyTorch when candidate has TensorFlow")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Visual data
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_VisualDataRadarChartLabels(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python"},
	}

	result := newAnalyzer().Analyze(profile, job)

	radar := result.VisualData.RadarChart
	if len(radar.Labels) == 0 {
		t.Error("expected non-empty radar chart labels")
	}
	if len(radar.CandidateScores) != len(radar.Labels) {
		t.Errorf("radar chart labels (%d) and candidate scores (%d) length mismatch",
			len(radar.Labels), len(radar.CandidateScores))
	}
	if len(radar.RequiredScores) != len(radar.Labels) {
		t.Errorf("radar chart labels (%d) and required scores (%d) length mismatch",
			len(radar.Labels), len(radar.RequiredScores))
	}
}

func TestAnalyze_VisualDataRadarScoresInRange(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Python", Proficiency: "advanced"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	for i, score := range result.VisualData.RadarChart.CandidateScores {
		if score < 0 || score > 1 {
			t.Errorf("radar candidate score[%d] out of [0,1]: %.4f", i, score)
		}
	}
	for i, score := range result.VisualData.RadarChart.RequiredScores {
		if score < 0 || score > 1 {
			t.Errorf("radar required score[%d] out of [0,1]: %.4f", i, score)
		}
	}
}

func TestAnalyze_VisualDataCategorySummary(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if len(result.VisualData.GapsByCategory) == 0 {
		t.Error("expected non-empty GapsByCategory")
	}

	// Find the critical category summary.
	var criticalSummary *CategorySummary
	for i := range result.VisualData.GapsByCategory {
		if result.VisualData.GapsByCategory[i].Category == "critical" {
			criticalSummary = &result.VisualData.GapsByCategory[i]
			break
		}
	}
	if criticalSummary == nil {
		t.Fatal("expected a 'critical' category in GapsByCategory")
	}
	if criticalSummary.GapCount != 2 {
		t.Errorf("expected critical GapCount=2, got %d", criticalSummary.GapCount)
	}
}

func TestAnalyze_VisualDataLearningTimeline(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	timeline := result.VisualData.LearningTimeline
	if len(timeline) == 0 {
		t.Error("expected non-empty learning timeline")
	}

	// Verify cumulative hours are monotonically increasing.
	for i := 1; i < len(timeline); i++ {
		if timeline[i].CumulativeHours < timeline[i-1].CumulativeHours {
			t.Errorf("cumulative hours not increasing: [%d]=%d < [%d]=%d",
				i, timeline[i].CumulativeHours,
				i-1, timeline[i-1].CumulativeHours)
		}
	}

	// Verify order is sequential.
	for i, entry := range timeline {
		if entry.Order != i+1 {
			t.Errorf("timeline entry[%d] has Order=%d, expected %d",
				i, entry.Order, i+1)
		}
	}
}

func TestAnalyze_VisualDataTimelineCriticalBeforeImportant(t *testing.T) {
	// Critical gaps should appear before important gaps in the timeline.
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python"},
		PreferredSkills: []string{"Docker"},
	}

	result := newAnalyzer().Analyze(profile, job)

	timeline := result.VisualData.LearningTimeline
	if len(timeline) < 2 {
		t.Skip("need at least 2 timeline entries")
	}

	// Find positions of critical and important entries.
	criticalPos := -1
	importantPos := -1
	for i, entry := range timeline {
		if entry.Category == string(GapCategoryCritical) && criticalPos == -1 {
			criticalPos = i
		}
		if entry.Category == string(GapCategoryImportant) && importantPos == -1 {
			importantPos = i
		}
	}

	if criticalPos != -1 && importantPos != -1 && criticalPos > importantPos {
		t.Errorf("critical gap (pos=%d) should appear before important gap (pos=%d) in timeline",
			criticalPos, importantPos)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Priority score formula
// ─────────────────────────────────────────────────────────────────────────────

func TestComputePriorityScore_Formula(t *testing.T) {
	// Verify the formula: priority = importance*0.50 + transferability*0.30 + ease*0.20
	importance := 1.0
	transferability := 0.9
	hours := 100 // ease = 1 - 100/500 = 0.8

	expected := importance*0.50 + transferability*0.30 + 0.8*0.20
	got := computePriorityScore(importance, transferability, hours)

	if !approxEqual(got, expected, 0.001) {
		t.Errorf("expected priority score %.4f, got %.4f", expected, got)
	}
}

func TestComputePriorityScore_ClampedToOne(t *testing.T) {
	score := computePriorityScore(1.0, 1.0, 0)
	if score > 1.0 {
		t.Errorf("priority score should be clamped to 1.0, got %.4f", score)
	}
}

func TestComputePriorityScore_ClampedToZero(t *testing.T) {
	score := computePriorityScore(0.0, 0.0, int(maxLearningHoursForNormalization)*10)
	if score < 0.0 {
		t.Errorf("priority score should be clamped to 0.0, got %.4f", score)
	}
}

func TestComputePriorityScore_HighHoursReducesPriority(t *testing.T) {
	// Same importance and transferability, but different hours.
	scoreEasy := computePriorityScore(0.8, 0.8, 10)
	scoreHard := computePriorityScore(0.8, 0.8, 400)

	if scoreEasy <= scoreHard {
		t.Errorf("easy skill (10h) should have higher priority than hard skill (400h): easy=%.4f, hard=%.4f",
			scoreEasy, scoreHard)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Skill metadata
// ─────────────────────────────────────────────────────────────────────────────

func TestGetSkillMetadata_KnownSkill(t *testing.T) {
	meta := getSkillMetadata("python")
	if meta.baseHours <= 0 {
		t.Error("expected positive base hours for python")
	}
	if meta.transferability <= 0 {
		t.Error("expected positive transferability for python")
	}
}

func TestGetSkillMetadata_AliasResolution(t *testing.T) {
	// "golang" should resolve to "go" metadata.
	metaGolang := getSkillMetadata("golang")
	metaGo := getSkillMetadata("go")

	if metaGolang.baseHours != metaGo.baseHours {
		t.Errorf("golang and go should have same metadata: golang=%d, go=%d",
			metaGolang.baseHours, metaGo.baseHours)
	}
}

func TestGetSkillMetadata_UnknownSkillUsesDefault(t *testing.T) {
	meta := getSkillMetadata("some_very_obscure_skill_xyz_123")
	if meta.baseHours != defaultSkillMetadata.baseHours {
		t.Errorf("unknown skill should use default hours: got %d, want %d",
			meta.baseHours, defaultSkillMetadata.baseHours)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Semantic similarity helpers
// ─────────────────────────────────────────────────────────────────────────────

func TestComputeSemanticSimilarity_SameSkill(t *testing.T) {
	meta := getSkillMetadata("python")
	sim := computeSemanticSimilarity("python", "python", meta)
	if sim != 1.0 {
		t.Errorf("same skill should have similarity=1.0, got %.4f", sim)
	}
}

func TestComputeSemanticSimilarity_AliasSkill(t *testing.T) {
	meta := getSkillMetadata("go")
	sim := computeSemanticSimilarity("go", "golang", meta)
	if sim < 0.9 {
		t.Errorf("alias skills should have high similarity: got %.4f", sim)
	}
}

func TestComputeSemanticSimilarity_RelatedSkill(t *testing.T) {
	meta := getSkillMetadata("pytorch")
	// TensorFlow is in PyTorch's related skills.
	sim := computeSemanticSimilarity("pytorch", "tensorflow", meta)
	if sim < 0.5 {
		t.Errorf("related skills should have moderate similarity: got %.4f", sim)
	}
}

func TestComputeSemanticSimilarity_UnrelatedSkill(t *testing.T) {
	meta := getSkillMetadata("python")
	sim := computeSemanticSimilarity("python", "excel", meta)
	if sim > 0.5 {
		t.Errorf("unrelated skills should have low similarity: got %.4f", sim)
	}
}

func TestComputeSemanticSimilarity_InRange(t *testing.T) {
	pairs := [][2]string{
		{"python", "r"},
		{"react", "vue"},
		{"docker", "kubernetes"},
		{"aws", "gcp"},
		{"python", "excel"},
	}
	for _, pair := range pairs {
		meta := getSkillMetadata(pair[0])
		sim := computeSemanticSimilarity(pair[0], pair[1], meta)
		if sim < 0 || sim > 1 {
			t.Errorf("similarity(%q, %q) = %.4f, want in [0,1]", pair[0], pair[1], sim)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Gap counts and totals
// ─────────────────────────────────────────────────────────────────────────────

func TestAnalyze_GapCountsConsistent(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Python", "Go", "Rust"},
		PreferredSkills: []string{"Docker", "Kubernetes"},
	}

	result := newAnalyzer().Analyze(profile, job)

	if result.CriticalGapCount != len(result.CriticalGaps) {
		t.Errorf("CriticalGapCount=%d != len(CriticalGaps)=%d",
			result.CriticalGapCount, len(result.CriticalGaps))
	}
	if result.ImportantGapCount != len(result.ImportantGaps) {
		t.Errorf("ImportantGapCount=%d != len(ImportantGaps)=%d",
			result.ImportantGapCount, len(result.ImportantGaps))
	}
	if result.TotalGaps != result.CriticalGapCount+result.ImportantGapCount+result.NiceToHaveGapCount {
		t.Errorf("TotalGaps=%d != critical(%d)+important(%d)+niceToHave(%d)",
			result.TotalGaps, result.CriticalGapCount, result.ImportantGapCount, result.NiceToHaveGapCount)
	}
}

func TestAnalyze_NoDuplicateGaps(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Python", "Go"}, // duplicate Python
	}

	result := newAnalyzer().Analyze(profile, job)

	// Should have at most 2 unique gaps (Python, Go).
	seen := map[string]bool{}
	for _, g := range result.CriticalGaps {
		if seen[g.SkillName] {
			t.Errorf("duplicate gap found: %q", g.SkillName)
		}
		seen[g.SkillName] = true
	}
}
