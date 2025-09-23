# 使用Node.js Alpine镜像作为基础镜像
FROM node:20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 配置npm使用中国镜像源
RUN npm config set registry https://registry.npmmirror.com

# 复制package.json和package-lock.json
COPY frontend/package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY frontend/ ./

# 构建项目
RUN npm run build

# 使用Nginx Alpine镜像作为生产镜像
FROM nginx:alpine

# 替换为USTC镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 安装必要的工具
RUN apk update && apk add --no-cache curl

# 复制构建产物到Nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制Nginx配置文件
COPY docker/nginx.conf /etc/nginx/nginx.conf

# 暴露80端口
EXPOSE 80

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"] 