# Local Setup

This project now uses Supabase natively for:

- Postgres
- Auth (`auth.users`)
- profile provisioning from auth metadata

Redis is still local via Docker.

## Prerequisites

- Node.js 20+
- Go 1.22+
- Docker
- Supabase CLI

Install the Supabase CLI if needed:

```bash
brew install supabase/tap/supabase
```

## 1. Copy Environment Files

```bash
cp .env.example .env
cp apps/api/.env.example apps/api/.env
cp apps/web/.env.example apps/web/.env.local
```

For local Supabase, the defaults already point at the standard local ports:

- Supabase API: `http://127.0.0.1:54321`
- Postgres: `postgresql://postgres:postgres@localhost:54322/postgres`
- Redis: `redis://localhost:6379`
- API: `http://localhost:8080`

## 2. Install App Dependencies

```bash
make setup
```

## 3. Start Supabase

From the repo root:

```bash
supabase start
```

This starts the local Supabase stack, including Postgres and Auth.

## 4. Start Redis

```bash
make dev
```

## 5. Apply Schema and Seed Data

```bash
export DATABASE_URL=postgresql://postgres:postgres@localhost:54322/postgres
psql "$DATABASE_URL" -f migrations/001_initial_schema.sql
psql "$DATABASE_URL" -f migrations/002_seed_data.sql
```

The schema expects the Supabase `auth` schema to exist. Use Supabase Postgres, not plain standalone Postgres.

## 6. Run the API

```bash
make dev-api
```

The API uses:

- direct Postgres access for application data
- Supabase Auth `/auth/v1/user` for bearer-token verification

## 7. Run the Web App

In another terminal:

```bash
make dev-web
```

## 8. Create Your First User

Open:

- Web app: `http://localhost:3000`
- Auth signup: `http://localhost:3000/auth/signup`

Sign up with:

- username
- email
- password

The signup flow now goes directly through Supabase Auth. A database trigger creates the corresponding row in `public.profiles`.

The first signed-up user becomes `is_admin = true`.

## 9. Test the Core Flows

After signing in, verify:

1. `/discover` loads and filters.
2. A work page lets you add a log.
3. A work page lets you submit a missing production.
4. `/profile/<username>` shows your logs.
5. `/admin` loads for the first user and can approve submissions.

## Useful Supabase Commands

Stop local Supabase:

```bash
supabase stop
```

Reset local Supabase database:

```bash
supabase db reset
```

If you use `supabase db reset`, reapply any repo-specific seed/migration steps if needed.

## Notes

- The frontend auth state is managed by `@supabase/supabase-js`.
- Protected API requests send Supabase access tokens as `Authorization: Bearer <token>`.
- The Go API does not own passwords or sessions anymore.
- `infra/docker-compose.yml` now exists only for Redis.
