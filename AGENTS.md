# Repository Instructions

## Project Shape
- VaporAuror is a small three-tier game-store app: static vanilla HTML/CSS/JS frontend, Go/Gin REST API, PostgreSQL database.
- Backend entrypoint is `backend/main.go`; route registration lives in `backend/routes/routes.go`.
- Frontend entrypoint is `frontend/index.html`; Caddy config is `frontend/Caddyfile`; shared API calls must go through `frontend/assets/js/api.js`, which uses same-origin relative URLs (`API_BASE = ''`).
- Frontend is static HTML/CSS/JS with Bulma loaded from `frontend/assets/css/style.css`; there is no frontend build step.
- Caddy is the browser-facing entrypoint: `/api/*` and `/media/*` reverse proxy to the backend service, so frontend code should not hardcode backend hosts.
- Backend serves images at `/media/images/*` from `backend/assets/images`; game files live under `backend/assets/game-files`, are stored as `/downloads/...` URLs in `game_media`, and are only returned by `GET /api/protected/library/{game_id}/download` after license checks.
- Developer game editing lives at `frontend/pages/dashboard/edit_game.html`; `dev_dashboard.html` should only list/create games and link to that editor.
- Game descriptions are stored in `games.description`, exposed as JSON `desc`, edited via `PUT /api/developer/games/:id`, and rendered as Markdown with `marked` + `DOMPurify` on detail/edit pages.
- PostgreSQL schema and seed data are `db/01_init_table.sql` then `db/02_init_data.sql`, mounted by Compose into `/docker-entrypoint-initdb.d/`.

## Commands
- Start the full stack from repo root: `docker compose up -d --build`; Caddy is the browser entrypoint at `http://localhost:3000`, Adminer is `:8080`, Postgres is `:5432`, and backend is internal as `backend:8000`.
- Run backend locally from `backend`: `go run .`; it reads `DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME` and defaults to `localhost:5432`, `admin/admin`, database `vapor_auror`.
- Verify backend compile/tests from `backend`: `go test ./...`.
- Frontend has no Node build step or npm scripts; use the Caddy service in Compose instead of `npm start`.

## Runtime Gotchas
- There are no frontend lint/test/typecheck scripts; do not invent them.
- There are currently no Go test files, but `go test ./...` is still the fastest backend verification.
- `backend/main.go` rewrites any non-bcrypt seed password hashes to the bcrypt hash for password `admin` on startup, so seeded accounts log in with `admin`.
- `backend/main.go` also ensures the `games.description` column exists and normalizes legacy media/download URLs on startup; check it before changing seeded media paths or game columns.
- Role gates are enforced by `middleware.RequireRole`; `ADMIN` bypasses role-specific checks.

## Source Of Truth
- Trust executable routes in `backend/routes/routes.go` over API prose if paths or response details disagree.
- Use `api/api_spec.md` and `api/api_list.txt` for endpoint intent, but confirm implementations in controllers before changing behavior.
- Architecture docs under `architecture/` are useful orientation, but contain stale details such as a non-existent CORS middleware and old file names.
