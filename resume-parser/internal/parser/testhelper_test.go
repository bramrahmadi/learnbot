package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"testing"

	"github.com/learnbot/resume-parser/internal/schema"
)

// buildMinimalDOCX creates a minimal valid DOCX file with the given text content.
func buildMinimalDOCX(text string) []byte {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// word/document.xml
	docXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
            xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>%s</w:body>
</w:document>`, textToParagraphs(text))

	f, _ := w.Create("word/document.xml")
	f.Write([]byte(docXML))

	// [Content_Types].xml (required for valid DOCX)
	ct, _ := w.Create("[Content_Types].xml")
	ct.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`))

	w.Close()
	return buf.Bytes()
}

// textToParagraphs converts plain text lines to Word XML paragraphs.
func textToParagraphs(text string) string {
	var sb bytes.Buffer
	lines := bytes.Split([]byte(text), []byte("\n"))
	for _, line := range lines {
		sb.WriteString(`<w:p><w:r><w:t>`)
		sb.Write(line)
		sb.WriteString(`</w:t></w:r></w:p>`)
	}
	return sb.String()
}

// sampleResumeText is a realistic resume for integration testing.
const sampleResumeText = `Jane Doe
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
Distributed Task Scheduler
Open-source job scheduling library built with Go and Redis
https://github.com/janedoe/scheduler`

// TestResumeParser_DOCXIntegration tests the full DOCX parsing pipeline.
func TestResumeParser_DOCXIntegration(t *testing.T) {
	rp := NewResumeParser()
	docxData := buildMinimalDOCX(sampleResumeText)

	result, err := rp.Parse(schema.ParseRequest{
		FileName:    "resume.docx",
		FileContent: docxData,
		FileType:    "docx",
		IncludeRaw:  true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate personal info
	if result.PersonalInfo.Email == "" {
		t.Error("expected email to be extracted from DOCX")
	}
	if result.PersonalInfo.Phone == "" {
		t.Error("expected phone to be extracted from DOCX")
	}

	// Validate sections found
	if len(result.SectionsFound) == 0 {
		t.Error("expected sections to be found")
	}

	// Validate skills
	if len(result.Skills) == 0 {
		t.Error("expected skills to be extracted")
	}

	// Validate raw text is included
	if result.RawText == "" {
		t.Error("expected raw text to be included when IncludeRaw=true")
	}

	// Validate parser metadata
	if result.ParserVersion == "" {
		t.Error("expected parser version to be set")
	}
	if result.FileType != "docx" {
		t.Errorf("expected file type 'docx', got %q", result.FileType)
	}
}

// TestResumeParser_DOCXIntegration_NoRaw tests that raw text is excluded by default.
func TestResumeParser_DOCXIntegration_NoRaw(t *testing.T) {
	rp := NewResumeParser()
	docxData := buildMinimalDOCX(sampleResumeText)

	result, err := rp.Parse(schema.ParseRequest{
		FileName:    "resume.docx",
		FileContent: docxData,
		FileType:    "docx",
		IncludeRaw:  false,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.RawText != "" {
		t.Error("expected raw text to be empty when IncludeRaw=false")
	}
}
