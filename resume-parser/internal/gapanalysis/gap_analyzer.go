// Package gapanalysis implements the skill gap analysis engine.
package gapanalysis

import (
	"math"
	"sort"
	"strings"

	"github.com/learnbot/resume-parser/internal/scorer"
)

// ─────────────────────────────────────────────────────────────────────────────
// Priority scoring weights
// ─────────────────────────────────────────────────────────────────────────────

const (
	// weightImportance is the weight given to how critical the skill is
	// to job acceptance in the priority score formula.
	weightImportance = 0.50

	// weightTransferability is the weight given to how broadly useful
	// the skill is across different roles.
	weightTransferability = 0.30

	// weightAcquisitionEase is the weight given to how quickly the skill
	// can be acquired (inverse of learning hours, normalized).
	weightAcquisitionEase = 0.20

	// importanceScoreCritical is the importance score for critical gaps.
	importanceScoreCritical = 1.0

	// importanceScoreImportant is the importance score for important gaps.
	importanceScoreImportant = 0.6

	// importanceScoreNiceToHave is the importance score for nice-to-have gaps.
	importanceScoreNiceToHave = 0.3

	// maxTopPriorityGaps is the maximum number of gaps in the TopPriorityGaps list.
	maxTopPriorityGaps = 5

	// maxLearningHoursForNormalization is used to normalize learning hours
	// into the [0, 1] range for priority scoring.
	maxLearningHoursForNormalization = 500.0
)

// ─────────────────────────────────────────────────────────────────────────────
// Skill metadata database
// ─────────────────────────────────────────────────────────────────────────────

// skillMetadata holds pre-defined metadata for common skills.
type skillMetadata struct {
	// baseHours is the estimated hours to reach job-ready proficiency
	// from zero knowledge.
	baseHours int

	// transferability is how broadly useful this skill is [0.0, 1.0].
	transferability float64

	// difficulty is the inherent difficulty of the skill.
	difficulty DifficultyLevel

	// targetLevel is the typical proficiency level required by employers.
	targetLevel string

	// relatedSkills lists skill names that are semantically similar.
	relatedSkills []string
}

