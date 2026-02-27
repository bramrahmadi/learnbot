package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ExperienceRepository provides CRUD operations for work experience, education,
// certifications, and projects.
type ExperienceRepository struct {
	db *sql.DB
}

// NewExperienceRepository creates a new ExperienceRepository.
func NewExperienceRepository(db *sql.DB) *ExperienceRepository {
	return &ExperienceRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// Work Experience
// ─────────────────────────────────────────────────────────────────────────────

// CreateWorkExperience adds a new work experience entry.
func (r *ExperienceRepository) CreateWorkExperience(ctx context.Context, userID uuid.UUID, input CreateWorkExperienceInput) (*WorkExperience, error) {
	exp := &WorkExperience{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO work_experience (
			user_id, resume_upload_id, company_name, job_title, employment_type,
			location_type, location, start_date, end_date, is_current,
			description, responsibilities, technologies_used, confidence
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, user_id, resume_upload_id, company_name, job_title,
		          employment_type, location_type, location, start_date, end_date,
		          is_current, description, responsibilities, technologies_used,
		          duration_months, confidence, display_order, created_at, updated_at`,
		userID, input.ResumeUploadID, input.CompanyName, input.JobTitle,
		input.EmploymentType, input.LocationType, input.Location,
		input.StartDate, input.EndDate, input.IsCurrent,
		input.Description, pq.Array(input.Responsibilities),
		pq.Array(input.TechnologiesUsed), input.Confidence,
	).Scan(
		&exp.ID, &exp.UserID, &exp.ResumeUploadID,
		&exp.CompanyName, &exp.JobTitle, &exp.EmploymentType,
		&exp.LocationType, &exp.Location, &exp.StartDate, &exp.EndDate,
		&exp.IsCurrent, &exp.Description, &exp.Responsibilities,
		&exp.TechnologiesUsed, &exp.DurationMonths, &exp.Confidence,
		&exp.DisplayOrder, &exp.CreatedAt, &exp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create work experience: %w", err)
	}

	// Recalculate years of experience on profile
	r.db.ExecContext(ctx, `
		UPDATE user_profiles
		SET years_of_experience = calculate_years_of_experience($1),
		    profile_completeness = calculate_profile_completeness($1)
		WHERE user_id = $1`, userID)

	// Record history
	r.db.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type, entity_id)
		VALUES ($1, $2, 'work_experience', $3)`,
		userID, EventExperienceAdded, exp.ID,
	)

	return exp, nil
}

// GetWorkExperienceByUserID retrieves all work experience entries for a user.
func (r *ExperienceRepository) GetWorkExperienceByUserID(ctx context.Context, userID uuid.UUID) ([]WorkExperience, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, resume_upload_id, company_name, job_title,
		       employment_type, location_type, location, start_date, end_date,
		       is_current, description, responsibilities, technologies_used,
		       duration_months, confidence, display_order, created_at, updated_at
		FROM work_experience
		WHERE user_id = $1
		ORDER BY is_current DESC, start_date DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get work experience: %w", err)
	}
	defer rows.Close()

	var experiences []WorkExperience
	for rows.Next() {
		var e WorkExperience
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.ResumeUploadID,
			&e.CompanyName, &e.JobTitle, &e.EmploymentType,
			&e.LocationType, &e.Location, &e.StartDate, &e.EndDate,
			&e.IsCurrent, &e.Description, &e.Responsibilities,
			&e.TechnologiesUsed, &e.DurationMonths, &e.Confidence,
			&e.DisplayOrder, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan work experience: %w", err)
		}
		experiences = append(experiences, e)
	}
	return experiences, rows.Err()
}

