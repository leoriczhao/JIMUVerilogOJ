# JIMUVerilogOJ - Verilog 在线判题系统

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![CI Status](https://github.com/leoriczhao/JIMUVerilogOJ/workflows/Backend%20CI/badge.svg)](https://github.com/leoriczhao/JIMUVerilogOJ/actions)

一个专为 Verilog HDL 设计的现代化在线判题平台，采用微服务架构，提供高效的代码编译、测试和评估服务。

[功能特性](#功能特性) • [快速开始](#快速开始) • [开发指南](#开发指南) • [API 文档](#api-文档) • [部署说明](#部署说明)

[English](README.md) | **[简体中文](README_zh.md)**

</div>

---

## 📑 目录

- [项目概述](#项目概述)
- [功能特性](#功能特性)
- [系统架构](#系统架构)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [开发指南](#开发指南)
- [API 文档](#api-文档)
- [测试](#测试)
- [部署说明](#部署说明)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 项目概述

**JIMUVerilogOJ** 是一个功能完善的 Verilog 硬件描述语言在线判题系统，旨在为硬件设计学习和教学提供便捷的代码验证平台。系统采用前后端分离的微服务架构，将判题引擎独立部署，确保高可用性和可扩展性。

### 核心优势

- 🎯 **专业性强** - 专为 Verilog HDL 设计，支持完整的编译和仿真流程
- 🚀 **高性能** - 异步队列处理，判题引擎独立部署，支持高并发
- 🛡️ **安全可靠** - 基于角色的权限控制（RBAC），完善的安全机制
- 📊 **易于扩展** - 微服务架构，模块化设计，便于功能扩展
- 📖 **开发友好** - 完整的 OpenAPI 文档，规范的代码结构

## 功能特性

### 🔐 用户系统
- 用户注册、登录与身份验证
- JWT Token 认证机制
- 基于角色的权限控制（管理员/普通用户）
- 用户资料管理和个性化设置

### 📚 题库管理
- 题目创建、编辑和分类管理
- 测试用例的增删改查
- 难度分级和标签系统
- 题目搜索和筛选功能

### ⚖️ 判题引擎
- Verilog 代码编译检查
- 自动化测试用例执行
- 波形对比和结果验证
- 异步队列处理判题任务
- 详细的错误反馈机制

### 💬 社区论坛
- 讨论帖发布和回复
- 点赞和互动功能
- 分类讨论区
- 内容管理和审核

### 📰 新闻公告
- 系统公告发布
- 新闻动态管理
- 分类和标签系统

### 📊 统计分析
- 用户提交统计
- 题目通过率分析
- 系统使用情况监控

## 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端层                                 │
│  ┌──────────────────┐              ┌──────────────────┐        │
│  │   用户前端界面    │              │   管理员后台      │        │
│  │     (Vue 3)      │              │     (Vue 3)      │        │
│  └──────────────────┘              └──────────────────┘        │
└────────────────┬──────────────────────────┬────────────────────┘
                 │                          │
                 └──────────┬───────────────┘
                            │ HTTP/REST API
                 ┌──────────▼───────────┐
                 │   Nginx (反向代理)    │
                 │                      │
                 └──────────┬───────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
┌────────▼────────┐  ┌──────▼──────┐  ┌───────▼────────┐
│   后端 API      │  │   判题服务   │  │     Redis      │
│     服务        │◄─┤             │  │  (缓存/队列)    │
│     (Go)        │  │    (Go)     │  └────────────────┘
└────────┬────────┘  └──────┬──────┘
         │                  │
         │                  │
    ┌────▼─────┐     ┌──────▼───────┐
    │PostgreSQL│     │Verilog 工具集│
    │  (数据库) │     │  (iverilog)  │
    └──────────┘     └──────────────┘
```

### 服务组件说明

- **用户前端**: 用户交互界面，提供题目浏览、代码提交等功能
- **管理后台**: 管理员界面，用于系统管理和内容审核
- **后端服务**: 核心业务逻辑，处理 API 请求
- **判题服务**: 独立的判题服务，负责代码编译和测试
- **PostgreSQL**: 主数据库，存储用户、题目等数据
- **Redis**: 缓存和消息队列，支持异步判题

## 技术栈

### 后端服务

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端开发语言 |
| Gin | Latest | HTTP 路由框架 |
| GORM | Latest | ORM 框架 |
| PostgreSQL | 15+ | 关系型数据库 |
| Redis | 7+ | 缓存和消息队列 |
| Wire | Latest | 依赖注入 |
| JWT | - | 身份认证 |

### 前端服务

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| TypeScript | Latest | 类型系统 |
| Vite | Latest | 构建工具 |
| Element Plus | Latest | UI 组件库 |
| Monaco Editor | Latest | 代码编辑器 |

### 判题环境

| 技术 | 用途 |
|------|------|
| Icarus Verilog (iverilog) | Verilog 编译器 |
| GTKWave | 波形查看工具 |
| Docker | 隔离的判题环境 |

### 开发运维

| 技术 | 用途 |
|------|------|
| Docker | 容器化 |
| Docker Compose | 服务编排 |
| Nginx | 反向代理 |
| GitHub Actions | CI/CD |
| golangci-lint | 代码质量检查 |

## 快速开始

### 环境要求

- **Docker** 20.0+
- **Docker Compose** 2.0+
- **Go** 1.21+ (本地开发)
- **Node.js** 18+ (前端开发)

### 一键部署

```bash
# 克隆项目
git clone https://github.com/leoriczhao/JIMUVerilogOJ.git
cd JIMUVerilogOJ

# 开发环境部署
./scripts/deploy.sh --dev

# 或者生产环境部署
./scripts/deploy.sh --prod
```

### 手动部署

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
```

### 访问地址

部署成功后，可以访问以下地址：

- **前端页面**: http://localhost:80
- **后端 API**: http://localhost:8080
- **API 文档**: http://localhost:8080/swagger/index.html
- **管理后台**: http://localhost:3000

默认管理员账号：
- 用户名: `admin`
- 密码: `admin123`

## 开发指南

### 后端开发

进入后端目录进行开发：

```bash
cd backend/

# 安装依赖
make deps

# 生成依赖注入代码 (修改 wire.go 后必须执行)
make wire-gen

# 运行服务
make run

# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 生成测试覆盖率
make test-coverage

# 运行所有检查
make check
```

### 前端开发

```bash
cd frontend/

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 代码检查
npm run lint

# 类型检查
npm run type-check

# 生产构建
npm run build
```

### 管理后台开发

```bash
cd admin-frontend/

# 与前端开发流程相同
npm install
npm run dev
```

### API 测试

使用 Python 测试套件：

```bash
cd tests/

# 使用 uv (推荐)
uv run python test_all.py

# 或使用 pip
pip install -r requirements.txt
python test_all.py
```

### 项目结构

```
JIMUVerilogOJ/
├── backend/                    # 后端服务
│   ├── cmd/                    # 应用入口
│   │   └── main.go
│   ├── internal/               # 内部模块
│   │   ├── config/            # 配置管理
│   │   ├── models/            # 数据模型
│   │   ├── handlers/          # HTTP 处理器
│   │   ├── services/          # 业务逻辑层
│   │   ├── repository/        # 数据访问层
│   │   ├── middleware/        # 中间件
│   │   └── wire.go           # 依赖注入配置
│   ├── Makefile              # 构建脚本
│   └── go.mod                # Go 模块依赖
│
├── judge-service/             # 判题服务
│   ├── cmd/
│   ├── internal/
│   │   ├── judge/            # 判题逻辑
│   │   └── queue/            # 队列处理
│   └── go.mod
│
├── frontend/                  # 用户前端
│   ├── src/
│   │   ├── components/       # Vue 组件
│   │   ├── views/           # 页面视图
│   │   ├── router/          # 路由配置
│   │   └── stores/          # 状态管理
│   └── package.json
│
├── admin-frontend/           # 管理后台
│   └── ...
│
├── tests/                    # API 测试
│   ├── test_user.py
│   ├── test_problem.py
│   └── test_submission.py
│
├── docs/                     # 文档
│   └── openapi/             # OpenAPI 规范
│       ├── user.yaml
│       ├── admin.yaml
│       ├── problem.yaml
│       ├── news.yaml
│       └── submission.yaml
│
├── scripts/                  # 部署脚本
│   └── deploy.sh
│
├── docker/                   # Docker 配置
│   ├── backend.Dockerfile
│   └── judge.Dockerfile
│
├── .github/                  # GitHub 配置
│   └── workflows/           # CI/CD 工作流
│
├── docker-compose.yml        # 服务编排（基础）
├── docker-compose.dev.yml    # 开发环境配置
├── docker-compose.prod.yml   # 生产环境配置
├── CLAUDE.md                 # Claude Code 项目说明
└── README.md                 # 本文件
```

## API 文档

### OpenAPI 规范

项目使用 OpenAPI 3.0 规范，文档分模块组织：

- **用户 API**: [docs/openapi/user.yaml](docs/openapi/user.yaml)
- **管理员 API**: [docs/openapi/admin.yaml](docs/openapi/admin.yaml)
- **题目 API**: [docs/openapi/problem.yaml](docs/openapi/problem.yaml)
- **提交 API**: [docs/openapi/submission.yaml](docs/openapi/submission.yaml)
- **新闻 API**: [docs/openapi/news.yaml](docs/openapi/news.yaml)

### 在线文档

启动服务后访问 Swagger UI：
```
http://localhost:8080/swagger/index.html
```

### 常用 API 示例

#### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 提交代码
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

## 测试

### 后端测试

```bash
cd backend/

# 运行所有测试
make test

# 运行特定服务测试
make test-user
make test-services

# 查看测试覆盖率
make test-coverage

# 查看详细输出
make test-verbose
```

### 集成测试

```bash
cd tests/

# 运行所有 API 测试
uv run python test_all.py

# 运行特定测试
uv run python test_user.py
```

### 代码质量检查

```bash
cd backend/

# 格式化代码
make fmt

# 运行 linter
make lint

# 运行 vet
make vet

# 运行所有检查
make check
```

## 部署说明

### 开发环境

```bash
# 使用部署脚本
./scripts/deploy.sh --dev

# 或手动启动
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### 生产环境

1. **配置环境变量**

创建 `.env.prod` 文件：

```bash
# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_secure_password
DB_DATABASE=verilog_oj

# Redis 配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT 配置
JWT_SECRET=your_jwt_secret_key_at_least_32_chars

# 服务器配置
GIN_MODE=release
SERVER_PORT=8080
```

2. **部署服务**

```bash
./scripts/deploy.sh --prod
```

3. **配置 Nginx 和 SSL**

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

### 运维命令

```bash
# 查看服务状态
./scripts/deploy.sh --status
docker-compose ps

# 查看日志
./scripts/deploy.sh --logs
docker-compose logs -f backend

# 重启服务
./scripts/deploy.sh --restart
docker-compose restart backend

# 停止服务
./scripts/deploy.sh --stop
docker-compose down

# 备份数据库
docker-compose exec postgres pg_dump -U postgres verilog_oj > backup.sql

# 恢复数据库
docker-compose exec -T postgres psql -U postgres verilog_oj < backup.sql
```

## 贡献指南

我们欢迎并感谢所有形式的贡献！

### 贡献流程

1. **Fork 项目** 到你的 GitHub 账号
2. **Clone 项目** 到本地：
   ```bash
   git clone https://github.com/YOUR_USERNAME/JIMUVerilogOJ.git
   ```
3. **创建功能分支**：
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **提交更改**：
   ```bash
   git commit -m "feat: add amazing feature"
   ```
5. **推送分支**：
   ```bash
   git push origin feature/amazing-feature
   ```
6. **创建 Pull Request**

### 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat:` 新功能
- `fix:` 错误修复
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `test:` 测试相关
- `chore:` 构建/工具相关

### 代码规范

- **Go**: 遵循 [Effective Go](https://golang.org/doc/effective_go) 和 golangci-lint 规则
- **Vue/TypeScript**: 遵循 ESLint 和 Prettier 配置
- **提交前检查**: 确保运行 `make check` (后端) 和 `npm run lint` (前端)

## 许可证

本项目采用 [Apache License 2.0](LICENSE) 开源协议。

## 联系方式

- **项目主页**: https://github.com/leoriczhao/JIMUVerilogOJ
- **问题反馈**: https://github.com/leoriczhao/JIMUVerilogOJ/issues
- **讨论区**: https://github.com/leoriczhao/JIMUVerilogOJ/discussions

## 致谢

感谢以下开源项目和工具：

- [Go](https://golang.org/) - 强大的后端开发语言
- [Gin](https://gin-gonic.com/) - 高性能 HTTP 框架
- [GORM](https://gorm.io/) - 优雅的 ORM 框架
- [Vue.js](https://vuejs.org/) - 渐进式前端框架
- [PostgreSQL](https://www.postgresql.org/) - 可靠的关系型数据库
- [Redis](https://redis.io/) - 高性能缓存和消息队列
- [Docker](https://www.docker.com/) - 容器化平台
- [Icarus Verilog](http://iverilog.icarus.com/) - Verilog 编译器

## Star History

如果这个项目对你有帮助，请给我们一个 ⭐️！

[![Star History Chart](https://api.star-history.com/svg?repos=leoriczhao/JIMUVerilogOJ&type=Date)](https://star-history.com/#leoriczhao/JIMUVerilogOJ&Date)

---

<div align="center">

**[⬆ 回到顶部](#jimuverilogoj---verilog-在线判题系统)**

Made with ❤️ by [leoriczhao](https://github.com/leoriczhao)

</div>
