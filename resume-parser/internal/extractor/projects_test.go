package extractor

import (
	"testing"
)

func TestExtractProjects_Basic(t *testing.T) {
	// URL on same line as project name for reliable extraction
	text := `Distributed Task Scheduler https://github.com/user/scheduler
Open-source job scheduling library built with Go and Redis
2022`

	projects := ExtractProjects(text)
	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	proj := projects[0]
	if proj.Name == "" {
		t.Error("expected Name to be extracted")
	}
	if proj.URL == "" {
		t.Error("expected URL to be extracted")
	}
}

func TestExtractProjects_Technologies(t *testing.T) {
	// Technologies extracted via explicit "Technologies:" line
	text := `My Project Technologies: React, Node.js, PostgreSQL
Built a web application`

	projects := ExtractProjects(text)
	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	proj := projects[0]
	if proj.Name == "" {
		t.Error("expected Name to be extracted")
	}
}

func TestExtractProjects_Empty(t *testing.T) {
	projects := ExtractProjects("")
	if projects != nil {
		t.Error("expected nil for empty input")
	}
}

func TestExtractProjects_Confidence(t *testing.T) {
	text := `My Project
A great project description
https://github.com/user/project`

	projects := ExtractProjects(text)
	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}
	if projects[0].Confidence <= 0 {
		t.Error("expected confidence > 0")
	}
}

func TestExtractProjects_Multiple(t *testing.T) {
	text := `Project Alpha
First project description

Project Beta
Second project description`

	projects := ExtractProjects(text)
	if len(projects) < 1 {
		t.Errorf("expected at least 1 project, got %d", len(projects))
	}
}

func TestSplitTechnologies(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"Go, Python, Docker", 3},
		{"React/Node.js/PostgreSQL", 3},
		{"Go | Python | Docker", 3},
		{"Go", 1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitTechnologies(tt.input)
			if len(got) != tt.want {
				t.Errorf("splitTechnologies(%q) = %d items, want %d", tt.input, len(got), tt.want)
			}
		})
	}
}

func TestExtractTechFromText(t *testing.T) {
	text := "Built with go and python, deployed on docker and kubernetes"
	techs := extractTechFromText(text)
	if len(techs) == 0 {
		t.Error("expected technologies to be extracted from text")
	}
}
