# LearnBot — AI-Powered Career Development Platform

[![Tests](https://github.com/learnbot/learnbot/actions/workflows/test.yml/badge.svg)](https://github.com/learnbot/learnbot/actions/workflows/test.yml)
[![Security](https://github.com/learnbot/learnbot/actions/workflows/security.yml/badge.svg)](https://github.com/learnbot/learnbot/actions/workflows/security.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

LearnBot is an AI-powered mentorship platform that analyzes your professional profile and provides personalized guidance to help you achieve your career goals through intelligent job matching, skill gap analysis, and tailored learning paths.

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [Architecture Guide](docs/ARCHITECTURE.md) | System design, component diagrams, data flows |
| [API Documentation](docs/API.md) | Endpoint reference, auth guide, error codes |
| [Developer Guide](docs/DEVELOPER_GUIDE.md) | Setup, workflow, code style, testing |
| [Deployment Guide](docs/DEPLOYMENT_GUIDE.md) | Infrastructure, CI/CD, configuration |
| [User Guide](docs/USER_GUIDE.md) | Feature walkthroughs for end users |
| [FAQ](docs/FAQ.md) | Frequently asked questions |
| [Testing Guide](docs/TESTING.md) | Test strategy, coverage, running tests |
| [Developer Wiki](docs/wiki/README.md) | In-depth technical reference |
| [User Help Center](docs/help/README.md) | User support documentation |

---

## 🎯 What LearnBot Does

LearnBot bridges the gap between where you are in your career and where you want to be:

1. **Resume Analysis** — Upload your resume (PDF/DOCX) and LearnBot extracts your skills, experience, education, and certifications into a structured profile.
2. **Intelligent Job Matching** — Browse jobs ranked by an *Acceptance Likelihood Score* (0–100%) that reflects how well your profile matches each role.
3. **Skill Gap Analysis** — See exactly which skills you're missing for your target role, categorized as critical, important, or nice-to-have.
4. **Personalized Learning Paths** — Get curated course and certification recommendations prioritized by impact and time investment.
5. **Career Goal Tracking** — Set career goals and track your progress over time.

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        LearnBot Platform                         │
│                                                                  │
│  ┌──────────────┐    ┌──────────────────────────────────────┐   │
│  │   Frontend   │    │           Backend Services            │   │
│  │  (Next.js)   │───▶│                                      │   │
│  │  :3000       │    │  ┌─────────────┐  ┌───────────────┐  │   │
│  └──────────────┘    │  │ API Gateway │  │ Resume Parser │  │   │
│                      │  │   :8090     │  │    :8080      │  │   │
│                      │  └──────┬──────┘  └───────────────┘  │   │
│                      │         │                             │   │
│                      │  ┌──────▼──────┐  ┌───────────────┐  │   │
│                      │  │Job Aggregat.│  │  Learning Res. │  │   │
│                      │  │   :8081     │  │    :8082      │  │   │
│                      │  └─────────────┘  └───────────────┘  │   │
│                      └──────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────┐    ┌──────────────────────────────────────┐   │
│  │  PostgreSQL  │    │           Monitoring                  │   │
│  │   :5432      │    │  Prometheus :9090 | Grafana :3001    │   │
│  └──────────────┘    └──────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

See the full [Architecture Documentation](docs/ARCHITECTURE.md) for detailed component diagrams and data flows.

---

## 🚀 Quick Start

### Prerequisites

- [Docker](https://www.docker.com/) >= 24.0
- [Docker Compose](https://docs.docker.com/compose/) >= 2.0
- [Go](https://go.dev/) >= 1.21 (for local development)
- [Node.js](https://nodejs.org/) >= 20 (for frontend development)

### 1. Clone the Repository

```bash
git clone https://github.com/learnbot/learnbot.git
cd learnbot
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your local values (JWT secret, etc.)
```

### 3. Start All Services

```bash
docker compose up -d
```

### 4. Run Database Migrations

```bash
docker compose exec postgres psql -U learnbot_admin -d learnbot \
  -f /docker-entrypoint-initdb.d/001_create_enums.sql
```

> **Note:** Migrations run automatically on first startup via Docker's `initdb` mechanism.

### 5. Access the Application

| Service | URL | Credentials |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | Register a new account |
| API Gateway | http://localhost:8090 | JWT Bearer token |
| Grafana | http://localhost:3001 | admin / admin |
| Prometheus | http://localhost:9090 | — |

---

## 📁 Repository Structure

```
learnbot/
├── api-gateway/          # Go — unified REST API gateway (:8090)
│   ├── cmd/server/       # Entry point
│   ├── internal/
│   │   ├── handler/      # HTTP handlers (auth, profile, jobs, analysis)
│   │   ├── middleware/   # Auth, logging, rate limiting, CORS
│   │   └── types/        # Shared request/response types
│   └── docs/             # OpenAPI spec + Postman collection
├── resume-parser/        # Go — PDF/DOCX resume parsing service (:8080)
├── job-aggregator/       # Go — multi-source job scraping service (:8081)
├── learning-resources/   # Go — curated learning resource catalog (:8082)
├── database/             # Go — shared PostgreSQL repository layer
│   ├── migrations/       # SQL migration files (run in order)
│   ├── repository/       # Repository implementations
│   └── docs/             # Data model + ER diagram documentation
├── frontend/             # Next.js 14 — React web application (:3000)
│   ├── src/app/          # App Router pages
│   ├── src/components/   # Reusable UI components
│   ├── src/store/        # Zustand state management
│   └── e2e/              # Playwright end-to-end tests
├── infrastructure/       # IaC, Kubernetes, monitoring, scripts
│   ├── terraform/        # AWS infrastructure (ECS, RDS, S3, ALB)
│   ├── kubernetes/       # K8s manifests (alternative deployment)
│   ├── monitoring/       # Prometheus + Grafana configuration
│   └── scripts/          # Deployment and database scripts
├── docs/                 # Project-wide documentation
├── docker-compose.yml    # Local development stack
└── .env.example          # Environment variable template
```

---

## 🛠️ Services

| Service | Language | Port | Description |
|---------|----------|------|-------------|
| [api-gateway](api-gateway/README.md) | Go | 8090 | Unified REST API, auth, routing |
| [resume-parser](resume-parser/README.md) | Go | 8080 | PDF/DOCX parsing, skill extraction |
| [job-aggregator](job-aggregator/README.md) | Go | 8081 | LinkedIn/Indeed job scraping |
| [learning-resources](learning-resources/README.md) | Go | 8082 | Learning resource catalog |
| [frontend](frontend/README.md) | Next.js/TypeScript | 3000 | Web application |
| [database](database/README.md) | Go/PostgreSQL | 5432 | Shared data layer |

---

## 🧪 Running Tests

```bash
# Backend — all services
cd api-gateway && go test ./... -v
cd resume-parser && go test ./... -v
cd job-aggregator && go test ./... -v
cd learning-resources && go test ./... -v
cd database && go test ./... -v

# Frontend — unit tests
cd frontend && npm test

# Frontend — E2E tests (requires running server)
cd frontend && npm run test:e2e

# Run with coverage
cd api-gateway && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

See the full [Testing Guide](docs/TESTING.md) for detailed test documentation.

---

## 🔐 Security

- **Authentication:** JWT Bearer tokens (HS256), 24-hour expiry
- **Rate Limiting:** 10 req/s per IP, burst of 30
- **Encryption:** TLS 1.3 in production, AES-256 for data at rest
- **Secrets:** AWS Secrets Manager in production, `.env` locally
- **Compliance:** GDPR-ready with soft deletes and audit logs

See [Security section](docs/DEPLOYMENT_GUIDE.md#security) in the Deployment Guide.

---

## 🚢 Deployment

LearnBot deploys to AWS using:
- **ECS Fargate** for containerized services
- **RDS PostgreSQL** (Multi-AZ) for the database
- **S3** for resume file storage
- **CloudFront** for frontend CDN
- **ALB** for load balancing and TLS termination

See the full [Deployment Guide](docs/DEPLOYMENT_GUIDE.md) for step-by-step instructions.

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Follow the [Developer Guide](docs/DEVELOPER_GUIDE.md) for code style and testing requirements
4. Submit a pull request to `main`

---

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

## 📞 Support

- **Documentation:** [docs/](docs/)
- **User Help Center:** [docs/help/](docs/help/)
- **Issues:** [GitHub Issues](https://github.com/learnbot/learnbot/issues)
- **Operations:** [Infrastructure Runbook](infrastructure/docs/RUNBOOK.md)
