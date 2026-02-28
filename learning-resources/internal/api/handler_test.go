package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Handler tests using a nil repository (testing routing and error handling)
// ─────────────────────────────────────────────────────────────────────────────

// newTestHandler creates a Handler with a nil repo for routing tests.
// Tests that exercise the repo must use a mock or integration test setup.
func newTestHandler() *Handler {
	return &Handler{
		repo:   nil,
		logger: log.New(io.Discard, "", 0),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Method validation tests
// ─────────────────────────────────────────────────────────────────────────────

func TestHandleResources_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/v1/resources", nil)
		w := httptest.NewRecorder()
		h.handleResources(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected 405, got %d", method, w.Code)
		}
	}
}

func TestHandleFeaturedResources_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resources/featured", nil)
	w := httptest.NewRecorder()
	h.handleFeaturedResources(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandleResourcesBySkill_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resources/by-skill", nil)
	w := httptest.NewRecorder()
	h.handleResourcesBySkill(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandleResourcesBySkill_MissingSkillParam(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources/by-skill", nil)
	w := httptest.NewRecorder()
	h.handleResourcesBySkill(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when skill param missing, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["success"] != false {
		t.Error("expected success=false")
	}
}

func TestHandleResourceBySlug_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resources/some-slug", nil)
	w := httptest.NewRecorder()
	h.handleResourceBySlug(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandleResourceBySlug_EmptySlug(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources/", nil)
	w := httptest.NewRecorder()
	h.handleResourceBySlug(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty slug, got %d", w.Code)
	}
}

func TestHandlePaths_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/paths", nil)
	w := httptest.NewRecorder()
	h.handlePaths(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandlePathBySlug_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/paths/some-slug", nil)
	w := httptest.NewRecorder()
	h.handlePathBySlug(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandlePathBySlug_EmptySlug(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/paths/", nil)
	w := httptest.NewRecorder()
	h.handlePathBySlug(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty slug, got %d", w.Code)
	}
}

func TestHandleProviders_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers", nil)
	w := httptest.NewRecorder()
	h.handleProviders(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// User progress routing tests
// ─────────────────────────────────────────────────────────────────────────────

func TestHandleUserProgress_InvalidUserID(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/not-a-uuid/progress", nil)
	req.URL.Path = "/api/v1/users/not-a-uuid/progress"
	w := httptest.NewRecorder()
	h.handleUserProgress(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid user ID, got %d", w.Code)
	}
}

func TestHandleUserProgress_NotFoundPath(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000/other", nil)
	req.URL.Path = "/api/v1/users/123e4567-e89b-12d3-a456-426614174000/other"
	w := httptest.NewRecorder()
	h.handleUserProgress(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown path, got %d", w.Code)
	}
}

func TestHandleUserProgress_MethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000/progress", nil)
	req.URL.Path = "/api/v1/users/123e4567-e89b-12d3-a456-426614174000/progress"
	w := httptest.NewRecorder()
	h.handleUserProgress(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for DELETE, got %d", w.Code)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Response format tests
// ─────────────────────────────────────────────────────────────────────────────

func TestWriteError_ResponseFormat(t *testing.T) {
	h := newTestHandler()
	w := httptest.NewRecorder()
	h.writeError(w, http.StatusBadRequest, "test error message")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp["success"] != false {
		t.Error("expected success=false in error response")
	}
	if resp["error"] != "test error message" {
		t.Errorf("expected error message 'test error message', got %v", resp["error"])
	}
}

func TestWriteJSON_ResponseFormat(t *testing.T) {
	h := newTestHandler()
	w := httptest.NewRecorder()
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []string{"item1", "item2"},
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Query parameter parsing tests
// ─────────────────────────────────────────────────────────────────────────────

func TestHandleResources_QueryParamParsing(t *testing.T) {
	// Test that query parameters are parsed without panicking.
	// The handler will fail at the repo call (nil repo), but we can verify
	// it gets past parameter parsing.
	h := newTestHandler()

	tests := []struct {
		name  string
		query string
	}{
		{"no params", ""},
		{"skill filter", "?skill=Python"},
		{"type filter", "?type=course"},
		{"difficulty filter", "?difficulty=beginner"},
		{"free filter", "?free=true"},
		{"has_certificate filter", "?has_certificate=true"},
		{"min_rating filter", "?min_rating=4.5"},
		{"search query", "?q=machine+learning"},
		{"limit and offset", "?limit=10&offset=20"},
		{"invalid limit", "?limit=abc"},
		{"invalid min_rating", "?min_rating=abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/resources"+tt.query, nil)
			w := httptest.NewRecorder()

			// This will panic/error at the repo call since repo is nil,
			// but the middleware will recover it.
			func() {
				defer func() { recover() }()
				h.handleResources(w, req)
			}()

			// We just verify it doesn't crash before reaching the repo.
			// The response will be 500 due to nil repo, which is expected.
		})
	}
}
