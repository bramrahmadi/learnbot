package gapanalysis

import (
	"testing"

	"github.com/learnbot/resume-parser/internal/scorer"
)

// ─────────────────────────────────────────────────────────────────────────────
// Performance benchmarks for gap analysis
// ─────────────────────────────────────────────────────────────────────────────

// BenchmarkAnalyze_EmptyInputs benchmarks gap analysis with empty inputs.
func BenchmarkAnalyze_EmptyInputs(b *testing.B) {
	analyzer := New()
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(profile, job)
	}
}

// BenchmarkAnalyze_ManyGaps benchmarks gap analysis with many missing skills.
func BenchmarkAnalyze_ManyGaps(b *testing.B) {
	analyzer := New()
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Excel", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{
			"Go", "Python", "Java", "JavaScript", "TypeScript",
			"PostgreSQL", "MongoDB", "Redis", "Elasticsearch",
			"Docker", "Kubernetes", "AWS", "Terraform",
		},
		PreferredSkills: []string{
			"React", "Node.js", "TensorFlow", "Kafka",
			"Ansible", "GCP", "Azure",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(profile, job)
	}
}

// BenchmarkAnalyze_NoGaps benchmarks gap analysis when candidate has all skills.
func BenchmarkAnalyze_NoGaps(b *testing.B) {
	analyzer := New()
	skills := []scorer.CandidateSkill{
		{Name: "Go", Proficiency: "expert"},
		{Name: "Python", Proficiency: "advanced"},
		{Name: "PostgreSQL", Proficiency: "advanced"},
		{Name: "Docker", Proficiency: "advanced"},
		{Name: "Kubernetes", Proficiency: "intermediate"},
		{Name: "AWS", Proficiency: "intermediate"},
	}
	profile := scorer.CandidateProfile{Skills: skills}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Go", "Python", "PostgreSQL"},
		PreferredSkills: []string{"Docker", "Kubernetes", "AWS"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(profile, job)
	}
}

// BenchmarkComputePriorityScore benchmarks the priority score computation.
func BenchmarkComputePriorityScore(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		computePriorityScore(1.0, 0.85, 120)
		computePriorityScore(0.6, 0.70, 80)
		computePriorityScore(0.3, 0.95, 30)
	}
}

