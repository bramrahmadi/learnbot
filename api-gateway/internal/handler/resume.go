// Package handler – resume.go implements resume upload and parsing endpoints.
package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/resume-parser/pkg/parse"
)

// ─────────────────────────────────────────────────────────────────────────────
// In-memory resume store (MVP)
// ─────────────────────────────────────────────────────────────────────────────

// resumeRecord stores a parsed resume.
type resumeRecord struct {
	ID         string
	UserID     string
	FileName   string
	ParsedAt   time.Time
	ParsedData *parse.ParsedResume
}

// resumeStore is a thread-safe in-memory resume store.
type resumeStore struct {
	mu      sync.RWMutex
	resumes map[string]*resumeRecord // keyed by userID (latest resume)
}

var globalResumeStore = &resumeStore{
	resumes: make(map[string]*resumeRecord),
}

func (s *resumeStore) save(userID string, rec *resumeRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resumes[userID] = rec
}

func (s *resumeStore) get(userID string) (*resumeRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.resumes[userID]
	return rec, ok
}

// ─────────────────────────────────────────────────────────────────────────────
// ResumeHandler
// ─────────────────────────────────────────────────────────────────────────────

// ResumeHandler handles resume upload and parsing endpoints.
type ResumeHandler struct {
	parser *parse.ResumeParser
}

// NewResumeHandler creates a new ResumeHandler.
func NewResumeHandler() *ResumeHandler {
	return &ResumeHandler{
		parser: parse.New(),
	}
}

// RegisterRoutes registers resume routes on the mux.
//
//	POST /api/resume/upload  – upload and parse a resume (PDF or DOCX)
func (h *ResumeHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/api/resume/upload",
		authMiddleware(http.HandlerFunc(h.Upload)))
}

// Upload handles POST /api/resume/upload.
//
// Accepts multipart/form-data with a "resume" file field.
// Parses the resume and stores the extracted data.
// Also updates the user's profile skills from the parsed resume.
//
// Response:
//
//	{
//	  "success": true,
//	  "data": {
//	    "resume_id": "...",
//	    "file_name": "resume.pdf",
//	    "parsed_at": "...",
//	    "personal": {...},
//	    "skills": [...],
//	    "experience": [...],
//	    "education": [...],
//	    "certifications": [...],
//	    "projects": [...]
//	  }
//	}
func (h *ResumeHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	// Parse multipart form (max 10MB).
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_FORM",
			"failed to parse multipart form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("resume")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "MISSING_FILE",
			"resume file is required (field name: 'resume')")
		return
	}
	defer file.Close()

	// Validate file type.
	fileName := header.Filename
	fileType := detectFileType(fileName)
	if fileType == "" {
		WriteError(w, http.StatusBadRequest, "UNSUPPORTED_FORMAT",
			"only PDF and DOCX files are supported")
		return
	}

	// Read file content.
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "READ_ERROR",
			"failed to read uploaded file")
		return
	}

	// Parse the resume.
	req := parse.ParseRequest{
		FileName:    fileName,
		FileContent: fileBytes,
		FileType:    fileType,
	}
	result, parseErr := h.parser.Parse(req)
	if parseErr != nil {
		WriteError(w, http.StatusUnprocessableEntity, "PARSE_ERROR",
			"failed to parse resume: "+parseErr.Error())
		return
	}

	// Store the parsed resume.
	rec := &resumeRecord{
		ID:         generateID(),
		UserID:     userID,
		FileName:   fileName,
		ParsedAt:   time.Now(),
		ParsedData: result,
	}
	globalResumeStore.save(userID, rec)

	// Auto-update profile skills from parsed resume.
	if len(result.Skills) > 0 {
		globalProfileStore.update(userID, func(p *profileRecord) {
			p.Skills = make([]skillRecord, len(result.Skills))
			for i, s := range result.Skills {
				p.Skills[i] = skillRecord{
					Name:        s.Name,
					Proficiency: s.Category, // use category as proficiency proxy
				}
			}
		})
	}

	WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"resume_id":      rec.ID,
		"file_name":      fileName,
		"parsed_at":      rec.ParsedAt,
		"personal":       result.PersonalInfo,
		"skills":         result.Skills,
		"experience":     result.WorkExperience,
		"education":      result.Education,
		"certifications": result.Certifications,
		"projects":       result.Projects,
		"summary": fmt.Sprintf("Parsed %d skills, %d experience entries, %d education entries",
			len(result.Skills), len(result.WorkExperience), len(result.Education)),
	})
}

// detectFileType returns "pdf" or "docx" based on the file extension.
func detectFileType(fileName string) string {
	lower := strings.ToLower(fileName)
	if strings.HasSuffix(lower, ".pdf") {
		return "pdf"
	}
	if strings.HasSuffix(lower, ".docx") {
		return "docx"
	}
	return ""
}
