# 后端单元测试

本目录包含后端服务的单元测试代码。

## 目录结构

```
tests/
├── README.md           # 测试说明文档
├── repository/         # 仓储层单元测试
│   ├── admin_test.go
│   ├── forum_test.go
│   ├── news_test.go
│   ├── problem_test.go
│   └── submission_test.go
└── services/           # 服务层单元测试
    ├── admin_test.go
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
```

### 生成测试覆盖率报告

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out

# 查看覆盖率详情
go tool cover -html=coverage.out
```

## 测试规范

### 1. 测试文件命名

- 测试文件以 `_test.go` 结尾
- 测试文件名应与被测试的文件名对应，如 `user.go` 对应 `user_test.go` (在 co-location 模式下)

### 2. 测试函数命名

- 测试函数以 `Test` 开头
- 函数名格式：`Test<StructName>_<MethodName>`
- 例如：`TestUserService_CreateUser`

### 3. 测试用例结构

使用表驱动测试（Table-Driven Tests）或场景驱动的独立函数。

### 4. Mock 使用

- 使用 `github.com/stretchr/testify/mock` 进行 Mock
- 为每个依赖接口创建对应的 Mock 实现
- 在测试中设置 Mock 期望和返回值

### 5. 断言

- 使用 `github.com/stretchr/testify/assert` 进行断言
- 优先使用语义化的断言方法，如 `assert.Equal`、`assert.NoError` 等

## 已实现的测试

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
- [ ] Handler 层测试

## 测试最佳实践

1. **独立性**：每个测试用例应该独立，不依赖其他测试的执行结果
2. **可重复性**：测试应该可以重复执行，每次结果一致
3. **快速执行**：单元测试应该快速执行，避免依赖外部资源
4. **清晰的命名**：测试用例名称应该清楚地描述测试场景
5. **充分的覆盖**：测试应该覆盖正常流程、边界条件和异常情况
6. **使用 Mock**：对外部依赖使用 Mock，确保测试的隔离性

## 持续集成

在 CI/CD 流程中，应该：

1. 在每次代码提交时运行所有单元测试
2. 确保测试覆盖率达到预期标准（建议 80% 以上）
3. 测试失败时阻止代码合并
4. 生成测试报告和覆盖率报告