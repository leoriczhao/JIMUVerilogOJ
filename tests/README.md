# Verilog OJ API 测试套件 (RBAC 版)

这是一个完全基于 **RBAC (基于角色的访问控制)** 的模块化 Python API 测试套件，用于全面测试 Verilog OJ 后端系统的各个功能模块及权限控制。

## 🎯 核心特性

### ✅ **RBAC 权限测试框架**
- 🔐 **统一用户池管理**: 自动创建和管理 student/teacher/admin 三种角色
- 🔄 **角色快速切换**: 一键切换不同角色进行测试
- 🛡️ **权限边界验证**: 全面测试各角色的权限限制
- 🧹 **自动资源清理**: 测试后自动清理创建的所有资源

### ✅ **完整的功能覆盖**
- 📊 **45 个测试用例，100% 通过率**
- 🔐 JWT 认证和授权测试
- 📊 数据验证和错误处理
- 🛡️ 安全性检查（401/403 状态码验证）
- 🔄 CRUD 操作完整测试

### ✅ **OpenAPI Schema 自动验证**
- 📝 基于 OpenAPI 规范自动验证 API 响应格式
- 🎯 确保 API 实现与文档规范一致
- 🔍 自动检测字段类型、格式、必需字段等
- 📊 详细的 schema 验证错误报告

### ✅ **美观的测试报告**
- 🌈 彩色输出和进度指示
- 📝 详细的请求/响应日志
- 📊 模块化测试结果统计
- ⏱️ 测试执行时间统计

## 🏗️ 项目结构

```
tests/
├── fixtures/                    # 测试固件目录（新增）
│   ├── __init__.py             # 包初始化
│   ├── permissions.py          # RBAC 权限映射（与 Go 后端同步）
│   └── users.py                # 用户池管理
├── base_test.py                # 测试基类（增强 RBAC 支持）
├── openapi_validator.py        # OpenAPI Schema 验证器
├── test_health.py              # 健康检查测试模块
├── test_user.py                # 用户管理测试模块（重构）
├── test_problem.py             # 题目管理测试模块（重构）
├── test_submission.py          # 提交管理测试模块（重构）
├── test_forum.py               # 论坛管理测试模块（重构）
├── test_news.py                # 新闻管理测试模块（重构）
├── test_all.py                 # 综合测试入口
├── pyproject.toml              # uv 项目配置
└── README.md                   # 本文档
```

## 🚀 快速开始

### 环境要求

- Python 3.8+
- uv (Python 包管理器)
- 运行中的 Verilog OJ 后端服务 (localhost:8080)
- **数据库中存在 admin 用户** (username: `admin`, password: `admin123`)

### 1. 安装依赖

```bash
cd tests
uv sync
```

### 2. 启动后端服务

```bash
# 在项目根目录
./scripts/deploy.sh --dev
```

### 3. 运行测试

#### 方法 1：运行所有模块测试
```bash
uv run python test_all.py
```

#### 方法 2：运行特定模块
```bash
# 用户管理测试（15 个测试）
uv run python test_user.py

# 论坛管理测试（13 个测试）
uv run python test_forum.py

# 提交管理测试（9 个测试）
uv run python test_submission.py

# 新闻管理测试（11 个测试）
uv run python test_news.py

# 题目管理测试（9 个测试）
uv run python test_problem.py
```

## 📋 测试模块详情

### 👤 用户管理模块 (`test_user.py`) - 13 个测试
- **测试覆盖**:
  - ✅ 公开接口：注册、登录、重复注册、无效登录、输入验证（7 个测试）
  - ✅ 学生角色：个人资料管理、密码修改（3 个测试）
  - ✅ 教师角色：个人资料管理（1 个测试）
  - ✅ 管理员角色：更新用户角色（1 个测试）
  - ✅ 权限边界：未授权访问、错误旧密码（1 个测试）

### 📚 题目管理模块 (`test_problem.py`) - 9 个测试
- **测试覆盖**:
  - ✅ 公开接口：查看题目列表（1 个测试）
  - ✅ 学生角色：查看题目、创建题目（应拒绝）（2 个测试）
  - ✅ 教师角色：创建、更新、删除自己的题目、添加测试用例（4 个测试）
  - ✅ 管理员角色：更新、删除任意题目（2 个测试）

### 💬 论坛管理模块 (`test_forum.py`) - 13 个测试
- **测试覆盖**:
  - ✅ 公开接口：查看帖子列表（1 个测试）
  - ✅ 学生角色：创建帖子、查看、回复、删除他人帖子（应拒绝）（5 个测试）
  - ✅ 教师角色：创建、更新、回复帖子（3 个测试）
  - ✅ 管理员角色：更新任意帖子、删除帖子（2 个测试）
  - ✅ 权限边界：未授权操作（2 个测试）

