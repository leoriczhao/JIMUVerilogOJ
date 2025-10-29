# E2E 测试框架更新日志

## 2025-10-24 - Schema 验证统计增强

### 新增功能

#### 1. Schema 验证问题分类追踪

在 `base_test.py` 中添加了两种验证问题的追踪：

- **`validation_errors`**: 记录 Schema 验证失败的情况（响应数据与定义不匹配）
- **`validation_skipped`**: 记录缺少 Schema 定义的 API 端点

```python
class BaseAPITester:
    def __init__(self):
        self.validation_errors = []    # 验证失败
        self.validation_skipped = []   # Schema未定义
```

#### 2. 单模块测试中的 Schema 统计

每个测试模块运行后会显示详细的 Schema 验证统计：

```
 📋 Schema验证统计 📋
验证通过: 12
验证失败: 0
跳过验证(未找到schema): 4

 ℹ️  Schema未定义 (4) ℹ️
以下API端点缺少OpenAPI Schema定义：
  1. PUT /problems/25 (状态码: 200)
  2. POST /problems/25/testcases (状态码: 201)
  3. DELETE /problems/25 (状态码: 200)
```

#### 3. 综合测试中的全局 Schema 汇总

`test_all.py` 会收集所有模块的 Schema 问题并进行汇总展示：

```
 📋 OpenAPI Schema验证汇总 📋
================================================================================
Schema验证错误: 0
Schema未定义: 4

ℹ️  缺少Schema定义的API端点：
  1. [problem] PUT /problems/23 (状态码: 200)
  2. [problem] POST /problems/23/testcases (状态码: 201)
  3. [problem] DELETE /problems/23 (状态码: 200)

💡 建议: 在 docs/openapi/ 目录下补充这些API的Schema定义
```

### 技术实现

#### `base_test.py` 更新

**1. 初始化时添加跟踪列表：**
```python
def __init__(self, enable_schema_validation: bool = True):
    self.validation_errors = []    # 验证失败
    self.validation_skipped = []   # Schema未定义
```

**2. 增强 `_validate_response_schema` 方法：**
```python
def _validate_response_schema(self, module, method, endpoint, status_code, response_data):
    is_valid, error_msg = self.validator.validate_response(...)

    if is_valid:
        if error_msg:  # Schema未找到
            self.validation_skipped.append({
                'module': module,
                'method': method,
                'endpoint': endpoint,
                'status_code': status_code,
                'reason': error_msg
            })
    else:  # 验证失败
        self.validation_errors.append({
            'module': module,
            'method': method,
            'endpoint': endpoint,
            'status_code': status_code,
            'error': error_msg
        })
```

**3. 更新 `print_test_summary` 方法：**

添加了 Schema 验证统计部分：
- 显示验证通过/失败/跳过的数量
- 列出验证失败的详细信息（前5个）
- 列出缺少定义的 API 端点（去重后前10个）

#### `test_all.py` 更新

**1. 添加全局收集器：**
```python
class ComprehensiveAPITester:
    def __init__(self):
        self.all_validation_errors = []
        self.all_validation_skipped = []
```

**2. 每个模块测试后收集数据：**
```python
def run_user_tests(self):
    tester = UserTester()
    success = tester.run_tests()
    # 收集schema验证问题
    self.all_validation_errors.extend(tester.validation_errors)
    self.all_validation_skipped.extend(tester.validation_skipped)
    return "user模块", success
```

**3. 添加 `_print_schema_summary` 方法：**

汇总并展示所有模块的 Schema 验证问题。

### 使用方法

#### 运行单个模块测试
```bash
uv run python test_problem.py
```

会在测试结果总结后显示该模块的 Schema 验证统计。

#### 运行完整测试套件
```bash
uv run python test_all.py
```

会在综合测试结果后显示所有模块的 Schema 验证汇总。

### 优势

1. **清晰分类**：区分"验证失败"和"未定义 Schema"两种情况
2. **详细报告**：提供具体的 API 端点、状态码和错误信息
3. **去重显示**：避免重复显示相同的未定义 Schema
4. **实用建议**：提示开发者补充缺失的 Schema 定义
5. **不影响测试**：Schema 验证问题不会导致功能测试失败

### 后续建议

1. **补充 OpenAPI Schema**：根据测试报告补充缺失的 API Schema 定义
2. **修复验证错误**：如果出现验证失败，检查 API 响应格式是否与 Schema 定义一致
3. **持续监控**：定期运行测试，确保新增 API 都有相应的 Schema 定义

### 示例输出

#### 无问题的情况
```
 📋 Schema验证统计 📋
验证通过: 15
验证失败: 0
跳过验证(未找到schema): 0

🎉 所有测试通过！
```

#### 有未定义 Schema 的情况
```
 📋 Schema验证统计 📋
验证通过: 12
验证失败: 0
跳过验证(未找到schema): 4

 ℹ️  Schema未定义 (4) ℹ️
以下API端点缺少OpenAPI Schema定义：
  1. PUT /problems/25 (状态码: 200)
  2. POST /problems/25/testcases (状态码: 201)
```

#### 有验证失败的情况
```
 📋 Schema验证统计 📋
验证通过: 10
验证失败: 2
跳过验证(未找到schema): 3

 ⚠️  Schema验证错误 (2) ⚠️

错误 #1:
  模块: problem
  请求: GET /problems
  状态码: 200
  详情: problems: None is not of type 'array'
```

### 相关文件

- `tests/base_test.py` - 基础测试类，添加验证追踪
- `tests/test_all.py` - 综合测试脚本，添加全局汇总
- `tests/openapi_validator.py` - OpenAPI Schema 验证器（未修改）
- `docs/openapi/*.yaml` - OpenAPI Schema 定义文件

### 版本兼容性

- ✅ 向后兼容：不影响现有测试逻辑
- ✅ 可选功能：可通过 `enable_schema_validation=False` 禁用
- ✅ 不影响测试结果：Schema 问题仅作为信息展示

## 历史版本

### 2025-10-23 - RBAC 框架重构
- 添加用户池管理
- 实现角色切换功能
- 完善权限验证测试

### 2025-10-22 - 初始版本
- 建立 E2E 测试框架
- 实现基础 API 测试
- 添加 OpenAPI Schema 验证
