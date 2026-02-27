package extractor

import (
	"regexp"
	"strings"
)

// SectionType identifies a resume section.
type SectionType string

const (
	SectionSummary        SectionType = "summary"
	SectionExperience     SectionType = "experience"
	SectionEducation      SectionType = "education"
	SectionSkills         SectionType = "skills"
	SectionCertifications SectionType = "certifications"
	SectionProjects       SectionType = "projects"
	SectionUnknown        SectionType = "unknown"
)

// Section holds the raw text of a resume section.
type Section struct {
	Type  SectionType
	Title string
	Text  string
}

var sectionHeaderRe = regexp.MustCompile(`(?im)^[\s\-=_*#]*([A-Z][A-Za-z\s&/]+?)[\s\-=_*#:]*$`)

// sectionKeywords maps normalized header keywords to section types.
var sectionKeywords = map[string]SectionType{
	"summary":                SectionSummary,
	"professional summary":   SectionSummary,
	"career summary":         SectionSummary,
	"objective":              SectionSummary,
	"career objective":       SectionSummary,
	"profile":                SectionSummary,
	"about":                  SectionSummary,
	"about me":               SectionSummary,
	"experience":             SectionExperience,
	"work experience":        SectionExperience,
	"professional experience": SectionExperience,
	"employment history":     SectionExperience,
	"employment":             SectionExperience,
	"work history":           SectionExperience,
	"career history":         SectionExperience,
	"positions held":         SectionExperience,
	"education":              SectionEducation,
	"educational background": SectionEducation,
	"academic background":    SectionEducation,
	"academic history":       SectionEducation,
	"qualifications":         SectionEducation,
	"skills":                 SectionSkills,
	"technical skills":       SectionSkills,
	"core competencies":      SectionSkills,
	"competencies":           SectionSkills,
	"key skills":             SectionSkills,
	"areas of expertise":     SectionSkills,
	"expertise":              SectionSkills,
	"technologies":           SectionSkills,
	"tools":                  SectionSkills,
	"certifications":         SectionCertifications,
	"licenses":               SectionCertifications,
	"certifications & licenses": SectionCertifications,
	"professional certifications": SectionCertifications,
	"credentials":            SectionCertifications,
	"projects":               SectionProjects,
	"personal projects":      SectionProjects,
	"key projects":           SectionProjects,
	"achievements":           SectionProjects,
	"accomplishments":        SectionProjects,
	"awards":                 SectionProjects,
	"publications":           SectionProjects,
	"portfolio":              SectionProjects,
}

// SplitSections splits raw resume text into labeled sections.
func SplitSections(text string) []Section {
	lines := strings.Split(text, "\n")
	var sections []Section
	var currentType SectionType = SectionUnknown
	var currentTitle string
	var currentLines []string

	flush := func() {
		if len(currentLines) > 0 {
			sections = append(sections, Section{
				Type:  currentType,
				Title: currentTitle,
				Text:  strings.TrimSpace(strings.Join(currentLines, "\n")),
			})
		}
		currentLines = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			currentLines = append(currentLines, line)
			continue
		}

		// Check if this line is a section header
		if st, title, ok := detectSectionHeader(trimmed); ok {
			flush()
			currentType = st
			currentTitle = title
			continue
		}

		currentLines = append(currentLines, line)
	}
	flush()

	return sections
}

// detectSectionHeader returns the section type if the line is a header.
func detectSectionHeader(line string) (SectionType, string, bool) {
	// Must be short (section headers are rarely > 50 chars)
	if len(line) > 60 {
		return SectionUnknown, "", false
	}

	// Strip decorators
	clean := strings.Trim(line, "-=_*#: \t")
	clean = strings.TrimSpace(clean)
	if clean == "" {
		return SectionUnknown, "", false
	}

	normalized := strings.ToLower(clean)
	if st, ok := sectionKeywords[normalized]; ok {
		return st, clean, true
	}

	// Partial match: check if any keyword is a prefix/suffix
	for kw, st := range sectionKeywords {
		if strings.HasPrefix(normalized, kw) || strings.HasSuffix(normalized, kw) {
			return st, clean, true
		}
	}

	// All-caps short line heuristic (e.g., "EXPERIENCE", "EDUCATION")
	if isAllCaps(clean) && len(strings.Fields(clean)) <= 4 {
		normalized = strings.ToLower(clean)
		if st, ok := sectionKeywords[normalized]; ok {
			return st, clean, true
		}
		// Try partial
		for kw, st := range sectionKeywords {
			if strings.Contains(normalized, kw) {
				return st, clean, true
			}
		}
	}

	return SectionUnknown, "", false
}

// isAllCaps returns true if all alphabetic characters are uppercase.
func isAllCaps(s string) bool {
	hasAlpha := false
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return false
		}
		if r >= 'A' && r <= 'Z' {
			hasAlpha = true
		}
	}
	return hasAlpha
}

// GetSectionText returns the combined text of all sections of a given type.
func GetSectionText(sections []Section, t SectionType) string {
	var parts []string
	for _, s := range sections {
		if s.Type == t {
			parts = append(parts, s.Text)
		}
	}
	return strings.Join(parts, "\n")
}

// ListFoundSections returns the unique section types found.
func ListFoundSections(sections []Section) []string {
	seen := map[SectionType]bool{}
	var result []string
	for _, s := range sections {
		if s.Type != SectionUnknown && !seen[s.Type] {
			seen[s.Type] = true
			result = append(result, string(s.Type))
		}
	}
	return result
}
