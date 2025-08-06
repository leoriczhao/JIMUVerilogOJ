# Verilog OJ 在线判题系统

## 项目概述

这是一个专门为 Verilog 硬件描述语言设计的在线判题系统，采用微服务架构，将判题功能与业务逻辑分离，确保系统的高可用性和可扩展性。

## 系统特性

✨ **核心功能**
- 🔐 用户管理：注册、登录、权限控制
- 📚 题库管理：题目创建、编辑、分类管理
- ⚖️ 智能判题：Verilog代码编译和测试
- 💬 论坛系统：讨论、回复、点赞功能
- 📰 新闻管理：公告和新闻发布
- 📊 统计分析：用户和题目数据统计

🚀 **技术亮点**
- 微服务架构，判题服务独立部署
- Redis消息队列，异步处理判题任务
- Docker容器化，一键部署
- OpenAPI文档，完整的接口规范
- 国内镜像源优化，快速构建

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Frontend    │    │     Backend     │    │  Judge Service  │
│     (Vue.js)    │◄──►│      (Go)       │◄──►│      (Go)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                 │                       │
                                 ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │      Redis      │    │   Verilog Tools │
│   (Database)    │    │  (Cache/Queue)  │    │   (iverilog)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 技术栈

### 后端服务
- **语言**: Go 1.21
- **框架**: Gin (HTTP路由)
- **数据库**: PostgreSQL 15 + GORM
- **缓存**: Redis 7
- **消息队列**: Redis
- **文档**: Swagger/OpenAPI 3.0

### 前端服务（规划中）
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI组件**: Element Plus
- **代码编辑器**: Monaco Editor

### 判题环境
- **编译器**: iverilog (Icarus Verilog)
- **波形查看**: GTKWave
- **容器**: Ubuntu 22.04

### 部署运维
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **监控**: 健康检查 + 日志管理

## 目录结构

```
verilog-oj/
├── backend/                    # 主后端服务
│   ├── cmd/main.go            # 应用入口
│   ├── internal/              # 内部模块
│   │   ├── config/           # 配置管理
│   │   ├── models/           # 数据模型
│   │   ├── handlers/         # HTTP处理器
│   │   ├── services/         # 业务逻辑
│   │   └── middleware/       # 中间件
│   ├── api/                  # API定义
│   ├── pkg/                  # 共享包
│   └── go.mod               # Go模块
├── judge-service/             # 判题服务
│   ├── cmd/main.go           # 应用入口
│   ├── internal/             # 内部模块
│   │   ├── config/          # 配置管理
│   │   ├── judge/           # 判题逻辑
│   │   └── queue/           # 队列处理
│   ├── pkg/                 # 共享包
│   └── go.mod              # Go模块
├── frontend/                 # 前端项目（待开发）
├── docker/                   # Docker配置
│   ├── backend.Dockerfile   # 后端镜像
│   ├── judge.Dockerfile     # 判题服务镜像
│   └── nginx.conf          # Nginx配置
├── docs/                    # 项目文档
│   ├── api.yaml            # OpenAPI规范
│   ├── architecture.md     # 架构设计
│   └── environment.md      # 环境配置
├── scripts/                 # 部署脚本
│   └── deploy.sh           # 部署脚本
├── docker-compose.yml       # 服务编排
└── README.md               # 项目说明
```

## 快速开始

### 环境要求
- Docker 20.0+
- Docker Compose 2.0+

### 一键部署

```bash
# 克隆项目
git clone <repository-url>
cd verilog-oj

# 开发环境部署
./scripts/deploy.sh --dev

# 生产环境部署
./scripts/deploy.sh --prod
```

### 手动部署

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 访问地址

- **前端页面**: http://localhost:80
- **后端API**: http://localhost:8080
- **API文档**: http://localhost:8080/docs
- **数据库**: localhost:5432
- **Redis**: localhost:6379

## 开发指南

### 环境配置

参考 [环境配置文档](docs/environment.md) 创建 `.env.dev` 文件：

```bash
# 数据库配置
DB_HOST=postgres
DB_USERNAME=postgres
DB_PASSWORD=password
DB_DATABASE=verilog_oj

# Redis配置
REDIS_HOST=redis
REDIS_PORT=6379

# JWT配置
JWT_SECRET=your-secret-key
```

### API文档

完整的API文档请查看：
- [OpenAPI规范](docs/api.yaml)
- Swagger UI: http://localhost:8080/docs

### 架构设计

详细的系统架构说明请查看：
- [架构设计文档](docs/architecture.md)

## 功能模块

### 用户管理
- ✅ 用户注册和登录
- ✅ JWT认证和授权
- ✅ 用户资料管理
- ✅ 角色权限控制

### 题库系统
- ✅ 题目创建和编辑
- ✅ 测试用例管理
- ✅ 难度分级和分类
- ✅ 题目搜索和筛选

### 判题引擎
- ✅ Verilog代码编译
- ✅ 测试用例执行
- ✅ 结果评估和反馈
- ✅ 异步队列处理

### 论坛系统
- ✅ 帖子发布和回复
- ✅ 点赞和互动功能
- ✅ 分类讨论区
- ✅ 内容审核机制

### 管理后台
- ✅ 用户管理
- ✅ 题目审核
- ✅ 系统监控
- ✅ 数据统计

## 部署说明

### 生产环境

1. **配置环境变量**
   ```bash
   cp docs/environment.md .env.prod
   # 编辑 .env.prod 设置生产参数
   ```

2. **部署服务**
   ```bash
   ./scripts/deploy.sh --prod
   ```

3. **SSL配置**
   - 配置域名解析
   - 申请SSL证书
   - 更新Nginx配置

### 监控运维

```bash
# 查看服务状态
./scripts/deploy.sh --status

# 查看日志
./scripts/deploy.sh --logs

# 重启服务
./scripts/deploy.sh --restart

# 停止服务
./scripts/deploy.sh --stop
```

## 贡献指南

欢迎贡献代码！请遵循以下流程：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 开源协议

本项目采用 Apache 2.0 协议开源，详见 [LICENSE](LICENSE) 文件。

## 联系我们

- **项目主页**: https://github.com/your-org/verilog-oj
- **问题反馈**: https://github.com/your-org/verilog-oj/issues
- **邮箱**: support@verilog-oj.com

## 致谢

感谢以下开源项目的支持：
- [Go](https://golang.org/) - 后端开发语言
- [Gin](https://gin-gonic.com/) - HTTP框架
- [GORM](https://gorm.io/) - ORM框架
- [Redis](https://redis.io/) - 缓存和队列
- [PostgreSQL](https://www.postgresql.org/) - 数据库
- [Docker](https://www.docker.com/) - 容器化平台
- [Icarus Verilog](http://iverilog.icarus.com/) - Verilog编译器 