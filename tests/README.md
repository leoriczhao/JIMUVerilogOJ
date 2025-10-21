# Verilog OJ API 测试套件

这是一个模块化的 Python API 测试套件，用于全面测试 Verilog OJ 后端系统的各个功能模块。

## 🏗️ 项目结构

```
tests/
├── base_test.py       # 测试基类，提供公共功能
├── test_health.py     # 健康检查测试模块
├── test_user.py       # 用户管理测试模块  
├── test_problem.py    # 题目管理测试模块
├── test_submission.py # 提交管理测试模块
├── test_forum.py      # 论坛管理测试模块
├── test_news.py       # 新闻管理测试模块
├── test_all.py        # 综合测试入口
├── run_tests.sh       # 测试运行脚本
├── pyproject.toml     # uv 项目配置
└── README.md          # 本文档
```

## 🚀 快速开始

### 环境要求

- Python 3.8+
- uv (Python 包管理器)
- 运行中的 Verilog OJ 后端服务 (localhost:8080)

### 1. 安装依赖

```bash
cd tests
uv sync
```

### 2. 启动后端服务

```bash
# 在项目根目录
./scripts/deploy.sh -d
```

### 3. 运行测试

#### 方法1：使用脚本（推荐）
```bash
# 运行所有模块测试
./run_tests.sh

# 或者直接使用 uv
uv run python test_all.py
```

#### 方法2：运行特定模块
```bash
# 只运行用户管理测试
uv run python test_user.py

# 只运行健康检查测试
uv run python test_health.py

# 运行指定的多个模块
uv run python test_all.py --modules user problem submission
```

#### 方法3：查看可用模块
```bash
uv run python test_all.py --list
```

## 📋 测试模块详情

### 🏥 健康检查模块 (`test_health.py`)
- **功能**: 验证 API 服务基本状态
- **测试项目**:
  - 健康检查接口 (`/health`)
  - API 根路径响应

### 👤 用户管理模块 (`test_user.py`)  
- **功能**: 测试用户注册、登录、信息管理
- **测试项目**:
  - 用户注册
  - 用户登录（JWT 认证）
  - 获取用户信息
  - 更新用户信息
  - 重复注册检查
  - 无效登录检查
  - 未授权访问保护

### 📚 题目管理模块 (`test_problem.py`)
- **功能**: 测试题目的增删改查操作
- **测试项目**:
  - 获取题目列表
  - 创建题目（需要认证）
  - 获取题目详情
  - 更新题目（需要认证）
  - 删除题目（需要认证）
  - 未授权操作保护

### 📝 提交管理模块 (`test_submission.py`)
- **功能**: 测试代码提交和判题功能
- **测试项目**:
  - 获取提交列表
  - 创建代码提交（需要认证）
  - 获取提交详情
  - 无效提交检查
  - 大代码提交测试
  - 未授权提交保护

### 💬 论坛管理模块 (`test_forum.py`)
- **功能**: 测试论坛帖子和回复功能
- **测试项目**:
  - 获取论坛帖子列表
  - 创建论坛帖子（需要认证）
  - 获取帖子详情
  - 更新论坛帖子（需要认证）
  - 获取帖子回复
  - 创建帖子回复（需要认证）
  - 删除论坛帖子（需要认证）
  - 未授权操作保护

### 📰 新闻管理模块 (`test_news.py`)
- **功能**: 测试新闻的增删改查功能
- **测试项目**:
  - 获取新闻列表
  - 创建新闻（需要认证）
  - 获取新闻详情
  - 更新新闻（需要认证）
  - 删除新闻（需要认证）
  - 创建草稿新闻
  - 未授权操作保护

## 🎯 测试特性

### ✅ **完整的功能覆盖**
- 🔐 JWT 认证和授权测试
- 📊 数据验证和错误处理
- 🛡️ 安全性检查（未授权访问）
- 🔄 CRUD 操作完整测试

### ✅ **OpenAPI Schema 自动验证** (新功能!)
- 📝 基于OpenAPI规范自动验证API响应格式
- 🎯 确保API实现与文档规范一致
- 🔍 自动检测字段类型、格式、必需字段等
- 📊 详细的schema验证错误报告
- 🚀 支持$ref引用的自动解析

### ✅ **美观的测试报告**
- 🌈 彩色输出和进度指示
- 📝 详细的请求/响应日志
- 📊 模块化测试结果统计
- ⏱️ 测试执行时间统计

