# Database

Shared PostgreSQL repository layer for the LearnBot platform. Contains database migrations, repository implementations, and data models used by multiple services.

Part of the [LearnBot AI Career Development Platform](../README.md).

## Overview

This module provides:
- **SQL Migrations** — Ordered migration files to set up the complete schema
- **Repository Layer** — Go interfaces and implementations for data access
- **Data Models** — Go structs representing database entities
- **Documentation** — Data model reference and ER diagram

## Directory Structure

```
database/
├── migrations/
│   ├── 001_create_enums.sql              # PostgreSQL enum types
│   ├── 002_create_users.sql              # users table
│   ├── 003_create_profiles.sql           # user_profiles, resume_uploads
│   ├── 004_create_skills_experience_education.sql  # user_skills, work_experience, education, certifications
│   ├── 005_create_preferences_goals.sql  # user_preferences, career_goals, skill_gaps
│   ├── 006_create_views.sql              # Materialized views
│   ├── 007_create_learning_resources.sql # learning_resources, company_career_pages
│   └── 008_seed_learning_resources.sql   # 60+ curated learning resources
├── repository/
│   ├── models.go                         # Core data models (User, Profile, Skill, etc.)
│   ├── user_repository.go                # User and profile CRUD
│   ├── skill_repository.go               # Skills management
│   ├── experience_repository.go          # Work experience management
│   ├── preferences_repository.go         # User preferences and career goals
│   ├── learning_resource_models.go       # Learning resource models
│   ├── learning_resource_repository.go   # Learning resource queries
│   └── learning_resource_models_test.go  # Repository tests
└── docs/
    ├── DATA_MODEL.md                     # Full schema documentation
    └── ER_DIAGRAM.md                     # Entity relationship diagram
```

## Running Migrations

### Using psql

```bash
# Run all migrations in order
for f in migrations/*.sql; do
  psql "postgres://learnbot_admin:password@localhost:5432/learnbot" -f "$f"
done
```

### Using golang-migrate

```bash
migrate -path ./migrations \
  -database "postgres://learnbot_admin:password@localhost:5432/learnbot?sslmode=disable" \
  up
```

### Via Docker Compose

Migrations run automatically on first startup:
```bash
docker compose up -d postgres
# Migrations are applied via docker-entrypoint-initdb.d
```

### Migration Order

| File | Description |
|------|-------------|
| `001_create_enums.sql` | PostgreSQL enum types (skill_proficiency, degree_level, etc.) |
| `002_create_users.sql` | Core users table with auth fields |
| `003_create_profiles.sql` | User profiles and resume uploads |
| `004_create_skills_experience_education.sql` | Skills, work experience, education, certifications |
| `005_create_preferences_goals.sql` | Job preferences, career goals, skill gaps, audit log |
| `006_create_views.sql` | Materialized views for common query patterns |
| `007_create_learning_resources.sql` | Learning resources catalog and company career pages |
| `008_seed_learning_resources.sql` | Seed data: 60+ curated learning resources |

## Schema Overview

### Core Tables

| Table | Description |
|-------|-------------|
| `users` | Authentication and identity |
| `user_profiles` | Professional profile data |
| `resume_uploads` | Versioned resume files |
| `user_skills` | Skills with proficiency levels |
| `work_experience` | Job history with auto-computed duration |
| `education` | Educational background |
| `certifications` | Professional certifications |
| `user_preferences` | Job search preferences |
| `career_goals` | Career objectives with progress tracking |
| `skill_gaps` | Identified skill gaps |
| `profile_history` | Immutable audit log |
| `learning_resources` | Curated learning resource catalog |
| `company_career_pages` | Configurable career page scrapers |

### Views

| View | Description |
|------|-------------|
| `v_user_full_profile` | Complete profile as JSONB (for RAG pipeline) |
| `v_user_skills_summary` | Aggregated skill statistics |
| `v_work_experience_summary` | Total experience and technologies |
| `v_active_skill_gaps` | Unaddressed gaps ordered by importance |
| `v_profile_completeness_breakdown` | Per-section completeness flags |

For full schema documentation, see [docs/DATA_MODEL.md](docs/DATA_MODEL.md).

## Repository Pattern

The repository layer uses Go interfaces for testability:

```go
// Example: UserRepository interface
type UserRepository interface {
    CreateUser(ctx context.Context, user *User) error
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
    UpdateUser(ctx context.Context, user *User) error
}
```

### Using the Repository

```go
import "github.com/learnbot/database/repository"

// Create a repository
repo := repository.NewUserRepository(db)

// Get user by email
user, err := repo.GetUserByEmail(ctx, "jane@example.com")
if err != nil {
    return fmt.Errorf("get user: %w", err)
}
```

## Running Tests

```bash
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Design Principles

1. **UUID primary keys** — All entities use `gen_random_uuid()` for globally unique IDs
2. **Soft deletes** — Users are deactivated (`is_active = FALSE`) rather than deleted
3. **Audit trail** — All profile changes recorded in `profile_history` with JSONB snapshots
4. **Computed columns** — `duration_months` and `is_expired` are PostgreSQL generated columns
5. **Enum types** — Custom PostgreSQL enums enforce valid values at the database level
6. **Array columns** — Native PostgreSQL arrays for multi-value fields
7. **JSONB** — Flexible storage for audit logs and full profile snapshots

## Database Connection

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

dsn := fmt.Sprintf(
    "host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
)

db, err := sql.Open("postgres", dsn)
if err != nil {
    log.Fatal(err)
}
```
