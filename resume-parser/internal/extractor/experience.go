package extractor

import (
	"regexp"
	"strings"

	"github.com/learnbot/resume-parser/internal/schema"
)

var (
	// Date range patterns: "Jan 2020 - Mar 2022", "2019 - Present", "01/2018 - 12/2020"
	dateRangeRe = regexp.MustCompile(
		`(?i)((?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|` +
			`Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?)` +
			`\.?\s+\d{4}|\d{1,2}/\d{4}|\d{4})` +
			`\s*[-–—to]+\s*` +
			`((?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|` +
			`Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?)` +
			`\.?\s+\d{4}|\d{1,2}/\d{4}|\d{4}|[Pp]resent|[Cc]urrent|[Nn]ow)`)

	// Bullet point patterns (using literal Unicode chars instead of \u escapes)
	bulletRe = regexp.MustCompile(`(?m)^[\s]*[•\-\*‣◦⁃∙>]\s+(.+)$`)

	// Job title keywords for heuristic detection
	titleKeywords = []string{
		"engineer", "developer", "manager", "director", "analyst", "designer",
		"architect", "consultant", "specialist", "coordinator", "lead", "senior",
		"junior", "intern", "associate", "officer", "executive", "president",
		"vice president", "vp", "cto", "ceo", "cfo", "head of", "principal",
		"staff", "scientist", "researcher", "administrator", "technician",
	}
)

// ExtractWorkExperience parses work experience entries from the experience section text.
func ExtractWorkExperience(text string) []schema.WorkExperience {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	blocks := splitExperienceBlocks(text)
	var experiences []schema.WorkExperience

	for _, block := range blocks {
		exp := parseExperienceBlock(block)
		if exp != nil {
			experiences = append(experiences, *exp)
		}
	}

	return experiences
}

// splitExperienceBlocks splits the experience section into individual job blocks.
// A new block starts when a non-indented, non-bullet line appears after a date line.
func splitExperienceBlocks(text string) []string {
	lines := strings.Split(text, "\n")
	var blocks []string
	var current []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			current = append(current, line)
			continue
		}

		// A new job block starts when we see a line that looks like a job header:
		// - Contains a separator like " | ", " at ", " - " with title keywords
		// - OR is a short non-bullet, non-indented line followed by a date
		isHeader := isJobHeaderLine(trimmed)
		if isHeader && len(current) > 0 {
			// Check if current block has meaningful content
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

// isJobHeaderLine returns true if the line looks like a job title/company header.
func isJobHeaderLine(line string) bool {
	// Must not be a bullet point
	if isBulletLine(line) {
		return false
	}
	// Must not be a date range itself
	if dateRangeRe.MatchString(line) && !strings.Contains(line, "|") && !strings.Contains(line, " at ") {
		return false
	}
	// Must contain a title keyword or separator
	lower := strings.ToLower(line)
	for _, kw := range titleKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	// Contains common job separators
	if strings.Contains(line, " | ") || strings.Contains(line, " at ") {
		return true
	}
	return false
}

// parseExperienceBlock extracts structured data from a single job block.
func parseExperienceBlock(block string) *schema.WorkExperience {
	lines := strings.Split(block, "\n")
	exp := &schema.WorkExperience{}

	var bodyLines []string
	headerLines := []string{}

	// First pass: collect header lines (before bullet points start)
	// and body lines (bullet points and descriptions)
	inBody := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if inBody {
				bodyLines = append(bodyLines, line)
			}
			continue
		}

		if isBulletLine(trimmed) {
			inBody = true
			bodyLines = append(bodyLines, line)
			continue
		}

		if inBody {
			bodyLines = append(bodyLines, line)
			continue
		}

		headerLines = append(headerLines, trimmed)
	}

	// Parse header lines for company, title, dates
	parseHeaderLines(headerLines, exp)

	// Extract responsibilities from bullet points
	bodyText := strings.Join(bodyLines, "\n")
	exp.Responsibilities = extractBullets(bodyText)

	// If no bullets, use non-empty body lines as responsibilities
	if len(exp.Responsibilities) == 0 {
		for _, l := range bodyLines {
			if l = strings.TrimSpace(l); l != "" {
				exp.Responsibilities = append(exp.Responsibilities, l)
			}
		}
	}

	// Skip entries with no meaningful data
	if exp.Company == "" && exp.Title == "" {
		return nil
	}

	// Confidence scoring
	score := 0.0
	total := 4.0
	if exp.Company != "" {
		score += 1.0
	}
	if exp.Title != "" {
		score += 1.0
	}
	if exp.StartDate != "" {
		score += 1.0
	}
	if len(exp.Responsibilities) > 0 {
		score += 0.8
	}

	exp.Confidence = schema.ConfidenceScore(score / total)

	return exp
}

// parseHeaderLines extracts company, title, and dates from the header lines of a job block.
func parseHeaderLines(lines []string, exp *schema.WorkExperience) {
	for _, line := range lines {
		// Check for date range
		if m := dateRangeRe.FindStringSubmatch(line); len(m) >= 3 {
			exp.StartDate = strings.TrimSpace(m[1])
			exp.EndDate = strings.TrimSpace(m[2])
			exp.IsCurrent = isCurrentDate(exp.EndDate)

			// Remove date from line to get remaining company/title info
			rest := strings.TrimSpace(dateRangeRe.ReplaceAllString(line, ""))
			rest = strings.Trim(rest, "|-–—,; \t")
			if rest != "" && exp.Company == "" && exp.Title == "" {
				parseCompanyTitle(rest, exp)
			}
			continue
		}

		// If no date yet, this is a company/title line
		if exp.Company == "" && exp.Title == "" {
			parseCompanyTitle(line, exp)
		} else if exp.Company != "" && exp.Title == "" {
			// Second header line might be the title
			exp.Title = line
		} else if exp.Title != "" && exp.Company == "" {
			exp.Company = line
		}
	}
}

// parseCompanyTitle attempts to split a line into company and title.
func parseCompanyTitle(line string, exp *schema.WorkExperience) {
	if line == "" {
		return
	}

	// Common separators: " at ", " | ", " - ", " – "
	separators := []string{" at ", " | ", " – ", " - ", ", "}
	for _, sep := range separators {
		parts := strings.SplitN(line, sep, 2)
		if len(parts) == 2 {
			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])
			// Determine which is title and which is company
			if looksLikeTitle(a) {
				exp.Title = a
				exp.Company = b
			} else {
				exp.Company = a
				exp.Title = b
			}
			return
		}
	}

	// No separator: if it looks like a title, assign as title; else company
	if looksLikeTitle(line) {
		exp.Title = line
	} else {
		exp.Company = line
	}
}

// looksLikeTitle returns true if the string resembles a job title.
func looksLikeTitle(s string) bool {
	lower := strings.ToLower(s)
	for _, kw := range titleKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// extractBullets extracts bullet point items from text.
func extractBullets(text string) []string {
	matches := bulletRe.FindAllStringSubmatch(text, -1)
	var bullets []string
	for _, m := range matches {
		if len(m) > 1 {
			b := strings.TrimSpace(m[1])
			if b != "" {
				bullets = append(bullets, b)
			}
		}
	}
	return bullets
}

// isCurrentDate returns true if the date string indicates a current position.
func isCurrentDate(s string) bool {
	lower := strings.ToLower(strings.TrimSpace(s))
	return lower == "present" || lower == "current" || lower == "now"
}
