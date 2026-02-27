package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// SkillRepository provides CRUD operations for user skills.
type SkillRepository struct {
	db *sql.DB
}

// NewSkillRepository creates a new SkillRepository.
func NewSkillRepository(db *sql.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

// UpsertSkill creates or updates a skill for a user.
// If a skill with the same normalized name already exists, it is updated.
func (r *SkillRepository) UpsertSkill(ctx context.Context, userID uuid.UUID, input UpsertSkillInput) (*UserSkill, error) {
	normalized := strings.ToLower(strings.TrimSpace(input.SkillName))
	if normalized == "" {
		return nil, fmt.Errorf("skill name is required")
	}

	skill := &UserSkill{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO user_skills (
			user_id, skill_name, normalized_name, category, proficiency,
			years_of_experience, is_primary, source, confidence, last_used_year
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id, normalized_name) DO UPDATE SET
			skill_name         = EXCLUDED.skill_name,
			category           = EXCLUDED.category,
			proficiency        = EXCLUDED.proficiency,
			years_of_experience = EXCLUDED.years_of_experience,
			is_primary         = EXCLUDED.is_primary,
			source             = EXCLUDED.source,
			confidence         = EXCLUDED.confidence,
			last_used_year     = EXCLUDED.last_used_year,
			updated_at         = NOW()
		RETURNING id, user_id, skill_taxonomy_id, skill_name, normalized_name,
		          category, proficiency, years_of_experience, is_primary, source,
		          confidence, last_used_year, created_at, updated_at`,
		userID, input.SkillName, normalized, input.Category, input.Proficiency,
		input.YearsOfExperience, input.IsPrimary, input.Source,
		input.Confidence, input.LastUsedYear,
	).Scan(
		&skill.ID, &skill.UserID, &skill.SkillTaxonomyID,
		&skill.SkillName, &skill.NormalizedName, &skill.Category,
		&skill.Proficiency, &skill.YearsOfExperience, &skill.IsPrimary,
		&skill.Source, &skill.Confidence, &skill.LastUsedYear,
		&skill.CreatedAt, &skill.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert skill: %w", err)
	}

	// Record history
	newData, _ := json.Marshal(skill)
	r.db.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type, entity_id, new_data)
		VALUES ($1, $2, 'skill', $3, $4)`,
		userID, EventSkillAdded, skill.ID, newData,
	)

	return skill, nil
}

// GetSkillsByUserID retrieves all skills for a user, ordered by proficiency and name.
func (r *SkillRepository) GetSkillsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSkill, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, skill_taxonomy_id, skill_name, normalized_name,
		       category, proficiency, years_of_experience, is_primary, source,
		       confidence, last_used_year, created_at, updated_at
		FROM user_skills
		WHERE user_id = $1
		ORDER BY is_primary DESC,
		         CASE proficiency
		             WHEN 'expert' THEN 1
		             WHEN 'advanced' THEN 2
		             WHEN 'intermediate' THEN 3
		             ELSE 4
		         END,
		         skill_name ASC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get skills: %w", err)
	}
	defer rows.Close()

	var skills []UserSkill
	for rows.Next() {
		var s UserSkill
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.SkillTaxonomyID,
			&s.SkillName, &s.NormalizedName, &s.Category,
			&s.Proficiency, &s.YearsOfExperience, &s.IsPrimary,
			&s.Source, &s.Confidence, &s.LastUsedYear,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		skills = append(skills, s)
	}
	return skills, rows.Err()
}

