package extractor

import (
	"testing"
)

func TestExtractPersonalInfo_Email(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantEmail string
	}{
		{
			name:      "simple email",
			text:      "John Doe\njohn.doe@example.com\n555-123-4567",
			wantEmail: "john.doe@example.com",
		},
		{
			name:      "email with plus",
			text:      "Contact: jane+work@company.org",
			wantEmail: "jane+work@company.org",
		},
		{
			name:      "no email",
			text:      "John Doe\nSoftware Engineer",
			wantEmail: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ExtractPersonalInfo(tt.text)
			if info.Email != tt.wantEmail {
				t.Errorf("Email = %q, want %q", info.Email, tt.wantEmail)
			}
		})
	}
}

func TestExtractPersonalInfo_Phone(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantPhone string
	}{
		{
			name:      "US phone with dashes",
			text:      "Phone: 555-123-4567",
			wantPhone: "555-123-4567",
		},
		{
			name:      "US phone with dots",
			text:      "555.123.4567",
			wantPhone: "555.123.4567",
		},
		{
			name:      "US phone with parens",
			text:      "(555) 123-4567",
			wantPhone: "(555) 123-4567",
		},
		{
			name:      "no phone",
			text:      "John Doe\nSoftware Engineer",
			wantPhone: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ExtractPersonalInfo(tt.text)
			if info.Phone != tt.wantPhone {
				t.Errorf("Phone = %q, want %q", info.Phone, tt.wantPhone)
			}
		})
	}
}

func TestExtractPersonalInfo_Name(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantName string
	}{
		{
			name:     "simple name first line",
			text:     "John Smith\njohn@example.com\n555-123-4567",
			wantName: "John Smith",
		},
		{
			name:     "three-part name",
			text:     "Mary Jane Watson\nmary@example.com",
			wantName: "Mary Jane Watson",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ExtractPersonalInfo(tt.text)
			if info.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", info.Name, tt.wantName)
			}
		})
	}
}

func TestExtractPersonalInfo_LinkedIn(t *testing.T) {
	text := "John Doe\nlinkedin.com/in/johndoe\ngithub.com/johndoe"
	info := ExtractPersonalInfo(text)
	if info.LinkedIn != "https://linkedin.com/in/johndoe" {
		t.Errorf("LinkedIn = %q, want %q", info.LinkedIn, "https://linkedin.com/in/johndoe")
	}
	if info.GitHub != "https://github.com/johndoe" {
		t.Errorf("GitHub = %q, want %q", info.GitHub, "https://github.com/johndoe")
	}
}

func TestExtractPersonalInfo_Location(t *testing.T) {
	text := "John Doe\nSan Francisco, CA\njohn@example.com"
	info := ExtractPersonalInfo(text)
	if info.Location == "" {
		t.Error("expected location to be extracted, got empty string")
	}
}

func TestExtractPersonalInfo_Confidence(t *testing.T) {
	// Full contact info should yield high confidence
	text := "John Smith\njohn.smith@example.com\n555-123-4567\nNew York, NY"
	info := ExtractPersonalInfo(text)
	if info.Confidence < 0.7 {
		t.Errorf("expected confidence >= 0.7, got %v", info.Confidence)
	}

	// No contact info should yield zero confidence
	empty := ExtractPersonalInfo("")
	if empty.Confidence != 0 {
		t.Errorf("expected confidence = 0 for empty text, got %v", empty.Confidence)
	}
}