// builtinSkillMetadata provides metadata for common skills.
// Skills not in this map get default values.
var builtinSkillMetadata = map[string]skillMetadata{
	// Programming languages
	"python":     {baseHours: 120, transferability: 0.95, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"r", "julia", "ruby"}},
	"go":         {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"rust", "c", "java"}},
	"golang":     {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"rust", "c", "java"}},
	"rust":       {baseHours: 250, transferability: 0.75, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"c", "c++", "go"}},
	"java":       {baseHours: 160, transferability: 0.90, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"kotlin", "scala", "c#"}},
	"kotlin":     {baseHours: 120, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"java", "scala"}},
	"scala":      {baseHours: 200, transferability: 0.75, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"java", "kotlin"}},
	"javascript": {baseHours: 100, transferability: 0.95, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"typescript", "node.js"}},
	"typescript": {baseHours: 80, transferability: 0.90, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"javascript"}},
	"c++":        {baseHours: 300, transferability: 0.80, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"c", "rust"}},
	"c":          {baseHours: 200, transferability: 0.75, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"c++", "rust"}},
	"c#":         {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"java", ".net"}},
	"ruby":       {baseHours: 100, transferability: 0.75, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"python", "rails"}},
	"php":        {baseHours: 80, transferability: 0.70, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"laravel"}},
	"swift":      {baseHours: 150, transferability: 0.65, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"objective-c", "kotlin"}},
	"r":          {baseHours: 100, transferability: 0.70, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"python", "julia"}},

	// Frameworks & libraries
	"react":          {baseHours: 80, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"vue", "angular", "javascript"}},
	"vue":            {baseHours: 70, transferability: 0.80, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"react", "angular"}},
	"angular":        {baseHours: 100, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"react", "vue"}},
	"node.js":        {baseHours: 80, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"javascript", "express"}},
	"django":         {baseHours: 80, transferability: 0.75, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"flask", "python"}},
	"flask":          {baseHours: 50, transferability: 0.70, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"django", "python"}},
	"spring":         {baseHours: 120, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"java", "spring boot"}},
	"spring boot":    {baseHours: 100, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"java", "spring"}},
	"tensorflow":     {baseHours: 150, transferability: 0.80, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"pytorch", "keras", "python"}},
	"pytorch":        {baseHours: 150, transferability: 0.80, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"tensorflow", "python"}},
	"keras":          {baseHours: 80, transferability: 0.75, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"tensorflow", "pytorch"}},

	// Databases
	"postgresql": {baseHours: 80, transferability: 0.90, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"mysql", "sqlite", "sql"}},
	"mysql":      {baseHours: 70, transferability: 0.85, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"postgresql", "sql"}},
	"mongodb":    {baseHours: 60, transferability: 0.80, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"redis", "cassandra"}},
	"redis":      {baseHours: 40, transferability: 0.85, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"memcached", "mongodb"}},
	"sql":        {baseHours: 60, transferability: 0.95, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"postgresql", "mysql"}},
	"elasticsearch": {baseHours: 80, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"opensearch", "solr"}},
	"cassandra":  {baseHours: 100, transferability: 0.70, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"mongodb", "dynamodb"}},

	// Cloud & DevOps
	"docker":     {baseHours: 60, transferability: 0.95, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"kubernetes", "podman"}},
	"kubernetes": {baseHours: 120, transferability: 0.90, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"docker", "helm"}},
	"aws":        {baseHours: 150, transferability: 0.90, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"gcp", "azure"}},
	"gcp":        {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"aws", "azure"}},
	"azure":      {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"aws", "gcp"}},
	"terraform":  {baseHours: 80, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"ansible", "pulumi"}},
	"ansible":    {baseHours: 60, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"terraform", "chef"}},
	"jenkins":    {baseHours: 60, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"github actions", "gitlab ci"}},
	"git":        {baseHours: 30, transferability: 0.99, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"github", "gitlab"}},
	"linux":      {baseHours: 80, transferability: 0.95, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"unix", "bash"}},
	"bash":       {baseHours: 40, transferability: 0.90, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"shell", "linux"}},

	// ML/AI
	"machine learning":  {baseHours: 200, transferability: 0.85, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"deep learning", "python"}},
	"deep learning":     {baseHours: 250, transferability: 0.80, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"machine learning", "tensorflow"}},
	"nlp":               {baseHours: 200, transferability: 0.75, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"machine learning", "python"}},
	"computer vision":   {baseHours: 200, transferability: 0.75, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"deep learning", "opencv"}},
	"data science":      {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"python", "machine learning"}},
	"data engineering":  {baseHours: 150, transferability: 0.85, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"spark", "kafka", "python"}},
	"spark":             {baseHours: 120, transferability: 0.80, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"hadoop", "kafka"}},
	"kafka":             {baseHours: 80, transferability: 0.80, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"rabbitmq", "spark"}},

	// Soft skills
	"communication":     {baseHours: 40, transferability: 1.0, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"presentation", "writing"}},
	"leadership":        {baseHours: 80, transferability: 1.0, difficulty: DifficultyIntermediate, targetLevel: "intermediate", relatedSkills: []string{"management", "mentoring"}},
	"agile":             {baseHours: 30, transferability: 0.95, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"scrum", "kanban"}},
	"scrum":             {baseHours: 20, transferability: 0.90, difficulty: DifficultyBeginner, targetLevel: "intermediate", relatedSkills: []string{"agile", "kanban"}},
	"system design":     {baseHours: 150, transferability: 0.90, difficulty: DifficultyAdvanced, targetLevel: "intermediate", relatedSkills: []string{"architecture", "distributed systems"}},
}

// defaultSkillMetadata is used for skills not in builtinSkillMetadata.
var defaultSkillMetadata = skillMetadata{
	baseHours:       100,
	transferability: 0.70,
	difficulty:      DifficultyIntermediate,
	targetLevel:     "intermediate",
}

// ─────────────────────────────────────────────────────────────────────────────
// Analyzer
// ─────────────────────────────────────────────────────────────────────────────

