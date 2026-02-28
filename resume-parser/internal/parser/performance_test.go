package parser

import (
	"strings"
	"testing"

	"github.com/learnbot/resume-parser/internal/schema"
)

// ─────────────────────────────────────────────────────────────────────────────
// Performance benchmarks for resume parser
// ─────────────────────────────────────────────────────────────────────────────

// buildLargeResumeText creates a large resume text for performance testing.
func buildLargeResumeText() string {
	var sb strings.Builder

	sb.WriteString("John Doe\njohn.doe@example.com\n+1-555-123-4567\nSan Francisco, CA\n\n")

	sb.WriteString("SUMMARY\n")
	sb.WriteString("Experienced software engineer with 10+ years of experience building scalable systems. ")
	sb.WriteString("Expert in Go, Python, and distributed systems. Passionate about clean code and performance.\n\n")

	sb.WriteString("EXPERIENCE\n")
	for i := 0; i < 5; i++ {
		sb.WriteString("Senior Software Engineer | TechCorp Inc | 2020 - Present\n")
		sb.WriteString("- Designed and implemented microservices architecture serving 10M+ users\n")
		sb.WriteString("- Led team of 8 engineers to deliver critical platform features\n")
		sb.WriteString("- Reduced API latency by 40% through query optimization\n")
		sb.WriteString("- Implemented CI/CD pipelines using GitHub Actions and Kubernetes\n\n")
	}

	sb.WriteString("EDUCATION\n")
	sb.WriteString("Master of Science in Computer Science | Stanford University | 2014\n")
	sb.WriteString("Bachelor of Science in Computer Science | UC Berkeley | 2012\n\n")

	sb.WriteString("SKILLS\n")
	sb.WriteString("Go, Python, Java, JavaScript, TypeScript, Rust, C++\n")
	sb.WriteString("React, Node.js, Django, Flask, Spring Boot\n")
	sb.WriteString("PostgreSQL, MySQL, MongoDB, Redis, Elasticsearch\n")
	sb.WriteString("AWS, GCP, Azure, Docker, Kubernetes, Terraform\n")
	sb.WriteString("Git, GitHub, GitLab, Jira, Confluence\n")
	sb.WriteString("Machine Learning, Deep Learning, NLP, TensorFlow, PyTorch\n\n")

	sb.WriteString("CERTIFICATIONS\n")
	sb.WriteString("AWS Certified Solutions Architect - Professional | Amazon Web Services | 2022\n")
	sb.WriteString("Google Cloud Professional Data Engineer | Google | 2021\n")
	sb.WriteString("Certified Kubernetes Administrator (CKA) | CNCF | 2020\n\n")

	sb.WriteString("PROJECTS\n")
	sb.WriteString("LearnBot Platform | 2023\n")
	sb.WriteString("AI-powered career development platform using RAG and LLMs\n")
	sb.WriteString("Technologies: Go, Python, PostgreSQL, Redis, Docker, Kubernetes\n\n")

	return sb.String()
}

