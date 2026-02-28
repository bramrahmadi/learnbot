# LearnBot Testing Documentation

## Overview

This document describes the comprehensive testing strategy for the LearnBot MVP platform. The test suite covers backend Go services, frontend React components, and end-to-end user flows.

## Test Architecture

```
LearnBot/
├── resume-parser/          # Go service - resume parsing & analysis
│   ├── internal/
│   │   ├── extractor/      # Unit tests for extraction functions
│   │   ├── gapanalysis/    # Unit + performance tests for gap analysis
│   │   ├── parser/         # Unit + performance tests for parsers
│   │   ├── recommendation/ # Unit tests for recommendation engine
│   │   └── scorer/         # Unit + performance tests for scoring
│   └── internal/api/       # Integration tests for API handlers
├── api-gateway/            # Go service - API gateway
│   └── internal/handler/   # Integration + performance/load tests
├── job-aggregator/         # Go service - job aggregation
│   └── internal/           # Unit tests for scrapers, HTTP client
├── learning-resources/     # Go service - learning resources
│   └── internal/api/       # Integration tests
├── database/               # Go - database models
│   └── repository/         # Unit tests for repository models
└── frontend/               # Next.js frontend
    ├── src/
    │   ├── components/     # Component unit tests
    │   └── store/          # State management tests
    └── e2e/                # Playwright E2E tests
```

## Backend Testing (Go)

### Test Frameworks

- **Go testing package** - Built-in testing framework
- **httptest** - HTTP handler testing
- **testify** (where applicable) - Assertions

### Running Backend Tests

```bash
# Run all tests in a service
cd resume-parser && go test ./... -v

# Run with coverage
cd resume-parser && go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
cd resume-parser && go test ./... -bench=. -benchmem -run=^$

# Run specific test
cd resume-parser && go test ./internal/scorer/... -run TestCalculate_PerfectMatch -v
```

### Unit Tests

#### Resume Parser (`resume-parser/internal/`)

**Extractor Tests** (`extractor/`)
- `skills_test.go` - Tests for skill extraction and classification
- `experience_test.go` - Tests for work experience extraction
- `education_test.go` - Tests for education extraction
- `personal_test.go` - Tests for personal info extraction
- `certifications_test.go` - Tests for certification extraction
- `projects_test.go` - Tests for project extraction
- `sections_test.go` - Tests for section splitting

**Parser Tests** (`parser/`)
- `resume_parser_test.go` - Core parser tests
- `resume_parser_extra_test.go` - Additional parser edge cases
- `performance_test.go` - Performance benchmarks and tests
  - `TestNormalizeFileType_Performance` - File type normalization
  - `TestComputeOverallConfidence_*` - Confidence calculation
  - `TestGenerateWarnings_*` - Warning generation
  - `BenchmarkNormalizeFileType` - Benchmark
  - `BenchmarkComputeOverallConfidence` - Benchmark
  - `BenchmarkGenerateWarnings` - Benchmark

**Scorer Tests** (`scorer/`)
- `scorer_test.go` - Acceptance likelihood scoring tests
- `handler_test.go` - HTTP handler tests
- `performance_test.go` - Performance tests
  - `TestCalculate_ScoreInRange` - Score bounds validation
  - `TestCalculate_SkillMatchImproves` - Skill match progression
  - `TestCalculate_ExperienceMatchImproves` - Experience match progression
  - `TestCalculate_MatchedSkillsTracked` - Skill tracking
  - `TestCalculate_PreferredSkillsBonus` - Preferred skills bonus
  - `TestCalculate_LocationFit_Remote` - Location matching
  - `TestCalculate_EducationMatch` - Education level matching
  - `TestCalculate_IndustryRelevance` - Industry matching
  - `BenchmarkCalculate_*` - Performance benchmarks

**Gap Analysis Tests** (`gapanalysis/`)
- `gap_analyzer_test.go` - Comprehensive gap analysis tests
- `performance_test.go` - Performance tests
  - `TestComputePriorityScore_*` - Priority score formula tests
  - `TestTokenOverlapSimilarity_*` - String similarity tests
  - `TestAdjustLearningHours_*` - Learning hours adjustment tests
  - `TestAnalyze_ReadinessScoreInRange` - Readiness score bounds
  - `TestAnalyze_VisualDataPopulated` - Visual data generation
  - `TestAnalyze_DuplicateSkillsDeduped` - Deduplication
  - `TestAnalyze_CaseInsensitiveSkillMatching` - Case handling
  - `BenchmarkAnalyze_*` - Performance benchmarks

### Integration Tests

#### API Gateway (`api-gateway/internal/handler/`)

