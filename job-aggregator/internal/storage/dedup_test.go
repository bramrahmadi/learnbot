package storage

import (
	"testing"

	"github.com/learnbot/job-aggregator/internal/model"
)

func TestComputeDedupHash_WithExternalID(t *testing.T) {
	job := &model.ScrapedJob{
		Source:     model.SourceLinkedIn,
		ExternalID: "12345",
		Title:      "Software Engineer",
		CompanyName: "Acme Corp",
	}

	hash1 := ComputeDedupHash(job)
	hash2 := ComputeDedupHash(job)

	if hash1 != hash2 {
		t.Error("same job should produce same hash")
	}
	if hash1 == "" {
		t.Error("hash should not be empty")
	}
}

func TestComputeDedupHash_WithoutExternalID(t *testing.T) {
	job := &model.ScrapedJob{
		Source:      model.SourceIndeed,
		Title:       "Backend Developer",
		CompanyName: "TechCorp",
		LocationRaw: "San Francisco, CA",
	}

	hash1 := ComputeDedupHash(job)
	hash2 := ComputeDedupHash(job)

	if hash1 != hash2 {
		t.Error("same job should produce same hash")
	}
}

func TestComputeDedupHash_DifferentJobs(t *testing.T) {
	job1 := &model.ScrapedJob{
		Source:      model.SourceLinkedIn,
		ExternalID:  "111",
		Title:       "Engineer",
		CompanyName: "Company A",
	}
	job2 := &model.ScrapedJob{
		Source:      model.SourceLinkedIn,
		ExternalID:  "222",
		Title:       "Engineer",
		CompanyName: "Company A",
	}

	hash1 := ComputeDedupHash(job1)
	hash2 := ComputeDedupHash(job2)

	if hash1 == hash2 {
		t.Error("different jobs should produce different hashes")
	}
}

func TestComputeDedupHash_SameJobDifferentSources(t *testing.T) {
	// Same job posted on different sources should have different hashes
	job1 := &model.ScrapedJob{
		Source:      model.SourceLinkedIn,
		ExternalID:  "abc123",
		Title:       "Engineer",
		CompanyName: "Company",
	}
	job2 := &model.ScrapedJob{
		Source:      model.SourceIndeed,
		ExternalID:  "abc123",
		Title:       "Engineer",
		CompanyName: "Company",
	}

	hash1 := ComputeDedupHash(job1)
	hash2 := ComputeDedupHash(job2)

	if hash1 == hash2 {
		t.Error("same external ID on different sources should produce different hashes")
	}
}

func TestNormalizeCompanyName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// The function strips known suffixes like ", inc.", " llc", ", ltd.", " corp."
		{"TechCorp, Inc.", "techcorp"},
		{"StartupXYZ LLC", "startupxyz"},
		{"BigCo, Ltd.", "bigco"},
		{"Company Corp.", "company"},
		// "Acme Corp." has no comma before "corp." so it strips " corp." suffix
		{"Acme Corp.", "acme"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeCompanyName(tt.input)
			if got != tt.want {
				t.Errorf("normalizeCompanyName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNullString(t *testing.T) {
	// Empty string should be null
	ns := nullString("")
	if ns.Valid {
		t.Error("empty string should produce null NullString")
	}

	// Non-empty string should be valid
	ns2 := nullString("hello")
	if !ns2.Valid {
		t.Error("non-empty string should produce valid NullString")
	}
	if ns2.String != "hello" {
		t.Errorf("NullString.String = %q, want %q", ns2.String, "hello")
	}
}
