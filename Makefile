.PHONY: dev dev-web dev-api build-api migrate-up migrate-down lint setup

# Local development
dev:
	@echo "Starting Redis..."
	docker compose -f infra/docker-compose.yml up -d
	@echo "Redis running on :6379. Start Supabase separately with 'supabase start'."

dev-web:
	cd apps/web && npm run dev

dev-api:
	cd apps/api && go run ./cmd/server

build-api:
	cd apps/api && go build -o bin/server ./cmd/server

# Database
migrate-up:
	@echo "Run migrations against your database"
	psql "$(DATABASE_URL)" -f migrations/001_initial_schema.sql

migrate-down:
	@echo "Run down migration"
	psql "$(DATABASE_URL)" -f migrations/001_initial_schema_down.sql

# Quality
lint:
	cd apps/web && npm run lint
	cd apps/api && golangci-lint run ./...

# Setup
setup:
	cd apps/web && npm install
	cd apps/api && go mod download
	cp -n .env.example .env || true
	@echo "Setup complete. Run 'make dev' to start local services."
