// Package handler provides HTTP handlers and response utilities for the API gateway.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/learnbot/api-gateway/internal/types"
)

// ─────────────────────────────────────────────────────────────────────────────
// Response helpers
// ─────────────────────────────────────────────────────────────────────────────

// WriteSuccess writes a successful JSON response.
func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, types.APIResponse{
		Success: true,
		Data:    data,
	})
}

// WriteSuccessWithMeta writes a successful JSON response with pagination metadata.
func WriteSuccessWithMeta(w http.ResponseWriter, status int, data interface{}, meta *types.ResponseMeta) {
	writeJSON(w, status, types.APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// WriteError writes a structured error JSON response.
func WriteError(w http.ResponseWriter, status int, code, message string, details ...types.FieldError) {
	apiErr := &types.APIError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		apiErr.Details = details
	}
	writeJSON(w, status, types.APIResponse{
		Success: false,
		Error:   apiErr,
	})
}

// WriteValidationError writes a 400 Bad Request with field-level errors.
func WriteValidationError(w http.ResponseWriter, details []types.FieldError) {
	WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR",
		"request validation failed", details...)
}

// WriteNotFound writes a 404 Not Found response.
func WriteNotFound(w http.ResponseWriter, resource string) {
	WriteError(w, http.StatusNotFound, "NOT_FOUND",
		resource+" not found")
}

// WriteInternalError writes a 500 Internal Server Error response.
func WriteInternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR",
		"an unexpected error occurred")
}

// WriteMethodNotAllowed writes a 405 Method Not Allowed response.
func WriteMethodNotAllowed(w http.ResponseWriter) {
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
		"HTTP method not allowed")
}

// writeJSON serializes v as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// ─────────────────────────────────────────────────────────────────────────────
// Request parsing helpers
// ─────────────────────────────────────────────────────────────────────────────

// DecodeJSON decodes the request body as JSON into v.
// Returns false and writes an error response if decoding fails.
func DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JSON",
			"invalid request body: "+err.Error())
		return false
	}
	return true
}

// ─────────────────────────────────────────────────────────────────────────────
// Validation helpers
// ─────────────────────────────────────────────────────────────────────────────

// Validator collects field validation errors.
type Validator struct {
	errors []types.FieldError
}

// Required validates that a string field is non-empty.
func (v *Validator) Required(field, value, message string) {
	if len(trimSpaceStr(value)) == 0 {
		v.errors = append(v.errors, types.FieldError{Field: field, Message: message})
	}
}

// MinLength validates that a string field meets a minimum length.
func (v *Validator) MinLength(field, value string, min int, message string) {
	if len(value) < min {
		v.errors = append(v.errors, types.FieldError{Field: field, Message: message})
	}
}

// ValidEmail validates that a string looks like an email address.
func (v *Validator) ValidEmail(field, value string) {
	if !isValidEmail(value) {
		v.errors = append(v.errors, types.FieldError{
			Field:   field,
			Message: "must be a valid email address",
		})
	}
}

// HasErrors returns true if any validation errors were collected.
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns the collected validation errors.
func (v *Validator) Errors() []types.FieldError {
	return v.errors
}

// WriteIfInvalid writes a validation error response if there are errors.
// Returns true if errors were written (caller should return).
func (v *Validator) WriteIfInvalid(w http.ResponseWriter) bool {
	if v.HasErrors() {
		WriteValidationError(w, v.errors)
		return true
	}
	return false
}

// isValidEmail performs a basic email format check.
func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	atIdx := -1
	for i, c := range email {
		if c == '@' {
			if atIdx >= 0 {
				return false // multiple @
			}
			atIdx = i
		}
	}
	if atIdx <= 0 || atIdx >= len(email)-1 {
		return false
	}
	// Check for dot in domain part.
	domain := email[atIdx+1:]
	hasDot := false
	for _, c := range domain {
		if c == '.' {
			hasDot = true
			break
		}
	}
	return hasDot
}

// trimSpaceStr removes leading and trailing whitespace.
func trimSpaceStr(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
