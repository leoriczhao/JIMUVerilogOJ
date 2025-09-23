# Verilog OJ 前端项目

## 项目概述

这是 Verilog OJ 系统的前端项目，基于 Vue 3 + TypeScript + Element Plus 开发。

## 技术栈

- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **UI 框架**: Element Plus
- **HTTP 客户端**: Axios
- **代码编辑器**: Monaco Editor（用于代码编辑）

## 项目结构

```
frontend/
├── src/
│   ├── api/           # API 服务
│   │   ├── index.ts   # Axios 配置
│   │   ├── problems.ts # 题目相关 API
│   │   ├── news.ts    # 新闻相关 API
│   │   └── forum.ts   # 论坛相关 API
│   ├── components/    # 公共组件
│   │   ├── Header.vue # 导航栏组件
│   │   └── Footer.vue # 页脚组件
│   ├── stores/        # 状态管理
│   │   └── user.ts    # 用户状态
│   ├── views/         # 页面组件
│   │   └── HomeView.vue # 首页
│   ├── router/        # 路由配置
│   ├── assets/        # 静态资源
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── public/            # 公共资源
├── docker/            # Docker 配置
│   ├── frontend.Dockerfile      # 生产环境 Dockerfile
│   └── frontend.dev.Dockerfile  # 开发环境 Dockerfile
└── package.json       # 依赖配置
```

## 功能模块

### 用户模块
- 用户注册/登录
- 个人信息管理
- 用户统计信息

### 题库模块
- 题目列表和搜索
- 题目详情展示
- 题目分类和标签

### 判题模块
- 代码编辑器
- 提交记录查看
- 实时判题结果

### 论坛模块
- 帖子发布和回复
- 点赞和收藏功能
- 分类讨论

### 新闻模块
- 公告和新闻展示
- 新闻分类和搜索

## 开发环境

### 本地开发

1. 安装依赖：
```bash
npm install
```

2. 启动开发服务器：
```bash
npm run dev
```

3. 构建生产版本：
```bash
npm run build
```

### Docker 开发环境

1. 构建并启动开发环境：
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

2. 访问前端应用：
   - 本地开发：http://localhost:3000
   - Docker 开发：http://localhost:3000

### Docker 生产环境

1. 构建并启动生产环境：
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build
```

2. 访问应用：
   - 通过 Nginx：http://localhost
   - 直接访问前端：http://localhost:3000

## API 接口

前端将调用后端提供的 RESTful API，详细接口文档请参考：
- Swagger 文档: http://localhost:8080/docs/
- API 基础地址: http://localhost:8080/api/v1/

## 部署说明

前端项目构建后的静态文件将通过 Nginx 服务器提供服务，配置在项目根目录的 `docker-compose.yml` 中。

### 开发环境特性

- 热重载支持
- 源码映射
- 实时错误提示
- 快速构建

### 生产环境特性

- 代码压缩和优化
- 静态资源缓存
- CDN 支持
- 性能监控

## 注意事项

1. 确保后端 API 服务已启动
2. 开发环境使用 npm 镜像源加速
3. Docker 环境已配置中国镜像源
4. 前端路由需要后端支持 history 模式



