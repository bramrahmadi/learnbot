package extractor

import (
	"testing"
)

func TestExtractWorkExperience_Basic(t *testing.T) {
	text := `Software Engineer | Acme Corp
Jan 2020 - Present
• Designed and implemented microservices using Go
• Reduced API latency by 40% through caching
• Led a team of 5 engineers

Backend Developer - TechStartup
March 2018 - December 2019
• Built REST APIs with Python/Django
• Managed PostgreSQL databases`

	experiences := ExtractWorkExperience(text)

	if len(experiences) == 0 {
		t.Fatal("expected at least one work experience entry")
	}

	// Check first entry
	exp := experiences[0]
	if exp.StartDate == "" {
		t.Error("expected StartDate to be extracted")
	}
	if len(exp.Responsibilities) == 0 {
		t.Error("expected responsibilities to be extracted")
	}
	if exp.Confidence <= 0 {
		t.Error("expected confidence > 0")
	}
}

func TestExtractWorkExperience_CurrentPosition(t *testing.T) {
	text := `Senior Engineer at BigCorp
June 2021 - Present
• Leading platform team`

	experiences := ExtractWorkExperience(text)
	if len(experiences) == 0 {
		t.Fatal("expected at least one work experience entry")
	}

	exp := experiences[0]
	if !exp.IsCurrent {
		t.Error("expected IsCurrent = true for 'Present' end date")
	}
	if exp.EndDate != "Present" {
		t.Errorf("expected EndDate = 'Present', got %q", exp.EndDate)
	}
}

func TestExtractWorkExperience_Empty(t *testing.T) {
	experiences := ExtractWorkExperience("")
	if experiences != nil {
		t.Error("expected nil for empty input")
	}
}

func TestExtractWorkExperience_MultipleEntries(t *testing.T) {
	text := `Software Engineer | Company A
Jan 2022 - Present
• Built features

Junior Developer | Company B
Jan 2020 - Dec 2021
• Wrote tests`

	experiences := ExtractWorkExperience(text)
	if len(experiences) < 1 {
		t.Errorf("expected at least 1 experience entry, got %d", len(experiences))
	}
}

func TestIsCurrentDate(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Present", true},
		{"present", true},
		{"Current", true},
		{"current", true},
		{"Now", true},
		{"2023", false},
		{"December 2022", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isCurrentDate(tt.input)
			if got != tt.want {
				t.Errorf("isCurrentDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractBullets(t *testing.T) {
	text := `• First bullet point
• Second bullet point
- Third with dash
* Fourth with asterisk`

	bullets := extractBullets(text)
	if len(bullets) < 2 {
		t.Errorf("expected at least 2 bullets, got %d", len(bullets))
	}
}
