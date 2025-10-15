# RBAC权限框架使用指南

## 概述

本项目的RBAC（基于角色的访问控制）权限框架提供了完整的权限管理能力，支持细粒度的权限控制和灵活的角色管理。

## 架构设计

### 核心组件

1. **Go权限检查中间件** (`backend/internal/middleware/`)
   - `rbac.go` - Gin中间件，处理HTTP权限检查
   - `rbac_roles.go` - 权限管理核心，包含缓存和匹配逻辑
   - `rbac_permissions.go` - 权限常量定义和角色权限分配
   - `rbac_test.go` - 单元测试

2. **OpenAPI权限定义** (`docs/openapi/`)
   - `models/common.yaml` - 统一的权限和角色定义
   - `user.yaml` - 用户API权限标记
   - `admin.yaml` - 管理员API权限标记

## 权限命名规范

### 格式
```
{resource}.{action}[.{condition}]
```

### 权限分类

#### 1. 用户权限 (`user.*`)
- `user.profile.read` - 查看个人信息
- `user.profile.update` - 更新个人信息
- `user.password.change` - 修改密码
- `user.avatar.upload` - 上传头像
- `user.create` - 创建用户
- `user.read` - 查看用户信息
- `user.update` - 更新用户信息
- `user.delete` - 删除用户
- `user.ban` - 封禁用户
- `user.unban` - 解封用户

#### 2. 题目权限 (`problem.*`)
- `problem.create` - 创建题目
- `problem.read` - 查看题目
- `problem.list` - 题目列表查看
- `problem.update.own` - 更新自己的题目
- `problem.update.all` - 更新任意题目
- `problem.delete.own` - 删除自己的题目
- `problem.delete.all` - 删除任意题目
- `problem.publish` - 发布题目
- `problem.archive` - 归档题目

#### 3. 测试用例权限 (`testcase.*`)
- `testcase.read.sample` - 查看样例测试用例
- `testcase.read.own` - 查看自己题目的测试用例
- `testcase.read.all` - 查看所有测试用例
- `testcase.create` - 创建测试用例
- `testcase.update` - 更新测试用例
- `testcase.delete` - 删除测试用例

#### 4. 提交权限 (`submission.*`)
- `submission.create` - 提交代码
- `submission.read` - 查看提交记录
- `submission.list` - 提交列表查看
- `submission.manage` - 管理提交记录
- `submission.delete` - 删除提交记录
- `submission.rejudge` - 重新判题

#### 5. 论坛权限 (`forum.*`)
- `forum.post.create` - 发帖
- `forum.post.read` - 查看帖子
- `forum.reply.create` - 回复
- `forum.reply.read` - 查看回复
- `forum.edit.own` - 编辑自己的帖子
- `forum.edit.all` - 编辑任意帖子
- `forum.moderate` - 管理论坛
- `forum.post.lock` - 锁定帖子
- `forum.post.sticky` - 置顶帖子
- `forum.delete` - 删除帖子

#### 6. 新闻权限 (`news.*`)
- `news.read` - 查看新闻
- `news.list` - 新闻列表查看
- `news.create` - 创建新闻
- `news.update` - 更新新闻
- `news.delete` - 删除新闻
- `news.publish` - 发布新闻
- `news.archive` - 归档新闻

#### 7. 系统管理权限 (`manage.*`)
- `manage.users` - 用户管理
- `manage.roles` - 角色管理
- `manage.permissions` - 权限管理
- `manage.system` - 系统管理
- `manage.config` - 系统配置管理
- `manage.logs` - 日志管理
- `manage.content` - 内容管理
- `manage.news` - 新闻管理
- `manage.forum` - 论坛管理
- `manage.submissions` - 提交管理

#### 8. 统计权限 (`stats.*`)
- `stats.read` - 查看统计信息
- `stats.admin` - 管理员统计信息

#### 9. 通配符权限
- `*` - 全权限
- `user.*` - 用户相关所有权限
- `problem.*` - 题目相关所有权限
- `testcase.*` - 测试用例相关所有权限
- `submission.*` - 提交相关所有权限
- `forum.*` - 论坛相关所有权限
- `news.*` - 新闻相关所有权限
- `manage.*` - 管理相关所有权限
- `stats.*` - 统计相关所有权限

## 角色定义

