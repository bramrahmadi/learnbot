// Package extractor provides regex and heuristic-based field extraction from resume text.
package extractor

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/learnbot/resume-parser/internal/schema"
)

var (
	emailRe    = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	phoneRe    = regexp.MustCompile(`(?:\+?1[\s\-.]?)?\(?\d{3}\)?[\s\-.]?\d{3}[\s\-.]?\d{4}`)
	linkedinRe = regexp.MustCompile(`(?i)(?:linkedin\.com/in/|linkedin:\s*)([a-z0-9\-]+)`)
	githubRe   = regexp.MustCompile(`(?i)(?:github\.com/|github:\s*)([a-z0-9\-]+)`)
	websiteRe  = regexp.MustCompile(`(?i)https?://[^\s,;]+`)

	// Location patterns: "City, State" or "City, Country"
	locationRe = regexp.MustCompile(`(?i)([A-Z][a-zA-Z\s]+),\s*([A-Z]{2}|[A-Z][a-zA-Z\s]+)`)

	// Name heuristics: first non-empty line that looks like a proper name
	nameLineRe = regexp.MustCompile(`^[A-Z][a-z]+(?:\s+[A-Z]\.?)?\s+[A-Z][a-z]+(?:\s+[A-Z][a-z]+)?$`)
)

// ExtractPersonalInfo extracts personal contact information from raw resume text.
func ExtractPersonalInfo(text string) schema.PersonalInfo {
	info := schema.PersonalInfo{}
	confidence := 0.0
	fields := 0

	// Email
	if m := emailRe.FindString(text); m != "" {
		info.Email = strings.ToLower(m)
		confidence += 1.0
		fields++
	}

	// Phone
	if m := phoneRe.FindString(text); m != "" {
		info.Phone = normalizePhone(m)
		confidence += 1.0
		fields++
	}

	// LinkedIn
	if m := linkedinRe.FindStringSubmatch(text); len(m) > 1 {
		info.LinkedIn = "https://linkedin.com/in/" + m[1]
		confidence += 0.5
		fields++
	}

	// GitHub
	if m := githubRe.FindStringSubmatch(text); len(m) > 1 {
		info.GitHub = "https://github.com/" + m[1]
		confidence += 0.5
		fields++
	}

	// Website (non-linkedin, non-github)
	for _, url := range websiteRe.FindAllString(text, -1) {
		lower := strings.ToLower(url)
		if !strings.Contains(lower, "linkedin") && !strings.Contains(lower, "github") {
			info.Website = url
			break
		}
	}

	// Name: scan first 10 non-empty lines for a proper name pattern
	info.Name = extractName(text)
	if info.Name != "" {
		confidence += 1.0
		fields++
	}

	// Location
	info.Location = extractLocation(text)
	if info.Location != "" {
		confidence += 0.8
		fields++
	}

	// Compute confidence
	if fields > 0 {
		info.Confidence = schema.ConfidenceScore(confidence / float64(fields))
	}

	return info
}

// extractName scans the first lines of the resume for a proper name.
func extractName(text string) string {
	lines := strings.Split(text, "\n")
	checked := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		checked++
		if checked > 10 {
			break
		}
		// Skip lines that look like section headers or contact info
		if emailRe.MatchString(line) || phoneRe.MatchString(line) {
			continue
		}
		if nameLineRe.MatchString(line) {
			return line
		}
		// Fallback: if line has 2-4 words, all title-cased, treat as name
		words := strings.Fields(line)
		if len(words) >= 2 && len(words) <= 4 && allTitleCase(words) && !containsKeyword(line) {
			return line
		}
	}
	return ""
}

// extractLocation finds a city/state or city/country pattern near contact info.
func extractLocation(text string) string {
	// Look for location near email/phone context (first 500 chars)
	searchArea := text
	if len(text) > 800 {
		searchArea = text[:800]
	}
	if m := locationRe.FindString(searchArea); m != "" {
		return m
	}
	return ""
}

// normalizePhone strips extra whitespace from phone numbers.
func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

// allTitleCase returns true if every word starts with an uppercase letter.
func allTitleCase(words []string) bool {
	for _, w := range words {
		if len(w) == 0 {
			continue
		}
		if !unicode.IsUpper(rune(w[0])) {
			return false
		}
	}
	return true
}

// containsKeyword returns true if the line contains common resume section keywords.
func containsKeyword(line string) bool {
	keywords := []string{
		"experience", "education", "skills", "summary", "objective",
		"certifications", "projects", "achievements", "references",
		"profile", "contact", "address", "phone", "email",
	}
	lower := strings.ToLower(line)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
