package scorer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// buildTestScorerHandler creates a Handler for testing.
func buildTestScorerHandler() *Handler {
	logger := log.New(os.Stderr, "[scorer-test] ", 0)
	return NewHandler(logger)
}

// buildScoreRequest creates a JSON-encoded ScoreRequest body.
func buildScoreRequest(t *testing.T, req ScoreRequest) *bytes.Buffer {
	t.Helper()
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}
	return bytes.NewBuffer(data)
}

// ─────────────────────────────────────────────────────────────────────────────
// ScoreHandler tests
// ─────────────────────────────────────────────────────────────────────────────

func TestScoreHandler_MethodNotAllowed(t *testing.T) {
	h := buildTestScorerHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/score", nil)
	w := httptest.NewRecorder()

	h.ScoreHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}

	var resp ScoreResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Success {
		t.Error("expected success=false")
	}
}

func TestScoreHandler_InvalidJSON(t *testing.T) {
	h := buildTestScorerHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/score",
		strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ScoreHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp ScoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Success {
		t.Error("expected success=false for invalid JSON")
	}
}

func TestScoreHandler_ValidRequest(t *testing.T) {
	h := buildTestScorerHandler()

	scoreReq := ScoreRequest{
		Profile: CandidateProfile{
			Skills: []CandidateSkill{
				{Name: "Go", Proficiency: "expert"},
				{Name: "PostgreSQL", Proficiency: "advanced"},
			},
			YearsOfExperience: 5,
			WorkHistory: []WorkHistoryEntry{
				{Title: "Software Engineer", Industry: "software", DurationMonths: 60},
			},
			Education: []EducationEntry{
				{DegreeLevel: "bachelor", FieldOfStudy: "Computer Science"},
			},
			RemotePreference: "remote",
		},
		Job: JobRequirements{
			Title:               "Software Engineer",
			RequiredSkills:      []string{"Go", "PostgreSQL"},
			MinYearsExperience:  3,
			RequiredDegreeLevel: "bachelor",
			LocationType:        "remote",
			Industry:            "software",
		},
	}

	body := buildScoreRequest(t, scoreReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/score", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ScoreHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ScoreResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true, got error: %s", resp.Error)
	}
	if resp.Data == nil {
		t.Fatal("expected data to be non-nil")
	}
	if resp.Data.OverallScore < 0 || resp.Data.OverallScore > 100 {
		t.Errorf("overall score out of range: %.2f", resp.Data.OverallScore)
	}
}

func TestScoreHandler_EmptyRequest(t *testing.T) {
	h := buildTestScorerHandler()

	// Empty but valid JSON object.
	body := strings.NewReader("{}")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/score", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ScoreHandler(w, req)

	// Empty request is valid – should return a score (likely low).
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for empty request, got %d", w.Code)
	}

	var resp ScoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Errorf("expected success=true for empty request, got: %s", resp.Error)
	}
}

func TestScoreHandler_ResponseContainsBreakdown(t *testing.T) {
	h := buildTestScorerHandler()

	scoreReq := ScoreRequest{
		Profile: CandidateProfile{
			Skills: []CandidateSkill{{Name: "Python", Proficiency: "expert"}},
		},
		Job: JobRequirements{
			RequiredSkills: []string{"Python", "Java"},
		},
	}

	body := buildScoreRequest(t, scoreReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/score", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ScoreHandler(w, req)

	var resp ScoreResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Data == nil {
		t.Fatal("expected data to be non-nil")
	}

	// Verify breakdown fields are present.
	if resp.Data.SkillMatchScore < 0 || resp.Data.SkillMatchScore > 1 {
		t.Errorf("skill match score out of range: %.2f", resp.Data.SkillMatchScore)
	}
	if resp.Data.ExperienceMatchScore < 0 || resp.Data.ExperienceMatchScore > 1 {
		t.Errorf("experience match score out of range: %.2f", resp.Data.ExperienceMatchScore)
	}
	if resp.Data.EducationMatchScore < 0 || resp.Data.EducationMatchScore > 1 {
		t.Errorf("education match score out of range: %.2f", resp.Data.EducationMatchScore)
	}
	if resp.Data.LocationFitScore < 0 || resp.Data.LocationFitScore > 1 {
		t.Errorf("location fit score out of range: %.2f", resp.Data.LocationFitScore)
	}
	if resp.Data.IndustryRelevanceScore < 0 || resp.Data.IndustryRelevanceScore > 1 {
		t.Errorf("industry relevance score out of range: %.2f", resp.Data.IndustryRelevanceScore)
	}

	// Python should be matched, Java should be missing.
	mustContain(t, resp.Data.MatchedRequiredSkills, "Python")
	mustContain(t, resp.Data.MissingRequiredSkills, "Java")
}

func TestScoreHandler_ContentTypeJSON(t *testing.T) {
	h := buildTestScorerHandler()

	body := strings.NewReader("{}")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/score", body)
	// No Content-Type header set.
	w := httptest.NewRecorder()

	h.ScoreHandler(w, req)

	// Should still work (Content-Type check is lenient for empty header).
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 when no Content-Type, got %d", w.Code)
	}
}

func TestRegisterScorerRoutes(t *testing.T) {
	h := buildTestScorerHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// POST to /api/v1/score should be handled.
	body := strings.NewReader("{}")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/score", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 from registered route, got %d", w.Code)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Benchmarks
