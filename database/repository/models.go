// Package repository provides data access objects and CRUD operations
// for the LearnBot user profile data model.
package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ─────────────────────────────────────────────────────────────────────────────
// Enum types (mirror PostgreSQL enums)
// ─────────────────────────────────────────────────────────────────────────────

// SkillProficiency represents the proficiency level of a skill.
type SkillProficiency string

const (
	ProficiencyBeginner     SkillProficiency = "beginner"
	ProficiencyIntermediate SkillProficiency = "intermediate"
	ProficiencyAdvanced     SkillProficiency = "advanced"
	ProficiencyExpert       SkillProficiency = "expert"
)

// DegreeLevel represents the level of an educational degree.
type DegreeLevel string

const (
	DegreeHighSchool  DegreeLevel = "high_school"
	DegreeAssociate   DegreeLevel = "associate"
	DegreeBachelor    DegreeLevel = "bachelor"
	DegreeMaster      DegreeLevel = "master"
	DegreeDoctorate   DegreeLevel = "doctorate"
	DegreeProfessional DegreeLevel = "professional"
	DegreeCertificate DegreeLevel = "certificate"
	DegreeDiploma     DegreeLevel = "diploma"
	DegreeOther       DegreeLevel = "other"
)

// SkillCategory represents the category of a skill.
type SkillCategory string

const (
	CategoryTechnical  SkillCategory = "technical"
	CategorySoft       SkillCategory = "soft"
	CategoryLanguage   SkillCategory = "language"
	CategoryTool       SkillCategory = "tool"
	CategoryFramework  SkillCategory = "framework"
	CategoryDatabase   SkillCategory = "database"
	CategoryCloud      SkillCategory = "cloud"
	CategoryOther      SkillCategory = "other"
)

// EmploymentType represents the type of employment.
type EmploymentType string

const (
	EmploymentFullTime  EmploymentType = "full_time"
	EmploymentPartTime  EmploymentType = "part_time"
	EmploymentContract  EmploymentType = "contract"
	EmploymentFreelance EmploymentType = "freelance"
	EmploymentInternship EmploymentType = "internship"
	EmploymentVolunteer EmploymentType = "volunteer"
	EmploymentOther     EmploymentType = "other"
)

// WorkLocationType represents the work location type.
type WorkLocationType string

const (
	LocationOnSite WorkLocationType = "on_site"
	LocationRemote WorkLocationType = "remote"
	LocationHybrid WorkLocationType = "hybrid"
)

// GoalStatus represents the status of a career goal.
type GoalStatus string

const (
	GoalActive    GoalStatus = "active"
	GoalAchieved  GoalStatus = "achieved"
	GoalPaused    GoalStatus = "paused"
	GoalAbandoned GoalStatus = "abandoned"
)

// ProfileEventType represents the type of profile change event.
type ProfileEventType string

const (
	EventCreated           ProfileEventType = "created"
	EventResumeUploaded    ProfileEventType = "resume_uploaded"
	EventManuallyUpdated   ProfileEventType = "manually_updated"
	EventSkillAdded        ProfileEventType = "skill_added"
	EventSkillRemoved      ProfileEventType = "skill_removed"
	EventExperienceAdded   ProfileEventType = "experience_added"
	EventExperienceUpdated ProfileEventType = "experience_updated"
	EventEducationAdded    ProfileEventType = "education_added"
	EventPreferenceUpdated ProfileEventType = "preference_updated"
)

// ─────────────────────────────────────────────────────────────────────────────
// Model structs
// ─────────────────────────────────────────────────────────────────────────────

// User represents a user account.
type User struct {
	ID            uuid.UUID      `db:"id" json:"id"`
	Email         string         `db:"email" json:"email"`
	EmailVerified bool           `db:"email_verified" json:"email_verified"`
	PasswordHash  sql.NullString `db:"password_hash" json:"-"`
	FullName      string         `db:"full_name" json:"full_name"`
	AvatarURL     sql.NullString `db:"avatar_url" json:"avatar_url,omitempty"`
	Timezone      string         `db:"timezone" json:"timezone"`
	Locale        string         `db:"locale" json:"locale"`
	IsActive      bool           `db:"is_active" json:"is_active"`
	IsAdmin       bool           `db:"is_admin" json:"is_admin"`
	LastLoginAt   sql.NullTime   `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}

// UserProfile represents a user's professional profile.
type UserProfile struct {
	ID                  uuid.UUID      `db:"id" json:"id"`
	UserID              uuid.UUID      `db:"user_id" json:"user_id"`
	Headline            sql.NullString `db:"headline" json:"headline,omitempty"`
	Summary             sql.NullString `db:"summary" json:"summary,omitempty"`
	LocationCity        sql.NullString `db:"location_city" json:"location_city,omitempty"`
	LocationState       sql.NullString `db:"location_state" json:"location_state,omitempty"`
	LocationCountry     string         `db:"location_country" json:"location_country"`
	Phone               sql.NullString `db:"phone" json:"phone,omitempty"`
	LinkedInURL         sql.NullString `db:"linkedin_url" json:"linkedin_url,omitempty"`
	GitHubURL           sql.NullString `db:"github_url" json:"github_url,omitempty"`
	WebsiteURL          sql.NullString `db:"website_url" json:"website_url,omitempty"`
	YearsOfExperience   sql.NullFloat64 `db:"years_of_experience" json:"years_of_experience,omitempty"`
	IsOpenToWork        bool           `db:"is_open_to_work" json:"is_open_to_work"`
	ProfileCompleteness int16          `db:"profile_completeness" json:"profile_completeness"`
	CreatedAt           time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at" json:"updated_at"`
}

