-- Migration 006: Create views for common read operations
-- Views optimize read-heavy queries by pre-joining related tables.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- v_user_full_profile: Complete profile with all sections
-- Used by the RAG pipeline and profile display
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_user_full_profile AS
SELECT
    u.id                    AS user_id,
    u.email,
    u.full_name,
    u.avatar_url,
    p.headline,
    p.summary,
    p.location_city,
    p.location_state,
    p.location_country,
    p.phone,
    p.linkedin_url,
    p.github_url,
    p.website_url,
    p.years_of_experience,
    p.is_open_to_work,
    p.profile_completeness,
    -- Skills as JSON array
    COALESCE(
        (SELECT jsonb_agg(
            jsonb_build_object(
                'id', s.id,
                'name', s.skill_name,
                'category', s.category,
                'proficiency', s.proficiency,
                'years_of_experience', s.years_of_experience,
                'is_primary', s.is_primary
            ) ORDER BY s.is_primary DESC, s.proficiency DESC, s.skill_name
        )
        FROM user_skills s WHERE s.user_id = u.id),
        '[]'::jsonb
    ) AS skills,
    -- Work experience as JSON array
    COALESCE(
        (SELECT jsonb_agg(
            jsonb_build_object(
                'id', e.id,
                'company', e.company_name,
                'title', e.job_title,
                'start_date', e.start_date,
                'end_date', e.end_date,
                'is_current', e.is_current,
                'duration_months', e.duration_months,
                'location', e.location,
                'responsibilities', e.responsibilities,
                'technologies', e.technologies_used
            ) ORDER BY e.start_date DESC
        )
        FROM work_experience e WHERE e.user_id = u.id),
        '[]'::jsonb
    ) AS work_experience,
    -- Education as JSON array
    COALESCE(
        (SELECT jsonb_agg(
            jsonb_build_object(
                'id', ed.id,
                'institution', ed.institution_name,
                'degree_level', ed.degree_level,
                'degree_name', ed.degree_name,
                'field', ed.field_of_study,
                'end_date', ed.end_date,
                'gpa', ed.gpa,
                'honors', ed.honors
            ) ORDER BY ed.end_date DESC NULLS FIRST
        )
        FROM education ed WHERE ed.user_id = u.id),
        '[]'::jsonb
    ) AS education,
    -- Certifications as JSON array
    COALESCE(
        (SELECT jsonb_agg(
            jsonb_build_object(
                'id', c.id,
                'name', c.name,
                'issuer', c.issuing_organization,
                'issue_date', c.issue_date,
                'expiry_date', c.expiry_date,
                'is_expired', c.is_expired,
                'credential_id', c.credential_id
            ) ORDER BY c.issue_date DESC NULLS LAST
        )
        FROM certifications c WHERE c.user_id = u.id),
        '[]'::jsonb
    ) AS certifications,
    p.updated_at            AS profile_updated_at,
    u.created_at            AS user_created_at
FROM users u
JOIN user_profiles p ON p.user_id = u.id
WHERE u.is_active = TRUE;

-- ─────────────────────────────────────────────────────────────────────────────
-- v_user_skills_summary: Aggregated skill statistics per user
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_user_skills_summary AS
SELECT
    user_id,
    COUNT(*)                                    AS total_skills,
    COUNT(*) FILTER (WHERE category = 'technical') AS technical_skills,
    COUNT(*) FILTER (WHERE category = 'soft')   AS soft_skills,
    COUNT(*) FILTER (WHERE proficiency = 'expert') AS expert_skills,
    COUNT(*) FILTER (WHERE proficiency = 'advanced') AS advanced_skills,
    COUNT(*) FILTER (WHERE proficiency = 'intermediate') AS intermediate_skills,
    COUNT(*) FILTER (WHERE proficiency = 'beginner') AS beginner_skills,
    ARRAY_AGG(skill_name ORDER BY is_primary DESC, proficiency DESC) FILTER (WHERE is_primary = TRUE)
        AS primary_skills,
    MAX(updated_at)                             AS last_updated
FROM user_skills
GROUP BY user_id;

-- ─────────────────────────────────────────────────────────────────────────────
-- v_work_experience_summary: Aggregated work experience per user
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_work_experience_summary AS
SELECT
    user_id,
    COUNT(*)                                    AS total_positions,
    SUM(duration_months)                        AS total_months,
    ROUND(SUM(duration_months)::NUMERIC / 12, 1) AS total_years,
    MAX(start_date)                             AS most_recent_start,
    BOOL_OR(is_current)                         AS has_current_position,
    ARRAY_AGG(DISTINCT company_name ORDER BY company_name) AS companies,
    ARRAY_AGG(DISTINCT UNNEST(technologies_used)) AS all_technologies
FROM work_experience
GROUP BY user_id;

-- ─────────────────────────────────────────────────────────────────────────────
-- v_active_skill_gaps: Unaddressed skill gaps with goal context
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_active_skill_gaps AS
SELECT
    sg.id,
    sg.user_id,
    sg.skill_name,
    sg.gap_type,
    sg.required_proficiency,
    sg.current_proficiency,
    sg.importance,
    cg.title                AS goal_title,
    cg.target_role          AS goal_target_role,
    cg.target_date          AS goal_target_date,
    sg.identified_at
FROM skill_gaps sg
LEFT JOIN career_goals cg ON cg.id = sg.career_goal_id
WHERE sg.is_addressed = FALSE
  AND (cg.id IS NULL OR cg.status = 'active')
ORDER BY
    CASE sg.importance
        WHEN 'critical' THEN 1
        WHEN 'important' THEN 2
        ELSE 3
    END,
    sg.identified_at DESC;

-- ─────────────────────────────────────────────────────────────────────────────
-- v_profile_completeness_breakdown: Detailed completeness per section
-- ─────────────────────────────────────────────────────────────────────────────
CREATE VIEW v_profile_completeness_breakdown AS
SELECT
    u.id                    AS user_id,
    u.full_name,
    p.profile_completeness,
    -- Section presence flags
    (p.headline IS NOT NULL AND LENGTH(p.headline) > 0)     AS has_headline,
    (p.summary IS NOT NULL AND LENGTH(p.summary) > 20)      AS has_summary,
    p.location_city IS NOT NULL                             AS has_location,
    p.linkedin_url IS NOT NULL                              AS has_linkedin,
    (SELECT COUNT(*) FROM user_skills s WHERE s.user_id = u.id) AS skill_count,
    (SELECT COUNT(*) FROM work_experience e WHERE e.user_id = u.id) AS experience_count,
    (SELECT COUNT(*) FROM education ed WHERE ed.user_id = u.id) AS education_count,
    (SELECT COUNT(*) FROM certifications c WHERE c.user_id = u.id) AS certification_count,
    (SELECT COUNT(*) FROM projects pr WHERE pr.user_id = u.id) AS project_count,
    -- Current resume
    (SELECT r.file_name FROM resume_uploads r
     WHERE r.user_id = u.id AND r.is_current = TRUE
     LIMIT 1)                                               AS current_resume_file
FROM users u
JOIN user_profiles p ON p.user_id = u.id
WHERE u.is_active = TRUE;

COMMIT;
