package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/learnbot/resume-parser/internal/extractor"
	"github.com/learnbot/resume-parser/internal/schema"
)

const parserVersion = "1.0.0"

// ResumeParser orchestrates document parsing and field extraction.
type ResumeParser struct {
	pdfParser  *PDFParser
	docxParser *DOCXParser
}

// NewResumeParser creates a new ResumeParser with all sub-parsers initialized.
func NewResumeParser() *ResumeParser {
	return &ResumeParser{
		pdfParser:  NewPDFParser(),
		docxParser: NewDOCXParser(),
	}
}

// Parse accepts a ParseRequest and returns a structured ParsedResume.
func (rp *ResumeParser) Parse(req schema.ParseRequest) (*schema.ParsedResume, error) {
	if len(req.FileContent) == 0 {
		return nil, &schema.ParseError{
			Code:    "EMPTY_REQUEST",
			Message: "file content is empty",
		}
	}

	fileType := normalizeFileType(req.FileType, req.FileName)
	if fileType == "" {
		return nil, &schema.ParseError{
			Code:    "UNSUPPORTED_FORMAT",
			Message: fmt.Sprintf("unsupported file type: %q (supported: pdf, docx)", req.FileType),
		}
	}

	// Step 1: Extract raw text from document
	rawText, err := rp.extractText(req.FileContent, fileType)
	if err != nil {
		return nil, err
	}

	// Step 2: Split into sections
	sections := extractor.SplitSections(rawText)

	// Step 3: Extract fields from sections
	result := &schema.ParsedResume{
		ParsedAt:      time.Now().UTC(),
		SourceFile:    req.FileName,
		FileType:      fileType,
		ParserVersion: parserVersion,
		SectionsFound: extractor.ListFoundSections(sections),
	}

	if req.IncludeRaw {
		result.RawText = rawText
	}

	// Personal info from full text (contact info can be anywhere in header)
	result.PersonalInfo = extractor.ExtractPersonalInfo(rawText)

	// Summary
	summaryText := extractor.GetSectionText(sections, extractor.SectionSummary)
	result.Summary = strings.TrimSpace(summaryText)

	// Work experience
	expText := extractor.GetSectionText(sections, extractor.SectionExperience)
	result.WorkExperience = extractor.ExtractWorkExperience(expText)

	// Education
	eduText := extractor.GetSectionText(sections, extractor.SectionEducation)
	result.Education = extractor.ExtractEducation(eduText)

	// Skills
	skillsText := extractor.GetSectionText(sections, extractor.SectionSkills)
	result.Skills = extractor.ExtractSkills(skillsText)

	// Certifications
	certText := extractor.GetSectionText(sections, extractor.SectionCertifications)
	result.Certifications = extractor.ExtractCertifications(certText)

	// Projects
	projText := extractor.GetSectionText(sections, extractor.SectionProjects)
	result.Projects = extractor.ExtractProjects(projText)

	// Compute overall confidence
	result.OverallConfidence = computeOverallConfidence(result)

	// Add warnings for missing sections
	result.Warnings = generateWarnings(result)

	return result, nil
}

// extractText dispatches to the appropriate parser based on file type.
func (rp *ResumeParser) extractText(data []byte, fileType string) (string, error) {
	switch fileType {
	case "pdf":
		text, err := rp.pdfParser.ExtractText(data)
		if err != nil {
			return "", convertParserError(err)
		}
		return text, nil
	case "docx":
		text, err := rp.docxParser.ExtractText(data)
		if err != nil {
			return "", convertParserError(err)
		}
		return text, nil
	default:
		return "", &schema.ParseError{
			Code:    "UNSUPPORTED_FORMAT",
			Message: fmt.Sprintf("unsupported file type: %s", fileType),
		}
	}
}

// normalizeFileType returns a canonical file type string from the provided type or filename.
func normalizeFileType(fileType, fileName string) string {
	ft := strings.ToLower(strings.TrimSpace(fileType))
	switch ft {
	case "pdf", "application/pdf":
		return "pdf"
	case "docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "docx"
	}

	// Infer from filename extension
	lower := strings.ToLower(fileName)
	if strings.HasSuffix(lower, ".pdf") {
		return "pdf"
	}
	if strings.HasSuffix(lower, ".docx") {
		return "docx"
	}

	return ""
}

// convertParserError converts internal parser errors to schema.ParseError.
func convertParserError(err error) *schema.ParseError {
	if pe, ok := err.(*ParseError); ok {
		return &schema.ParseError{
			Code:    pe.Code,
			Message: pe.Message,
		}
	}
	return &schema.ParseError{
		Code:    "PARSE_ERROR",
		Message: err.Error(),
	}
}

// computeOverallConfidence calculates a weighted average confidence score.
func computeOverallConfidence(r *schema.ParsedResume) schema.ConfidenceScore {
	scores := []float64{}

	if r.PersonalInfo.Name != "" || r.PersonalInfo.Email != "" {
		scores = append(scores, float64(r.PersonalInfo.Confidence))
	}

	if len(r.WorkExperience) > 0 {
		total := 0.0
		for _, e := range r.WorkExperience {
			total += float64(e.Confidence)
		}
		scores = append(scores, total/float64(len(r.WorkExperience)))
	}

	if len(r.Education) > 0 {
		total := 0.0
		for _, e := range r.Education {
			total += float64(e.Confidence)
		}
		scores = append(scores, total/float64(len(r.Education)))
	}

	if len(r.Skills) > 0 {
		scores = append(scores, 0.85) // Skills section is generally reliable
	}

	if len(scores) == 0 {
		return 0
	}

	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	return schema.ConfidenceScore(sum / float64(len(scores)))
}

// generateWarnings produces warnings for missing or low-confidence sections.
func generateWarnings(r *schema.ParsedResume) []string {
	var warnings []string

	if r.PersonalInfo.Name == "" {
		warnings = append(warnings, "could not extract candidate name")
	}
	if r.PersonalInfo.Email == "" {
		warnings = append(warnings, "could not extract email address")
	}
	if len(r.WorkExperience) == 0 {
		warnings = append(warnings, "no work experience section found")
	}
	if len(r.Education) == 0 {
		warnings = append(warnings, "no education section found")
	}
	if len(r.Skills) == 0 {
		warnings = append(warnings, "no skills section found")
	}
	if r.OverallConfidence < 0.5 {
		warnings = append(warnings, "low overall confidence - resume format may be non-standard")
	}

	return warnings
}
