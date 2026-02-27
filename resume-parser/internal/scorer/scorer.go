// Package scorer implements the acceptance likelihood scoring algorithm.
package scorer

import (
	"math"
	"strings"
)

// degreeLevelRank maps degree level strings to a numeric rank for comparison.
// Higher rank = higher degree.
var degreeLevelRank = map[string]int{
	"high_school":  1,
	"certificate":  2,
	"diploma":      3,
	"associate":    4,
	"bachelor":     5,
	"master":       6,
	"professional": 7,
	"doctorate":    8,
	"other":        0,
}

// proficiencyWeight maps proficiency levels to a weight multiplier used when
// computing skill match scores.
var proficiencyWeight = map[string]float64{
	"beginner":     0.5,
	"intermediate": 0.75,
	"advanced":     0.9,
	"expert":       1.0,
	"":             0.7, // unknown proficiency – assume intermediate-ish
}

// experienceLevelYears maps experience level labels to approximate midpoint
// years of experience, used when the job specifies a level but no explicit
// year range.
var experienceLevelYears = map[string]float64{
	"internship": 0,
	"entry":      1,
	"mid":        3,
	"senior":     6,
	"lead":       8,
	"executive":  12,
}

// ─────────────────────────────────────────────────────────────────────────────
// Public API
// ─────────────────────────────────────────────────────────────────────────────

