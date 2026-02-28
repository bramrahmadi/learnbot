// Package repository provides data access objects for the LearnBot platform.
// This file defines models for the learning resource catalog.
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

// ResourceType classifies the type of learning resource.
type ResourceType string

const (
	ResourceTypeCourse        ResourceType = "course"
	ResourceTypeCertification ResourceType = "certification"
	ResourceTypeDocumentation ResourceType = "documentation"
	ResourceTypeVideo         ResourceType = "video"
	ResourceTypeBook          ResourceType = "book"
	ResourceTypePractice      ResourceType = "practice"
	ResourceTypeArticle       ResourceType = "article"
	ResourceTypeProject       ResourceType = "project"
	ResourceTypeOther         ResourceType = "other"
)

// ResourceDifficulty classifies the difficulty level of a resource.
type ResourceDifficulty string

const (
	ResourceDifficultyBeginner     ResourceDifficulty = "beginner"
	ResourceDifficultyIntermediate ResourceDifficulty = "intermediate"
	ResourceDifficultyAdvanced     ResourceDifficulty = "advanced"
	ResourceDifficultyExpert       ResourceDifficulty = "expert"
	ResourceDifficultyAllLevels    ResourceDifficulty = "all_levels"
)

// ResourceCostType classifies the cost model of a resource.
type ResourceCostType string

const (
	ResourceCostFree             ResourceCostType = "free"
	ResourceCostFreemium         ResourceCostType = "freemium"
	ResourceCostPaid             ResourceCostType = "paid"
	ResourceCostSubscription     ResourceCostType = "subscription"
	ResourceCostFreeAudit        ResourceCostType = "free_audit"
	ResourceCostEmployerSponsored ResourceCostType = "employer_sponsored"
)

// UserResourceStatus tracks a user's progress status for a resource.
type UserResourceStatus string

const (
	UserResourceStatusSaved      UserResourceStatus = "saved"
	UserResourceStatusInProgress UserResourceStatus = "in_progress"
	UserResourceStatusCompleted  UserResourceStatus = "completed"
	UserResourceStatusAbandoned  UserResourceStatus = "abandoned"
)

// ─────────────────────────────────────────────────────────────────────────────
// Model structs
// ─────────────────────────────────────────────────────────────────────────────

