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