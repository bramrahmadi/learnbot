package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/learnbot/resume-parser/internal/parser"
	"github.com/learnbot/resume-parser/internal/schema"
)

// buildTestHandler creates a Handler for testing.
func buildTestHandler() *Handler {
	p := parser.NewResumeParser()
	logger := log.New(os.Stderr, "[test] ", 0)
	return NewHandler(p, logger)
}

// buildMinimalDOCX creates a minimal valid DOCX file with the given text content.
func buildMinimalDOCX(text string) []byte {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	docXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>%s</w:body>
</w:document>`, textToParagraphs(text))

	f, _ := w.Create("word/document.xml")
	f.Write([]byte(docXML))

	ct, _ := w.Create("[Content_Types].xml")
	ct.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`))

	w.Close()
	return buf.Bytes()
}

func textToParagraphs(text string) string {
	var sb bytes.Buffer
	lines := bytes.Split([]byte(text), []byte("\n"))
	for _, line := range lines {
		sb.WriteString(`<w:p><w:r><w:t>`)
		sb.Write(line)
		sb.WriteString(`</w:t></w:r></w:p>`)
	}
	return sb.String()
}

// createMultipartRequest creates a multipart/form-data request with a file upload.
func createMultipartRequest(t *testing.T, filename string, content []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("resume", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		t.Fatalf("failed to write file content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/parse", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// TestHealthCheck tests the health check endpoint.
func TestHealthCheck(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}
}

// TestParseResume_MethodNotAllowed tests that GET is rejected.
func TestParseResume_MethodNotAllowed(t *testing.T) {
	h := buildTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/parse", nil)
	w := httptest.NewRecorder()

	h.ParseResume(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

// TestParseResume_MissingFile tests that missing file field returns 400.
func TestParseResume_MissingFile(t *testing.T) {
	h := buildTestHandler()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/parse", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.ParseResume(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp schema.ParseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Success {
		t.Error("expected success=false")
	}
}

// TestParseResume_InvalidPDF tests that invalid PDF returns error.
func TestParseResume_InvalidPDF(t *testing.T) {
	h := buildTestHandler()
	req := createMultipartRequest(t, "bad.pdf", []byte("not a pdf"))
	w := httptest.NewRecorder()

	h.ParseResume(w, req)

	if w.Code == http.StatusOK {
		t.Error("expected non-200 status for invalid PDF")
	}

	var resp schema.ParseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Success {
		t.Error("expected success=false for invalid PDF")
	}
}

// TestParseResume_UnsupportedFormat tests that unsupported file types return 400.
func TestParseResume_UnsupportedFormat(t *testing.T) {
	h := buildTestHandler()
	req := createMultipartRequest(t, "resume.txt", []byte("some text"))
	w := httptest.NewRecorder()

	h.ParseResume(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestParseResume_ValidDOCX tests successful DOCX parsing.
func TestParseResume_ValidDOCX(t *testing.T) {
	h := buildTestHandler()

	resumeText := `Jane Doe
jane@example.com
555-123-4567

SKILLS
Go, Python, Docker`

	docxData := buildMinimalDOCX(resumeText)
	req := createMultipartRequest(t, "resume.docx", docxData)
	w := httptest.NewRecorder()

	h.ParseResume(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp schema.ParseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true, got error: %v", resp.Error)
	}
	if resp.Data == nil {
		t.Fatal("expected data to be non-nil")
	}
}

// TestParseResume_IncludeRaw tests the include_raw query parameter.
func TestParseResume_IncludeRaw(t *testing.T) {
	h := buildTestHandler()

	resumeText := `Jane Doe
jane@example.com

SKILLS
Go, Python`

	docxData := buildMinimalDOCX(resumeText)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("resume", "resume.docx")
	io.Copy(part, bytes.NewReader(docxData))
	writer.WriteField("include_raw", "true")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/parse?include_raw=true", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.ParseResume(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp schema.ParseResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Data != nil && resp.Data.RawText == "" {
		t.Error("expected raw text to be included")
	}
}

// TestDetectFileType tests file type detection.
func TestDetectFileType(t *testing.T) {
	tests := []struct {
		filename    string
		contentType string
		want        string
	}{
		{"resume.pdf", "", "pdf"},
		{"resume.docx", "", "docx"},
		{"resume.PDF", "", "pdf"},
		{"resume.DOCX", "", "docx"},
		{"resume", "application/pdf", "pdf"},
		{"resume", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "docx"},
		{"resume.txt", "text/plain", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := detectFileType(tt.filename, tt.contentType)
			if got != tt.want {
				t.Errorf("detectFileType(%q, %q) = %q, want %q",
					tt.filename, tt.contentType, got, tt.want)
			}
		})
	}
}

// TestParseErrorToHTTPStatus tests error code to HTTP status mapping.
func TestParseErrorToHTTPStatus(t *testing.T) {
	tests := []struct {
		code string
		want int
	}{
		{"EMPTY_FILE", http.StatusBadRequest},
		{"INVALID_FORMAT", http.StatusBadRequest},
		{"UNSUPPORTED_FORMAT", http.StatusBadRequest},
		{"NO_TEXT_CONTENT", http.StatusUnprocessableEntity},
		{"UNKNOWN_ERROR", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got := parseErrorToHTTPStatus(tt.code)
			if got != tt.want {
				t.Errorf("parseErrorToHTTPStatus(%q) = %d, want %d", tt.code, got, tt.want)
			}
		})
	}
}

// TestRegisterRoutes tests that routes are registered correctly.
func TestRegisterRoutes(t *testing.T) {
	h := buildTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Test health endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected health endpoint to return 200, got %d", w.Code)
	}
}