// Calculate computes the acceptance likelihood score for a candidate against
// a job posting. It returns a ScoreBreakdown with the overall score (0–100)
// and the individual component scores (0–1).
//
// The weighted formula is:
//
//	overall = (skill_match * 0.35 + experience_match * 0.25 +
//	           education_match * 0.15 + location_fit * 0.10 +
//	           industry_relevance * 0.15) * 100
func Calculate(profile CandidateProfile, job JobRequirements) ScoreBreakdown {
	skillScore, matched, missing, matchedPref := scoreSkillMatch(profile, job)
	expScore := scoreExperienceMatch(profile, job)
	eduScore := scoreEducationMatch(profile, job)
	locScore := scoreLocationFit(profile, job)
	indScore := scoreIndustryRelevance(profile, job)

	overall := (skillScore*WeightSkillMatch +
		expScore*WeightExperienceMatch +
		eduScore*WeightEducationMatch +
		locScore*WeightLocationFit +
		indScore*WeightIndustryRelevance) * 100.0

	// Clamp to [0, 100]
	overall = math.Max(0, math.Min(100, overall))

	return ScoreBreakdown{
		OverallScore:           roundTo2(overall),
		SkillMatchScore:        roundTo2(skillScore),
		ExperienceMatchScore:   roundTo2(expScore),
		EducationMatchScore:    roundTo2(eduScore),
		LocationFitScore:       roundTo2(locScore),
		IndustryRelevanceScore: roundTo2(indScore),
		MatchedRequiredSkills:  matched,
		MissingRequiredSkills:  missing,
		MatchedPreferredSkills: matchedPref,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Component scorers
// ─────────────────────────────────────────────────────────────────────────────

// scoreSkillMatch computes the skill match component score [0, 1].
//
// Algorithm:
//  1. Build a normalised lookup of candidate skills → proficiency weight.
//  2. For each required skill, check for an exact or fuzzy match.
//     - Matched required skills contribute to a weighted numerator.
//     - Missing required skills are tracked separately.
//  3. Preferred skills add a bonus (up to 20% of the required score).
//  4. Final score = required_score * 0.80 + preferred_bonus * 0.20
//     (capped at 1.0).
func scoreSkillMatch(profile CandidateProfile, job JobRequirements) (
	score float64,
	matchedRequired []string,
	missingRequired []string,
	matchedPreferred []string,
) {
	if len(job.RequiredSkills) == 0 {
		// No required skills specified – full score by default.
		return 1.0, nil, nil, nil
	}

	// Build candidate skill index: normalised name → proficiency weight.
	candidateIndex := buildSkillIndex(profile.Skills)

	// Score required skills.
	var requiredWeightedSum float64
	var requiredMaxSum float64

	for _, req := range job.RequiredSkills {
		reqNorm := normalizeSkillName(req)
		weight, found := lookupSkill(reqNorm, candidateIndex)
		requiredMaxSum += 1.0
		if found {
			requiredWeightedSum += weight
			matchedRequired = append(matchedRequired, req)
		} else {
			missingRequired = append(missingRequired, req)
		}
	}

	requiredScore := 0.0
	if requiredMaxSum > 0 {
		requiredScore = requiredWeightedSum / requiredMaxSum
	}

	// Score preferred skills (bonus).
	var preferredMatched int
	for _, pref := range job.PreferredSkills {
		prefNorm := normalizeSkillName(pref)
		_, found := lookupSkill(prefNorm, candidateIndex)
		if found {
			preferredMatched++
			matchedPreferred = append(matchedPreferred, pref)
		}
	}

	preferredBonus := 0.0
	if len(job.PreferredSkills) > 0 {
		preferredBonus = float64(preferredMatched) / float64(len(job.PreferredSkills))
	}

	// Combine: required skills are 80% of the score, preferred are 20%.
	score = requiredScore*0.80 + preferredBonus*0.20
	score = math.Min(1.0, score)
	return score, matchedRequired, missingRequired, matchedPreferred
}

// scoreExperienceMatch computes the experience match component score [0, 1].
//
// Algorithm:
//  1. Determine the target years from MinYearsExperience / ExperienceLevel.
//  2. Compute a ratio of candidate years to target years.
//  3. Apply a soft penalty for over-qualification (> 2× target).
//  4. Also consider title/role similarity from work history.
func scoreExperienceMatch(profile CandidateProfile, job JobRequirements) float64 {
	targetMin := job.MinYearsExperience

	// If no explicit minimum, infer from experience level.
	if targetMin == 0 && job.ExperienceLevel != "" {
		targetMin = experienceLevelYears[strings.ToLower(job.ExperienceLevel)]
	}

	candidateYears := profile.YearsOfExperience

	// If candidate years not set, estimate from work history.
	if candidateYears == 0 && len(profile.WorkHistory) > 0 {
		var totalMonths int
		for _, w := range profile.WorkHistory {
			totalMonths += w.DurationMonths
		}
		candidateYears = float64(totalMonths) / 12.0
	}

	yearsScore := computeYearsScore(candidateYears, targetMin, job.MaxYearsExperience)

	// Title similarity bonus: check if any past title matches the job title.
	titleBonus := computeTitleSimilarity(profile.WorkHistory, job.Title)

	// Combine: years are 70% of experience score, title similarity 30%.
	score := yearsScore*0.70 + titleBonus*0.30
	return math.Min(1.0, score)
}

// scoreEducationMatch computes the education match component score [0, 1].
//
// Algorithm:
//  1. If no degree is required, return 1.0.
//  2. Find the candidate's highest degree level.
//  3. Score based on whether the candidate meets or exceeds the requirement.
//  4. Apply a field-of-study bonus if the field matches preferred fields.
func scoreEducationMatch(profile CandidateProfile, job JobRequirements) float64 {
	if job.RequiredDegreeLevel == "" {
		return 1.0 // No education requirement.
	}

	requiredRank := degreeLevelRank[strings.ToLower(job.RequiredDegreeLevel)]
	if requiredRank == 0 {
		return 1.0 // Unknown requirement – don't penalise.
	}

	// Find candidate's highest degree rank.
	highestRank := 0
	var highestFieldOfStudy string
	for _, edu := range profile.Education {
		rank := degreeLevelRank[strings.ToLower(edu.DegreeLevel)]
		if rank > highestRank {
			highestRank = rank
			highestFieldOfStudy = edu.FieldOfStudy
		}
	}

	var degreeScore float64
	switch {
	case highestRank == 0:
		// No education info – partial credit.
		degreeScore = 0.3
	case highestRank >= requiredRank:
		degreeScore = 1.0
	case highestRank == requiredRank-1:
		// One level below – partial credit.
		degreeScore = 0.6
	default:
		// Two or more levels below.
		degreeScore = 0.2
	}

	// Field-of-study bonus (up to 0.2 added to degree score, capped at 1.0).
	fieldBonus := computeFieldBonus(highestFieldOfStudy, job.PreferredFields)
	return math.Min(1.0, degreeScore+fieldBonus*0.2)
}

// scoreLocationFit computes the location fit component score [0, 1].
//
// Algorithm:
//  1. If the job is fully remote, check candidate's remote preference.
//  2. If the job is on-site or hybrid, check city/country match or relocation.
func scoreLocationFit(profile CandidateProfile, job JobRequirements) float64 {
	jobLocType := strings.ToLower(job.LocationType)
	candidatePref := strings.ToLower(profile.RemotePreference)

	switch jobLocType {
	case "remote":
		// Remote job: candidate who prefers remote or "any" is a perfect fit.
		if candidatePref == "remote" || candidatePref == "any" || candidatePref == "" {
			return 1.0
		}
		// Candidate prefers on-site but job is remote – partial fit.
		return 0.7

	case "hybrid":
		if candidatePref == "hybrid" || candidatePref == "any" || candidatePref == "" {
			return 1.0
		}
		if candidatePref == "remote" {
			return 0.6
		}
		// on_site preference for hybrid role – reasonable fit.
		return 0.8

	default: // "on_site" or unknown
		// Check geographic match.
		return computeGeoScore(profile, job)
	}
}

// scoreIndustryRelevance computes the industry relevance component score [0, 1].
//
// Algorithm:
//  1. If no industry is specified in the job, return 1.0.
//  2. Check if any of the candidate's past industries match the target.
//  3. Check related industries for partial credit.
func scoreIndustryRelevance(profile CandidateProfile, job JobRequirements) float64 {
	if job.Industry == "" {
		return 1.0 // No industry requirement.
	}

	targetIndustry := normalizeIndustry(job.Industry)

	// Build set of related industries (including the target itself).
	relatedSet := map[string]bool{targetIndustry: true}
	for _, rel := range job.RelatedIndustries {
		relatedSet[normalizeIndustry(rel)] = true
	}

	// Check candidate's work history industries.
	var directMatch bool
	var relatedMatch bool

	for _, w := range profile.WorkHistory {
		if w.Industry == "" {
			continue
		}
		norm := normalizeIndustry(w.Industry)
		if norm == targetIndustry {
			directMatch = true
			break
		}
		if relatedSet[norm] {
			relatedMatch = true
		}
	}

	switch {
	case directMatch:
		return 1.0
	case relatedMatch:
		return 0.7
	case len(profile.WorkHistory) == 0:
		// No work history – neutral score.
		return 0.5
	default:
		// No industry match at all.
		return 0.2
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────────────────────

// buildSkillIndex creates a map from normalised skill name to proficiency weight.
func buildSkillIndex(skills []CandidateSkill) map[string]float64 {
	index := make(map[string]float64, len(skills))
	for _, s := range skills {
		norm := normalizeSkillName(s.Name)
		if norm == "" {
			continue
		}
		w := proficiencyWeight[strings.ToLower(s.Proficiency)]
		// Keep the highest weight if the skill appears multiple times.
		if existing, ok := index[norm]; !ok || w > existing {
			index[norm] = w
		}
	}
	return index
}

// lookupSkill checks whether a required skill exists in the candidate index.
// It first tries an exact normalised match, then a substring/alias match.
// Returns the proficiency weight and whether a match was found.
func lookupSkill(reqNorm string, index map[string]float64) (float64, bool) {
	// Exact match.
	if w, ok := index[reqNorm]; ok {
		return w, true
	}

	// Alias / substring match: e.g. "golang" matches "go", "node" matches "node.js".
	for candidateNorm, w := range index {
		if skillsAreAliases(reqNorm, candidateNorm) {
			return w, true
		}
	}
	return 0, false
}

// skillsAreAliases returns true if two normalised skill names are considered
// equivalent (e.g. "golang" and "go", "javascript" and "js").
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

	// Resolve aliases.
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
		// Only allow if the shorter string is at least 3 chars to avoid false positives.
		shorter := a
		if len(b) < len(a) {
			shorter = b
		}
		return len(shorter) >= 3
	}

	return false
}

// normalizeSkillName lowercases and trims a skill name.
func normalizeSkillName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// computeYearsScore returns a score [0, 1] based on candidate years vs target.
func computeYearsScore(candidateYears, targetMin, targetMax float64) float64 {
	if targetMin == 0 {
		return 1.0 // No minimum requirement.
	}

	if candidateYears >= targetMin {
		// Check for over-qualification if a max is specified.
		if targetMax > 0 && candidateYears > targetMax*2 {
			// Significantly over-qualified – slight penalty.
			return 0.8
		}
		return 1.0
	}

	// Under-qualified: linear decay from targetMin down to 0.
	// At 0 years, score = 0.1 (not zero, to avoid harsh penalisation of
	// candidates who are close to the requirement).
	ratio := candidateYears / targetMin
	return math.Max(0.1, ratio)
}

// computeTitleSimilarity returns a bonus [0, 1] based on how closely the
// candidate's past job titles match the target job title.
func computeTitleSimilarity(history []WorkHistoryEntry, targetTitle string) float64 {
	if targetTitle == "" || len(history) == 0 {
		return 0.5 // Neutral when no data.
	}

	targetNorm := normalizeJobTitle(targetTitle)
	targetWords := strings.Fields(targetNorm)

	bestScore := 0.0
	for _, w := range history {
		candidateNorm := normalizeJobTitle(w.Title)
		candidateWords := strings.Fields(candidateNorm)
		sim := wordOverlapScore(targetWords, candidateWords)
		if sim > bestScore {
			bestScore = sim
		}
	}
	return bestScore
}

// normalizeJobTitle lowercases and removes common filler words from a title.
func normalizeJobTitle(title string) string {
	title = strings.ToLower(strings.TrimSpace(title))
	// Remove common filler words that don't add signal.
	fillers := []string{"and", "the", "of", "a", "an", "in", "at", "for", "to"}
	words := strings.Fields(title)
	var filtered []string
	for _, w := range words {
		isFiller := false
		for _, f := range fillers {
			if w == f {
				isFiller = true
				break
			}
		}
		if !isFiller {
			filtered = append(filtered, w)
		}
	}
	return strings.Join(filtered, " ")
}

// wordOverlapScore computes the Jaccard similarity between two word sets.
func wordOverlapScore(a, b []string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	setA := make(map[string]bool, len(a))
	for _, w := range a {
		setA[w] = true
	}
	intersection := 0
	for _, w := range b {
		if setA[w] {
			intersection++
		}
	}
	union := len(setA)
	for _, w := range b {
		if !setA[w] {
			union++
		}
	}
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

// computeFieldBonus returns a bonus [0, 1] if the candidate's field of study
// matches any of the preferred fields.
func computeFieldBonus(candidateField string, preferredFields []string) float64 {
	if candidateField == "" || len(preferredFields) == 0 {
		return 0
	}
	candidateNorm := strings.ToLower(strings.TrimSpace(candidateField))
	for _, pf := range preferredFields {
		pfNorm := strings.ToLower(strings.TrimSpace(pf))
		if candidateNorm == pfNorm ||
			strings.Contains(candidateNorm, pfNorm) ||
			strings.Contains(pfNorm, candidateNorm) {
			return 1.0
		}
	}
	return 0
}

// computeGeoScore returns a location score for on-site / hybrid jobs.
func computeGeoScore(profile CandidateProfile, job JobRequirements) float64 {
	// Same city → perfect match.
	if profile.LocationCity != "" && job.LocationCity != "" {
		if strings.EqualFold(profile.LocationCity, job.LocationCity) {
			return 1.0
		}
	}

	// Same country → good match (commutable or willing to relocate within country).
	if profile.LocationCountry != "" && job.LocationCountry != "" {
		if strings.EqualFold(profile.LocationCountry, job.LocationCountry) {
			return 0.8
		}
	}

	// Willing to relocate internationally → partial credit.
	if profile.WillingToRelocate {
		return 0.6
	}

	// Different country, not willing to relocate.
	return 0.2
}

// normalizeIndustry lowercases and trims an industry string.
func normalizeIndustry(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// roundTo2 rounds a float64 to 2 decimal places.
func roundTo2(v float64) float64 {
	return math.Round(v*100) / 100
}
