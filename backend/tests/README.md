# 后端单元测试

本目录包含后端服务的单元测试代码。

## 目录结构

```
tests/
├── README.md           # 测试说明文档
├── services/           # 服务层单元测试
│   └── user_test.go    # 用户服务测试
└── ...
```

## 运行测试

### 运行所有测试

```bash
# 在backend根目录下运行
go test ./tests/... -v
```

### 运行特定模块测试

```bash
# 运行服务层测试
go test ./tests/services/... -v

# 运行用户服务测试
go test ./tests/services/ -run TestUserService -v
```

### 生成测试覆盖率报告

```bash
# 生成覆盖率报告
go test ./tests/... -coverprofile=coverage.out

# 查看覆盖率详情
go tool cover -html=coverage.out
```

## 测试规范

### 1. 测试文件命名

- 测试文件以 `_test.go` 结尾
- 测试文件名应与被测试的文件名对应，如 `user.go` 对应 `user_test.go`

### 2. 测试函数命名

- 测试函数以 `Test` 开头
- 函数名格式：`Test<StructName>_<MethodName>`
- 例如：`TestUserService_CreateUser`

### 3. 测试用例结构

使用表驱动测试（Table-Driven Tests）：

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name          string
        input         interface{}
        mockSetup     func(*MockRepository)
        expectedError string
        expectedResult interface{}
    }{
        {
            name: "成功案例",
            // ...
        },
        {
            name: "失败案例",
            // ...
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

### 4. Mock 使用

- 使用 `github.com/stretchr/testify/mock` 进行 Mock
- 为每个依赖接口创建对应的 Mock 实现
- 在测试中设置 Mock 期望和返回值

### 5. 断言

- 使用 `github.com/stretchr/testify/assert` 进行断言
- 优先使用语义化的断言方法，如 `assert.Equal`、`assert.NoError` 等

## 已实现的测试

### UserService 测试

- ✅ `TestUserService_CreateUser` - 测试用户创建功能
  - 成功创建用户
  - 用户名已存在
  - 邮箱已存在
  - 数据库创建失败

- ✅ `TestUserService_GetUserByUsername` - 测试根据用户名获取用户
  - 成功获取用户
  - 用户不存在
  - 数据库错误

- ✅ `TestUserService_GetUserByID` - 测试根据ID获取用户
  - 成功获取用户
  - 用户不存在
  - 数据库错误

- ✅ `TestUserService_ValidatePassword` - 测试密码验证
  - 密码正确
  - 密码错误

- ✅ `TestUserService_UpdateUser` - 测试用户更新
  - 成功更新用户
  - 更新失败

- ✅ `TestUserService_UpdateUserStats` - 测试用户统计信息更新
  - 成功更新统计信息
  - 更新失败

### ProblemService 测试

- ✅ `TestProblemService_CreateProblem` - 测试题目创建功能
  - 成功创建题目
  - 题目标题已存在
  - 数据库创建失败

- ✅ `TestProblemService_GetProblem` - 测试获取题目详情
  - 成功获取题目
  - 题目不存在
  - 数据库错误

- ✅ `TestProblemService_ListProblems` - 测试获取题目列表
  - 成功获取列表
  - 分页查询
  - 按难度筛选

- ✅ `TestProblemService_UpdateProblem` - 测试题目更新
  - 成功更新题目
  - 题目不存在
  - 更新失败

- ✅ `TestProblemService_DeleteProblem` - 测试题目删除
  - 成功删除题目
  - 题目不存在
  - 删除失败

- ✅ `TestProblemService_AddTestCase` - 测试添加测试用例
  - 成功添加测试用例
  - 题目不存在
  - 添加失败

- ✅ `TestProblemService_UpdateProblemStats` - 测试更新题目统计
  - 成功更新统计信息
  - 更新失败

## 待实现的测试

- [ ] ProblemService 测试
- [ ] SubmissionService 测试
- [ ] ForumService 测试
- [ ] NewsService 测试
- [ ] Repository 层测试
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