#!/bin/bash

# Verilog OJ 项目部署脚本

set -e

echo "=== Verilog OJ 部署脚本 ==="

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装（优先使用 V2 语法）
if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    DOCKER_COMPOSE_CMD="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_CMD="docker-compose"
    echo "警告: 使用旧版 docker-compose 命令，建议升级到 Docker Compose V2"
else
    echo "错误: Docker Compose 未安装，请先安装 Docker 和 Docker Compose"
    exit 1
fi

# 函数：显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help         显示帮助信息"
    echo "  -d, --dev          开发环境部署"
    echo "  -p, --prod         生产环境部署"
    echo "  --stop             停止所有服务"
    echo "  --restart          重启所有服务"
    echo "  --logs             查看服务日志"
    echo "  --status           查看服务状态"
    echo "  --reset-passwords  重置密码并重新创建容器"
    echo ""
}

# 函数：创建开发环境模板文件
create_dev_template() {
    cat > .env.dev << 'EOF'
# === 数据库配置 ===
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=CHANGE_THIS_STRONG_PASSWORD
DB_DATABASE=verilog_oj
DB_SSL_MODE=disable

# === Redis 配置 ===
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=CHANGE_THIS_REDIS_PASSWORD
REDIS_DB=0

# === 服务器配置 ===
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=debug

# === JWT 配置 ===
JWT_SECRET=CHANGE_THIS_TO_STRONG_JWT_SECRET
JWT_EXPIRES_IN=24

# === 消息队列配置 ===
QUEUE_TYPE=redis
QUEUE_HOST=redis
QUEUE_PORT=6379
QUEUE_PASSWORD=CHANGE_THIS_REDIS_PASSWORD
QUEUE_NAME=judge_queue

# === 判题服务配置 ===
JUDGE_WORK_DIR=/tmp/judge

# === 其他配置 ===
TZ=Asia/Shanghai
LOG_LEVEL=info
EOF
}

# 函数：检查不安全的密码
check_insecure_passwords() {
    local insecure_found=false

    if grep -q "CHANGE_THIS" .env.dev 2>/dev/null; then
        echo "⚠️  发现未修改的模板密码！"
        insecure_found=true
    fi

    # 检查一些常见的弱密码
    if grep -E "(password|123|test|dev_password|simple)" .env.dev 2>/dev/null | grep -v "#" | grep -q "="; then
        echo "⚠️  发现可能的弱密码！"
        insecure_found=true
    fi

    if [ "$insecure_found" = true ]; then
        echo ""
        echo "🔒 安全建议："
        echo "  1. 使用强密码（至少12位，包含字母数字特殊字符）"
        echo "  2. 数据库密码和Redis密码应该不同"
        echo "  3. JWT密钥应该是随机生成的长字符串"
        echo ""
        read -p "是否继续部署？可能存在安全风险 (y/N): " continue_deploy

        if [[ ! $continue_deploy =~ ^[Yy]$ ]]; then
            echo "部署取消。请先修改密码。"
            exit 1
        fi
    fi
}

# 函数：开发环境部署
deploy_dev() {
    echo "=== 开发环境部署 ==="

    # 检查开发环境配置文件
    if [ ! -f ".env.dev" ]; then
        echo "❌ 缺少开发环境配置文件 .env.dev"
        echo ""
        echo "请手动创建 .env.dev 文件，或使用以下命令创建模板："
        echo "  cp .env.dev.example .env.dev"
        echo ""
        echo "为了安全起见，请设置强密码而不是使用默认值！"
        echo ""
        read -p "是否创建一个模板文件？(y/N): " create_template

        if [[ $create_template =~ ^[Yy]$ ]]; then
            create_dev_template
            echo ""
            echo "⚠️  已创建 .env.dev 模板文件"
            echo "⚠️  请立即编辑此文件并修改所有密码！"
            echo "⚠️  当前使用的是示例密码，不安全！"
            echo ""
            read -p "按回车键继续部署，或 Ctrl+C 取消..."
        else
            echo "部署取消。请先创建 .env.dev 配置文件。"
            exit 1
        fi
    fi

    # 检查是否使用了默认的不安全密码
    check_insecure_passwords
    
    # 构建并启动服务
    $DOCKER_COMPOSE_CMD --env-file .env.dev -f docker-compose.yml -f docker-compose.dev.yml up --build -d
    
    echo "开发环境部署完成！"
    echo "后端 API: http://localhost:8080"
    echo "前端页面: http://localhost:80"
    echo "API 文档: http://localhost:8080/docs"
    echo "数据库: localhost:5432 (用户名: postgres)"
    echo "Redis: localhost:6379 "
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
    $DOCKER_COMPOSE_CMD -f docker-compose.yml -f docker-compose.prod.yml up --build -d
    
    echo "生产环境部署完成！"
    echo "请确保防火墙已开放相应端口"
}

# 函数：停止所有服务
stop_services() {
    echo "=== 停止所有服务 ==="
    $DOCKER_COMPOSE_CMD down
    echo "所有服务已停止"
}

# 函数：重启所有服务
restart_services() {
    echo "=== 重启所有服务 ==="
    $DOCKER_COMPOSE_CMD restart
    echo "所有服务已重启"
}

# 函数：查看服务日志
show_logs() {
    echo "=== 服务日志 ==="
    $DOCKER_COMPOSE_CMD logs -f
}

# 函数：查看服务状态
show_status() {
    echo "=== 服务状态 ==="
    $DOCKER_COMPOSE_CMD ps
}

# 函数：初始化数据库
init_database() {
    echo "=== 初始化数据库 ==="
    
    # 等待数据库启动
    echo "等待数据库启动..."
    sleep 10
    
    # 运行数据库迁移（如果需要）
    $DOCKER_COMPOSE_CMD exec backend ./main migrate
    
    echo "数据库初始化完成"
}

# 函数：重置密码并重新创建容器
reset_passwords() {
    echo "=== 重置密码并重新创建容器 ==="
    echo ""
    echo "⚠️  警告：此操作将："
    echo "  1. 停止并删除所有容器"
    echo "  2. 删除数据库和Redis的数据卷"
    echo "  3. 重新创建容器（所有数据将丢失）"
    echo ""
    read -p "确定要继续吗？这将删除所有数据！(y/N): " confirm_reset

    if [[ ! $confirm_reset =~ ^[Yy]$ ]]; then
        echo "操作取消"
        exit 0
    fi

    echo "停止并删除容器..."
    $DOCKER_COMPOSE_CMD down

    echo "删除数据卷..."
    docker volume rm $(docker volume ls -q | grep -E "(postgres|redis)") 2>/dev/null || true

    echo "清理未使用的资源..."
    docker system prune -f

    echo ""
    echo "✅ 容器和数据已清理完成"
    echo "现在可以修改 .env.dev 中的密码，然后重新部署"
    echo ""
    echo "建议的操作："
    echo "  1. 编辑 .env.dev 文件，修改所有密码"
    echo "  2. 运行 $0 --dev 重新部署"
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
        --reset-passwords)
            reset_passwords
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