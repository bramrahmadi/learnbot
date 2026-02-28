-- Migration 007: Create learning resources database
-- Implements the learning resource catalog for the LearnBot platform.
-- Supports courses, certifications, free resources, books, and practice platforms.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- Enum types for learning resources
-- ─────────────────────────────────────────────────────────────────────────────

-- Resource type classification
CREATE TYPE resource_type AS ENUM (
    'course',           -- Structured online course (Coursera, Udemy, edX, Pluralsight)
    'certification',    -- Official certification program (AWS, GCP, Azure, etc.)
    'documentation',    -- Official documentation or tutorial
    'video',            -- YouTube channel, video series
    'book',             -- Book or study guide
    'practice',         -- Practice platform (LeetCode, HackerRank, Exercism)
    'article',          -- Blog post, article, or written tutorial
    'project',          -- Hands-on project or workshop
    'other'
);

-- Difficulty level for resources
CREATE TYPE resource_difficulty AS ENUM (
    'beginner',
    'intermediate',
    'advanced',
    'expert',
    'all_levels'
);

-- Cost model for resources
CREATE TYPE resource_cost_type AS ENUM (
    'free',             -- Completely free
    'freemium',         -- Free tier with paid upgrades
    'paid',             -- One-time purchase
    'subscription',     -- Monthly/annual subscription
    'free_audit',       -- Free to audit, paid for certificate
    'employer_sponsored' -- Typically employer-paid
);

