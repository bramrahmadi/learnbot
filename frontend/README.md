# LearnBot Frontend

The LearnBot web application built with Next.js 14 (App Router), TypeScript, and Tailwind CSS.

Part of the [LearnBot AI Career Development Platform](../README.md).

## Features

- **Landing Page** — Marketing page with feature overview
- **Authentication** — Registration and login with JWT
- **Onboarding** — Resume upload and profile setup flow
- **Dashboard** — Profile overview, job recommendations, learning progress
- **Jobs** — Browse and search job listings with acceptance likelihood scores
- **Job Detail** — Full job description with match breakdown
- **Analysis** — Skill gap analysis with visual breakdown
- **Learning** — Personalized learning plan with curated resources

## Tech Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| Next.js | 14 | React framework (App Router) |
| TypeScript | 5 | Type safety |
| Tailwind CSS | 3 | Utility-first styling |
| Zustand | 4 | State management |
| Jest | 29 | Unit testing |
| React Testing Library | 16 | Component testing |
| Playwright | Latest | E2E testing |

## Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── layout.tsx          # Root layout (Navbar, global styles)
│   │   ├── page.tsx            # Landing page
│   │   ├── globals.css         # Global CSS + Tailwind directives
│   │   ├── (auth)/             # Auth route group (no shared layout)
│   │   │   ├── login/page.tsx
│   │   │   └── register/page.tsx
│   │   ├── dashboard/page.tsx  # Main dashboard
│   │   ├── onboarding/page.tsx # Resume upload + profile setup
│   │   ├── jobs/
│   │   │   ├── page.tsx        # Job listings
│   │   │   └── [id]/page.tsx   # Job detail
│   │   ├── analysis/page.tsx   # Skill gap analysis
│   │   └── learning/page.tsx   # Learning plan
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Navbar.tsx      # Navigation bar
│   │   │   └── __tests__/
│   │   └── ui/                 # Reusable UI primitives
│   │       ├── Badge.tsx       # Status badges, score badges
│   │       ├── Button.tsx      # Button with variants and loading state
│   │       ├── Card.tsx        # Card container with header
│   │       ├── Input.tsx       # Input, textarea, select
│   │       ├── LoadingSpinner.tsx  # Spinner, PageLoader, InlineLoader
│   │       └── __tests__/
│   ├── lib/
│   │   └── api.ts              # Centralized API client
│   └── store/
│       ├── authStore.ts        # Zustand auth state
│       └── __tests__/
├── e2e/                        # Playwright E2E tests
│   ├── auth.spec.ts
│   ├── jobs.spec.ts
│   ├── navigation.spec.ts
│   └── resume-upload.spec.ts
├── public/                     # Static assets
├── jest.config.js
├── jest.polyfill.js            # React 19 + testing-library compatibility
├── jest.setup.ts
├── playwright.config.ts
├── next.config.js
├── tailwind.config.ts
└── tsconfig.json
```

## Quick Start

### Prerequisites

- Node.js >= 20 LTS
- npm >= 10

### Development

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env.local
# Set NEXT_PUBLIC_API_URL=http://localhost:8090

# Start development server
npm run dev
```

Open http://localhost:3000.

### Production Build

```bash
npm run build
npm start
```

### Docker

```bash
# From repo root
docker compose up -d frontend
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8090` | API Gateway URL |
| `NODE_ENV` | `development` | Node environment |

## Running Tests

### Unit Tests

```bash
# Run all unit tests
npm test

# Watch mode
npm test -- --watch

# With coverage
npm test -- --coverage

# Specific file
npm test -- --testPathPattern="Button"
```

### E2E Tests

```bash
# Requires running frontend server (npm run dev)

# Run all E2E tests
npm run test:e2e

# Run with Playwright UI
npm run test:e2e:ui

# Run specific test file
npx playwright test e2e/auth.spec.ts
```

### Type Checking

```bash
npx tsc --noEmit
```

### Linting

```bash
npm run lint
```

## Component Library

### UI Components

| Component | Props | Description |
|-----------|-------|-------------|
| `Button` | `variant`, `size`, `loading`, `disabled` | Button with primary/secondary/danger/ghost variants |
| `Card` | `hover`, `padding` | Card container |
| `CardHeader` | `title`, `subtitle`, `action`, `icon` | Card header section |
| `Input` | `label`, `error`, `hint`, `required` | Text input with label and validation |
| `Textarea` | Same as Input | Multi-line text input |
| `Select` | `label`, `options`, `error` | Dropdown select |
| `LoadingSpinner` | `size`, `label` | Animated loading indicator |
| `PageLoader` | `message` | Full-page loading overlay |
| `InlineLoader` | `message` | Inline loading indicator |
| `Badge` | `color` | Status badge |
| `ScoreBadge` | `score` | Acceptance likelihood score badge |
| `GapCategoryBadge` | `category` | Skill gap priority badge |
| `DifficultyBadge` | `difficulty` | Learning resource difficulty badge |
| `CostBadge` | `isFree`, `price` | Resource cost badge |

### State Management

The app uses Zustand for global state. The auth store (`src/store/authStore.ts`) manages:

```typescript