// UpdateWorkExperience updates a work experience entry.
func (r *ExperienceRepository) UpdateWorkExperience(ctx context.Context, userID, expID uuid.UUID, input CreateWorkExperienceInput) (*WorkExperience, error) {
	exp := &WorkExperience{}
	err := r.db.QueryRowContext(ctx, `
		UPDATE work_experience SET
			company_name      = $3,
			job_title         = $4,
			employment_type   = $5,
			location_type     = $6,
			location          = $7,
			start_date        = $8,
			end_date          = $9,
			is_current        = $10,
			description       = $11,
			responsibilities  = $12,
			technologies_used = $13,
			confidence        = $14
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, resume_upload_id, company_name, job_title,
		          employment_type, location_type, location, start_date, end_date,
		          is_current, description, responsibilities, technologies_used,
		          duration_months, confidence, display_order, created_at, updated_at`,
		expID, userID,
		input.CompanyName, input.JobTitle, input.EmploymentType,
		input.LocationType, input.Location, input.StartDate, input.EndDate,
		input.IsCurrent, input.Description,
		pq.Array(input.Responsibilities), pq.Array(input.TechnologiesUsed),
		input.Confidence,
	).Scan(
		&exp.ID, &exp.UserID, &exp.ResumeUploadID,
		&exp.CompanyName, &exp.JobTitle, &exp.EmploymentType,
		&exp.LocationType, &exp.Location, &exp.StartDate, &exp.EndDate,
		&exp.IsCurrent, &exp.Description, &exp.Responsibilities,
		&exp.TechnologiesUsed, &exp.DurationMonths, &exp.Confidence,
		&exp.DisplayOrder, &exp.CreatedAt, &exp.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update work experience: %w", err)
	}

	// Recalculate years of experience
	r.db.ExecContext(ctx, `
		UPDATE user_profiles
		SET years_of_experience = calculate_years_of_experience($1)
		WHERE user_id = $1`, userID)

	r.db.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type, entity_id)
		VALUES ($1, $2, 'work_experience', $3)`,
		userID, EventExperienceUpdated, expID,
	)

	return exp, nil
}

// DeleteWorkExperience removes a work experience entry.
func (r *ExperienceRepository) DeleteWorkExperience(ctx context.Context, userID, expID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM work_experience WHERE id = $1 AND user_id = $2`, expID, userID)
	if err != nil {
		return fmt.Errorf("delete work experience: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	// Recalculate years of experience
	r.db.ExecContext(ctx, `
		UPDATE user_profiles
		SET years_of_experience = calculate_years_of_experience($1),
		    profile_completeness = calculate_profile_completeness($1)
		WHERE user_id = $1`, userID)

	return nil
}

// GetTotalExperienceMonths returns the total work experience in months for a user.
func (r *ExperienceRepository) GetTotalExperienceMonths(ctx context.Context, userID uuid.UUID) (int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(duration_months), 0)
		FROM work_experience WHERE user_id = $1`, userID,
	).Scan(&total)
	return total, err
}

// ─────────────────────────────────────────────────────────────────────────────
// Education
// ─────────────────────────────────────────────────────────────────────────────

// CreateEducation adds a new education entry.
func (r *ExperienceRepository) CreateEducation(ctx context.Context, userID uuid.UUID, input CreateEducationInput) (*Education, error) {
	gpaScale := input.GPAScale
	if gpaScale == 0 {
		gpaScale = 4.0
	}

	edu := &Education{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO education (
			user_id, resume_upload_id, institution_name, degree_level, degree_name,
			field_of_study, start_date, end_date, is_current, gpa, gpa_scale,
			honors, activities, confidence
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, user_id, resume_upload_id, institution_name, degree_level,
		          degree_name, field_of_study, start_date, end_date, is_current,
		          gpa, gpa_scale, honors, activities, confidence,
		          display_order, created_at, updated_at`,
		userID, input.ResumeUploadID, input.InstitutionName, input.DegreeLevel,
		input.DegreeName, input.FieldOfStudy, input.StartDate, input.EndDate,
		input.IsCurrent, input.GPA, gpaScale, input.Honors,
		pq.Array(input.Activities), input.Confidence,
	).Scan(
		&edu.ID, &edu.UserID, &edu.ResumeUploadID,
		&edu.InstitutionName, &edu.DegreeLevel, &edu.DegreeName,
		&edu.FieldOfStudy, &edu.StartDate, &edu.EndDate, &edu.IsCurrent,
		&edu.GPA, &edu.GPAScale, &edu.Honors, &edu.Activities,
		&edu.Confidence, &edu.DisplayOrder, &edu.CreatedAt, &edu.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create education: %w", err)
	}

	// Recalculate completeness
	r.db.ExecContext(ctx, `
		UPDATE user_profiles
		SET profile_completeness = calculate_profile_completeness($1)
		WHERE user_id = $1`, userID)

	r.db.ExecContext(ctx, `
		INSERT INTO profile_history (user_id, event_type, entity_type, entity_id)
		VALUES ($1, $2, 'education', $3)`,
		userID, EventEducationAdded, edu.ID,
	)

	return edu, nil
}

