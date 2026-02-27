// Package taxonomy provides a standardized skill taxonomy system for the
// LearnBot platform. It supports:
//
//   - A hierarchical skill ontology (domain → category → skill)
//   - Synonym/alias grouping (e.g. "JS", "JavaScript", "ECMAScript")
//   - Skill prerequisite relationships
//   - NLP-based extraction from free-form job description text
//   - Fuzzy matching for skill normalization
//   - Mapping of raw user skill strings to canonical taxonomy entries
package taxonomy

// ─────────────────────────────────────────────────────────────────────────────
// Taxonomy node types
// ─────────────────────────────────────────────────────────────────────────────

// Domain is the top-level grouping (e.g. "Engineering", "Data Science").
type Domain string

const (
	DomainEngineering    Domain = "engineering"
	DomainDataScience    Domain = "data_science"
	DomainDevOps         Domain = "devops"
	DomainDesign         Domain = "design"
	DomainManagement     Domain = "management"
	DomainCommunication  Domain = "communication"
	DomainDomain         Domain = "domain_knowledge"
)

// Category is the second-level grouping within a domain
// (e.g. "Frontend", "Backend", "Database").
type Category string

const (
	// Engineering categories
	CategoryLanguage    Category = "language"
	CategoryFrontend    Category = "frontend"
	CategoryBackend     Category = "backend"
	CategoryMobile      Category = "mobile"
	CategoryDatabase    Category = "database"
	CategoryCloud       Category = "cloud"
	CategoryDevOps      Category = "devops"
	CategorySecurity    Category = "security"
	CategoryTesting     Category = "testing"
	CategoryAPI         Category = "api"
	CategoryMessaging   Category = "messaging"

	// Data Science categories
	CategoryMLFramework Category = "ml_framework"
	CategoryDataTools   Category = "data_tools"
	CategoryMLConcept   Category = "ml_concept"

	// Soft skill categories
	CategoryLeadership  Category = "leadership"
	CategoryCollaboration Category = "collaboration"
	CategoryCommunication Category = "communication_skill"
	CategoryProblemSolving Category = "problem_solving"
	CategoryProjectMgmt Category = "project_management"

	// Domain knowledge categories
	CategoryFinance     Category = "finance"
	CategoryHealthcare  Category = "healthcare"
	CategoryEcommerce   Category = "ecommerce"
	CategoryLegal       Category = "legal"
)