### ✅ **灵活的测试控制**
- 🎯 可选择性运行特定模块
- 🔧 命令行参数支持
- 🚀 独立模块测试能力
- 📋 测试模块列表查看

## 📊 示例输出

```
🚀 Verilog OJ 完整API测试套件 🚀
API 基础地址: http://localhost:8080/api/v1
测试时间: 2025-01-23 16:48:00
================================================================================

🏥 健康检查测试模块
✅ 健康检查接口正常

👤 用户管理测试模块  
✅ 用户注册成功
✅ 用户登录成功
✅ 获取用户信息成功
...

📊 综合测试结果总结 📊
================================================================================
健康检查模块          ✅ 通过
用户管理模块          ✅ 通过
题目管理模块          ✅ 通过
提交管理模块          ✅ 通过
论坛管理模块          ✅ 通过
新闻管理模块          ✅ 通过
================================================================================
总模块数: 6
通过模块: 6
失败模块: 0
通过率: 100.0%
测试耗时: 12.45秒

🎉 所有模块测试通过！Verilog OJ API 完全正常！
🚀 系统已准备好投入使用！
```

## 🛠️ 自定义配置

### API 地址配置
在 `base_test.py` 中修改基础 URL：
```python
BASE_URL = "http://localhost:8080"  # 修改为你的 API 地址
```

### 测试数据配置
各个测试模块都会自动生成唯一的测试数据，避免冲突。

## 🔧 开发和扩展

### 添加新的测试模块

1. **创建测试文件**: `test_新模块.py`
2. **继承基类**: 从 `BaseAPITester` 继承
3. **实现测试方法**: 添加具体的测试函数
4. **集成到总测试**: 在 `test_all.py` 中添加模块

### 测试模块模板
```python
#!/usr/bin/env python3
from base_test import BaseAPITester
from colorama import Back

class 新模块Tester(BaseAPITester):
    def test_功能(self):
        self.print_section_header("测试功能", Back.COLOR)
        response = self.make_request("GET", "/endpoint")
        return response is not None
    
    def run_tests(self):
        test_results = [("功能测试", self.test_功能())]
        return self.print_test_summary(test_results)

if __name__ == "__main__":
    tester = 新模块Tester()
    success = tester.run_tests()
    exit(0 if success else 1)
```

## 🔬 OpenAPI Schema 验证详解

### 工作原理

测试框架会自动：
1. 从 `../docs/openapi/` 加载各模块的OpenAPI规范
2. 解析所有endpoint的请求/响应schema定义
3. 自动解析 `$ref` 引用，构建完整的schema
4. 在每次API调用后验证响应是否符合schema
5. 收集并报告所有验证错误

### 使用方法

#### 启用Schema验证（默认）

```python
from base_test import BaseAPITester

class UserTester(BaseAPITester):
    def __init__(self):
        super().__init__(enable_schema_validation=True)  # 默认启用

    def test_user_registration(self):
        # 添加 module="user" 参数启用该请求的schema验证
        response = self.make_request(
            "POST",
            "/users/register",
            data={
                "username": "testuser",
                "email": "test@example.com",
                "password": "password123"
            },
            expect_status=201,
            module="user"  # 指定OpenAPI模块名
        )
        return response is not None
```

#### 禁用Schema验证

```python
# 方法1: 全局禁用
class MyTester(BaseAPITester):
    def __init__(self):
        super().__init__(enable_schema_validation=False)

# 方法2: 单个请求禁用
response = self.make_request(
    "GET", "/users/profile",
    module="user",
    validate_schema=False  # 仅此请求不验证
)
```

### OpenAPI规范目录结构

```
docs/openapi/
├── user.yaml          # 用户API规范
├── problem.yaml       # 题目API规范
├── submission.yaml    # 提交API规范
├── forum.yaml         # 论坛API规范
├── news.yaml          # 新闻API规范
├── admin.yaml         # 管理API规范
└── models/
    ├── common.yaml    # 通用schema定义
    ├── user.yaml      # 用户相关schema
    ├── problem.yaml   # 题目相关schema
    └── ...
```

### 编写OpenAPI规范

#### 主API文件示例 (`user.yaml`)

