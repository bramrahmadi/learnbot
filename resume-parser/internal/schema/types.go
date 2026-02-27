// Package schema defines the core data structures for the resume parser.
package schema

import "time"

// ConfidenceScore represents the confidence level (0.0 - 1.0) of an extracted field.
type ConfidenceScore float64

const (
	ConfidenceHigh   ConfidenceScore = 0.9
	ConfidenceMedium ConfidenceScore = 0.7
	ConfidenceLow    ConfidenceScore = 0.5
)

// PersonalInfo holds extracted personal contact information.
type PersonalInfo struct {
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	Location   string          `json:"location"`
	LinkedIn   string          `json:"linkedin,omitempty"`
	GitHub     string          `json:"github,omitempty"`
	Website    string          `json:"website,omitempty"`
	Confidence ConfidenceScore `json:"confidence"`
}

// WorkExperience represents a single job entry.
type WorkExperience struct {
	Company          string          `json:"company"`
	Title            string          `json:"title"`
	StartDate        string          `json:"start_date"`
	EndDate          string          `json:"end_date"` // "Present" if current
	IsCurrent        bool            `json:"is_current"`
	Location         string          `json:"location,omitempty"`
	Responsibilities []string        `json:"responsibilities"`
	Confidence       ConfidenceScore `json:"confidence"`
}

// Education represents a single educational entry.
type Education struct {
	Institution string          `json:"institution"`
	Degree      string          `json:"degree"`
	Field       string          `json:"field,omitempty"`
	StartDate   string          `json:"start_date,omitempty"`
	EndDate     string          `json:"end_date,omitempty"`
	GPA         string          `json:"gpa,omitempty"`
	Honors      string          `json:"honors,omitempty"`
	Confidence  ConfidenceScore `json:"confidence"`
}

// Skill represents a single skill with its category.
type Skill struct {
	Name       string          `json:"name"`
	Category   string          `json:"category"` // "technical", "soft", "language", "tool"
	Confidence ConfidenceScore `json:"confidence"`
}

// Certification represents a professional certification or license.
type Certification struct {
	Name       string          `json:"name"`
	Issuer     string          `json:"issuer,omitempty"`
	Date       string          `json:"date,omitempty"`
	ExpiryDate string          `json:"expiry_date,omitempty"`
	ID         string          `json:"id,omitempty"`
	Confidence ConfidenceScore `json:"confidence"`
}

// Project represents a project or achievement entry.
type Project struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Technologies []string        `json:"technologies,omitempty"`
	URL          string          `json:"url,omitempty"`
	Date         string          `json:"date,omitempty"`
	Confidence   ConfidenceScore `json:"confidence"`
}

// ParsedResume is the top-level structured output of the resume parser.
type ParsedResume struct {
	// Metadata
	ParsedAt      time.Time `json:"parsed_at"`
	SourceFile    string    `json:"source_file"`
	FileType      string    `json:"file_type"` // "pdf" or "docx"
	ParserVersion string    `json:"parser_version"`

	// Extracted sections
	PersonalInfo   PersonalInfo     `json:"personal_info"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Education      []Education      `json:"education"`
	Skills         []Skill          `json:"skills"`
	Certifications []Certification  `json:"certifications"`
	Projects       []Project        `json:"projects"`
	Summary        string           `json:"summary,omitempty"`

	// Overall quality metrics
	OverallConfidence ConfidenceScore `json:"overall_confidence"`
	SectionsFound     []string        `json:"sections_found"`
	Warnings          []string        `json:"warnings,omitempty"`
	RawText           string          `json:"raw_text,omitempty"`
}

// ParseError represents a structured parsing error.
type ParseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Section string `json:"section,omitempty"`
}

func (e *ParseError) Error() string {
	if e.Section != "" {
		return e.Code + ": " + e.Message + " (section: " + e.Section + ")"
	}
	return e.Code + ": " + e.Message
}

// ParseRequest is the input to the parser API.
type ParseRequest struct {
	FileName    string `json:"file_name"`
	FileContent []byte `json:"file_content"`
	FileType    string `json:"file_type"` // "pdf" or "docx"
	IncludeRaw  bool   `json:"include_raw,omitempty"`
}

// ParseResponse wraps the parsed resume and any errors.
type ParseResponse struct {
	Success bool         `json:"success"`
	Data    *ParsedResume `json:"data,omitempty"`
	Error   *ParseError  `json:"error,omitempty"`
}

// UserProfile is the enriched profile built from a ParsedResume,
// used downstream by the RAG pipeline.
type UserProfile struct {
	ID             string           `json:"id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	PersonalInfo   PersonalInfo     `json:"personal_info"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Education      []Education      `json:"education"`
	Skills         []Skill          `json:"skills"`
	Certifications []Certification  `json:"certifications"`
	Projects       []Project        `json:"projects"`
	Summary        string           `json:"summary,omitempty"`
	YearsOfExp     float64          `json:"years_of_experience"`
	SkillTaxonomy  map[string][]string `json:"skill_taxonomy,omitempty"`
}
