package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/learnbot/api-gateway/internal/handler"
	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/api-gateway/internal/types"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test setup helpers
// ─────────────────────────────────────────────────────────────────────────────

// testServer creates a test HTTP server with all routes registered.
func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	jwtCfg := middleware.DefaultJWTConfig("test-secret")
	authMiddleware := middleware.RequireAuth(jwtCfg)

	mux := http.NewServeMux()

	authH := handler.NewAuthHandler(jwtCfg)
	profileH := handler.NewProfileHandler(jwtCfg)
	jobsH := handler.NewJobsHandler()
	analysisH := handler.NewAnalysisHandler()
	resourcesH := handler.NewResourcesHandler()

	authH.RegisterRoutes(mux)
	profileH.RegisterRoutes(mux, authMiddleware)
	jobsH.RegisterRoutes(mux, authMiddleware)
	analysisH.RegisterRoutes(mux, authMiddleware)
	resourcesH.RegisterRoutes(mux)

	return httptest.NewServer(mux)
}

// doRequest performs an HTTP request and returns the response.
func doRequest(t *testing.T, srv *httptest.Server, method, path string, body interface{}, token string) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, srv.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

// decodeResponse decodes a JSON response body into v.
func decodeResponse(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

// registerAndLogin registers a user and returns the JWT token.
func registerAndLogin(t *testing.T, srv *httptest.Server, email, password, fullName string) string {
	t.Helper()

	// Register.
	regResp := doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	}, "")
	if regResp.StatusCode != http.StatusCreated {
		t.Fatalf("register failed with status %d", regResp.StatusCode)
	}

	var regResult map[string]interface{}
	decodeResponse(t, regResp, &regResult)

	data, ok := regResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("register response missing data")
	}
	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatal("register response missing token")
	}
	return token
}

// ─────────────────────────────────────────────────────────────────────────────
// Auth integration tests
// ─────────────────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}, "")

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)

	if result["success"] != true {
		t.Error("expected success=true")
	}
	data := result["data"].(map[string]interface{})
	if data["token"] == "" {
		t.Error("expected non-empty token")
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    "not-an-email",
		Password: "password123",
		FullName: "Test User",
	}, "")

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid email, got %d", resp.StatusCode)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    "test2@example.com",
		Password: "short",
		FullName: "Test User",
	}, "")

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for short password, got %d", resp.StatusCode)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	email := "dup@example.com"
	// First registration.
	doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    email,
		Password: "password123",
		FullName: "User One",
	}, "")

	// Second registration with same email.
	resp := doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    email,
		Password: "password456",
		FullName: "User Two",
	}, "")

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 for duplicate email, got %d", resp.StatusCode)
	}
}

func TestLogin_Success(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	email := "login@example.com"
	password := "password123"

	// Register first.
	doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: "Login User",
	}, "")

	// Login.
	resp := doRequest(t, srv, http.MethodPost, "/api/auth/login", types.LoginRequest{
		Email:    email,
		Password: password,
	}, "")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	email := "wrongpw@example.com"
	doRequest(t, srv, http.MethodPost, "/api/auth/register", types.RegisterRequest{
		Email:    email,
		Password: "correctpassword",
		FullName: "User",
	}, "")

	resp := doRequest(t, srv, http.MethodPost, "/api/auth/login", types.LoginRequest{
		Email:    email,
		Password: "wrongpassword",
	}, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong password, got %d", resp.StatusCode)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Auth middleware tests
// ─────────────────────────────────────────────────────────────────────────────

func TestProtectedEndpoint_NoToken(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/users/profile", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", resp.StatusCode)
	}
}

func TestProtectedEndpoint_InvalidToken(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/users/profile", nil, "invalid.token.here")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 with invalid token, got %d", resp.StatusCode)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Profile integration tests
// ─────────────────────────────────────────────────────────────────────────────

func TestGetProfile_Authenticated(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "profile@example.com", "password123", "Profile User")

	resp := doRequest(t, srv, http.MethodGet, "/api/users/profile", nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "updateprofile@example.com", "password123", "Update User")

	headline := "Senior Software Engineer"
	resp := doRequest(t, srv, http.MethodPut, "/api/users/profile", types.ProfileUpdateRequest{
		Headline: &headline,
	}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].(map[string]interface{})
	if data["headline"] != headline {
		t.Errorf("expected headline=%q, got %v", headline, data["headline"])
	}
}

func TestUpdateSkills_Success(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "skills@example.com", "password123", "Skills User")

	resp := doRequest(t, srv, http.MethodPut, "/api/profile/skills", types.SkillUpdateRequest{
		Skills: []types.SkillInput{
			{Name: "Go", Proficiency: "advanced", IsPrimary: true},
			{Name: "Python", Proficiency: "intermediate"},
		},
	}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].(map[string]interface{})
	if data["count"].(float64) != 2 {
		t.Errorf("expected 2 skills, got %v", data["count"])
	}
}

