# Repository Guidelines

Use this guide to onboard quickly and keep updates consistent across services.

## Project Structure & Module Organization
- `backend/` hosts the Go API; start from `cmd/main.go` and keep request logic in `internal/handlers`, business rules in `internal/services`, middleware under `internal/middleware`, and shared helpers in `pkg/`.
- `judge-service/` mirrors the backend layout for asynchronous judging; queue clients live in `internal/queue`, grading strategies in `internal/grader`.
- `frontend/` is reserved for the upcoming Vue 3 SPA; place route modules inside `frontend/src/router` and shared UI assets inside `frontend/src/assets`.
- `tests/` contains the Python API acceptance suite; `docs/`, `docker/`, and `scripts/` collect reference material, container specs, and automation such as `scripts/deploy.sh`.

## Build, Test, and Development Commands
- `cd backend && make build` compiles the Go API into `backend/main`; use `make run` during local iteration.
- `cd backend && make test` (or `make test-verbose`) runs unit tests; `make test-coverage` writes `coverage.out`.
- `cd backend && make check` chains `fmt`, `vet`, `lint`, and `test`. Run it before commits land in main.
- `make dev-start` and `make dev-stop` boot or tear down the Docker compose stack defined in `docker-compose.dev.yml`.
- `make security-check` scans dependencies for known vulnerabilities.

## Coding Style & Naming Conventions
Run `cd backend && make fmt` (gofmt) to enforce formatting; `golangci-lint` backs `make lint`. Use `CamelCase` for exported Go symbols, `camelCase` for private helpers, and name tests `*_test.go`. Python fixtures in `tests/` follow PEP 8 `snake_case`.

## Testing Guidelines
Go tests sit beside their packages; define cases as `func TestXxx(t *testing.T)` so `go test ./...` discovers them. For API flows, run `cd tests && uv sync` once, then `./run_tests.sh` or `uv run python test_all.py --modules submissions` against a running backend. Keep generated coverage reports untracked.

## Commit & Pull Request Guidelines
Follow Conventional Commits such as `feat: add judge timeout` or `fix: handle submission retry`. Pull requests should summarize the change, link issues, include verification steps (`make check`, `./run_tests.sh`), and attach screenshots for UI updates. Keep each PR focused and update `docs/` when behavior or endpoints change.

## Security & Configuration Tips
Base environment values on `docs/environment.md`, copying into `.env.dev` or `.env.prod` without committing secrets. Audit container edits for exposed ports or credentials, and prefer `make security-check` before deploying new dependencies.
