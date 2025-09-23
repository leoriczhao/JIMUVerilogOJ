# Makefile for Verilog OJ Backend

.PHONY: test test-verbose test-coverage test-services test-user test-problem test-forum test-news \
        clean build run deps fmt vet check help \
        docker-build docker-run docker-stop docker-clean \
        dev-setup dev-start dev-stop \
        lint security-check mod-tidy wire-gen

# 默认目标
all: check

# 运行所有测试
test:
	@echo "Running all tests..."
	go test ./tests/... -short

# 运行测试并显示详细输出
test-verbose:
	@echo "Running all tests with verbose output..."
	go test ./tests/... -v

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	go test ./tests/... -coverprofile=coverage.out
	@echo "Coverage report generated: coverage.out"
	@echo "To view HTML coverage report, run: make coverage-html"

# 生成HTML覆盖率报告
coverage-html: coverage.out
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

# 运行服务层测试
test-services:
	@echo "Running services tests..."
	go test ./tests/services/... -v

# 运行用户服务测试
test-user:
	@echo "Running user service tests..."
	go test ./tests/services/ -run TestUserService -v

# 运行题目服务测试
test-problem:
	@echo "Running problem service tests..."
	go test ./tests/services/ -run TestProblemService -v

# 运行论坛服务测试
test-forum:
	@echo "Running forum service tests..."
	go test ./tests/services/ -run TestForumService -v

# 运行新闻服务测试
test-news:
	@echo "Running news service tests..."
	go test ./tests/services/ -run TestNewsService -v

# 清理生成的文件
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html
	rm -f main

# 构建项目
build:
	@echo "Building project..."
	go build -o main cmd/main.go

# 运行项目
run: build
	@echo "Running project..."
	./main

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 代码检查
vet:
	@echo "Running go vet..."
	go vet ./...

# 运行所有检查（格式化、vet、测试）
check: fmt vet lint test
	@echo "All checks passed!"

# 代码质量检查
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# 安全检查
security-check:
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed, skipping..."; \
	fi

# 整理模块依赖
mod-tidy:
	@echo "Tidying module dependencies..."
	go mod tidy
	go mod verify

# 生成Wire依赖注入代码
wire-gen:
	@echo "Generating Wire dependency injection code..."
	@if command -v wire >/dev/null 2>&1; then \
		wire ./internal/...; \
	else \
		echo "wire not installed, run: go install github.com/google/wire/cmd/wire@latest"; \
	fi

# Docker相关操作
docker-build:
	@echo "Building Docker image..."
	docker build -f ../docker/backend.Dockerfile -t verilog-oj-backend .

docker-run:
	@echo "Running Docker container..."
	docker run -d --name verilog-oj-backend -p 8080:8080 verilog-oj-backend

docker-stop:
	@echo "Stopping Docker container..."
	docker stop verilog-oj-backend || true
	docker rm verilog-oj-backend || true

docker-clean:
	@echo "Cleaning Docker images..."
	docker rmi verilog-oj-backend || true

# 开发环境管理
dev-setup:
	@echo "Setting up development environment..."
	@echo "Installing development tools..."
	@go install github.com/google/wire/cmd/wire@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development environment setup complete!"

dev-start:
	@echo "Starting development environment..."
	docker-compose -f ../docker-compose.dev.yml up -d
	@echo "Development environment started!"

dev-stop:
	@echo "Stopping development environment..."
	docker-compose -f ../docker-compose.dev.yml down
	@echo "Development environment stopped!"

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run all tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  coverage-html  - Generate HTML coverage report"
	@echo "  test-services  - Run services layer tests"
	@echo "  test-user      - Run user service tests"
	@echo "  test-problem   - Run problem service tests"
	@echo "  test-forum     - Run forum service tests"
	@echo "  test-news      - Run news service tests"
	@echo ""
	@echo "Development:"
	@echo "  build          - Build the project"
	@echo "  run            - Build and run the project"
	@echo "  clean          - Clean generated files"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  security-check - Run security checks"
	@echo "  check          - Run all checks (fmt, vet, lint, test)"
	@echo "  mod-tidy       - Tidy module dependencies"
	@echo "  wire-gen       - Generate Wire dependency injection code"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-stop    - Stop Docker container"
	@echo "  docker-clean   - Clean Docker images"
	@echo ""
	@echo "Environment:"
	@echo "  dev-setup      - Setup development environment"
	@echo "  dev-start      - Start development environment"
	@echo "  dev-stop       - Stop development environment"
	@echo ""
	@echo "  help           - Show this help message"