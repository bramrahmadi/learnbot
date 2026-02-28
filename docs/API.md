# LearnBot — API Documentation

## Table of Contents

1. [Overview](#1-overview)
2. [Authentication Guide](#2-authentication-guide)
3. [Base URL & Versioning](#3-base-url--versioning)
4. [Request/Response Format](#4-requestresponse-format)
5. [Rate Limiting](#5-rate-limiting)
6. [Endpoints Reference](#6-endpoints-reference)
   - [Authentication](#61-authentication)
   - [Profile](#62-profile)
   - [Resume](#63-resume)
   - [Jobs](#64-jobs)
   - [Analysis](#65-analysis)
   - [Training](#66-training)
   - [Resources](#67-resources)
   - [Health](#68-health)
7. [Error Codes Reference](#7-error-codes-reference)
8. [Postman Collection](#8-postman-collection)

---

## 1. Overview

The LearnBot API is a RESTful HTTP API that provides access to all platform features: user authentication, profile management, resume parsing, job matching, skill gap analysis, and personalized training recommendations.

**OpenAPI Specification:** [`api-gateway/docs/openapi.yaml`](../api-gateway/docs/openapi.yaml)  
**Postman Collection:** [`api-gateway/docs/learnbot.postman_collection.json`](../api-gateway/docs/learnbot.postman_collection.json)

---

## 2. Authentication Guide

### Overview

LearnBot uses **JWT Bearer tokens** for authentication. Tokens are obtained via the `/api/auth/register` or `/api/auth/login` endpoints and must be included in the `Authorization` header of all protected requests.

### Step 1: Register or Login

**Register a new account:**
```bash
curl -X POST http://localhost:8090/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "securepassword123",
    "full_name": "Jane Doe"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "jane@example.com",
      "full_name": "Jane Doe"
    }
  }
}
```

### Step 2: Use the Token

Include the token in the `Authorization` header:

```bash
curl http://localhost:8090/api/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token Details

| Property | Value |
|----------|-------|
| Algorithm | HS256 |
| Expiry | 24 hours |
| Claims | `sub` (user ID), `email`, `exp`, `iat` |

### Token Expiry

When a token expires, you will receive a `401 Unauthorized` response:
```json
{
  "success": false,
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "token has expired"
  }
}
```

Re-authenticate via `/api/auth/login` to obtain a new token.

---

## 3. Base URL & Versioning

| Environment | Base URL |
|-------------|----------|
| Local Development | `http://localhost:8090` |
| Staging | `https://api-staging.learnbot.example.com` |
| Production | `https://api.learnbot.example.com` |

The API is currently at **v1**. The version is implicit in the URL path (no `/v1/` prefix in MVP).

---

## 4. Request/Response Format

### Content Type

All requests with a body must include:
```
Content-Type: application/json
```

Exception: Resume upload uses `multipart/form-data`.

### Success Response Envelope

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

- `data` — The response payload (object or array)
- `meta` — Pagination metadata (only present on list endpoints)

### Error Response Envelope

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "request validation failed",
    "details": [
      {
        "field": "email",
        "message": "must be a valid email address"
      }
    ]
  }
}
```

- `code` — Machine-readable error code (see [Error Codes Reference](#7-error-codes-reference))
- `message` — Human-readable error description
- `details` — Array of field-level validation errors (only on `VALIDATION_ERROR`)

---

## 5. Rate Limiting

The API enforces rate limiting per IP address:

| Property | Value |
|----------|-------|
| Rate | 10 requests/second |
| Burst | 30 requests |
| Algorithm | Token bucket |

When the rate limit is exceeded:
- **Status:** `429 Too Many Requests`
- **Header:** `Retry-After: <seconds>`

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "too many requests, please slow down"
  }
}
```

---

## 6. Endpoints Reference

### 6.1 Authentication

#### `POST /api/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "email": "jane@example.com",
  "password": "securepassword123",
  "full_name": "Jane Doe"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| `email` | string | ✓ | Valid email format |
| `password` | string | ✓ | Minimum 8 characters |
| `full_name` | string | ✓ | Non-empty |

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "jane@example.com",
      "full_name": "Jane Doe",
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` — Validation error (invalid email, short password)
- `409 Conflict` — Email already registered

---

#### `POST /api/auth/login`

Authenticate with email and password.

**Request Body:**
```json
{
  "email": "jane@example.com",
  "password": "securepassword123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "jane@example.com",
      "full_name": "Jane Doe"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` — Missing fields
- `401 Unauthorized` — Invalid credentials

---

### 6.2 Profile

> **Authentication required** for all profile endpoints.

#### `GET /api/users/profile`

Retrieve the current user's profile.

**Headers:** `Authorization: Bearer <token>`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "headline": "Senior Software Engineer",
    "summary": "Experienced backend engineer with 5+ years in Go and Python.",
    "location_city": "San Francisco",
    "location_country": "US",
    "years_of_experience": 5.5,
    "is_open_to_work": true,
    "profile_completeness": 85,
    "linkedin_url": "https://linkedin.com/in/janedoe",
    "github_url": "https://github.com/janedoe"
  }
}
```

---

#### `PUT /api/users/profile`

Update the current user's profile.

**Request Body:**
```json
{
  "headline": "Senior Software Engineer",
  "summary": "Experienced backend engineer with 5+ years in Go and Python.",
  "location_city": "San Francisco",
  "location_country": "US",
  "is_open_to_work": true,
  "linkedin_url": "https://linkedin.com/in/janedoe",
  "github_url": "https://github.com/janedoe"
}
```

All fields are optional. Only provided fields are updated.

**Success Response (200):** Updated profile object (same as GET response).

---

#### `GET /api/profile/skills`

Retrieve the current user's skills.

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "skills": [
      {
        "id": "skill-uuid",
        "skill_name": "Go",
        "category": "technical",
        "proficiency": "advanced",
        "years_of_experience": 4,
        "is_primary": true,
        "source": "resume_parser"
      },
      {
        "id": "skill-uuid-2",
        "skill_name": "Python",
        "category": "technical",
        "proficiency": "intermediate",
        "years_of_experience": 2,
        "is_primary": false,
        "source": "manual"
      }
    ]
  }
}
```

---

#### `PUT /api/profile/skills`

Replace the current user's skill list.

**Request Body:**
```json
{
  "skills": [
    {
      "name": "Go",
      "proficiency": "advanced",
      "years_of_experience": 4,
      "is_primary": true
    },
    {
      "name": "Python",
      "proficiency": "intermediate",
      "years_of_experience": 2
    }
  ]
}
```

| Proficiency Values | Description |
|-------------------|-------------|
| `beginner` | 0–1 year, basic awareness |
| `intermediate` | 1–3 years, works independently |
| `advanced` | 3–5 years, deep knowledge |
| `expert` | 5+ years, can teach others |

**Success Response (200):** Updated skills list.

**Error Responses:**
- `400 Bad Request` — Invalid proficiency value

---

### 6.3 Resume

> **Authentication required.**

#### `POST /api/resume/upload`

Upload and parse a resume file.

**Content-Type:** `multipart/form-data`

**Form Fields:**
| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| `resume` | file | ✓ | PDF or DOCX, max 10MB |

**Example (curl):**
```bash
curl -X POST http://localhost:8090/api/resume/upload \
  -H "Authorization: Bearer <token>" \
  -F "resume=@/path/to/resume.pdf"
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "upload_id": "upload-uuid",
    "file_name": "jane_doe_resume.pdf",
    "parse_status": "done",
    "overall_confidence": 0.92,
    "extracted": {
      "skills": [
        { "name": "Go", "proficiency": "advanced", "confidence": 0.95 },
        { "name": "PostgreSQL", "proficiency": "intermediate", "confidence": 0.88 }
      ],
      "work_experience": [
        {
          "company_name": "Acme Corp",
          "job_title": "Software Engineer",
          "start_date": "2020-01-01",
          "end_date": null,
          "is_current": true,
          "duration_months": 48
        }
      ],
      "education": [
        {
          "institution_name": "State University",
          "degree_level": "bachelor",
          "field_of_study": "Computer Science"
        }
      ],
      "certifications": []
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` — No file provided, invalid file type
- `422 Unprocessable Entity` — File could not be parsed

---

### 6.4 Jobs

#### `POST /api/jobs/search`

Search jobs with filters. **Authentication required.**

**Request Body:**
```json
{
  "query": "backend engineer",
  "skills": ["Go", "PostgreSQL"],
  "location_type": "remote",
  "experience_level": "senior",
  "limit": 20,
  "offset": 0
}
```

| Field | Type | Description |
|-------|------|-------------|
| `query` | string | Full-text search on job title |
| `skills` | string[] | Filter by required skills |
| `location_type` | string | `on_site`, `remote`, `hybrid` |
| `experience_level` | string | `entry`, `mid`, `senior`, `lead`, `executive` |
| `limit` | integer | Results per page (default: 20, max: 100) |
| `offset` | integer | Pagination offset (default: 0) |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "jobs": [
      {
        "id": "job-001",
        "title": "Senior Backend Engineer",
        "company": "TechCorp",
        "location": "Remote",
        "location_type": "remote",
        "experience_level": "senior",
        "required_skills": ["Go", "PostgreSQL", "Docker"],
        "preferred_skills": ["Kubernetes", "AWS"],
        "posted_at": "2024-01-10T00:00:00Z",
        "apply_url": "https://techcorp.com/jobs/123",
        "acceptance_likelihood": 78
      }
    ]
  },
  "meta": {
    "total": 145,
    "limit": 20,
    "offset": 0
  }
}
```

---

#### `GET /api/jobs/recommendations`

Get personalized job recommendations ranked by acceptance likelihood. **Authentication required.**

**Success Response (200):** Same format as job search, jobs sorted by `acceptance_likelihood` descending.

---

#### `GET /api/jobs/{id}`

Get details for a specific job. **No authentication required.**

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Job UUID |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "job-001",
    "title": "Senior Backend Engineer",
    "company": "TechCorp",
    "description": "We are looking for a senior backend engineer...",
    "location": "Remote",
    "location_type": "remote",
    "experience_level": "senior",
    "employment_type": "full_time",
    "required_skills": ["Go", "PostgreSQL", "Docker"],
    "preferred_skills": ["Kubernetes", "AWS"],
    "min_years_experience": 5,
    "salary_min": 120000,
    "salary_max": 160000,
    "salary_currency": "USD",
    "source": "linkedin",
    "apply_url": "https://techcorp.com/jobs/123",
    "posted_at": "2024-01-10T00:00:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found` — Job not found

---

#### `GET /api/jobs/{id}/match`

Calculate acceptance likelihood for a specific job. **Authentication required.**

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "job_id": "job-001",
    "acceptance_likelihood": 78,
    "breakdown": {
      "skill_match": 0.85,
      "experience_match": 0.80,
      "education_match": 0.90,
      "location_fit": 1.0,
      "industry_relevance": 0.70
    },
    "matched_skills": ["Go", "PostgreSQL"],
    "missing_required_skills": ["Docker"],
    "missing_preferred_skills": ["Kubernetes", "AWS"]
  }
}
```

---

### 6.5 Analysis

> **Authentication required.**

#### `POST /api/analysis/gaps`

Analyze skill gaps between the current user's profile and a target job.

**Request Body (by job ID):**
```json
{
  "job_id": "job-001"
}
```

**Request Body (inline job description):**
```json
{
  "job": {
    "title": "Senior Go Engineer",
    "required_skills": ["Go", "PostgreSQL", "Docker", "Kubernetes"],
    "preferred_skills": ["Terraform", "AWS"],
    "min_years_experience": 5,
    "experience_level": "senior"
  }
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "readiness_score": 72,
    "gaps": [
      {
        "skill": "Docker",
        "gap_type": "missing",
        "importance": "critical",
        "required_proficiency": "intermediate",
        "current_proficiency": null,
        "priority_score": 95
      },
      {
        "skill": "Kubernetes",
        "gap_type": "missing",
        "importance": "important",
        "required_proficiency": "beginner",
        "current_proficiency": null,
        "priority_score": 72
      }
    ],
    "strengths": ["Go", "PostgreSQL"],
    "visual_data": {
      "skill_radar": [...],
      "gap_distribution": {
        "critical": 1,
        "important": 1,
        "nice_to_have": 0
      }
    }
  }
}
```

---

### 6.6 Training

> **Authentication required.**

#### `GET /api/training/recommendations?job_id={id}`

Get a personalized training plan for a target job.

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `job_id` | string | Target job ID (optional) |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "plan": {
      "total_estimated_hours": 40,
      "estimated_weeks": 4,
      "phases": [
        {
          "phase": 1,
          "focus": "Critical Skills",
          "resources": [
            {
              "id": "resource-uuid",
              "title": "Docker for Developers",
              "provider": "Udemy",
              "type": "course",
              "difficulty": "intermediate",
              "estimated_hours": 8,
              "is_free": false,
              "price_usd": 14.99,
              "has_certificate": true,
              "url": "https://udemy.com/docker-for-developers",
              "skill": "Docker",
              "rating": 4.7
            }
          ]
        }
      ]
    }
  }
}
```

---

#### `POST /api/training/recommendations`

Get a personalized training plan with learning preferences.

**Request Body:**
```json
{
  "job_id": "job-001",
  "preferences": {
    "prefer_free": false,
    "max_budget_usd": 100,
    "weekly_hours_available": 10,
    "prefer_hands_on": true,
    "prefer_certificates": false
  }
}
```

**Success Response (200):** Same format as GET response.

---

### 6.7 Resources

> **No authentication required.**

#### `GET /api/resources/search`

Search the curated learning resource catalog.

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `skill` | string | Filter by skill name (e.g., `Python`, `Go`) |
| `type` | string | `course`, `certification`, `documentation`, `video`, `book`, `practice`, `article`, `project` |
| `difficulty` | string | `beginner`, `intermediate`, `advanced`, `expert`, `all_levels` |
| `free` | boolean | Return only free resources |
| `has_certificate` | boolean | Return only resources with certificates |
| `has_hands_on` | boolean | Return only resources with hands-on exercises |
| `min_rating` | number | Minimum rating (0.0–5.0) |
| `limit` | integer | Results per page (default: 20) |
| `offset` | integer | Pagination offset |

**Example:**
```bash
curl "http://localhost:8090/api/resources/search?skill=Go&free=true&difficulty=intermediate"
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "resources": [
      {
        "id": "resource-uuid",
        "title": "A Tour of Go",
        "provider": "golang.org",
        "type": "documentation",
        "difficulty": "beginner",
        "skills": ["Go"],
        "estimated_hours": 4,
        "is_free": true,
        "has_certificate": false,
        "has_hands_on": true,
        "rating": 4.8,
        "url": "https://tour.golang.org"
      }
    ]
  },
  "meta": {
    "total": 12,
    "limit": 20,
    "offset": 0
  }
}
```

---

### 6.8 Health

#### `GET /health`

Service health check. No authentication required.

**Success Response (200):**
```json
{
  "status": "ok",
  "service": "api-gateway",
  "version": "1.0.0"
}
```

---

## 7. Error Codes Reference

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `VALIDATION_ERROR` | Request body failed validation |
| 400 | `INVALID_REQUEST` | Malformed JSON or missing required fields |
| 401 | `UNAUTHORIZED` | Missing or invalid Authorization header |
| 401 | `TOKEN_EXPIRED` | JWT token has expired |
| 401 | `INVALID_TOKEN` | JWT token signature is invalid |
| 401 | `INVALID_CREDENTIALS` | Wrong email or password |
| 403 | `FORBIDDEN` | Authenticated but not authorized for this resource |
| 404 | `NOT_FOUND` | Requested resource does not exist |
| 409 | `CONFLICT` | Resource already exists (e.g., duplicate email) |
| 413 | `FILE_TOO_LARGE` | Uploaded file exceeds 10MB limit |
| 415 | `UNSUPPORTED_MEDIA_TYPE` | File type not supported (must be PDF or DOCX) |
| 422 | `PARSE_ERROR` | Resume file could not be parsed |
| 429 | `RATE_LIMIT_EXCEEDED` | Too many requests |
| 500 | `INTERNAL_ERROR` | Unexpected server error |
| 503 | `SERVICE_UNAVAILABLE` | Downstream service unavailable |

### Validation Error Details

When `code` is `VALIDATION_ERROR`, the `details` array contains field-level errors:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "request validation failed",
    "details": [
      {
        "field": "email",
        "message": "must be a valid email address"
      },
      {
        "field": "password",
        "message": "must be at least 8 characters"
      }
    ]
  }
}
```

---

## 8. Postman Collection

A complete Postman collection is available at [`api-gateway/docs/learnbot.postman_collection.json`](../api-gateway/docs/learnbot.postman_collection.json).

### Importing the Collection

1. Open Postman
2. Click **Import** → **File**
3. Select `api-gateway/docs/learnbot.postman_collection.json`
4. Set the `base_url` variable to `http://localhost:8090`

### Collection Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `base_url` | `http://localhost:8090` | API base URL |
| `token` | — | JWT token (auto-set by login/register requests) |
| `job_id` | `job-001` | Sample job ID for testing |

### Using the Collection

The collection includes pre-request scripts that automatically set the `token` variable after successful login/register requests. All authenticated endpoints use `{{token}}` in the Authorization header.

---

## Appendix: SDK Examples

### JavaScript/TypeScript

```typescript
const API_BASE = 'http://localhost:8090';

// Login
const loginResponse = await fetch(`${API_BASE}/api/auth/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email: 'jane@example.com', password: 'password123' })
});
const { data: { token } } = await loginResponse.json();

// Get profile
const profileResponse = await fetch(`${API_BASE}/api/users/profile`, {
  headers: { 'Authorization': `Bearer ${token}` }
});
const { data: profile } = await profileResponse.json();
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    // Login
    body, _ := json.Marshal(map[string]string{
        "email":    "jane@example.com",
        "password": "password123",
    })
    resp, _ := http.Post("http://localhost:8090/api/auth/login",
        "application/json", bytes.NewBuffer(body))
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    token := result["data"].(map[string]interface{})["token"].(string)
    
    // Get profile
    req, _ := http.NewRequest("GET", "http://localhost:8090/api/users/profile", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    client := &http.Client{}
    profileResp, _ := client.Do(req)
    fmt.Println(profileResp.Status)
}
```

### Python

```python
import requests

BASE_URL = "http://localhost:8090"

# Login
resp = requests.post(f"{BASE_URL}/api/auth/login", json={
    "email": "jane@example.com",
    "password": "password123"
})
token = resp.json()["data"]["token"]

# Get profile
headers = {"Authorization": f"Bearer {token}"}
profile = requests.get(f"{BASE_URL}/api/users/profile", headers=headers)
print(profile.json())
```