```yaml
openapi: 3.0.3
info:
  title: User API
  version: 1.0.0

servers:
  - url: http://localhost:8080/api/v1

components:
  $ref: './models/common.yaml#/components'

paths:
  /users/register:
    post:
      tags:
        - 用户管理
      summary: 用户注册
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: './models/user.yaml#/components/schemas/RegisterRequest'
      responses:
        '201':
          description: 注册成功
          content:
            application/json:
              schema:
                $ref: './models/user.yaml#/components/schemas/RegisterResponse'
        '400':
          description: 请求参数错误
          content:
            application/json:
              schema:
                $ref: './models/common.yaml#/components/schemas/Error'
```

#### Schema定义文件示例 (`models/user.yaml`)

```yaml
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        username:
          type: string
        email:
          type: string
          format: email
        role:
          type: string
          enum: [student, teacher, admin, super_admin]
      required: [id, username, email, role]

    RegisterRequest:
      type: object
      required:
        - username
        - email
        - password
      properties:
        username:
          type: string
          minLength: 3
          maxLength: 20
        email:
          type: string
          format: email
        password:
          type: string
          minLength: 6
        nickname:
          type: string
          maxLength: 50

    RegisterResponse:
      type: object
      properties:
        message:
          type: string
        user:
          $ref: '#/components/schemas/User'
```

### Schema验证的好处

1. **自动化测试**: 无需手动编写断言检查每个字段
2. **完整性保证**: 确保响应包含所有必需字段
3. **类型安全**: 自动验证字段类型、格式、范围等
4. **文档一致性**: 确保实现与文档保持同步
5. **回归测试**: 防止API breaking changes
6. **开发效率**: 快速发现API实现与规范的偏差

### 验证错误示例

如果API响应不符合schema，会显示详细错误：

```
⚠️  Schema验证错误 (2) ⚠️

错误 #1:
  模块: user
  请求: POST /users/register
  状态码: 201
  详情: Schema验证失败:
  - user.email: 'invalid-email' is not a 'email'
  - user: 'role' is a required property

错误 #2:
  模块: problem
  请求: GET /problems/1
  状态码: 200
  详情: Schema验证失败:
  - problem.difficulty: 'extreme' is not one of ['easy', 'medium', 'hard']
```

### 常见问题

**Q: Schema验证失败怎么办？**

A: 检查以下几点：
1. 确认OpenAPI规范定义是否正确
2. 确认后端API实现是否符合规范
3. 查看详细的错误信息定位具体字段问题
4. 更新规范或修复实现以保持一致

**Q: 如何跳过特定endpoint的验证？**

A: 两种方法：
```python
# 方法1: 不传module参数（跳过验证）
response = self.make_request("GET", "/some/endpoint")

# 方法2: 显式禁用
response = self.make_request(
    "GET", "/some/endpoint",
    module="user",
    validate_schema=False
)
```

**Q: "Schema未找到"警告是什么意思？**

A: 这表示OpenAPI规范中没有定义该endpoint的schema，验证器会自动跳过。不是错误，只是提醒。可以：
1. 在OpenAPI规范中添加该endpoint的定义
2. 忽略警告（如果不需要验证该endpoint）

**Q: 如何查看可用的schema？**

A: 使用验证器的辅助方法：
```python
from openapi_validator import get_validator

validator = get_validator()
# 查看所有模块的schema
schemas = validator.get_available_schemas()
for module, keys in schemas.items():
    print(f"{module}: {keys}")

# 查看特定模块的schema
user_schemas = validator.get_available_schemas("user")
print(user_schemas)
```

## 🐛 故障排除

### 常见问题

1. **连接被拒绝**
   ```
   Connection refused
   ```
   - 确保后端服务在 localhost:8080 运行
   - 检查 Docker 容器状态

2. **认证失败**
   ```
   JWT 认证失败
   ```
   - 检查用户注册和登录流程
   - 验证 JWT 中间件配置

3. **依赖包问题**
   ```
   Module not found
   ```
   - 运行 `uv sync` 重新安装依赖
   - 检查 Python 版本兼容性

### 调试模式

在测试文件中启用详细日志：
```python
import logging
logging.basicConfig(level=logging.DEBUG)
```

## 📈 性能优化

- 测试用例设计为独立运行，避免相互依赖
- 使用连接池复用 HTTP 连接
- 并发测试支持（可扩展）

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进测试套件！

### 贡献流程
1. Fork 项目
2. 创建功能分支
3. 添加测试用例
4. 提交 Pull Request

## �� 许可证

MIT License 