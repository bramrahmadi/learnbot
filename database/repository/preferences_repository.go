package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PreferencesRepository provides CRUD operations for user preferences and career goals.
type PreferencesRepository struct {
	db *sql.DB
}

// NewPreferencesRepository creates a new PreferencesRepository.
func NewPreferencesRepository(db *sql.DB) *PreferencesRepository {
	return &PreferencesRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// User Preferences
// ─────────────────────────────────────────────────────────────────────────────

// GetPreferences retrieves a user's preferences.
func (r *PreferencesRepository) GetPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error) {
	p := &UserPreferences{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, desired_job_titles, desired_industries, desired_company_sizes,
		       desired_location_types, desired_locations, is_willing_to_relocate,
		       relocation_locations, salary_currency, salary_min, salary_max,
		       include_equity, job_search_urgency, available_from, career_stage,
		       email_job_alerts, email_weekly_digest, email_training_recs,
		       created_at, updated_at
		FROM user_preferences WHERE user_id = $1`, userID,
	).Scan(
		&p.ID, &p.UserID,
		&p.DesiredJobTitles, &p.DesiredIndustries, &p.DesiredCompanySizes,
		&p.DesiredLocationTypes, &p.DesiredLocations, &p.IsWillingToRelocate,
		&p.RelocationLocations, &p.SalaryCurrency, &p.SalaryMin, &p.SalaryMax,
		&p.IncludeEquity, &p.JobSearchUrgency, &p.AvailableFrom, &p.CareerStage,
		&p.EmailJobAlerts, &p.EmailWeeklyDigest, &p.EmailTrainingRecs,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	return p, nil
}

// UpdatePreferences updates a user's job search preferences.
func (r *PreferencesRepository) UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs UserPreferences) (*UserPreferences, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE user_preferences SET
			desired_job_titles      = $2,
			desired_industries      = $3,
			desired_company_sizes   = $4,
			desired_location_types  = $5,
			desired_locations       = $6,
			is_willing_to_relocate  = $7,
			relocation_locations    = $8,
			salary_currency         = $9,
			salary_min              = $10,
			salary_max              = $11,
			include_equity          = $12,
			job_search_urgency      = $13,
			available_from          = $14,
			career_stage            = $15,
			email_job_alerts        = $16,
			email_weekly_digest     = $17,
			email_training_recs     = $18
		WHERE user_id = $1
		RETURNING id, user_id, desired_job_titles, desired_industries, desired_company_sizes,
		          desired_location_types, desired_locations, is_willing_to_relocate,
		          relocation_locations, salary_currency, salary_min, salary_max,
		          include_equity, job_search_urgency, available_from, career_stage,
		          email_job_alerts, email_weekly_digest, email_training_recs,
		          created_at, updated_at`,
		userID,
		pq.Array(prefs.DesiredJobTitles), pq.Array(prefs.DesiredIndustries),
		pq.Array(prefs.DesiredCompanySizes), pq.Array(prefs.DesiredLocationTypes),
		pq.Array(prefs.DesiredLocations), prefs.IsWillingToRelocate,
		pq.Array(prefs.RelocationLocations), prefs.SalaryCurrency,
		prefs.SalaryMin, prefs.SalaryMax, prefs.IncludeEquity,
		prefs.JobSearchUrgency, prefs.AvailableFrom, prefs.CareerStage,
		prefs.EmailJobAlerts, prefs.EmailWeeklyDigest, prefs.EmailTrainingRecs,
	).Scan(
		&prefs.ID, &prefs.UserID,
		&prefs.DesiredJobTitles, &prefs.DesiredIndustries, &prefs.DesiredCompanySizes,
		&prefs.DesiredLocationTypes, &prefs.DesiredLocations, &prefs.IsWillingToRelocate,
		&prefs.RelocationLocations, &prefs.SalaryCurrency, &prefs.SalaryMin, &prefs.SalaryMax,
		&prefs.IncludeEquity, &prefs.JobSearchUrgency, &prefs.AvailableFrom, &prefs.CareerStage,
		&prefs.EmailJobAlerts, &prefs.EmailWeeklyDigest, &prefs.EmailTrainingRecs,
		&prefs.CreatedAt, &prefs.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update preferences: %w", err)
	}

	// Record history
	r.db.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type)
		VALUES ($1, $2, 'preferences')`,
		userID, EventPreferenceUpdated,
	)

	return &prefs, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Career Goals
// ─────────────────────────────────────────────────────────────────────────────

// CreateCareerGoal creates a new career goal.
func (r *PreferencesRepository) CreateCareerGoal(ctx context.Context, userID uuid.UUID, goal CareerGoal) (*CareerGoal, error) {
	goal.UserID = userID
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO career_goals (
			user_id, title, description, target_role, target_industry,
			target_date, status, priority, progress_percentage, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, user_id, title, description, target_role, target_industry,
		          target_date, status, priority, progress_percentage, notes,
		          achieved_at, created_at, updated_at`,
		userID, goal.Title, goal.Description, goal.TargetRole,
		goal.TargetIndustry, goal.TargetDate, goal.Status,
		goal.Priority, goal.ProgressPercentage, goal.Notes,
	).Scan(
		&goal.ID, &goal.UserID, &goal.Title, &goal.Description,
		&goal.TargetRole, &goal.TargetIndustry, &goal.TargetDate,
		&goal.Status, &goal.Priority, &goal.ProgressPercentage,
		&goal.Notes, &goal.AchievedAt, &goal.CreatedAt, &goal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create career goal: %w", err)
	}
	return &goal, nil
}

