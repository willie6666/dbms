# Repository Instructions

## Project Shape
- VaporAuror is a small three-tier game-store app: static vanilla HTML/CSS/JS frontend, Go/Gin REST API, PostgreSQL database.
- Backend entrypoint is `backend/main.go`; route registration and inline CORS live in `backend/routes/routes.go`.
- Frontend entrypoint is `frontend/index.html`; Caddy config is `frontend/Caddyfile`; shared API calls must go through `frontend/assets/js/api.js`, which hardcodes `API_BASE = 'http://localhost:8000'`.
- PostgreSQL schema and seed data are `db/01_init_table.sql` then `db/02_init_data.sql`, mounted by Compose into `/docker-entrypoint-initdb.d/`.

## Commands
- Start frontend/database/Adminer from repo root: `docker compose up -d`; Caddy serves the frontend at `http://localhost:3000`.
- Run backend from `backend`: `go run .`; it listens on `:8000` and expects Postgres at `localhost:5432` with `admin/admin`, database `vapor_auror`.
- Verify backend compile/tests from `backend`: `go test ./...`.
- Frontend has no Node build step or npm scripts; use the Caddy service in Compose instead of `npm start`.

## Runtime Gotchas
- There are no frontend lint/test/typecheck scripts; do not invent them.
- There are currently no Go test files, but `go test ./...` is still the fastest backend verification.
- `backend/main.go` rewrites any non-bcrypt seed password hashes to the bcrypt hash for password `admin` on startup, so seeded accounts log in with `admin`.
- Role gates are enforced by `middleware.RequireRole`; `ADMIN` bypasses role-specific checks.

## Source Of Truth
- Trust executable routes in `backend/routes/routes.go` over API prose if paths or response details disagree.
- Use `api/api_spec.md` and `api/api_list.txt` for endpoint intent, but confirm implementations in controllers before changing behavior.
- Architecture docs under `architecture/` are useful orientation, but they mention files that may not exist exactly as written.
