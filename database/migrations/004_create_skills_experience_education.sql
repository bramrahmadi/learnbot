-- Migration 004: Create skills, work experience, and education tables

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- skill_taxonomy: Master list of canonical skills (shared across all users)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE skill_taxonomy (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    normalized_name TEXT NOT NULL,                     -- lowercase, trimmed
    category        skill_category NOT NULL DEFAULT 'other',
    aliases         TEXT[],                            -- alternative names
    parent_skill_id UUID REFERENCES skill_taxonomy(id), -- for skill hierarchies
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,    -- curated by admins
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT skill_taxonomy_name_unique UNIQUE (normalized_name)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- user_skills: Skills associated with a user profile
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE user_skills (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_taxonomy_id   UUID REFERENCES skill_taxonomy(id),
    skill_name          TEXT NOT NULL,                 -- raw name (may differ from taxonomy)
    normalized_name     TEXT NOT NULL,
    category            skill_category NOT NULL DEFAULT 'other',
    proficiency         skill_proficiency NOT NULL DEFAULT 'beginner',
    years_of_experience NUMERIC(4,1),
    is_primary          BOOLEAN NOT NULL DEFAULT FALSE, -- highlighted skill
    source              TEXT NOT NULL DEFAULT 'manual', -- 'resume_parser', 'manual', 'assessment'
    confidence          NUMERIC(4,3),                  -- parser confidence (0-1)
    last_used_year      SMALLINT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_skills_unique UNIQUE (user_id, normalized_name),
    CONSTRAINT user_skills_years_positive CHECK (years_of_experience IS NULL OR years_of_experience >= 0),
    CONSTRAINT user_skills_confidence_range CHECK (confidence IS NULL OR confidence BETWEEN 0 AND 1),
    CONSTRAINT user_skills_last_used_year CHECK (
        last_used_year IS NULL OR last_used_year BETWEEN 1970 AND EXTRACT(YEAR FROM NOW())::SMALLINT + 1
    ),
    CONSTRAINT user_skills_source_valid CHECK (
        source IN ('resume_parser', 'manual', 'assessment', 'linkedin_import')
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- work_experience: Job history entries
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE work_experience (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resume_upload_id    UUID REFERENCES resume_uploads(id) ON DELETE SET NULL,
    company_name        TEXT NOT NULL,
    job_title           TEXT NOT NULL,
    employment_type     employment_type NOT NULL DEFAULT 'full_time',
    location_type       work_location_type,
    location            TEXT,
    start_date          DATE NOT NULL,
    end_date            DATE,                          -- NULL = current position
    is_current          BOOLEAN NOT NULL DEFAULT FALSE,
    description         TEXT,
    responsibilities    TEXT[],                        -- bullet points
    technologies_used   TEXT[],                        -- skills used in this role
    -- Computed column: duration in months (updated by trigger)
    duration_months     INTEGER GENERATED ALWAYS AS (
        CASE
            WHEN end_date IS NOT NULL
                THEN (EXTRACT(YEAR FROM end_date) - EXTRACT(YEAR FROM start_date)) * 12
                   + (EXTRACT(MONTH FROM end_date) - EXTRACT(MONTH FROM start_date))
            ELSE (EXTRACT(YEAR FROM NOW()) - EXTRACT(YEAR FROM start_date)) * 12
               + (EXTRACT(MONTH FROM NOW()) - EXTRACT(MONTH FROM start_date))
        END
    ) STORED,
    confidence          NUMERIC(4,3),
    display_order       SMALLINT NOT NULL DEFAULT 0,   -- for manual ordering
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT work_exp_dates_valid CHECK (end_date IS NULL OR end_date >= start_date),
    CONSTRAINT work_exp_current_no_end CHECK (
        NOT (is_current = TRUE AND end_date IS NOT NULL)
    ),
    CONSTRAINT work_exp_confidence_range CHECK (confidence IS NULL OR confidence BETWEEN 0 AND 1),
    CONSTRAINT work_exp_company_not_empty CHECK (LENGTH(TRIM(company_name)) > 0),
    CONSTRAINT work_exp_title_not_empty CHECK (LENGTH(TRIM(job_title)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- education: Educational background entries
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE education (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resume_upload_id    UUID REFERENCES resume_uploads(id) ON DELETE SET NULL,
    institution_name    TEXT NOT NULL,
    degree_level        degree_level NOT NULL DEFAULT 'bachelor',
    degree_name         TEXT,                          -- e.g. "Bachelor of Science"
    field_of_study      TEXT,                          -- e.g. "Computer Science"
    start_date          DATE,
    end_date            DATE,
    is_current          BOOLEAN NOT NULL DEFAULT FALSE,
    gpa                 NUMERIC(3,2),
    gpa_scale           NUMERIC(3,1) NOT NULL DEFAULT 4.0,
    honors              TEXT,                          -- "Magna Cum Laude", etc.
    activities          TEXT[],
    confidence          NUMERIC(4,3),
    display_order       SMALLINT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT education_dates_valid CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    CONSTRAINT education_gpa_range CHECK (gpa IS NULL OR (gpa >= 0 AND gpa <= gpa_scale)),
    CONSTRAINT education_gpa_scale_positive CHECK (gpa_scale > 0),
    CONSTRAINT education_confidence_range CHECK (confidence IS NULL OR confidence BETWEEN 0 AND 1),
    CONSTRAINT education_institution_not_empty CHECK (LENGTH(TRIM(institution_name)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- certifications: Professional certifications and licenses
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE certifications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resume_upload_id    UUID REFERENCES resume_uploads(id) ON DELETE SET NULL,
    name                TEXT NOT NULL,
    issuing_organization TEXT,
    issue_date          DATE,
    expiry_date         DATE,
    credential_id       TEXT,
    credential_url      TEXT,
    is_expired          BOOLEAN GENERATED ALWAYS AS (
        expiry_date IS NOT NULL AND expiry_date < CURRENT_DATE
    ) STORED,
    confidence          NUMERIC(4,3),
    display_order       SMALLINT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT cert_dates_valid CHECK (expiry_date IS NULL OR issue_date IS NULL OR expiry_date >= issue_date),
    CONSTRAINT cert_confidence_range CHECK (confidence IS NULL OR confidence BETWEEN 0 AND 1),
    CONSTRAINT cert_name_not_empty CHECK (LENGTH(TRIM(name)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- projects: Portfolio projects and achievements
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE projects (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resume_upload_id    UUID REFERENCES resume_uploads(id) ON DELETE SET NULL,
    name                TEXT NOT NULL,
    description         TEXT,
    project_url         TEXT,
    repository_url      TEXT,
    technologies        TEXT[],
    start_date          DATE,
    end_date            DATE,
    is_ongoing          BOOLEAN NOT NULL DEFAULT FALSE,
    confidence          NUMERIC(4,3),
    display_order       SMALLINT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT project_dates_valid CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    CONSTRAINT project_confidence_range CHECK (confidence IS NULL OR confidence BETWEEN 0 AND 1),
    CONSTRAINT project_name_not_empty CHECK (LENGTH(TRIM(name)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- Indexes (optimized for read-heavy operations)
-- ─────────────────────────────────────────────────────────────────────────────

-- skill_taxonomy
CREATE INDEX idx_skill_taxonomy_normalized ON skill_taxonomy(normalized_name);
CREATE INDEX idx_skill_taxonomy_category ON skill_taxonomy(category);
CREATE INDEX idx_skill_taxonomy_aliases ON skill_taxonomy USING GIN(aliases);

-- user_skills
CREATE INDEX idx_user_skills_user_id ON user_skills(user_id);
CREATE INDEX idx_user_skills_taxonomy ON user_skills(skill_taxonomy_id);
CREATE INDEX idx_user_skills_proficiency ON user_skills(user_id, proficiency);
CREATE INDEX idx_user_skills_category ON user_skills(user_id, category);
CREATE INDEX idx_user_skills_normalized ON user_skills(normalized_name);
-- Full-text search on skill names
CREATE INDEX idx_user_skills_name_fts ON user_skills USING GIN(to_tsvector('english', skill_name));

-- work_experience
CREATE INDEX idx_work_exp_user_id ON work_experience(user_id);
CREATE INDEX idx_work_exp_current ON work_experience(user_id, is_current) WHERE is_current = TRUE;
CREATE INDEX idx_work_exp_dates ON work_experience(user_id, start_date DESC, end_date DESC NULLS FIRST);
CREATE INDEX idx_work_exp_company ON work_experience(company_name);
CREATE INDEX idx_work_exp_technologies ON work_experience USING GIN(technologies_used);

-- education
CREATE INDEX idx_education_user_id ON education(user_id);
CREATE INDEX idx_education_degree_level ON education(user_id, degree_level);
CREATE INDEX idx_education_dates ON education(user_id, end_date DESC NULLS FIRST);

-- certifications
CREATE INDEX idx_certifications_user_id ON certifications(user_id);
CREATE INDEX idx_certifications_active ON certifications(user_id, is_expired) WHERE is_expired = FALSE;
CREATE INDEX idx_certifications_expiry ON certifications(expiry_date) WHERE expiry_date IS NOT NULL;

-- projects
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_technologies ON projects USING GIN(technologies);

-- ─────────────────────────────────────────────────────────────────────────────
-- Triggers
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TRIGGER user_skills_updated_at
    BEFORE UPDATE ON user_skills
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER work_experience_updated_at
    BEFORE UPDATE ON work_experience
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER education_updated_at
    BEFORE UPDATE ON education
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER certifications_updated_at
    BEFORE UPDATE ON certifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ─────────────────────────────────────────────────────────────────────────────
-- Function: Recalculate profile completeness score
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION calculate_profile_completeness(p_user_id UUID)
RETURNS SMALLINT AS $$
DECLARE
    score INTEGER := 0;
    v_profile user_profiles%ROWTYPE;
BEGIN
    SELECT * INTO v_profile FROM user_profiles WHERE user_id = p_user_id;
    IF NOT FOUND THEN RETURN 0; END IF;

    -- Basic info (30 points)
    IF v_profile.headline IS NOT NULL AND LENGTH(v_profile.headline) > 0 THEN score := score + 10; END IF;
    IF v_profile.summary IS NOT NULL AND LENGTH(v_profile.summary) > 20 THEN score := score + 10; END IF;
    IF v_profile.location_city IS NOT NULL THEN score := score + 5; END IF;
    IF v_profile.linkedin_url IS NOT NULL THEN score := score + 5; END IF;

    -- Skills (20 points)
    IF (SELECT COUNT(*) FROM user_skills WHERE user_id = p_user_id) >= 5 THEN score := score + 20;
    ELSIF (SELECT COUNT(*) FROM user_skills WHERE user_id = p_user_id) >= 1 THEN score := score + 10;
    END IF;

    -- Work experience (25 points)
    IF (SELECT COUNT(*) FROM work_experience WHERE user_id = p_user_id) >= 2 THEN score := score + 25;
    ELSIF (SELECT COUNT(*) FROM work_experience WHERE user_id = p_user_id) >= 1 THEN score := score + 15;
    END IF;

    -- Education (15 points)
    IF (SELECT COUNT(*) FROM education WHERE user_id = p_user_id) >= 1 THEN score := score + 15; END IF;

    -- Certifications (5 points)
    IF (SELECT COUNT(*) FROM certifications WHERE user_id = p_user_id) >= 1 THEN score := score + 5; END IF;

    -- Projects (5 points)
    IF (SELECT COUNT(*) FROM projects WHERE user_id = p_user_id) >= 1 THEN score := score + 5; END IF;

    RETURN LEAST(score, 100)::SMALLINT;
END;
$$ LANGUAGE plpgsql;

-- Function: Calculate total years of experience from work_experience table
CREATE OR REPLACE FUNCTION calculate_years_of_experience(p_user_id UUID)
RETURNS NUMERIC AS $$
DECLARE
    total_months INTEGER;
BEGIN
    SELECT COALESCE(SUM(duration_months), 0)
    INTO total_months
    FROM work_experience
    WHERE user_id = p_user_id;

    RETURN ROUND(total_months::NUMERIC / 12, 1);
END;
$$ LANGUAGE plpgsql;

COMMIT;