### 1. 学生 (student)
**权限等级**: 1
**描述**: 基础用户权限

**权限列表**:
- `user.profile.*` - 个人资料管理
- `problem.read`, `problem.list` - 查看题目
- `testcase.read.sample` - 查看样例测试用例
- `submission.create`, `submission.read`, `submission.list` - 代码提交
- `forum.post.*`, `forum.reply.*` - 论坛参与
- `news.read`, `news.list` - 查看新闻
- `stats.read` - 基础统计

### 2. 教师 (teacher)
**权限等级**: 2
**描述**: 继承学生权限 + 题目管理权限

**额外权限**:
- `problem.create`, `problem.update.own`, `problem.delete.own` - 题目管理
- `testcase.read.own`, `testcase.*` - 测试用例管理
- `news.create`, `news.update`, `news.delete` - 新闻管理

### 3. 管理员 (admin)
**权限等级**: 3
**描述**: 继承教师权限 + 系统管理权限

**额外权限**:
- `user.*` - 用户管理
- `problem.update.all`, `problem.delete.all`, `problem.publish`, `problem.archive` - 高级题目管理
- `testcase.read.all` - 查看所有测试用例
- `submission.*` - 提交管理
- `forum.edit.all`, `forum.*` - 论坛管理
- `manage.*` - 系统管理
- `stats.admin` - 高级统计

### 4. 超级管理员 (super_admin)
**权限等级**: 4
**描述**: 拥有所有权限

**权限**: `*`

## Go代码使用指南

### 1. 基础权限检查

```go
// 单个权限检查
router.GET("/admin/users",
    middleware.RequirePermission("manage.users"),
    handler.GetUsers,
)

// 多个权限检查（需要任意一个权限）
router.GET("/admin/content",
    middleware.RequireAnyPermission([]string{"manage.news", "manage.forum"}),
    handler.GetContent,
)

// 多个权限检查（需要所有权限）
router.POST("/admin/users/create",
    middleware.RequireAllPermissions([]string{"user.create", "manage.users"}),
    handler.CreateUser,
)
```

### 2. 资源所有权检查

```go
// 需要权限或资源所有权
router.PUT("/problems/:id",
    middleware.RequireOwnershipOrPermission(
        "problem.update.all",
        middleware.GetProblemOwner("id"),
    ),
    handler.UpdateProblem,
)

// 需要角色或资源所有权
router.DELETE("/forum/posts/:id",
    middleware.RequireRoleOrOwnership(
        []string{"admin", "super_admin"},
        middleware.GetForumPostOwner("id"),
    ),
    handler.DeletePost,
)
```

### 3. 可选认证检查

```go
// 可选认证的权限检查
router.GET("/public/stats",
    middleware.OptionalAuthPermission("stats.read"),
    handler.GetPublicStats,
)
```

## OpenAPI使用指南

### 1. 权限标记格式

在OpenAPI路径定义中使用以下标记：

```yaml
# 单个权限检查
x-rbac-permissions: [user.read]

# 多个权限检查（AND关系）
x-rbac-permissions: [user.read, user.update]

# 角色检查
x-rbac-roles: [admin, super_admin]

# 混合检查（高级用法）
x-rbac-require:
  permissions: [user.read]
  roles: [admin]
  operator: "OR"
```

### 2. 实际示例

```yaml
paths:
  /admin/users:
    get:
      summary: 获取用户列表
      security:
        - BearerAuth: []
      x-rbac-permissions: [manage.users]
      responses:
        '200':
          description: 获取成功

  /admin/users/{id}:
    put:
      summary: 更新用户信息
      security:
        - BearerAuth: []
      x-rbac-permissions: [user.update]
      responses:
        '200':
          description: 更新成功

  /problems/{id}:
    put:
      summary: 更新题目
      security:
        - BearerAuth: []
      x-rbac-require:
        permissions: [problem.update.own]
        roles: [admin, teacher]
        operator: "OR"
      responses:
        '200':
          description: 更新成功
```

## 常见使用场景

### 1. 用户自己修改资料

```go
// Go代码
router.PUT("/users/profile",
    middleware.RequirePermission("user.profile.update"),
    handler.UpdateProfile,
)
```

```yaml
# OpenAPI
/users/profile:
  put:
    x-rbac-permissions: [user.profile.update]
```

### 2. 管理员查看所有用户

