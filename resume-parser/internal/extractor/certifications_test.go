package extractor

import (
	"testing"
)

func TestExtractCertifications_Basic(t *testing.T) {
	text := `AWS Certified Solutions Architect - Associate
Amazon Web Services
2022
Credential ID: AWS-SAA-12345`

	certs := ExtractCertifications(text)
	if len(certs) == 0 {
		t.Fatal("expected at least one certification")
	}

	cert := certs[0]
	if cert.Name == "" {
		t.Error("expected Name to be extracted")
	}
	if cert.ID == "" {
		t.Error("expected credential ID to be extracted")
	}
}

func TestExtractCertifications_Multiple(t *testing.T) {
	text := `AWS Certified Solutions Architect
Amazon
2022

Certified Kubernetes Administrator
Linux Foundation
2021`

	certs := ExtractCertifications(text)
	if len(certs) < 2 {
		t.Errorf("expected at least 2 certifications, got %d", len(certs))
	}
}

func TestExtractCertifications_Empty(t *testing.T) {
	certs := ExtractCertifications("")
	if certs != nil {
		t.Error("expected nil for empty input")
	}
}

func TestExtractCertifications_Confidence(t *testing.T) {
	text := `AWS Certified Developer
Amazon Web Services
2023`

	certs := ExtractCertifications(text)
	if len(certs) == 0 {
		t.Fatal("expected at least one certification")
	}
	if certs[0].Confidence <= 0 {
		t.Error("expected confidence > 0")
	}
}

func TestExtractCertifications_BulletList(t *testing.T) {
	text := `• AWS Certified Solutions Architect
• Google Cloud Professional
• Certified Kubernetes Administrator`

	certs := ExtractCertifications(text)
	if len(certs) == 0 {
		t.Error("expected certifications to be extracted from bullet list")
	}
}

func TestExtractIssuer(t *testing.T) {
	tests := []struct {
		input    string
		wantNonEmpty bool
	}{
		{"Issued by Amazon Web Services", true},
		{"Google Cloud certification", true},
		{"Microsoft Azure exam", true},
		{"Unknown issuer xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractIssuer(tt.input)
			if tt.wantNonEmpty && got == "" {
				t.Errorf("extractIssuer(%q) = empty, want non-empty", tt.input)
			}
			if !tt.wantNonEmpty && got != "" {
				t.Errorf("extractIssuer(%q) = %q, want empty", tt.input, got)
			}
		})
	}
}

func TestIsBulletLine(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"• bullet", true},
		{"- dash bullet", true},
		{"* star bullet", true},
		{"> arrow bullet", true},
		{"normal line", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isBulletLine(tt.input)
			if got != tt.want {
				t.Errorf("isBulletLine(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