// ResourceProvider represents a platform or organization that offers resources.
type ResourceProvider struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	Name           string         `db:"name" json:"name"`
	NormalizedName string         `db:"normalized_name" json:"normalized_name"`
	WebsiteURL     sql.NullString `db:"website_url" json:"website_url,omitempty"`
	LogoURL        sql.NullString `db:"logo_url" json:"logo_url,omitempty"`
	Description    sql.NullString `db:"description" json:"description,omitempty"`
	IsActive       bool           `db:"is_active" json:"is_active"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// LearningResource represents a single learning resource in the catalog.
type LearningResource struct {
	ID               uuid.UUID          `db:"id" json:"id"`
	Title            string             `db:"title" json:"title"`
	Slug             string             `db:"slug" json:"slug"`
	Description      sql.NullString     `db:"description" json:"description,omitempty"`
	URL              string             `db:"url" json:"url"`
	ProviderID       *uuid.UUID         `db:"provider_id" json:"provider_id,omitempty"`
	ResourceType     ResourceType       `db:"resource_type" json:"resource_type"`
	Difficulty       ResourceDifficulty `db:"difficulty" json:"difficulty"`
	CostType         ResourceCostType   `db:"cost_type" json:"cost_type"`
	CostAmount       sql.NullFloat64    `db:"cost_amount" json:"cost_amount,omitempty"`
	CostCurrency     string             `db:"cost_currency" json:"cost_currency"`
	DurationHours    sql.NullFloat64    `db:"duration_hours" json:"duration_hours,omitempty"`
	DurationLabel    sql.NullString     `db:"duration_label" json:"duration_label,omitempty"`
	Language         string             `db:"language" json:"language"`
	IsActive         bool               `db:"is_active" json:"is_active"`
	IsFeatured       bool               `db:"is_featured" json:"is_featured"`
	IsVerified       bool               `db:"is_verified" json:"is_verified"`
	HasCertificate   bool               `db:"has_certificate" json:"has_certificate"`
	HasHandsOn       bool               `db:"has_hands_on" json:"has_hands_on"`
	Rating           sql.NullFloat64    `db:"rating" json:"rating,omitempty"`
	RatingCount      int                `db:"rating_count" json:"rating_count"`
	EnrollmentCount  sql.NullInt32      `db:"enrollment_count" json:"enrollment_count,omitempty"`
	LastUpdatedDate  sql.NullTime       `db:"last_updated_date" json:"last_updated_date,omitempty"`
	CreatedAt        time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `db:"updated_at" json:"updated_at"`
}

// ResourceSkill represents a skill covered by a learning resource.
type ResourceSkill struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	ResourceID     uuid.UUID          `db:"resource_id" json:"resource_id"`
	SkillName      string             `db:"skill_name" json:"skill_name"`
	NormalizedName string             `db:"normalized_name" json:"normalized_name"`
	IsPrimary      bool               `db:"is_primary" json:"is_primary"`
	CoverageLevel  *ResourceDifficulty `db:"coverage_level" json:"coverage_level,omitempty"`
	CreatedAt      time.Time          `db:"created_at" json:"created_at"`
}

// ResourcePrerequisite represents a prerequisite for a learning resource.
type ResourcePrerequisite struct {
	ID                     uuid.UUID      `db:"id" json:"id"`
	ResourceID             uuid.UUID      `db:"resource_id" json:"resource_id"`
	PrerequisiteSkill      sql.NullString `db:"prerequisite_skill" json:"prerequisite_skill,omitempty"`
	PrerequisiteResourceID *uuid.UUID     `db:"prerequisite_resource_id" json:"prerequisite_resource_id,omitempty"`
	IsRequired             bool           `db:"is_required" json:"is_required"`
	Description            sql.NullString `db:"description" json:"description,omitempty"`
	CreatedAt              time.Time      `db:"created_at" json:"created_at"`
}

// LearningPath represents a curated sequence of resources for a learning goal.
type LearningPath struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	Title          string             `db:"title" json:"title"`
	Slug           string             `db:"slug" json:"slug"`
	Description    sql.NullString     `db:"description" json:"description,omitempty"`
	TargetRole     sql.NullString     `db:"target_role" json:"target_role,omitempty"`
	TargetSkill    sql.NullString     `db:"target_skill" json:"target_skill,omitempty"`
	Difficulty     ResourceDifficulty `db:"difficulty" json:"difficulty"`
	EstimatedHours sql.NullFloat64    `db:"estimated_hours" json:"estimated_hours,omitempty"`
	IsActive       bool               `db:"is_active" json:"is_active"`
	IsFeatured     bool               `db:"is_featured" json:"is_featured"`
	CreatedBy      *uuid.UUID         `db:"created_by" json:"created_by,omitempty"`
	CreatedAt      time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `db:"updated_at" json:"updated_at"`
}

// LearningPathResource represents a resource within a learning path.
type LearningPathResource struct {
	ID         uuid.UUID `db:"id" json:"id"`
	PathID     uuid.UUID `db:"path_id" json:"path_id"`
	ResourceID uuid.UUID `db:"resource_id" json:"resource_id"`
	StepOrder  int16     `db:"step_order" json:"step_order"`
	IsRequired bool      `db:"is_required" json:"is_required"`
	Notes      sql.NullString `db:"notes" json:"notes,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// UserResourceProgress tracks a user's progress through a resource.
type UserResourceProgress struct {
	ID                 uuid.UUID          `db:"id" json:"id"`
	UserID             uuid.UUID          `db:"user_id" json:"user_id"`
	ResourceID         uuid.UUID          `db:"resource_id" json:"resource_id"`
	Status             UserResourceStatus `db:"status" json:"status"`
	ProgressPercentage int16              `db:"progress_percentage" json:"progress_percentage"`
	StartedAt          sql.NullTime       `db:"started_at" json:"started_at,omitempty"`
	CompletedAt        sql.NullTime       `db:"completed_at" json:"completed_at,omitempty"`
	UserRating         sql.NullInt16      `db:"user_rating" json:"user_rating,omitempty"`
	UserNotes          sql.NullString     `db:"user_notes" json:"user_notes,omitempty"`
	CreatedAt          time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `db:"updated_at" json:"updated_at"`
}