```go
// Go代码
router.GET("/admin/users",
    middleware.RequirePermission("manage.users"),
    handler.AdminListUsers,
)
```

```yaml
# OpenAPI
/admin/users:
  get:
    x-rbac-permissions: [manage.users]
```

### 3. 教师管理自己的题目

```go
// Go代码
router.PUT("/problems/:id",
    middleware.RequireOwnershipOrPermission(
        "problem.update.all",
        middleware.GetProblemOwner("id"),
    ),
    handler.UpdateProblem,
)
```

```yaml
# OpenAPI
/problems/{id}:
  put:
    x-rbac-require:
      permissions: [problem.update.own]
      roles: [admin, teacher]
      operator: "OR"
```

### 4. 论坛帖子管理

```go
// Go代码
router.DELETE("/forum/posts/:id",
    middleware.RequireRoleOrOwnership(
        []string{"admin", "super_admin"},
        middleware.GetForumPostOwner("id"),
    ),
    handler.DeletePost,
)
```

```yaml
# OpenAPI
/forum/posts/{id}:
  delete:
    x-rbac-require:
      roles: [admin, super_admin]
      permissions: [forum.edit.own]
      operator: "OR"
```

## 权限匹配算法

### 支持的匹配模式

1. **精确匹配**: `user.read` 匹配 `user.read`
2. **全权限通配符**: `*` 匹配所有权限
3. **资源级通配符**: `user.*` 匹配 `user.read`, `user.update` 等
4. **操作级通配符**: `problem.update.*` 匹配 `problem.update.own`, `problem.update.all`

### 匹配优先级

1. 精确匹配 > 资源通配符 > 全权限通配符
2. 支持短路评估，找到匹配即返回

## 缓存机制

### 缓存策略

- **用户权限缓存**: 基于用户ID和角色生成缓存键
- **权限匹配结果缓存**: 避免重复计算
- **并发安全**: 使用读写锁保证线程安全

### 缓存管理

```go
// 清除特定用户缓存
DefaultRBAC.ClearUserCache(userID)

// 清除所有缓存
DefaultRBAC.ClearAllCache()

// 运行时添加权限
DefaultRBAC.AddRolePermission("teacher", "new.permission")

// 运行时移除权限
DefaultRBAC.RemoveRolePermission("teacher", "old.permission")
```

## 最佳实践

### 1. 权限设计原则

- **最小权限原则**: 默认拒绝，明确允许
- **权限分离**: 不同角色权限清晰分离
- **继承性**: 高级角色继承低级角色权限
- **可扩展性**: 支持新增权限和角色

### 2. 代码组织

- **权限常量**: 集中在 `rbac_permissions.go` 中定义
- **中间件使用**: 在路由定义时明确指定权限要求
- **错误处理**: 提供清晰的权限不足错误信息

### 3. OpenAPI文档

- **统一标记**: 使用 `x-rbac-permissions` 标记
- **权限描述**: 在 `common.yaml` 中维护权限定义
- **一致性**: 确保OpenAPI与Go代码权限定义一致

### 4. 测试策略

- **单元测试**: 测试权限匹配逻辑
- **集成测试**: 测试中间件权限检查
- **边界测试**: 测试权限边界情况

## 故障排除

### 常见问题

1. **权限不生效**: 检查权限常量拼写和角色权限分配
2. **缓存问题**: 调用 `ClearUserCache()` 清理缓存
3. **通配符不匹配**: 检查通配符使用格式
4. **OpenAPI文档不同步**: 确保 `common.yaml` 与Go代码一致

### 调试技巧

1. **启用详细日志**: 在中间件中添加权限检查日志
2. **权限验证**: 使用测试工具验证权限配置
3. **缓存监控**: 监控权限缓存命中率

## 扩展指南

### 新增权限

1. 在 `rbac_permissions.go` 中定义权限常量
2. 在角色权限数组中添加新权限
3. 在 `common.yaml` 中添加权限描述
4. 编写相应的测试用例

### 新增角色

1. 在 `rbac_permissions.go` 中定义角色权限数组
2. 在 `common.yaml` 中添加角色定义
3. 更新相关的中间件和API文档
4. 测试角色权限配置

这个RBAC框架为系统提供了完整、安全、高性能的权限管理能力，支持复杂的业务需求和未来的功能扩展。