### 📝 提交管理模块 (`test_submission.py`) - 9 个测试
- **测试覆盖**:
  - ✅ 准备环境：创建测试题目（自动）
  - ✅ 公开接口：查看提交列表（1 个测试）
  - ✅ 学生角色：创建提交、查看详情、删除（应拒绝）（4 个测试）
  - ✅ 教师角色：创建提交、查看统计（2 个测试）
  - ✅ 管理员角色：删除提交（1 个测试）
  - ✅ 权限边界：未授权提交（1 个测试）

### 📰 新闻管理模块 (`test_news.py`) - 11 个测试
- **测试覆盖**:
  - ✅ 公开接口：查看新闻列表（1 个测试）
  - ✅ 学生角色：查看新闻、创建新闻（应拒绝）（2 个测试）
  - ✅ 教师角色：创建、查看、更新、删除自己的新闻（4 个测试）
  - ✅ 管理员角色：创建、更新任意新闻、删除（3 个测试）
  - ✅ 权限边界：未授权操作（1 个测试）

## 🎯 RBAC 测试框架详解

### 用户池管理 (`fixtures/users.py`)

测试框架自动管理三种角色的用户：

```python
from base_test import BaseAPITester

class MyTester(BaseAPITester):
    def __init__(self):
        super().__init__()
        # 自动初始化用户池（student/teacher/admin）
        self.setup_user_pool()
```

**用户池工作流程**：
1. 使用预置的 `admin/admin123` 账户登录
2. 创建 `test_teacher_xxxx` 用户并提升为 teacher 角色
3. 创建 `test_student_xxxx` 用户（默认 student 角色）
4. 所有用户登录并保存 token 供后续使用

### 角色切换

```python
# 切换到学生角色
self.login_as('student')

# 切换到教师角色
self.login_as('teacher')

# 切换到管理员角色
self.login_as('admin')

# 获取当前角色
current_role = self.get_current_role()

# 检查当前用户权限
has_perm = self.check_permission('problem.create')
```

### 权限映射 (`fixtures/permissions.py`)

权限定义与 Go 后端 RBAC 系统完全同步：

```python
ROLE_PERMISSIONS = {
    'student': [
        'user.profile.read',
        'user.profile.update',
        'problem.read',
        'submission.create',
        'forum.post.create',
        # ... 共 15 个权限
    ],
    'teacher': [
        # 继承所有 student 权限 +
        'problem.create',
        'problem.update.own',
        'testcase.create',
        'news.create',
        # ... 共 25 个权限
    ],
    'admin': [
        # 继承所有 teacher 权限 +
        'user.update',
        'problem.update.all',
        'problem.delete.all',
        'manage.system',
        # ... 共 60+ 个权限
    ]
}
```

**辅助函数**：
```python
from fixtures.permissions import has_permission, get_minimum_role

# 检查角色权限
has_permission('teacher', 'problem.create')  # True
has_permission('student', 'problem.create')  # False

# 获取拥有某权限的最低角色
get_minimum_role('problem.create')  # 'teacher'
```

### 资源清理

测试框架会自动清理创建的所有资源：

```python
class MyTester(BaseAPITester):
    def run_tests(self):
        test_results = []

        try:
            # 创建资源并标记清理
            response = self.make_request("POST", "/problems", data=problem_data)
            if response:
                problem_id = response['problem']['id']
                self.mark_for_cleanup('problem', problem_id)

            # 运行测试...

        finally:
            # 自动清理所有标记的资源
            self.cleanup()
```

**支持的资源类型**：
- `problem` - 题目
- `submission` - 提交
- `post` - 论坛帖子
- `news` - 新闻

### 权限断言

框架提供专用的权限断言方法：

```python
# 断言返回 401 Unauthorized
response = self.make_request("POST", "/problems", expect_status=401)
self.assert_unauthorized(response)

# 断言返回 403 Forbidden
response = self.make_request("DELETE", "/problems/1", expect_status=403)
self.assert_forbidden(response)

# 断言当前用户有权限
self.assert_has_permission('problem.create')
```

## 📝 编写测试用例最佳实践

### 1. 标准测试流程

