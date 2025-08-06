# 环境变量配置说明

## 概述

Verilog OJ 系统使用环境变量进行配置管理。系统支持多环境配置，包括开发环境和生产环境。

## 配置文件

创建 `.env.dev` 或 `.env.prod` 文件，内容如下：

```bash
# === 数据库配置 ===
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_DATABASE=verilog_oj
DB_SSL_MODE=disable

# === Redis 配置 ===
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0

# === 服务器配置 ===
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=debug

# === JWT 配置 ===
# 生产环境请使用强密码
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24

# === 消息队列配置 ===
QUEUE_TYPE=redis
QUEUE_HOST=redis
QUEUE_PORT=6379
QUEUE_PASSWORD=""
QUEUE_NAME=judge_queue

# === 判题服务配置 ===
JUDGE_WORK_DIR=/tmp/judge

# === 邮件配置（可选） ===
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@verilog-oj.com

# === 文件上传配置 ===
UPLOAD_PATH=/uploads
MAX_UPLOAD_SIZE=10MB

# === 其他配置 ===
TZ=Asia/Shanghai
LOG_LEVEL=info
```

## 配置项说明

### 数据库配置
- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口
- `DB_USERNAME`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_DATABASE`: 数据库名称
- `DB_SSL_MODE`: SSL 模式（disable/require/verify-ca/verify-full）

### Redis 配置
- `REDIS_HOST`: Redis 主机地址
- `REDIS_PORT`: Redis 端口
- `REDIS_PASSWORD`: Redis 密码（可选）
- `REDIS_DB`: Redis 数据库编号

### 服务器配置
- `SERVER_HOST`: 服务器绑定地址
- `SERVER_PORT`: 服务器端口
- `GIN_MODE`: Gin 框架模式（debug/release/test）

### JWT 配置
- `JWT_SECRET`: JWT 签名密钥（生产环境必须使用强密码）
- `JWT_EXPIRES_IN`: Token 过期时间（小时）

### 消息队列配置
- `QUEUE_TYPE`: 队列类型（redis/rabbitmq）
- `QUEUE_HOST`: 队列主机地址
- `QUEUE_PORT`: 队列端口
- `QUEUE_PASSWORD`: 队列密码（可选）
- `QUEUE_NAME`: 队列名称

### 判题服务配置
- `JUDGE_WORK_DIR`: 判题工作目录

## 安全建议

### 生产环境
1. **JWT_SECRET**: 使用至少 32 个字符的随机字符串
2. **数据库密码**: 使用强密码，包含大小写字母、数字和特殊字符
3. **Redis密码**: 设置 Redis 密码保护
4. **SSL**: 在生产环境中启用 SSL/TLS

### 开发环境
- 可以使用默认配置进行快速开发
- 建议使用独立的开发数据库

## Docker Compose 集成

配置文件通过 Docker Compose 的 `env_file` 指令自动加载：

```yaml
services:
  backend:
    env_file:
      - .env.prod  # 或 .env.dev
```

## 配置验证

系统启动时会验证必需的环境变量是否设置，如果缺少关键配置会报错并停止启动。 