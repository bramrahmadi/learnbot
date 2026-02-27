-- Migration 001: Create custom enum types
-- These enums enforce valid values at the database level.

BEGIN;

-- Skill proficiency levels
CREATE TYPE skill_proficiency AS ENUM (
    'beginner',       -- 0-1 year, basic awareness
    'intermediate',   -- 1-3 years, can work independently
    'advanced',       -- 3-5 years, deep knowledge
    'expert'          -- 5+ years, can teach others
);

-- Degree levels for education
CREATE TYPE degree_level AS ENUM (
    'high_school',
    'associate',
    'bachelor',
    'master',
    'doctorate',
    'professional',   -- MD, JD, etc.
    'certificate',
    'diploma',
    'other'
);

-- Skill categories
CREATE TYPE skill_category AS ENUM (
    'technical',
    'soft',
    'language',
    'tool',
    'framework',
    'database',
    'cloud',
    'other'
);

-- Employment type
CREATE TYPE employment_type AS ENUM (
    'full_time',
    'part_time',
    'contract',
    'freelance',
    'internship',
    'volunteer',
    'other'
);

-- Work location type
CREATE TYPE work_location_type AS ENUM (
    'on_site',
    'remote',
    'hybrid'
);

-- Career goal status
CREATE TYPE goal_status AS ENUM (
    'active',
    'achieved',
    'paused',
    'abandoned'
);

-- Profile change event types (for history tracking)
CREATE TYPE profile_event_type AS ENUM (
    'created',
    'resume_uploaded',
    'manually_updated',
    'skill_added',
    'skill_removed',
    'experience_added',
    'experience_updated',
    'education_added',
    'preference_updated'
);

COMMIT;
