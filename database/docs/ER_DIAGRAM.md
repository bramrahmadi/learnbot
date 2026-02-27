# LearnBot User Profile Data Model — ER Diagram

## Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         LearnBot Database Schema                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────────┐
│      users       │         │  user_oauth_accounts  │
├──────────────────┤         ├──────────────────────┤
│ PK id (UUID)     │◄────────│ FK user_id (UUID)     │
│    email         │  1:N    │    provider           │
│    email_verified│         │    provider_id        │
│    password_hash │         │    access_token       │
│    full_name     │         │    refresh_token      │
│    avatar_url    │         │    token_expires_at   │
│    timezone      │         └──────────────────────┘
│    locale        │
│    is_active     │
│    is_admin      │
│    last_login_at │
│    created_at    │
│    updated_at    │
└────────┬─────────┘
         │ 1
         │
    ┌────┴──────────────────────────────────────────────────────────────┐
    │                                                                    │
    │ 1:1                                                                │ 1:N
    ▼                                                                    ▼
┌──────────────────────┐    ┌──────────────────────┐    ┌──────────────────────┐
│    user_profiles     │    │    resume_uploads     │    │    profile_history   │
├──────────────────────┤    ├──────────────────────┤    ├──────────────────────┤
│ PK id (UUID)         │    │ PK id (UUID)          │    │ PK id (BIGSERIAL)    │
│ FK user_id (UUID)    │    │ FK user_id (UUID)     │    │ FK user_id (UUID)    │
│    headline          │    │    file_name          │    │    event_type        │
│    summary           │    │    file_type          │    │    entity_type       │
│    location_city     │    │    file_size_bytes    │    │    entity_id         │
│    location_state    │    │    storage_key        │    │    old_data (JSONB)  │
│    location_country  │    │    version            │    │    new_data (JSONB)  │
│    phone             │    │    is_current         │    │    changed_by        │
│    linkedin_url      │    │    parsed_at          │    │    ip_address        │
│    github_url        │    │    parse_status       │    │    user_agent        │
│    website_url       │    │    parse_error        │    │    created_at        │
│    years_of_exp      │    │    raw_text           │    └──────────────────────┘
│    is_open_to_work   │    │    parser_version     │
│    profile_complete  │    │    overall_confidence │
│    created_at        │    │    created_at         │
│    updated_at        │    └──────────────────────┘
└──────────────────────┘
         │
         │ 1:N (via user_id)
    ┌────┴──────────────────────────────────────────────────────────────────┐
    │                │                │                │                    │
    ▼                ▼                ▼                ▼                    ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ user_skills  │ │work_experience│ │  education   │ │certifications│ │   projects   │
├──────────────┤ ├──────────────┤ ├──────────────┤ ├──────────────┤ ├──────────────┤
│PK id         │ │PK id         │ │PK id         │ │PK id         │ │PK id         │
│FK user_id    │ │FK user_id    │ │FK user_id    │ │FK user_id    │ │FK user_id    │
│FK taxonomy_id│ │FK upload_id  │ │FK upload_id  │ │FK upload_id  │ │FK upload_id  │
│   skill_name │ │  company_name│ │  institution │ │  name        │ │  name        │
│   normalized │ │  job_title   │ │  degree_level│ │  issuer      │ │  description │
│   category   │ │  emp_type    │ │  degree_name │ │  issue_date  │ │  project_url │
│   proficiency│ │  location    │ │  field       │ │  expiry_date │ │  repo_url    │
│   years_exp  │ │  start_date  │ │  start_date  │ │  credential_id│ │  technologies│
│   is_primary │ │  end_date    │ │  end_date    │ │  is_expired  │ │  start_date  │
│   source     │ │  is_current  │ │  gpa         │ │  confidence  │ │  end_date    │
│   confidence │ │  description │ │  honors      │ │  created_at  │ │  is_ongoing  │
│   created_at │ │  responsib.  │ │  confidence  │ │  updated_at  │ │  confidence  │
│   updated_at │ │  technologies│ │  created_at  │ └──────────────┘ │  created_at  │
└──────┬───────┘ │  duration_mo │ │  updated_at  │                  │  updated_at  │
       │         │  confidence  │ └──────────────┘                  └──────────────┘
       │         │  created_at  │
       │         │  updated_at  │
       │         └──────────────┘
       │
       │ N:1
       ▼
