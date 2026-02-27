package extractor

import (
	"regexp"
	"strings"

	"github.com/learnbot/resume-parser/internal/schema"
)

var (
	gpaRe = regexp.MustCompile(`(?i)GPA[:\s]+(\d+\.\d+)(?:\s*/\s*\d+\.\d+)?`)

	degreeKeywords = []string{
		"bachelor", "b.s.", "b.a.", "bs", "ba",
		"master", "m.s.", "m.a.", "ms", "ma", "mba",
		"doctor", "ph.d.", "phd", "md", "jd",
		"associate", "a.s.", "a.a.",
		"diploma", "certificate", "certification",
		"b.eng", "m.eng", "b.tech", "m.tech",
	}

	honorKeywords = []string{
		"cum laude", "magna cum laude", "summa cum laude",
		"with honors", "with distinction", "honors",
	}

	// Year-only date pattern
	yearRe = regexp.MustCompile(`\b(19|20)\d{2}\b`)
)

// ExtractEducation parses education entries from the education section text.
func ExtractEducation(text string) []schema.Education {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	blocks := splitEducationBlocks(text)
	var educations []schema.Education

	for _, block := range blocks {
		edu := parseEducationBlock(block)
		if edu != nil {
			educations = append(educations, *edu)
		}
	}

	return educations
}

// splitEducationBlocks splits education text into individual entries.
// A new block starts when we see a line that looks like an institution name
// (not a degree, not a date, not a GPA line) after some content.
func splitEducationBlocks(text string) []string {
	lines := strings.Split(text, "\n")
	var blocks []string
	var current []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			// Empty line may separate entries
			if len(current) > 0 {
				current = append(current, line)
			}
			continue
		}

		// A new block starts when we see a line that looks like an institution
		// (not a degree, not a date, not GPA) and we already have content
		if len(current) > 0 && looksLikeInstitution(trimmed) {
			if block := strings.TrimSpace(strings.Join(current, "\n")); block != "" {
				blocks = append(blocks, block)
			}
			current = []string{line}
			continue
		}

		current = append(current, line)
	}

	if block := strings.TrimSpace(strings.Join(current, "\n")); block != "" {
		blocks = append(blocks, block)
	}

	return blocks
}

// looksLikeInstitution returns true if the line looks like a university/school name.
func looksLikeInstitution(s string) bool {
	// Must not be a degree line
	if looksLikeDegree(s) {
		return false
	}
	// Must not be a date line
	if dateRangeRe.MatchString(s) {
		return false
	}
	// Must not be a GPA line
	if gpaRe.MatchString(s) {
		return false
	}
	// Must not be a bullet
	if isBulletLine(s) {
		return false
	}
	// Must not be an honor line
	lower := strings.ToLower(s)
	for _, h := range honorKeywords {
		if strings.Contains(lower, h) {
			return false
		}
	}

	// Institution keywords
	institutionKeywords := []string{
		"university", "college", "institute", "school", "academy",
		"polytechnic", "faculty", "campus",
	}
	for _, kw := range institutionKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}

	// Heuristic: title-cased, 2-6 words, no numbers
	words := strings.Fields(s)
	if len(words) >= 2 && len(words) <= 6 && allTitleCase(words) && !containsDigit(s) {
		return true
	}

	return false
}

// parseEducationBlock extracts structured data from a single education block.
func parseEducationBlock(block string) *schema.Education {
	lines := strings.Split(block, "\n")
	edu := &schema.Education{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// GPA
		if m := gpaRe.FindStringSubmatch(trimmed); len(m) > 1 {
			edu.GPA = m[1]
			continue
		}

		// Honors
		if h := extractHonors(trimmed); h != "" {
			edu.Honors = h
			continue
		}

		// Date range
		if m := dateRangeRe.FindStringSubmatch(trimmed); len(m) >= 3 {
			edu.StartDate = strings.TrimSpace(m[1])
			edu.EndDate = strings.TrimSpace(m[2])
			// Remove date from line to get remaining info
			rest := strings.TrimSpace(dateRangeRe.ReplaceAllString(trimmed, ""))
			rest = strings.Trim(rest, "|-–—,; \t")
			if rest != "" && edu.Institution == "" {
				edu.Institution = rest
			}
			continue
		}

		// Year only (e.g., "2020" or "2016 - 2020" already handled above)
		if years := yearRe.FindAllString(trimmed, -1); len(years) >= 1 && edu.StartDate == "" && edu.EndDate == "" {
			// Only treat as date if the line is mostly just years
			nonYearPart := yearRe.ReplaceAllString(trimmed, "")
			nonYearPart = strings.Trim(nonYearPart, " -–—to")
			if strings.TrimSpace(nonYearPart) == "" {
				if len(years) == 2 {
					edu.StartDate = years[0]
					edu.EndDate = years[1]
				} else {
					edu.EndDate = years[0]
				}
				continue
			}
		}

		// Degree line
		if looksLikeDegree(trimmed) {
			if edu.Degree == "" {
				edu.Degree, edu.Field = parseDegreeField(trimmed)
			}
			continue
		}

		// Institution (first non-degree, non-date, non-GPA line)
		if edu.Institution == "" {
			edu.Institution = trimmed
		}
	}

	if edu.Institution == "" && edu.Degree == "" {
		return nil
	}

	// Confidence scoring
	score := 0.0
	total := 3.0
	if edu.Institution != "" {
		score += 1.0
	}
	if edu.Degree != "" {
		score += 1.0
	}
	if edu.EndDate != "" || edu.StartDate != "" {
		score += 1.0
	}
	edu.Confidence = schema.ConfidenceScore(score / total)

	return edu
}

// looksLikeDegree returns true if the line contains a degree keyword.
// Uses word-boundary matching to avoid false positives like "technology" matching "tech".
func looksLikeDegree(s string) bool {
	lower := strings.ToLower(s)
	for _, kw := range degreeKeywords {
		// Use word-boundary matching
		pattern := `(?i)\b` + regexp.QuoteMeta(kw) + `\b`
		if matched, _ := regexp.MatchString(pattern, lower); matched {
			return true
		}
	}
	return false
}

// parseDegreeField splits a degree line into degree type and field of study.
func parseDegreeField(line string) (degree, field string) {
	// Patterns: "Bachelor of Science in Computer Science"
	//           "B.S. Computer Science"
	//           "Master of Business Administration"
	inRe := regexp.MustCompile(`(?i)\bin\b`)
	parts := inRe.Split(line, 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	// Try comma split: "B.S., Computer Science"
	commaParts := strings.SplitN(line, ",", 2)
	if len(commaParts) == 2 {
		return strings.TrimSpace(commaParts[0]), strings.TrimSpace(commaParts[1])
	}

	return strings.TrimSpace(line), ""
}

// extractHonors returns any honors string found in the line.
func extractHonors(line string) string {
	lower := strings.ToLower(line)
	for _, h := range honorKeywords {
		if strings.Contains(lower, h) {
			return h
		}
	}
	return ""
}

// containsDigit returns true if the string contains any digit.
func containsDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
