# Training Recommendation Engine

## Overview

The training recommendation engine generates personalized learning plans for candidates based on their skill gap analysis. It uses a rule-based MVP approach to match skill gaps with curated learning resources and produce a structured, phased learning timeline.

## Architecture

```
POST /api/v1/recommendations
         │
         ▼
  RecommendationRequest
  ┌─────────────────────┐
  │ CandidateProfile    │
  │ JobRequirements     │
  │ UserPreferences     │
  └─────────────────────┘
         │
         ▼
  Engine.Generate()
         │
         ├─► GapAnalyzer.Analyze()     → GapAnalysisResult
         │         (from gapanalysis package)
         │
         ├─► buildSkillRecommendations()
         │         │
         │         ├─► findMatchingResources()  → filtered catalog entries
         │         ├─► scoreResources()         → ranked RecommendedResources
         │         └─► buildSkillRecommendation() → SkillRecommendation
         │
         ├─► buildPhases()             → []LearningPhase
         │
         ├─► buildTimeline()           → LearningTimeline
         │
         └─► buildSummary()            → LearningPlanSummary
                   │
                   ▼
            LearningPlan (JSON response)
```

## API Endpoint

### `POST /api/v1/recommendations`

Generates a personalized learning plan.

**Request Body:**

```json
{
  "profile": {
    "skills": [
      {"name": "Python", "proficiency": "advanced"},
      {"name": "SQL", "proficiency": "intermediate"}
    ],
    "years_of_experience": 3,
    "work_history": [...],
    "education": [...]
  },
  "job": {
    "title": "Senior Backend Engineer",
    "required_skills": ["Go", "PostgreSQL", "Docker", "Kubernetes"],
    "preferred_skills": ["Terraform", "AWS"],
    "min_years_experience": 5,
    "experience_level": "senior"
  },
  "preferences": {
    "prefer_free": false,
    "max_budget_usd": 100,
    "weekly_hours_available": 10,
    "prefer_hands_on": true,
    "prefer_certificates": false,
    "target_date": "2025-09-01",
    "preferred_resource_types": ["course", "documentation"],
    "excluded_providers": []
  }
}
```

**Response Body:**

```json
{
  "success": true,
  "data": {
    "job_title": "Senior Backend Engineer",
    "readiness_score": 42.5,
    "total_gaps": 4,
    "total_estimated_hours": 145.0,
    "phases": [
      {
        "phase_number": 1,
        "phase_name": "Critical Skills",
        "phase_description": "Master the must-have skills required for this role.",
        "skills": [
          {
            "skill_name": "Go",
            "gap_category": "critical",
            "priority_score": 0.8750,
            "primary_resource": {
              "resource": {
                "id": "go-complete-guide",
                "title": "Go: The Complete Developer's Guide",
                "url": "https://www.udemy.com/course/go-the-complete-developers-guide/",
                "provider": "Udemy",
                "resource_type": "course",
                "difficulty": "intermediate",
                "cost_type": "paid",
                "cost_usd": 19.99,
                "duration_hours": 9,
                "rating": 4.6,
                "has_certificate": true,
                "has_hands_on": true
              },
              "relevance_score": 0.8234,
              "is_alternative": false,
              "recommendation_reason": "Recommended because it directly covers Go, highly rated (4.6/5), curated resource.",
              "estimated_completion_hours": 9.0
            },
            "alternative_resources": [...],
            "estimated_hours_to_job_ready": 30,
            "target_level": "intermediate"
          }
        ],
        "total_hours": 85.0,
        "estimated_weeks": 8.5,
        "milestone": "✅ Job-ready in Go, PostgreSQL, Docker – cleared all critical requirements"
      },
      {
        "phase_number": 2,
        "phase_name": "Preferred Skills",
        ...
      }
    ],
    "timeline": {
      "total_weeks": 12,
      "total_hours": 145.0,
      "weekly_hours": 10,
      "weeks": [
        {
          "week_number": 1,
          "phase_number": 1,
          "skill_focus": "Go",
          "resource_title": "Go: The Complete Developer's Guide",
          "hours_planned": 10.0,
          "cumulative_hours": 10.0,
          "activities": [
            "Start 'Go: The Complete Developer's Guide'",
            "Set up development environment for Go",
            "Complete introductory exercises",
            "Supplement with LeetCode/HackerRank problems for Go"
          ],
          "is_checkpoint": false
        }
      ],
      "target_completion_date": "2025-09-01"
    },
    "matched_skills": ["Python", "SQL"],
    "summary": {
      "headline": "📚 4 critical gaps to close for Senior Backend Engineer – estimated 15 weeks",
      "critical_gap_count": 4,
      "important_gap_count": 2,
      "free_resource_count": 2,
      "paid_resource_count": 4,
      "estimated_total_cost_usd": 79.96,
      "top_skills_to_learn": ["Go", "PostgreSQL", "Docker"],
      "quick_wins": ["Docker"]
    }
  }
}
```

## Recommendation Algorithm

### 1. Gap Analysis

The engine first runs the gap analysis module (`gapanalysis.Analyzer`) to identify:
- **Critical gaps**: Must-have skills the candidate is missing
- **Important gaps**: Preferred skills the candidate is missing
- **Nice-to-have gaps**: Optional skills the candidate is missing

### 2. Resource Matching

For each skill gap, the engine queries the built-in resource catalog:

```
findMatchingResources(skillName, preferences)
  1. Normalize skill name (lowercase, trim)
  2. Resolve aliases (e.g. "golang" → "go")
  3. For each catalog entry:
     a. Check if primary skill matches
     b. Check if any skill in skills[] matches
     c. Check substring containment for compound skills
  4. Apply preference filters:
     - prefer_free: exclude paid resources
     - max_budget_usd: exclude resources above budget
     - excluded_providers: exclude specific providers
     - preferred_resource_types: only include specified types
```

