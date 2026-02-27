// Package taxonomy – taxonomy.go implements the in-memory taxonomy database,
// skill normalization (exact + alias + fuzzy matching), and NLP-based skill
// extraction from free-form text.
package taxonomy

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

// ─────────────────────────────────────────────────────────────────────────────
// Taxonomy (in-memory database)
// ─────────────────────────────────────────────────────────────────────────────

// Taxonomy is the in-memory skill ontology database.
// It provides O(1) lookup by ID, O(1) lookup by alias, and linear search
// for fuzzy matching.
type Taxonomy struct {
	// byID maps canonical skill ID → SkillNode.
	byID map[string]*SkillNode

	// byAlias maps normalised alias → canonical skill ID.
	byAlias map[string]string

	// all is the ordered list of all skill nodes (for iteration).
	all []*SkillNode
}

// New creates a Taxonomy populated with the built-in skill ontology.
func New() *Taxonomy {
	t := &Taxonomy{
		byID:    make(map[string]*SkillNode, len(builtinSkills)),
		byAlias: make(map[string]string),
	}
	for i := range builtinSkills {
		node := &builtinSkills[i]
		t.byID[node.ID] = node
		t.all = append(t.all, node)

		// Index the canonical name as an alias.
		t.byAlias[normalise(node.CanonicalName)] = node.ID
		// Index the ID itself.
		t.byAlias[normalise(node.ID)] = node.ID
		// Index all declared aliases.
		for _, alias := range node.Aliases {
			t.byAlias[normalise(alias)] = node.ID
		}
	}
	return t
}

// Lookup returns the SkillNode for the given canonical ID, or nil if not found.
func (t *Taxonomy) Lookup(id string) *SkillNode {
	return t.byID[id]
}

// All returns all skill nodes in the taxonomy.
func (t *Taxonomy) All() []*SkillNode {
	return t.all
}

// Search returns skill nodes whose canonical name, ID, or aliases contain the
// query string. Results are filtered by domain and category when non-empty.
// Limit controls the maximum number of results (0 = no limit).
func (t *Taxonomy) Search(query string, domain Domain, category Category, limit int) []SkillNode {
	q := normalise(query)
	var results []SkillNode
	for _, node := range t.all {
		if domain != "" && node.Domain != domain {
			continue
		}
		if category != "" && node.Category != category {
			continue
		}
		if q == "" ||
			strings.Contains(normalise(node.CanonicalName), q) ||
			strings.Contains(normalise(node.ID), q) ||
			aliasContains(node.Aliases, q) {
			results = append(results, *node)
		}
		if limit > 0 && len(results) >= limit {
			break
		}
	}
	return results
}

