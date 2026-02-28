// Package repository provides data access for the learning resource catalog.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LearningResourceRepository provides CRUD operations for learning resources.
type LearningResourceRepository struct {
	db *sql.DB
}

// NewLearningResourceRepository creates a new LearningResourceRepository.
func NewLearningResourceRepository(db *sql.DB) *LearningResourceRepository {
	return &LearningResourceRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// Resource queries
// ─────────────────────────────────────────────────────────────────────────────

// GetByID returns a learning resource by its UUID.
func (r *LearningResourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*LearningResource, error) {
	const q = `
		SELECT id, title, slug, description, url, provider_id, resource_type,
		       difficulty, cost_type, cost_amount, cost_currency, duration_hours,
		       duration_label, language, is_active, is_featured, is_verified,
		       has_certificate, has_hands_on, rating, rating_count, enrollment_count,
		       last_updated_date, created_at, updated_at
		FROM learning_resources
		WHERE id = $1 AND is_active = TRUE`

	var res LearningResource
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.URL,
		&res.ProviderID, &res.ResourceType, &res.Difficulty, &res.CostType,
		&res.CostAmount, &res.CostCurrency, &res.DurationHours, &res.DurationLabel,
		&res.Language, &res.IsActive, &res.IsFeatured, &res.IsVerified,
		&res.HasCertificate, &res.HasHandsOn, &res.Rating, &res.RatingCount,
		&res.EnrollmentCount, &res.LastUpdatedDate, &res.CreatedAt, &res.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resource by id: %w", err)
	}
	return &res, nil
}

// GetBySlug returns a learning resource by its slug.
func (r *LearningResourceRepository) GetBySlug(ctx context.Context, slug string) (*LearningResource, error) {
	const q = `
		SELECT id, title, slug, description, url, provider_id, resource_type,
		       difficulty, cost_type, cost_amount, cost_currency, duration_hours,
		       duration_label, language, is_active, is_featured, is_verified,
		       has_certificate, has_hands_on, rating, rating_count, enrollment_count,
		       last_updated_date, created_at, updated_at
		FROM learning_resources
		WHERE slug = $1 AND is_active = TRUE`

	var res LearningResource
	err := r.db.QueryRowContext(ctx, q, slug).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.URL,
		&res.ProviderID, &res.ResourceType, &res.Difficulty, &res.CostType,
		&res.CostAmount, &res.CostCurrency, &res.DurationHours, &res.DurationLabel,
		&res.Language, &res.IsActive, &res.IsFeatured, &res.IsVerified,
		&res.HasCertificate, &res.HasHandsOn, &res.Rating, &res.RatingCount,
		&res.EnrollmentCount, &res.LastUpdatedDate, &res.CreatedAt, &res.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resource by slug: %w", err)
	}
	return &res, nil
}

