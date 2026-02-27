# LearnBot User Profile Data Model Documentation

## Overview

The LearnBot user profile data model is a PostgreSQL relational schema designed for the AI-powered career development platform. It stores user accounts, parsed resume data with versioning, skills with proficiency levels, work experience with duration calculations, education with degree levels, and user preferences and career goals.

The schema is optimized for **read-heavy operations** through strategic indexing, materialized views, and denormalized computed columns.

---

## Design Principles

1. **UUID primary keys** — All entities use `gen_random_uuid()` for globally unique, non-sequential IDs that are safe to expose in APIs
2. **Soft deletes** — Users are deactivated (`is_active = FALSE`) rather than deleted to preserve referential integrity
3. **Audit trail** — All profile changes are recorded in `profile_history` with before/after JSONB snapshots
4. **Computed columns** — `duration_months` (work experience) and `is_expired` (certifications) are PostgreSQL `GENERATED ALWAYS AS ... STORED` columns
5. **Enum types** — Custom PostgreSQL enums enforce valid values at the database level
6. **Array columns** — PostgreSQL native arrays (`TEXT[]`) for multi-value fields like skills, responsibilities, and technologies
7. **JSONB** — Used for flexible audit log snapshots in `profile_history`

---

## Tables

### `users`
Core user account information. Separate from profile data to support multiple authentication methods.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK, DEFAULT gen_random_uuid() | Unique user identifier |
| `email` | TEXT | UNIQUE, NOT NULL, CHECK format | Email address (lowercase) |
| `email_verified` | BOOLEAN | NOT NULL DEFAULT FALSE | Email verification status |
| `password_hash` | TEXT | NULL | Bcrypt hash; NULL for OAuth-only accounts |
| `full_name` | TEXT | NOT NULL, CHECK not empty | Display name |
| `avatar_url` | TEXT | NULL | Profile picture URL |
| `timezone` | TEXT | NOT NULL DEFAULT 'UTC' | IANA timezone string |
| `locale` | TEXT | NOT NULL DEFAULT 'en' | BCP 47 locale code |
| `is_active` | BOOLEAN | NOT NULL DEFAULT TRUE | Soft delete flag |
| `is_admin` | BOOLEAN | NOT NULL DEFAULT FALSE | Admin access flag |
| `last_login_at` | TIMESTAMPTZ | NULL | Last successful login |
| `created_at` | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Account creation time |
| `updated_at` | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Last update (auto-trigger) |

**Indexes:** `email`, `is_active` (partial), `created_at DESC`

---

### `user_profiles`
Professional profile data. One-to-one with `users`.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK | Profile identifier |
| `user_id` | UUID | FK users, UNIQUE | Owner |
| `headline` | TEXT | NULL | Professional headline |
| `summary` | TEXT | NULL | Professional summary |
| `location_city` | TEXT | NULL | City |
| `location_state` | TEXT | NULL | State/Province |
| `location_country` | TEXT | NOT NULL DEFAULT 'US' | Country code |
| `phone` | TEXT | NULL | Phone number |
| `linkedin_url` | TEXT | NULL, CHECK format | LinkedIn profile URL |
| `github_url` | TEXT | NULL, CHECK format | GitHub profile URL |
| `website_url` | TEXT | NULL | Personal website |
| `years_of_experience` | NUMERIC(4,1) | NULL, CHECK ≥ 0 | Total years (computed or manual) |
| `is_open_to_work` | BOOLEAN | NOT NULL DEFAULT FALSE | Job search status |
| `profile_completeness` | SMALLINT | NOT NULL DEFAULT 0, CHECK 0-100 | Completeness score (%) |
| `created_at` | TIMESTAMPTZ | NOT NULL | Creation time |
| `updated_at` | TIMESTAMPTZ | NOT NULL | Last update |

**Indexes:** `user_id`, `location`, `is_open_to_work` (partial), `updated_at DESC`

**Functions:**
- `calculate_profile_completeness(user_id)` — Returns 0-100 score based on filled sections
- `calculate_years_of_experience(user_id)` — Sums `duration_months` from work_experience

---

### `resume_uploads`
Versioned resume file uploads. Supports multiple versions per user.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK | Upload identifier |
| `user_id` | UUID | FK users | Owner |
| `file_name` | TEXT | NOT NULL | Original filename |
| `file_type` | TEXT | NOT NULL, CHECK 'pdf'/'docx' | File format |
| `file_size_bytes` | INTEGER | NOT NULL, CHECK > 0 | File size |
| `storage_key` | TEXT | NOT NULL | S3/GCS object key |
| `version` | INTEGER | NOT NULL DEFAULT 1, CHECK > 0 | Version number |
| `is_current` | BOOLEAN | NOT NULL DEFAULT TRUE | Active version flag |
| `parsed_at` | TIMESTAMPTZ | NULL | When parsing completed |
| `parse_status` | TEXT | NOT NULL DEFAULT 'pending' | pending/processing/done/failed |
| `parse_error` | TEXT | NULL | Error message if failed |
| `raw_text` | TEXT | NULL | Extracted text content |
| `parser_version` | TEXT | NULL | Parser version used |
| `overall_confidence` | NUMERIC(4,3) | NULL, CHECK 0-1 | Parser confidence score |
| `created_at` | TIMESTAMPTZ | NOT NULL | Upload time |

