# LearnBot Developer Wiki

Welcome to the LearnBot developer wiki. This is the in-depth technical reference for contributors and maintainers.

## Contents

| Article | Description |
|---------|-------------|
| [Architecture Deep Dive](architecture.md) | Detailed component internals |
| [Database Internals](database.md) | Schema design decisions, query patterns |
| [Resume Parser Internals](resume-parser.md) | Parsing pipeline, scoring algorithms |
| [Job Aggregator Internals](job-aggregator.md) | Scraping architecture, deduplication |
| [API Gateway Internals](api-gateway.md) | Middleware chain, handler patterns |
| [Frontend Architecture](frontend.md) | Next.js patterns, state management |
| [Testing Strategy](testing.md) | Test philosophy, patterns, coverage |
| [Security Model](security.md) | Auth, encryption, threat model |
| [Performance Guide](performance.md) | Benchmarks, optimization techniques |
| [Observability Guide](observability.md) | Metrics, logging, tracing |
| [Contributing Guide](contributing.md) | How to contribute effectively |
| [ADR Index](adr/README.md) | Architecture Decision Records |

---

## Quick Reference

### Service Ports

| Service | Port | Health Check |
|---------|------|-------------|
| API Gateway | 8090 | `GET /health` |
| Resume Parser | 8080 | `GET /health` |
| Job Aggregator | 8081 | `GET /admin/health` |
| Learning Resources | 8082 | `GET /health` |
| Frontend | 3000 | `GET /` |
| PostgreSQL | 5432 | `pg_isready` |
| Redis | 6379 | `PING` |
| Prometheus | 9090 | `GET /-/healthy` |
| Grafana | 3001 | `GET /api/health` |

### Key Files

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Local development stack |
| `.env.example` | Environment variable template |
| `api-gateway/docs/openapi.yaml` | OpenAPI specification |
| `database/migrations/` | SQL migration files |
| `infrastructure/terraform/` | AWS infrastructure |
| `.github/workflows/` | CI/CD pipelines |

### Go Module Structure

```
learnbot/
├── api-gateway/          github.com/learnbot/api-gateway
│   └── go.mod            replace: ../database, ../resume-parser
├── resume-parser/        github.com/learnbot/resume-parser
│   └── go.mod
├── job-aggregator/       github.com/learnbot/job-aggregator
│   └── go.mod
├── learning-resources/   github.com/learnbot/learning-resources
│   └── go.mod            replace: ../database
└── database/             github.com/learnbot/database
    └── go.mod
```

> **Important:** `api-gateway` and `learning-resources` use `replace` directives to reference local modules. Always build from the repo root when using Docker.

---

## Architecture Decision Records (ADRs)

ADRs document significant technical decisions. See [adr/README.md](adr/README.md) for the full index.

Key decisions:
- [ADR-001](adr/001-go-for-backend.md) — Go for all backend services
- [ADR-002](adr/002-postgresql-primary-db.md) — PostgreSQL as primary database
- [ADR-003](adr/003-jwt-authentication.md) — JWT for stateless authentication
- [ADR-004](adr/004-microservices-architecture.md) — Microservices over monolith
- [ADR-005](adr/005-ecs-fargate-deployment.md) — AWS ECS Fargate for deployment

---

## Development Standards

### Code Review Checklist

Before approving a PR, verify:

**Correctness**
- [ ] Logic is correct and handles edge cases
- [ ] Error handling is complete (no silent failures)
- [ ] No data races (run `go test -race ./...`)

**Testing**
- [ ] New code has tests
- [ ] Tests cover happy path and error cases
- [ ] Coverage doesn't decrease significantly

**Security**
- [ ] No secrets in code or logs
- [ ] Input validation on all user-provided data
- [ ] SQL queries use parameterized statements
- [ ] File uploads validate type and size

**Performance**
- [ ] No N+1 query patterns
- [ ] Large result sets are paginated
- [ ] Expensive operations are not in hot paths

**Documentation**
- [ ] Public functions have doc comments
- [ ] API changes update OpenAPI spec
- [ ] Breaking changes are noted in PR description

### Dependency Management

**Go:**
```bash
# Add a dependency
go get github.com/some/package@v1.2.3

# Update all dependencies
go get -u ./...

# Tidy (remove unused)
go mod tidy

# Verify checksums
go mod verify
```

**Node.js:**
```bash
# Add a dependency
npm install package-name

# Add dev dependency
npm install --save-dev package-name

# Update all
npm update

# Audit for vulnerabilities
npm audit
```

### Database Migration Guidelines

1. **Never modify existing migrations** — Create a new migration file
2. **Always test rollback** — Ensure migrations can be reversed
3. **Use transactions** — Wrap DDL changes in transactions where possible
4. **Add indexes** — Consider query patterns when adding new tables
5. **Document changes** — Update `database/docs/DATA_MODEL.md`

Migration file naming:
```
NNN_description.sql
001_create_enums.sql
002_create_users.sql
```

### API Design Guidelines

1. **Consistent envelope** — Always use `{ "success": bool, "data": ..., "meta": ... }`
2. **Meaningful error codes** — Use specific codes from the error code table
3. **Pagination** — All list endpoints must support `limit` and `offset`
4. **Validation** — Validate all inputs, return field-level errors
5. **Idempotency** — PUT operations should be idempotent
6. **Versioning** — Breaking changes require a new API version

---

## Glossary

| Term | Definition |
|------|-----------|
| **Acceptance Likelihood Score** | 0–100% score estimating how well a user's profile matches a job |
| **Gap Analysis** | Comparison of user skills vs. job requirements to identify missing skills |
| **RAG** | Retrieval-Augmented Generation — AI technique combining search with LLM generation |
| **Skill Gap** | A skill required by a target job that the user doesn't have or has at insufficient proficiency |
| **Skill Proficiency** | Level of expertise: beginner, intermediate, advanced, expert |
| **Profile Completeness** | 0–100% score measuring how much profile information has been provided |
| **Readiness Score** | 0–100% score measuring overall readiness for a specific role |
| **ECS Fargate** | AWS serverless container service |
| **JWT** | JSON Web Token — stateless authentication mechanism |
| **Soft Delete** | Marking records as inactive rather than physically deleting them |
| **Deduplication** | Process of identifying and removing duplicate job listings |
| **Token Bucket** | Rate limiting algorithm allowing burst traffic up to a configured limit |

---

## Getting Help

- **Slack:** `#learnbot-dev` for development questions
- **Slack:** `#learnbot-ops` for operational issues
- **GitHub Issues:** Bug reports and feature requests
- **GitHub Discussions:** Architecture discussions and RFCs