```python
#!/usr/bin/env python3
from base_test import BaseAPITester
from colorama import Back

class ModuleTester(BaseAPITester):
    """模块测试类（基于 RBAC）"""

    def __init__(self):
        super().__init__()
        # 初始化用户池
        self.setup_user_pool()
        # 存储资源ID
        self.resource_id = None

    def run_tests(self):
        """主测试流程"""
        test_results = []

        try:
            # 1. 公开接口测试
            self.print_section_header("公开接口测试", Back.CYAN)
            test_results.append(("获取列表(公开)", self.test_list_public()))

            # 2. 学生角色测试
            self.print_section_header("学生角色测试", Back.BLUE)
            self.login_as('student')
            test_results.append(("学生-操作", self.test_student_action()))

            # 3. 教师角色测试
            self.print_section_header("教师角色测试", Back.GREEN)
            self.login_as('teacher')
            test_results.append(("教师-操作", self.test_teacher_action()))

            # 4. 管理员角色测试
            self.print_section_header("管理员角色测试", Back.MAGENTA)
            self.login_as('admin')
            test_results.append(("管理员-操作", self.test_admin_action()))

            # 5. 权限边界测试
            self.print_section_header("权限边界测试", Back.RED)
            test_results.append(("未登录-操作(应拒绝)", self.test_unauthorized()))

        finally:
            # 6. 清理测试数据
            self.cleanup()

        # 7. 打印测试报告
        return self.print_test_summary(test_results)

    def test_list_public(self):
        """测试公开列表（无需登录）"""
        self.clear_token()
        response = self.make_request("GET", "/resources", expect_status=200)
        return response is not None

    def test_student_action(self):
        """学生角色测试"""
        response = self.make_request(
            "POST", "/resources",
            data={"name": "test"},
            expect_status=201
        )
        if response and "id" in response:
            self.mark_for_cleanup('resource', response['id'])
            return True
        return False

    def test_unauthorized(self):
        """未授权测试"""
        self.clear_token()
        response = self.make_request(
            "POST", "/resources",
            data={"name": "test"},
            expect_status=401
        )
        return self.assert_unauthorized(response)

def main():
    """主函数"""
    print("\n" + "=" * 60)
    print(" 模块测试（重构版 - 基于 RBAC）")
    print("=" * 60 + "\n")

    tester = ModuleTester()
    success = tester.run_tests()

    return 0 if success else 1

if __name__ == "__main__":
    import sys
    sys.exit(main())
```

### 2. 权限测试模式

```python
# 模式 1: 期望成功
def test_allowed_action(self):
    """角色有权限，应该成功"""
    response = self.make_request("POST", "/endpoint", expect_status=201)
    return response is not None

# 模式 2: 期望被拒绝（403 Forbidden）
def test_forbidden_action(self):
    """角色无权限，应该返回 403"""
    response = self.make_request("DELETE", "/endpoint/1", expect_status=403)
    return self.assert_forbidden(response)

# 模式 3: 期望未认证（401 Unauthorized）
def test_unauthorized_action(self):
    """未登录，应该返回 401"""
    self.clear_token()
    response = self.make_request("POST", "/endpoint", expect_status=401)
    return self.assert_unauthorized(response)
```

### 3. 资源依赖处理

```python
def test_sequence(self):
    """测试序列：创建 -> 使用 -> 清理"""
    # 第 1 步：创建依赖资源
    response = self.make_request("POST", "/problems", data=problem_data)
    if not response:
        return False

    problem_id = response['problem']['id']
    self.mark_for_cleanup('problem', problem_id)

    # 第 2 步：使用资源
    submission_data = {"problem_id": problem_id, "code": "..."}
    response = self.make_request("POST", "/submissions", data=submission_data)
    if not response:
        return False

    submission_id = response['submission']['id']
    self.mark_for_cleanup('submission', submission_id)

    # cleanup() 会自动倒序删除（submission 先删，problem 后删）
    return True
```

## 📊 示例输出

```
============================================================
 用户管理测试模块（重构版 - 基于 RBAC）
============================================================

✅ 用户池初始化完成

 公开接口测试

✅ 请求成功 - POST /users/register
✅ 用户注册成功，ID: 21
✅ 请求成功 - POST /users/login
✅ 登录成功，Token: eyJhbGciOiJIUzI1NiI...
✅ 请求成功 - POST /users/register
✅ 重复注册正确被拒绝

 学生角色测试

ℹ️  已切换到角色: student (用户ID: 20)
✅ 请求成功 - GET /users/profile
✅ 成功获取资料: test_student_5658
✅ 请求成功 - PUT /users/profile
✅ 请求成功 - PUT /users/password

 教师角色测试

ℹ️  已切换到角色: teacher (用户ID: 19)
✅ 请求成功 - GET /users/profile
✅ 成功获取资料: test_teacher_5657

 管理员角色测试

ℹ️  已切换到角色: admin (用户ID: 3)
✅ 请求成功 - PUT /admin/users/20/role
✅ 成功更新并回滚用户角色

 权限边界测试

✅ 请求成功 - GET /users/profile
✅ ✓ 认证检查通过：正确返回 401 Unauthorized

⚠️  开始清理测试数据 (0 个资源)...
✅ 用户池清理完成
✅ 测试数据清理完成

 📊 测试结果总结 📊
============================================================
新用户注册                ✅ 通过
新用户登录                ✅ 通过
重复注册(应拒绝)            ✅ 通过
无效登录(应拒绝)            ✅ 通过
无效邮箱格式(应拒绝)          ✅ 通过
用户名过短(应拒绝)           ✅ 通过
密码过短(应拒绝)            ✅ 通过
学生-查看个人资料            ✅ 通过
学生-更新个人资料            ✅ 通过
学生-修改密码              ✅ 通过
教师-查看个人资料            ✅ 通过
教师-更新个人资料            ✅ 通过
管理员-更新用户角色           ✅ 通过
未登录-访问资料(应拒绝)        ✅ 通过
错误旧密码(应拒绝)           ✅ 通过
============================================================
总测试数: 15
通过数: 15
失败数: 0
通过率: 100.0%

🎉 所有测试通过！
```