// Analyzer performs skill gap analysis between a candidate profile and job requirements.
type Analyzer struct{}

// New creates a new gap Analyzer.
func New() *Analyzer {
	return &Analyzer{}
}

// Analyze computes the skill gap analysis for a candidate against a job.
// It returns a GapAnalysisResult with prioritized gaps and recommendations.
func (a *Analyzer) Analyze(profile scorer.CandidateProfile, job scorer.JobRequirements) GapAnalysisResult {
	// Build a normalized index of candidate skills for fast lookup.
	candidateIndex := buildCandidateIndex(profile.Skills)

	// Identify gaps in each category.
	criticalGaps := a.identifyGaps(job.RequiredSkills, GapCategoryCritical, candidateIndex, profile.Skills)
	importantGaps := a.identifyGaps(job.PreferredSkills, GapCategoryImportant, candidateIndex, profile.Skills)

	// Collect matched skills (required + preferred that the candidate has).
	matchedSkills := collectMatchedSkills(job.RequiredSkills, job.PreferredSkills, candidateIndex)

	// Sort each category by priority score descending.
	sortGapsByPriority(criticalGaps)
	sortGapsByPriority(importantGaps)

	// Build top priority gaps across all categories.
	allGaps := append(append([]SkillGap{}, criticalGaps...), importantGaps...)
	sortGapsByPriority(allGaps)
	topPriority := topN(allGaps, maxTopPriorityGaps)

	// Calculate total learning hours.
	totalHours := sumLearningHours(criticalGaps) + sumLearningHours(importantGaps)

	// Calculate readiness score: 100 - penalty for critical gaps.
	readinessScore := calculateReadinessScore(criticalGaps, importantGaps, job)

	// Build visual data.
	visualData := buildVisualData(criticalGaps, importantGaps, profile, job)

	return GapAnalysisResult{
		CriticalGaps:                criticalGaps,
		ImportantGaps:               importantGaps,
		NiceToHaveGaps:              []SkillGap{}, // populated from context in future
		TotalGaps:                   len(criticalGaps) + len(importantGaps),
		CriticalGapCount:            len(criticalGaps),
		ImportantGapCount:           len(importantGaps),
		NiceToHaveGapCount:          0,
		TotalEstimatedLearningHours: totalHours,
		ReadinessScore:              readinessScore,
		TopPriorityGaps:             topPriority,
		MatchedSkills:               matchedSkills,
		VisualData:                  visualData,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Gap identification
// ─────────────────────────────────────────────────────────────────────────────

// identifyGaps finds missing skills from a list of required/preferred skills.
// Duplicate skill names in the input list are deduplicated.
func (a *Analyzer) identifyGaps(
	requiredSkills []string,
	category GapCategory,
	candidateIndex map[string]scorer.CandidateSkill,
	allCandidateSkills []scorer.CandidateSkill,
) []SkillGap {
	var gaps []SkillGap
	// Track already-processed normalized skill names to avoid duplicates.
	seen := make(map[string]bool)

	for _, skillName := range requiredSkills {
		norm := normalizeSkill(skillName)
		// Skip duplicates in the input list.
		if seen[norm] {
			continue
		}
		seen[norm] = true

		if _, found := lookupInIndex(norm, candidateIndex); found {
			continue // Candidate already has this skill.
		}

		// Candidate is missing this skill – build a gap entry.
		meta := getSkillMetadata(norm)
		importanceScore := categoryImportanceScore(category)

		// Find the closest existing skill for semantic similarity.
		closestSkill, simScore := findClosestSkill(norm, allCandidateSkills)

		// Adjust learning hours based on semantic similarity
		// (if candidate has a related skill, reduce hours).
		adjustedHours := adjustLearningHours(meta.baseHours, simScore, "")

		// Compute priority score.
		priorityScore := computePriorityScore(importanceScore, meta.transferability, adjustedHours)

		gap := SkillGap{
			SkillName:               skillName,
			Category:                category,
			PriorityScore:           roundTo4(priorityScore),
			ImportanceScore:         importanceScore,
			EstimatedLearningHours:  adjustedHours,
			TransferabilityScore:    meta.transferability,
			TargetLevel:             meta.targetLevel,
			SemanticSimilarityScore: roundTo4(simScore),
			ClosestExistingSkill:    closestSkill,
			Difficulty:              meta.difficulty,
			Recommendations:         buildRecommendations(skillName, meta, simScore),
		}

		gaps = append(gaps, gap)
	}

	return gaps
}

// ─────────────────────────────────────────────────────────────────────────────
// Priority scoring
// ─────────────────────────────────────────────────────────────────────────────

// computePriorityScore calculates the composite priority score for a gap.
//
// Formula:
//
//	priority = (importance * 0.50) + (transferability * 0.30) + (acquisition_ease * 0.20)
//
// where acquisition_ease = 1 - (hours / maxHours), clamped to [0, 1].
func computePriorityScore(importance, transferability float64, learningHours int) float64 {
	acquisitionEase := 1.0 - math.Min(1.0, float64(learningHours)/maxLearningHoursForNormalization)
	score := importance*weightImportance +
		transferability*weightTransferability +
		acquisitionEase*weightAcquisitionEase
	return math.Min(1.0, math.Max(0.0, score))
}

// categoryImportanceScore returns the importance score for a gap category.
func categoryImportanceScore(category GapCategory) float64 {
	switch category {
	case GapCategoryCritical:
		return importanceScoreCritical
	case GapCategoryImportant:
		return importanceScoreImportant
	case GapCategoryNiceToHave:
		return importanceScoreNiceToHave
	default:
		return importanceScoreNiceToHave
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Semantic similarity
// ─────────────────────────────────────────────────────────────────────────────

// findClosestSkill finds the candidate's skill most semantically similar to
// the target skill. Returns the skill name and similarity score [0, 1].
func findClosestSkill(targetNorm string, candidateSkills []scorer.CandidateSkill) (string, float64) {
	if len(candidateSkills) == 0 {
		return "", 0.0
	}

	targetMeta := getSkillMetadata(targetNorm)
	bestScore := 0.0
	bestSkill := ""

	for _, cs := range candidateSkills {
		csNorm := normalizeSkill(cs.Name)
		sim := computeSemanticSimilarity(targetNorm, csNorm, targetMeta)
		if sim > bestScore {
			bestScore = sim
			bestSkill = cs.Name
		}
	}

	return bestSkill, bestScore
}

// computeSemanticSimilarity computes the semantic similarity between two skills.
//
// Algorithm:
//  1. Check if they are the same skill (exact/alias match) → 1.0
//  2. Check if the candidate skill is in the target's related skills list → 0.7
//  3. Check if the target skill is in the candidate's related skills list → 0.7
//  4. Check for shared category/domain via metadata → 0.4
//  5. Use string similarity as a fallback → scaled score
func computeSemanticSimilarity(targetNorm, candidateNorm string, targetMeta skillMetadata) float64 {
	if targetNorm == candidateNorm {
		return 1.0
	}

	// Check alias equivalence.
	if skillsAreAliases(targetNorm, candidateNorm) {
		return 0.95
	}

	// Check if candidate skill is in target's related skills.
	for _, rel := range targetMeta.relatedSkills {
		if normalizeSkill(rel) == candidateNorm || skillsAreAliases(normalizeSkill(rel), candidateNorm) {
			return 0.70
		}
	}

	// Check if target is in candidate's related skills.
	candidateMeta := getSkillMetadata(candidateNorm)
	for _, rel := range candidateMeta.relatedSkills {
		if normalizeSkill(rel) == targetNorm || skillsAreAliases(normalizeSkill(rel), targetNorm) {
			return 0.65
		}
	}

	// Check shared difficulty/domain as a weak signal.
	if targetMeta.difficulty == candidateMeta.difficulty && targetMeta.transferability > 0.8 {
		return 0.25
	}

	// String-based similarity as a last resort.
	strSim := tokenOverlapSimilarity(targetNorm, candidateNorm)
	return strSim * 0.5
}

// tokenOverlapSimilarity computes a simple token overlap similarity [0, 1].
func tokenOverlapSimilarity(a, b string) float64 {
	tokensA := strings.Fields(a)
	tokensB := strings.Fields(b)
	if len(tokensA) == 0 || len(tokensB) == 0 {
		return 0.0
	}
	setA := make(map[string]bool, len(tokensA))
	for _, t := range tokensA {
		setA[t] = true
	}
	intersection := 0
	for _, t := range tokensB {
		if setA[t] {
			intersection++
		}
	}
	union := len(tokensA) + len(tokensB) - intersection
	if union == 0 {
		return 0.0
	}
	return float64(intersection) / float64(union)
}

// ─────────────────────────────────────────────────────────────────────────────
// Learning hours estimation
// ─────────────────────────────────────────────────────────────────────────────

// adjustLearningHours adjusts the base learning hours based on:
//   - Semantic similarity to existing skills (higher sim → fewer hours)
//   - Candidate's current level in the skill (if partial knowledge)
func adjustLearningHours(baseHours int, simScore float64, currentLevel string) int {
	hours := float64(baseHours)

	// Reduce hours if candidate has a related skill.
	// At sim=1.0 (same skill), hours → 0 (already knows it).
	// At sim=0.7 (related skill), hours → 30% reduction.
	// At sim=0.0 (unrelated), no reduction.
	if simScore > 0 {
		reduction := simScore * 0.40 // max 40% reduction from related skills
		hours = hours * (1.0 - reduction)
	}

	// Reduce hours if candidate has partial knowledge.
	switch strings.ToLower(currentLevel) {
	case "beginner":
		hours *= 0.70
	case "intermediate":
		hours *= 0.40
	case "advanced":
		hours *= 0.15
	}

	return int(math.Max(1, math.Round(hours)))
}

// ─────────────────────────────────────────────────────────────────────────────
// Recommendations
// ─────────────────────────────────────────────────────────────────────────────

// buildRecommendations generates actionable recommendations for a skill gap.
func buildRecommendations(skillName string, meta skillMetadata, simScore float64) []Recommendation {
	var recs []Recommendation
	priority := 1

	// If candidate has a related skill, suggest bridging first.
	if simScore >= 0.5 {
		recs = append(recs, Recommendation{
			Title:          "Leverage your existing related skills",
			Description:    "You already have related knowledge. Focus on the differences and new concepts specific to " + skillName + ".",
			ResourceType:   "documentation",
			EstimatedHours: int(math.Max(1, float64(meta.baseHours)*0.15)),
			Priority:       priority,
		})
		priority++
	}

	// Primary learning resource based on difficulty.
	switch meta.difficulty {
	case DifficultyBeginner:
		recs = append(recs, Recommendation{
			Title:          "Complete an introductory course on " + skillName,
			Description:    "Start with a structured beginner course to build foundational knowledge. Platforms like Coursera, Udemy, or official documentation are excellent starting points.",
			ResourceType:   "course",
			EstimatedHours: int(math.Max(1, float64(meta.baseHours)*0.40)),
			Priority:       priority,
		})
	case DifficultyIntermediate:
		recs = append(recs, Recommendation{
			Title:          "Take an intermediate-level course on " + skillName,
			Description:    "Enroll in a structured course covering core concepts and practical applications. Look for project-based courses that include hands-on exercises.",
			ResourceType:   "course",
			EstimatedHours: int(math.Max(1, float64(meta.baseHours)*0.40)),
			Priority:       priority,
		})
	case DifficultyAdvanced, DifficultyExpert:
		recs = append(recs, Recommendation{
			Title:          "Study " + skillName + " through official documentation and advanced resources",
			Description:    "This is an advanced skill. Start with official documentation, then progress to advanced courses or books. Consider mentorship from an expert.",
			ResourceType:   "documentation",
			EstimatedHours: int(math.Max(1, float64(meta.baseHours)*0.30)),
			Priority:       priority,
		})
	}
	priority++

	// Hands-on practice recommendation.
	recs = append(recs, Recommendation{
		Title:          "Build a hands-on project using " + skillName,
		Description:    "Apply your learning by building a real project. This solidifies understanding and creates portfolio evidence of your skills.",
		ResourceType:   "project",
		EstimatedHours: int(math.Max(1, float64(meta.baseHours)*0.35)),
		Priority:       priority,
	})
	priority++

	// Certification recommendation for high-transferability skills.
	if meta.transferability >= 0.85 {
		recs = append(recs, Recommendation{
			Title:          "Obtain a recognized certification in " + skillName,
			Description:    "A certification validates your skills to employers and demonstrates commitment. Look for industry-recognized certifications.",
			ResourceType:   "certification",
			EstimatedHours: int(math.Max(1, float64(meta.baseHours)*0.25)),
			Priority:       priority,
		})
	}

	return recs
}

// ─────────────────────────────────────────────────────────────────────────────
// Readiness score
// ─────────────────────────────────────────────────────────────────────────────

// calculateReadinessScore computes a readiness score [0, 100] for the candidate.
//
// Algorithm:
//   - Start at 100.
//   - Deduct points for each critical gap (weighted by priority score).
//   - Deduct smaller points for important gaps.
//   - Clamp to [0, 100].
func calculateReadinessScore(criticalGaps, importantGaps []SkillGap, job scorer.JobRequirements) float64 {
	score := 100.0

	// Each critical gap deducts up to 20 points (scaled by priority).
	for _, g := range criticalGaps {
		deduction := g.PriorityScore * 20.0
		score -= deduction
	}

	// Each important gap deducts up to 8 points.
	for _, g := range importantGaps {
		deduction := g.PriorityScore * 8.0
		score -= deduction
	}

	return math.Max(0, math.Min(100, roundTo2(score)))
}

// ─────────────────────────────────────────────────────────────────────────────
// Visual data
// ─────────────────────────────────────────────────────────────────────────────

// buildVisualData constructs the visual representation of the gap analysis.
func buildVisualData(
	criticalGaps, importantGaps []SkillGap,
	profile scorer.CandidateProfile,
	job scorer.JobRequirements,
) GapVisualData {
	// Build radar chart data by skill category.
	radarData := buildRadarChart(criticalGaps, importantGaps, profile, job)

	// Build category summary.
	categorySummary := []CategorySummary{
		buildCategorySummary("critical", criticalGaps),
		buildCategorySummary("important", importantGaps),
	}

	// Build learning timeline.
	timeline := buildLearningTimeline(criticalGaps, importantGaps)

	return GapVisualData{
		RadarChart:      radarData,
		GapsByCategory:  categorySummary,
		LearningTimeline: timeline,
	}
}

// buildRadarChart creates radar chart data from gap analysis results.
func buildRadarChart(
	criticalGaps, importantGaps []SkillGap,
	profile scorer.CandidateProfile,
	job scorer.JobRequirements,
) RadarChartData {
	// Define skill categories for the radar chart.
	categories := []string{
		"Required Skills",
		"Preferred Skills",
		"Experience",
		"Education",
	}

	// Calculate candidate coverage for each category.
	totalRequired := len(job.RequiredSkills)
	totalPreferred := len(job.PreferredSkills)

	missingRequired := len(criticalGaps)
	missingPreferred := len(importantGaps)

	matchedRequired := totalRequired - missingRequired
	matchedPreferred := totalPreferred - missingPreferred

	candidateRequiredScore := 0.0
	if totalRequired > 0 {
		candidateRequiredScore = float64(matchedRequired) / float64(totalRequired)
	} else {
		candidateRequiredScore = 1.0
	}

	candidatePreferredScore := 0.0
	if totalPreferred > 0 {
		candidatePreferredScore = float64(matchedPreferred) / float64(totalPreferred)
	} else {
		candidatePreferredScore = 1.0
	}

	// Experience score: simple ratio of candidate years to required years.
	expScore := 0.5 // default neutral
	if job.MinYearsExperience > 0 && profile.YearsOfExperience > 0 {
		expScore = math.Min(1.0, profile.YearsOfExperience/job.MinYearsExperience)
	} else if job.MinYearsExperience == 0 {
		expScore = 1.0
	}

	// Education score: simplified.
	eduScore := 1.0
	if job.RequiredDegreeLevel != "" && len(profile.Education) == 0 {
		eduScore = 0.3
	}

	candidateScores := []float64{
		roundTo2(candidateRequiredScore),
		roundTo2(candidatePreferredScore),
		roundTo2(expScore),
		roundTo2(eduScore),
	}

	requiredScores := []float64{1.0, 1.0, 1.0, 1.0}

	return RadarChartData{
		Labels:          categories,
		CandidateScores: candidateScores,
		RequiredScores:  requiredScores,
	}
}

// buildCategorySummary creates a summary for a single gap category.
func buildCategorySummary(category string, gaps []SkillGap) CategorySummary {
	totalHours := 0
	totalPriority := 0.0
	for _, g := range gaps {
		totalHours += g.EstimatedLearningHours
		totalPriority += g.PriorityScore
	}
	avgPriority := 0.0
	if len(gaps) > 0 {
		avgPriority = totalPriority / float64(len(gaps))
	}
	return CategorySummary{
		Category:           category,
		GapCount:           len(gaps),
		TotalLearningHours: totalHours,
		AveragePriority:    roundTo4(avgPriority),
	}
}

// buildLearningTimeline creates a suggested learning order.
// Critical gaps come first (sorted by priority), then important gaps.
func buildLearningTimeline(criticalGaps, importantGaps []SkillGap) []TimelineEntry {
	// Combine all gaps, critical first.
	ordered := append(append([]SkillGap{}, criticalGaps...), importantGaps...)

	var timeline []TimelineEntry
	cumulative := 0
	for i, g := range ordered {
		cumulative += g.EstimatedLearningHours
		rationale := buildTimelineRationale(g, i)
		timeline = append(timeline, TimelineEntry{
			Order:           i + 1,
			SkillName:       g.SkillName,
			Category:        string(g.Category),
			EstimatedHours:  g.EstimatedLearningHours,
			CumulativeHours: cumulative,
			Rationale:       rationale,
		})
	}
	return timeline
}

// buildTimelineRationale generates a rationale string for a timeline entry.
func buildTimelineRationale(gap SkillGap, position int) string {
	if gap.Category == GapCategoryCritical {
		if position == 0 {
			return "Highest priority: this is a must-have skill for the role with the highest impact on job acceptance."
		}
		return "Critical skill required for the role. Address this before applying."
	}
	if gap.SemanticSimilarityScore >= 0.5 {
		return "You have related skills that will accelerate learning this preferred skill."
	}
	return "Preferred skill that will strengthen your application once critical gaps are addressed."
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────────────────────

// buildCandidateIndex creates a map from normalized skill name to CandidateSkill.
func buildCandidateIndex(skills []scorer.CandidateSkill) map[string]scorer.CandidateSkill {
	index := make(map[string]scorer.CandidateSkill, len(skills))
	for _, s := range skills {
		norm := normalizeSkill(s.Name)
		if norm == "" {
			continue
		}
		// Keep the highest proficiency if the skill appears multiple times.
		if existing, ok := index[norm]; !ok || proficiencyRank(s.Proficiency) > proficiencyRank(existing.Proficiency) {
			index[norm] = s
		}
	}
	return index
}

// lookupInIndex checks whether a skill exists in the candidate index.
// It tries exact match first, then alias matching.
func lookupInIndex(norm string, index map[string]scorer.CandidateSkill) (scorer.CandidateSkill, bool) {
	if s, ok := index[norm]; ok {
		return s, true
	}
	// Alias matching.
	for candidateNorm, s := range index {
		if skillsAreAliases(norm, candidateNorm) {
			return s, true
		}
	}
	return scorer.CandidateSkill{}, false
}

// collectMatchedSkills returns skills from required and preferred lists that
// the candidate already has.
func collectMatchedSkills(required, preferred []string, index map[string]scorer.CandidateSkill) []string {
	var matched []string
	seen := map[string]bool{}
	for _, s := range append(required, preferred...) {
		norm := normalizeSkill(s)
		if _, found := lookupInIndex(norm, index); found && !seen[norm] {
			seen[norm] = true
			matched = append(matched, s)
		}
	}
	return matched
}

// getSkillMetadata returns metadata for a skill, falling back to defaults.
func getSkillMetadata(norm string) skillMetadata {
	if meta, ok := builtinSkillMetadata[norm]; ok {
		return meta
	}
	// Try alias resolution.
	aliases := map[string]string{
		"golang":     "go",
		"js":         "javascript",
		"ts":         "typescript",
		"node":       "node.js",
		"nodejs":     "node.js",
		"react.js":   "react",
		"reactjs":    "react",
		"vue.js":     "vue",
		"vuejs":      "vue",
		"angular.js": "angular",
		"angularjs":  "angular",
		"postgres":   "postgresql",
		"psql":       "postgresql",
		"k8s":        "kubernetes",
		"py":         "python",
		"rb":         "ruby",
		"cpp":        "c++",
		"csharp":     "c#",
		"dotnet":     ".net",
		"net":        ".net",
		"ml":         "machine learning",
		"dl":         "deep learning",
	}
	if canonical, ok := aliases[norm]; ok {
		if meta, ok := builtinSkillMetadata[canonical]; ok {
			return meta
		}
	}
	return defaultSkillMetadata
}

// skillsAreAliases returns true if two normalized skill names are equivalent.
func skillsAreAliases(a, b string) bool {
	aliases := map[string]string{
		"golang":     "go",
		"js":         "javascript",
		"ts":         "typescript",
		"node":       "node.js",
		"nodejs":     "node.js",
		"react.js":   "react",
		"reactjs":    "react",
		"vue.js":     "vue",
		"vuejs":      "vue",
		"angular.js": "angular",
		"angularjs":  "angular",
		"postgres":   "postgresql",
		"psql":       "postgresql",
		"k8s":        "kubernetes",
		"py":         "python",
		"rb":         "ruby",
		"cpp":        "c++",
		"csharp":     "c#",
		"dotnet":     ".net",
		"net":        ".net",
	}
	if canonical, ok := aliases[a]; ok {
		a = canonical
	}
	if canonical, ok := aliases[b]; ok {
		b = canonical
	}
	if a == b {
		return true
	}
	// Substring containment (e.g. "spring boot" contains "spring").
	if strings.Contains(a, b) || strings.Contains(b, a) {
		shorter := a
		if len(b) < len(a) {
			shorter = b
		}
		return len(shorter) >= 3
	}
	return false
}

// normalizeSkill lowercases and trims a skill name.
func normalizeSkill(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// proficiencyRank returns a numeric rank for a proficiency level.
func proficiencyRank(p string) int {
	switch strings.ToLower(p) {
	case "beginner":
		return 1
	case "intermediate":
		return 2
	case "advanced":
		return 3
	case "expert":
		return 4
	default:
		return 0
	}
}

// sortGapsByPriority sorts gaps by PriorityScore descending.
func sortGapsByPriority(gaps []SkillGap) {
	sort.Slice(gaps, func(i, j int) bool {
		return gaps[i].PriorityScore > gaps[j].PriorityScore
	})
}

// topN returns the first n elements of a slice.
func topN(gaps []SkillGap, n int) []SkillGap {
	if len(gaps) <= n {
		return gaps
	}
	return gaps[:n]
}

// sumLearningHours returns the total estimated learning hours for a list of gaps.
func sumLearningHours(gaps []SkillGap) int {
	total := 0
	for _, g := range gaps {
		total += g.EstimatedLearningHours
	}
	return total
}

// roundTo2 rounds a float64 to 2 decimal places.
func roundTo2(f float64) float64 {
	return math.Round(f*100) / 100
}

// roundTo4 rounds a float64 to 4 decimal places.
func roundTo4(f float64) float64 {
	return math.Round(f*10000) / 10000
}