// SkillNode represents a single skill entry in the taxonomy.
type SkillNode struct {
	// ID is the canonical identifier (lowercase, hyphenated).
	ID string `json:"id"`

	// CanonicalName is the preferred display name.
	CanonicalName string `json:"canonical_name"`

	// Domain is the top-level grouping.
	Domain Domain `json:"domain"`

	// Category is the second-level grouping.
	Category Category `json:"category"`

	// Aliases lists all known synonyms and alternate spellings.
	Aliases []string `json:"aliases,omitempty"`

	// Prerequisites lists IDs of skills that are typically learned before this one.
	Prerequisites []string `json:"prerequisites,omitempty"`

	// RelatedSkills lists IDs of skills that are commonly used alongside this one.
	RelatedSkills []string `json:"related_skills,omitempty"`

	// Description is a short human-readable description.
	Description string `json:"description,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Extraction types
// ─────────────────────────────────────────────────────────────────────────────

// ExtractedSkill is a skill found in a piece of text, with its taxonomy mapping.
type ExtractedSkill struct {
	// RawText is the exact text span that was matched.
	RawText string `json:"raw_text"`

	// CanonicalID is the taxonomy node ID this skill maps to (empty if unknown).
	CanonicalID string `json:"canonical_id,omitempty"`

	// CanonicalName is the preferred display name from the taxonomy.
	CanonicalName string `json:"canonical_name,omitempty"`

	// Domain is the top-level grouping.
	Domain Domain `json:"domain,omitempty"`

	// Category is the second-level grouping.
	Category Category `json:"category,omitempty"`

	// Confidence is the extraction confidence [0.0, 1.0].
	Confidence float64 `json:"confidence"`

	// MatchType describes how the skill was matched:
	// "exact", "alias", "fuzzy", "pattern", "unknown".
	MatchType string `json:"match_type"`
}

// ExtractionResult holds all skills extracted from a piece of text.
type ExtractionResult struct {
	// Skills is the deduplicated list of extracted skills.
	Skills []ExtractedSkill `json:"skills"`

	// TechnicalSkills is the subset of technical skills.
	TechnicalSkills []ExtractedSkill `json:"technical_skills"`

	// SoftSkills is the subset of soft skills.
	SoftSkills []ExtractedSkill `json:"soft_skills"`

	// DomainSkills is the subset of domain knowledge skills.
	DomainSkills []ExtractedSkill `json:"domain_skills"`

	// UnknownSkills is the subset of skills not found in the taxonomy.
	UnknownSkills []ExtractedSkill `json:"unknown_skills,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Normalization types
// ─────────────────────────────────────────────────────────────────────────────

// NormalizeResult is the output of normalizing a raw skill string.
type NormalizeResult struct {
	// Input is the original raw skill string.
	Input string `json:"input"`

	// CanonicalID is the matched taxonomy node ID (empty if no match).
	CanonicalID string `json:"canonical_id,omitempty"`

	// CanonicalName is the preferred display name.
	CanonicalName string `json:"canonical_name,omitempty"`

	// Domain is the top-level grouping.
	Domain Domain `json:"domain,omitempty"`

	// Category is the second-level grouping.
	Category Category `json:"category,omitempty"`

	// MatchType is how the match was found: "exact", "alias", "fuzzy", "none".
	MatchType string `json:"match_type"`

	// FuzzyScore is the similarity score [0.0, 1.0] for fuzzy matches.
	FuzzyScore float64 `json:"fuzzy_score,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// API request/response types
// ─────────────────────────────────────────────────────────────────────────────

// ExtractRequest is the input to the skill extraction API.
type ExtractRequest struct {
	// Text is the free-form text to extract skills from (e.g. job description).
	Text string `json:"text"`

	// IncludeUnknown controls whether skills not in the taxonomy are included.
	IncludeUnknown bool `json:"include_unknown,omitempty"`
}

// ExtractResponse is the output of the skill extraction API.
type ExtractResponse struct {
	Success bool              `json:"success"`
	Data    *ExtractionResult `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// NormalizeRequest is the input to the skill normalization API.
type NormalizeRequest struct {
	// Skills is the list of raw skill strings to normalize.
	Skills []string `json:"skills"`
}

// NormalizeResponse is the output of the skill normalization API.
type NormalizeResponse struct {
	Success bool              `json:"success"`
	Data    []NormalizeResult `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// LookupRequest is the input to the taxonomy lookup API.
type LookupRequest struct {
	// ID is the canonical skill ID to look up.
	ID string `json:"id"`
}

// LookupResponse is the output of the taxonomy lookup API.
type LookupResponse struct {
	Success bool       `json:"success"`
	Data    *SkillNode `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}

// SearchRequest is the input to the taxonomy search API.
type SearchRequest struct {
	// Query is the search term.
	Query string `json:"query"`

	// Domain filters results to a specific domain (optional).
	Domain Domain `json:"domain,omitempty"`

	// Category filters results to a specific category (optional).
	Category Category `json:"category,omitempty"`

	// Limit is the maximum number of results to return (default 20).
	Limit int `json:"limit,omitempty"`
}

// SearchResponse is the output of the taxonomy search API.
type SearchResponse struct {
	Success bool        `json:"success"`
	Data    []SkillNode `json:"data,omitempty"`
	Total   int         `json:"total"`
	Error   string      `json:"error,omitempty"`
}