// GetSkillsByCategory retrieves skills for a user filtered by category.
func (r *SkillRepository) GetSkillsByCategory(ctx context.Context, userID uuid.UUID, category SkillCategory) ([]UserSkill, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, skill_taxonomy_id, skill_name, normalized_name,
		       category, proficiency, years_of_experience, is_primary, source,
		       confidence, last_used_year, created_at, updated_at
		FROM user_skills
		WHERE user_id = $1 AND category = $2
		ORDER BY proficiency DESC, skill_name ASC`, userID, category,
	)
	if err != nil {
		return nil, fmt.Errorf("get skills by category: %w", err)
	}
	defer rows.Close()

	var skills []UserSkill
	for rows.Next() {
		var s UserSkill
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.SkillTaxonomyID,
			&s.SkillName, &s.NormalizedName, &s.Category,
			&s.Proficiency, &s.YearsOfExperience, &s.IsPrimary,
			&s.Source, &s.Confidence, &s.LastUsedYear,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		skills = append(skills, s)
	}
	return skills, rows.Err()
}

// DeleteSkill removes a skill from a user's profile.
func (r *SkillRepository) DeleteSkill(ctx context.Context, userID, skillID uuid.UUID) error {
	// Capture for history
	var oldData []byte
	r.db.QueryRowContext(ctx, `
		SELECT row_to_json(s) FROM user_skills s WHERE id = $1 AND user_id = $2`,
		skillID, userID,
	).Scan(&oldData)

	result, err := r.db.ExecContext(ctx,
		`DELETE FROM user_skills WHERE id = $1 AND user_id = $2`, skillID, userID)
	if err != nil {
		return fmt.Errorf("delete skill: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	// Record history
	r.db.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type, entity_id, old_data)
		VALUES ($1, $2, 'skill', $3, $4)`,
		userID, EventSkillRemoved, skillID, oldData,
	)

	return nil
}

// BulkUpsertSkills creates or updates multiple skills in a single transaction.
// Used by the resume parser to import all skills at once.
func (r *SkillRepository) BulkUpsertSkills(ctx context.Context, userID uuid.UUID, inputs []UpsertSkillInput) ([]UserSkill, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var skills []UserSkill
	for _, input := range inputs {
		normalized := strings.ToLower(strings.TrimSpace(input.SkillName))
		if normalized == "" {
			continue
		}

		skill := &UserSkill{}
		err := tx.QueryRowContext(ctx, `
			INSERT INTO user_skills (
				user_id, skill_name, normalized_name, category, proficiency,
				years_of_experience, is_primary, source, confidence, last_used_year
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (user_id, normalized_name) DO UPDATE SET
				proficiency        = EXCLUDED.proficiency,
				source             = EXCLUDED.source,
				confidence         = EXCLUDED.confidence,
				updated_at         = NOW()
			RETURNING id, user_id, skill_taxonomy_id, skill_name, normalized_name,
			          category, proficiency, years_of_experience, is_primary, source,
			          confidence, last_used_year, created_at, updated_at`,
			userID, input.SkillName, normalized, input.Category, input.Proficiency,
			input.YearsOfExperience, input.IsPrimary, input.Source,
			input.Confidence, input.LastUsedYear,
		).Scan(
			&skill.ID, &skill.UserID, &skill.SkillTaxonomyID,
			&skill.SkillName, &skill.NormalizedName, &skill.Category,
			&skill.Proficiency, &skill.YearsOfExperience, &skill.IsPrimary,
			&skill.Source, &skill.Confidence, &skill.LastUsedYear,
			&skill.CreatedAt, &skill.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("upsert skill %q: %w", input.SkillName, err)
		}
		skills = append(skills, *skill)
	}

	// Recalculate profile completeness
	tx.ExecContext(ctx, `
		UPDATE user_profiles
		SET profile_completeness = calculate_profile_completeness($1)
		WHERE user_id = $1`, userID)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return skills, nil
}

// SearchSkillsByName searches for skills matching a name pattern (for autocomplete).
func (r *SkillRepository) SearchSkillsByName(ctx context.Context, query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT normalized_name
		FROM skill_taxonomy
		WHERE normalized_name ILIKE $1
		   OR $1 = ANY(aliases)
		ORDER BY normalized_name
		LIMIT $2`,
		"%"+strings.ToLower(query)+"%", limit,
	)
	if err != nil {
		return nil, fmt.Errorf("search skills: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// ensure pq is used (for array support)
var _ = pq.StringArray{}
