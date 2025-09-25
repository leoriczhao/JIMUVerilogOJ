# 后端单元测试

本目录包含后端服务的单元测试代码。

## 目录结构

```
tests/
├── README.md           # 测试说明文档
├── handlers/           # 控制器层单元测试
│   ├── admin_test.go
│   └── user_test.go
├── repository/         # 仓储层单元测试
│   ├── admin_test.go
│   ├── forum_test.go
│   ├── news_test.go
│   ├── problem_test.go
│   └── submission_test.go
└── services/           # 服务层单元测试
    └── ...
```

## 运行测试

### 运行所有测试

```bash
# 在backend根目录下运行
go test ./... -v
```

### 运行特定模块测试

```bash
# 运行服务层测试
go test ./tests/services/... -v

# 运行仓储层测试
go test ./tests/repository/... -v

# 运行控制器层测试
go test ./tests/handlers/... -v
```

### 生成测试覆盖率报告

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out

# 查看覆盖率详情
go tool cover -html=coverage.out
```

## 测试规范

- **位置**: 所有单元测试都应放在 `backend/tests` 目录下，并根据被测试的层（`services`, `repository`, `handlers`）进行组织。
- **包名**: 测试文件应使用 `_test` 后缀的包名，例如 `services_test`，以进行黑盒测试。 (*注意：根据用户反馈，部分现有测试可能不遵循此规则，但新测试应尽量遵循。*)
- **文件名**: 测试文件以 `_test.go` 结尾。
- **函数名**: 测试函数以 `Test` 开头，格式为 `Test<StructName>_<MethodName>`。

## 已实现的测试

### Handler 层测试
- ✅ `AdminHandler`
- ✅ `UserHandler`

### Repository 层测试
- ✅ `AdminRepository`
- ✅ `ForumRepository`
- ✅ `NewsRepository`
- ✅ `ProblemRepository`
- ✅ `SubmissionRepository`
- ✅ `UserRepository` (位于 `internal/repository` 目录下)

### Service 层测试
- ✅ `AdminService`
- ✅ `UserService`
- ✅ `ProblemService`
- ✅ `SubmissionService`
- ✅ `ForumService`
- ✅ `NewsService`


## 待实现的测试

- [x] Repository 层测试
- [x] Handler 层测试
- [ ] ... (其他待完善部分)

## 测试最佳实践

1. **独立性**：每个测试用例应该独立，不依赖其他测试的执行结果。
2. **可重复性**：测试应该可以重复执行，每次结果一致。
3. **快速执行**：单元测试应该快速执行，避免依赖外部资源。
4. **清晰的命名**：测试用例名称应该清楚地描述测试场景。
5. **充分的覆盖**：测试应该覆盖正常流程、边界条件和异常情况。
6. **使用 Mock**：对外部依赖使用 Mock，确保测试的隔离性。