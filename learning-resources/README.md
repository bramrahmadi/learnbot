# Learning Resources Service

A Go microservice that manages and serves a curated catalog of learning resources for the LearnBot platform. Provides search and filtering capabilities for courses, certifications, documentation, videos, books, and practice platforms.

Part of the [LearnBot AI Career Development Platform](../README.md).

## Features

- **Curated Catalog** — 60+ hand-picked learning resources across 8 resource types
- **Skill-Based Search** — Filter resources by skill name (e.g., "Go", "Python", "Docker")
- **Multi-Dimensional Filtering** — Filter by type, difficulty, cost, certificate availability, hands-on exercises
- **Rating Filter** — Minimum rating threshold
- **Pagination** — Efficient paginated results
- **Admin Dashboard** — Resource management endpoints

## Architecture

```
learning-resources/
├── cmd/
│   └── server/
│       └── main.go          # Entry point, server setup
├── internal/
│   ├── api/
│   │   ├── handler.go       # GET /api/resources/search
│   │   └── handler_test.go  # Integration tests
│   └── admin/
│       └── handler.go       # Admin resource management
├── Dockerfile
├── go.mod                   # replace: ../database
└── go.sum
```

## Quick Start

### Prerequisites

- Go >= 1.21
- PostgreSQL 15+ with learning resources seeded

### Run Locally

```bash
# Ensure database is running with seed data
docker compose up -d postgres

# Run the service
cd learning-resources
go run ./cmd/server
```

### Run with Docker Compose

```bash
# From repo root
docker compose up -d learning-resources
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | `8082` | HTTP server port |
| `ENVIRONMENT` | `development` | `development` or `production` |
| `DB_HOST` | `postgres` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `learnbot` | Database name |
| `DB_USER` | `learnbot_admin` | Database user |
| `DB_PASSWORD` | — | Database password |

## API Endpoints

### `GET /api/resources/search`

Search the learning resource catalog. **No authentication required.**

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `skill` | string | Filter by skill name (e.g., `Python`, `Go`, `Docker`) |
| `type` | string | `course`, `certification`, `documentation`, `video`, `book`, `practice`, `article`, `project` |
| `difficulty` | string | `beginner`, `intermediate`, `advanced`, `expert`, `all_levels` |
| `free` | boolean | Return only free resources |
| `has_certificate` | boolean | Return only resources with certificates |
| `has_hands_on` | boolean | Return only resources with hands-on exercises |
| `min_rating` | number | Minimum rating (0.0–5.0) |
| `limit` | integer | Results per page (default: 20, max: 100) |
| `offset` | integer | Pagination offset (default: 0) |

**Example Requests:**

```bash
# All Go resources
curl "http://localhost:8082/api/resources/search?skill=Go"

# Free Python courses for beginners
curl "http://localhost:8082/api/resources/search?skill=Python&type=course&difficulty=beginner&free=true"

# Resources with certificates, rated 4.5+
curl "http://localhost:8082/api/resources/search?has_certificate=true&min_rating=4.5"
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
        "url": "https://tour.golang.org",
        "description": "Interactive introduction to Go"
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

### `GET /health`

Health check endpoint.

```json
{ "status": "ok", "service": "learning-resources", "version": "1.0.0" }
```

## Resource Catalog

The catalog is seeded via `database/migrations/008_seed_learning_resources.sql` and includes resources for:

| Skill Category | Example Resources |
|----------------|------------------|
| Go | A Tour of Go, Go by Example, Effective Go |
| Python | Python.org Tutorial, Automate the Boring Stuff, Real Python |
| JavaScript/TypeScript | MDN Web Docs, TypeScript Handbook, JavaScript.info |
| React/Next.js | React Docs, Next.js Learn, Scrimba React Course |
| Docker/Kubernetes | Docker Getting Started, Kubernetes Docs, KodeKloud |
| AWS/Cloud | AWS Training, Cloud Practitioner Essentials |
| PostgreSQL | PostgreSQL Tutorial, Use The Index, Luke |
| Data Structures | LeetCode, HackerRank, Cracking the Coding Interview |
| System Design | System Design Primer, Designing Data-Intensive Applications |

## Running Tests

```bash
# All tests
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Integration tests only
go test ./internal/api/... -v
```

### Test Coverage Target

| Package | Target |
|---------|--------|
| `internal/api` | 70%+ |

## Adding New Resources

To add resources to the catalog, add `INSERT` statements to `database/migrations/008_seed_learning_resources.sql`:

```sql
INSERT INTO learning_resources (
    title, provider, resource_type, url, description,
    skills, difficulty, estimated_hours,
    is_free, price_usd, has_certificate, has_hands_on,
    rating, rating_count
) VALUES (
    'New Resource Title',
    'Provider Name',
    'course',  -- course, certification, documentation, video, book, practice, article, project
    'https://example.com/resource',
    'Brief description of the resource',
    ARRAY['Skill1', 'Skill2'],
    'intermediate',  -- beginner, intermediate, advanced, expert, all_levels
    10,  -- estimated hours
    false,  -- is_free
    29.99,  -- price_usd (NULL if free)
    true,  -- has_certificate
    true,  -- has_hands_on
    4.7,  -- rating (0.0-5.0)
    1500  -- rating_count
);
```

After adding, run the migration:
```bash
psql "postgres://learnbot_admin:password@localhost:5432/learnbot" \
  -f database/migrations/008_seed_learning_resources.sql
```