**Trigger:** `enforce_single_current_resume` — Ensures only one `is_current = TRUE` per user

---

### `user_skills`
Skills associated with a user profile, with proficiency levels.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK | Skill entry identifier |
| `user_id` | UUID | FK users | Owner |
| `skill_taxonomy_id` | UUID | FK skill_taxonomy, NULL | Canonical taxonomy reference |
| `skill_name` | TEXT | NOT NULL | Display name |
| `normalized_name` | TEXT | NOT NULL | Lowercase, trimmed |
| `category` | skill_category | NOT NULL DEFAULT 'other' | technical/soft/language/tool/etc. |
| `proficiency` | skill_proficiency | NOT NULL DEFAULT 'beginner' | beginner/intermediate/advanced/expert |
| `years_of_experience` | NUMERIC(4,1) | NULL, CHECK ≥ 0 | Years using this skill |
| `is_primary` | BOOLEAN | NOT NULL DEFAULT FALSE | Highlighted skill |
| `source` | TEXT | NOT NULL DEFAULT 'manual' | resume_parser/manual/assessment |
| `confidence` | NUMERIC(4,3) | NULL, CHECK 0-1 | Parser confidence |
| `last_used_year` | SMALLINT | NULL, CHECK range | Last year skill was used |
| `created_at` | TIMESTAMPTZ | NOT NULL | Creation time |
| `updated_at` | TIMESTAMPTZ | NOT NULL | Last update |

**Unique constraint:** `(user_id, normalized_name)` — One entry per skill per user

**Indexes:** `user_id`, `skill_taxonomy_id`, `proficiency`, `category`, `normalized_name`, GIN full-text on `skill_name`

---

### `work_experience`
Job history entries with automatic duration calculation.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK | Entry identifier |
| `user_id` | UUID | FK users | Owner |
| `resume_upload_id` | UUID | FK resume_uploads, NULL | Source upload |
| `company_name` | TEXT | NOT NULL, CHECK not empty | Employer name |
| `job_title` | TEXT | NOT NULL, CHECK not empty | Position title |
| `employment_type` | employment_type | NOT NULL DEFAULT 'full_time' | full_time/part_time/contract/etc. |
| `location_type` | work_location_type | NULL | on_site/remote/hybrid |
| `location` | TEXT | NULL | City, Country |
| `start_date` | DATE | NOT NULL | Start date |
| `end_date` | DATE | NULL | End date (NULL = current) |
| `is_current` | BOOLEAN | NOT NULL DEFAULT FALSE | Current position flag |
| `description` | TEXT | NULL | Role description |
| `responsibilities` | TEXT[] | NULL | Bullet point list |
| `technologies_used` | TEXT[] | NULL | Skills used in role |
| `duration_months` | INTEGER | GENERATED ALWAYS AS STORED | Auto-computed duration |
| `confidence` | NUMERIC(4,3) | NULL, CHECK 0-1 | Parser confidence |
| `display_order` | SMALLINT | NOT NULL DEFAULT 0 | Manual sort order |
| `created_at` | TIMESTAMPTZ | NOT NULL | Creation time |
| `updated_at` | TIMESTAMPTZ | NOT NULL | Last update |

**Check constraints:**
- `end_date >= start_date` (when both present)
- `NOT (is_current = TRUE AND end_date IS NOT NULL)`

**Duration formula:**
```sql
CASE
    WHEN end_date IS NOT NULL
        THEN (YEAR(end_date) - YEAR(start_date)) * 12 + (MONTH(end_date) - MONTH(start_date))
    ELSE (YEAR(NOW()) - YEAR(start_date)) * 12 + (MONTH(NOW()) - MONTH(start_date))
END
```

---

### `education`
Educational background with degree levels.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK | Entry identifier |
| `user_id` | UUID | FK users | Owner |
| `institution_name` | TEXT | NOT NULL, CHECK not empty | School/university name |
| `degree_level` | degree_level | NOT NULL DEFAULT 'bachelor' | high_school/associate/bachelor/master/doctorate/etc. |
| `degree_name` | TEXT | NULL | Full degree name |
| `field_of_study` | TEXT | NULL | Major/concentration |
| `start_date` | DATE | NULL | Enrollment date |
| `end_date` | DATE | NULL | Graduation date |
| `is_current` | BOOLEAN | NOT NULL DEFAULT FALSE | Currently enrolled |
| `gpa` | NUMERIC(3,2) | NULL, CHECK 0 ≤ gpa ≤ gpa_scale | Grade point average |
| `gpa_scale` | NUMERIC(3,1) | NOT NULL DEFAULT 4.0, CHECK > 0 | GPA scale (4.0, 5.0, 10.0) |
| `honors` | TEXT | NULL | Graduation honors |
| `activities` | TEXT[] | NULL | Clubs, activities |
| `confidence` | NUMERIC(4,3) | NULL, CHECK 0-1 | Parser confidence |

