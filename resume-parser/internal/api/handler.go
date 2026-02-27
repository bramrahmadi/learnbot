// Package api provides the HTTP API layer for the resume parser.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/learnbot/resume-parser/internal/parser"
	"github.com/learnbot/resume-parser/internal/schema"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
)

// Handler holds the HTTP handler dependencies.
type Handler struct {
	parser *parser.ResumeParser
	logger *log.Logger
}

// NewHandler creates a new API Handler.
func NewHandler(p *parser.ResumeParser, logger *log.Logger) *Handler {
	return &Handler{
		parser: p,
		logger: logger,
	}
}

// RegisterRoutes registers all API routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/parse", h.withMiddleware(h.ParseResume))
	mux.HandleFunc("/api/v1/health", h.withMiddleware(h.HealthCheck))
}

// withMiddleware wraps a handler with logging and recovery middleware.
func (h *Handler) withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				h.logger.Printf("PANIC: %v", rec)
				h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred", "")
			}
		}()

		h.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		h.logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// ParseResume handles POST /api/v1/parse
// Accepts multipart/form-data with a "resume" file field.
// Optional query param: include_raw=true to include raw text in response.
//
// Example:
//
//	curl -X POST http://localhost:8080/api/v1/parse \
//	  -F "resume=@/path/to/resume.pdf"
func (h *Handler) ParseResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
			"only POST is supported", "")
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID_REQUEST",
			fmt.Sprintf("failed to parse multipart form: %v", err), "")
		return
	}

	file, header, err := r.FormFile("resume")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "MISSING_FILE",
			"'resume' file field is required", "")
		return
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "READ_ERROR",
			"failed to read uploaded file", "")
		return
	}

	// Determine file type
	fileType := detectFileType(header.Filename, header.Header.Get("Content-Type"))
	if fileType == "" {
		h.writeError(w, http.StatusBadRequest, "UNSUPPORTED_FORMAT",
			"unsupported file type; only PDF and DOCX are supported", "")
		return
	}

	includeRaw := strings.ToLower(r.FormValue("include_raw")) == "true"

	req := schema.ParseRequest{
		FileName:    header.Filename,
		FileContent: data,
		FileType:    fileType,
		IncludeRaw:  includeRaw,
	}

	parsed, err := h.parser.Parse(req)
	if err != nil {
		if pe, ok := err.(*schema.ParseError); ok {
			statusCode := parseErrorToHTTPStatus(pe.Code)
			h.writeError(w, statusCode, pe.Code, pe.Message, pe.Section)
			return
		}
		h.writeError(w, http.StatusInternalServerError, "PARSE_ERROR", err.Error(), "")
		return
	}

	resp := schema.ParseResponse{
		Success: true,
		Data:    parsed,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// HealthCheck handles GET /api/v1/health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

// writeJSON serializes v as JSON and writes it to the response.
func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Printf("failed to encode JSON response: %v", err)
	}
}

// writeError writes a structured error response.
func (h *Handler) writeError(w http.ResponseWriter, status int, code, message, section string) {
	resp := schema.ParseResponse{
		Success: false,
		Error: &schema.ParseError{
			Code:    code,
			Message: message,
			Section: section,
		},
	}
	h.writeJSON(w, status, resp)
}

// detectFileType returns the canonical file type from filename and content-type.
func detectFileType(filename, contentType string) string {
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".pdf") {
		return "pdf"
	}
	if strings.HasSuffix(lower, ".docx") {
		return "docx"
	}

	ct := strings.ToLower(contentType)
	if strings.Contains(ct, "pdf") {
		return "pdf"
	}
	if strings.Contains(ct, "wordprocessingml") || strings.Contains(ct, "docx") {
		return "docx"
	}

	return ""
}

// parseErrorToHTTPStatus maps error codes to HTTP status codes.
func parseErrorToHTTPStatus(code string) int {
	switch code {
	case "EMPTY_FILE", "INVALID_FORMAT", "UNSUPPORTED_FORMAT",
		"MISSING_FILE", "INVALID_REQUEST":
		return http.StatusBadRequest
	case "NO_TEXT_CONTENT":
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
