package parser

import (
	"fmt"
	"testing"

	"github.com/learnbot/resume-parser/internal/schema"
)

// TestResumeParser_FileTypeFromName tests file type detection from filename.
func TestResumeParser_FileTypeFromName(t *testing.T) {
	rp := NewResumeParser()

	// Test with docx extension but no explicit type
	docxData := buildMinimalDOCX("Jane Doe\njane@example.com\n\nSKILLS\nGo, Python")
	result, err := rp.Parse(schema.ParseRequest{
		FileName:    "my_resume.docx",
		FileContent: docxData,
		// FileType intentionally empty - should be inferred from filename
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FileType != "docx" {
		t.Errorf("expected file type 'docx', got %q", result.FileType)
	}
}

// TestResumeParser_ParsedAt tests that ParsedAt is set.
func TestResumeParser_ParsedAt(t *testing.T) {
	rp := NewResumeParser()
	docxData := buildMinimalDOCX("Jane Doe\njane@example.com")

	result, err := rp.Parse(schema.ParseRequest{
		FileName:    "resume.docx",
		FileContent: docxData,
		FileType:    "docx",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ParsedAt.IsZero() {
		t.Error("expected ParsedAt to be set")
	}
}

// TestResumeParser_SourceFile tests that SourceFile is set correctly.
func TestResumeParser_SourceFile(t *testing.T) {
	rp := NewResumeParser()
	docxData := buildMinimalDOCX("Jane Doe\njane@example.com")

	result, err := rp.Parse(schema.ParseRequest{
		FileName:    "my_resume.docx",
		FileContent: docxData,
		FileType:    "docx",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SourceFile != "my_resume.docx" {
		t.Errorf("expected SourceFile = 'my_resume.docx', got %q", result.SourceFile)
	}
}

// TestResumeParser_Warnings tests that warnings are generated for missing sections.
func TestResumeParser_Warnings(t *testing.T) {
	rp := NewResumeParser()
	// Minimal resume with no sections
	docxData := buildMinimalDOCX("Some random text without any structure")

	result, err := rp.Parse(schema.ParseRequest{
		FileName:    "minimal.docx",
		FileContent: docxData,
		FileType:    "docx",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should have warnings about missing sections
	if len(result.Warnings) == 0 {
		t.Error("expected warnings for minimal resume")
	}
}

// TestResumeParser_FullResume tests parsing a comprehensive resume.
func TestResumeParser_FullResume(t *testing.T) {
	rp := NewResumeParser()

	fullResume := `Jane Doe
jane.doe@email.com
(555) 987-6543
San Francisco, CA
linkedin.com/in/janedoe

SUMMARY
Experienced software engineer with 6+ years building scalable backend systems.

EXPERIENCE
Senior Software Engineer | TechCorp Inc.
March 2021 - Present
Built distributed job scheduling system handling 1M+ tasks/day
Reduced infrastructure costs by 35%

Software Engineer | StartupXYZ
January 2019 - February 2021
Built RESTful APIs serving 500K daily active users using Go and PostgreSQL

EDUCATION
University of California, Berkeley
Bachelor of Science in Computer Science
2013 - 2017
GPA: 3.7

SKILLS
Go, Python, JavaScript, Docker, Kubernetes, PostgreSQL, Redis, AWS

CERTIFICATIONS
AWS Certified Solutions Architect - Associate
Amazon Web Services
2022

PROJECTS
Distributed Task Scheduler https://github.com/janedoe/scheduler
Open-source job scheduling library`

	docxData := buildMinimalDOCX(fullResume)
	result, err := rp.Parse(schema.ParseRequest{
		FileName:   "full_resume.docx",
		FileContent: docxData,
		FileType:   "docx",
		IncludeRaw: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate all sections
	if result.PersonalInfo.Email == "" {
		t.Error("expected email")
	}
	if result.PersonalInfo.Phone == "" {
		t.Error("expected phone")
	}
	if len(result.WorkExperience) == 0 {
		t.Error("expected work experience")
	}
	if len(result.Education) == 0 {
		t.Error("expected education")
	}
	if len(result.Skills) == 0 {
		t.Error("expected skills")
	}
	if len(result.Certifications) == 0 {
		t.Error("expected certifications")
	}
	if len(result.Projects) == 0 {
		t.Error("expected projects")
	}
	if result.Summary == "" {
		t.Error("expected summary")
	}
	if result.RawText == "" {
		t.Error("expected raw text")
	}
	if result.OverallConfidence <= 0 {
		t.Error("expected overall confidence > 0")
	}
	if len(result.SectionsFound) == 0 {
		t.Error("expected sections found")
	}
}

// TestConvertParserError tests error conversion.
func TestConvertParserError(t *testing.T) {
	pe := &ParseError{Code: "TEST", Message: "test message"}
	schemaErr := convertParserError(pe)
	if schemaErr.Code != "TEST" {
		t.Errorf("expected code 'TEST', got %q", schemaErr.Code)
	}

	// Test with generic error
	genericErr := fmt.Errorf("generic error")
	schemaErr2 := convertParserError(genericErr)
	if schemaErr2.Code != "PARSE_ERROR" {
		t.Errorf("expected code 'PARSE_ERROR', got %q", schemaErr2.Code)
	}
}