// GetEducationByUserID retrieves all education entries for a user.
func (r *ExperienceRepository) GetEducationByUserID(ctx context.Context, userID uuid.UUID) ([]Education, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, resume_upload_id, institution_name, degree_level,
		       degree_name, field_of_study, start_date, end_date, is_current,
		       gpa, gpa_scale, honors, activities, confidence,
		       display_order, created_at, updated_at
		FROM education
		WHERE user_id = $1
		ORDER BY end_date DESC NULLS FIRST, start_date DESC NULLS FIRST`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get education: %w", err)
	}
	defer rows.Close()

	var educations []Education
	for rows.Next() {
		var e Education
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.ResumeUploadID,
			&e.InstitutionName, &e.DegreeLevel, &e.DegreeName,
			&e.FieldOfStudy, &e.StartDate, &e.EndDate, &e.IsCurrent,
			&e.GPA, &e.GPAScale, &e.Honors, &e.Activities,
			&e.Confidence, &e.DisplayOrder, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan education: %w", err)
		}
		educations = append(educations, e)
	}
	return educations, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Certifications
// ─────────────────────────────────────────────────────────────────────────────

// CreateCertification adds a new certification entry.
func (r *ExperienceRepository) CreateCertification(ctx context.Context, userID uuid.UUID, cert Certification) (*Certification, error) {
	cert.UserID = userID
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO certifications (
			user_id, resume_upload_id, name, issuing_organization,
			issue_date, expiry_date, credential_id, credential_url, confidence
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, resume_upload_id, name, issuing_organization,
		          issue_date, expiry_date, credential_id, credential_url,
		          is_expired, confidence, display_order, created_at, updated_at`,
		userID, cert.ResumeUploadID, cert.Name, cert.IssuingOrganization,
		cert.IssueDate, cert.ExpiryDate, cert.CredentialID,
		cert.CredentialURL, cert.Confidence,
	).Scan(
		&cert.ID, &cert.UserID, &cert.ResumeUploadID,
		&cert.Name, &cert.IssuingOrganization,
		&cert.IssueDate, &cert.ExpiryDate, &cert.CredentialID,
		&cert.CredentialURL, &cert.IsExpired, &cert.Confidence,
		&cert.DisplayOrder, &cert.CreatedAt, &cert.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create certification: %w", err)
	}
	return &cert, nil
}

// GetCertificationsByUserID retrieves all certifications for a user.
func (r *ExperienceRepository) GetCertificationsByUserID(ctx context.Context, userID uuid.UUID) ([]Certification, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, resume_upload_id, name, issuing_organization,
		       issue_date, expiry_date, credential_id, credential_url,
		       is_expired, confidence, display_order, created_at, updated_at
		FROM certifications
		WHERE user_id = $1
		ORDER BY is_expired ASC, issue_date DESC NULLS LAST`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get certifications: %w", err)
	}
	defer rows.Close()

	var certs []Certification
	for rows.Next() {
		var c Certification
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.ResumeUploadID,
			&c.Name, &c.IssuingOrganization,
			&c.IssueDate, &c.ExpiryDate, &c.CredentialID,
			&c.CredentialURL, &c.IsExpired, &c.Confidence,
			&c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan certification: %w", err)
		}
		certs = append(certs, c)
	}
	return certs, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Projects
// ─────────────────────────────────────────────────────────────────────────────

// CreateProject adds a new project entry.
func (r *ExperienceRepository) CreateProject(ctx context.Context, userID uuid.UUID, proj Project) (*Project, error) {
	proj.UserID = userID
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO projects (
			user_id, resume_upload_id, name, description, project_url,
			repository_url, technologies, start_date, end_date, is_ongoing, confidence
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, resume_upload_id, name, description, project_url,
		          repository_url, technologies, start_date, end_date, is_ongoing,
		          confidence, display_order, created_at, updated_at`,
		userID, proj.ResumeUploadID, proj.Name, proj.Description,
		proj.ProjectURL, proj.RepositoryURL, pq.Array(proj.Technologies),
		proj.StartDate, proj.EndDate, proj.IsOngoing, proj.Confidence,
	).Scan(
		&proj.ID, &proj.UserID, &proj.ResumeUploadID,
		&proj.Name, &proj.Description, &proj.ProjectURL,
		&proj.RepositoryURL, &proj.Technologies, &proj.StartDate,
		&proj.EndDate, &proj.IsOngoing, &proj.Confidence,
		&proj.DisplayOrder, &proj.CreatedAt, &proj.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return &proj, nil
}

// GetProjectsByUserID retrieves all projects for a user.
func (r *ExperienceRepository) GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, resume_upload_id, name, description, project_url,
		       repository_url, technologies, start_date, end_date, is_ongoing,
		       confidence, display_order, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY is_ongoing DESC, end_date DESC NULLS FIRST`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.ResumeUploadID,
			&p.Name, &p.Description, &p.ProjectURL,
			&p.RepositoryURL, &p.Technologies, &p.StartDate,
			&p.EndDate, &p.IsOngoing, &p.Confidence,
			&p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// ensure time is used
var _ = time.Now
