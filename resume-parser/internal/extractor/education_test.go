package extractor

import (
	"testing"
)

func TestExtractEducation_Basic(t *testing.T) {
	text := `Massachusetts Institute of Technology
Bachelor of Science in Computer Science
2016 - 2020
GPA: 3.8`

	educations := ExtractEducation(text)
	if len(educations) == 0 {
		t.Fatal("expected at least one education entry")
	}

	edu := educations[0]
	if edu.Institution == "" {
		t.Error("expected Institution to be extracted")
	}
	if edu.GPA != "3.8" {
		t.Errorf("expected GPA = '3.8', got %q", edu.GPA)
	}
}

func TestExtractEducation_DegreeField(t *testing.T) {
	text := `Stanford University
Master of Science in Machine Learning
2020 - 2022`

	educations := ExtractEducation(text)
	if len(educations) == 0 {
		t.Fatal("expected at least one education entry")
	}

	edu := educations[0]
	if edu.Degree == "" {
		t.Error("expected Degree to be extracted")
	}
}

func TestExtractEducation_Honors(t *testing.T) {
	text := `Harvard University
Bachelor of Arts in Economics
2015 - 2019
Graduated Magna Cum Laude`

	educations := ExtractEducation(text)
	if len(educations) == 0 {
		t.Fatal("expected at least one education entry")
	}

	edu := educations[0]
	if edu.Honors == "" {
		t.Error("expected Honors to be extracted")
	}
}

func TestExtractEducation_Empty(t *testing.T) {
	educations := ExtractEducation("")
	if educations != nil {
		t.Error("expected nil for empty input")
	}
}

func TestExtractEducation_Confidence(t *testing.T) {
	text := `MIT
B.S. Computer Science
2016 - 2020`

	educations := ExtractEducation(text)
	if len(educations) == 0 {
		t.Fatal("expected at least one education entry")
	}

	edu := educations[0]
	if edu.Confidence <= 0 {
		t.Error("expected confidence > 0")
	}
}

func TestParseDegreeField(t *testing.T) {
	tests := []struct {
		input     string
		wantDeg   string
		wantField string
	}{
		{
			input:     "Bachelor of Science in Computer Science",
			wantDeg:   "Bachelor of Science",
			wantField: "Computer Science",
		},
		{
			input:     "B.S., Computer Science",
			wantDeg:   "B.S.",
			wantField: "Computer Science",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			deg, field := parseDegreeField(tt.input)
			if deg != tt.wantDeg {
				t.Errorf("degree = %q, want %q", deg, tt.wantDeg)
			}
			if field != tt.wantField {
				t.Errorf("field = %q, want %q", field, tt.wantField)
			}
		})
	}
}