// ResumeUpload represents a versioned resume file upload.
type ResumeUpload struct {
	ID                 uuid.UUID      `db:"id" json:"id"`
	UserID             uuid.UUID      `db:"user_id" json:"user_id"`
	FileName           string         `db:"file_name" json:"file_name"`
	FileType           string         `db:"file_type" json:"file_type"`
	FileSizeBytes      int            `db:"file_size_bytes" json:"file_size_bytes"`
	StorageKey         string         `db:"storage_key" json:"storage_key"`
	Version            int            `db:"version" json:"version"`
	IsCurrent          bool           `db:"is_current" json:"is_current"`
	ParsedAt           sql.NullTime   `db:"parsed_at" json:"parsed_at,omitempty"`
	ParseStatus        string         `db:"parse_status" json:"parse_status"`
	ParseError         sql.NullString `db:"parse_error" json:"parse_error,omitempty"`
	RawText            sql.NullString `db:"raw_text" json:"-"`
	ParserVersion      sql.NullString `db:"parser_version" json:"parser_version,omitempty"`
	OverallConfidence  sql.NullFloat64 `db:"overall_confidence" json:"overall_confidence,omitempty"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
}

// UserSkill represents a skill associated with a user.
type UserSkill struct {
	ID                uuid.UUID       `db:"id" json:"id"`
	UserID            uuid.UUID       `db:"user_id" json:"user_id"`
	SkillTaxonomyID   *uuid.UUID      `db:"skill_taxonomy_id" json:"skill_taxonomy_id,omitempty"`
	SkillName         string          `db:"skill_name" json:"skill_name"`
	NormalizedName    string          `db:"normalized_name" json:"normalized_name"`
	Category          SkillCategory   `db:"category" json:"category"`
	Proficiency       SkillProficiency `db:"proficiency" json:"proficiency"`
	YearsOfExperience sql.NullFloat64 `db:"years_of_experience" json:"years_of_experience,omitempty"`
	IsPrimary         bool            `db:"is_primary" json:"is_primary"`
	Source            string          `db:"source" json:"source"`
	Confidence        sql.NullFloat64 `db:"confidence" json:"confidence,omitempty"`
	LastUsedYear      sql.NullInt16   `db:"last_used_year" json:"last_used_year,omitempty"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
}

