#!/bin/bash

# Verilog OJ 项目部署脚本

set -e

echo "=== Verilog OJ 部署脚本 ==="

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 函数：显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -d, --dev      开发环境部署"
    echo "  -p, --prod     生产环境部署"
    echo "  --stop         停止所有服务"
    echo "  --restart      重启所有服务"
    echo "  --logs         查看服务日志"
    echo "  --status       查看服务状态"
    echo ""
}

# 函数：开发环境部署
deploy_dev() {
    echo "=== 开发环境部署 ==="
    
    # 创建开发环境配置文件
    if [ ! -f ".env.dev" ]; then
        echo "创建开发环境配置文件 .env.dev"
        cat > .env.dev << 'EOF'
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

# === 其他配置 ===
TZ=Asia/Shanghai
LOG_LEVEL=info
EOF
    fi
    
    # 构建并启动服务
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d
    
    echo "开发环境部署完成！"
    echo "后端 API: http://localhost:8080"
    echo "前端页面: http://localhost:80"
    echo "API 文档: http://localhost:8080/docs"
}

# 函数：生产环境部署
deploy_prod() {
    echo "=== 生产环境部署 ==="
    
    # 检查生产环境配置文件
    if [ ! -f ".env.prod" ]; then
        echo "创建生产环境配置文件 .env.prod"
        cat > .env.prod << 'EOF'
# === 数据库配置 ===
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=CHANGE_THIS_IN_PRODUCTION
DB_DATABASE=verilog_oj
DB_SSL_MODE=require

# === Redis 配置 ===
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=CHANGE_THIS_IN_PRODUCTION
REDIS_DB=0

# === 服务器配置 ===
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=release

# === JWT 配置 ===
JWT_SECRET=CHANGE_THIS_TO_STRONG_SECRET_IN_PRODUCTION
JWT_EXPIRES_IN=24

# === 消息队列配置 ===
QUEUE_TYPE=redis
QUEUE_HOST=redis
QUEUE_PORT=6379
QUEUE_PASSWORD=CHANGE_THIS_IN_PRODUCTION
QUEUE_NAME=judge_queue

# === 判题服务配置 ===
JUDGE_WORK_DIR=/tmp/judge

# === 其他配置 ===
TZ=Asia/Shanghai
LOG_LEVEL=info
EOF
        echo "生产环境配置文件已创建！"
        echo "请编辑 .env.prod 文件，修改所有 CHANGE_THIS_IN_PRODUCTION 的值"
        read -p "按回车键继续..."
    fi
    
    # 构建并启动服务
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
    
    echo "生产环境部署完成！"
    echo "请确保防火墙已开放相应端口"
}

# 函数：停止所有服务
stop_services() {
    echo "=== 停止所有服务 ==="
    docker-compose down
    echo "所有服务已停止"
}

# 函数：重启所有服务
restart_services() {
    echo "=== 重启所有服务 ==="
    docker-compose restart
    echo "所有服务已重启"
}

# 函数：查看服务日志
show_logs() {
    echo "=== 服务日志 ==="
    docker-compose logs -f
}

# 函数：查看服务状态
show_status() {
    echo "=== 服务状态 ==="
    docker-compose ps
}

# 函数：初始化数据库
init_database() {
    echo "=== 初始化数据库 ==="
    
    # 等待数据库启动
    echo "等待数据库启动..."
    sleep 10
    
    # 运行数据库迁移（如果需要）
    docker-compose exec backend ./main migrate
    
    echo "数据库初始化完成"
}

# 函数：清理未使用的 Docker 资源
cleanup() {
    echo "=== 清理 Docker 资源 ==="
    docker system prune -f
    docker volume prune -f
    echo "清理完成"
}

# 主函数
main() {
    case "$1" in
        -h|--help)
            show_help
            ;;
        -d|--dev)
            deploy_dev
            ;;
        -p|--prod)
            deploy_prod
            ;;
        --stop)
            stop_services
            ;;
        --restart)
            restart_services
            ;;
        --logs)
            show_logs
            ;;
        --status)
            show_status
            ;;
        --init-db)
            init_database
            ;;
        --cleanup)
            cleanup
            ;;
        "")
            echo "请指定部署选项，使用 -h 查看帮助"
            ;;
        *)
            echo "未知选项: $1"
            echo "使用 -h 查看帮助"
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 