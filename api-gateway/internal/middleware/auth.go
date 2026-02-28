// Package middleware provides HTTP middleware for the API gateway.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	// ContextKeyUserID is the context key for the authenticated user's ID.
	ContextKeyUserID contextKey = "user_id"

	// ContextKeyEmail is the context key for the authenticated user's email.
	ContextKeyEmail contextKey = "email"

	// ContextKeyIsAdmin is the context key for the admin flag.
	ContextKeyIsAdmin contextKey = "is_admin"
)

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	// SecretKey is the HMAC secret used to sign tokens.
	SecretKey []byte

	// TokenDuration is how long tokens are valid.
	TokenDuration time.Duration
}

// DefaultJWTConfig returns a JWTConfig with sensible defaults.
// In production, SecretKey must be set from environment variables.
func DefaultJWTConfig(secret string) JWTConfig {
	if secret == "" {
		secret = "learnbot-dev-secret-change-in-production"
	}
	return JWTConfig{
		SecretKey:     []byte(secret),
		TokenDuration: 24 * time.Hour,
	}
}

// jwtClaims represents the JWT payload.
type jwtClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token for the given user.
func GenerateToken(cfg JWTConfig, userID, email string, isAdmin bool) (string, time.Time, error) {
	expiresAt := time.Now().Add(cfg.TokenDuration)
	claims := jwtClaims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "learnbot",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(cfg.SecretKey)
	return signed, expiresAt, err
}

// ParseToken validates and parses a JWT token string.
func ParseToken(cfg JWTConfig, tokenStr string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return cfg.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

// RequireAuth is a middleware that validates the JWT Bearer token.
// It sets user context values and calls next on success.
// On failure it returns 401 Unauthorized.
func RequireAuth(cfg JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				writeAuthError(w, "missing or invalid Authorization header")
				return
			}

			claims, err := ParseToken(cfg, tokenStr)
			if err != nil {
				writeAuthError(w, "invalid or expired token")
				return
			}

			// Inject claims into context.
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)
			ctx = context.WithValue(ctx, ContextKeyIsAdmin, claims.IsAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin is a middleware that requires the user to be an admin.
// Must be used after RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, _ := r.Context().Value(ContextKeyIsAdmin).(bool)
		if !isAdmin {
			writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts the authenticated user ID from the request context.
// Returns empty string if not authenticated.
func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(ContextKeyUserID).(string)
	return id
}

// GetEmail extracts the authenticated user email from the request context.
func GetEmail(r *http.Request) string {
	email, _ := r.Context().Value(ContextKeyEmail).(string)
	return email
}

// extractBearerToken extracts the token from the Authorization header.
// Expects format: "Bearer <token>"
func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// writeAuthError writes a 401 Unauthorized response.
func writeAuthError(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// writeJSONError writes a structured JSON error response.
func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
