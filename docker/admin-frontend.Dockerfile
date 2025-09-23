# 构建阶段
FROM node:20-alpine AS builder
WORKDIR /app

# 配置npm使用中国镜像源
RUN npm config set registry https://registry.npmmirror.com

# 复制package.json和package-lock.json
COPY admin-frontend/package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY admin-frontend/ ./

# 构建项目
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 替换为USTC镜像源，并安装常用工具
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && apk add --no-cache curl

# 复制构建产物到Nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 暴露80端口
EXPOSE 80

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]

