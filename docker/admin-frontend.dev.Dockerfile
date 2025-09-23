# 使用Node.js Alpine镜像作为基础镜像
FROM node:20-alpine

# 设置工作目录
WORKDIR /app

# 替换为USTC镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 配置npm使用中国镜像源
RUN npm config set registry https://registry.npmmirror.com

# 安装必要的工具
RUN apk update && apk add --no-cache curl

# 复制package.json和package-lock.json
COPY admin-frontend/package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY admin-frontend/ ./

# 暴露5174端口
EXPOSE 5174

# 启动开发服务器
CMD ["npm", "run", "dev", "--", "--port", "5174", "--host", "0.0.0.0"]