---

### `certifications`
Professional certifications and licenses.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK | Entry identifier |
| `user_id` | UUID | FK users | Owner |
| `name` | TEXT | NOT NULL, CHECK not empty | Certification name |
| `issuing_organization` | TEXT | NULL | Issuer name |
| `issue_date` | DATE | NULL | Issue date |
| `expiry_date` | DATE | NULL | Expiry date |
| `credential_id` | TEXT | NULL | Credential ID |
| `credential_url` | TEXT | NULL | Verification URL |
| `is_expired` | BOOLEAN | GENERATED ALWAYS AS STORED | Auto-computed from expiry_date |
| `confidence` | NUMERIC(4,3) | NULL, CHECK 0-1 | Parser confidence |

---

### `user_preferences`
Job search and career preferences. One-to-one with `users`.

Key fields:
- `desired_job_titles TEXT[]` — Target roles
- `desired_industries TEXT[]` — Target industries
- `salary_min / salary_max INTEGER` — Annual salary range
- `job_search_urgency TEXT` — passive/active/urgent/not_looking
- `career_stage TEXT` — entry/mid/senior/lead/principal/executive

---

### `career_goals`
Specific career objectives with progress tracking.

| Column | Type | Description |
|--------|------|-------------|
| `status` | goal_status | active/achieved/paused/abandoned |
| `priority` | SMALLINT | 1 = highest priority |
| `progress_percentage` | SMALLINT | 0-100% |
| `achieved_at` | TIMESTAMPTZ | Auto-set when progress = 100% |

---

### `skill_gaps`
Identified gaps between current skills and target role requirements.

| Column | Type | Description |
|--------|------|-------------|
| `gap_type` | TEXT | missing/needs_improvement |
| `required_proficiency` | skill_proficiency | Required level for target role |
| `current_proficiency` | skill_proficiency | Current level (NULL = not present) |
| `importance` | TEXT | critical/important/nice_to_have |
| `is_addressed` | BOOLEAN | Whether gap has been resolved |

---

### `profile_history`
Immutable audit log of all profile changes.

| Column | Type | Description |
|--------|------|-------------|
| `event_type` | profile_event_type | Type of change event |
| `entity_type` | TEXT | Which table was changed |
| `entity_id` | UUID | ID of changed record |
| `old_data` | JSONB | Snapshot before change |
| `new_data` | JSONB | Snapshot after change |
| `changed_by` | UUID | User who made the change (NULL = system) |

---

## Enum Types

### `skill_proficiency`
| Value | Description |
|-------|-------------|
| `beginner` | 0-1 year, basic awareness |
| `intermediate` | 1-3 years, works independently |
| `advanced` | 3-5 years, deep knowledge |
| `expert` | 5+ years, can teach others |

### `degree_level`
`high_school`, `associate`, `bachelor`, `master`, `doctorate`, `professional`, `certificate`, `diploma`, `other`

### `skill_category`
`technical`, `soft`, `language`, `tool`, `framework`, `database`, `cloud`, `other`

### `employment_type`
`full_time`, `part_time`, `contract`, `freelance`, `internship`, `volunteer`, `other`

### `work_location_type`
`on_site`, `remote`, `hybrid`

### `goal_status`
`active`, `achieved`, `paused`, `abandoned`

---

## Views

### `v_user_full_profile`
Complete profile with all sections as JSONB arrays. Used by the RAG pipeline.

### `v_user_skills_summary`
Aggregated skill statistics per user (counts by proficiency and category).

### `v_work_experience_summary`
Total experience months/years, company list, all technologies used.

### `v_active_skill_gaps`
Unaddressed skill gaps with career goal context, ordered by importance.

### `v_profile_completeness_breakdown`
Per-section completeness flags for UI display.

---

## Indexing Strategy

The schema is optimized for these common read patterns:

| Query Pattern | Index |
|---|---|
| Get profile by user ID | `idx_profiles_user_id` |
| Find open-to-work users | `idx_profiles_open_to_work` (partial) |
| Get current resume | `idx_resume_uploads_current` (partial) |
| Get skills by proficiency | `idx_user_skills_proficiency` |
| Search skills by name | `idx_user_skills_name_fts` (GIN full-text) |
| Get current job | `idx_work_exp_current` (partial) |
| Find expiring certifications | `idx_certifications_expiry` |
| Get unaddressed skill gaps | `idx_skill_gaps_unaddressed` (partial) |
| Search skill taxonomy | `idx_skill_taxonomy_aliases` (GIN array) |

---

## Migration Execution Order

```bash
psql -d learnbot -f migrations/001_create_enums.sql
psql -d learnbot -f migrations/002_create_users.sql
psql -d learnbot -f migrations/003_create_profiles.sql
psql -d learnbot -f migrations/004_create_skills_experience_education.sql
psql -d learnbot -f migrations/005_create_preferences_goals.sql
psql -d learnbot -f migrations/006_create_views.sql
```

Or using a migration tool:
```bash
# Using golang-migrate
migrate -path ./migrations -database "postgres://user:pass@localhost/learnbot?sslmode=disable" up
```
