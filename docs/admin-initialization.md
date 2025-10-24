# 管理员初始化方案

## 问题背景

项目引入了 RBAC 权限管理系统，需要管理员用户才能创建其他高权限用户。这产生了"鸡生蛋"问题：
- 没有管理员无法创建其他用户
- 但系统需要初始管理员才能启动

## 解决方案

采用**环境配置驱动的种子数据**方案，根据运行环境自动处理：

### 1. 开发环境（Debug Mode）

**自动创建默认管理员**，用于本地开发和 E2E 测试：

```yaml
用户名: admin
密码: admin123
邮箱: admin@verilogoj.local
角色: admin
```

**特点：**
- ✅ 开箱即用，无需手动配置
- ✅ E2E 测试无缝集成
- ✅ 重启不会重复创建（检查是否已存在）
- ⚠️ 仅用于开发环境，不适合生产

### 2. 生产环境（Release Mode）

**通过环境变量配置首个管理员**：

```bash
# 在 .env.prod 中配置
INIT_ADMIN_USERNAME=your_admin
INIT_ADMIN_EMAIL=admin@yourdomain.com
INIT_ADMIN_PASSWORD=secure_password_123
```

**特点：**
- ✅ 安全可控的密码
- ✅ 自定义用户名和邮箱
- ✅ 仅在数据库无管理员时创建
- ✅ 创建后可删除环境变量

## 使用指南

### 开发环境

1. 启动服务：
```bash
./scripts/deploy.sh --dev
```

2. 系统自动创建 `admin/admin123` 账户

3. 运行 E2E 测试：
```bash
cd tests
uv run python test_all.py
```

### 生产环境首次部署

1. 配置环境变量（`.env.prod`）：
```bash
GIN_MODE=release
INIT_ADMIN_USERNAME=superadmin
INIT_ADMIN_EMAIL=admin@company.com
INIT_ADMIN_PASSWORD=YourSecurePassword123!
```

2. 启动服务：
```bash
./scripts/deploy.sh --prod
```

3. 验证管理员创建：
```bash
docker logs verilog_oj_backend | grep "Custom admin"
# 应该看到：✅ Custom admin user created: superadmin
```

4. **安全措施**：首次启动后，删除或清空环境变量：
```bash
# 在 .env.prod 中
INIT_ADMIN_USERNAME=
INIT_ADMIN_EMAIL=
INIT_ADMIN_PASSWORD=
```

### 生产环境后续部署

如果数据库中已有管理员，系统会跳过创建：
```
2025/10/24 16:14:14 Admin user already exists, skipping seed
```

## 代码实现

### 文件结构

```
backend/
├── cmd/main.go                  # 主程序，调用 initializeAdmin()
├── internal/
│   ├── config/config.go        # 配置定义，包含 InitAdminConfig
│   └── seed/seed.go            # 种子数据逻辑
```

### 核心逻辑

```go
// 开发环境：创建默认管理员
if cfg.Server.Mode == "debug" {
    return seed.SeedDefaultAdmin(db)
}

// 生产环境：根据环境变量创建
if cfg.InitAdmin.Username != "" {
    return seed.SeedCustomAdmin(db, ...)
}
```

### 安全检查

```go
// 检查是否已存在管理员
var count int64
db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

if count > 0 {
    log.Println("Admin user already exists, skipping seed")
    return nil
}
```

## E2E 测试集成

测试框架会自动使用预置管理员（`tests/fixtures/users.py`）：

```python
def _setup_admin(self, api_tester):
    admin_login = {
        "username": "admin",
        "password": "admin123"
    }

    response = api_tester.make_request(
        "POST", "/users/login",
        data=admin_login,
        expect_status=200
    )

    if not response:
        raise Exception("无法登录管理员账户")
```

**无需修改测试代码**，开发环境会自动提供这个账户！

## 安全建议

### 开发环境
- ✅ 使用默认的 `admin/admin123`
- ⚠️ 不要将开发数据库暴露到公网

### 生产环境
- ✅ 使用强密码（至少 12 位，包含大小写、数字、特殊字符）
- ✅ 首次部署后删除环境变量配置
- ✅ 定期更换管理员密码
- ✅ 使用独立的管理员账户，不要共享
- ⚠️ 不要在代码仓库中提交 `.env.prod` 文件

## 故障排除

### 问题 1：E2E 测试报错 "无法登录管理员账户"

**原因**：开发环境未正确创建默认管理员

**解决**：
```bash
# 检查日志
docker logs verilog_oj_backend | grep admin

# 应该看到
✅ Default admin user created (username: admin, password: admin123)

# 如果没有，检查 GIN_MODE 是否为 debug
docker exec verilog_oj_backend env | grep GIN_MODE
```

### 问题 2：生产环境未创建自定义管理员

**原因**：环境变量未正确设置或数据库已有管理员

**解决**：
```bash
# 检查环境变量
docker exec verilog_oj_backend env | grep INIT_ADMIN

# 检查数据库
docker exec verilog_oj_db psql -U postgres -d verilog_oj \
  -c "SELECT id, username, role FROM users WHERE role='admin';"
```

### 问题 3：重复创建管理员

**原因**：代码逻辑错误（不应该发生）

**验证**：
```bash
# 查看数据库中管理员数量
docker exec verilog_oj_db psql -U postgres -d verilog_oj \
  -c "SELECT count(*) FROM users WHERE role='admin';"

# 应该只有 1 个
```

## 总结

| 环境 | 创建方式 | 凭据 | 安全性 |
|------|----------|------|--------|
| 开发 | 自动创建 | admin/admin123 | ⚠️ 仅开发使用 |
| 测试 | 自动创建 | admin/admin123 | ⚠️ E2E 测试 |
| 生产 | 环境变量 | 自定义强密码 | ✅ 高安全性 |

**核心优势**：
- ✅ 开发和测试无缝集成
- ✅ 生产环境安全可控
- ✅ 避免"鸡生蛋"问题
- ✅ 支持多次部署和重启
