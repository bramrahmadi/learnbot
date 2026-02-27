-- Migration 005: Create user preferences and career goals tables

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- user_preferences: Career and job search preferences
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE user_preferences (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Job search preferences
    desired_job_titles      TEXT[],                    -- e.g. ['Senior Engineer', 'Staff Engineer']
    desired_industries      TEXT[],                    -- e.g. ['fintech', 'healthcare']
    desired_company_sizes   TEXT[],                    -- 'startup', 'mid', 'enterprise'
    desired_employment_types employment_type[],
    desired_location_types  work_location_type[],
    desired_locations       TEXT[],                    -- city/country preferences
    is_willing_to_relocate  BOOLEAN NOT NULL DEFAULT FALSE,
    relocation_locations    TEXT[],

    -- Salary preferences
    salary_currency         TEXT NOT NULL DEFAULT 'USD',
    salary_min              INTEGER,                   -- annual, in currency units
    salary_max              INTEGER,
    include_equity          BOOLEAN NOT NULL DEFAULT FALSE,

    -- Career timeline
    job_search_urgency      TEXT NOT NULL DEFAULT 'passive',  -- 'passive', 'active', 'urgent'
    available_from          DATE,
    career_stage            TEXT,                      -- 'entry', 'mid', 'senior', 'lead', 'executive'

    -- Notification preferences
    email_job_alerts        BOOLEAN NOT NULL DEFAULT TRUE,
    email_weekly_digest     BOOLEAN NOT NULL DEFAULT TRUE,
    email_training_recs     BOOLEAN NOT NULL DEFAULT TRUE,

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT preferences_user_unique UNIQUE (user_id),
    CONSTRAINT preferences_salary_range CHECK (
        salary_min IS NULL OR salary_max IS NULL OR salary_min <= salary_max
    ),
    CONSTRAINT preferences_salary_positive CHECK (
        (salary_min IS NULL OR salary_min >= 0) AND
        (salary_max IS NULL OR salary_max >= 0)
    ),
    CONSTRAINT preferences_urgency_valid CHECK (
        job_search_urgency IN ('passive', 'active', 'urgent', 'not_looking')
    ),
    CONSTRAINT preferences_career_stage_valid CHECK (
        career_stage IS NULL OR
        career_stage IN ('entry', 'mid', 'senior', 'lead', 'principal', 'executive')
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- career_goals: Specific career objectives with tracking
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE career_goals (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title               TEXT NOT NULL,
    description         TEXT,
    target_role         TEXT,
    target_industry     TEXT,
    target_date         DATE,
    status              goal_status NOT NULL DEFAULT 'active',
    priority            SMALLINT NOT NULL DEFAULT 1,   -- 1=highest
    progress_percentage SMALLINT NOT NULL DEFAULT 0,   -- 0-100
    notes               TEXT,
    achieved_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT goals_priority_positive CHECK (priority > 0),
    CONSTRAINT goals_progress_range CHECK (progress_percentage BETWEEN 0 AND 100),
    CONSTRAINT goals_title_not_empty CHECK (LENGTH(TRIM(title)) > 0),
    CONSTRAINT goals_achieved_requires_status CHECK (
        achieved_at IS NULL OR status = 'achieved'
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- skill_gaps: Identified gaps between current skills and target role requirements
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE skill_gaps (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    career_goal_id      UUID REFERENCES career_goals(id) ON DELETE CASCADE,
    skill_name          TEXT NOT NULL,
    normalized_name     TEXT NOT NULL,
    skill_taxonomy_id   UUID REFERENCES skill_taxonomy(id),
    gap_type            TEXT NOT NULL DEFAULT 'missing', -- 'missing', 'needs_improvement'
    required_proficiency skill_proficiency NOT NULL DEFAULT 'intermediate',
    current_proficiency skill_proficiency,              -- NULL = not present
    importance          TEXT NOT NULL DEFAULT 'nice_to_have', -- 'critical', 'important', 'nice_to_have'
    is_addressed        BOOLEAN NOT NULL DEFAULT FALSE,
    identified_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    addressed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT skill_gaps_gap_type_valid CHECK (gap_type IN ('missing', 'needs_improvement')),
    CONSTRAINT skill_gaps_importance_valid CHECK (
        importance IN ('critical', 'important', 'nice_to_have')
    ),
    CONSTRAINT skill_gaps_addressed_check CHECK (
        addressed_at IS NULL OR is_addressed = TRUE
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- Indexes
-- ─────────────────────────────────────────────────────────────────────────────
CREATE INDEX idx_preferences_user_id ON user_preferences(user_id);
CREATE INDEX idx_preferences_job_titles ON user_preferences USING GIN(desired_job_titles);
CREATE INDEX idx_preferences_industries ON user_preferences USING GIN(desired_industries);
CREATE INDEX idx_preferences_locations ON user_preferences USING GIN(desired_locations);
CREATE INDEX idx_preferences_urgency ON user_preferences(job_search_urgency);

CREATE INDEX idx_career_goals_user_id ON career_goals(user_id);
CREATE INDEX idx_career_goals_status ON career_goals(user_id, status);
CREATE INDEX idx_career_goals_target_date ON career_goals(target_date) WHERE target_date IS NOT NULL;

CREATE INDEX idx_skill_gaps_user_id ON skill_gaps(user_id);
CREATE INDEX idx_skill_gaps_goal ON skill_gaps(career_goal_id);
CREATE INDEX idx_skill_gaps_unaddressed ON skill_gaps(user_id, is_addressed) WHERE is_addressed = FALSE;
CREATE INDEX idx_skill_gaps_importance ON skill_gaps(user_id, importance);

-- ─────────────────────────────────────────────────────────────────────────────
-- Triggers
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TRIGGER preferences_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER career_goals_updated_at
    BEFORE UPDATE ON career_goals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER skill_gaps_updated_at
    BEFORE UPDATE ON skill_gaps
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
