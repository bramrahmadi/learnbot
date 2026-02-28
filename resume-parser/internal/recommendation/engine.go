// Package recommendation – engine.go implements the core recommendation engine.
// It takes a candidate profile and job requirements, runs gap analysis,
// and produces a personalized learning plan.
package recommendation

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/learnbot/resume-parser/internal/gapanalysis"
	"github.com/learnbot/resume-parser/internal/scorer"
)

// ─────────────────────────────────────────────────────────────────────────────
// Engine
// ─────────────────────────────────────────────────────────────────────────────

// Engine is the training recommendation engine.
type Engine struct {
	gapAnalyzer *gapanalysis.Analyzer
	catalog     []ResourceEntry
}

// New creates a new recommendation Engine with the built-in resource catalog.
func New() *Engine {
	return &Engine{
		gapAnalyzer: gapanalysis.New(),
		catalog:     builtinCatalog,
	}
}

// NewWithCatalog creates a new Engine with a custom resource catalog.
// Useful for testing.
func NewWithCatalog(catalog []ResourceEntry) *Engine {
	return &Engine{
		gapAnalyzer: gapanalysis.New(),
		catalog:     catalog,
	}
}

// Generate produces a personalized learning plan for the given profile, job,
// and user preferences.
func (e *Engine) Generate(
	profile scorer.CandidateProfile,
	job scorer.JobRequirements,
	prefs UserPreferences,
) LearningPlan {
	// Apply defaults to preferences.
	prefs = applyPreferenceDefaults(prefs)

	// Run gap analysis.
	gapResult := e.gapAnalyzer.Analyze(profile, job)

	// Build skill recommendations for each gap category.
	criticalRecs := e.buildSkillRecommendations(gapResult.CriticalGaps, prefs)
	importantRecs := e.buildSkillRecommendations(gapResult.ImportantGaps, prefs)
	niceToHaveRecs := e.buildSkillRecommendations(gapResult.NiceToHaveGaps, prefs)

	// Build learning phases.
	phases := buildPhases(criticalRecs, importantRecs, niceToHaveRecs, prefs)

	// Calculate totals.
	totalHours := sumPhaseHours(phases)

	// Build timeline.
	timeline := buildTimeline(phases, prefs, job.Title)

	// Build summary.
	summary := buildSummary(gapResult, phases, job.Title)

	return LearningPlan{
		JobTitle:            job.Title,
		ReadinessScore:      gapResult.ReadinessScore,
		TotalGaps:           gapResult.TotalGaps,
		TotalEstimatedHours: totalHours,
		Phases:              phases,
		Timeline:            timeline,
		MatchedSkills:       gapResult.MatchedSkills,
		Summary:             summary,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Skill recommendation building
// ─────────────────────────────────────────────────────────────────────────────

// buildSkillRecommendations creates SkillRecommendation entries for a list of gaps.
func (e *Engine) buildSkillRecommendations(gaps []gapanalysis.SkillGap, prefs UserPreferences) []SkillRecommendation {
	var recs []SkillRecommendation
	for _, gap := range gaps {
		rec := e.buildSkillRecommendation(gap, prefs)
		recs = append(recs, rec)
	}
	return recs
}

// buildSkillRecommendation creates a SkillRecommendation for a single gap.
func (e *Engine) buildSkillRecommendation(gap gapanalysis.SkillGap, prefs UserPreferences) SkillRecommendation {
	// Find matching resources from the catalog.
	candidates := e.findMatchingResources(gap.SkillName, prefs)

	// Score and rank candidates.
	scored := e.scoreResources(candidates, gap, prefs)

	// Sort by relevance score descending.
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].RelevanceScore > scored[j].RelevanceScore
	})

	var primary *RecommendedResource
	var alternatives []RecommendedResource

	if len(scored) > 0 {
		p := scored[0]
		p.IsAlternative = false
		p.RecommendationReason = buildRecommendationReason(p.Resource, gap, prefs)
		primary = &p
	}

	// Add up to 2 alternatives (different type or provider from primary).
	for i := 1; i < len(scored) && len(alternatives) < 2; i++ {
		alt := scored[i]
		if primary != nil &&
			(alt.Resource.ResourceType != primary.Resource.ResourceType ||
				alt.Resource.Provider != primary.Resource.Provider) {
			alt.IsAlternative = true
			alt.RecommendationReason = buildRecommendationReason(alt.Resource, gap, prefs)
			alternatives = append(alternatives, alt)
		}
	}

	return SkillRecommendation{
		SkillName:                gap.SkillName,
		GapCategory:              string(gap.Category),
		PriorityScore:            gap.PriorityScore,
		PrimaryResource:          primary,
		AlternativeResources:     alternatives,
		EstimatedHoursToJobReady: gap.EstimatedLearningHours,
		CurrentLevel:             gap.CurrentLevel,
		TargetLevel:              gap.TargetLevel,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource matching and scoring
// ─────────────────────────────────────────────────────────────────────────────

// findMatchingResources returns catalog resources that cover the given skill.
func (e *Engine) findMatchingResources(skillName string, prefs UserPreferences) []ResourceEntry {
	norm := normalizeSkillName(skillName)
	canonical := resolveAlias(norm)

	var matches []ResourceEntry
	for _, res := range e.catalog {
		if !resourceMatchesSkill(res, canonical, norm) {
			continue
		}
		if !passesPreferenceFilter(res, prefs) {
			continue
		}
		matches = append(matches, res)
	}
	return matches
}

// resourceMatchesSkill returns true if a resource covers the given skill.
func resourceMatchesSkill(res ResourceEntry, canonical, norm string) bool {
	// Check primary skill.
	if normalizeSkillName(res.PrimarySkill) == canonical ||
		normalizeSkillName(res.PrimarySkill) == norm {
		return true
	}
	// Check all skills.
	for _, s := range res.Skills {
		sNorm := normalizeSkillName(s)
		if sNorm == canonical || sNorm == norm {
			return true
		}
		// Substring match for compound skills (e.g. "spring boot" matches "spring").
		if len(canonical) >= 3 && strings.Contains(sNorm, canonical) {
			return true
		}
		if len(norm) >= 3 && strings.Contains(sNorm, norm) {
			return true
		}
	}
	return false
}

// passesPreferenceFilter returns true if a resource passes the user's filters.
func passesPreferenceFilter(res ResourceEntry, prefs UserPreferences) bool {
	// Free preference filter.
	if prefs.PreferFree && res.CostType != "free" && res.CostType != "free_audit" {
		return false
	}

	// Budget filter.
	if prefs.MaxBudgetUSD > 0 && res.CostUSD > prefs.MaxBudgetUSD {
		return false
	}

	// Excluded providers.
	for _, excluded := range prefs.ExcludedProviders {
		if strings.EqualFold(res.Provider, excluded) {
			return false
		}
	}

	// Preferred resource types filter (if specified).
	if len(prefs.PreferredResourceTypes) > 0 {
		found := false
		for _, t := range prefs.PreferredResourceTypes {
			if strings.EqualFold(res.ResourceType, t) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// scoreResources scores and ranks resources for a skill gap.
func (e *Engine) scoreResources(resources []ResourceEntry, gap gapanalysis.SkillGap, prefs UserPreferences) []RecommendedResource {
	var scored []RecommendedResource
	for _, res := range resources {
		score := computeRelevanceScore(res, gap, prefs)
		hours := estimateCompletionHours(res, gap)
		scored = append(scored, RecommendedResource{
			Resource:                 res,
			RelevanceScore:           roundTo4(score),
			EstimatedCompletionHours: hours,
		})
	}
	return scored
}

// computeRelevanceScore computes a composite relevance score [0, 1] for a resource.
//
// Factors and weights:
//   - Skill match quality (primary vs secondary): 0.30
//   - Difficulty fit (matches gap target level): 0.20
//   - Quality (rating, verification): 0.20
//   - User preference alignment: 0.20
//   - Popularity (rating count): 0.10
func computeRelevanceScore(res ResourceEntry, gap gapanalysis.SkillGap, prefs UserPreferences) float64 {
	// 1. Skill match quality.
	skillScore := 0.5 // secondary skill match
	if normalizeSkillName(res.PrimarySkill) == resolveAlias(normalizeSkillName(gap.SkillName)) {
		skillScore = 1.0 // primary skill match
	}

	// 2. Difficulty fit.
	diffScore := computeDifficultyFit(res.Difficulty, gap.TargetLevel, gap.CurrentLevel)

	// 3. Quality score.
	qualityScore := 0.5 // default
	if res.Rating > 0 {
		qualityScore = res.Rating / 5.0
	}
	if res.IsVerified {
		qualityScore = math.Min(1.0, qualityScore+0.1)
	}

	// 4. Preference alignment.
	prefScore := computePreferenceAlignment(res, prefs)

	// 5. Popularity (log-normalized rating count).
	popularityScore := 0.0
	if res.RatingCount > 0 {
		popularityScore = math.Min(1.0, math.Log10(float64(res.RatingCount))/7.0) // 10M = 1.0
	}

	return skillScore*0.30 +
		diffScore*0.20 +
		qualityScore*0.20 +
		prefScore*0.20 +
		popularityScore*0.10
}

// computeDifficultyFit returns a score [0, 1] for how well the resource
// difficulty matches the user's current level and target level.
func computeDifficultyFit(resDifficulty, targetLevel, currentLevel string) float64 {
	if resDifficulty == "all_levels" {
		return 0.9
	}

	// Map levels to numeric ranks.
	levelRank := map[string]int{
		"beginner":     1,
		"intermediate": 2,
		"advanced":     3,
		"expert":       4,
		"all_levels":   2, // treat as intermediate
	}

	targetRank := levelRank[strings.ToLower(targetLevel)]
	if targetRank == 0 {
		targetRank = 2 // default to intermediate
	}

	currentRank := levelRank[strings.ToLower(currentLevel)]
	// If no current level, assume one below target.
	if currentRank == 0 {
		currentRank = max(1, targetRank-1)
	}

	resRank := levelRank[strings.ToLower(resDifficulty)]
	if resRank == 0 {
		resRank = 2
	}

	// Ideal: resource difficulty is between current and target level.
	if resRank >= currentRank && resRank <= targetRank {
		return 1.0
	}
	// One level off.
	diff := abs(resRank - targetRank)
	if diff == 1 {
		return 0.7
	}
	return 0.4
}

// computePreferenceAlignment returns a score [0, 1] for how well a resource
// aligns with user preferences.
func computePreferenceAlignment(res ResourceEntry, prefs UserPreferences) float64 {
	score := 0.5 // neutral baseline

	// Free preference.
	if prefs.PreferFree && (res.CostType == "free" || res.CostType == "free_audit") {
		score += 0.3
	}

	// Hands-on preference.
	if prefs.PreferHandsOn && res.HasHandsOn {
		score += 0.1
	}

	// Certificate preference.
	if prefs.PreferCertificates && res.HasCertificate {
		score += 0.1
	}

	return math.Min(1.0, score)
}

// estimateCompletionHours estimates the hours to complete a resource given
// the user's current level.
func estimateCompletionHours(res ResourceEntry, gap gapanalysis.SkillGap) float64 {
	if res.DurationHours <= 0 {
		// Use gap estimate if resource has no duration.
		return float64(gap.EstimatedLearningHours)
	}

	hours := res.DurationHours

	// Reduce hours if user has partial knowledge.
	switch strings.ToLower(gap.CurrentLevel) {
	case "beginner":
		hours *= 0.70
	case "intermediate":
		hours *= 0.40
	case "advanced":
		hours *= 0.15
	}

	return math.Max(1.0, math.Round(hours*10)/10)
}

// ─────────────────────────────────────────────────────────────────────────────
// Phase building
// ─────────────────────────────────────────────────────────────────────────────

// buildPhases creates the ordered learning phases from skill recommendations.
func buildPhases(
	criticalRecs, importantRecs, niceToHaveRecs []SkillRecommendation,
	prefs UserPreferences,
) []LearningPhase {
	var phases []LearningPhase

	if len(criticalRecs) > 0 {
		phase := buildPhase(1, "Critical Skills", criticalRecs, prefs,
			"Master the must-have skills required for this role. These are non-negotiable for job acceptance.")
		phases = append(phases, phase)
	}

	if len(importantRecs) > 0 {
		phase := buildPhase(2, "Preferred Skills", importantRecs, prefs,
			"Strengthen your profile with preferred skills that significantly improve your candidacy.")
		phases = append(phases, phase)
	}

	if len(niceToHaveRecs) > 0 {
		phase := buildPhase(3, "Nice-to-Have Skills", niceToHaveRecs, prefs,
			"Optional skills that differentiate you from other candidates and expand your career options.")
		phases = append(phases, phase)
	}

	return phases
}

// buildPhase creates a single learning phase.
func buildPhase(
	phaseNum int,
	name string,
	recs []SkillRecommendation,
	prefs UserPreferences,
	description string,
) LearningPhase {
	totalHours := 0.0
	for _, rec := range recs {
		totalHours += float64(rec.EstimatedHoursToJobReady)
	}

	weeklyHours := prefs.WeeklyHoursAvailable
	if weeklyHours <= 0 {
		weeklyHours = 10
	}

	estimatedWeeks := totalHours / weeklyHours
	milestone := buildPhaseMilestone(phaseNum, name, recs)

	return LearningPhase{
		PhaseNumber:      phaseNum,
		PhaseName:        name,
		PhaseDescription: description,
		Skills:           recs,
		TotalHours:       roundTo2(totalHours),
		EstimatedWeeks:   roundTo1(estimatedWeeks),
		Milestone:        milestone,
	}
}

// buildPhaseMilestone generates a milestone description for a phase.
func buildPhaseMilestone(phaseNum int, phaseName string, recs []SkillRecommendation) string {
	if len(recs) == 0 {
		return ""
	}

	skillNames := make([]string, 0, len(recs))
	for _, r := range recs {
		skillNames = append(skillNames, r.SkillName)
	}

	switch phaseNum {
	case 1:
		if len(skillNames) == 1 {
			return fmt.Sprintf("✅ Job-ready in %s – ready to apply for the role", skillNames[0])
		}
		return fmt.Sprintf("✅ Job-ready in %s – cleared all critical requirements", strings.Join(skillNames[:min(3, len(skillNames))], ", "))
	case 2:
		return fmt.Sprintf("⭐ Strong candidate – proficient in %d preferred skills", len(recs))
	default:
		return fmt.Sprintf("🚀 Standout candidate – mastered %d additional skills", len(recs))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Summary building
// ─────────────────────────────────────────────────────────────────────────────

// buildSummary creates a high-level summary of the learning plan.
func buildSummary(
	gapResult gapanalysis.GapAnalysisResult,
	phases []LearningPhase,
	jobTitle string,
) LearningPlanSummary {
	freeCount := 0
	paidCount := 0
	totalCost := 0.0
	topSkills := make([]string, 0, 3)
	quickWins := make([]string, 0)

	for _, phase := range phases {
		for _, rec := range phase.Skills {
			if rec.PrimaryResource != nil {
				if rec.PrimaryResource.Resource.CostType == "free" ||
					rec.PrimaryResource.Resource.CostType == "free_audit" {
					freeCount++
				} else {
					paidCount++
					totalCost += rec.PrimaryResource.Resource.CostUSD
				}
			}
			if len(topSkills) < 3 {
				topSkills = append(topSkills, rec.SkillName)
			}
			if rec.EstimatedHoursToJobReady <= 20 {
				quickWins = append(quickWins, rec.SkillName)
			}
		}
	}

	// Build headline.
	headline := buildHeadline(gapResult, phases, jobTitle)

	return LearningPlanSummary{
		Headline:              headline,
		CriticalGapCount:      gapResult.CriticalGapCount,
		ImportantGapCount:     gapResult.ImportantGapCount,
		FreeResourceCount:     freeCount,
		PaidResourceCount:     paidCount,
		EstimatedTotalCostUSD: roundTo2(totalCost),
		TopSkillsToLearn:      topSkills,
		QuickWins:             quickWins,
	}
}

// buildHeadline generates a one-line summary headline.
func buildHeadline(gapResult gapanalysis.GapAnalysisResult, phases []LearningPhase, jobTitle string) string {
	if gapResult.TotalGaps == 0 {
		return fmt.Sprintf("🎉 You're ready to apply for %s!", jobTitle)
	}

	totalWeeks := 0.0
	for _, p := range phases {
		totalWeeks += p.EstimatedWeeks
	}

	if gapResult.CriticalGapCount == 0 {
		return fmt.Sprintf("✨ Strong candidate for %s – %d preferred skills to strengthen in %.0f weeks",
			jobTitle, gapResult.ImportantGapCount, totalWeeks)
	}

	return fmt.Sprintf("📚 %d critical gap%s to close for %s – estimated %.0f weeks",
		gapResult.CriticalGapCount,
		pluralize(gapResult.CriticalGapCount),
		jobTitle,
		totalWeeks)
}

// ─────────────────────────────────────────────────────────────────────────────
// Recommendation reason
// ─────────────────────────────────────────────────────────────────────────────

// buildRecommendationReason generates a human-readable reason for recommending
// a resource.
func buildRecommendationReason(res ResourceEntry, gap gapanalysis.SkillGap, prefs UserPreferences) string {
	var reasons []string

	// Primary skill match.
	if normalizeSkillName(res.PrimarySkill) == resolveAlias(normalizeSkillName(gap.SkillName)) {
		reasons = append(reasons, "directly covers "+gap.SkillName)
	}

	// Quality.
	if res.Rating >= 4.7 {
		reasons = append(reasons, fmt.Sprintf("highly rated (%.1f/5)", res.Rating))
	}

	// Free.
	if res.CostType == "free" || res.CostType == "free_audit" {
		reasons = append(reasons, "free to access")
	}

	// Certificate.
	if prefs.PreferCertificates && res.HasCertificate {
		reasons = append(reasons, "includes certificate")
	}

	// Hands-on.
	if prefs.PreferHandsOn && res.HasHandsOn {
		reasons = append(reasons, "hands-on learning")
	}

	// Verified.
	if res.IsVerified {
		reasons = append(reasons, "curated resource")
	}

	if len(reasons) == 0 {
		return "Recommended based on skill coverage and quality."
	}

	return "Recommended because it " + strings.Join(reasons, ", ") + "."
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────────────────────

// applyPreferenceDefaults fills in default values for missing preferences.
func applyPreferenceDefaults(prefs UserPreferences) UserPreferences {
	if prefs.WeeklyHoursAvailable <= 0 {
		prefs.WeeklyHoursAvailable = 10
	}
	return prefs
}

// normalizeSkillName lowercases and trims a skill name.
func normalizeSkillName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// resolveAlias resolves a skill alias to its canonical name.
func resolveAlias(norm string) string {
	if canonical, ok := skillAliases[norm]; ok {
		return canonical
	}
	return norm
}

// sumPhaseHours returns the total hours across all phases.
func sumPhaseHours(phases []LearningPhase) float64 {
	total := 0.0
	for _, p := range phases {
		total += p.TotalHours
	}
	return total
}

// roundTo2 rounds a float64 to 2 decimal places.
func roundTo2(f float64) float64 {
	return math.Round(f*100) / 100
}

// roundTo1 rounds a float64 to 1 decimal place.
func roundTo1(f float64) float64 {
	return math.Round(f*10) / 10
}

// roundTo4 rounds a float64 to 4 decimal places.
func roundTo4(f float64) float64 {
	return math.Round(f*10000) / 10000
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// pluralize returns "s" if n != 1, else "".
func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