// List returns resources matching the given filter.
func (r *LearningResourceRepository) List(ctx context.Context, filter ResourceQueryFilter) ([]LearningResourceWithSkills, int, error) {
	// Normalize limit.
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// Build WHERE clauses dynamically.
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, "lr.is_active = TRUE")

	if filter.SkillName != "" {
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM resource_skills rs2 WHERE rs2.resource_id = lr.id AND rs2.normalized_name = $%d)", argIdx))
		args = append(args, strings.ToLower(strings.TrimSpace(filter.SkillName)))
		argIdx++
	}
	if filter.ResourceType != "" {
		conditions = append(conditions, fmt.Sprintf("lr.resource_type = $%d", argIdx))
		args = append(args, string(filter.ResourceType))
		argIdx++
	}
	if filter.Difficulty != "" {
		conditions = append(conditions, fmt.Sprintf("lr.difficulty = $%d", argIdx))
		args = append(args, string(filter.Difficulty))
		argIdx++
	}
	if filter.CostType != "" {
		conditions = append(conditions, fmt.Sprintf("lr.cost_type = $%d", argIdx))
		args = append(args, string(filter.CostType))
		argIdx++
	}
	if filter.IsFree {
		conditions = append(conditions, "lr.cost_type IN ('free', 'free_audit')")
	}
	if filter.HasCertificate {
		conditions = append(conditions, "lr.has_certificate = TRUE")
	}
	if filter.HasHandsOn {
		conditions = append(conditions, "lr.has_hands_on = TRUE")
	}
	if filter.MinRating > 0 {
		conditions = append(conditions, fmt.Sprintf("lr.rating >= $%d", argIdx))
		args = append(args, filter.MinRating)
		argIdx++
	}
	if filter.ProviderID != nil {
		conditions = append(conditions, fmt.Sprintf("lr.provider_id = $%d", argIdx))
		args = append(args, *filter.ProviderID)
		argIdx++
	}
	if filter.SearchQuery != "" {
		conditions = append(conditions, fmt.Sprintf(
			"to_tsvector('english', lr.title || ' ' || COALESCE(lr.description, '')) @@ plainto_tsquery('english', $%d)", argIdx))
		args = append(args, filter.SearchQuery)
		argIdx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	// Count query.
	countQ := fmt.Sprintf(`
		SELECT COUNT(DISTINCT lr.id)
		FROM learning_resources lr
		%s`, where)

	var total int
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count resources: %w", err)
	}

	// Data query using the view.
	dataQ := fmt.Sprintf(`
		SELECT
			lr.id, lr.title, lr.slug, lr.description, lr.url, lr.provider_id,
			lr.resource_type, lr.difficulty, lr.cost_type, lr.cost_amount,
			lr.cost_currency, lr.duration_hours, lr.duration_label, lr.language,
			lr.is_active, lr.is_featured, lr.is_verified, lr.has_certificate,
			lr.has_hands_on, lr.rating, lr.rating_count, lr.enrollment_count,
			lr.last_updated_date, lr.created_at, lr.updated_at,
			COALESCE(rp.name, '') AS provider_name,
			rp.website_url AS provider_url,
			ARRAY_AGG(DISTINCT rs.skill_name ORDER BY rs.skill_name) FILTER (WHERE rs.skill_name IS NOT NULL) AS skills,
			ARRAY_AGG(DISTINCT rs.normalized_name ORDER BY rs.normalized_name) FILTER (WHERE rs.normalized_name IS NOT NULL) AS skill_ids
		FROM learning_resources lr
		LEFT JOIN resource_providers rp ON lr.provider_id = rp.id
		LEFT JOIN resource_skills rs ON lr.id = rs.resource_id
		%s
		GROUP BY lr.id, rp.name, rp.website_url
		ORDER BY lr.is_featured DESC, lr.rating DESC NULLS LAST, lr.rating_count DESC
		LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list resources: %w", err)
	}
	defer rows.Close()

	var resources []LearningResourceWithSkills
	for rows.Next() {
		var res LearningResourceWithSkills
		if err := rows.Scan(
			&res.ID, &res.Title, &res.Slug, &res.Description, &res.URL,
			&res.ProviderID, &res.ResourceType, &res.Difficulty, &res.CostType,
			&res.CostAmount, &res.CostCurrency, &res.DurationHours, &res.DurationLabel,
			&res.Language, &res.IsActive, &res.IsFeatured, &res.IsVerified,
			&res.HasCertificate, &res.HasHandsOn, &res.Rating, &res.RatingCount,
			&res.EnrollmentCount, &res.LastUpdatedDate, &res.CreatedAt, &res.UpdatedAt,
			&res.ProviderName, &res.ProviderURL, &res.Skills, &res.SkillIDs,
		); err != nil {
			return nil, 0, fmt.Errorf("scan resource: %w", err)
		}
		resources = append(resources, res)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate resources: %w", err)
	}

	return resources, total, nil
}

// GetBySkill returns resources that cover a specific skill, ordered by rating.
func (r *LearningResourceRepository) GetBySkill(ctx context.Context, skillName string, limit int) ([]LearningResourceWithSkills, error) {
	if limit <= 0 {
		limit = 10
	}
	filter := ResourceQueryFilter{
		SkillName: skillName,
		Limit:     limit,
	}
	resources, _, err := r.List(ctx, filter)
	return resources, err
}

// GetFeatured returns featured resources.
func (r *LearningResourceRepository) GetFeatured(ctx context.Context, limit int) ([]LearningResourceWithSkills, error) {
	if limit <= 0 {
		limit = 10
	}
	const q = `
		SELECT
			lr.id, lr.title, lr.slug, lr.description, lr.url, lr.provider_id,
			lr.resource_type, lr.difficulty, lr.cost_type, lr.cost_amount,
			lr.cost_currency, lr.duration_hours, lr.duration_label, lr.language,
			lr.is_active, lr.is_featured, lr.is_verified, lr.has_certificate,
			lr.has_hands_on, lr.rating, lr.rating_count, lr.enrollment_count,
			lr.last_updated_date, lr.created_at, lr.updated_at,
			COALESCE(rp.name, '') AS provider_name,
			rp.website_url AS provider_url,
			ARRAY_AGG(DISTINCT rs.skill_name ORDER BY rs.skill_name) FILTER (WHERE rs.skill_name IS NOT NULL) AS skills,
			ARRAY_AGG(DISTINCT rs.normalized_name ORDER BY rs.normalized_name) FILTER (WHERE rs.normalized_name IS NOT NULL) AS skill_ids
		FROM learning_resources lr
		LEFT JOIN resource_providers rp ON lr.provider_id = rp.id
		LEFT JOIN resource_skills rs ON lr.id = rs.resource_id
		WHERE lr.is_active = TRUE AND lr.is_featured = TRUE
		GROUP BY lr.id, rp.name, rp.website_url
		ORDER BY lr.rating DESC NULLS LAST, lr.rating_count DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("get featured resources: %w", err)
	}
	defer rows.Close()

	var resources []LearningResourceWithSkills
	for rows.Next() {
		var res LearningResourceWithSkills
		if err := rows.Scan(
			&res.ID, &res.Title, &res.Slug, &res.Description, &res.URL,
			&res.ProviderID, &res.ResourceType, &res.Difficulty, &res.CostType,
			&res.CostAmount, &res.CostCurrency, &res.DurationHours, &res.DurationLabel,
			&res.Language, &res.IsActive, &res.IsFeatured, &res.IsVerified,
			&res.HasCertificate, &res.HasHandsOn, &res.Rating, &res.RatingCount,
			&res.EnrollmentCount, &res.LastUpdatedDate, &res.CreatedAt, &res.UpdatedAt,
			&res.ProviderName, &res.ProviderURL, &res.Skills, &res.SkillIDs,
		); err != nil {
			return nil, fmt.Errorf("scan featured resource: %w", err)
		}
		resources = append(resources, res)
	}
	return resources, rows.Err()
}