┌──────────────────────┐
│   skill_taxonomy     │
├──────────────────────┤
│ PK id (UUID)         │
│    name              │
│    normalized_name   │
│    category          │
│    aliases (TEXT[])  │
│ FK parent_skill_id   │ ◄─── self-referential (skill hierarchy)
│    is_verified       │
│    created_at        │
└──────────────────────┘

         │ 1:1 (via user_id)
    ┌────┴──────────────────────────────────────────────────────────────────┐
    │                                                                        │
    ▼                                                                        ▼
┌──────────────────────┐                                    ┌──────────────────────┐
│  user_preferences    │                                    │    career_goals      │
├──────────────────────┤                                    ├──────────────────────┤
│ PK id (UUID)         │                                    │ PK id (UUID)         │
│ FK user_id (UUID)    │                                    │ FK user_id (UUID)    │
│    desired_titles    │                                    │    title             │
│    desired_industries│                                    │    description       │
│    desired_co_sizes  │                                    │    target_role       │
│    desired_loc_types │                                    │    target_industry   │
│    desired_locations │                                    │    target_date       │
│    willing_relocate  │                                    │    status            │
│    salary_currency   │                                    │    priority          │
│    salary_min        │                                    │    progress_pct      │
│    salary_max        │                                    │    achieved_at       │
│    include_equity    │                                    │    created_at        │
│    job_search_urgency│                                    │    updated_at        │
│    available_from    │                                    └──────────┬───────────┘
│    career_stage      │                                               │ 1:N
│    email_*           │                                               ▼
│    created_at        │                                    ┌──────────────────────┐
│    updated_at        │                                    │     skill_gaps       │
└──────────────────────┘                                    ├──────────────────────┤
                                                            │ PK id (UUID)         │
                                                            │ FK user_id (UUID)    │
                                                            │ FK career_goal_id    │
                                                            │ FK taxonomy_id       │
                                                            │    skill_name        │
                                                            │    gap_type          │
                                                            │    required_prof     │
                                                            │    current_prof      │
                                                            │    importance        │
                                                            │    is_addressed      │
                                                            │    identified_at     │
                                                            │    addressed_at      │
                                                            └──────────────────────┘
```

## Cardinality Summary

| Relationship | Type | Description |
|---|---|---|
| `users` → `user_oauth_accounts` | 1:N | One user can have multiple OAuth providers |
| `users` → `user_profiles` | 1:1 | Each user has exactly one profile |
| `users` → `resume_uploads` | 1:N | One user can upload multiple resume versions |
| `users` → `profile_history` | 1:N | Full audit trail of all changes |
| `users` → `user_skills` | 1:N | One user has many skills |
| `users` → `work_experience` | 1:N | One user has many job entries |
| `users` → `education` | 1:N | One user has many education entries |
| `users` → `certifications` | 1:N | One user has many certifications |
| `users` → `projects` | 1:N | One user has many projects |
| `users` → `user_preferences` | 1:1 | Each user has one preferences record |
| `users` → `career_goals` | 1:N | One user can have multiple career goals |
| `career_goals` → `skill_gaps` | 1:N | One goal can have many skill gaps |
| `user_skills` → `skill_taxonomy` | N:1 | Many user skills map to one taxonomy entry |
| `skill_taxonomy` → `skill_taxonomy` | N:1 | Self-referential for skill hierarchy |
| `resume_uploads` → `work_experience` | 1:N | One upload can produce many experience entries |
| `resume_uploads` → `education` | 1:N | One upload can produce many education entries |
| `resume_uploads` → `certifications` | 1:N | One upload can produce many certifications |
| `resume_uploads` → `projects` | 1:N | One upload can produce many projects |
