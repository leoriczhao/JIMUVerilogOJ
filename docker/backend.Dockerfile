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
COPY backend/go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY backend/ .

# 重新整理依赖并编译应用
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 使用阿里云 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的依赖
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 