### 3. Relevance Scoring

Each matching resource is scored using a composite formula:

```
relevance_score = 
  skill_match_quality  × 0.30  +  // primary vs secondary skill match
  difficulty_fit       × 0.20  +  // matches user's current → target level
  quality_score        × 0.20  +  // rating + verification bonus
  preference_alignment × 0.20  +  // free/hands-on/certificate preferences
  popularity_score     × 0.10     // log-normalized rating count
```

**Skill Match Quality:**
- Primary skill match: 1.0
- Secondary skill match: 0.5

**Difficulty Fit:**
- Resource difficulty within [current_level, target_level]: 1.0
- One level off: 0.7
- Two+ levels off: 0.4
- `all_levels` resource: 0.9

**Quality Score:**
- `rating / 5.0` (base)
- +0.1 bonus for verified/curated resources (capped at 1.0)

**Preference Alignment:**
- Base: 0.5
- +0.3 if user prefers free and resource is free/free_audit
- +0.1 if user prefers hands-on and resource has hands-on
- +0.1 if user prefers certificates and resource has certificate

**Popularity Score:**
- `min(1.0, log10(rating_count) / 7.0)` (10M ratings = 1.0)

### 4. Learning Hours Estimation

Estimated completion hours are adjusted based on the user's current level:

```
hours = resource.duration_hours
if current_level == "beginner":    hours *= 0.70
if current_level == "intermediate": hours *= 0.40
if current_level == "advanced":    hours *= 0.15
```

### 5. Phase Building

Gaps are organized into phases:
- **Phase 1: Critical Skills** — must-have requirements
- **Phase 2: Preferred Skills** — preferred requirements
- **Phase 3: Nice-to-Have Skills** — optional skills

Each phase includes:
- Total estimated hours
- Estimated weeks (total_hours / weekly_hours_available)
- A milestone achievement description

### 6. Timeline Generation

The timeline converts phases into a week-by-week schedule:

1. For each skill in each phase:
   - Calculate weeks needed: `ceil(resource_hours / weekly_hours)`
   - Assign activities for each week (start/continue/complete)
   - Mark the last week of critical skills as a checkpoint
2. Add a review week after phases with multiple skills
3. Calculate cumulative hours
4. Estimate target completion date

### 7. Summary Generation

The summary provides a high-level overview:
- **Headline**: One-line description of the learning plan
- **Quick wins**: Skills that can be learned in < 20 hours
- **Cost breakdown**: Free vs paid resource counts and total cost
- **Top skills**: First 3 skills to focus on

## Resource Catalog

The built-in catalog contains 60+ curated resources covering:

| Category | Resources |
|----------|-----------|
| Languages | Python, Go, JavaScript, TypeScript, Java, Rust, C#, Kotlin, Swift |
| Frameworks | React, Vue, Angular, Node.js, Spring Boot, Django |
| Databases | PostgreSQL, MongoDB, Redis, Elasticsearch, Kafka |
| Cloud | AWS (SAA, Developer), GCP (ACE, Architect), Azure (AZ-900, AZ-204) |
| DevOps | Docker, Kubernetes (CKA), Terraform, Ansible, GitHub Actions |
| ML/AI | Machine Learning, Deep Learning, PyTorch, TensorFlow, NLP, LangChain |
| Algorithms | LeetCode, HackerRank, Stanford Algorithms |
| Other | Git, Linux, System Design, Agile/Scrum, Security |

Each resource includes:
- Provider, type, difficulty, cost model
- Duration in hours
- Rating and rating count
- Skills covered (primary + secondary)
- Certificate and hands-on flags

## User Preferences

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `prefer_free` | bool | false | Only show free/free_audit resources |
| `max_budget_usd` | float | 0 (no limit) | Maximum cost per resource |
| `weekly_hours_available` | float | 10 | Hours per week for studying |
| `prefer_hands_on` | bool | false | Boost hands-on resources |
| `prefer_certificates` | bool | false | Boost resources with certificates |
| `target_date` | string | "" | Target completion date (ISO 8601) |
| `preferred_resource_types` | []string | [] (all) | Filter by resource type |
| `excluded_providers` | []string | [] | Exclude specific providers |

## Integration with Gap Analysis

The recommendation engine is built on top of the gap analysis module:

```go
// Gap analysis provides the input gaps.
gapResult := gapAnalyzer.Analyze(profile, job)

// Recommendation engine processes each gap.
for _, gap := range gapResult.CriticalGaps {
    rec := engine.buildSkillRecommendation(gap, prefs)
    // ...
}
```

The gap analysis provides:
- `SkillGap.SkillName` — used for catalog lookup
- `SkillGap.Category` — determines phase placement
- `SkillGap.PriorityScore` — used for ordering within phases
- `SkillGap.EstimatedLearningHours` — fallback when no resource duration
- `SkillGap.CurrentLevel` — used for hours adjustment
- `SkillGap.TargetLevel` — used for difficulty fit scoring

## Example Usage

```bash
curl -X POST http://localhost:8080/api/v1/recommendations \
  -H "Content-Type: application/json" \
  -d '{
    "profile": {
      "skills": [
        {"name": "Python", "proficiency": "advanced"},
        {"name": "SQL", "proficiency": "intermediate"}
      ],
      "years_of_experience": 3
    },
    "job": {
      "title": "ML Engineer",
      "required_skills": ["Python", "TensorFlow", "PyTorch", "Kubernetes"],
      "preferred_skills": ["Docker", "AWS"],
      "min_years_experience": 3,
      "experience_level": "mid"
    },
    "preferences": {
      "weekly_hours_available": 15,
      "prefer_hands_on": true,
      "prefer_free": false
    }
  }'
```