## 🛠️ 自定义配置

### API 地址配置

在 `base_test.py` 中修改基础 URL：
```python
BASE_URL = "http://localhost:8080"  # 修改为你的 API 地址
API_BASE = f"{BASE_URL}/api/v1"
```

### 用户池配置

在 `fixtures/users.py` 中修改用户池配置：
```python
# 修改 admin 账户信息
admin_login = {
    "username": "admin",      # 修改管理员用户名
    "password": "admin123"    # 修改管理员密码
}

# 修改测试用户密码
teacher_data = {
    "password": "test123456"  # 修改测试用户密码
}
```

## 🔬 OpenAPI Schema 验证详解

### 启用 Schema 验证

```python
# 方法 1: 全局启用（默认）
class UserTester(BaseAPITester):
    def __init__(self):
        super().__init__(enable_schema_validation=True)

# 方法 2: 单个请求验证
response = self.make_request(
    "POST", "/users/register",
    data=user_data,
    expect_status=201,
    module="user"  # 指定模块名启用验证
)

# 方法 3: 禁用特定请求验证
response = self.make_request(
    "POST", "/admin/users/1/role",
    data={"role": "teacher"},
    expect_status=200,
    module="admin",
    validate_schema=False  # 禁用验证
)
```

### Schema 验证错误示例

```
⚠️  Schema验证错误 (1) ⚠️

错误 #1:
  模块: user
  请求: POST /users/register
  状态码: 201
  详情: Schema验证失败:
  - user.email: 'invalid-email' is not a 'email'
  - user: 'role' is a required property
```

## 🐛 故障排除

### 常见问题

**1. "无法登录管理员账户"**
```
Exception: 无法登录管理员账户，请确保数据库中存在 admin 用户
```
**解决方案**: 在数据库中手动创建 admin 用户
```sql
INSERT INTO users (username, password_hash, email, role, created_at, updated_at)
VALUES ('admin', '$2a$10$...', 'admin@verilogoj.com', 'admin', NOW(), NOW());
```

**2. "连接被拒绝"**
```
Connection refused
```
**解决方案**: 确保后端服务运行在 localhost:8080
```bash
./scripts/deploy.sh --dev
./scripts/deploy.sh --status
```

**3. "测试用户创建失败"**
```
无法创建 teacher 测试用户
```
**解决方案**: 检查 admin API 是否正常工作
```bash
# 测试 admin API
curl -X PUT http://localhost:8080/api/v1/admin/users/1/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "teacher"}'
```

## 📈 测试统计

### 测试覆盖率

| 模块 | 测试数量 | 通过率 | 覆盖的角色 |
|------|----------|--------|------------|
| test_user.py | 13 | 100% | student, teacher, admin |
| test_forum.py | 13 | 100% | student, teacher, admin |
| test_submission.py | 9 | 100% | student, teacher, admin |
| test_news.py | 11 | 100% | student, teacher, admin |
| test_problem.py | 9 | 100% | student, teacher, admin |
| **总计** | **45** | **100%** | **完整 RBAC 覆盖** |

### 权限测试覆盖

- ✅ 公开接口测试（无需登录）
- ✅ Student 角色权限测试
- ✅ Teacher 角色权限测试
- ✅ Admin 角色权限测试
- ✅ 401 Unauthorized 边界测试
- ✅ 403 Forbidden 边界测试

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进测试套件！

### 贡献流程
1. Fork 项目
2. 创建功能分支
3. 遵循 RBAC 测试框架编写测试
4. 确保所有测试通过
5. 提交 Pull Request

### 代码规范
- 所有测试继承自 `BaseAPITester`
- 使用 `setup_user_pool()` 初始化用户
- 使用 `login_as()` 切换角色
- 使用 `mark_for_cleanup()` 标记资源
- 使用 `assert_forbidden()` 和 `assert_unauthorized()` 验证权限

## 📄 许可证

MIT License

---

**注意**: 本测试套件与 Go 后端的 RBAC 权限系统完全同步，权限定义位于 `fixtures/permissions.py`，源自 `backend/internal/middleware/rbac_permissions.go`。
