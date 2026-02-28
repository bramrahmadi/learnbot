package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/learnbot/api-gateway/internal/handler"
	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/api-gateway/internal/types"
)

// newTestServer creates a test HTTP server without requiring *testing.T.
// Used for benchmarks and load tests.
func newTestServer() *httptest.Server {
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

// registerUserForBench registers a user and returns the JWT token for benchmarks.
func registerUserForBench(srv *httptest.Server, email string) (string, error) {
	body, _ := json.Marshal(types.RegisterRequest{
		Email:    email,
		Password: "password123",
		FullName: "Bench User",
	})
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no data in response")
	}
	token, ok := data["token"].(string)
	if !ok {
		return "", fmt.Errorf("no token in response")
	}
	return token, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Performance benchmarks for API endpoints
// ─────────────────────────────────────────────────────────────────────────────

// BenchmarkRegister benchmarks the registration endpoint.
func BenchmarkRegister(b *testing.B) {
	srv := newTestServer()
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(types.RegisterRequest{
			Email:    fmt.Sprintf("bench%d@example.com", i),
			Password: "password123",
			FullName: "Bench User",
		})
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

// BenchmarkJobSearch benchmarks the job search endpoint.
func BenchmarkJobSearch(b *testing.B) {
	srv := newTestServer()
	defer srv.Close()

	token, err := registerUserForBench(srv, "benchjobs@example.com")
	if err != nil {
		b.Skipf("could not register user for benchmark: %v", err)
		return
	}

	searchBody, _ := json.Marshal(types.JobSearchRequest{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/jobs/search", bytes.NewReader(searchBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

// BenchmarkResourceSearch benchmarks the resource search endpoint.
func BenchmarkResourceSearch(b *testing.B) {
	srv := newTestServer()
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/resources/search", nil)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

// BenchmarkGapAnalysis benchmarks the gap analysis endpoint.
func BenchmarkGapAnalysis(b *testing.B) {
	srv := newTestServer()
	defer srv.Close()

	token, err := registerUserForBench(srv, "benchgap@example.com")
	if err != nil {
		b.Skipf("could not register user for benchmark: %v", err)
		return
	}

	// Set skills first.
	skillsBody, _ := json.Marshal(types.SkillUpdateRequest{
		Skills: []types.SkillInput{
			{Name: "Go", Proficiency: "advanced"},
			{Name: "Python", Proficiency: "intermediate"},
		},
	})
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/api/profile/skills", bytes.NewReader(skillsBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	gapBody, _ := json.Marshal(types.GapAnalysisRequest{JobID: "job-001"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/analysis/gaps", bytes.NewReader(gapBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Load tests (concurrent requests)
// ─────────────────────────────────────────────────────────────────────────────

// TestConcurrentRegistrations tests that concurrent registrations work correctly.
func TestConcurrentRegistrations(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	const numConcurrent = 10
	var wg sync.WaitGroup
	errors := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body, _ := json.Marshal(types.RegisterRequest{
				Email:    fmt.Sprintf("concurrent%d@example.com", idx),
				Password: "password123",
				FullName: fmt.Sprintf("Concurrent User %d", idx),
			})
			req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errors <- fmt.Errorf("request %d failed: %v", idx, err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				errors <- fmt.Errorf("request %d got status %d", idx, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// TestConcurrentJobSearches tests that concurrent job searches work correctly.
func TestConcurrentJobSearches(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	// Register a user.
	body, _ := json.Marshal(types.RegisterRequest{
		Email:    "concjobs@example.com",
		Password: "password123",
		FullName: "Concurrent Jobs User",
	})
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var regResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&regResult)
	resp.Body.Close()

	data, ok := regResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("could not get registration data")
	}
	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatal("could not get token")
	}

	const numConcurrent = 20
	var wg sync.WaitGroup
	errors := make(chan error, numConcurrent)

	searchBody, _ := json.Marshal(types.JobSearchRequest{})

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/jobs/search", bytes.NewReader(searchBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errors <- fmt.Errorf("request %d failed: %v", idx, err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("request %d got status %d", idx, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// TestAPIResponseTime tests that API responses are within acceptable time limits.
func TestAPIResponseTime(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	// Register a user.
	body, _ := json.Marshal(types.RegisterRequest{
		Email:    "timing@example.com",
		Password: "password123",
		FullName: "Timing User",
	})
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var regResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&regResult)
	resp.Body.Close()

	data, ok := regResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("could not get registration data")
	}
	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatal("could not get token")
	}

	endpoints := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		token   string
		maxTime time.Duration
	}{
		{
			name:    "resource search",
			method:  http.MethodGet,
			path:    "/api/resources/search",
			maxTime: 500 * time.Millisecond,
		},
		{
			name:    "job search",
			method:  http.MethodPost,
			path:    "/api/jobs/search",
			body:    types.JobSearchRequest{},
			token:   token,
			maxTime: 500 * time.Millisecond,
		},
		{
			name:    "get profile",
			method:  http.MethodGet,
			path:    "/api/users/profile",
			token:   token,
			maxTime: 200 * time.Millisecond,
		},
	}

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			var bodyBytes []byte
			if ep.body != nil {
				bodyBytes, _ = json.Marshal(ep.body)
			}

			req, _ := http.NewRequest(ep.method, srv.URL+ep.path, bytes.NewReader(bodyBytes))
			if ep.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			if ep.token != "" {
				req.Header.Set("Authorization", "Bearer "+ep.token)
			}

			start := time.Now()
			resp, err := http.DefaultClient.Do(req)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if elapsed > ep.maxTime {
				t.Errorf("%s took %v, expected < %v", ep.name, elapsed, ep.maxTime)
			}
		})
	}
}