**Authentication Tests** (`integration_test.go`)
- `TestRegister_Success` - User registration
- `TestRegister_InvalidEmail` - Email validation
- `TestRegister_ShortPassword` - Password validation
- `TestRegister_DuplicateEmail` - Duplicate detection
- `TestLogin_Success` - User login
- `TestLogin_WrongPassword` - Invalid credentials

**Profile Tests**
- `TestGetProfile_Authenticated` - Profile retrieval
- `TestUpdateProfile_Success` - Profile update
- `TestUpdateSkills_Success` - Skills update
- `TestUpdateSkills_InvalidProficiency` - Proficiency validation

**Job Tests**
- `TestJobSearch_NoFilters` - Basic job search
- `TestJobSearch_WithLocationFilter` - Filtered search
- `TestJobDetail_ValidID` - Job detail retrieval
- `TestJobDetail_InvalidID` - 404 handling
- `TestJobMatch_ValidID` - Job matching
- `TestJobRecommendations_Authenticated` - Recommendations

**Gap Analysis Tests**
- `TestGapAnalysis_WithJobID` - Gap analysis by job ID
- `TestGapAnalysis_WithInlineJob` - Inline job analysis
- `TestGapAnalysis_NoJob` - Missing job validation

**Training Tests**
- `TestTrainingRecommendations_GET` - GET recommendations
- `TestTrainingRecommendations_POST` - POST recommendations

**Resource Tests**
- `TestResourceSearch_NoFilters` - Resource search
- `TestResourceSearch_BySkill` - Skill-filtered search
- `TestResourceSearch_FreeOnly` - Free resources filter

### Performance/Load Tests

#### API Gateway (`api-gateway/internal/handler/performance_test.go`)

**Load Tests**
- `TestConcurrentRegistrations` - 10 concurrent registrations
- `TestConcurrentJobSearches` - 20 concurrent job searches
- `TestAPIResponseTime` - Response time validation
  - Resource search: < 500ms
  - Job search: < 500ms
  - Profile retrieval: < 200ms

**Benchmarks**
- `BenchmarkRegister` - Registration endpoint
- `BenchmarkJobSearch` - Job search endpoint
- `BenchmarkResourceSearch` - Resource search endpoint
- `BenchmarkGapAnalysis` - Gap analysis endpoint

## Frontend Testing (TypeScript/React)

### Test Frameworks

- **Jest 29** - Test runner
- **@testing-library/react 16** - React component testing
- **@testing-library/jest-dom** - DOM assertions
- **@testing-library/user-event** - User interaction simulation
- **Playwright** - E2E testing

### Running Frontend Tests

```bash
cd frontend

# Run all unit tests
npm test

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test -- --testPathPattern="Button"

# Run E2E tests (requires running server)
npm run test:e2e

# Run E2E tests with UI
npm run test:e2e:ui
```

### Component Tests

#### UI Components (`src/components/ui/__tests__/`)

**Button Tests** (`Button.test.tsx`)
- Renders children correctly
- Applies variant classes (primary, secondary, danger, ghost)
- Shows loading spinner when loading=true
- Disables button when disabled or loading
- Calls onClick handler
- Applies size classes (sm, md, lg)

**Card Tests** (`Card.test.tsx`)
- Renders children
- Applies card/card-hover classes
- Applies padding variants (none, sm, md, lg)
- Applies custom className
- CardHeader: renders title, subtitle, action, icon

**Input Tests** (`Input.test.tsx`)
- Renders with/without label
- Shows required asterisk
- Shows error message with aria-invalid
- Shows hint text (hidden when error present)
- Calls onChange handler
- Generates ID from label
- Textarea: resize-none class
- Select: renders all options

**LoadingSpinner Tests** (`LoadingSpinner.test.tsx`)
- Renders with role="status"
- Custom aria-label
- Shows/hides label text by size
- Applies size classes (sm, md, lg)
- PageLoader: full-page layout
- InlineLoader: inline layout

**Badge Tests** (`Badge.test.tsx`)
- Renders children
- Applies color variants
- ScoreBadge: color by score range
- GapCategoryBadge: category labels
- DifficultyBadge: difficulty labels
- CostBadge: cost type labels

#### Layout Components (`src/components/layout/__tests__/`)

**Navbar Tests** (`Navbar.test.tsx`)
- Unauthenticated: shows Sign in/Get started links
- Unauthenticated: links to home page
- Authenticated: shows navigation links
- Authenticated: shows user name
- Authenticated: Sign out button calls logout
- Authenticated: links to dashboard
- Mobile menu toggle
- Active link highlighting

### State Management Tests

#### Auth Store (`src/store/__tests__/authStore.test.ts`)

**Initial State**
- Null token, user, profile, error
- isLoading=false

**Login**
- Sets token and user on success
- Sets error on failure
- Throws error on failure
- Calls API with correct credentials

