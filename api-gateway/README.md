# API Gateway

The LearnBot API Gateway is the unified entry point for all client requests. It handles authentication, user profiles, job matching, skill gap analysis, and personalized training recommendations.

Part of the [LearnBot AI Career Development Platform](../README.md).

## Features

- **JWT Authentication** — Registration, login, and token validation
- **Profile Management** — User profile and skills CRUD
- **Resume Orchestration** — Forwards resume uploads to the Resume Parser service
- **Job Matching** — Job search, recommendations, and acceptance likelihood scoring
- **Skill Gap Analysis** — Compares user skills against job requirements
- **Training Recommendations** — Generates personalized learning plans
- **Resource Search** — Searches the curated learning resource catalog
- **Rate Limiting** — 10 req/s per IP with burst of 30 (token bucket)
- **Middleware Chain** — Recovery, logging, CORS, rate limiting, auth

## Architecture

```
api-gateway/
├── cmd/
│   └── server/
│       └── main.go          # Entry point, server setup, route registration
├── internal/
│   ├── handler/
│   │   ├── auth.go          # POST /api/auth/register, POST /api/auth/login
│   │   ├── profile.go       # GET/PUT /api/users/profile, GET/PUT /api/profile/skills
│   │   ├── resume.go        # POST /api/resume/upload
│   │   ├── jobs.go          # POST /api/jobs/search, GET /api/jobs/recommendations, etc.
│   │   ├── analysis.go      # POST /api/analysis/gaps
│   │   ├── response.go      # Shared response helpers (writeJSON, writeError)
│   │   ├── integration_test.go  # Integration tests
│   │   └── performance_test.go  # Load and benchmark tests
│   ├── middleware/
│   │   ├── auth.go          # JWT validation middleware
│   │   ├── logging.go       # Request/response logging
│   │   └── ratelimit.go     # Token bucket rate limiter
│   └── types/
│       └── types.go         # Shared request/response types
└── docs/
    ├── openapi.yaml         # OpenAPI 3.0 specification
    └── learnbot.postman_collection.json
```

## Quick Start

### Prerequisites

- Go >= 1.21
- PostgreSQL 15+ (or Docker)

### Run Locally

```bash
# Start dependencies
docker compose up -d postgres redis

# Run the API gateway
cd api-gateway
go run ./cmd/server \
  -addr :8090 \
  -jwt-secret "your-local-dev-secret-at-least-32-chars"
```

Or with environment variables:
```bash
JWT_SECRET=your-secret go run ./cmd/server
```

### Run with Docker Compose

```bash
# From repo root
docker compose up -d api-gateway
```

## Configuration

| Flag / Env Var | Default | Description |
|----------------|---------|-------------|
| `-addr` / `PORT` | `:8090` | HTTP server address |
| `-jwt-secret` / `JWT_SECRET` | — | JWT signing secret (required, min 32 chars) |
| `DB_HOST` | `postgres` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `learnbot` | Database name |
| `DB_USER` | `learnbot_admin` | Database user |
| `DB_PASSWORD` | — | Database password |
| `RESUME_PARSER_URL` | `http://resume-parser:8080` | Resume parser service URL |
| `JOB_AGGREGATOR_URL` | `http://job-aggregator:8081` | Job aggregator service URL |
| `LEARNING_RESOURCES_URL` | `http://learning-resources:8082` | Learning resources service URL |
| `REDIS_URL` | `redis://:password@redis:6379/0` | Redis connection URL |

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/auth/register` | — | Register new user |
| `POST` | `/api/auth/login` | — | Login |
| `GET` | `/api/users/profile` | ✓ | Get profile |
| `PUT` | `/api/users/profile` | ✓ | Update profile |
| `GET` | `/api/profile/skills` | ✓ | Get skills |
| `PUT` | `/api/profile/skills` | ✓ | Update skills |
| `POST` | `/api/resume/upload` | ✓ | Upload resume |
| `POST` | `/api/jobs/search` | ✓ | Search jobs |
| `GET` | `/api/jobs/recommendations` | ✓ | Get recommended jobs |
| `GET` | `/api/jobs/{id}` | — | Get job details |
| `GET` | `/api/jobs/{id}/match` | ✓ | Get acceptance likelihood |
| `POST` | `/api/analysis/gaps` | ✓ | Analyze skill gaps |
| `GET` | `/api/training/recommendations` | ✓ | Get training plan |
| `POST` | `/api/training/recommendations` | ✓ | Get training plan with preferences |
| `GET` | `/api/resources/search` | — | Search learning resources |
| `GET` | `/health` | — | Health check |

Full API documentation: [`docs/openapi.yaml`](docs/openapi.yaml) | [API Docs](../docs/API.md)

## Running Tests

```bash
# All tests
go test ./... -v

# Integration tests only
go test ./internal/handler/... -v -run TestIntegration

# Performance/load tests
go test ./internal/handler/... -v -run TestConcurrent

# Benchmarks
go test ./internal/handler/... -bench=. -benchmem -run=^$

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Coverage Targets

| Package | Target |
|---------|--------|
| `internal/handler` | 75%+ |
| `internal/middleware` | 80%+ |

## Middleware Chain

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
CORS (configurable origins, default: *)
  │
  ▼
Rate Limiter (10 req/s per IP, burst 30)
  │
  ▼
[Route-specific Auth Middleware]
  │
  ▼
Handler
```

## Response Format

All responses use a consistent JSON envelope:

```json
// Success
{
  "success": true,
  "data": { ... },
  "meta": { "total": 100, "limit": 20, "offset": 0 }
}

// Error
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "request validation failed",
    "details": [{ "field": "email", "message": "must be a valid email" }]
  }
}
```

## Health Check

```bash
curl http://localhost:8090/health
# {"status":"ok","service":"api-gateway","version":"1.0.0"}
```