func TestUpdateSkills_InvalidProficiency(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "badskills@example.com", "password123", "Bad Skills User")

	resp := doRequest(t, srv, http.MethodPut, "/api/profile/skills", types.SkillUpdateRequest{
		Skills: []types.SkillInput{
			{Name: "Go", Proficiency: "super-expert"},
		},
	}, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid proficiency, got %d", resp.StatusCode)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Job matching integration tests
// ─────────────────────────────────────────────────────────────────────────────

func TestJobSearch_NoFilters(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "jobsearch@example.com", "password123", "Job Search User")

	resp := doRequest(t, srv, http.MethodPost, "/api/jobs/search", types.JobSearchRequest{}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	// Should return all sample jobs.
	data := result["data"].([]interface{})
	if len(data) == 0 {
		t.Error("expected at least one job in results")
	}
}

func TestJobSearch_WithLocationFilter(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "jobsearch2@example.com", "password123", "Job Search User 2")

	resp := doRequest(t, srv, http.MethodPost, "/api/jobs/search", types.JobSearchRequest{
		LocationType: "remote",
	}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].([]interface{})
	for _, item := range data {
		job := item.(map[string]interface{})
		if job["location_type"] != "remote" {
			t.Errorf("expected location_type=remote, got %v", job["location_type"])
		}
	}
}

func TestJobDetail_ValidID(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/jobs/job-001", nil, "")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestJobDetail_InvalidID(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/jobs/nonexistent-job", nil, "")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestJobMatch_ValidID(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "jobmatch@example.com", "password123", "Job Match User")

	// Set some skills first.
	doRequest(t, srv, http.MethodPut, "/api/profile/skills", types.SkillUpdateRequest{
		Skills: []types.SkillInput{
			{Name: "Go", Proficiency: "advanced"},
			{Name: "PostgreSQL", Proficiency: "intermediate"},
		},
	}, token)

	resp := doRequest(t, srv, http.MethodGet, "/api/jobs/job-001/match", nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].(map[string]interface{})
	if data["overall_score"] == nil {
		t.Error("expected overall_score in response")
	}
	if data["recommendation"] == "" {
		t.Error("expected non-empty recommendation")
	}
}

func TestJobRecommendations_Authenticated(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "jobrec@example.com", "password123", "Job Rec User")

	resp := doRequest(t, srv, http.MethodGet, "/api/jobs/recommendations", nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Gap analysis integration tests
// ─────────────────────────────────────────────────────────────────────────────

func TestGapAnalysis_WithJobID(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "gapanalysis@example.com", "password123", "Gap Analysis User")

	// Set some skills.
	doRequest(t, srv, http.MethodPut, "/api/profile/skills", types.SkillUpdateRequest{
		Skills: []types.SkillInput{
			{Name: "Go", Proficiency: "advanced"},
		},
	}, token)

	resp := doRequest(t, srv, http.MethodPost, "/api/analysis/gaps", types.GapAnalysisRequest{
		JobID: "job-001",
	}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	data := result["data"].(map[string]interface{})
	if data["readiness_score"] == nil {
		t.Error("expected readiness_score in response")
	}
}

func TestGapAnalysis_WithInlineJob(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "gapanalysis2@example.com", "password123", "Gap Analysis User 2")

	resp := doRequest(t, srv, http.MethodPost, "/api/analysis/gaps", types.GapAnalysisRequest{
		Job: &types.JobRequirementsInput{
			Title:          "Python Developer",
			RequiredSkills: []string{"Python", "Django", "PostgreSQL"},
		},
	}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGapAnalysis_NoJob(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "gapanalysis3@example.com", "password123", "Gap Analysis User 3")

	resp := doRequest(t, srv, http.MethodPost, "/api/analysis/gaps", types.GapAnalysisRequest{}, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 when no job specified, got %d", resp.StatusCode)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Training recommendations integration tests
// ─────────────────────────────────────────────────────────────────────────────

func TestTrainingRecommendations_GET(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "training@example.com", "password123", "Training User")

	resp := doRequest(t, srv, http.MethodGet, "/api/training/recommendations?job_id=job-001", nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestTrainingRecommendations_POST(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	token := registerAndLogin(t, srv, "training2@example.com", "password123", "Training User 2")

	resp := doRequest(t, srv, http.MethodPost, "/api/training/recommendations", types.TrainingRecommendationRequest{
		JobID: "job-002",
		Preferences: types.LearningPreferencesInput{
			WeeklyHoursAvailable: 10,
			PreferFree:           true,
		},
	}, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].(map[string]interface{})
	if data["job_title"] == nil {
		t.Error("expected job_title in response")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource search integration tests
// ─────────────────────────────────────────────────────────────────────────────

func TestResourceSearch_NoFilters(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/resources/search", nil, "")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	data := result["data"].([]interface{})
	if len(data) == 0 {
		t.Error("expected at least one resource")
	}
}

func TestResourceSearch_BySkill(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/resources/search?skill=python", nil, "")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].([]interface{})
	if len(data) == 0 {
		t.Error("expected Python resources")
	}
}

func TestResourceSearch_FreeOnly(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/resources/search?free=true", nil, "")

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeResponse(t, resp, &result)
	data := result["data"].([]interface{})
	for _, item := range data {
		res := item.(map[string]interface{})
		costType := res["cost_type"].(string)
		if costType != "free" && costType != "free_audit" {
			t.Errorf("expected only free resources, got cost_type=%s", costType)
		}
	}
}

func TestResourceSearch_MethodNotAllowed(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/resources/search", nil, "")

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Response format tests
// ─────────────────────────────────────────────────────────────────────────────

func TestResponseFormat_ContentType(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/resources/search", nil, "")

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestResponseFormat_ErrorStructure(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/users/profile", nil, "")

	var result map[string]interface{}
	decodeResponse(t, resp, &result)

	if result["success"] != false {
		t.Error("expected success=false for error response")
	}
	if result["error"] == nil {
		t.Error("expected error field in error response")
	}
	errObj := result["error"].(map[string]interface{})
	if errObj["code"] == "" {
		t.Error("expected non-empty error code")
	}
	if errObj["message"] == "" {
		t.Error("expected non-empty error message")
	}
}