// ─────────────────────────────────────────────────────────────────────────────

// BenchmarkCalculate_Typical benchmarks a typical scoring calculation.
func BenchmarkCalculate_Typical(b *testing.B) {
	profile := CandidateProfile{
		Skills: []CandidateSkill{
			{Name: "Go", Proficiency: "expert"},
			{Name: "Python", Proficiency: "advanced"},
			{Name: "PostgreSQL", Proficiency: "advanced"},
			{Name: "Docker", Proficiency: "intermediate"},
			{Name: "Kubernetes", Proficiency: "intermediate"},
			{Name: "AWS", Proficiency: "intermediate"},
			{Name: "REST", Proficiency: "expert"},
			{Name: "Microservices", Proficiency: "advanced"},
		},
		YearsOfExperience: 6,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Senior Software Engineer", Industry: "software", DurationMonths: 36, IsCurrent: true},
			{Title: "Software Engineer", Industry: "fintech", DurationMonths: 24},
			{Title: "Junior Developer", Industry: "software", DurationMonths: 12},
		},
		Education: []EducationEntry{
			{DegreeLevel: "master", FieldOfStudy: "Computer Science"},
		},
		LocationCity:     "San Francisco",
		LocationCountry:  "US",
		RemotePreference: "hybrid",
	}

	job := JobRequirements{
		Title:               "Senior Software Engineer",
		RequiredSkills:      []string{"Go", "PostgreSQL", "Docker", "Kubernetes"},
		PreferredSkills:     []string{"Python", "AWS", "Terraform"},
		MinYearsExperience:  5,
		RequiredDegreeLevel: "bachelor",
		PreferredFields:     []string{"Computer Science", "Software Engineering"},
		LocationCity:        "San Francisco",
		LocationCountry:     "US",
		LocationType:        "hybrid",
		Industry:            "software",
		RelatedIndustries:   []string{"fintech", "saas"},
		ExperienceLevel:     "senior",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Calculate(profile, job)
	}
}

// BenchmarkCalculate_LargeSkillSet benchmarks scoring with many skills.
func BenchmarkCalculate_LargeSkillSet(b *testing.B) {
	// Generate a large skill set.
	skillNames := []string{
		"Go", "Python", "Java", "JavaScript", "TypeScript", "Rust", "C++", "C#",
		"Ruby", "PHP", "Swift", "Kotlin", "Scala", "R", "MATLAB",
		"React", "Angular", "Vue", "Node.js", "Django", "Flask", "FastAPI",
		"Spring", "Rails", "Laravel", "Gin", "Echo",
		"PostgreSQL", "MySQL", "MongoDB", "Redis", "Elasticsearch", "Cassandra",
		"AWS", "Azure", "GCP", "Docker", "Kubernetes", "Terraform", "Ansible",
		"Git", "GitHub", "GitLab", "Jira", "Confluence",
		"Kafka", "RabbitMQ", "Celery", "Airflow",
		"TensorFlow", "PyTorch", "Keras", "scikit-learn",
		"REST", "GraphQL", "gRPC", "Microservices", "CI/CD", "DevOps",
	}

	skills := make([]CandidateSkill, len(skillNames))
	for i, name := range skillNames {
		skills[i] = CandidateSkill{Name: name, Proficiency: "advanced"}
	}

	profile := CandidateProfile{
		Skills:            skills,
		YearsOfExperience: 10,
		WorkHistory: []WorkHistoryEntry{
			{Title: "Principal Engineer", Industry: "software", DurationMonths: 60},
		},
		Education: []EducationEntry{
			{DegreeLevel: "master", FieldOfStudy: "Computer Science"},
		},
		RemotePreference: "remote",
	}

	job := JobRequirements{
		Title:               "Staff Engineer",
		RequiredSkills:      skillNames[:20],
		PreferredSkills:     skillNames[20:40],
		MinYearsExperience:  8,
		RequiredDegreeLevel: "bachelor",
		LocationType:        "remote",
		Industry:            "software",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Calculate(profile, job)
	}
}

// BenchmarkCalculate_MinimalInput benchmarks scoring with minimal input.
func BenchmarkCalculate_MinimalInput(b *testing.B) {
	profile := CandidateProfile{}
	job := JobRequirements{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Calculate(profile, job)
	}
}

// BenchmarkScoreHandler_HTTP benchmarks the full HTTP handler path.
func BenchmarkScoreHandler_HTTP(b *testing.B) {
	h := buildTestScorerHandler()

	scoreReq := ScoreRequest{
		Profile: CandidateProfile{
			Skills: []CandidateSkill{
				{Name: "Go", Proficiency: "expert"},
				{Name: "PostgreSQL", Proficiency: "advanced"},
			},
			YearsOfExperience: 5,
			WorkHistory: []WorkHistoryEntry{
				{Title: "Software Engineer", Industry: "software", DurationMonths: 60},
			},
			Education: []EducationEntry{
				{DegreeLevel: "bachelor", FieldOfStudy: "Computer Science"},
			},
			RemotePreference: "remote",
		},
		Job: JobRequirements{
			Title:               "Software Engineer",
			RequiredSkills:      []string{"Go", "PostgreSQL"},
			MinYearsExperience:  3,
			RequiredDegreeLevel: "bachelor",
			LocationType:        "remote",
			Industry:            "software",
		},
	}

	body, _ := json.Marshal(scoreReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/score", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ScoreHandler(w, req)
	}
}
