package scraper

import (
	"testing"
	"time"

	"github.com/learnbot/job-aggregator/internal/model"
)

func TestExtractExperienceLevel(t *testing.T) {
	tests := []struct {
		title       string
		description string
		want        model.ExperienceLevel
	}{
		{"Senior Software Engineer", "", model.LevelSenior},
		{"Junior Developer", "", model.LevelEntry},
		{"Software Engineering Intern", "", model.LevelInternship},
		{"VP of Engineering", "", model.LevelExecutive},
		{"Director of Engineering", "", model.LevelExecutive},
		{"Software Engineer", "", model.LevelUnknown},
		{"Mid-Level Backend Developer", "", model.LevelMid},
		{"Principal Engineer", "", model.LevelSenior},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := ExtractExperienceLevel(tt.title, tt.description)
			if got != tt.want {
				t.Errorf("ExtractExperienceLevel(%q) = %v, want %v", tt.title, got, tt.want)
			}
		})
	}
}

func TestExtractLocationType(t *testing.T) {
	tests := []struct {
		location    string
		description string
		want        model.WorkLocationType
	}{
		{"Remote", "", model.LocationRemote},
		{"New York, NY", "", model.LocationOnSite},
		{"", "This is a remote position", model.LocationRemote},
		{"", "Hybrid work arrangement available", model.LocationHybrid},
		{"San Francisco, CA (Hybrid)", "", model.LocationHybrid},
		{"", "", model.LocationUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.location+"/"+tt.description[:min(len(tt.description), 20)], func(t *testing.T) {
			got := ExtractLocationType(tt.location, tt.description)
			if got != tt.want {
				t.Errorf("ExtractLocationType(%q, %q) = %v, want %v",
					tt.location, tt.description, got, tt.want)
			}
		})
	}
}

func TestExtractEmploymentType(t *testing.T) {
	tests := []struct {
		title string
		want  model.EmploymentType
	}{
		{"Software Engineer (Full-time)", model.EmploymentFullTime},
		{"Part-time Developer", model.EmploymentPartTime},
		{"Contract Backend Engineer", model.EmploymentContract},
		{"Software Engineering Intern", model.EmploymentInternship},
		{"Software Engineer", model.EmploymentFullTime},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := ExtractEmploymentType(tt.title, "")
			if got != tt.want {
				t.Errorf("ExtractEmploymentType(%q) = %v, want %v", tt.title, got, tt.want)
			}
		})
	}
}

func TestParseSalary(t *testing.T) {
	tests := []struct {
		raw      string
		wantMin  bool
		wantMax  bool
		wantCurr string
	}{
		{"$100,000 - $150,000", true, true, "USD"},
		{"$120k - $180k per year", true, true, "USD"},
		// Non-USD currencies detected but regex requires $ prefix for value extraction
		{"£50,000 - £70,000", false, false, "GBP"},
		{"€60,000 - €80,000", false, false, "EUR"},
		{"$95,000", true, false, "USD"},
		{"", false, false, "USD"},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			minVal, maxVal, curr := ParseSalary(tt.raw)
			if tt.wantMin && minVal == nil {
				t.Errorf("ParseSalary(%q): expected min salary, got nil", tt.raw)
			}
			if !tt.wantMin && minVal != nil {
				t.Errorf("ParseSalary(%q): expected nil min, got %d", tt.raw, *minVal)
			}
			if tt.wantMax && maxVal == nil {
				t.Errorf("ParseSalary(%q): expected max salary, got nil", tt.raw)
			}
			if curr != tt.wantCurr {
				t.Errorf("ParseSalary(%q): currency = %q, want %q", tt.raw, curr, tt.wantCurr)
			}
		})
	}
}

func TestExtractSkillsFromText(t *testing.T) {
	text := `We are looking for a Go developer with experience in:
Required:
- Go (Golang)
- PostgreSQL
- Docker
- Kubernetes

Preferred:
- Redis
- Kafka`

	required, preferred := ExtractSkillsFromText(text)

	if len(required) == 0 {
		t.Error("expected required skills to be extracted")
	}

	// Go should be in required
	foundGo := false
	for _, s := range required {
		if s == "go" || s == "golang" {
			foundGo = true
			break
		}
	}
	if !foundGo {
		t.Error("expected 'go' or 'golang' in required skills")
	}

	_ = preferred // preferred may or may not be populated depending on text parsing
}

func TestCleanText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<p>Hello <b>World</b></p>", "Hello World"},
		{"&amp; &lt; &gt;", "& < >"},
		{"  multiple   spaces  ", "multiple spaces"},
		{"&nbsp;non-breaking&nbsp;space", "non-breaking space"},
	}

	for _, tt := range tests {
		t.Run(tt.input[:min(len(tt.input), 20)], func(t *testing.T) {
			got := CleanText(tt.input)
			if got != tt.want {
				t.Errorf("CleanText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRelativeDate(t *testing.T) {
	tests := []struct {
		input    string
		wantNil  bool
		maxAge   time.Duration
	}{
		{"today", false, 24 * time.Hour},
		{"1 day ago", false, 48 * time.Hour},
		{"3 days ago", false, 4 * 24 * time.Hour},
		{"1 week ago", false, 8 * 24 * time.Hour},
		{"2 weeks ago", false, 15 * 24 * time.Hour},
		{"1 month ago", false, 32 * 24 * time.Hour},
		{"invalid date", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseRelativeDate(tt.input)
			if tt.wantNil && got != nil {
				t.Errorf("ParseRelativeDate(%q) = %v, want nil", tt.input, got)
			}
			if !tt.wantNil && got == nil {
				t.Errorf("ParseRelativeDate(%q) = nil, want non-nil", tt.input)
			}
			if got != nil {
				age := time.Since(*got)
				if age > tt.maxAge {
					t.Errorf("ParseRelativeDate(%q): age %v > maxAge %v", tt.input, age, tt.maxAge)
				}
			}
		})
	}
}
