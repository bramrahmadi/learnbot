// Package handler – auth.go implements authentication endpoints.
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/learnbot/api-gateway/internal/middleware"
	"github.com/learnbot/api-gateway/internal/types"
)

// ─────────────────────────────────────────────────────────────────────────────
// In-memory user store (MVP – replace with database in production)
// ─────────────────────────────────────────────────────────────────────────────

// userRecord stores a user in the in-memory store.
type userRecord struct {
	ID           string
	Email        string
	PasswordHash string // bcrypt hash (simplified: SHA-256 hex for MVP)
	FullName     string
	IsAdmin      bool
	CreatedAt    time.Time
}

// userStore is a thread-safe in-memory user store.
type userStore struct {
	mu    sync.RWMutex
	users map[string]*userRecord // keyed by email (lowercase)
	byID  map[string]*userRecord // keyed by ID
}

var globalUserStore = &userStore{
	users: make(map[string]*userRecord),
	byID:  make(map[string]*userRecord),
}

func (s *userStore) create(email, passwordHash, fullName string) *userRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := generateID()
	u := &userRecord{
		ID:           id,
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}
	s.users[u.Email] = u
	s.byID[id] = u
	return u
}

func (s *userStore) findByEmail(email string) (*userRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[strings.ToLower(strings.TrimSpace(email))]
	return u, ok
}

func (s *userStore) findByID(id string) (*userRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.byID[id]
	return u, ok
}

// generateID generates a random hex ID.
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// hashPassword creates a simple hash for the MVP.
// In production, use bcrypt.
func hashPassword(password string) string {
	// Simple deterministic hash for MVP testing.
	// Production: use golang.org/x/crypto/bcrypt
	h := 0
	for _, c := range password {
		h = h*31 + int(c)
	}
	return hex.EncodeToString([]byte{
		byte(h >> 24), byte(h >> 16), byte(h >> 8), byte(h),
		byte(len(password)), byte(h ^ len(password)),
	})
}

func checkPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

// ─────────────────────────────────────────────────────────────────────────────
// AuthHandler
// ─────────────────────────────────────────────────────────────────────────────

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	jwtCfg middleware.JWTConfig
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(jwtCfg middleware.JWTConfig) *AuthHandler {
	return &AuthHandler{jwtCfg: jwtCfg}
}

// RegisterRoutes registers auth routes on the mux.
//
//	POST /api/auth/register  – register a new user
//	POST /api/auth/login     – login and get JWT token
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/register", h.Register)
	mux.HandleFunc("/api/auth/login", h.Login)
}

// Register handles POST /api/auth/register.
//
// Request body:
//
//	{"email": "user@example.com", "password": "secret123", "full_name": "Jane Doe"}
//
// Response:
//
//	{"success": true, "data": {"token": "...", "expires_at": "...", "user": {...}}}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var req types.RegisterRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	// Validate input.
	var v Validator
	v.Required("email", req.Email, "email is required")
	v.ValidEmail("email", req.Email)
	v.Required("password", req.Password, "password is required")
	v.MinLength("password", req.Password, 8, "password must be at least 8 characters")
	v.Required("full_name", req.FullName, "full_name is required")
	if v.WriteIfInvalid(w) {
		return
	}

	// Check if email already exists.
	if _, exists := globalUserStore.findByEmail(req.Email); exists {
		WriteError(w, http.StatusConflict, "EMAIL_TAKEN",
			"an account with this email already exists")
		return
	}

	// Create user.
	hash := hashPassword(req.Password)
	user := globalUserStore.create(req.Email, hash, req.FullName)

	// Generate token.
	token, expiresAt, err := middleware.GenerateToken(h.jwtCfg, user.ID, user.Email, user.IsAdmin)
	if err != nil {
		WriteInternalError(w)
		return
	}

	WriteSuccess(w, http.StatusCreated, types.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: types.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	})
}

// Login handles POST /api/auth/login.
//
// Request body:
//
//	{"email": "user@example.com", "password": "secret123"}
//
// Response:
//
//	{"success": true, "data": {"token": "...", "expires_at": "...", "user": {...}}}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var req types.LoginRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	var v Validator
	v.Required("email", req.Email, "email is required")
	v.Required("password", req.Password, "password is required")
	if v.WriteIfInvalid(w) {
		return
	}

	// Find user.
	user, exists := globalUserStore.findByEmail(req.Email)
	if !exists || !checkPassword(req.Password, user.PasswordHash) {
		WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS",
			"invalid email or password")
		return
	}

	// Generate token.
	token, expiresAt, err := middleware.GenerateToken(h.jwtCfg, user.ID, user.Email, user.IsAdmin)
	if err != nil {
		WriteInternalError(w)
		return
	}

	WriteSuccess(w, http.StatusOK, types.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: types.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	})
}
