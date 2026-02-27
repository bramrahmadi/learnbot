package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// UserRepository provides CRUD operations for users and profiles.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// User CRUD
// ─────────────────────────────────────────────────────────────────────────────

// CreateUser inserts a new user and initializes their profile and preferences.
// Uses a transaction to ensure atomicity.
func (r *UserRepository) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Normalize email
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	timezone := input.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	locale := input.Locale
	if locale == "" {
		locale = "en"
	}

	user := &User{}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (email, full_name, password_hash, avatar_url, timezone, locale)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, email_verified, full_name, avatar_url, timezone, locale,
		          is_active, is_admin, last_login_at, created_at, updated_at`,
		email, input.FullName, input.PasswordHash, input.AvatarURL, timezone, locale,
	).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.FullName,
		&user.AvatarURL, &user.Timezone, &user.Locale,
		&user.IsActive, &user.IsAdmin, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	// Initialize empty profile
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_profiles (user_id) VALUES ($1)`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}

	// Initialize empty preferences
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_preferences (user_id) VALUES ($1)`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("create preferences: %w", err)
	}

	// Record creation event
	_, err = tx.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type)
		VALUES ($1, $2)`, user.ID, EventCreated)
	if err != nil {
		return nil, fmt.Errorf("record history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their UUID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, email_verified, full_name, avatar_url, timezone, locale,
		       is_active, is_admin, last_login_at, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = TRUE`, id,
	).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.FullName,
		&user.AvatarURL, &user.Timezone, &user.Locale,
		&user.IsActive, &user.IsAdmin, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, email_verified, password_hash, full_name, avatar_url,
		       timezone, locale, is_active, is_admin, last_login_at, created_at, updated_at
		FROM users WHERE email = $1 AND is_active = TRUE`,
		strings.ToLower(strings.TrimSpace(email)),
	).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.FullName, &user.AvatarURL, &user.Timezone, &user.Locale,
		&user.IsActive, &user.IsAdmin, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

// UpdateLastLogin updates the user's last login timestamp.
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET last_login_at = $1 WHERE id = $2`,
		time.Now(), userID,
	)
	return err
}

// DeactivateUser soft-deletes a user by setting is_active = false.
func (r *UserRepository) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE users SET is_active = FALSE WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Profile CRUD
// ─────────────────────────────────────────────────────────────────────────────

// GetProfile retrieves a user's profile by user ID.
func (r *UserRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	p := &UserProfile{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, headline, summary, location_city, location_state,
		       location_country, phone, linkedin_url, github_url, website_url,
		       years_of_experience, is_open_to_work, profile_completeness,
		       created_at, updated_at
		FROM user_profiles WHERE user_id = $1`, userID,
	).Scan(
		&p.ID, &p.UserID, &p.Headline, &p.Summary,
		&p.LocationCity, &p.LocationState, &p.LocationCountry,
		&p.Phone, &p.LinkedInURL, &p.GitHubURL, &p.WebsiteURL,
		&p.YearsOfExperience, &p.IsOpenToWork, &p.ProfileCompleteness,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return p, nil
}

// UpdateProfile updates a user's profile fields and recalculates completeness.
func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput, changedBy *uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Capture old data for history
	var oldData []byte
	tx.QueryRowContext(ctx, `
		SELECT row_to_json(p) FROM user_profiles p WHERE user_id = $1`, userID,
	).Scan(&oldData)

	// Build dynamic UPDATE query
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if input.Headline != nil {
		setClauses = append(setClauses, fmt.Sprintf("headline = $%d", argIdx))
		args = append(args, *input.Headline)
		argIdx++
	}
	if input.Summary != nil {
		setClauses = append(setClauses, fmt.Sprintf("summary = $%d", argIdx))
		args = append(args, *input.Summary)
		argIdx++
	}
	if input.LocationCity != nil {
		setClauses = append(setClauses, fmt.Sprintf("location_city = $%d", argIdx))
		args = append(args, *input.LocationCity)
		argIdx++
	}
	if input.LocationState != nil {
		setClauses = append(setClauses, fmt.Sprintf("location_state = $%d", argIdx))
		args = append(args, *input.LocationState)
		argIdx++
	}
	if input.LocationCountry != nil {
		setClauses = append(setClauses, fmt.Sprintf("location_country = $%d", argIdx))
		args = append(args, *input.LocationCountry)
		argIdx++
	}
	if input.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argIdx))
		args = append(args, *input.Phone)
		argIdx++
	}
	if input.LinkedInURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("linkedin_url = $%d", argIdx))
		args = append(args, *input.LinkedInURL)
		argIdx++
	}
	if input.GitHubURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("github_url = $%d", argIdx))
		args = append(args, *input.GitHubURL)
		argIdx++
	}
	if input.WebsiteURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("website_url = $%d", argIdx))
		args = append(args, *input.WebsiteURL)
		argIdx++
	}
	if input.YearsOfExperience != nil {
		setClauses = append(setClauses, fmt.Sprintf("years_of_experience = $%d", argIdx))
		args = append(args, *input.YearsOfExperience)
		argIdx++
	}
	if input.IsOpenToWork != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_open_to_work = $%d", argIdx))
		args = append(args, *input.IsOpenToWork)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil // nothing to update
	}

	args = append(args, userID)
	query := fmt.Sprintf(
		"UPDATE user_profiles SET %s WHERE user_id = $%d",
		strings.Join(setClauses, ", "), argIdx,
	)

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update profile: %w", err)
	}

	// Recalculate completeness
	if _, err := tx.ExecContext(ctx, `
		UPDATE user_profiles
		SET profile_completeness = calculate_profile_completeness($1)
		WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("recalculate completeness: %w", err)
	}

	// Record history
	var newData []byte
	tx.QueryRowContext(ctx, `
		SELECT row_to_json(p) FROM user_profiles p WHERE user_id = $1`, userID,
	).Scan(&newData)

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type, old_data, new_data, changed_by)
		VALUES ($1, $2, 'profile', $3, $4, $5)`,
		userID, EventManuallyUpdated, oldData, newData, changedBy,
	); err != nil {
		return fmt.Errorf("record history: %w", err)
	}

	return tx.Commit()
}

// GetProfileHistory returns the audit log for a user's profile.
func (r *UserRepository) GetProfileHistory(ctx context.Context, userID uuid.UUID, limit int) ([]ProfileHistory, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, event_type, entity_type, entity_id,
		       old_data, new_data, changed_by, ip_address, user_agent, created_at
		FROM profile_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get profile history: %w", err)
	}
	defer rows.Close()

	var history []ProfileHistory
	for rows.Next() {
		var h ProfileHistory
		if err := rows.Scan(
			&h.ID, &h.UserID, &h.EventType, &h.EntityType, &h.EntityID,
			&h.OldData, &h.NewData, &h.ChangedBy, &h.IPAddress, &h.UserAgent,
			&h.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan history row: %w", err)
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Sentinel errors
// ─────────────────────────────────────────────────────────────────────────────

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = fmt.Errorf("record not found")

// ErrDuplicate is returned when a unique constraint is violated.
var ErrDuplicate = fmt.Errorf("record already exists")