// ResourceReview represents a user review for a learning resource.
type ResourceReview struct {
	ID           uuid.UUID      `db:"id" json:"id"`
	ResourceID   uuid.UUID      `db:"resource_id" json:"resource_id"`
	UserID       uuid.UUID      `db:"user_id" json:"user_id"`
	Rating       int16          `db:"rating" json:"rating"`
	Title        sql.NullString `db:"title" json:"title,omitempty"`
	Body         sql.NullString `db:"body" json:"body,omitempty"`
	IsVerified   bool           `db:"is_verified" json:"is_verified"`
	HelpfulCount int            `db:"helpful_count" json:"helpful_count"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Enriched view types (for API responses)
// ─────────────────────────────────────────────────────────────────────────────

// LearningResourceWithSkills is a resource enriched with its skills and provider.
type LearningResourceWithSkills struct {
	LearningResource
	ProviderName string         `db:"provider_name" json:"provider_name,omitempty"`
	ProviderURL  sql.NullString `db:"provider_url" json:"provider_url,omitempty"`
	Skills       pq.StringArray `db:"skills" json:"skills,omitempty"`
	SkillIDs     pq.StringArray `db:"skill_ids" json:"skill_ids,omitempty"`
}

// LearningPathWithResources is a learning path enriched with its resources.
type LearningPathWithResources struct {
	LearningPath
	ResourceCount         int `db:"resource_count" json:"resource_count"`
	RequiredResourceCount int `db:"required_resource_count" json:"required_resource_count"`
	// Resources is populated separately (not from DB view).
	Resources []LearningPathStep `json:"resources,omitempty"`
}

// LearningPathStep is a resource within a learning path with ordering info.
type LearningPathStep struct {
	StepOrder  int16             `json:"step_order"`
	IsRequired bool              `json:"is_required"`
	Notes      string            `json:"notes,omitempty"`
	Resource   LearningResource  `json:"resource"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Input/request types for CRUD operations
// ─────────────────────────────────────────────────────────────────────────────

// CreateResourceInput holds the data for creating a new learning resource.
type CreateResourceInput struct {
	Title           string
	Slug            string
	Description     *string
	URL             string
	ProviderID      *uuid.UUID
	ResourceType    ResourceType
	Difficulty      ResourceDifficulty
	CostType        ResourceCostType
	CostAmount      *float64
	CostCurrency    string
	DurationHours   *float64
	DurationLabel   *string
	Language        string
	HasCertificate  bool
	HasHandsOn      bool
	Skills          []ResourceSkillInput
}

// ResourceSkillInput holds skill data for a resource.
type ResourceSkillInput struct {
	SkillName     string
	IsPrimary     bool
	CoverageLevel *ResourceDifficulty
}

// UpdateResourceInput holds the data for updating a learning resource.
type UpdateResourceInput struct {
	Title          *string
	Description    *string
	URL            *string
	Difficulty     *ResourceDifficulty
	CostType       *ResourceCostType
	CostAmount     *float64
	DurationHours  *float64
	DurationLabel  *string
	IsActive       *bool
	IsFeatured     *bool
	HasCertificate *bool
	HasHandsOn     *bool
}

// ResourceQueryFilter holds filter parameters for resource queries.
type ResourceQueryFilter struct {
	// SkillName filters resources by skill name (case-insensitive).
	SkillName string

	// ResourceType filters by resource type.
	ResourceType ResourceType

	// Difficulty filters by difficulty level.
	Difficulty ResourceDifficulty

	// CostType filters by cost model.
	CostType ResourceCostType

	// IsFree returns only free resources when true.
	IsFree bool

	// HasCertificate returns only resources with certificates when true.
	HasCertificate bool

	// HasHandsOn returns only resources with hands-on exercises when true.
	HasHandsOn bool

	// MinRating filters resources with rating >= MinRating.
	MinRating float64

	// ProviderID filters by provider.
	ProviderID *uuid.UUID

	// SearchQuery is a full-text search query.
	SearchQuery string

	// Limit is the maximum number of results (default 20, max 100).
	Limit int

	// Offset is the pagination offset.
	Offset int
}

// UpsertProgressInput holds data for updating user resource progress.
type UpsertProgressInput struct {
	Status             UserResourceStatus
	ProgressPercentage int16
	UserRating         *int16
	UserNotes          *string
}

// CreateReviewInput holds data for creating a resource review.
type CreateReviewInput struct {
	Rating int16
	Title  *string
	Body   *string
}

// CreateProviderInput holds data for creating a resource provider.
type CreateProviderInput struct {
	Name        string
	WebsiteURL  *string
	LogoURL     *string
	Description *string
}

// CreateLearningPathInput holds data for creating a learning path.
type CreateLearningPathInput struct {
	Title          string
	Slug           string
	Description    *string
	TargetRole     *string
	TargetSkill    *string
	Difficulty     ResourceDifficulty
	EstimatedHours *float64
	Resources      []LearningPathResourceInput
}

// LearningPathResourceInput holds data for adding a resource to a path.
type LearningPathResourceInput struct {
	ResourceID uuid.UUID
	StepOrder  int16
	IsRequired bool
	Notes      *string
}