// GetCareerGoalsByUserID retrieves all career goals for a user.
func (r *PreferencesRepository) GetCareerGoalsByUserID(ctx context.Context, userID uuid.UUID) ([]CareerGoal, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, title, description, target_role, target_industry,
		       target_date, status, priority, progress_percentage, notes,
		       achieved_at, created_at, updated_at
		FROM career_goals
		WHERE user_id = $1
		ORDER BY
		    CASE status WHEN 'active' THEN 1 WHEN 'paused' THEN 2 ELSE 3 END,
		    priority ASC, created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get career goals: %w", err)
	}
	defer rows.Close()

	var goals []CareerGoal
	for rows.Next() {
		var g CareerGoal
		if err := rows.Scan(
			&g.ID, &g.UserID, &g.Title, &g.Description,
			&g.TargetRole, &g.TargetIndustry, &g.TargetDate,
			&g.Status, &g.Priority, &g.ProgressPercentage,
			&g.Notes, &g.AchievedAt, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan career goal: %w", err)
		}
		goals = append(goals, g)
	}
	return goals, rows.Err()
}

// UpdateGoalProgress updates the progress percentage of a career goal.
func (r *PreferencesRepository) UpdateGoalProgress(ctx context.Context, userID, goalID uuid.UUID, progress int16) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE career_goals
		SET progress_percentage = $3,
		    status = CASE WHEN $3 = 100 THEN 'achieved'::goal_status ELSE status END,
		    achieved_at = CASE WHEN $3 = 100 THEN NOW() ELSE achieved_at END
		WHERE id = $1 AND user_id = $2`,
		goalID, userID, progress,
	)
	if err != nil {
		return fmt.Errorf("update goal progress: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// GetActiveSkillGaps retrieves unaddressed skill gaps for a user.
func (r *PreferencesRepository) GetActiveSkillGaps(ctx context.Context, userID uuid.UUID) ([]SkillGap, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT sg.id, sg.user_id, sg.career_goal_id, sg.skill_name, sg.normalized_name,
		       sg.skill_taxonomy_id, sg.gap_type, sg.required_proficiency,
		       sg.current_proficiency, sg.importance, sg.is_addressed,
		       sg.identified_at, sg.addressed_at, sg.created_at, sg.updated_at
		FROM skill_gaps sg
		WHERE sg.user_id = $1 AND sg.is_addressed = FALSE
		ORDER BY
		    CASE sg.importance WHEN 'critical' THEN 1 WHEN 'important' THEN 2 ELSE 3 END,
		    sg.identified_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get skill gaps: %w", err)
	}
	defer rows.Close()

	var gaps []SkillGap
	for rows.Next() {
		var g SkillGap
		if err := rows.Scan(
			&g.ID, &g.UserID, &g.CareerGoalID, &g.SkillName, &g.NormalizedName,
			&g.SkillTaxonomyID, &g.GapType, &g.RequiredProficiency,
			&g.CurrentProficiency, &g.Importance, &g.IsAddressed,
			&g.IdentifiedAt, &g.AddressedAt, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan skill gap: %w", err)
		}
		gaps = append(gaps, g)
	}
	return gaps, rows.Err()
}

// SkillGap represents an identified skill gap.
type SkillGap struct {
	ID                  uuid.UUID        `db:"id" json:"id"`
	UserID              uuid.UUID        `db:"user_id" json:"user_id"`
	CareerGoalID        *uuid.UUID       `db:"career_goal_id" json:"career_goal_id,omitempty"`
	SkillName           string           `db:"skill_name" json:"skill_name"`
	NormalizedName      string           `db:"normalized_name" json:"normalized_name"`
	SkillTaxonomyID     *uuid.UUID       `db:"skill_taxonomy_id" json:"skill_taxonomy_id,omitempty"`
	GapType             string           `db:"gap_type" json:"gap_type"`
	RequiredProficiency SkillProficiency `db:"required_proficiency" json:"required_proficiency"`
	CurrentProficiency  *SkillProficiency `db:"current_proficiency" json:"current_proficiency,omitempty"`
	Importance          string           `db:"importance" json:"importance"`
	IsAddressed         bool             `db:"is_addressed" json:"is_addressed"`
	IdentifiedAt        interface{}      `db:"identified_at" json:"identified_at"`
	AddressedAt         sql.NullTime     `db:"addressed_at" json:"addressed_at,omitempty"`
	CreatedAt           interface{}      `db:"created_at" json:"created_at"`
	UpdatedAt           interface{}      `db:"updated_at" json:"updated_at"`
}
