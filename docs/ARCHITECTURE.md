# LearnBot — Architecture Documentation

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Architecture Diagram](#2-architecture-diagram)
3. [Component Descriptions](#3-component-descriptions)
4. [Component Interaction Flows](#4-component-interaction-flows)
5. [Database Schema](#5-database-schema)
6. [API Architecture](#6-api-architecture)
7. [Security Architecture](#7-security-architecture)
8. [Observability Architecture](#8-observability-architecture)
9. [Technology Decisions](#9-technology-decisions)

---

## 1. System Overview

LearnBot is a microservices-based AI career development platform. The system is composed of five Go backend services, a Next.js frontend, a PostgreSQL database, and a Redis cache — all orchestrated via Docker Compose locally and AWS ECS Fargate in production.

### Design Principles

| Principle | Implementation |
|-----------|---------------|
| **Separation of Concerns** | Each service owns a single domain (auth, parsing, jobs, learning) |
| **API-First** | All inter-service communication via HTTP REST |
| **Stateless Services** | JWT auth; no server-side sessions |
| **Read-Optimized DB** | Materialized views, partial indexes, GIN full-text indexes |
| **Graceful Degradation** | Services fail independently; gateway returns partial data |
| **Security by Default** | JWT required on all non-public endpoints; rate limiting on all routes |

---

## 2. Architecture Diagram

### High-Level System Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              Internet                                     │
└──────────────────────────────────┬───────────────────────────────────────┘
                                   │ HTTPS
                    ┌──────────────▼──────────────┐
                    │     CloudFront CDN / ALB      │
                    │   (TLS termination, WAF)      │
                    └──────┬───────────────┬────────┘
                           │               │
              ┌────────────▼───┐   ┌───────▼────────────┐
              │   Frontend     │   │    API Gateway      │
              │  (Next.js)     │   │    (Go :8090)       │
              │  :3000         │   │                     │
              └────────────────┘   └──┬──────────────────┘
                                      │ Internal HTTP
              ┌───────────────────────┼──────────────────────┐
              │                       │                      │
   ┌──────────▼──────┐   ┌────────────▼────┐   ┌────────────▼────┐
   │  Resume Parser  │   │ Job Aggregator  │   │Learning Resources│
   │  (Go :8080)     │   │  (Go :8081)     │   │  (Go :8082)     │
   └──────────┬──────┘   └────────────┬────┘   └────────────┬────┘
              │                       │                      │
              └───────────────────────┼──────────────────────┘
                                      │
                    ┌─────────────────▼──────────────────┐
                    │         PostgreSQL :5432             │
                    │         (RDS Multi-AZ in prod)       │
                    └────────────────────────────────────┘
                                      │
                    ┌─────────────────▼──────────────────┐
                    │           Redis :6379               │
                    │      (Rate limiting, caching)        │
                    └────────────────────────────────────┘
```

### Local Development Stack

```
docker compose up -d
│
├── postgres:5432          PostgreSQL 15 (Alpine)
├── redis:6379             Redis 7 (Alpine)
├── api-gateway:8090       Go API Gateway
├── resume-parser:8080     Go Resume Parser
├── job-aggregator:8081    Go Job Aggregator
├── learning-resources:8082 Go Learning Resources
├── frontend:3000          Next.js 14
├── prometheus:9090        Metrics collection
├── grafana:3001           Metrics dashboards
└── nginx:80               Local reverse proxy
```

---

## 3. Component Descriptions

### 3.1 API Gateway (`api-gateway/`)

The central entry point for all client requests. Handles:

- **Authentication** — JWT registration, login, token validation
- **Profile Management** — User profile CRUD, skills management
- **Resume Orchestration** — Forwards resume uploads to the Resume Parser service
- **Job Matching** — Proxies job search/recommendations from Job Aggregator
- **Skill Gap Analysis** — Computes gap analysis using user profile vs. job requirements
- **Training Recommendations** — Generates personalized learning plans
- **Resource Search** — Searches the curated learning resource catalog

**Key design decisions:**
- Uses Go's standard `net/http` (no external framework) for minimal overhead
- Middleware chain: Recovery → Logger → CORS → Rate Limiter → Auth
- Rate limit: 10 req/s per IP, burst of 30 (token bucket via `golang.org/x/time/rate`)
- All responses use a consistent JSON envelope: `{ "success": bool, "data": ..., "meta": ... }`

### 3.2 Resume Parser (`resume-parser/`)

Parses PDF and DOCX resume files and extracts structured data:

- **Sections:** Skills, work experience, education, certifications, projects, personal info
- **Scoring:** Acceptance likelihood calculation (weighted: skills 35%, experience 25%, education 15%, location 10%, industry 15%)
- **Gap Analysis:** Identifies missing/weak skills vs. job requirements
- **Recommendations:** Suggests learning resources to close skill gaps

**Acceptance Likelihood Model:**
```
score = (skill_match × 0.35) + (experience_match × 0.25) +
        (education_match × 0.15) + (location_fit × 0.10) +
        (industry_relevance × 0.15)
```

### 3.3 Job Aggregator (`job-aggregator/`)

Collects job postings from multiple sources:

- **LinkedIn Jobs** — Public job search API (5 req/min, max 250 results)
- **Indeed** — Public job search (6 req/min, max 300 results)
- **Company Career Pages** — Configurable CSS selector or JSON API scraping

**Key features:**
- SHA-256 hash-based deduplication (by `source:external_id` or `source:title:company:location`)
- Daily scheduler at 2am UTC with configurable worker pool
- `robots.txt` compliance checking before scraping
- Exponential backoff retry logic
- Admin dashboard at `/admin/*` for monitoring

### 3.4 Learning Resources (`learning-resources/`)

Manages a curated catalog of 60+ learning resources:

- Courses, certifications, documentation, videos, books, practice platforms
- Filterable by skill, type, difficulty, cost, certificate availability
- Seeded via `database/migrations/008_seed_learning_resources.sql`

### 3.5 Frontend (`frontend/`)

Next.js 14 App Router application with:

- **Pages:** Landing, Login, Register, Onboarding, Dashboard, Jobs, Job Detail, Analysis, Learning
- **State Management:** Zustand (`authStore`) for authentication state
- **API Client:** Centralized `src/lib/api.ts` with JWT injection
- **UI Components:** Custom design system (Button, Card, Input, Badge, LoadingSpinner)
- **Testing:** Jest + React Testing Library (unit), Playwright (E2E)

### 3.6 Database (`database/`)

Shared PostgreSQL repository layer used by multiple services:

- Repository pattern with interfaces for testability
- UUID primary keys throughout
- Soft deletes (no hard deletes on user data)
- Audit trail via `profile_history` table
- Computed columns for `duration_months` and `is_expired`

---

## 4. Component Interaction Flows

### 4.1 User Registration Flow

```
Client                API Gateway           PostgreSQL
  │                       │                     │
  │── POST /api/auth/register ──────────────────▶│
  │                       │                     │
  │                       │── INSERT users ─────▶│
  │                       │◀── user_id ──────────│
  │                       │                     │
  │                       │── INSERT user_profiles ▶│
  │                       │◀── profile_id ───────│
  │                       │                     │
  │                       │── Sign JWT ──────────│
  │◀── 201 { token, user } ─────────────────────│
```

### 4.2 Resume Upload & Profile Building Flow

```
Client          API Gateway      Resume Parser      PostgreSQL
  │                 │                 │                 │
  │── POST /api/resume/upload ────────▶│                 │
  │                 │                 │                 │
  │                 │── Forward file ─▶│                 │
  │                 │                 │── Parse PDF/DOCX │
  │                 │                 │── Extract skills │
  │                 │                 │── Score confidence│
  │                 │◀── ParsedResume ─│                 │
  │                 │                 │                 │
  │                 │── UPDATE user_skills ─────────────▶│
  │                 │── INSERT resume_uploads ───────────▶│
  │                 │── INSERT work_experience ──────────▶│
  │                 │── INSERT education ────────────────▶│
  │                 │                 │                 │
  │◀── 200 { parsed_data } ──────────────────────────────│
```

### 4.3 Job Matching & Acceptance Likelihood Flow

```
Client          API Gateway      Job Aggregator     PostgreSQL
  │                 │                 │                 │
  │── GET /api/jobs/recommendations ──▶│                 │
  │                 │                 │                 │
  │                 │── GET user profile ───────────────▶│
  │                 │◀── UserProfile ───────────────────│
  │                 │                 │                 │
  │                 │── GET /jobs?skills=... ───────────▶│
  │                 │◀── JobList ──────────────────────│
  │                 │                 │                 │
  │                 │── Calculate acceptance likelihood  │
  │                 │   for each job (in-process)        │
  │                 │                 │                 │
  │◀── 200 { jobs with scores } ─────────────────────────│
```

### 4.4 Skill Gap Analysis Flow

```
Client          API Gateway      Resume Parser      PostgreSQL
  │                 │                 │                 │
  │── POST /api/analysis/gaps ────────▶│                 │
  │   { job_id: "job-001" }            │                 │
  │                 │                 │                 │
  │                 │── GET user skills ────────────────▶│
  │                 │◀── UserSkills ────────────────────│
  │                 │                 │                 │
  │                 │── GET job requirements ───────────▶│
  │                 │◀── JobRequirements ───────────────│
  │                 │                 │                 │
  │                 │── POST /analyze/gaps ─────────────▶│
  │                 │   { user_skills, job_requirements } │
  │                 │◀── GapAnalysis ──────────────────│
  │                 │                 │                 │
  │◀── 200 { gaps, readiness_score } ────────────────────│
```

### 4.5 Job Aggregation Flow (Background)

```
Scheduler (2am UTC)    Scrapers           PostgreSQL
       │                  │                   │
       │── Trigger run ───▶│                   │
       │                  │── LinkedIn scrape  │
       │                  │── Indeed scrape    │
       │                  │── Career pages     │
       │                  │                   │
       │                  │── Deduplicate ─────▶│
       │                  │   (SHA-256 hash)    │
       │                  │                   │
       │                  │── UPSERT jobs ─────▶│
       │                  │── Mark stale expired▶│
       │◀── Run complete ──│                   │
```

---

## 5. Database Schema

### Entity Relationship Overview

```
users (1) ──────────────── (1) user_profiles
  │                              │
  │ (1)                          │ (1)
  │                              │
  ├── (many) resume_uploads      ├── (many) user_skills
  ├── (many) work_experience     ├── (many) certifications
  ├── (many) education           ├── (many) skill_gaps
  ├── (1) user_preferences       └── (many) career_goals
  └── (many) profile_history

jobs (separate schema, managed by job-aggregator)
  └── company_career_pages

learning_resources (separate schema, managed by learning-resources)
```

### Core Tables

| Table | Purpose | Key Columns |
|-------|---------|-------------|
| `users` | Authentication & identity | `id`, `email`, `password_hash`, `is_active` |
| `user_profiles` | Professional profile | `headline`, `summary`, `location_*`, `years_of_experience` |
| `resume_uploads` | Versioned resume files | `storage_key`, `parse_status`, `is_current` |
| `user_skills` | Skills with proficiency | `skill_name`, `proficiency`, `years_of_experience`, `is_primary` |
| `work_experience` | Job history | `company_name`, `job_title`, `start_date`, `duration_months` |
| `education` | Educational background | `institution_name`, `degree_level`, `field_of_study` |
| `certifications` | Professional certs | `name`, `issuing_organization`, `expiry_date`, `is_expired` |
| `user_preferences` | Job search preferences | `desired_job_titles[]`, `salary_min/max`, `job_search_urgency` |
| `career_goals` | Career objectives | `status`, `priority`, `progress_percentage` |
| `skill_gaps` | Identified skill gaps | `gap_type`, `importance`, `is_addressed` |
| `profile_history` | Immutable audit log | `event_type`, `old_data`, `new_data` (JSONB) |

### Database Views

| View | Purpose |
|------|---------|
| `v_user_full_profile` | Complete profile as JSONB (used by RAG pipeline) |
| `v_user_skills_summary` | Aggregated skill stats per user |
| `v_work_experience_summary` | Total experience, company list, technologies |
| `v_active_skill_gaps` | Unaddressed gaps ordered by importance |
| `v_profile_completeness_breakdown` | Per-section completeness flags |

### Enum Types

| Enum | Values |
|------|--------|
| `skill_proficiency` | `beginner`, `intermediate`, `advanced`, `expert` |
| `skill_category` | `technical`, `soft`, `language`, `tool`, `framework`, `database`, `cloud`, `other` |
| `employment_type` | `full_time`, `part_time`, `contract`, `freelance`, `internship`, `volunteer`, `other` |
| `work_location_type` | `on_site`, `remote`, `hybrid` |
| `degree_level` | `high_school`, `associate`, `bachelor`, `master`, `doctorate`, `professional`, `certificate`, `diploma`, `other` |
| `goal_status` | `active`, `achieved`, `paused`, `abandoned` |

For full schema documentation, see [database/docs/DATA_MODEL.md](../database/docs/DATA_MODEL.md) and [database/docs/ER_DIAGRAM.md](../database/docs/ER_DIAGRAM.md).

---

## 6. API Architecture

### Request/Response Envelope

All API responses use a consistent JSON envelope:

**Success:**
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "total": 100,
    "limit": 20,
    "offset": 0
  }
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "request validation failed",
    "details": [
      { "field": "email", "message": "must be a valid email" }
    ]
  }
}
```

### Authentication Architecture

```
Client                    API Gateway
  │                           │
  │── POST /api/auth/login ───▶│
  │◀── { token: "eyJ..." } ───│
  │                           │
  │── GET /api/users/profile ─▶│
  │   Authorization: Bearer eyJ...
  │                           │── Validate JWT signature
  │                           │── Check expiry (24h)
  │                           │── Extract user_id from claims
  │◀── 200 { profile } ───────│
```

JWT Claims structure:
```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Rate Limiting Architecture

- **Algorithm:** Token bucket (via `golang.org/x/time/rate`)
- **Limit:** 10 requests/second per IP
- **Burst:** 30 requests
- **Response on limit:** `429 Too Many Requests`
- **Headers:** `Retry-After` included in 429 responses

### Middleware Chain

```
Request
  │
  ▼
Recovery (panic → 500)
  │
  ▼
Logger (structured request/response logging)
  │
  ▼
CORS (configurable origins)
  │
  ▼
Rate Limiter (per-IP token bucket)
  │
  ▼
[Route-specific Auth Middleware]
  │
  ▼
Handler
```

---

## 7. Security Architecture

### Authentication & Authorization

- **Mechanism:** JWT (HS256) with 24-hour expiry
- **Secret:** Minimum 32-character random string, stored in AWS Secrets Manager
- **Validation:** Signature, expiry, and issuer checked on every protected request
- **No refresh tokens in MVP** — users re-authenticate after 24 hours

### Data Protection

| Layer | Mechanism |
|-------|-----------|
| Transport | TLS 1.3 (enforced by ALB in production) |
| Data at rest | AES-256 via AWS KMS (RDS + S3) |
| Passwords | bcrypt (cost factor 12) |
| Resume files | S3 with SSE-KMS, pre-signed URLs |
| Secrets | AWS Secrets Manager (never in env files) |

### Network Security

- Services run in private VPC subnets (no public IPs)
- ALB is the only internet-facing component
- Security groups follow least-privilege (ECS → RDS only, no direct DB access)
- VPC Flow Logs enabled

### GDPR Compliance

- Soft deletes preserve referential integrity while allowing data removal
- `profile_history` audit log tracks all changes
- User data export supported via `v_user_full_profile` view
- No third-party analytics without consent

---

## 8. Observability Architecture

### Metrics (Prometheus + Grafana)

```
Services ──── /metrics ────▶ Prometheus ────▶ Grafana Dashboards
                              (scrape every 15s)
```

**Key metrics collected:**
- HTTP request rate, latency (P50/P95/P99), error rate
- Database connection pool utilization
- Resume parse success/failure rate
- Job scraping run duration and job count
- Go runtime metrics (GC, goroutines, memory)

### Alerting Rules

| Alert | Condition | Severity |
|-------|-----------|----------|
| High error rate | > 5% errors for 5 min | Warning |
| Critical error rate | > 20% errors for 2 min | Critical |
| High latency | P99 > 2s for 5 min | Warning |
| Service down | No scrape for 1 min | Critical |
| High CPU | > 85% for 10 min | Warning |
| Low disk | < 15% free | Warning |

### Logging

- Structured JSON logs to stdout (collected by CloudWatch in production)
- Log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`
- Request logs include: method, path, status, duration, user_id (if authenticated)
- Correlation IDs for request tracing (planned for Phase 2)

---

## 9. Technology Decisions

### Why Go for Backend Services?

- **Performance:** Compiled, low memory footprint, excellent concurrency via goroutines
- **Simplicity:** Standard library covers most needs (HTTP, JSON, testing)
- **Deployment:** Single static binary, minimal Docker image size
- **Concurrency:** Goroutines ideal for concurrent job scraping and parallel request handling

### Why Next.js for Frontend?

- **SSR/SSG:** Server-side rendering for SEO and initial load performance
- **App Router:** File-based routing with React Server Components
- **TypeScript:** Type safety across the entire frontend codebase
- **Ecosystem:** Rich component library ecosystem, excellent testing support

### Why PostgreSQL?

- **JSONB:** Flexible storage for audit logs and parsed resume data
- **Arrays:** Native array columns for skills, technologies, activities
- **Generated Columns:** Automatic `duration_months` and `is_expired` computation
- **Full-Text Search:** GIN indexes for skill name search
- **Reliability:** ACID compliance, mature ecosystem, excellent Go drivers

### Why Redis?

- **Rate Limiting:** Distributed rate limiting across multiple API gateway instances
- **Caching:** Job search results, user profile caching (planned)
- **Session Store:** Future use for refresh tokens

### Why AWS ECS Fargate?

- **Serverless containers:** No EC2 instance management
- **Auto-scaling:** Scale individual services independently
- **Cost:** Pay only for running containers
- **Integration:** Native integration with ECR, ALB, CloudWatch, Secrets Manager
