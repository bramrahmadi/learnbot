# LearnBot — Developer Guide

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Setup Instructions](#2-setup-instructions)
3. [Development Workflow](#3-development-workflow)
4. [Project Structure](#4-project-structure)
5. [Code Style Guide](#5-code-style-guide)
6. [Testing Guidelines](#6-testing-guidelines)
7. [Environment Variables](#7-environment-variables)
8. [Database Development](#8-database-development)
9. [Adding a New Feature](#9-adding-a-new-feature)
10. [Debugging](#10-debugging)
11. [Common Issues](#11-common-issues)

---

## 1. Prerequisites

### Required Tools

| Tool | Version | Install |
|------|---------|---------|
| Go | >= 1.21 | https://go.dev/dl/ |
| Node.js | >= 20 LTS | https://nodejs.org/ |
| Docker | >= 24.0 | https://docs.docker.com/get-docker/ |
| Docker Compose | >= 2.0 | Included with Docker Desktop |
| Git | >= 2.40 | https://git-scm.com/ |
| PostgreSQL client | >= 15 | `brew install postgresql` / `apt install postgresql-client` |

### Optional Tools

| Tool | Purpose |
|------|---------|
| `golangci-lint` | Go linting (`brew install golangci-lint`) |
| `migrate` | Database migrations CLI (`go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`) |
| `air` | Go hot reload (`go install github.com/air-verse/air@latest`) |
| Postman | API testing (import `api-gateway/docs/learnbot.postman_collection.json`) |

---

## 2. Setup Instructions

### 2.1 Clone the Repository

```bash
git clone https://github.com/learnbot/learnbot.git
cd learnbot
```

### 2.2 Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your local values. The defaults work for local Docker development:

```bash
# Minimum required changes for local dev:
JWT_SECRET=your-local-dev-secret-at-least-32-chars
```

### 2.3 Start Infrastructure Services

Start only the database and Redis (without application services):

```bash
docker compose up -d postgres redis
```

Wait for services to be healthy:
```bash
docker compose ps
# postgres should show "healthy"
```

### 2.4 Run Database Migrations

Migrations run automatically when PostgreSQL starts (via `docker-entrypoint-initdb.d`). To run them manually:

```bash
# Run all migrations in order
for f in database/migrations/*.sql; do
  docker compose exec -T postgres psql -U learnbot_admin -d learnbot -f /dev/stdin < "$f"
done
```

Or using `golang-migrate`:
```bash
migrate -path ./database/migrations \
  -database "postgres://learnbot_admin:localdevpassword@localhost:5432/learnbot?sslmode=disable" \
  up
```

### 2.5 Start Backend Services

**Option A: Docker Compose (recommended for full stack)**
```bash
docker compose up -d
```

**Option B: Run services locally (for active development)**

```bash
# Terminal 1: API Gateway
cd api-gateway
go run ./cmd/server -addr :8090

# Terminal 2: Resume Parser
cd resume-parser
go run ./cmd/server -addr :8080

# Terminal 3: Job Aggregator
cd job-aggregator
go run ./cmd/server -addr :8081

# Terminal 4: Learning Resources
cd learning-resources
go run ./cmd/server -addr :8082
```

### 2.6 Start Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at http://localhost:3000.

### 2.7 Verify Setup

```bash
# Check API Gateway health
curl http://localhost:8090/health
# Expected: {"status":"ok","service":"api-gateway","version":"1.0.0"}

# Check frontend
open http://localhost:3000

# Check monitoring
open http://localhost:3001  # Grafana (admin/admin)
open http://localhost:9090  # Prometheus
```

---

## 3. Development Workflow

### 3.1 Branch Strategy

```
main                    ← Production-ready code
  └── feature/*         ← New features
  └── fix/*             ← Bug fixes
  └── chore/*           ← Maintenance, docs, refactoring
  └── hotfix/*          ← Critical production fixes
```

**Branch naming:**
- `feature/add-linkedin-oauth`
- `fix/resume-parser-pdf-encoding`
- `chore/update-go-dependencies`
- `hotfix/jwt-expiry-bug`

### 3.2 Development Cycle

```bash
# 1. Create a feature branch
git checkout main
git pull origin main
git checkout -b feature/your-feature

# 2. Make changes
# ... edit files ...

# 3. Run tests
cd api-gateway && go test ./... -v
cd frontend && npm test

# 4. Lint
cd api-gateway && go vet ./...
cd frontend && npm run lint

# 5. Commit
git add .
git commit -m "feat: add LinkedIn OAuth integration"

# 6. Push and create PR
git push origin feature/your-feature
# Open PR on GitHub targeting main
```

### 3.3 Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
| Type | When to Use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no logic change |
| `refactor` | Code restructuring, no feature/fix |
| `test` | Adding or fixing tests |
| `chore` | Build process, dependencies |
| `perf` | Performance improvement |

**Examples:**
```
feat(api-gateway): add job recommendations endpoint
fix(resume-parser): handle UTF-8 encoded PDFs correctly
docs: add API authentication guide
test(frontend): add E2E tests for job search flow
chore(deps): update Go to 1.22
```

### 3.4 Pull Request Process

1. Ensure all tests pass locally
2. Update documentation if adding/changing features
3. Fill out the PR template
4. Request review from at least one team member
5. Address review comments
6. Squash and merge to `main`

### 3.5 Hot Reload Development

For faster Go development, use `air`:

```bash
# Install air
go install github.com/air-verse/air@latest

# Run with hot reload (from service directory)
cd api-gateway
air
```

For frontend, `npm run dev` already includes hot reload via Next.js.

---

## 4. Project Structure

### Backend Services (Go)

Each Go service follows the same structure:

```
service-name/
├── cmd/
│   └── server/
│       └── main.go          # Entry point, server setup
├── internal/
│   ├── handler/             # HTTP handlers (one file per domain)
│   ├── middleware/          # HTTP middleware
│   ├── model/               # Domain models/types
│   ├── service/             # Business logic (if complex)
│   └── repository/          # Data access layer
├── Dockerfile
├── go.mod
└── go.sum
```

### Frontend (Next.js)

```
frontend/
├── src/
│   ├── app/                 # Next.js App Router pages
│   │   ├── layout.tsx       # Root layout
│   │   ├── page.tsx         # Landing page
│   │   ├── (auth)/          # Auth route group (no layout)
│   │   │   ├── login/
│   │   │   └── register/
│   │   ├── dashboard/
│   │   ├── jobs/
│   │   ├── analysis/
│   │   └── learning/
│   ├── components/
│   │   ├── layout/          # Layout components (Navbar, etc.)
│   │   └── ui/              # Reusable UI primitives
│   ├── lib/
│   │   └── api.ts           # API client
│   └── store/
│       └── authStore.ts     # Zustand auth state
├── e2e/                     # Playwright E2E tests
└── public/                  # Static assets
```

---

## 5. Code Style Guide

### Go

#### Formatting

- Use `gofmt` (enforced by CI): `gofmt -w .`
- Use `goimports` for import organization: `goimports -w .`
- Line length: 100 characters (soft limit)

#### Naming Conventions

```go
// Packages: lowercase, single word
package handler

// Exported types: PascalCase
type UserProfile struct { ... }

// Unexported types: camelCase
type authConfig struct { ... }

// Constants: PascalCase for exported, camelCase for unexported
const MaxFileSize = 10 * 1024 * 1024
const defaultTimeout = 30 * time.Second

// Interfaces: noun or adjective + "er"
type JobRepository interface { ... }
type Embedder interface { ... }

// Error variables: ErrXxx
var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
```

#### Error Handling

```go
// Always handle errors explicitly
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething: %w", err)
}

// Use sentinel errors for known conditions
if errors.Is(err, ErrNotFound) {
    http.Error(w, "not found", http.StatusNotFound)
    return
}

// Wrap errors with context
if err := db.QueryRow(...).Scan(&id); err != nil {
    return fmt.Errorf("getUserByEmail %q: %w", email, err)
}
```

#### HTTP Handlers

```go
// Handler functions follow this pattern:
func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
    // 1. Parse and validate request
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON")
        return
    }
    
    // 2. Validate fields
    if req.Email == "" || req.Password == "" {
        writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "email and password required")
        return
    }
    
    // 3. Business logic
    token, user, err := h.authenticate(req.Email, req.Password)
    if err != nil {
        if errors.Is(err, ErrInvalidCredentials) {
            writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
            return
        }
        writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "authentication failed")
        return
    }
    
    // 4. Write response
    writeJSON(w, http.StatusOK, map[string]interface{}{
        "token": token,
        "user":  user,
    })
}
```

#### Testing

```go
// Test function naming: TestFunctionName_Scenario
func TestLogin_Success(t *testing.T) { ... }
func TestLogin_InvalidCredentials(t *testing.T) { ... }
func TestLogin_MissingFields(t *testing.T) { ... }

// Use table-driven tests for multiple scenarios
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"empty", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
            }
        })
    }
}
```

### TypeScript/React

#### Formatting

- Use Prettier (configured in `frontend/.prettierrc` if present)
- ESLint: `npm run lint`
- TypeScript strict mode enabled

#### Component Structure

```tsx
// components/ui/Button.tsx

import React from 'react';

// 1. Types/interfaces first
interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
  className?: string;
}

// 2. Component function
export function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  onClick,
  children,
  className = '',
}: ButtonProps) {
  // 3. Derived state/computed values
  const isDisabled = disabled || loading;
  
  // 4. Event handlers
  const handleClick = () => {
    if (!isDisabled && onClick) {
      onClick();
    }
  };
  
  // 5. Render
  return (
    <button
      onClick={handleClick}
      disabled={isDisabled}
      className={`btn btn-${variant} btn-${size} ${className}`}
    >
      {loading ? <LoadingSpinner size="sm" /> : children}
    </button>
  );
}
```

#### API Calls

Always use the centralized API client in `src/lib/api.ts`:

```typescript
// ✅ Correct
import { api } from '@/lib/api';
const profile = await api.getProfile();

// ❌ Incorrect — don't use fetch directly in components
const response = await fetch('/api/users/profile', {
  headers: { Authorization: `Bearer ${token}` }
});
```

#### State Management

Use Zustand for global state. Keep component-local state in `useState`:

```typescript
// Global state (auth, user profile) → Zustand store
import { useAuthStore } from '@/store/authStore';
const { user, token, login } = useAuthStore();

// Local UI state → useState
const [isMenuOpen, setIsMenuOpen] = useState(false);
```

---

## 6. Testing Guidelines

### Backend (Go)

#### Running Tests

```bash
# All tests in a service
cd api-gateway && go test ./... -v

# With coverage
cd api-gateway && go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
cd resume-parser && go test ./... -bench=. -benchmem -run=^$

# Run specific test
cd api-gateway && go test ./internal/handler/... -run TestLogin_Success -v
```

#### Test Coverage Targets

| Service | Target |
|---------|--------|
| `resume-parser` | 80%+ |
| `api-gateway` | 75%+ |
| `job-aggregator` | 70%+ |
| `learning-resources` | 70%+ |

#### Writing Tests

- Use table-driven tests for multiple input scenarios
- Test both happy path and error cases
- Use `httptest.NewRecorder()` for HTTP handler tests
- No external dependencies in unit tests (use in-memory stores)

```go
// Integration test example (api-gateway)
func TestRegister_Success(t *testing.T) {
    srv := newTestServer(t)
    
    body := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
    req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    srv.ServeHTTP(w, req)
    
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
    }
    
    var resp map[string]interface{}
    json.NewDecoder(w.Body).Decode(&resp)
    
    data := resp["data"].(map[string]interface{})
    if data["token"] == nil {
        t.Error("expected token in response")
    }
}
```

### Frontend (TypeScript)

#### Running Tests

```bash
cd frontend

# Unit tests
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage

# Specific file
npm test -- --testPathPattern="Button"

# E2E tests (requires running server)
npm run test:e2e

# E2E with UI
npm run test:e2e:ui
```

#### Writing Component Tests

```tsx
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from '../Button';

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const onClick = jest.fn();
    render(<Button onClick={onClick}>Click me</Button>);
    fireEvent.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('does not call onClick when disabled', () => {
    const onClick = jest.fn();
    render(<Button onClick={onClick} disabled>Click me</Button>);
    fireEvent.click(screen.getByRole('button'));
    expect(onClick).not.toHaveBeenCalled();
  });
});
```

#### Writing E2E Tests

```typescript
import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('user can register and login', async ({ page }) => {
    await page.goto('/register');
    
    await page.getByLabel('Full Name').fill('Jane Doe');
    await page.getByLabel('Email').fill('jane@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Create Account' }).click();
    
    await expect(page).toHaveURL('/dashboard');
    await expect(page.getByText('Welcome, Jane')).toBeVisible();
  });
});
```

---

## 7. Environment Variables

### Root `.env` (Docker Compose)

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | `local-dev-jwt-secret-change-in-production` | JWT signing secret (min 32 chars) |
| `DB_PASSWORD` | `localdevpassword` | PostgreSQL password |
| `REDIS_PASSWORD` | `localredispassword` | Redis password |
| `GRAFANA_USER` | `admin` | Grafana admin username |
| `GRAFANA_PASSWORD` | `admin` | Grafana admin password |

### API Gateway

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8090` | HTTP server port |
| `ENVIRONMENT` | `development` | `development` or `production` |
| `JWT_SECRET` | — | JWT signing secret |
| `DB_HOST` | `postgres` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `learnbot` | Database name |
| `DB_USER` | `learnbot_admin` | Database user |
| `DB_PASSWORD` | — | Database password |
| `RESUME_PARSER_URL` | `http://resume-parser:8080` | Resume parser service URL |
| `JOB_AGGREGATOR_URL` | `http://job-aggregator:8081` | Job aggregator service URL |
| `LEARNING_RESOURCES_URL` | `http://learning-resources:8082` | Learning resources service URL |
| `REDIS_URL` | `redis://:password@redis:6379/0` | Redis connection URL |

### Frontend

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8090` | API Gateway URL (public) |
| `NODE_ENV` | `development` | Node environment |

---

## 8. Database Development

### Running Migrations

```bash
# Apply all migrations
for f in database/migrations/*.sql; do
  psql "postgres://learnbot_admin:localdevpassword@localhost:5432/learnbot" -f "$f"
done

# Apply specific migration
psql "postgres://learnbot_admin:localdevpassword@localhost:5432/learnbot" \
  -f database/migrations/001_create_enums.sql
```

### Migration Naming Convention

```
NNN_description.sql
001_create_enums.sql
002_create_users.sql
003_create_profiles.sql
```

Always increment the number. Never modify existing migrations — create a new one.

### Connecting to the Database

```bash
# Via Docker
docker compose exec postgres psql -U learnbot_admin -d learnbot

# Direct connection
psql "postgres://learnbot_admin:localdevpassword@localhost:5432/learnbot"
```

### Useful Queries

```sql
-- Check profile completeness
SELECT user_id, profile_completeness FROM user_profiles ORDER BY profile_completeness DESC;

-- View all skills for a user
SELECT skill_name, proficiency, years_of_experience 
FROM user_skills 
WHERE user_id = 'your-user-uuid'
ORDER BY is_primary DESC, proficiency DESC;

-- Check active skill gaps
SELECT * FROM v_active_skill_gaps WHERE user_id = 'your-user-uuid';

-- View full profile
SELECT * FROM v_user_full_profile WHERE user_id = 'your-user-uuid';
```

---

## 9. Adding a New Feature

### Example: Adding a New API Endpoint

**Step 1: Define the types** in `api-gateway/internal/types/types.go`:
```go
type NewFeatureRequest struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

type NewFeatureResponse struct {
    Result string `json:"result"`
}
```

**Step 2: Create the handler** in `api-gateway/internal/handler/newfeature.go`:
```go
package handler

type NewFeatureHandler struct {
    // dependencies
}

func NewNewFeatureHandler() *NewFeatureHandler {
    return &NewFeatureHandler{}
}

func (h *NewFeatureHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
    mux.Handle("POST /api/newfeature", auth(http.HandlerFunc(h.handleNewFeature)))
}

func (h *NewFeatureHandler) handleNewFeature(w http.ResponseWriter, r *http.Request) {
    var req types.NewFeatureRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON")
        return
    }
    // ... business logic ...
    writeJSON(w, http.StatusOK, types.NewFeatureResponse{Result: "success"})
}
```

**Step 3: Register in main.go**:
```go
newFeatureHandler := handler.NewNewFeatureHandler()
newFeatureHandler.RegisterRoutes(mux, authMiddleware)
```

**Step 4: Write tests** in `api-gateway/internal/handler/integration_test.go`:
```go
func TestNewFeature_Success(t *testing.T) {
    // ...
}
```

**Step 5: Update OpenAPI spec** in `api-gateway/docs/openapi.yaml`.

**Step 6: Update API documentation** in `docs/API.md`.

---

## 10. Debugging

### Go Services

```bash
# Enable verbose logging
ENVIRONMENT=development go run ./cmd/server

# Check goroutine leaks
go test -race ./...

# Profile CPU
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Profile memory
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

### Docker Services

```bash
# View logs
docker compose logs -f api-gateway
docker compose logs -f resume-parser

# Exec into container
docker compose exec api-gateway sh

# Restart a service
docker compose restart api-gateway

# Rebuild after code changes
docker compose up -d --build api-gateway
```

### Frontend

```bash
# Check TypeScript errors
cd frontend && npx tsc --noEmit

# Debug Next.js build
cd frontend && npm run build

# Analyze bundle size
cd frontend && ANALYZE=true npm run build
```

---

## 11. Common Issues

### `go: module not found`

The `api-gateway` and `learning-resources` modules use `replace` directives in `go.mod`. Always build from the **repo root** when using Docker:

```bash
# ✅ Correct (from repo root)
docker compose build api-gateway

# ❌ Incorrect (from service directory)
cd api-gateway && docker build .
```

### Database connection refused

Ensure PostgreSQL is healthy before starting application services:

```bash
docker compose up -d postgres
docker compose ps  # Wait for "healthy" status
docker compose up -d api-gateway
```

### JWT token invalid after restart

The JWT secret must be consistent. If you restart with a different `JWT_SECRET`, all existing tokens become invalid. Set a stable secret in `.env`.

### Frontend can't reach API

Check `NEXT_PUBLIC_API_URL` in `frontend/.env.local`:
```bash
# For local development
NEXT_PUBLIC_API_URL=http://localhost:8090
```

### React 19 + Testing Library compatibility

If you see `React.act is not a function` in tests, ensure `jest.polyfill.js` is loaded. This is configured in `frontend/jest.config.js` via `setupFiles`.

### Port already in use

```bash
# Find and kill process on port 8090
lsof -ti:8090 | xargs kill -9

# Or use different ports
PORT=8091 go run ./cmd/server
```