-- ─────────────────────────────────────────────────────────────────────────────
-- resource_providers: Platforms and organizations that offer resources
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE resource_providers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    website_url     TEXT,
    logo_url        TEXT,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT resource_providers_name_unique UNIQUE (normalized_name),
    CONSTRAINT resource_providers_name_not_empty CHECK (LENGTH(TRIM(name)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- learning_resources: Core resource catalog
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE learning_resources (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    slug                TEXT NOT NULL,                  -- URL-friendly identifier
    description         TEXT,
    url                 TEXT NOT NULL,
    provider_id         UUID REFERENCES resource_providers(id) ON DELETE SET NULL,
    resource_type       resource_type NOT NULL DEFAULT 'course',
    difficulty          resource_difficulty NOT NULL DEFAULT 'intermediate',
    cost_type           resource_cost_type NOT NULL DEFAULT 'paid',
    cost_amount         NUMERIC(10,2),                  -- NULL = free or variable
    cost_currency       TEXT NOT NULL DEFAULT 'USD',
    duration_hours      NUMERIC(6,1),                   -- Estimated completion hours
    duration_label      TEXT,                           -- Human-readable: "8 weeks", "40 hours"
    language            TEXT NOT NULL DEFAULT 'en',
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    is_featured         BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified         BOOLEAN NOT NULL DEFAULT FALSE,  -- Curated by admin
    has_certificate     BOOLEAN NOT NULL DEFAULT FALSE,  -- Offers completion certificate
    has_hands_on        BOOLEAN NOT NULL DEFAULT FALSE,  -- Includes practical exercises
    rating              NUMERIC(3,2),                   -- Average rating 0.00-5.00
    rating_count        INTEGER NOT NULL DEFAULT 0,
    enrollment_count    INTEGER,                        -- Approximate enrollment
    last_updated_date   DATE,                           -- When content was last updated
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT learning_resources_slug_unique UNIQUE (slug),
    CONSTRAINT learning_resources_title_not_empty CHECK (LENGTH(TRIM(title)) > 0),
    CONSTRAINT learning_resources_url_not_empty CHECK (LENGTH(TRIM(url)) > 0),
    CONSTRAINT learning_resources_rating_range CHECK (rating IS NULL OR rating BETWEEN 0 AND 5),
    CONSTRAINT learning_resources_rating_count_positive CHECK (rating_count >= 0),
    CONSTRAINT learning_resources_cost_positive CHECK (cost_amount IS NULL OR cost_amount >= 0),
    CONSTRAINT learning_resources_duration_positive CHECK (duration_hours IS NULL OR duration_hours > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- resource_skills: Many-to-many mapping of resources to skills covered
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE resource_skills (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id     UUID NOT NULL REFERENCES learning_resources(id) ON DELETE CASCADE,
    skill_name      TEXT NOT NULL,                      -- Canonical skill name
    normalized_name TEXT NOT NULL,                      -- Lowercase, trimmed
    is_primary      BOOLEAN NOT NULL DEFAULT FALSE,     -- Main skill taught
    coverage_level  resource_difficulty,                -- Level covered for this skill
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT resource_skills_unique UNIQUE (resource_id, normalized_name),
    CONSTRAINT resource_skills_skill_not_empty CHECK (LENGTH(TRIM(skill_name)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- resource_prerequisites: Skills or resources required before taking a resource
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE resource_prerequisites (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id             UUID NOT NULL REFERENCES learning_resources(id) ON DELETE CASCADE,
    prerequisite_skill      TEXT,                       -- Skill name prerequisite
    prerequisite_resource_id UUID REFERENCES learning_resources(id) ON DELETE SET NULL,
    is_required             BOOLEAN NOT NULL DEFAULT TRUE, -- FALSE = recommended
    description             TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT resource_prereqs_has_content CHECK (
        prerequisite_skill IS NOT NULL OR prerequisite_resource_id IS NOT NULL
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- learning_paths: Curated sequences of resources for a learning goal
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE learning_paths (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    slug                TEXT NOT NULL,
    description         TEXT,
    target_role         TEXT,                           -- e.g. "Backend Engineer"
    target_skill        TEXT,                           -- e.g. "Python"
    difficulty          resource_difficulty NOT NULL DEFAULT 'intermediate',
    estimated_hours     NUMERIC(6,1),                   -- Total hours for the path
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    is_featured         BOOLEAN NOT NULL DEFAULT FALSE,
    created_by          UUID,                           -- Admin user ID (nullable)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT learning_paths_slug_unique UNIQUE (slug),
    CONSTRAINT learning_paths_title_not_empty CHECK (LENGTH(TRIM(title)) > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- learning_path_resources: Ordered resources within a learning path
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE learning_path_resources (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path_id         UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    resource_id     UUID NOT NULL REFERENCES learning_resources(id) ON DELETE CASCADE,
    step_order      SMALLINT NOT NULL,                  -- 1-based ordering
    is_required     BOOLEAN NOT NULL DEFAULT TRUE,
    notes           TEXT,                               -- Why this resource is in the path
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT learning_path_resources_unique UNIQUE (path_id, resource_id),
    CONSTRAINT learning_path_resources_order_positive CHECK (step_order > 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- user_resource_progress: Track user progress through resources
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE user_resource_progress (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_id         UUID NOT NULL REFERENCES learning_resources(id) ON DELETE CASCADE,
    status              TEXT NOT NULL DEFAULT 'saved',  -- saved/in_progress/completed/abandoned
    progress_percentage SMALLINT NOT NULL DEFAULT 0,
    started_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    user_rating         SMALLINT,                       -- User's personal rating 1-5
    user_notes          TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_resource_progress_unique UNIQUE (user_id, resource_id),
    CONSTRAINT user_resource_progress_pct_range CHECK (progress_percentage BETWEEN 0 AND 100),
    CONSTRAINT user_resource_progress_rating_range CHECK (user_rating IS NULL OR user_rating BETWEEN 1 AND 5),
    CONSTRAINT user_resource_progress_status_valid CHECK (
        status IN ('saved', 'in_progress', 'completed', 'abandoned')
    )
);

-- ─────────────────────────────────────────────────────────────────────────────
-- resource_reviews: User reviews and ratings for resources
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE resource_reviews (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id     UUID NOT NULL REFERENCES learning_resources(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating          SMALLINT NOT NULL,
    title           TEXT,
    body            TEXT,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,     -- Verified completion
    helpful_count   INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT resource_reviews_unique UNIQUE (resource_id, user_id),
    CONSTRAINT resource_reviews_rating_range CHECK (rating BETWEEN 1 AND 5),
    CONSTRAINT resource_reviews_helpful_positive CHECK (helpful_count >= 0)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- Indexes
-- ─────────────────────────────────────────────────────────────────────────────

-- resource_providers
CREATE INDEX idx_resource_providers_name ON resource_providers(normalized_name);
CREATE INDEX idx_resource_providers_active ON resource_providers(is_active) WHERE is_active = TRUE;

-- learning_resources
CREATE INDEX idx_learning_resources_type ON learning_resources(resource_type);
CREATE INDEX idx_learning_resources_difficulty ON learning_resources(difficulty);
CREATE INDEX idx_learning_resources_cost ON learning_resources(cost_type);
CREATE INDEX idx_learning_resources_provider ON learning_resources(provider_id);
CREATE INDEX idx_learning_resources_active ON learning_resources(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_learning_resources_featured ON learning_resources(is_featured) WHERE is_featured = TRUE;
CREATE INDEX idx_learning_resources_rating ON learning_resources(rating DESC NULLS LAST);
-- Full-text search on title and description
CREATE INDEX idx_learning_resources_fts ON learning_resources
    USING GIN(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- resource_skills
CREATE INDEX idx_resource_skills_resource ON resource_skills(resource_id);
CREATE INDEX idx_resource_skills_name ON resource_skills(normalized_name);
CREATE INDEX idx_resource_skills_primary ON resource_skills(resource_id, is_primary) WHERE is_primary = TRUE;

-- resource_prerequisites
CREATE INDEX idx_resource_prereqs_resource ON resource_prerequisites(resource_id);
CREATE INDEX idx_resource_prereqs_skill ON resource_prerequisites(prerequisite_skill);

-- learning_paths
CREATE INDEX idx_learning_paths_role ON learning_paths(target_role);
CREATE INDEX idx_learning_paths_skill ON learning_paths(target_skill);
CREATE INDEX idx_learning_paths_active ON learning_paths(is_active) WHERE is_active = TRUE;

-- learning_path_resources
CREATE INDEX idx_lpr_path ON learning_path_resources(path_id, step_order);
CREATE INDEX idx_lpr_resource ON learning_path_resources(resource_id);

-- user_resource_progress
CREATE INDEX idx_urp_user ON user_resource_progress(user_id);
CREATE INDEX idx_urp_resource ON user_resource_progress(resource_id);
CREATE INDEX idx_urp_status ON user_resource_progress(user_id, status);
CREATE INDEX idx_urp_in_progress ON user_resource_progress(user_id, status)
    WHERE status = 'in_progress';

-- resource_reviews
CREATE INDEX idx_resource_reviews_resource ON resource_reviews(resource_id);
CREATE INDEX idx_resource_reviews_user ON resource_reviews(user_id);
CREATE INDEX idx_resource_reviews_rating ON resource_reviews(resource_id, rating DESC);

-- ─────────────────────────────────────────────────────────────────────────────
-- Triggers
-- ─────────────────────────────────────────────────────────────────────────────

CREATE TRIGGER resource_providers_updated_at
    BEFORE UPDATE ON resource_providers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER learning_resources_updated_at
    BEFORE UPDATE ON learning_resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER learning_paths_updated_at
    BEFORE UPDATE ON learning_paths
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER user_resource_progress_updated_at
    BEFORE UPDATE ON user_resource_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER resource_reviews_updated_at
    BEFORE UPDATE ON resource_reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ─────────────────────────────────────────────────────────────────────────────
-- Function: Update resource rating when a review is added/updated/deleted
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION refresh_resource_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE learning_resources
    SET
        rating = (
            SELECT ROUND(AVG(rating)::NUMERIC, 2)
            FROM resource_reviews
            WHERE resource_id = COALESCE(NEW.resource_id, OLD.resource_id)
        ),
        rating_count = (
            SELECT COUNT(*)
            FROM resource_reviews
            WHERE resource_id = COALESCE(NEW.resource_id, OLD.resource_id)
        )
    WHERE id = COALESCE(NEW.resource_id, OLD.resource_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER resource_reviews_refresh_rating
    AFTER INSERT OR UPDATE OR DELETE ON resource_reviews
    FOR EACH ROW EXECUTE FUNCTION refresh_resource_rating();

-- ─────────────────────────────────────────────────────────────────────────────
-- View: Resources with their primary skills (for quick lookup)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_resources_with_skills AS
SELECT
    lr.id,
    lr.title,
    lr.slug,
    lr.url,
    lr.resource_type,
    lr.difficulty,
    lr.cost_type,
    lr.cost_amount,
    lr.duration_hours,
    lr.rating,
    lr.rating_count,
    lr.has_certificate,
    lr.has_hands_on,
    lr.is_featured,
    rp.name AS provider_name,
    rp.website_url AS provider_url,
    ARRAY_AGG(rs.skill_name ORDER BY rs.is_primary DESC, rs.skill_name) AS skills,
    ARRAY_AGG(rs.normalized_name ORDER BY rs.is_primary DESC, rs.normalized_name) AS skill_ids
FROM learning_resources lr
LEFT JOIN resource_providers rp ON lr.provider_id = rp.id
LEFT JOIN resource_skills rs ON lr.id = rs.resource_id
WHERE lr.is_active = TRUE
GROUP BY lr.id, rp.name, rp.website_url;

-- ─────────────────────────────────────────────────────────────────────────────
-- View: Learning paths with resource count and total hours
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_learning_paths_summary AS
SELECT
    lp.id,
    lp.title,
    lp.slug,
    lp.description,
    lp.target_role,
    lp.target_skill,
    lp.difficulty,
    lp.estimated_hours,
    lp.is_featured,
    COUNT(lpr.id) AS resource_count,
    COUNT(lpr.id) FILTER (WHERE lpr.is_required = TRUE) AS required_resource_count
FROM learning_paths lp
LEFT JOIN learning_path_resources lpr ON lp.id = lpr.path_id
WHERE lp.is_active = TRUE
GROUP BY lp.id;

COMMIT;
