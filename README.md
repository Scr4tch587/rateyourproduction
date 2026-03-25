# RateYourProduction

A RateYourMusic-inspired theatre logging and discovery platform for plays, musicals, and operas.

## Stack

- **Frontend**: Next.js (TypeScript, Tailwind, shadcn/ui)
- **Backend**: Go (chi router)
- **Database**: PostgreSQL (Supabase)
- **Cache**: Redis
- **Auth**: Supabase Auth
- **Hosting**: Vercel (frontend), Railway (backend)

## Project Structure

```
apps/web/        → Next.js frontend
apps/api/        → Go API server
migrations/      → SQL migration files
infra/           → Docker Compose for local dev
docs/            → Project documentation
```

## Getting Started

### Prerequisites

- Node.js 20+
- Go 1.22+
- Docker & Docker Compose

### Setup

```bash
make setup           # Install dependencies, copy .env
make dev             # Start Postgres + Redis via Docker
make dev-api         # Start Go API server (separate terminal)
make dev-web         # Start Next.js dev server (separate terminal)
```

### Environment

Copy `.env.example` to `.env` and fill in values. Each app also has its own `.env.example`.

### Database

```bash
make migrate-up      # Run migrations against DATABASE_URL
```

## API Routes

```
GET  /health
GET  /api/v1/works
GET  /api/v1/works/:slug
GET  /api/v1/productions
GET  /api/v1/productions/:id
GET  /api/v1/logs
POST /api/v1/logs
GET  /api/v1/discover
GET  /api/v1/submissions
POST /api/v1/submissions
GET  /api/v1/admin/works
GET  /api/v1/admin/productions
GET  /api/v1/admin/submissions
```

## Frontend Routes

```
/                        → Home
/discover                → Discovery with filters
/works/:slug             → Work page
/productions/:id         → Production page
/profile/:username       → User profile
/auth/login              → Login
/auth/signup             → Sign up
```