// BenchmarkNormalizeFileType benchmarks the normalizeFileType function.
func BenchmarkNormalizeFileType(b *testing.B) {
	inputs := []struct{ fileType, fileName string }{
		{"pdf", "resume.pdf"},
		{"application/pdf", "resume.pdf"},
		{"docx", "resume.docx"},
		{"", "resume.pdf"},
		{"", "resume.docx"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, inp := range inputs {
			normalizeFileType(inp.fileType, inp.fileName)
		}
	}
}

// BenchmarkComputeOverallConfidence benchmarks confidence computation.
func BenchmarkComputeOverallConfidence(b *testing.B) {
	resume := &schema.ParsedResume{
		PersonalInfo: schema.PersonalInfo{
			Name:       "John Doe",
			Email:      "john@example.com",
			Confidence: schema.ConfidenceHigh,
		},
		WorkExperience: []schema.WorkExperience{
			{Title: "Engineer", Company: "TechCorp", Confidence: schema.ConfidenceHigh},
			{Title: "Developer", Company: "StartupCo", Confidence: schema.ConfidenceMedium},
		},
		Education: []schema.Education{
			{Degree: "BS Computer Science", Institution: "MIT", Confidence: schema.ConfidenceHigh},
		},
		Skills: []schema.Skill{
			{Name: "Go", Category: "technical"},
			{Name: "Python", Category: "technical"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		computeOverallConfidence(resume)
	}
}

// BenchmarkGenerateWarnings benchmarks warning generation.
func BenchmarkGenerateWarnings(b *testing.B) {
	resume := &schema.ParsedResume{
		PersonalInfo: schema.PersonalInfo{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		WorkExperience: []schema.WorkExperience{
			{Title: "Engineer"},
		},
		Education: []schema.Education{
			{Degree: "BS CS"},
		},
		Skills: []schema.Skill{
			{Name: "Go"},
		},
		OverallConfidence: 0.85,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateWarnings(resume)
	}
}

// TestNormalizeFileType_Performance tests that normalizeFileType handles many inputs quickly.
func TestNormalizeFileType_Performance(t *testing.T) {
	inputs := []struct {
		fileType string
		fileName string
		expected string
	}{
		{"pdf", "", "pdf"},
		{"application/pdf", "", "pdf"},
		{"PDF", "", "pdf"},
		{"docx", "", "docx"},
		{"application/vnd.openxmlformats-officedocument.wordprocessingml.document", "", "docx"},
		{"", "resume.pdf", "pdf"},
		{"", "resume.docx", "docx"},
		{"", "RESUME.PDF", "pdf"},
		{"", "RESUME.DOCX", "docx"},
		{"txt", "", ""},
		{"", "resume.txt", ""},
	}

	for _, inp := range inputs {
		result := normalizeFileType(inp.fileType, inp.fileName)
		if result != inp.expected {
			t.Errorf("normalizeFileType(%q, %q) = %q, want %q",
				inp.fileType, inp.fileName, result, inp.expected)
		}
	}
}

// TestComputeOverallConfidence_EmptyResume tests confidence for empty resume.
func TestComputeOverallConfidence_EmptyResume(t *testing.T) {
	resume := &schema.ParsedResume{}
	confidence := computeOverallConfidence(resume)
	if confidence != 0 {
		t.Errorf("expected 0 confidence for empty resume, got %v", confidence)
	}
}

// TestComputeOverallConfidence_FullResume tests confidence for complete resume.
func TestComputeOverallConfidence_FullResume(t *testing.T) {
	resume := &schema.ParsedResume{
		PersonalInfo: schema.PersonalInfo{
			Name:       "John Doe",
			Email:      "john@example.com",
			Confidence: schema.ConfidenceHigh,
		},
		WorkExperience: []schema.WorkExperience{
			{Title: "Engineer", Confidence: schema.ConfidenceHigh},
		},
		Education: []schema.Education{
			{Degree: "BS CS", Confidence: schema.ConfidenceHigh},
		},
		Skills: []schema.Skill{
			{Name: "Go"},
		},
	}

	confidence := computeOverallConfidence(resume)
	if confidence <= 0 {
		t.Errorf("expected positive confidence for full resume, got %v", confidence)
	}
	if confidence > 1 {
		t.Errorf("expected confidence <= 1, got %v", confidence)
	}
}

// TestGenerateWarnings_CompleteResume tests that no warnings for complete resume.
func TestGenerateWarnings_CompleteResume(t *testing.T) {
	resume := &schema.ParsedResume{
		PersonalInfo: schema.PersonalInfo{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		WorkExperience: []schema.WorkExperience{
			{Title: "Engineer"},
		},
		Education: []schema.Education{
			{Degree: "BS CS"},
		},
		Skills: []schema.Skill{
			{Name: "Go"},
		},
		OverallConfidence: 0.85,
	}

	warnings := generateWarnings(resume)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for complete resume, got: %v", warnings)
	}
}

// TestGenerateWarnings_EmptyResume tests that warnings are generated for empty resume.
func TestGenerateWarnings_EmptyResume(t *testing.T) {
	resume := &schema.ParsedResume{}
	warnings := generateWarnings(resume)

	expectedWarnings := []string{
		"could not extract candidate name",
		"could not extract email address",
		"no work experience section found",
		"no education section found",
		"no skills section found",
	}

	for _, expected := range expectedWarnings {
		found := false
		for _, w := range warnings {
			if w == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected warning %q not found in %v", expected, warnings)
		}
	}
}

// TestGenerateWarnings_LowConfidence tests warning for low confidence.
func TestGenerateWarnings_LowConfidence(t *testing.T) {
	resume := &schema.ParsedResume{
		PersonalInfo: schema.PersonalInfo{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		WorkExperience: []schema.WorkExperience{{Title: "Engineer"}},
		Education:      []schema.Education{{Degree: "BS CS"}},
		Skills:         []schema.Skill{{Name: "Go"}},
		OverallConfidence: 0.3, // Low confidence
	}

	warnings := generateWarnings(resume)
	found := false
	for _, w := range warnings {
		if w == "low overall confidence - resume format may be non-standard" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected low confidence warning, got: %v", warnings)
	}
}
