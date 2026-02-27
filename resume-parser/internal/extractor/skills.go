package extractor

import (
	"regexp"
	"strings"

	"github.com/learnbot/resume-parser/internal/schema"
)

// technicalSkills is a curated set of known technical skills for classification.
var technicalSkills = map[string]bool{
	// Languages
	"go": true, "golang": true, "rust": true, "python": true, "java": true,
	"javascript": true, "typescript": true, "c": true, "c++": true, "c#": true,
	"ruby": true, "php": true, "swift": true, "kotlin": true, "scala": true,
	"r": true, "matlab": true, "perl": true, "bash": true, "shell": true,
	"powershell": true, "sql": true, "html": true, "css": true, "sass": true,
	"less": true, "xml": true, "json": true, "yaml": true, "graphql": true,

	// Frameworks & Libraries
	"react": true, "angular": true, "vue": true, "next.js": true, "nuxt": true,
	"node.js": true, "express": true, "django": true, "flask": true, "fastapi": true,
	"spring": true, "spring boot": true, "rails": true, "laravel": true,
	"gin": true, "echo": true, "fiber": true, "actix": true, "rocket": true,
	"tensorflow": true, "pytorch": true, "keras": true, "scikit-learn": true,
	"pandas": true, "numpy": true, "matplotlib": true, "seaborn": true,

	// Databases
	"postgresql": true, "mysql": true, "sqlite": true, "mongodb": true,
	"redis": true, "elasticsearch": true, "cassandra": true, "dynamodb": true,
	"oracle": true, "mssql": true, "mariadb": true, "neo4j": true,
	"pinecone": true, "weaviate": true, "qdrant": true,

	// Cloud & DevOps
	"aws": true, "azure": true, "gcp": true, "google cloud": true,
	"docker": true, "kubernetes": true, "k8s": true, "terraform": true,
	"ansible": true, "jenkins": true, "github actions": true, "gitlab ci": true,
	"circleci": true, "travis ci": true, "helm": true, "istio": true,
	"prometheus": true, "grafana": true, "datadog": true, "splunk": true,

	// Tools
	"git": true, "github": true, "gitlab": true, "bitbucket": true,
	"jira": true, "confluence": true, "slack": true, "figma": true,
	"postman": true, "swagger": true, "openapi": true, "grpc": true,
	"kafka": true, "rabbitmq": true, "celery": true, "airflow": true,
	"spark": true, "hadoop": true, "hive": true, "flink": true,

	// Concepts
	"rest": true, "restful": true, "microservices": true, "api": true,
	"ci/cd": true, "devops": true, "agile": true, "scrum": true,
	"tdd": true, "bdd": true, "oop": true, "functional programming": true,
	"machine learning": true, "deep learning": true, "nlp": true,
	"computer vision": true, "data science": true, "big data": true,
	"blockchain": true, "web3": true, "rag": true, "llm": true,
}

// softSkills is a curated set of known soft skills.
var softSkills = map[string]bool{
	"leadership": true, "communication": true, "teamwork": true,
	"problem solving": true, "critical thinking": true, "creativity": true,
	"adaptability": true, "time management": true, "project management": true,
	"collaboration": true, "mentoring": true, "coaching": true,
	"presentation": true, "negotiation": true, "conflict resolution": true,
	"analytical": true, "detail-oriented": true, "self-motivated": true,
	"organized": true, "multitasking": true, "fast learner": true,
	"customer service": true, "stakeholder management": true,
}

var (
	// Skill list separators
	skillSepRe = regexp.MustCompile(`[,;|•\t]+`)
	// Parenthetical content
	parenRe = regexp.MustCompile(`\([^)]*\)`)
)

// ExtractSkills parses skills from the skills section text.
func ExtractSkills(text string) []schema.Skill {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	rawSkills := tokenizeSkills(text)
	seen := map[string]bool{}
	var skills []schema.Skill

	for _, raw := range rawSkills {
		normalized := normalizeSkill(raw)
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true

		skill := schema.Skill{
			Name:       normalized,
			Category:   classifySkill(normalized),
			Confidence: computeSkillConfidence(normalized),
		}
		skills = append(skills, skill)
	}

	return skills
}

// tokenizeSkills splits the skills section into individual skill tokens.
func tokenizeSkills(text string) []string {
	// Remove parenthetical notes
	text = parenRe.ReplaceAllString(text, "")

	// Split by common separators
	parts := skillSepRe.Split(text, -1)
	var tokens []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Remove bullet characters
		part = strings.TrimLeft(part, "•-*>·")
		part = strings.TrimSpace(part)

		if part == "" || len(part) < 2 {
			continue
		}

		// If the part is a line with multiple words that looks like a category header, skip
		if isSkillCategoryHeader(part) {
			continue
		}

		// Split by newlines too
		for _, line := range strings.Split(part, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && len(line) >= 2 {
				tokens = append(tokens, line)
			}
		}
	}

	return tokens
}

// normalizeSkill cleans and normalizes a skill string.
func normalizeSkill(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".,;:\"'")
	s = strings.TrimSpace(s)

	// Skip if too long (likely a sentence, not a skill)
	if len(s) > 50 {
		return ""
	}

	// Skip if it's a number
	if len(s) <= 4 {
		allDigits := true
		for _, r := range s {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return ""
		}
	}

	return s
}

// classifySkill returns the category of a skill.
func classifySkill(skill string) string {
	lower := strings.ToLower(skill)

	if technicalSkills[lower] {
		return "technical"
	}
	if softSkills[lower] {
		return "soft"
	}

	// Heuristic: if it contains version numbers or is all-caps acronym, likely technical
	if regexp.MustCompile(`\d+\.\d+`).MatchString(skill) {
		return "technical"
	}
	if regexp.MustCompile(`^[A-Z]{2,6}$`).MatchString(skill) {
		return "technical"
	}

	return "other"
}

// computeSkillConfidence returns confidence based on whether the skill is in known lists.
func computeSkillConfidence(skill string) schema.ConfidenceScore {
	lower := strings.ToLower(skill)
	if technicalSkills[lower] || softSkills[lower] {
		return schema.ConfidenceHigh
	}
	return schema.ConfidenceMedium
}

// isSkillCategoryHeader returns true if the string looks like a category label.
func isSkillCategoryHeader(s string) bool {
	headers := []string{
		"technical skills", "soft skills", "programming languages",
		"frameworks", "tools", "databases", "cloud", "languages",
		"core competencies", "areas of expertise",
	}
	lower := strings.ToLower(s)
	for _, h := range headers {
		if lower == h {
			return true
		}
	}
	return false
}