**Register**
- Sets token and user on success
- Sets error on failure
- Calls API with correct parameters

**Logout**
- Clears token, user, profile, error

**Load Profile**
- Does nothing without token
- Loads profile with token
- Handles failure gracefully

**Update Profile**
- Does nothing without token
- Updates profile with token
- Sets error on failure

**Clear Error**
- Clears error state

### E2E Tests (Playwright)

#### Authentication (`e2e/auth.spec.ts`)
- Landing page loads with hero section
- Navigation links visible
- Register page loads with form fields
- Register form validates required fields
- Register form validates email format
- Register form validates password length
- Login page loads
- Login form validates required fields
- Cross-page navigation links

#### Navigation (`e2e/navigation.spec.ts`)
- Navbar shows correct links for unauthenticated users
- Navbar shows correct links for authenticated users

#### Resume Upload (`e2e/resume-upload.spec.ts`)
- Onboarding page redirects when not authenticated
- Dashboard shows upload prompt for new users

#### Jobs (`e2e/jobs.spec.ts`)
- Jobs page redirects when not authenticated
- Job detail page redirects when not authenticated
- Analysis page redirects when not authenticated
- Learning page redirects when not authenticated

## Test Coverage Targets

| Service | Target Coverage |
|---------|----------------|
| resume-parser | 80%+ |
| api-gateway | 75%+ |
| job-aggregator | 70%+ |
| learning-resources | 70%+ |
| frontend (unit) | 60%+ |

## CI/CD Integration

Tests run automatically on:
- Push to `main` branch
- Push to `feature/**` branches
- Pull requests to `main`

### GitHub Actions Workflow (`.github/workflows/test.yml`)

**Jobs:**
1. `test-resume-parser` - Unit + benchmark tests
2. `test-api-gateway` - Integration + load tests
3. `test-job-aggregator` - Unit tests
4. `test-learning-resources` - Integration tests
5. `test-database` - Repository tests
6. `test-frontend-unit` - Component + state tests
7. `test-frontend-e2e` - Playwright E2E tests
8. `lint-go` - Go vet on all modules
9. `typecheck-frontend` - TypeScript type checking
10. `all-tests-pass` - Summary gate

## Test Configuration

### Backend (Go)

Tests use the standard Go testing package. No external test database is required - all tests use in-memory stores.

### Frontend (Jest)

Configuration in `frontend/jest.config.js`:
- Test environment: `jest-environment-jsdom`
- React 19 compatibility polyfill in `jest.polyfill.js`
- Module aliases: `@/` → `src/`
- Coverage threshold: 50% (branches, functions, lines, statements)

### E2E (Playwright)

Configuration in `frontend/playwright.config.ts`:
- Browsers: Chromium (default), Firefox, WebKit
- Base URL: `http://localhost:3000`
- Timeout: 30 seconds per test

## Writing New Tests

### Backend Unit Test Template

```go
package mypackage

import (
    "testing"
)

func TestMyFunction_Description(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    if result != expected {
        t.Errorf("expected %q, got %q", expected, result)
    }
}

// Table-driven test
func TestMyFunction_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"empty input", "", ""},
        {"normal input", "hello", "HELLO"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}

// Benchmark
func BenchmarkMyFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MyFunction("test input")
    }
}
```

### Frontend Component Test Template

```tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { MyComponent } from "../MyComponent";

describe("MyComponent", () => {
  it("renders correctly", () => {
    render(<MyComponent title="Test" />);
    expect(screen.getByText("Test")).toBeInTheDocument();
  });
  
  it("handles user interaction", () => {
    const onClick = jest.fn();
    render(<MyComponent onClick={onClick} />);
    fireEvent.click(screen.getByRole("button"));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
```

### E2E Test Template

```typescript
import { test, expect } from "@playwright/test";

test.describe("Feature Name", () => {
  test("user can perform action", async ({ page }) => {
    await page.goto("/feature-page");
    await expect(page.getByText("Expected Text")).toBeVisible();
    await page.getByRole("button", { name: "Action" }).click();
    await expect(page).toHaveURL("/result-page");
  });
});
```

## Troubleshooting

### Common Issues

**React 19 + @testing-library/react compatibility**
- Issue: `React.act is not a function`
- Fix: The `jest.polyfill.js` file patches React to add `act` when `NODE_ENV=test`
- The `jest.config.js` sets `process.env.NODE_ENV = "test"` before `next/jest` overrides it

**Go test timeout**
- Issue: Tests timeout on slow machines
- Fix: Increase timeout with `-timeout 120s` flag

**E2E tests failing**
- Issue: Server not ready
- Fix: Use `wait-on` to wait for server before running tests