// aliasContains returns true if any alias contains the query.
func aliasContains(aliases []string, q string) bool {
	for _, a := range aliases {
		if strings.Contains(normalise(a), q) {
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// Normalisation & Fuzzy Matching
// ─────────────────────────────────────────────────────────────────────────────

// Normalize maps a raw skill string to its canonical taxonomy entry.
// It tries, in order:
//  1. Exact match on normalised canonical name / ID / alias.
//  2. Fuzzy match using Jaro-Winkler similarity (threshold 0.85).
//
// Returns a NormalizeResult with MatchType "exact", "alias", "fuzzy", or "none".
func (t *Taxonomy) Normalize(raw string) NormalizeResult {
	result := NormalizeResult{Input: raw}
	if strings.TrimSpace(raw) == "" {
		result.MatchType = "none"
		return result
	}

	norm := normalise(raw)

	// 1. Exact / alias match.
	if id, ok := t.byAlias[norm]; ok {
		node := t.byID[id]
		result.CanonicalID = node.ID
		result.CanonicalName = node.CanonicalName
		result.Domain = node.Domain
		result.Category = node.Category
		result.FuzzyScore = 1.0
		// Distinguish exact ID/name match from alias match.
		if norm == normalise(node.ID) || norm == normalise(node.CanonicalName) {
			result.MatchType = "exact"
		} else {
			result.MatchType = "alias"
		}
		return result
	}

	// 2. Fuzzy match.
	bestScore := 0.0
	var bestNode *SkillNode

	for _, node := range t.all {
		// Compare against canonical name.
		score := jaroWinkler(norm, normalise(node.CanonicalName))
		if score > bestScore {
			bestScore = score
			bestNode = node
		}
		// Compare against each alias.
		for _, alias := range node.Aliases {
			score = jaroWinkler(norm, normalise(alias))
			if score > bestScore {
				bestScore = score
				bestNode = node
			}
		}
	}

	const fuzzyThreshold = 0.85
	if bestScore >= fuzzyThreshold && bestNode != nil {
		result.CanonicalID = bestNode.ID
		result.CanonicalName = bestNode.CanonicalName
		result.Domain = bestNode.Domain
		result.Category = bestNode.Category
		result.FuzzyScore = roundTo4(bestScore)
		result.MatchType = "fuzzy"
		return result
	}

	result.MatchType = "none"
	return result
}

// NormalizeMany normalizes a slice of raw skill strings.
func (t *Taxonomy) NormalizeMany(raws []string) []NormalizeResult {
	results := make([]NormalizeResult, len(raws))
	for i, raw := range raws {
		results[i] = t.Normalize(raw)
	}
	return results
}

// ─────────────────────────────────────────────────────────────────────────────
// NLP-based Skill Extraction
// ─────────────────────────────────────────────────────────────────────────────

// Extractor extracts skills from free-form text using a combination of:
//   - Taxonomy alias lookup (exact + alias matching)
//   - Pattern-based NER (named entity recognition heuristics)
//   - Fuzzy matching for near-matches
type Extractor struct {
	taxonomy *Taxonomy

	// multiWordPatterns is a list of compiled regexes for multi-word skill
	// patterns that are hard to catch with simple tokenisation.
	multiWordPatterns []*regexp.Regexp

	// sentenceSplitter splits text into sentences/clauses.
	sentenceSplitter *regexp.Regexp

	// tokenSplitter splits a sentence into tokens.
	tokenSplitter *regexp.Regexp
}

// NewExtractor creates an Extractor backed by the given taxonomy.
func NewExtractor(t *Taxonomy) *Extractor {
	e := &Extractor{
		taxonomy:         t,
		sentenceSplitter: regexp.MustCompile(`[.!?\n;]+`),
		tokenSplitter:    regexp.MustCompile(`[\s,/|•\t]+`),
	}
	e.buildMultiWordPatterns()
	return e
}

// buildMultiWordPatterns compiles regex patterns for multi-word skills that
// appear frequently in job descriptions.
func (e *Extractor) buildMultiWordPatterns() {
	// Collect all multi-word aliases and canonical names from the taxonomy.
	seen := map[string]bool{}
	var patterns []string

	for _, node := range e.taxonomy.all {
		candidates := append([]string{node.CanonicalName}, node.Aliases...)
		for _, c := range candidates {
			if strings.Contains(c, " ") && !seen[c] {
				seen[c] = true
				// Escape for regex and allow flexible whitespace.
				escaped := regexp.QuoteMeta(strings.ToLower(c))
				escaped = strings.ReplaceAll(escaped, `\ `, `[\s\-]+`)
				patterns = append(patterns, escaped)
			}
		}
	}

	// Sort by length descending so longer patterns match first.
	sortByLengthDesc(patterns)

	for _, p := range patterns {
		re, err := regexp.Compile(`(?i)\b` + p + `\b`)
		if err == nil {
			e.multiWordPatterns = append(e.multiWordPatterns, re)
		}
	}
}

// Extract extracts skills from the given text.
// It returns an ExtractionResult with deduplicated skills grouped by type.
func (e *Extractor) Extract(text string, includeUnknown bool) ExtractionResult {
	if strings.TrimSpace(text) == "" {
		return ExtractionResult{}
	}

	// Track seen canonical IDs to avoid duplicates.
	seenIDs := map[string]bool{}
	// Track seen raw texts for unknown skills.
	seenRaw := map[string]bool{}

	var allSkills []ExtractedSkill

	// Step 1: Multi-word pattern matching (highest priority).
	remaining := text
	for _, re := range e.multiWordPatterns {
		matches := re.FindAllString(remaining, -1)
		for _, match := range matches {
			norm := e.taxonomy.Normalize(match)
			if norm.MatchType != "none" && !seenIDs[norm.CanonicalID] {
				seenIDs[norm.CanonicalID] = true
				allSkills = append(allSkills, ExtractedSkill{
					RawText:       match,
					CanonicalID:   norm.CanonicalID,
					CanonicalName: norm.CanonicalName,
					Domain:        norm.Domain,
					Category:      norm.Category,
					Confidence:    confidenceFromMatchType(norm.MatchType, norm.FuzzyScore),
					MatchType:     norm.MatchType,
				})
			}
		}
		// Blank out matched text to avoid double-counting.
		remaining = re.ReplaceAllStringFunc(remaining, func(s string) string {
			return strings.Repeat(" ", len(s))
		})
	}

	// Step 2: Single-token matching on the remaining text.
	tokens := e.tokenSplitter.Split(remaining, -1)
	for _, token := range tokens {
		token = cleanToken(token)
		if len(token) < 2 {
			continue
		}
		norm := e.taxonomy.Normalize(token)
		if norm.MatchType != "none" && !seenIDs[norm.CanonicalID] {
			seenIDs[norm.CanonicalID] = true
			allSkills = append(allSkills, ExtractedSkill{
				RawText:       token,
				CanonicalID:   norm.CanonicalID,
				CanonicalName: norm.CanonicalName,
				Domain:        norm.Domain,
				Category:      norm.Category,
				Confidence:    confidenceFromMatchType(norm.MatchType, norm.FuzzyScore),
				MatchType:     norm.MatchType,
			})
		} else if norm.MatchType == "none" && includeUnknown {
			// Collect unknown skills (tokens that look like skill names).
			if looksLikeSkill(token) && !seenRaw[strings.ToLower(token)] {
				seenRaw[strings.ToLower(token)] = true
				allSkills = append(allSkills, ExtractedSkill{
					RawText:    token,
					Confidence: 0.3,
					MatchType:  "unknown",
				})
			}
		}
	}

	return groupSkills(allSkills)
}

// groupSkills partitions extracted skills into technical, soft, domain, and unknown.
func groupSkills(skills []ExtractedSkill) ExtractionResult {
	result := ExtractionResult{Skills: skills}
	for _, s := range skills {
		switch {
		case s.MatchType == "unknown":
			result.UnknownSkills = append(result.UnknownSkills, s)
		case s.Domain == DomainManagement || s.Domain == DomainCommunication:
			result.SoftSkills = append(result.SoftSkills, s)
		case s.Domain == DomainDomain:
			result.DomainSkills = append(result.DomainSkills, s)
		default:
			result.TechnicalSkills = append(result.TechnicalSkills, s)
		}
	}
	return result
}

// confidenceFromMatchType returns a confidence score based on match type.
func confidenceFromMatchType(matchType string, fuzzyScore float64) float64 {
	switch matchType {
	case "exact":
		return 1.0
	case "alias":
		return 0.95
	case "fuzzy":
		return fuzzyScore * 0.9
	default:
		return 0.3
	}
}

// looksLikeSkill returns true if a token looks like it could be a skill name
// (heuristic: not a common English word, not a number, reasonable length).
func looksLikeSkill(token string) bool {
	if len(token) < 2 || len(token) > 50 {
		return false
	}
	// Skip pure numbers.
	allDigits := true
	for _, r := range token {
		if !unicode.IsDigit(r) {
			allDigits = false
			break
		}
	}
	if allDigits {
		return false
	}
	// Skip common English stop words.
	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "with": true, "that": true,
		"this": true, "are": true, "you": true, "have": true, "will": true,
		"from": true, "they": true, "been": true, "has": true, "not": true,
		"but": true, "what": true, "all": true, "were": true, "when": true,
		"your": true, "can": true, "said": true, "there": true, "use": true,
		"each": true, "which": true, "she": true, "how": true, "their": true,
		"about": true, "out": true, "many": true, "then": true, "them": true,
		"these": true, "some": true, "her": true, "would": true, "make": true,
		"like": true, "into": true, "time": true, "look": true, "two": true,
		"more": true, "write": true, "see": true, "number": true, "way": true,
		"could": true, "people": true, "than": true, "first": true, "water": true,
		"call": true, "who": true, "oil": true, "its": true,
		"now": true, "find": true, "long": true, "down": true, "day": true,
		"did": true, "get": true, "come": true, "made": true, "may": true,
		"part": true, "over": true, "new": true, "sound": true, "take": true,
		"only": true, "little": true, "work": true, "know": true, "place": true,
		"years": true, "live": true, "back": true, "give": true, "most": true,
		"very": true, "after": true, "thing": true, "our": true, "just": true,
		"name": true, "good": true, "sentence": true, "man": true, "think": true,
		"say": true, "great": true, "where": true, "help": true, "through": true,
		"much": true, "before": true, "line": true, "right": true, "too": true,
		"mean": true, "old": true, "any": true, "same": true, "tell": true,
		"boy": true, "follow": true, "came": true, "want": true, "show": true,
		"also": true, "around": true, "form": true, "three": true, "small": true,
		"set": true, "put": true, "end": true, "does": true, "another": true,
		"well": true, "large": true, "need": true, "big": true, "high": true,
		"such": true, "even": true, "because": true, "turn": true, "here": true,
		"why": true, "ask": true, "went": true, "men": true, "read": true,
		"land": true, "different": true, "home": true, "move": true, "try": true,
		"kind": true, "hand": true, "picture": true, "again": true, "change": true,
		"off": true, "play": true, "spell": true, "air": true, "away": true,
		"animal": true, "house": true, "point": true, "page": true, "letter": true,
		"mother": true, "answer": true, "found": true, "study": true, "still": true,
		"learn": true, "plant": true, "cover": true, "food": true, "sun": true,
		"four": true, "between": true, "state": true, "keep": true, "eye": true,
		"never": true, "last": true, "let": true, "thought": true, "city": true,
		"tree": true, "cross": true, "farm": true, "hard": true, "start": true,
		"might": true, "story": true, "saw": true, "far": true, "sea": true,
		"draw": true, "left": true, "late": true, "run": true, "while": true,
		"press": true, "close": true, "night": true, "real": true, "life": true,
		"few": true, "north": true, "open": true, "seem": true, "together": true,
		"next": true, "white": true, "children": true, "begin": true, "got": true,
		"walk": true, "example": true, "ease": true, "paper": true, "group": true,
		"always": true, "music": true, "those": true, "both": true, "mark": true,
		"often": true, "until": true, "mile": true, "river": true,
		"car": true, "feet": true, "care": true, "second": true, "enough": true,
		"plain": true, "girl": true, "usual": true, "young": true, "ready": true,
		"above": true, "ever": true, "red": true, "list": true, "though": true,
		"feel": true, "talk": true, "bird": true, "soon": true, "body": true,
		"dog": true, "family": true, "direct": true, "pose": true, "leave": true,
		"song": true, "measure": true, "door": true, "product": true, "black": true,
		"short": true, "numeral": true, "class": true, "wind": true, "question": true,
		"happen": true, "complete": true, "ship": true, "area": true, "half": true,
		"rock": true, "order": true, "fire": true, "south": true, "problem": true,
		"piece": true, "told": true, "knew": true, "pass": true, "since": true,
		"top": true, "whole": true, "king": true, "space": true, "heard": true,
		"best": true, "hour": true, "better": true, "true": true, "during": true,
		"hundred": true, "five": true, "remember": true, "step": true, "early": true,
		"hold": true, "west": true, "ground": true, "interest": true, "reach": true,
		"fast": true, "verb": true, "sing": true, "listen": true, "six": true,
		"table": true, "travel": true, "less": true, "morning": true, "ten": true,
		"simple": true, "several": true, "vowel": true, "toward": true, "war": true,
		"lay": true, "against": true, "pattern": true, "slow": true, "center": true,
		"love": true, "person": true, "money": true, "serve": true, "appear": true,
		"road": true, "map": true, "rain": true, "rule": true, "govern": true,
		"pull": true, "cold": true, "notice": true, "voice": true, "unit": true,
		"power": true, "town": true, "fine": true, "drive": true, "lead": true,
		"cry": true, "dark": true, "machine": true, "note": true, "wait": true,
		"plan": true, "figure": true, "star": true, "box": true, "noun": true,
		"field": true, "correct": true, "able": true, "pound": true,
		"done": true, "beauty": true, "stood": true, "contain": true,
		"front": true, "teach": true, "week": true, "final": true, "gave": true,
		"green": true, "oh": true, "quick": true, "develop": true, "ocean": true,
		"warm": true, "free": true, "minute": true, "strong": true, "special": true,
		"mind": true, "behind": true, "clear": true, "tail": true, "produce": true,
		"fact": true, "street": true, "inch": true, "multiply": true, "nothing": true,
		"course": true, "stay": true, "wheel": true, "full": true, "force": true,
		"blue": true, "object": true, "decide": true, "surface": true, "deep": true,
		"moon": true, "island": true, "foot": true, "system": true, "busy": true,
		"test": true, "record": true, "boat": true, "common": true, "gold": true,
		"possible": true, "plane": true, "stead": true, "dry": true, "wonder": true,
		"laugh": true, "thousand": true, "ago": true, "ran": true, "check": true,
		"game": true, "shape": true, "equate": true, "hot": true, "miss": true,
		"brought": true, "heat": true, "snow": true, "tire": true, "bring": true,
		"yes": true, "distant": true, "fill": true, "east": true, "paint": true,
		"language": true, "among": true,
	}
	return !stopWords[strings.ToLower(token)]
}

// cleanToken removes punctuation from the edges of a token.
func cleanToken(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".,;:\"'()[]{}!?-_/\\")
	return strings.TrimSpace(s)
}

// ─────────────────────────────────────────────────────────────────────────────
// Jaro-Winkler similarity
// ─────────────────────────────────────────────────────────────────────────────

// jaroWinkler computes the Jaro-Winkler similarity between two strings.
// Returns a value in [0.0, 1.0] where 1.0 is an exact match.
func jaroWinkler(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	jaro := jaroSimilarity(s1, s2)

	// Compute common prefix length (up to 4 characters).
	prefixLen := 0
	maxPrefix := 4
	if len(s1) < maxPrefix {
		maxPrefix = len(s1)
	}
	if len(s2) < maxPrefix {
		maxPrefix = len(s2)
	}
	for i := 0; i < maxPrefix; i++ {
		if s1[i] == s2[i] {
			prefixLen++
		} else {
			break
		}
	}

	// Jaro-Winkler scaling factor p = 0.1.
	return jaro + float64(prefixLen)*0.1*(1.0-jaro)
}

// jaroSimilarity computes the Jaro similarity between two strings.
func jaroSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	len1, len2 := len(s1), len(s2)
	matchDist := max(len1, len2)/2 - 1
	if matchDist < 0 {
		matchDist = 0
	}

	s1Matches := make([]bool, len1)
	s2Matches := make([]bool, len2)

	matches := 0
	transpositions := 0

	for i := 0; i < len1; i++ {
		start := max(0, i-matchDist)
		end := min(i+matchDist+1, len2)
		for j := start; j < end; j++ {
			if s2Matches[j] || s1[i] != s2[j] {
				continue
			}
			s1Matches[i] = true
			s2Matches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	k := 0
	for i := 0; i < len1; i++ {
		if !s1Matches[i] {
			continue
		}
		for !s2Matches[k] {
			k++
		}
		if s1[i] != s2[k] {
			transpositions++
		}
		k++
	}

	m := float64(matches)
	return (m/float64(len1) + m/float64(len2) + (m-float64(transpositions)/2)/m) / 3.0
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// normalise lowercases and trims a string for comparison.
func normalise(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// roundTo4 rounds a float64 to 4 decimal places.
func roundTo4(v float64) float64 {
	return math.Round(v*10000) / 10000
}

// max returns the larger of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// sortByLengthDesc sorts a string slice by length descending (in-place).
func sortByLengthDesc(ss []string) {
	// Simple insertion sort – the slice is typically small.
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && len(ss[j]) > len(ss[j-1]); j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
