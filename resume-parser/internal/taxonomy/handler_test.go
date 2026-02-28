package taxonomy

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

// buildTestHandler creates a Handler for testing.
func buildTestHandler() *Handler {
	logger := log.New(os.Stderr, "[taxonomy-test] ", 0)
	return NewHandler(logger)
}

// ─────────────────────────────────────────────────────────────────────────────
// ExtractHandler
// ─────────────────────────────────────────────────────────────────────────────

func TestExtractHandler_MethodNotAllowed(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/extract", nil)
	w := httptest.NewRecorder()
	h.ExtractHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestExtractHandler_InvalidJSON(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/extract",
		strings.NewReader("{invalid}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ExtractHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExtractHandler_ValidRequest(t *testing.T) {
	h := buildTestHandler()

	body, _ := json.Marshal(ExtractRequest{
		Text: "We need a Go developer with Kubernetes and PostgreSQL experience.",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ExtractHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ExtractResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true, got error: %s", resp.Error)
	}
	if resp.Data == nil {
		t.Fatal("expected data to be non-nil")
	}
	if len(resp.Data.TechnicalSkills) == 0 {
		t.Error("expected technical skills to be extracted")
	}
}

func TestExtractHandler_EmptyText(t *testing.T) {
	h := buildTestHandler()

	body, _ := json.Marshal(ExtractRequest{Text: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ExtractHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for empty text, got %d", w.Code)
	}

	var resp ExtractResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Errorf("expected success=true for empty text")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// NormalizeHandler
// ─────────────────────────────────────────────────────────────────────────────

func TestNormalizeHandler_MethodNotAllowed(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/normalize", nil)
	w := httptest.NewRecorder()
	h.NormalizeHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestNormalizeHandler_EmptySkills(t *testing.T) {
	h := buildTestHandler()
	body, _ := json.Marshal(NormalizeRequest{Skills: []string{}})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/normalize", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.NormalizeHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty skills, got %d", w.Code)
	}
}

func TestNormalizeHandler_ValidRequest(t *testing.T) {
	h := buildTestHandler()

	body, _ := json.Marshal(NormalizeRequest{
		Skills: []string{"golang", "k8s", "postgres", "unknown-skill-xyz"},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/normalize", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.NormalizeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp NormalizeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true, got error: %s", resp.Error)
	}
	if len(resp.Data) != 4 {
		t.Errorf("expected 4 results, got %d", len(resp.Data))
	}

	// Check specific normalizations.
	if resp.Data[0].CanonicalID != "go" {
		t.Errorf("expected 'go' for 'golang', got %q", resp.Data[0].CanonicalID)
	}
	if resp.Data[1].CanonicalID != "kubernetes" {
		t.Errorf("expected 'kubernetes' for 'k8s', got %q", resp.Data[1].CanonicalID)
	}
	if resp.Data[2].CanonicalID != "postgresql" {
		t.Errorf("expected 'postgresql' for 'postgres', got %q", resp.Data[2].CanonicalID)
	}
	if resp.Data[3].MatchType != "none" {
		t.Errorf("expected 'none' for unknown skill, got %q", resp.Data[3].MatchType)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// LookupHandler
// ─────────────────────────────────────────────────────────────────────────────

func TestLookupHandler_MethodNotAllowed(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/lookup", nil)
	w := httptest.NewRecorder()
	h.LookupHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestLookupHandler_MissingID(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/lookup", nil)
	w := httptest.NewRecorder()
	h.LookupHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing id, got %d", w.Code)
	}
}

func TestLookupHandler_NotFound(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/lookup?id=nonexistent", nil)
	w := httptest.NewRecorder()
	h.LookupHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent skill, got %d", w.Code)
	}
}

func TestLookupHandler_Found(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/lookup?id=go", nil)
	w := httptest.NewRecorder()
	h.LookupHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp LookupResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true, got error: %s", resp.Error)
	}
	if resp.Data == nil {
		t.Fatal("expected data to be non-nil")
	}
	if resp.Data.ID != "go" {
		t.Errorf("expected ID 'go', got %q", resp.Data.ID)
	}
	if resp.Data.CanonicalName != "Go" {
		t.Errorf("expected canonical name 'Go', got %q", resp.Data.CanonicalName)
	}
	if len(resp.Data.Aliases) == 0 {
		t.Error("expected aliases to be non-empty")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SearchHandler
// ─────────────────────────────────────────────────────────────────────────────

func TestSearchHandler_MethodNotAllowed(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/skills/search", nil)
	w := httptest.NewRecorder()
	h.SearchHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestSearchHandler_EmptyQuery(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/search", nil)
	w := httptest.NewRecorder()
	h.SearchHandler(w, req)

	// Empty query should return all skills (up to limit).
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for empty query, got %d", w.Code)
	}

	var resp SearchResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Errorf("expected success=true")
	}
	if resp.Total == 0 {
		t.Error("expected results for empty query")
	}
}

func TestSearchHandler_WithQuery(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/search?q=python", nil)
	w := httptest.NewRecorder()
	h.SearchHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp SearchResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Total == 0 {
		t.Error("expected results for query 'python'")
	}
}

func TestSearchHandler_WithDomainFilter(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/skills/search?domain=data_science", nil)
	w := httptest.NewRecorder()
	h.SearchHandler(w, req)

	var resp SearchResponse
	json.NewDecoder(w.Body).Decode(&resp)

	for _, skill := range resp.Data {
		if skill.Domain != DomainDataScience {
			t.Errorf("expected domain %q, got %q for skill %q",
				DomainDataScience, skill.Domain, skill.ID)
		}
	}
}

func TestSearchHandler_WithLimit(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills/search?limit=3", nil)
	w := httptest.NewRecorder()
	h.SearchHandler(w, req)

	var resp SearchResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp.Data) > 3 {
		t.Errorf("expected at most 3 results, got %d", len(resp.Data))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// RegisterRoutes
// ─────────────────────────────────────────────────────────────────────────────

func TestRegisterTaxonomyRoutes(t *testing.T) {
	h := buildTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	routes := []struct {
		method string
		path   string
		body   string
		want   int
	}{
		{http.MethodPost, "/api/v1/skills/extract", `{"text":"Go Python"}`, http.StatusOK},
		{http.MethodPost, "/api/v1/skills/normalize", `{"skills":["golang"]}`, http.StatusOK},
		{http.MethodGet, "/api/v1/skills/lookup?id=go", "", http.StatusOK},
		{http.MethodGet, "/api/v1/skills/search?q=python", "", http.StatusOK},
	}

	for _, tt := range routes {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			var body *strings.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			} else {
				body = strings.NewReader("")
			}
			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Errorf("expected %d, got %d: %s", tt.want, w.Code, w.Body.String())
			}
		})
	}
}
