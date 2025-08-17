# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Verilog Online Judge system (在线判题系统) - a microservices-based platform for evaluating Verilog HDL code submissions. The system separates judging functionality from business logic for scalability and reliability.

## Architecture

The system consists of multiple services:
- **Backend** (Go): Main API service handling users, problems, submissions, forum, news
- **Judge Service** (Go): Independent microservice for code compilation and testing
- **Frontend** (Vue.js): User-facing web interface  
- **Admin Frontend** (Vue.js): Administrative interface
- **Database**: PostgreSQL for persistent data
- **Cache/Queue**: Redis for caching and async task processing

The backend uses **dependency injection via Google Wire** - run `make wire-gen` after modifying wire.go files.

## Common Development Commands

### Backend (Go)
Navigate to `backend/` directory:

```bash
# Build and run
make build
make run

# Testing
make test                    # Run all tests
make test-verbose           # Run tests with verbose output  
make test-coverage          # Run tests with coverage report
make test-services          # Run service layer tests
make test-user              # Run user service tests specifically

# Code quality
make check                  # Run all checks (fmt, vet, lint, test)
make fmt                    # Format code
make lint                   # Run golangci-lint
make vet                    # Run go vet

# Development tools
make wire-gen              # Generate dependency injection code
make deps                  # Install dependencies
make dev-setup             # Install development tools
```

### Frontend (Vue.js)
Navigate to `frontend/` directory:

```bash
npm run dev                # Start development server
npm run build              # Build for production
npm run lint               # Run ESLint with auto-fix
npm run format             # Format code with Prettier
npm run type-check         # TypeScript type checking
```

### Admin Frontend
Navigate to `admin-frontend/` directory - same commands as frontend.

### API Testing (Python)
Navigate to `tests/` directory:

```bash
# Using uv (recommended)
uv run python test_all.py

# Or with pip
pip install -r requirements.txt
python test_all.py
```

### Docker Deployment

```bash
# Development environment
./scripts/deploy.sh --dev

# Production environment  
./scripts/deploy.sh --prod

# Service management
./scripts/deploy.sh --status      # Check service status
./scripts/deploy.sh --logs        # View logs
./scripts/deploy.sh --restart     # Restart services
./scripts/deploy.sh --stop        # Stop all services
```

## Code Architecture

### Backend Structure
- **cmd/main.go**: Application entry point
- **internal/models/**: Database models (User, Problem, Submission, etc.)
- **internal/handlers/**: HTTP request handlers organized by domain
- **internal/services/**: Business logic layer
- **internal/repository/**: Data access layer
- **internal/middleware/**: HTTP middleware (auth, logging, rate limiting)
- **internal/wire.go**: Dependency injection setup

### Key Patterns
- **Clean Architecture**: Separation of handlers -> services -> repositories
- **Wire DI**: Dependency injection using Google Wire
- **Domain-Driven**: Code organized by business domains (user, problem, submission, forum, news)
- **Middleware Pipeline**: Auth, logging, rate limiting applied via Gin middleware

### Database Models
Core entities: User, Problem, TestCase, Submission, ForumPost, ForumReply, News

### Judge Service Communication
- Backend queues judge tasks in Redis
- Judge service processes tasks asynchronously  
- Results published back via Redis pubsub
- Backend updates submission status

## Environment Configuration

Development and production configs use `.env.dev` and `.env.prod` files. Key settings:
- Database connection (PostgreSQL)
- Redis connection  
- JWT secrets
- Server ports and modes
- Queue configuration

## Testing Strategy

- **Go tests**: Service layer unit tests in `backend/tests/`
- **Python API tests**: End-to-end API testing in `tests/`
- **Coverage**: Use `make test-coverage` for coverage reports

## API Documentation

OpenAPI specs available:
- Main API: `docs/api.yaml`
- Admin API: `docs/openapi/admin.yaml`
- Live docs: http://localhost:8080/docs (when running)

## Development Workflow

### Backend Development
1. Make changes to code
2. Run `make wire-gen` if modifying dependency injection
3. Run `make check` to verify code quality
4. Test with `make test` or API tests
5. Use `./scripts/deploy.sh --dev` for integration testing

### Frontend Development
1. Make changes to code (in `frontend/` or `admin-frontend/`)
2. Run `npm run lint` to check code style
3. Run `npm run type-check` for TypeScript validation
4. Test locally with `npm run dev`
5. Build with `npm run build` to verify production readiness

### Full System Testing
1. Use `./scripts/deploy.sh --dev` for end-to-end testing
2. Run Python API tests with `uv run python test_all.py`
3. Verify all services are working via `./scripts/deploy.sh --status`