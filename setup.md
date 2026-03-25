You are setting up the initial foundation for a new full-stack project in an EMPTY repo.

I will also provide you with the MVP spec separately. That spec is the source of truth. Follow it strictly. Do not add features beyond the spec unless they are necessary scaffolding.

Your job in this pass is NOT to build the full product. Your job is to create a clean, production-minded initial setup that makes the project easy to build on.

Project context:
- Project name: RateYourProduction
- Product: a RateYourMusic-inspired theatre logging/discovery platform for plays, musicals, and operas
- Frontend: Next.js
- Backend: Go
- Database/Auth/Storage: Supabase
- Cache/queue: Redis
- Backend hosting: Railway
- Frontend hosting: Vercel
- Observability: Sentry
- UI direction: dense, information-first, dark-first, clearly inspired by RYM but modernized

Your goals for this setup:
1. Create a monorepo structure that is simple, clean, and easy to extend
2. Set up the frontend app
3. Set up the Go backend app
4. Set up shared project conventions and documentation
5. Set up local development infrastructure
6. Set up the minimum schema/migrations/models needed to begin implementing the MVP
7. Do not overengineer

Important constraints:
- Prioritize speed, clarity, and clean foundations over completeness
- Make decisions that are resume-legible and production-plausible
- Avoid unnecessary abstractions
- Avoid placeholder complexity
- Prefer boring, standard tools
- Do not implement product features unless they are necessary to validate the setup
- If something from the spec is better deferred, scaffold for it rather than fully implementing it

What I want you to produce:

1. Repo structure
Create a monorepo with a structure similar to:
- apps/web
- apps/api
- apps/worker if needed, but only if justified
- packages/shared only if clearly useful
- infra
- docs

If you think a slightly different structure is better, explain why briefly and use it.

2. Frontend setup
In apps/web:
- initialize a Next.js app with TypeScript
- configure Tailwind
- configure shadcn/ui if appropriate
- set up a dark-first global theme foundation
- create a minimal app shell/layout
- create placeholder routes/pages for the MVP’s major surfaces:
  - home
  - discover
  - work page route
  - production page route
  - profile page route
  - auth area if needed
- create a minimal, dense, RYM-inspired UI foundation
- do NOT spend time polishing visuals beyond a solid base

3. Backend setup
In apps/api:
- initialize a Go service
- choose a lightweight, standard HTTP framework or router
- set up a clean folder structure
- include:
  - config loading
  - server startup
  - health route
  - versioned API routing
  - database connection setup
  - Redis connection setup scaffold
  - basic middleware setup
  - graceful shutdown
- structure the code so it will support:
  - works
  - productions
  - logs
  - discovery
  - submissions
  - admin
- but do not fully implement these yet

4. Database and schema foundation
Using the MVP spec as source of truth:
- define the initial relational schema needed for MVP
- create migrations
- include the core entities and relationships
- include sensible indexes
- include enums/constants where appropriate
- include enough structure for:
  - works
  - productions
  - users/profiles as needed
  - genres
  - people
  - credits
  - logs
  - companies
  - venues
- if there are parts of the final schema that are too speculative for initial setup, note that and choose the most stable foundation

5. Supabase integration foundation
- assume Supabase for Postgres/Auth/Storage
- wire environment/config expectations clearly
- do not build full auth flows yet unless they are very lightweight
- make sure the project is clearly ready for Supabase-backed auth and DB usage

6. Redis foundation
- set up Redis connection/config scaffolding
- if background jobs are not needed on day 1, still prepare a clean place for them
- do not build queue logic yet unless it is minimal and justified

7. Sentry foundation
- add Sentry scaffolding for both frontend and backend if reasonable
- if full setup is too much for initial pass, at least prepare the config hooks and env expectations

8. Developer experience
Set up:
- root README with setup instructions
- per-app README if helpful
- .env.example files
- Makefile or equivalent commands if useful
- basic lint/format setup
- scripts for local dev
- Docker Compose for local dependencies if appropriate
- clear instructions for running the project locally

9. API foundation
Create an initial REST API structure with placeholder routes for:
- /health
- /api/v1/works
- /api/v1/productions
- /api/v1/logs
- /api/v1/discover
- /api/v1/submissions
- /api/v1/admin

Do not implement business logic beyond what is necessary to validate the plumbing.

10. Output expectations
I want:
- the actual files created
- a brief explanation of the architecture decisions
- a list of what was intentionally left unimplemented
- a list of the next 5 highest-leverage implementation steps

Engineering preferences:
- Keep the setup clean and practical
- Optimize for real deployment later on Vercel + Railway + Supabase
- Keep the Go service idiomatic and easy to grow
- Keep the frontend easy to iterate on quickly
- Keep naming consistent with the MVP spec
- Avoid comments unless they are genuinely useful
- Avoid adding tests unless they are setup-level smoke tests that meaningfully help

Decision guidance:
- When uncertain, prefer the simpler option
- When uncertain, align to the MVP spec
- When uncertain, choose the option that helps me ship faster while still sounding credible in interviews

Important:
- This is an EMPTY repo
- I want an initial setup only
- Do not build the entire product
- Do not invent scope outside the MVP spec
- Make the repo feel like a serious project from day 1