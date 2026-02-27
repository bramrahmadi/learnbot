-- Migration 003: Create user profiles and resume versioning
-- Stores the structured profile data extracted from resumes.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- user_profiles: The canonical profile for each user
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE user_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    headline            TEXT,                          -- e.g. "Senior Software Engineer"
    summary             TEXT,                          -- professional summary
    location_city       TEXT,
    location_state      TEXT,
    location_country    TEXT NOT NULL DEFAULT 'US',
    phone               TEXT,
    linkedin_url        TEXT,
    github_url          TEXT,
    website_url         TEXT,
    years_of_experience NUMERIC(4,1),                  -- computed or manually set
    is_open_to_work     BOOLEAN NOT NULL DEFAULT FALSE,
    profile_completeness SMALLINT NOT NULL DEFAULT 0,  -- 0-100 percentage
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT profiles_user_unique UNIQUE (user_id),
    CONSTRAINT profiles_years_exp_positive CHECK (years_of_experience IS NULL OR years_of_experience >= 0),
    CONSTRAINT profiles_completeness_range CHECK (profile_completeness BETWEEN 0 AND 100),
    CONSTRAINT profiles_linkedin_format CHECK (
        linkedin_url IS NULL OR linkedin_url ~* '^https?://(www\.)?linkedin\.com/'
    ),
    CONSTRAINT profiles_github_format CHECK (
        github_url IS NULL OR github_url ~* '^https?://(www\.)?github\.com/'
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- resume_uploads: Versioned resume file uploads
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE resume_uploads (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name       TEXT NOT NULL,
    file_type       TEXT NOT NULL,                     -- 'pdf' or 'docx'
    file_size_bytes INTEGER NOT NULL,
    storage_key     TEXT NOT NULL,                     -- S3/GCS object key
    version         INTEGER NOT NULL DEFAULT 1,
    is_current      BOOLEAN NOT NULL DEFAULT TRUE,
    parsed_at       TIMESTAMPTZ,
    parse_status    TEXT NOT NULL DEFAULT 'pending',   -- 'pending','processing','done','failed'
    parse_error     TEXT,
    raw_text        TEXT,                              -- extracted text (optional)
    parser_version  TEXT,
    overall_confidence NUMERIC(4,3),                   -- 0.000 - 1.000
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT resume_file_type_valid CHECK (file_type IN ('pdf', 'docx')),
    CONSTRAINT resume_file_size_positive CHECK (file_size_bytes > 0),
    CONSTRAINT resume_version_positive CHECK (version > 0),
    CONSTRAINT resume_confidence_range CHECK (
        overall_confidence IS NULL OR overall_confidence BETWEEN 0 AND 1
    ),
    CONSTRAINT resume_parse_status_valid CHECK (
        parse_status IN ('pending', 'processing', 'done', 'failed')
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- profile_history: Audit log of all profile changes
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE profile_history (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type      profile_event_type NOT NULL,
    entity_type     TEXT,                              -- 'skill', 'experience', 'education', etc.
    entity_id       UUID,                              -- ID of the changed entity
    old_data        JSONB,                             -- snapshot before change
    new_data        JSONB,                             -- snapshot after change
    changed_by      UUID REFERENCES users(id),         -- NULL = system (e.g. resume parser)
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ─────────────────────────────────────────────────────────────────────────────
-- Indexes
-- ─────────────────────────────────────────────────────────────────────────────
CREATE INDEX idx_profiles_user_id ON user_profiles(user_id);
CREATE INDEX idx_profiles_location ON user_profiles(location_country, location_state, location_city);
CREATE INDEX idx_profiles_open_to_work ON user_profiles(is_open_to_work) WHERE is_open_to_work = TRUE;
CREATE INDEX idx_profiles_updated_at ON user_profiles(updated_at DESC);

CREATE INDEX idx_resume_uploads_user_id ON resume_uploads(user_id);
CREATE INDEX idx_resume_uploads_current ON resume_uploads(user_id, is_current) WHERE is_current = TRUE;
CREATE INDEX idx_resume_uploads_parse_status ON resume_uploads(parse_status) WHERE parse_status IN ('pending', 'processing');

CREATE INDEX idx_profile_history_user_id ON profile_history(user_id);
CREATE INDEX idx_profile_history_created_at ON profile_history(created_at DESC);
CREATE INDEX idx_profile_history_entity ON profile_history(entity_type, entity_id);

-- Triggers
CREATE TRIGGER profiles_updated_at
    BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to enforce only one current resume per user
CREATE OR REPLACE FUNCTION enforce_single_current_resume()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_current = TRUE THEN
        UPDATE resume_uploads
        SET is_current = FALSE
        WHERE user_id = NEW.user_id
          AND id != NEW.id
          AND is_current = TRUE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER resume_single_current
    AFTER INSERT OR UPDATE OF is_current ON resume_uploads
    FOR EACH ROW
    WHEN (NEW.is_current = TRUE)
    EXECUTE FUNCTION enforce_single_current_resume();

COMMIT;
