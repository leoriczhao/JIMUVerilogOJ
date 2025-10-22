# JIMUVerilogOJ - Verilog Online Judge System

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![CI Status](https://github.com/leoriczhao/JIMUVerilogOJ/workflows/Backend%20CI/badge.svg)](https://github.com/leoriczhao/JIMUVerilogOJ/actions)

A modern online judge platform specifically designed for Verilog HDL, featuring microservices architecture and efficient code compilation, testing, and evaluation services.

[Features](#features) • [Quick Start](#quick-start) • [Development](#development) • [API Documentation](#api-documentation) • [Deployment](#deployment)

**[English](README.md)** | [简体中文](README_zh.md)

</div>

---

## 📑 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [System Architecture](#system-architecture)
- [Tech Stack](#tech-stack)
- [Quick Start](#quick-start)
- [Development Guide](#development-guide)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Overview

**JIMUVerilogOJ** is a fully-featured online judge system for Verilog Hardware Description Language, designed to provide a convenient code verification platform for hardware design learning and teaching. The system adopts a microservices architecture with front-end and back-end separation, deploying the judge engine independently to ensure high availability and scalability.

### Core Advantages

- 🎯 **Specialized** - Specifically designed for Verilog HDL, supporting complete compilation and simulation workflows
- 🚀 **High Performance** - Asynchronous queue processing with independently deployed judge engine, supporting high concurrency
- 🛡️ **Secure & Reliable** - Role-Based Access Control (RBAC) with comprehensive security mechanisms
- 📊 **Easily Extensible** - Microservices architecture with modular design for easy feature expansion
- 📖 **Developer Friendly** - Complete OpenAPI documentation with well-structured codebase

## Features

### 🔐 User System
- User registration, login, and authentication
- JWT Token authentication mechanism
- Role-Based Access Control (Admin/User)
- User profile management and personalization

### 📚 Problem Management
- Problem creation, editing, and categorization
- Test case CRUD operations
- Difficulty levels and tagging system
- Problem search and filtering

### ⚖️ Judge Engine
- Verilog code compilation checking
- Automated test case execution
- Waveform comparison and result verification
- Asynchronous queue processing for judge tasks
- Detailed error feedback

### 💬 Community Forum
- Discussion posts and replies
- Like and interaction features
- Categorized discussion areas
- Content management and moderation

### 📰 News & Announcements
- System announcements
- News management
- Categories and tagging

### 📊 Statistics & Analytics
- User submission statistics
- Problem pass rate analysis
- System usage monitoring

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│  ┌──────────────────┐              ┌──────────────────┐        │
│  │   User Frontend  │              │  Admin Frontend  │        │
│  │     (Vue 3)      │              │     (Vue 3)      │        │
│  └──────────────────┘              └──────────────────┘        │
└────────────────┬──────────────────────────┬────────────────────┘
                 │                          │
                 └──────────┬───────────────┘
                            │ HTTP/REST API
                 ┌──────────▼───────────┐
                 │   Nginx (Reverse     │
                 │      Proxy)          │
                 └──────────┬───────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
┌────────▼────────┐  ┌──────▼──────┐  ┌───────▼────────┐
│  Backend API    │  │   Judge     │  │   Redis        │
│    Service      │◄─┤   Service   │  │  (Cache/Queue) │
│     (Go)        │  │    (Go)     │  └────────────────┘
└────────┬────────┘  └──────┬──────┘
         │                  │
         │                  │
    ┌────▼─────┐     ┌──────▼───────┐
    │PostgreSQL│     │Verilog Tools │
    │(Database)│     │  (iverilog)  │
    └──────────┘     └──────────────┘
```

### Service Components

- **Frontend**: User interface for problem browsing and code submission
- **Admin Frontend**: Administrative dashboard for system management and content moderation
- **Backend Service**: Core business logic handling API requests
- **Judge Service**: Independent judging service for code compilation and testing
- **PostgreSQL**: Primary database storing users, problems, and other data
- **Redis**: Cache and message queue supporting asynchronous judging

## Tech Stack

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.21+ | Backend programming language |
| Gin | Latest | HTTP routing framework |
| GORM | Latest | ORM framework |
| PostgreSQL | 15+ | Relational database |
| Redis | 7+ | Cache and message queue |
| Wire | Latest | Dependency injection |
| JWT | - | Authentication |

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| Vue | 3.x | Frontend framework |
| TypeScript | Latest | Type system |
| Vite | Latest | Build tool |
| Element Plus | Latest | UI component library |
| Monaco Editor | Latest | Code editor |

### Judge Environment

| Technology | Purpose |
|------------|---------|
| Icarus Verilog (iverilog) | Verilog compiler |
| GTKWave | Waveform viewer |
| Docker | Isolated judging environment |

### DevOps

| Technology | Purpose |
|------------|---------|
| Docker | Containerization |
| Docker Compose | Service orchestration |
| Nginx | Reverse proxy |
| GitHub Actions | CI/CD |
| golangci-lint | Code quality checking |

## Quick Start

### Prerequisites

- **Docker** 20.0+
- **Docker Compose** 2.0+
- **Go** 1.21+ (for local development)
- **Node.js** 18+ (for frontend development)

### One-Click Deployment

```bash
# Clone the repository
git clone https://github.com/leoriczhao/JIMUVerilogOJ.git
cd JIMUVerilogOJ

# Deploy development environment
./scripts/deploy.sh --dev

# Or deploy production environment
./scripts/deploy.sh --prod
```

### Manual Deployment

```bash
# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f backend
```

### Access URLs

After successful deployment:

- **Frontend**: http://localhost:80
- **Backend API**: http://localhost:8080
- **API Docs**: http://localhost:8080/swagger/index.html
- **Admin Panel**: http://localhost:3000

Default admin credentials:
- Username: `admin`
- Password: `admin123`

## Development Guide

### Backend Development

Navigate to the backend directory:

```bash
cd backend/

# Install dependencies
make deps

# Generate dependency injection code (required after modifying wire.go)
make wire-gen

# Run service
make run

# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Generate test coverage
make test-coverage

# Run all checks
make check
```

### Frontend Development

```bash
cd frontend/

# Install dependencies
npm install

# Start development server
npm run dev

# Lint code
npm run lint

# Type checking
npm run type-check

# Production build
npm run build
```

### Admin Frontend Development

```bash
cd admin-frontend/

# Same workflow as frontend
npm install
npm run dev
```

### API Testing

Using Python test suite:

```bash
cd tests/

# Using uv (recommended)
uv run python test_all.py

# Or using pip
pip install -r requirements.txt
python test_all.py
```

### Project Structure

```
JIMUVerilogOJ/
├── backend/                    # Backend service
│   ├── cmd/                    # Application entry point
│   │   └── main.go
│   ├── internal/               # Internal modules
│   │   ├── config/            # Configuration management
│   │   ├── models/            # Data models
│   │   ├── handlers/          # HTTP handlers
│   │   ├── services/          # Business logic layer
│   │   ├── repository/        # Data access layer
│   │   ├── middleware/        # Middleware
│   │   └── wire.go           # Dependency injection config
│   ├── Makefile              # Build scripts
│   └── go.mod                # Go module dependencies
│
├── judge-service/             # Judge service
│   ├── cmd/
│   ├── internal/
│   │   ├── judge/            # Judge logic
│   │   └── queue/            # Queue processing
│   └── go.mod
│
├── frontend/                  # User frontend
│   ├── src/
│   │   ├── components/       # Vue components
│   │   ├── views/           # Page views
│   │   ├── router/          # Router config
│   │   └── stores/          # State management
│   └── package.json
│
├── admin-frontend/           # Admin dashboard
│   └── ...
│
├── tests/                    # API tests
│   ├── test_user.py
│   ├── test_problem.py
│   └── test_submission.py
│
├── docs/                     # Documentation
│   └── openapi/             # OpenAPI specifications
│       ├── user.yaml
│       ├── admin.yaml
│       ├── problem.yaml
│       ├── news.yaml
│       └── submission.yaml
│
├── scripts/                  # Deployment scripts
│   └── deploy.sh
│
├── docker/                   # Docker configurations
│   ├── backend.Dockerfile
│   └── judge.Dockerfile
│
├── .github/                  # GitHub configurations
│   └── workflows/           # CI/CD workflows
│
├── docker-compose.yml        # Base service orchestration
├── docker-compose.dev.yml    # Development environment
├── docker-compose.prod.yml   # Production environment
├── CLAUDE.md                 # Claude Code project guide
└── README.md                 # This file
```

## API Documentation

### OpenAPI Specifications

The project uses OpenAPI 3.0 specifications, organized by modules:

- **User API**: [docs/openapi/user.yaml](docs/openapi/user.yaml)
- **Admin API**: [docs/openapi/admin.yaml](docs/openapi/admin.yaml)
- **Problem API**: [docs/openapi/problem.yaml](docs/openapi/problem.yaml)
- **Submission API**: [docs/openapi/submission.yaml](docs/openapi/submission.yaml)
- **News API**: [docs/openapi/news.yaml](docs/openapi/news.yaml)

### Online Documentation

Access Swagger UI after starting the service:
```
http://localhost:8080/swagger/index.html
```

### Common API Examples

#### User Registration
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### User Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### Submit Code
```bash
curl -X POST http://localhost:8080/api/v1/submissions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "problem_id": 1,
    "code": "module test; ... endmodule",
    "language": "verilog"
  }'
```

## Testing

### Backend Tests

```bash
cd backend/

# Run all tests
make test

# Run specific service tests
make test-user
make test-services

# View test coverage
make test-coverage

# View verbose output
make test-verbose
```

### Integration Tests

```bash
cd tests/

# Run all API tests
uv run python test_all.py

# Run specific tests
uv run python test_user.py
```

### Code Quality Checks

```bash
cd backend/

# Format code
make fmt

# Run linter
make lint

# Run vet
make vet

# Run all checks
make check
```

## Deployment

### Development Environment

```bash
# Using deployment script
./scripts/deploy.sh --dev

# Or manually
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### Production Environment

1. **Configure Environment Variables**

Create `.env.prod` file:

```bash
# Database configuration
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_secure_password
DB_DATABASE=verilog_oj

# Redis configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT configuration
JWT_SECRET=your_jwt_secret_key_at_least_32_chars

# Server configuration
GIN_MODE=release
SERVER_PORT=8080
```

2. **Deploy Services**

```bash
./scripts/deploy.sh --prod
```

3. **Configure Nginx and SSL**

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        proxy_pass http://frontend:3000;
    }
}
```

### Operations Commands

```bash
# Check service status
./scripts/deploy.sh --status
docker-compose ps

# View logs
./scripts/deploy.sh --logs
docker-compose logs -f backend

# Restart services
./scripts/deploy.sh --restart
docker-compose restart backend

# Stop services
./scripts/deploy.sh --stop
docker-compose down

# Backup database
docker-compose exec postgres pg_dump -U postgres verilog_oj > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres verilog_oj < backup.sql
```

## Contributing

We welcome and appreciate all forms of contributions!

### Contribution Workflow

1. **Fork the project** to your GitHub account
2. **Clone the project** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/JIMUVerilogOJ.git
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **Commit your changes**:
   ```bash
   git commit -m "feat: add amazing feature"
   ```
5. **Push the branch**:
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Create a Pull Request**

### Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation updates
- `style:` Code formatting
- `refactor:` Code refactoring
- `test:` Test-related
- `chore:` Build/tooling

### Code Standards

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go) and golangci-lint rules
- **Vue/TypeScript**: Follow ESLint and Prettier configurations
- **Pre-commit**: Ensure `make check` (backend) and `npm run lint` (frontend) pass

## License

This project is licensed under the [Apache License 2.0](LICENSE).

## Contact

- **Project Homepage**: https://github.com/leoriczhao/JIMUVerilogOJ
- **Issue Tracker**: https://github.com/leoriczhao/JIMUVerilogOJ/issues
- **Discussions**: https://github.com/leoriczhao/JIMUVerilogOJ/discussions

## Acknowledgments

Thanks to the following open-source projects and tools:

- [Go](https://golang.org/) - Powerful backend programming language
- [Gin](https://gin-gonic.com/) - High-performance HTTP framework
- [GORM](https://gorm.io/) - Elegant ORM framework
- [Vue.js](https://vuejs.org/) - Progressive frontend framework
- [PostgreSQL](https://www.postgresql.org/) - Reliable relational database
- [Redis](https://redis.io/) - High-performance cache and message queue
- [Docker](https://www.docker.com/) - Containerization platform
- [Icarus Verilog](http://iverilog.icarus.com/) - Verilog compiler

## Star History

If this project helps you, please give us a ⭐️!

[![Star History Chart](https://api.star-history.com/svg?repos=leoriczhao/JIMUVerilogOJ&type=Date)](https://star-history.com/#leoriczhao/JIMUVerilogOJ&Date)

---

<div align="center">

**[⬆ Back to Top](#jimuverilogoj---verilog-online-judge-system)**

Made with ❤️ by [leoriczhao](https://github.com/leoriczhao)

</div>