// WorkExperience represents a job entry in a user's work history.
type WorkExperience struct {
	ID               uuid.UUID       `db:"id" json:"id"`
	UserID           uuid.UUID       `db:"user_id" json:"user_id"`
	ResumeUploadID   *uuid.UUID      `db:"resume_upload_id" json:"resume_upload_id,omitempty"`
	CompanyName      string          `db:"company_name" json:"company_name"`
	JobTitle         string          `db:"job_title" json:"job_title"`
	EmploymentType   EmploymentType  `db:"employment_type" json:"employment_type"`
	LocationType     *WorkLocationType `db:"location_type" json:"location_type,omitempty"`
	Location         sql.NullString  `db:"location" json:"location,omitempty"`
	StartDate        time.Time       `db:"start_date" json:"start_date"`
	EndDate          sql.NullTime    `db:"end_date" json:"end_date,omitempty"`
	IsCurrent        bool            `db:"is_current" json:"is_current"`
	Description      sql.NullString  `db:"description" json:"description,omitempty"`
	Responsibilities pq.StringArray  `db:"responsibilities" json:"responsibilities,omitempty"`
	TechnologiesUsed pq.StringArray  `db:"technologies_used" json:"technologies_used,omitempty"`
	DurationMonths   int             `db:"duration_months" json:"duration_months"`
	Confidence       sql.NullFloat64 `db:"confidence" json:"confidence,omitempty"`
	DisplayOrder     int16           `db:"display_order" json:"display_order"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
}

// Education represents an educational entry.
type Education struct {
	ID              uuid.UUID       `db:"id" json:"id"`
	UserID          uuid.UUID       `db:"user_id" json:"user_id"`
	ResumeUploadID  *uuid.UUID      `db:"resume_upload_id" json:"resume_upload_id,omitempty"`
	InstitutionName string          `db:"institution_name" json:"institution_name"`
	DegreeLevel     DegreeLevel     `db:"degree_level" json:"degree_level"`
	DegreeName      sql.NullString  `db:"degree_name" json:"degree_name,omitempty"`
	FieldOfStudy    sql.NullString  `db:"field_of_study" json:"field_of_study,omitempty"`
	StartDate       sql.NullTime    `db:"start_date" json:"start_date,omitempty"`
	EndDate         sql.NullTime    `db:"end_date" json:"end_date,omitempty"`
	IsCurrent       bool            `db:"is_current" json:"is_current"`
	GPA             sql.NullFloat64 `db:"gpa" json:"gpa,omitempty"`
	GPAScale        float64         `db:"gpa_scale" json:"gpa_scale"`
	Honors          sql.NullString  `db:"honors" json:"honors,omitempty"`
	Activities      pq.StringArray  `db:"activities" json:"activities,omitempty"`
	Confidence      sql.NullFloat64 `db:"confidence" json:"confidence,omitempty"`
	DisplayOrder    int16           `db:"display_order" json:"display_order"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

// Certification represents a professional certification.
type Certification struct {
	ID                  uuid.UUID      `db:"id" json:"id"`
	UserID              uuid.UUID      `db:"user_id" json:"user_id"`
	ResumeUploadID      *uuid.UUID     `db:"resume_upload_id" json:"resume_upload_id,omitempty"`
	Name                string         `db:"name" json:"name"`
	IssuingOrganization sql.NullString `db:"issuing_organization" json:"issuing_organization,omitempty"`
	IssueDate           sql.NullTime   `db:"issue_date" json:"issue_date,omitempty"`
	ExpiryDate          sql.NullTime   `db:"expiry_date" json:"expiry_date,omitempty"`
	CredentialID        sql.NullString `db:"credential_id" json:"credential_id,omitempty"`
	CredentialURL       sql.NullString `db:"credential_url" json:"credential_url,omitempty"`
	IsExpired           bool           `db:"is_expired" json:"is_expired"`
	Confidence          sql.NullFloat64 `db:"confidence" json:"confidence,omitempty"`
	DisplayOrder        int16          `db:"display_order" json:"display_order"`
	CreatedAt           time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at" json:"updated_at"`
}

// Project represents a portfolio project.
type Project struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	UserID         uuid.UUID      `db:"user_id" json:"user_id"`
	ResumeUploadID *uuid.UUID     `db:"resume_upload_id" json:"resume_upload_id,omitempty"`
	Name           string         `db:"name" json:"name"`
	Description    sql.NullString `db:"description" json:"description,omitempty"`
	ProjectURL     sql.NullString `db:"project_url" json:"project_url,omitempty"`
	RepositoryURL  sql.NullString `db:"repository_url" json:"repository_url,omitempty"`
	Technologies   pq.StringArray `db:"technologies" json:"technologies,omitempty"`
	StartDate      sql.NullTime   `db:"start_date" json:"start_date,omitempty"`
	EndDate        sql.NullTime   `db:"end_date" json:"end_date,omitempty"`
	IsOngoing      bool           `db:"is_ongoing" json:"is_ongoing"`
	Confidence     sql.NullFloat64 `db:"confidence" json:"confidence,omitempty"`
	DisplayOrder   int16          `db:"display_order" json:"display_order"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// UserPreferences represents a user's job search and career preferences.
type UserPreferences struct {
	ID                    uuid.UUID        `db:"id" json:"id"`
	UserID                uuid.UUID        `db:"user_id" json:"user_id"`
	DesiredJobTitles      pq.StringArray   `db:"desired_job_titles" json:"desired_job_titles,omitempty"`
	DesiredIndustries     pq.StringArray   `db:"desired_industries" json:"desired_industries,omitempty"`
	DesiredCompanySizes   pq.StringArray   `db:"desired_company_sizes" json:"desired_company_sizes,omitempty"`
	DesiredLocationTypes  pq.StringArray   `db:"desired_location_types" json:"desired_location_types,omitempty"`
	DesiredLocations      pq.StringArray   `db:"desired_locations" json:"desired_locations,omitempty"`
	IsWillingToRelocate   bool             `db:"is_willing_to_relocate" json:"is_willing_to_relocate"`
	RelocationLocations   pq.StringArray   `db:"relocation_locations" json:"relocation_locations,omitempty"`
	SalaryCurrency        string           `db:"salary_currency" json:"salary_currency"`
	SalaryMin             sql.NullInt32    `db:"salary_min" json:"salary_min,omitempty"`
	SalaryMax             sql.NullInt32    `db:"salary_max" json:"salary_max,omitempty"`
	IncludeEquity         bool             `db:"include_equity" json:"include_equity"`
	JobSearchUrgency      string           `db:"job_search_urgency" json:"job_search_urgency"`
	AvailableFrom         sql.NullTime     `db:"available_from" json:"available_from,omitempty"`
	CareerStage           sql.NullString   `db:"career_stage" json:"career_stage,omitempty"`
	EmailJobAlerts        bool             `db:"email_job_alerts" json:"email_job_alerts"`
	EmailWeeklyDigest     bool             `db:"email_weekly_digest" json:"email_weekly_digest"`
	EmailTrainingRecs     bool             `db:"email_training_recs" json:"email_training_recs"`
	CreatedAt             time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time        `db:"updated_at" json:"updated_at"`
}

// CareerGoal represents a specific career objective.
type CareerGoal struct {
	ID                 uuid.UUID      `db:"id" json:"id"`
	UserID             uuid.UUID      `db:"user_id" json:"user_id"`
	Title              string         `db:"title" json:"title"`
	Description        sql.NullString `db:"description" json:"description,omitempty"`
	TargetRole         sql.NullString `db:"target_role" json:"target_role,omitempty"`
	TargetIndustry     sql.NullString `db:"target_industry" json:"target_industry,omitempty"`
	TargetDate         sql.NullTime   `db:"target_date" json:"target_date,omitempty"`
	Status             GoalStatus     `db:"status" json:"status"`
	Priority           int16          `db:"priority" json:"priority"`
	ProgressPercentage int16          `db:"progress_percentage" json:"progress_percentage"`
	Notes              sql.NullString `db:"notes" json:"notes,omitempty"`
	AchievedAt         sql.NullTime   `db:"achieved_at" json:"achieved_at,omitempty"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at" json:"updated_at"`
}

