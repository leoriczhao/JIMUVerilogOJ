# RBAC权限框架架构设计

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        RBAC 权限框架                          │
├─────────────────────────────────────────────────────────────┤
│  中间件层 (Middleware Layer)                                │
│  ├── RequirePermission          │ RequireAnyPermission     │
│  ├── RequireAllPermissions      │ RequireOwnershipOrPermission │
│  ├── RequireRoleOrOwnership     │ OptionalAuthPermission   │
├─────────────────────────────────────────────────────────────┤
│  权限检查层 (Permission Check Layer)                        │
│  ├── DefaultRBAC (全局权限管理器)                           │
│  ├── 权限匹配引擎 (通配符支持)                              │
│  ├── 权限缓存系统 (用户缓存 + 角色缓存)                      │
│  └── 并发安全控制 (读写锁)                                  │
├─────────────────────────────────────────────────────────────┤
│  权限定义层 (Permission Definition Layer)                   │
│  ├── 权限常量定义 (71个权限)                                │
│  ├── 角色权限映射 (4个角色)                                 │
│  └── 权限继承机制 (角色层级)                                │
└─────────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. 中间件层 (`rbac.go`)
**作用**：Gin框架集成层，处理HTTP请求的权限检查

**核心中间件**：
- `RequirePermission` - 单个权限检查
- `RequireAnyPermission` - 任意权限检查
- `RequireAllPermissions` - 所有权限检查
- `RequireOwnershipOrPermission` - 所有权或权限检查
- `RequireRoleOrOwnership` - 角色或所有权检查
- `OptionalAuthPermission` - 可选认证检查

**工作流程**：
```
HTTP请求 → 认证中间件 → RBAC中间件 → 权限检查 → 业务处理
```

### 2. 权限管理器 (`rbac_roles.go`)
**作用**：核心权限检查逻辑和缓存管理

**主要结构**：
```go
type RBAC struct {
    rolePermissions map[string][]string  // 角色权限映射
    userCache       map[string][]string  // 用户权限缓存（按用户+角色隔离）
    cacheMutex      sync.RWMutex         // 读写锁
}
```

**核心方法**：
- `HasPermission` - 检查单个权限
- `HasAnyPermission` - 检查任意权限
- `HasAllPermissions` - 检查所有权限
- `GetUserPermissions` - 获取用户权限列表
- `ClearUserCache` - 清理用户缓存

### 3. 权限定义层 (`rbac_permissions.go`)
**作用**：定义系统中所有权限和角色权限分配

**权限命名规范**：
```
{resource}.{action}[.{condition}]
```

**权限分类**：
- **用户权限** (user.*)：个人资料管理
- **题目权限** (problem.*)：题目CRUD操作
- **测试用例权限** (testcase.*)：测试用例管理
- **提交权限** (submission.*)：代码提交管理
- **论坛权限** (forum.*)：论坛发帖回复
- **新闻权限** (news.*)：新闻内容管理
- **系统管理权限** (manage.*)：系统配置管理

**角色体系**：
```go
Student (16权限)     ← 基础用户权限
Teacher (25权限)     ← 继承Student + 题目管理
Admin (42权限)       ← 继承Teacher + 系统管理
SuperAdmin (1权限)   ← 全权限 (*)
```


## 权限检查流程

### 1. 基础权限检查
```
请求 → 获取用户ID和角色 → 检查权限 → 返回结果
  ↓         ↓              ↓         ↓
Token → user_id, role → HasPermission → ALLOW/DENY
```

### 2. 所有权权限检查
```
请求 → 检查权限 → 检查所有权 → 返回结果
  ↓         ↓          ↓         ↓
Token → HasPermission → IsOwner → ALLOW/DENY
```

### 3. 权限匹配算法
```go
func matchPermission(userPerm, requiredPerm string) bool {
    // 1. 精确匹配
    if userPerm == requiredPerm { return true }

    // 2. 全权限通配符
    if userPerm == "*" { return true }

    // 3. 资源级通配符
    if userPerm == "user.*" && strings.HasPrefix(requiredPerm, "user.") {
        return true
    }

    return false
}
```

## 缓存策略

### 1. 缓存结构
```
用户缓存: map[string][]string
├── key: "<userID>|<role>" → [权限列表]
└── 支持并发读写（按用户和角色隔离）

角色缓存: map[string][]string
├── role_name → [权限列表]
└── 启动时初始化
```

