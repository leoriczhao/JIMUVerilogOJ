# Repository Guidelines

## Project Structure & Module Organization
- `backend/` hosts the main Go service; use `cmd/main.go` as the entry point and keep application logic inside `internal/` (handlers, services, middleware) with shared helpers under `pkg/`.
- `judge-service/` contains the asynchronous judging pipeline; mirror the backend layout when adding queue or grading logic.
- `frontend/` is reserved for the Vue 3 client; store UI assets and route modules there when work begins.
- `tests/` provides a Python API test suite managed by `uv`; `docs/`, `docker/`, and `scripts/` hold reference material, container specs, and automation such as `deploy.sh`.

## Build, Test, and Development Commands
- `cd backend && make build` compiles the Go API to `backend/main`; `make run` builds and starts it locally.
- `make test`, `make test-verbose`, and `make test-coverage` execute Go unit tests with optional verbosity or coverage output.
- `make check` chains `fmt`, `vet`, `lint`, and `test`; run it before every push.
- `make dev-start` / `make dev-stop` spin up or tear down the Docker-based dev stack defined in `docker-compose.dev.yml`.

## Coding Style & Naming Conventions
- Format Go code with `make fmt` (gofmt) before committing; `golangci-lint` backs the `make lint` target.
- Follow Go package norms: exported items use `CamelCase`, private helpers use `camelCase`, and test files end with `_test.go`.
- Python fixtures in `tests/` follow PEP 8 `snake_case`; keep modules focused on a single API area.

## Testing Guidelines
- Go tests live alongside their packages; add new cases using `func TestXxx(t *testing.T)` names so `go test ./...` discovers them.
- For end-to-end checks, `cd tests && uv sync` once, then `./run_tests.sh` or `uv run python test_all.py --modules <area>` against a running backend.
- Store coverage artifacts locally: `backend/coverage.out` and HTML reports remain untracked.

## Commit & Pull Request Guidelines
- Match the existing Conventional Commit style (`feat:`, `fix:`, `chore:`) with concise, imperative subjects under 72 characters.
- Open pull requests with a clear summary, linked issues, verification steps (`make check`, `./run_tests.sh`), and screenshots for UI work.
- Keep changes scoped to one feature or bugfix; update `docs/` when behavior or endpoints shift.

## Security & Configuration Tips
- Base environment variables on `docs/environment.md`; copy to `.env.dev` or `.env.prod` and never commit secrets.
- Validate new dependencies for vulnerabilities with `make security-check`, and review Docker changes for exposed ports or credentials.
