package extractor

import (
	"testing"
)

func TestSplitSections_BasicResume(t *testing.T) {
	text := `John Smith
john@example.com

EXPERIENCE
Software Engineer at Acme Corp
Jan 2020 - Present
• Built microservices

EDUCATION
Bachelor of Science in Computer Science
MIT
2016 - 2020

SKILLS
Go, Python, Docker, Kubernetes`

	sections := SplitSections(text)
	found := ListFoundSections(sections)

	wantSections := map[string]bool{
		"experience": false,
		"education":  false,
		"skills":     false,
	}

	for _, s := range found {
		if _, ok := wantSections[s]; ok {
			wantSections[s] = true
		}
	}

	for section, found := range wantSections {
		if !found {
			t.Errorf("expected section %q to be found", section)
		}
	}
}

func TestSplitSections_AllCapsHeaders(t *testing.T) {
	text := `WORK EXPERIENCE
Engineer at Company
2020 - 2022

EDUCATION
BS Computer Science
2016 - 2020`

	sections := SplitSections(text)
	found := ListFoundSections(sections)

	hasExp := false
	hasEdu := false
	for _, s := range found {
		if s == "experience" {
			hasExp = true
		}
		if s == "education" {
			hasEdu = true
		}
	}

	if !hasExp {
		t.Error("expected 'experience' section to be found with all-caps header")
	}
	if !hasEdu {
		t.Error("expected 'education' section to be found with all-caps header")
	}
}

func TestGetSectionText(t *testing.T) {
	text := `SKILLS
Go, Python, Docker

EXPERIENCE
Engineer at Company`

	sections := SplitSections(text)
	skillsText := GetSectionText(sections, SectionSkills)

	if skillsText == "" {
		t.Error("expected skills section text to be non-empty")
	}
}

func TestDetectSectionHeader(t *testing.T) {
	tests := []struct {
		line     string
		wantType SectionType
		wantOK   bool
	}{
		{"EXPERIENCE", SectionExperience, true},
		{"Work Experience", SectionExperience, true},
		{"EDUCATION", SectionEducation, true},
		{"Skills", SectionSkills, true},
		{"Certifications", SectionCertifications, true},
		{"Projects", SectionProjects, true},
		{"Summary", SectionSummary, true},
		{"John Smith", SectionUnknown, false},
		{"Software Engineer at Google", SectionUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			st, _, ok := detectSectionHeader(tt.line)
			if ok != tt.wantOK {
				t.Errorf("detectSectionHeader(%q) ok = %v, want %v", tt.line, ok, tt.wantOK)
			}
			if ok && st != tt.wantType {
				t.Errorf("detectSectionHeader(%q) type = %v, want %v", tt.line, st, tt.wantType)
			}
		})
	}
}
