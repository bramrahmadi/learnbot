package httpclient

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestClient_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer server.Close()

	cfg := DefaultConfig()
	cfg.RequestsPerMinute = 60
	client, err := New(cfg, log.New(os.Stderr, "", 0))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Get(context.Background(), server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_GetBody_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test body"))
	}))
	defer server.Close()

	cfg := DefaultConfig()
	cfg.RequestsPerMinute = 60
	client, err := New(cfg, log.New(os.Stderr, "", 0))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	body, err := client.GetBody(context.Background(), server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "test body" {
		t.Errorf("expected body 'test body', got %q", body)
	}
}

func TestClient_Retry_On_500(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	cfg := DefaultConfig()
	cfg.RequestsPerMinute = 600
	cfg.MaxRetries = 3
	cfg.RetryDelay = 10 * time.Millisecond
	client, err := New(cfg, log.New(os.Stderr, "", 0))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Get(context.Background(), server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 after retries, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestClient_MaxRetries_Exceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	cfg := DefaultConfig()
	cfg.RequestsPerMinute = 600
	cfg.MaxRetries = 2
	cfg.RetryDelay = 10 * time.Millisecond
	client, err := New(cfg, log.New(os.Stderr, "", 0))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Get(context.Background(), server.URL, nil)
	if err == nil {
		t.Error("expected error when max retries exceeded")
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := DefaultConfig()
	cfg.RequestsPerMinute = 600
	client, err := New(cfg, log.New(os.Stderr, "", 0))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = client.Get(ctx, server.URL, nil)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
}

func TestClient_CustomHeaders(t *testing.T) {
	var receivedUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := DefaultConfig()
	cfg.RequestsPerMinute = 600
	cfg.UserAgent = "TestAgent/1.0"
	client, err := New(cfg, log.New(os.Stderr, "", 0))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Get(context.Background(), server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if receivedUA != "TestAgent/1.0" {
		t.Errorf("expected User-Agent 'TestAgent/1.0', got %q", receivedUA)
	}
}

func TestIsRetryableStatus(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{200, false},
		{400, false},
		{404, false},
		{429, true},
		{500, true},
		{502, true},
		{503, true},
		{504, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := isRetryableStatus(tt.code)
			if got != tt.want {
				t.Errorf("isRetryableStatus(%d) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestRobotsChecker_IsAllowed(t *testing.T) {
	robotsTxt := `User-agent: *
Disallow: /private/
Disallow: /admin/

User-agent: TestBot
Disallow: /testonly/`

	tests := []struct {
		path      string
		userAgent string
		want      bool
	}{
		{"/public/page", "*", true},
		{"/private/secret", "*", false},
		{"/admin/panel", "*", false},
		{"/testonly/page", "TestBot", false},
		{"/testonly/page", "OtherBot", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isPathAllowed(robotsTxt, tt.userAgent, tt.path)
			if got != tt.want {
				t.Errorf("isPathAllowed(%q, %q, %q) = %v, want %v",
					robotsTxt[:20], tt.userAgent, tt.path, got, tt.want)
			}
		})
	}
}