// ProfileHistory represents an audit log entry for profile changes.
type ProfileHistory struct {
	ID          int64            `db:"id" json:"id"`
	UserID      uuid.UUID        `db:"user_id" json:"user_id"`
	EventType   ProfileEventType `db:"event_type" json:"event_type"`
	EntityType  sql.NullString   `db:"entity_type" json:"entity_type,omitempty"`
	EntityID    *uuid.UUID       `db:"entity_id" json:"entity_id,omitempty"`
	OldData     []byte           `db:"old_data" json:"old_data,omitempty"`
	NewData     []byte           `db:"new_data" json:"new_data,omitempty"`
	ChangedBy   *uuid.UUID       `db:"changed_by" json:"changed_by,omitempty"`
	IPAddress   sql.NullString   `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent   sql.NullString   `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Input/request types for CRUD operations
// ─────────────────────────────────────────────────────────────────────────────

// CreateUserInput holds the data needed to create a new user.
type CreateUserInput struct {
	Email        string
	PasswordHash *string
	FullName     string
	AvatarURL    *string
	Timezone     string
	Locale       string
}

// UpdateProfileInput holds the data for updating a user profile.
type UpdateProfileInput struct {
	Headline          *string
	Summary           *string
	LocationCity      *string
	LocationState     *string
	LocationCountry   *string
	Phone             *string
	LinkedInURL       *string
	GitHubURL         *string
	WebsiteURL        *string
	YearsOfExperience *float64
	IsOpenToWork      *bool
}

// UpsertSkillInput holds the data for creating or updating a skill.
type UpsertSkillInput struct {
	SkillName         string
	Category          SkillCategory
	Proficiency       SkillProficiency
	YearsOfExperience *float64
	IsPrimary         bool
	Source            string
	Confidence        *float64
	LastUsedYear      *int16
}

// CreateWorkExperienceInput holds the data for creating a work experience entry.
type CreateWorkExperienceInput struct {
	ResumeUploadID   *uuid.UUID
	CompanyName      string
	JobTitle         string
	EmploymentType   EmploymentType
	LocationType     *WorkLocationType
	Location         *string
	StartDate        time.Time
	EndDate          *time.Time
	IsCurrent        bool
	Description      *string
	Responsibilities []string
	TechnologiesUsed []string
	Confidence       *float64
}

// CreateEducationInput holds the data for creating an education entry.
type CreateEducationInput struct {
	ResumeUploadID  *uuid.UUID
	InstitutionName string
	DegreeLevel     DegreeLevel
	DegreeName      *string
	FieldOfStudy    *string
	StartDate       *time.Time
	EndDate         *time.Time
	IsCurrent       bool
	GPA             *float64
	GPAScale        float64
	Honors          *string
	Activities      []string
	Confidence      *float64
}
