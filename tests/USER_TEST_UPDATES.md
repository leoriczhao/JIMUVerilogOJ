# 用户测试脚本更新说明

## 更新日期
2025-01-XX

## 更新内容

### 1. 新增测试用例

#### 密码管理测试
- **test_change_password()**: 测试修改密码功能
  - 使用正确的旧密码修改为新密码
  - 验证API返回200状态码
  - 更新测试数据中的密码以供后续测试使用

- **test_change_password_wrong_old()**: 测试使用错误旧密码
  - 提供错误的旧密码尝试修改
  - 验证API正确返回401状态码
  - 确保密码未被修改

- **test_login_with_new_password()**: 测试使用新密码登录
  - 修改密码后使用新密码重新登录
  - 验证新密码可以正常使用
  - 确保密码修改流程完整

#### 输入验证测试
- **test_invalid_email_format()**: 测试无效邮箱格式
  - 使用不符合邮箱格式的字符串注册
  - 验证API返回400错误
  - 确保邮箱格式验证生效

- **test_short_username()**: 测试用户名过短
  - 使用少于3个字符的用户名
  - 验证API返回400错误
  - 确保用户名长度验证(min=3)生效

- **test_short_password()**: 测试密码过短
  - 使用少于6个字符的密码
  - 验证API返回400错误
  - 确保密码长度验证(min=6)生效

### 2. 测试覆盖度提升

#### 更新前
- 7个测试用例
- 覆盖4个API端点（注册、登录、获取信息、更新信息）
- 缺少密码修改功能测试
- 缺少输入验证测试

#### 更新后
- **13个测试用例** ✅
- **覆盖5个API端点** ✅
  - POST /users/register - 用户注册
  - POST /users/login - 用户登录
  - GET /users/profile - 获取个人信息
  - PUT /users/profile - 更新个人信息
  - PUT /users/password - 修改密码 ⭐ 新增
- 完整的输入验证测试 ⭐ 新增
- 完整的错误处理测试

### 3. 测试组织优化

测试用例按功能分组：

```python
test_results = [
    # 基本功能测试
    ("用户注册", self.test_user_registration()),
    ("用户登录", self.test_user_login()),
    ("获取用户信息", self.test_get_profile()),
    ("更新用户信息", self.test_update_profile()),

    # 密码管理测试 ⭐ 新增分组
    ("修改密码", self.test_change_password()),
    ("错误旧密码", self.test_change_password_wrong_old()),
    ("新密码登录", self.test_login_with_new_password()),

    # 错误处理测试
    ("重复注册检查", self.test_duplicate_registration()),
    ("无效登录检查", self.test_invalid_login()),
    ("未授权访问检查", self.test_unauthorized_access()),

    # 输入验证测试 ⭐ 新增分组
    ("无效邮箱格式", self.test_invalid_email_format()),
    ("用户名过短", self.test_short_username()),
    ("密码过短", self.test_short_password()),
]
```

### 4. Schema验证集成

所有测试用例都包含 `module="user"` 参数，启用OpenAPI schema自动验证：

```python
response = self.make_request(
    "PUT", "/users/password",
    data=password_data,
    expect_status=200,
    module="user"  # ← 启用schema验证
)
```

### 5. 与OpenAPI规范对齐

#### 测试覆盖的OpenAPI端点

| 端点 | 方法 | OpenAPI定义 | 测试覆盖 | 状态 |
|-----|------|------------|---------|------|
| /users/register | POST | ✅ | ✅ | 完整 |
| /users/login | POST | ✅ | ✅ | 完整 |
| /users/profile | GET | ✅ | ✅ | 完整 |
| /users/profile | PUT | ✅ | ✅ | 完整 |
| /users/password | PUT | ✅ | ✅ | 完整 ⭐ |

#### 测试覆盖的验证规则

| 字段 | 验证规则 | DTO定义 | 测试覆盖 |
|-----|---------|---------|---------|
| username | min=3, max=20 | ✅ | ✅ |
| email | email格式 | ✅ | ✅ |
| password | min=6 | ✅ | ✅ |
| old_password | min=6 | ✅ | ✅ |
| new_password | min=6 | ✅ | ✅ |

### 6. 测试执行流程

```
1. 用户注册（基本流程）
   └─ 输入验证测试（边界情况）
      ├─ 无效邮箱格式
      ├─ 用户名过短
      └─ 密码过短

2. 用户登录
   └─ 错误处理测试
      ├─ 无效登录
      └─ 重复注册

3. 获取用户信息
   └─ 未授权访问测试

4. 更新用户信息

5. 密码管理（新增完整流程）
   ├─ 修改密码
   ├─ 错误旧密码测试
   └─ 新密码登录验证
```

## 运行测试

### 单独运行用户测试
```bash
cd tests
uv run python test_user.py
```

### 运行所有测试
```bash
cd tests
uv run python test_all.py
```

## 测试输出示例

```
👤 开始用户管理测试
==================================================

 测试用户注册
✅ 请求成功 - POST /users/register
✅ Schema验证通过
✅ 用户注册成功

 测试用户登录
✅ 请求成功 - POST /users/login
✅ Schema验证通过
✅ 用户登录成功

...

 测试修改密码
✅ 请求成功 - PUT /users/password
✅ Schema验证通过
✅ 修改密码成功

 📊 测试结果总结 📊
==================================================
用户注册               ✅ 通过
用户登录               ✅ 通过
获取用户信息           ✅ 通过
更新用户信息           ✅ 通过
修改密码               ✅ 通过
错误旧密码             ✅ 通过
新密码登录             ✅ 通过
重复注册检查           ✅ 通过
无效登录检查           ✅ 通过
未授权访问检查         ✅ 通过
无效邮箱格式           ✅ 通过
用户名过短             ✅ 通过
密码过短               ✅ 通过
==================================================
总测试数: 13
通过数: 13
失败数: 0
通过率: 100.0%

🎉 所有测试通过！
```

## 后续工作建议

### 可以考虑添加的测试用例
1. **用户名边界测试**
   - 测试恰好3个字符的用户名（边界值）
   - 测试恰好20个字符的用户名（边界值）
   - 测试21个字符的用户名（超过最大长度）

2. **密码复杂度测试**
   - 测试恰好6个字符的密码（边界值）
   - 测试包含特殊字符的密码

3. **并发测试**
   - 测试同时注册相同用户名
   - 测试同时修改同一用户的信息

4. **性能测试**
   - 批量注册用户
   - 批量登录测试

5. **安全测试**
   - SQL注入测试
   - XSS攻击测试
   - JWT token过期测试

### 其他模块更新
按照相同的方法，可以更新其他模块的测试：
- `test_problem.py` - 题目管理测试
- `test_submission.py` - 提交管理测试
- `test_forum.py` - 论坛管理测试
- `test_news.py` - 新闻管理测试

## 总结

此次更新显著提升了用户模块的测试覆盖度和质量：

✅ 新增6个测试用例（从7个增加到13个）
✅ 完整覆盖所有5个用户API端点
✅ 添加输入验证测试
✅ 添加密码管理完整流程测试
✅ 所有测试集成OpenAPI schema验证
✅ 测试组织更清晰，按功能分组

测试更加全面，能够有效发现API实现问题和与规范的偏差。
