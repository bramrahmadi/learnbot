package extractor

import (
	"testing"
)

func TestExtractSkills_TechnicalSkills(t *testing.T) {
	text := `Go, Python, JavaScript, Docker, Kubernetes, PostgreSQL, Redis`

	skills := ExtractSkills(text)
	if len(skills) == 0 {
		t.Fatal("expected skills to be extracted")
	}

	// Check that known technical skills are classified correctly
	techCount := 0
	for _, s := range skills {
		if s.Category == "technical" {
			techCount++
		}
	}
	if techCount == 0 {
		t.Error("expected at least one technical skill")
	}
}

func TestExtractSkills_BulletList(t *testing.T) {
	text := `• Go
• Python
• Docker
• Kubernetes`

	skills := ExtractSkills(text)
	if len(skills) < 3 {
		t.Errorf("expected at least 3 skills, got %d", len(skills))
	}
}

func TestExtractSkills_SoftSkills(t *testing.T) {
	text := `Leadership, Communication, Teamwork, Problem Solving`

	skills := ExtractSkills(text)
	softCount := 0
	for _, s := range skills {
		if s.Category == "soft" {
			softCount++
		}
	}
	if softCount == 0 {
		t.Error("expected at least one soft skill")
	}
}

func TestExtractSkills_Empty(t *testing.T) {
	skills := ExtractSkills("")
	if skills != nil {
		t.Error("expected nil for empty input")
	}
}

func TestExtractSkills_NoDuplicates(t *testing.T) {
	text := `Go, Go, Python, Python, Docker`

	skills := ExtractSkills(text)
	seen := map[string]int{}
	for _, s := range skills {
		seen[s.Name]++
	}
	for name, count := range seen {
		if count > 1 {
			t.Errorf("skill %q appears %d times, expected 1", name, count)
		}
	}
}

func TestExtractSkills_Confidence(t *testing.T) {
	text := `Go, Python, Docker`

	skills := ExtractSkills(text)
	for _, s := range skills {
		if s.Confidence <= 0 {
			t.Errorf("skill %q has confidence <= 0", s.Name)
		}
	}
}

func TestClassifySkill(t *testing.T) {
	tests := []struct {
		skill    string
		wantCat  string
	}{
		{"go", "technical"},
		{"python", "technical"},
		{"docker", "technical"},
		{"leadership", "soft"},
		{"communication", "soft"},
		{"unknownskill123", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.skill, func(t *testing.T) {
			got := classifySkill(tt.skill)
			if got != tt.wantCat {
				t.Errorf("classifySkill(%q) = %q, want %q", tt.skill, got, tt.wantCat)
			}
		})
	}
}