// Create inserts a new learning resource and its skills.
func (r *LearningResourceRepository) Create(ctx context.Context, input CreateResourceInput) (*LearningResource, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const q = `
		INSERT INTO learning_resources (
			title, slug, description, url, provider_id, resource_type,
			difficulty, cost_type, cost_amount, cost_currency, duration_hours,
			duration_label, language, has_certificate, has_hands_on
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id, title, slug, description, url, provider_id, resource_type,
		          difficulty, cost_type, cost_amount, cost_currency, duration_hours,
		          duration_label, language, is_active, is_featured, is_verified,
		          has_certificate, has_hands_on, rating, rating_count, enrollment_count,
		          last_updated_date, created_at, updated_at`

	currency := input.CostCurrency
	if currency == "" {
		currency = "USD"
	}
	lang := input.Language
	if lang == "" {
		lang = "en"
	}

	var res LearningResource
	err = tx.QueryRowContext(ctx, q,
		input.Title, input.Slug, input.Description, input.URL, input.ProviderID,
		string(input.ResourceType), string(input.Difficulty), string(input.CostType),
		input.CostAmount, currency, input.DurationHours, input.DurationLabel,
		lang, input.HasCertificate, input.HasHandsOn,
	).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.URL,
		&res.ProviderID, &res.ResourceType, &res.Difficulty, &res.CostType,
		&res.CostAmount, &res.CostCurrency, &res.DurationHours, &res.DurationLabel,
		&res.Language, &res.IsActive, &res.IsFeatured, &res.IsVerified,
		&res.HasCertificate, &res.HasHandsOn, &res.Rating, &res.RatingCount,
		&res.EnrollmentCount, &res.LastUpdatedDate, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert resource: %w", err)
	}

	// Insert skills.
	for _, skill := range input.Skills {
		norm := strings.ToLower(strings.TrimSpace(skill.SkillName))
		if norm == "" {
			continue
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO resource_skills (resource_id, skill_name, normalized_name, is_primary, coverage_level)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (resource_id, normalized_name) DO NOTHING`,
			res.ID, skill.SkillName, norm, skill.IsPrimary, skill.CoverageLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("insert resource skill %q: %w", skill.SkillName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return &res, nil
}

// Update modifies an existing learning resource.
func (r *LearningResourceRepository) Update(ctx context.Context, id uuid.UUID, input UpdateResourceInput) (*LearningResource, error) {
	var setClauses []string
	var args []interface{}
	argIdx := 1

	if input.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *input.Title)
		argIdx++
	}
	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *input.Description)
		argIdx++
	}
	if input.URL != nil {
		setClauses = append(setClauses, fmt.Sprintf("url = $%d", argIdx))
		args = append(args, *input.URL)
		argIdx++
	}
	if input.Difficulty != nil {
		setClauses = append(setClauses, fmt.Sprintf("difficulty = $%d", argIdx))
		args = append(args, string(*input.Difficulty))
		argIdx++
	}
	if input.CostType != nil {
		setClauses = append(setClauses, fmt.Sprintf("cost_type = $%d", argIdx))
		args = append(args, string(*input.CostType))
		argIdx++
	}
	if input.CostAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("cost_amount = $%d", argIdx))
		args = append(args, *input.CostAmount)
		argIdx++
	}
	if input.DurationHours != nil {
		setClauses = append(setClauses, fmt.Sprintf("duration_hours = $%d", argIdx))
		args = append(args, *input.DurationHours)
		argIdx++
	}
	if input.DurationLabel != nil {
		setClauses = append(setClauses, fmt.Sprintf("duration_label = $%d", argIdx))
		args = append(args, *input.DurationLabel)
		argIdx++
	}
	if input.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *input.IsActive)
		argIdx++
	}
	if input.IsFeatured != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_featured = $%d", argIdx))
		args = append(args, *input.IsFeatured)
		argIdx++
	}
	if input.HasCertificate != nil {
		setClauses = append(setClauses, fmt.Sprintf("has_certificate = $%d", argIdx))
		args = append(args, *input.HasCertificate)
		argIdx++
	}
	if input.HasHandsOn != nil {
		setClauses = append(setClauses, fmt.Sprintf("has_hands_on = $%d", argIdx))
		args = append(args, *input.HasHandsOn)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now())
	argIdx++

	args = append(args, id)
	q := fmt.Sprintf(`
		UPDATE learning_resources
		SET %s
		WHERE id = $%d
		RETURNING id, title, slug, description, url, provider_id, resource_type,
		          difficulty, cost_type, cost_amount, cost_currency, duration_hours,
		          duration_label, language, is_active, is_featured, is_verified,
		          has_certificate, has_hands_on, rating, rating_count, enrollment_count,
		          last_updated_date, created_at, updated_at`,
		strings.Join(setClauses, ", "), argIdx)

	var res LearningResource
	err := r.db.QueryRowContext(ctx, q, args...).Scan(
		&res.ID, &res.Title, &res.Slug, &res.Description, &res.URL,
		&res.ProviderID, &res.ResourceType, &res.Difficulty, &res.CostType,
		&res.CostAmount, &res.CostCurrency, &res.DurationHours, &res.DurationLabel,
		&res.Language, &res.IsActive, &res.IsFeatured, &res.IsVerified,
		&res.HasCertificate, &res.HasHandsOn, &res.Rating, &res.RatingCount,
		&res.EnrollmentCount, &res.LastUpdatedDate, &res.CreatedAt, &res.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}
	return &res, nil
}

// Delete soft-deletes a resource by setting is_active = FALSE.
func (r *LearningResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE learning_resources SET is_active = FALSE, updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete resource: %w", err)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Learning path queries
// ─────────────────────────────────────────────────────────────────────────────

// GetPathBySlug returns a learning path by slug with its resources.
func (r *LearningResourceRepository) GetPathBySlug(ctx context.Context, slug string) (*LearningPathWithResources, error) {
	const q = `
		SELECT lp.id, lp.title, lp.slug, lp.description, lp.target_role,
		       lp.target_skill, lp.difficulty, lp.estimated_hours, lp.is_active,
		       lp.is_featured, lp.created_by, lp.created_at, lp.updated_at,
		       COUNT(lpr.id) AS resource_count,
		       COUNT(lpr.id) FILTER (WHERE lpr.is_required = TRUE) AS required_resource_count
		FROM learning_paths lp
		LEFT JOIN learning_path_resources lpr ON lp.id = lpr.path_id
		WHERE lp.slug = $1 AND lp.is_active = TRUE
		GROUP BY lp.id`

	var path LearningPathWithResources
	err := r.db.QueryRowContext(ctx, q, slug).Scan(
		&path.ID, &path.Title, &path.Slug, &path.Description, &path.TargetRole,
		&path.TargetSkill, &path.Difficulty, &path.EstimatedHours, &path.IsActive,
		&path.IsFeatured, &path.CreatedBy, &path.CreatedAt, &path.UpdatedAt,
		&path.ResourceCount, &path.RequiredResourceCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get path by slug: %w", err)
	}

	// Load path resources.
	steps, err := r.getPathResources(ctx, path.ID)
	if err != nil {
		return nil, err
	}
	path.Resources = steps

	return &path, nil
}

// ListPaths returns all active learning paths.
func (r *LearningResourceRepository) ListPaths(ctx context.Context, targetSkill, targetRole string) ([]LearningPathWithResources, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, "lp.is_active = TRUE")

	if targetSkill != "" {
		conditions = append(conditions, fmt.Sprintf("lp.target_skill ILIKE $%d", argIdx))
		args = append(args, "%"+targetSkill+"%")
		argIdx++
	}
	if targetRole != "" {
		conditions = append(conditions, fmt.Sprintf("lp.target_role ILIKE $%d", argIdx))
		args = append(args, "%"+targetRole+"%")
		argIdx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	q := fmt.Sprintf(`
		SELECT lp.id, lp.title, lp.slug, lp.description, lp.target_role,
		       lp.target_skill, lp.difficulty, lp.estimated_hours, lp.is_active,
		       lp.is_featured, lp.created_by, lp.created_at, lp.updated_at,
		       COUNT(lpr.id) AS resource_count,
		       COUNT(lpr.id) FILTER (WHERE lpr.is_required = TRUE) AS required_resource_count
		FROM learning_paths lp
		LEFT JOIN learning_path_resources lpr ON lp.id = lpr.path_id
		%s
		GROUP BY lp.id
		ORDER BY lp.is_featured DESC, lp.title`, where)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list paths: %w", err)
	}
	defer rows.Close()

	var paths []LearningPathWithResources
	for rows.Next() {
		var path LearningPathWithResources
		if err := rows.Scan(
			&path.ID, &path.Title, &path.Slug, &path.Description, &path.TargetRole,
			&path.TargetSkill, &path.Difficulty, &path.EstimatedHours, &path.IsActive,
			&path.IsFeatured, &path.CreatedBy, &path.CreatedAt, &path.UpdatedAt,
			&path.ResourceCount, &path.RequiredResourceCount,
		); err != nil {
			return nil, fmt.Errorf("scan path: %w", err)
		}
		paths = append(paths, path)
	}
	return paths, rows.Err()
}

// getPathResources returns the ordered resources for a learning path.
func (r *LearningResourceRepository) getPathResources(ctx context.Context, pathID uuid.UUID) ([]LearningPathStep, error) {
	const q = `
		SELECT lpr.step_order, lpr.is_required, lpr.notes,
		       lr.id, lr.title, lr.slug, lr.description, lr.url, lr.provider_id,
		       lr.resource_type, lr.difficulty, lr.cost_type, lr.cost_amount,
		       lr.cost_currency, lr.duration_hours, lr.duration_label, lr.language,
		       lr.is_active, lr.is_featured, lr.is_verified, lr.has_certificate,
		       lr.has_hands_on, lr.rating, lr.rating_count, lr.enrollment_count,
		       lr.last_updated_date, lr.created_at, lr.updated_at
		FROM learning_path_resources lpr
		JOIN learning_resources lr ON lpr.resource_id = lr.id
		WHERE lpr.path_id = $1
		ORDER BY lpr.step_order`

	rows, err := r.db.QueryContext(ctx, q, pathID)
	if err != nil {
		return nil, fmt.Errorf("get path resources: %w", err)
	}
	defer rows.Close()

	var steps []LearningPathStep
	for rows.Next() {
		var step LearningPathStep
		var notes sql.NullString
		if err := rows.Scan(
			&step.StepOrder, &step.IsRequired, &notes,
			&step.Resource.ID, &step.Resource.Title, &step.Resource.Slug,
			&step.Resource.Description, &step.Resource.URL, &step.Resource.ProviderID,
			&step.Resource.ResourceType, &step.Resource.Difficulty, &step.Resource.CostType,
			&step.Resource.CostAmount, &step.Resource.CostCurrency, &step.Resource.DurationHours,
			&step.Resource.DurationLabel, &step.Resource.Language, &step.Resource.IsActive,
			&step.Resource.IsFeatured, &step.Resource.IsVerified, &step.Resource.HasCertificate,
			&step.Resource.HasHandsOn, &step.Resource.Rating, &step.Resource.RatingCount,
			&step.Resource.EnrollmentCount, &step.Resource.LastUpdatedDate,
			&step.Resource.CreatedAt, &step.Resource.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan path step: %w", err)
		}
		if notes.Valid {
			step.Notes = notes.String
		}
		steps = append(steps, step)
	}
	return steps, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// User progress queries
// ─────────────────────────────────────────────────────────────────────────────

// GetUserProgress returns a user's progress for a specific resource.
func (r *LearningResourceRepository) GetUserProgress(ctx context.Context, userID, resourceID uuid.UUID) (*UserResourceProgress, error) {
	const q = `
		SELECT id, user_id, resource_id, status, progress_percentage,
		       started_at, completed_at, user_rating, user_notes, created_at, updated_at
		FROM user_resource_progress
		WHERE user_id = $1 AND resource_id = $2`

	var p UserResourceProgress
	err := r.db.QueryRowContext(ctx, q, userID, resourceID).Scan(
		&p.ID, &p.UserID, &p.ResourceID, &p.Status, &p.ProgressPercentage,
		&p.StartedAt, &p.CompletedAt, &p.UserRating, &p.UserNotes,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user progress: %w", err)
	}
	return &p, nil
}

// UpsertUserProgress creates or updates a user's progress for a resource.
func (r *LearningResourceRepository) UpsertUserProgress(ctx context.Context, userID, resourceID uuid.UUID, input UpsertProgressInput) (*UserResourceProgress, error) {
	now := time.Now()

	var startedAt *time.Time
	var completedAt *time.Time

	if input.Status == UserResourceStatusInProgress {
		startedAt = &now
	}
	if input.Status == UserResourceStatusCompleted {
		completedAt = &now
		pct := int16(100)
		input.ProgressPercentage = pct
	}

	const q = `
		INSERT INTO user_resource_progress (
			user_id, resource_id, status, progress_percentage,
			started_at, completed_at, user_rating, user_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, resource_id) DO UPDATE SET
			status = EXCLUDED.status,
			progress_percentage = EXCLUDED.progress_percentage,
			started_at = COALESCE(user_resource_progress.started_at, EXCLUDED.started_at),
			completed_at = EXCLUDED.completed_at,
			user_rating = COALESCE(EXCLUDED.user_rating, user_resource_progress.user_rating),
			user_notes = COALESCE(EXCLUDED.user_notes, user_resource_progress.user_notes),
			updated_at = NOW()
		RETURNING id, user_id, resource_id, status, progress_percentage,
		          started_at, completed_at, user_rating, user_notes, created_at, updated_at`

	var p UserResourceProgress
	err := r.db.QueryRowContext(ctx, q,
		userID, resourceID, string(input.Status), input.ProgressPercentage,
		startedAt, completedAt, input.UserRating, input.UserNotes,
	).Scan(
		&p.ID, &p.UserID, &p.ResourceID, &p.Status, &p.ProgressPercentage,
		&p.StartedAt, &p.CompletedAt, &p.UserRating, &p.UserNotes,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert user progress: %w", err)
	}
	return &p, nil
}

// ListUserProgress returns all resources a user has interacted with.
func (r *LearningResourceRepository) ListUserProgress(ctx context.Context, userID uuid.UUID, status UserResourceStatus) ([]UserResourceProgress, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
	args = append(args, userID)
	argIdx++

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, string(status))
		argIdx++
	}

	q := fmt.Sprintf(`
		SELECT id, user_id, resource_id, status, progress_percentage,
		       started_at, completed_at, user_rating, user_notes, created_at, updated_at
		FROM user_resource_progress
		WHERE %s
		ORDER BY updated_at DESC`, strings.Join(conditions, " AND "))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list user progress: %w", err)
	}
	defer rows.Close()

	var progress []UserResourceProgress
	for rows.Next() {
		var p UserResourceProgress
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.ResourceID, &p.Status, &p.ProgressPercentage,
			&p.StartedAt, &p.CompletedAt, &p.UserRating, &p.UserNotes,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user progress: %w", err)
		}
		progress = append(progress, p)
	}
	return progress, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Provider queries
// ─────────────────────────────────────────────────────────────────────────────

// ListProviders returns all active resource providers.
func (r *LearningResourceRepository) ListProviders(ctx context.Context) ([]ResourceProvider, error) {
	const q = `
		SELECT id, name, normalized_name, website_url, logo_url, description,
		       is_active, created_at, updated_at
		FROM resource_providers
		WHERE is_active = TRUE
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list providers: %w", err)
	}
	defer rows.Close()

	var providers []ResourceProvider
	for rows.Next() {
		var p ResourceProvider
		if err := rows.Scan(
			&p.ID, &p.Name, &p.NormalizedName, &p.WebsiteURL, &p.LogoURL,
			&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan provider: %w", err)
		}
		providers = append(providers, p)
	}
	return providers, rows.Err()
}

// CreateProvider inserts a new resource provider.
func (r *LearningResourceRepository) CreateProvider(ctx context.Context, input CreateProviderInput) (*ResourceProvider, error) {
	norm := strings.ToLower(strings.TrimSpace(input.Name))
	const q = `
		INSERT INTO resource_providers (name, normalized_name, website_url, logo_url, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, normalized_name, website_url, logo_url, description,
		          is_active, created_at, updated_at`

	var p ResourceProvider
	err := r.db.QueryRowContext(ctx, q,
		input.Name, norm, input.WebsiteURL, input.LogoURL, input.Description,
	).Scan(
		&p.ID, &p.Name, &p.NormalizedName, &p.WebsiteURL, &p.LogoURL,
		&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create provider: %w", err)
	}
	return &p, nil
}
