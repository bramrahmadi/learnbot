package extractor

import (
	"regexp"
	"strings"

	"github.com/learnbot/resume-parser/internal/schema"
)

var (
	urlRe = regexp.MustCompile(`https?://[^\s,;]+`)

	// Technology extraction from project descriptions
	techInProjectRe = regexp.MustCompile(`(?i)(?:built\s+with|using|technologies?:|tech\s+stack:|stack:)\s*([^.\n]+)`)
)

// ExtractProjects parses project and achievement entries from the projects section text.
func ExtractProjects(text string) []schema.Project {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	blocks := splitProjectBlocks(text)
	var projects []schema.Project

	for _, block := range blocks {
		proj := parseProjectBlock(block)
		if proj != nil {
			projects = append(projects, *proj)
		}
	}

	return projects
}

// splitProjectBlocks splits the projects section into individual project entries.
func splitProjectBlocks(text string) []string {
	lines := strings.Split(text, "\n")
	var blocks []string
	var current []string
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if inBlock {
				current = append(current, line)
			}
			continue
		}

		// A new project starts with a bullet or a short title line (not indented)
		isNewProject := isBulletLine(trimmed) ||
			(!strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") &&
				len(trimmed) < 80 && !strings.HasPrefix(trimmed, "-"))

		if isNewProject && inBlock {
			if block := strings.TrimSpace(strings.Join(current, "\n")); block != "" {
				blocks = append(blocks, block)
			}
			current = []string{line}
			inBlock = true
			continue
		}

		current = append(current, line)
		inBlock = true
	}

	if block := strings.TrimSpace(strings.Join(current, "\n")); block != "" {
		blocks = append(blocks, block)
	}

	return blocks
}

// parseProjectBlock extracts structured data from a single project block.
func parseProjectBlock(block string) *schema.Project {
	lines := strings.Split(block, "\n")
	proj := &schema.Project{}

	var descLines []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimLeft(line, "•-*>·")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// URL
		if u := urlRe.FindString(line); u != "" {
			proj.URL = u
			// Remove URL from line
			line = strings.TrimSpace(urlRe.ReplaceAllString(line, ""))
		}

		// Date
		if m := dateRangeRe.FindStringSubmatch(line); len(m) >= 3 {
			proj.Date = strings.TrimSpace(m[1]) + " - " + strings.TrimSpace(m[2])
			line = strings.TrimSpace(dateRangeRe.ReplaceAllString(line, ""))
		} else if years := yearRe.FindAllString(line, -1); len(years) > 0 && proj.Date == "" {
			proj.Date = years[len(years)-1]
		}

		// Technologies
		if m := techInProjectRe.FindStringSubmatch(line); len(m) > 1 {
			proj.Technologies = splitTechnologies(m[1])
			line = strings.TrimSpace(techInProjectRe.ReplaceAllString(line, ""))
		}

		// First non-empty line is the project name
		if i == 0 || proj.Name == "" {
			if line != "" {
				proj.Name = line
			}
			continue
		}

		if line != "" {
			descLines = append(descLines, line)
		}
	}

	proj.Description = strings.Join(descLines, " ")

	// If no technologies found yet, try to extract from description
	if len(proj.Technologies) == 0 && proj.Description != "" {
		proj.Technologies = extractTechFromText(proj.Description)
	}

	if proj.Name == "" {
		return nil
	}

	// Confidence scoring
	score := 0.0
	total := 3.0
	if proj.Name != "" {
		score += 1.0
	}
	if proj.Description != "" {
		score += 1.0
	}
	if len(proj.Technologies) > 0 || proj.URL != "" {
		score += 1.0
	}
	proj.Confidence = schema.ConfidenceScore(score / total)

	return proj
}

// splitTechnologies splits a comma/slash-separated technology string.
func splitTechnologies(s string) []string {
	parts := regexp.MustCompile(`[,/|]+`).Split(s, -1)
	var techs []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			techs = append(techs, p)
		}
	}
	return techs
}

// extractTechFromText finds known technical skills mentioned in a description.
func extractTechFromText(text string) []string {
	lower := strings.ToLower(text)
	var found []string
	seen := map[string]bool{}

	for skill := range technicalSkills {
		// Use word boundary matching
		pattern := `\b` + regexp.QuoteMeta(skill) + `\b`
		if matched, _ := regexp.MatchString(pattern, lower); matched {
			if !seen[skill] {
				seen[skill] = true
				found = append(found, skill)
			}
		}
	}
	return found
}
