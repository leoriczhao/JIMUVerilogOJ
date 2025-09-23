# 使用官方 Go 基础镜像
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go模块代理为国内源
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV GOTOOLCHAIN=local

# 使用阿里云 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY judge-service/go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY judge-service/ .

# 重新整理依赖并编译应用
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o judge ./cmd/main.go

# 运行阶段 - 使用 Ubuntu 以便安装 Verilog 工具
FROM ubuntu:22.04

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 使用清华大学镜像源
RUN sed -i 's@//.*archive.ubuntu.com@//mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list && \
    sed -i 's@//.*security.ubuntu.com@//mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list

# 更新包列表并安装必要的工具
RUN apt-get update && apt-get install -y \
    iverilog \
    gtkwave \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/judge .

# 创建工作目录
RUN mkdir -p /tmp/judge

# 设置环境变量
ENV JUDGE_WORK_DIR=/tmp/judge

# 运行应用
CMD ["./judge"] 