// BenchmarkTokenOverlapSimilarity benchmarks string similarity computation.
func BenchmarkTokenOverlapSimilarity(b *testing.B) {
	pairs := [][2]string{
		{"machine learning", "deep learning"},
		{"node.js", "javascript"},
		{"spring boot", "spring framework"},
		{"kubernetes", "docker"},
		{"tensorflow", "pytorch"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pair := range pairs {
			tokenOverlapSimilarity(pair[0], pair[1])
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Additional unit tests for gap analysis
// ─────────────────────────────────────────────────────────────────────────────

// TestComputePriorityScore_Formula verifies the priority score formula.
func TestComputePriorityScore_Formula(t *testing.T) {
	// With importance=1.0, transferability=1.0, hours=0 (ease=1.0):
	// score = 1.0*0.50 + 1.0*0.30 + 1.0*0.20 = 1.0
	score := computePriorityScore(1.0, 1.0, 0)
	if !approxEqual(score, 1.0, 0.001) {
		t.Errorf("expected priority score=1.0 for max inputs, got %.4f", score)
	}
}

// TestComputePriorityScore_ZeroInputs verifies zero inputs give zero score.
func TestComputePriorityScore_ZeroInputs(t *testing.T) {
	// With importance=0, transferability=0, hours=maxHours (ease=0):
	// score = 0*0.50 + 0*0.30 + 0*0.20 = 0
	score := computePriorityScore(0.0, 0.0, int(maxLearningHoursForNormalization))
	if !approxEqual(score, 0.0, 0.001) {
		t.Errorf("expected priority score=0.0 for zero inputs, got %.4f", score)
	}
}

// TestComputePriorityScore_InRange verifies score is always in [0, 1].
func TestComputePriorityScore_InRange(t *testing.T) {
	testCases := []struct {
		importance      float64
		transferability float64
		hours           int
	}{
		{1.0, 1.0, 0},
		{0.0, 0.0, 500},
		{0.5, 0.5, 100},
		{1.0, 0.0, 0},
		{0.0, 1.0, 500},
	}

	for _, tc := range testCases {
		score := computePriorityScore(tc.importance, tc.transferability, tc.hours)
		if score < 0 || score > 1 {
			t.Errorf("priority score out of [0,1]: importance=%.2f, transferability=%.2f, hours=%d → %.4f",
				tc.importance, tc.transferability, tc.hours, score)
		}
	}
}

// TestTokenOverlapSimilarity_IdenticalStrings verifies identical strings return 1.0.
func TestTokenOverlapSimilarity_IdenticalStrings(t *testing.T) {
	score := tokenOverlapSimilarity("machine learning", "machine learning")
	if !approxEqual(score, 1.0, 0.001) {
		t.Errorf("expected 1.0 for identical strings, got %.4f", score)
	}
}

// TestTokenOverlapSimilarity_NoOverlap verifies no overlap returns 0.0.
func TestTokenOverlapSimilarity_NoOverlap(t *testing.T) {
	score := tokenOverlapSimilarity("python", "kubernetes")
	if score != 0.0 {
		t.Errorf("expected 0.0 for no overlap, got %.4f", score)
	}
}

// TestTokenOverlapSimilarity_PartialOverlap verifies partial overlap returns intermediate score.
func TestTokenOverlapSimilarity_PartialOverlap(t *testing.T) {
	// "machine learning" and "deep learning" share "learning"
	score := tokenOverlapSimilarity("machine learning", "deep learning")
	if score <= 0 || score >= 1 {
		t.Errorf("expected partial overlap score in (0,1), got %.4f", score)
	}
}

// TestTokenOverlapSimilarity_EmptyStrings verifies empty strings return 0.0.
func TestTokenOverlapSimilarity_EmptyStrings(t *testing.T) {
	score := tokenOverlapSimilarity("", "python")
	if score != 0.0 {
		t.Errorf("expected 0.0 for empty string, got %.4f", score)
	}

	score = tokenOverlapSimilarity("python", "")
	if score != 0.0 {
		t.Errorf("expected 0.0 for empty string, got %.4f", score)
	}
}

// TestAdjustLearningHours_NoSimilarity verifies no reduction without similarity.
func TestAdjustLearningHours_NoSimilarity(t *testing.T) {
	baseHours := 100
	adjusted := adjustLearningHours(baseHours, 0.0, "")
	if adjusted != baseHours {
		t.Errorf("expected no reduction with 0 similarity: base=%d, adjusted=%d",
			baseHours, adjusted)
	}
}

// TestAdjustLearningHours_HighSimilarity verifies reduction with high similarity.
func TestAdjustLearningHours_HighSimilarity(t *testing.T) {
	baseHours := 100
	adjusted := adjustLearningHours(baseHours, 0.7, "")
	if adjusted >= baseHours {
		t.Errorf("expected reduction with high similarity: base=%d, adjusted=%d",
			baseHours, adjusted)
	}
}

// TestAdjustLearningHours_BeginnerLevel verifies reduction for beginner level.
func TestAdjustLearningHours_BeginnerLevel(t *testing.T) {
	baseHours := 100
	adjusted := adjustLearningHours(baseHours, 0.0, "beginner")
	if adjusted >= baseHours {
		t.Errorf("expected reduction for beginner level: base=%d, adjusted=%d",
			baseHours, adjusted)
	}
}

// TestAdjustLearningHours_AdvancedLevel verifies larger reduction for advanced level.
func TestAdjustLearningHours_AdvancedLevel(t *testing.T) {
	baseHours := 100
	adjustedBeginner := adjustLearningHours(baseHours, 0.0, "beginner")
	adjustedAdvanced := adjustLearningHours(baseHours, 0.0, "advanced")
	if adjustedAdvanced >= adjustedBeginner {
		t.Errorf("advanced level should have fewer hours than beginner: beginner=%d, advanced=%d",
			adjustedBeginner, adjustedAdvanced)
	}
}

// TestAdjustLearningHours_MinimumOne verifies minimum of 1 hour.
func TestAdjustLearningHours_MinimumOne(t *testing.T) {
	// Even with maximum reduction, should be at least 1 hour
	adjusted := adjustLearningHours(1, 1.0, "advanced")
	if adjusted < 1 {
		t.Errorf("expected minimum 1 hour, got %d", adjusted)
	}
}

// TestAnalyze_ReadinessScoreInRange verifies readiness score is in [0, 100].
func TestAnalyze_ReadinessScoreInRange(t *testing.T) {
	testCases := []struct {
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
			name: "all skills missing",
			profile: scorer.CandidateProfile{
				Skills: []scorer.CandidateSkill{{Name: "Excel"}},
			},
			job: scorer.JobRequirements{
				RequiredSkills: []string{"Go", "Python", "Kubernetes", "TensorFlow", "Rust"},
			},
		},
		{
			name: "all skills present",
			profile: scorer.CandidateProfile{
				Skills: []scorer.CandidateSkill{
					{Name: "Go"}, {Name: "Python"}, {Name: "Docker"},
				},
			},
			job: scorer.JobRequirements{
				RequiredSkills: []string{"Go", "Python", "Docker"},
			},
		},
	}

	analyzer := New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := analyzer.Analyze(tc.profile, tc.job)
			if result.ReadinessScore < 0 || result.ReadinessScore > 100 {
				t.Errorf("readiness score out of [0, 100]: %.2f", result.ReadinessScore)
			}
		})
	}
}

// TestAnalyze_VisualDataPopulated verifies visual data is populated.
func TestAnalyze_VisualDataPopulated(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills:  []string{"Go", "Python"},
		PreferredSkills: []string{"Docker"},
	}

	result := New().Analyze(profile, job)

	if len(result.VisualData.GapsByCategory) == 0 {
		t.Error("expected GapsByCategory to be populated")
	}
}

// TestAnalyze_DuplicateSkillsDeduped verifies duplicate skills in job are deduped.
func TestAnalyze_DuplicateSkillsDeduped(t *testing.T) {
	profile := scorer.CandidateProfile{}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"Python", "Python", "Python"},
	}

	result := New().Analyze(profile, job)

	if len(result.CriticalGaps) != 1 {
		t.Errorf("expected 1 critical gap (deduped), got %d", len(result.CriticalGaps))
	}
}

// TestAnalyze_CaseInsensitiveSkillMatching verifies case-insensitive skill matching.
func TestAnalyze_CaseInsensitiveSkillMatching(t *testing.T) {
	profile := scorer.CandidateProfile{
		Skills: []scorer.CandidateSkill{
			{Name: "PYTHON", Proficiency: "expert"},
		},
	}
	job := scorer.JobRequirements{
		RequiredSkills: []string{"python"},
	}

	result := New().Analyze(profile, job)

	if len(result.CriticalGaps) != 0 {
		t.Errorf("expected 0 critical gaps (case-insensitive match), got %d: %v",
			len(result.CriticalGaps), result.CriticalGaps)
	}
}