### 2. 缓存键生成
```go
// 用户角色组合生成缓存键
cacheKey := hash(userID + "_" + role)
```

### 3. 缓存管理
- **缓存清理**：`ClearUserCache(userID)` - 清理特定用户缓存
- **全量清理**：`ClearAllCache()` - 清理所有缓存
- **并发安全**：使用 `sync.RWMutex` 保证线程安全

## 性能优化

### 1. 权限缓存
- **命中率高**：用户权限列表缓存，避免重复计算
- **内存效率**：合理的缓存键生成策略
- **并发支持**：读写锁优化高并发场景

### 2. 权限匹配优化
- **O(1)复杂度**：权限匹配算法优化
- **短路评估**：找到匹配即返回
- **通配符优化**：高效的字符串前缀匹配

### 3. 内存管理
- **缓存控制**：支持缓存大小控制
- **垃圾回收**：及时清理无用缓存
- **内存监控**：避免内存泄漏

## 安全设计

### 1. 权限验证原则
- **最小权限原则**：默认拒绝，明确允许
- **权限分离**：不同角色权限隔离
- **权限继承**：角色层级权限传递

### 2. 安全防护
- **权限提升防护**：防止用户自行提升权限
- **资源所有权验证**：支持"作者OR管理员"模式
- **访问控制**：JWT认证集成

### 3. 错误处理
```go
// 权限不足响应
{
    "error": "forbidden",
    "message": "权限不足",
    "details": {
        "required_permission": "problem.create",
        "user_role": "student"
    }
}
```

## 扩展性设计

### 1. 新增权限
```go
// 1. 定义权限常量
const PermNewFeature = "new.feature.action"

// 2. 分配给角色
var RolePermissions = map[string][]string{
    "admin": {
        // ... 其他权限
        PermNewFeature, // 添加新权限
    },
}

// 3. 使用权限中间件
r.POST("/new-feature",
    middleware.RequirePermission(PermNewFeature),
    handler.NewFeature,
)
```

### 2. 新增角色
```go
// 1. 定义角色权限
ModeratorPermissions := append(StudentPermissions,
    PermForumModerate,
    PermNewsManage,
)

// 2. 注册角色
DefaultRBAC.rolePermissions["moderator"] = ModeratorPermissions
```

### 3. 复杂权限逻辑
```go
// 自定义权限检查函数
func RequireCustomPermission(customCheck func(uint, string) bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, role := getUserInfo(c)
        if !customCheck(userID, role) {
            c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
            return
        }
        c.Next()
    }
}
```

## 测试架构

### 1. 测试分层
```
单元测试 (rbac_test.go)
├── 权限匹配逻辑测试
├── 角色权限获取测试
├── 用户权限检查测试
├── 权限缓存测试
└── 性能基准测试

集成测试
├── 中间件集成测试
├── 业务逻辑权限测试
└── 边界条件测试
```

### 2. 测试覆盖
- **功能测试**：所有权限检查功能
- **性能测试**：权限检查基准测试
- **并发测试**：多线程安全验证
- **边界测试**：异常情况处理

## 部署架构

### 1. 文件结构
```
backend/internal/middleware/
├── rbac.go               # 权限检查中间件
├── rbac_permissions.go   # 权限常量定义
├── rbac_roles.go         # 权限管理核心
└── rbac_test.go          # 单元测试
```

### 2. 依赖关系
```
业务代码 → RBAC中间件 → 权限管理器 → 权限定义
    ↑         ↓           ↓         ↓
  处理请求 → 权限检查 → 缓存查询 → 权限匹配
```

### 3. 配置管理
```yaml
rbac:
  cache_enabled: true
  cache_ttl: 300s
  audit_enabled: true
  default_role: "student"
```

## 总结

RBAC权限框架采用分层架构设计，具有以下特点：

1. **清晰的职责分离**：中间件层、权限检查层、定义层各司其职
2. **高性能设计**：权限缓存 + 优化算法保证性能
3. **高度可扩展**：支持新增权限、角色和复杂权限逻辑
4. **安全可靠**：最小权限原则 + 完善的安全防护
5. **易于维护**：集中权限管理 + 详细文档

该架构为系统提供了完整、安全、高性能的权限管理能力，支持未来复杂业务需求的扩展。