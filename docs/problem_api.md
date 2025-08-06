# 题目管理模块 API 文档

## 概述

题目管理模块提供了完整的题目CRUD操作和测试用例管理功能，支持权限控制和数据验证。

## 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token (JWT)
- **内容类型**: `application/json`

## 接口列表

### 1. 获取题目列表

**接口**: `GET /problems`

**描述**: 获取题目列表，支持分页和过滤

**请求参数**:
- `page` (可选): 页码，默认1
- `limit` (可选): 每页数量，默认20，最大100
- `difficulty` (可选): 难度过滤 (Easy/Medium/Hard)
- `category` (可选): 分类过滤

**响应示例**:
```json
{
  "problems": [
    {
      "id": 1,
      "title": "两数之和",
      "description": "实现一个加法器",
      "difficulty": "Easy",
      "category": "数字逻辑",
      "tags": "[基础,逻辑门]",
      "time_limit": 1000,
      "memory_limit": 128,
      "submit_count": 10,
      "accepted_count": 8,
      "is_public": true,
      "author": {
        "id": 1,
        "username": "teacher001",
        "nickname": "张老师"
      },
      "created_at": "2025-07-23T15:30:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

### 2. 获取题目详情

**接口**: `GET /problems/{id}`

**描述**: 获取指定题目的详细信息

**路径参数**:
- `id`: 题目ID

**响应示例**:
```json
{
  "problem": {
    "id": 1,
    "title": "两数之和",
    "description": "实现一个加法器，输入两个数字，输出它们的和。",
    "input_desc": "输入两个整数a和b",
    "output_desc": "输出a+b的结果",
    "difficulty": "Easy",
    "category": "数字逻辑",
    "tags": "[基础,逻辑门]",
    "time_limit": 1000,
    "memory_limit": 128,
    "submit_count": 10,
    "accepted_count": 8,
    "is_public": true,
    "author": {
      "id": 1,
      "username": "teacher001",
      "nickname": "张老师"
    },
    "created_at": "2025-07-23T15:30:00Z"
  }
}
```

### 3. 创建题目

**接口**: `POST /problems`

**描述**: 创建新题目（需要认证）

**认证**: 需要Bearer Token

**请求体**:
```json
{
  "title": "题目标题",
  "description": "题目描述",
  "input_desc": "输入说明",
  "output_desc": "输出说明",
  "difficulty": "Easy",
  "category": "数字逻辑",
  "tags": ["基础", "逻辑门"],
  "time_limit": 1000,
  "memory_limit": 128,
  "test_cases": [
    {
      "input": "1 1",
      "output": "2",
      "is_sample": true
    }
  ]
}
```

**响应示例**:
```json
{
  "message": "题目创建成功",
  "problem": {
    "id": 2,
    "title": "题目标题",
    "description": "题目描述",
    "difficulty": "Easy",
    "category": "数字逻辑",
    "tags": "[基础,逻辑门]",
    "time_limit": 1000,
    "memory_limit": 128,
    "is_public": false,
    "author": {
      "id": 1,
      "username": "teacher001",
      "nickname": "张老师"
    },
    "created_at": "2025-07-23T15:30:00Z"
  }
}
```

### 4. 更新题目

**接口**: `PUT /problems/{id}`

**描述**: 更新题目信息（需要认证，仅作者或管理员）

**认证**: 需要Bearer Token

**路径参数**:
- `id`: 题目ID

**请求体**:
```json
{
  "title": "更新后的标题",
  "description": "更新后的描述",
  "difficulty": "Medium",
  "is_public": true
}
```

**响应示例**:
```json
{
  "message": "题目更新成功",
  "problem": {
    "id": 1,
    "title": "更新后的标题",
    "description": "更新后的描述",
    "difficulty": "Medium",
    "is_public": true
  }
}
```

### 5. 删除题目

**接口**: `DELETE /problems/{id}`

**描述**: 删除题目（需要认证，仅作者或管理员）

**认证**: 需要Bearer Token

**路径参数**:
- `id`: 题目ID

**响应示例**:
```json
{
  "message": "题目删除成功"
}
```

### 6. 获取测试用例

**接口**: `GET /problems/{id}/testcases`

**描述**: 获取题目的测试用例列表

**路径参数**:
- `id`: 题目ID

**响应示例**:
```json
{
  "test_cases": [
    {
      "id": 1,
      "problem_id": 1,
      "input": "1 1",
      "output": "2",
      "is_sample": true,
      "created_at": "2025-07-23T15:30:00Z"
    }
  ]
}
```

### 7. 添加测试用例

**接口**: `POST /problems/{id}/testcases`

**描述**: 为题目添加测试用例（需要认证）

**认证**: 需要Bearer Token

**路径参数**:
- `id`: 题目ID

**请求体**:
```json
{
  "input": "1 1",
  "output": "2",
  "is_sample": true
}
```

**响应示例**:
```json
{
  "message": "测试用例添加成功",
  "test_case": {
    "id": 2,
    "problem_id": 1,
    "input": "1 1",
    "output": "2",
    "is_sample": true,
    "created_at": "2025-07-23T15:30:00Z"
  }
}
```

## 错误响应

### 通用错误格式
```json
{
  "error": "错误代码",
  "message": "错误描述"
}
```

### 常见错误代码

| 错误代码 | HTTP状态码 | 描述 |
|---------|-----------|------|
| `invalid_request` | 400 | 请求参数错误 |
| `invalid_id` | 400 | 无效的题目ID |
| `unauthorized` | 401 | 未认证或认证失败 |
| `problem_not_found` | 404 | 题目不存在 |
| `creation_failed` | 500 | 创建失败 |
| `update_failed` | 500 | 更新失败 |
| `delete_failed` | 500 | 删除失败 |
| `internal_error` | 500 | 内部服务器错误 |

## 权限说明

### 题目操作权限

| 操作 | 学生 | 教师 | 管理员 |
|------|------|------|--------|
| 查看公开题目 | ✅ | ✅ | ✅ |
| 查看私有题目 | ❌ | 仅自己创建 | ✅ |
| 创建题目 | ✅ | ✅ | ✅ |
| 更新题目 | 仅自己创建 | 仅自己创建 | ✅ |
| 删除题目 | 仅自己创建 | 仅自己创建 | ✅ |
| 管理测试用例 | 仅自己创建 | 仅自己创建 | ✅ |

### 角色说明

- **学生 (student)**: 可以创建题目，但只能管理自己创建的题目
- **教师 (teacher)**: 可以创建题目，但只能管理自己创建的题目
- **管理员 (admin)**: 可以管理所有题目

## 数据模型

### Problem 模型
```go
type Problem struct {
    ID            uint           `json:"id"`
    Title         string         `json:"title"`
    Description   string         `json:"description"`
    InputDesc     string         `json:"input_desc"`
    OutputDesc    string         `json:"output_desc"`
    Difficulty    string         `json:"difficulty"`
    Category      string         `json:"category"`
    Tags          string         `json:"tags"`
    TimeLimit     int            `json:"time_limit"`
    MemoryLimit   int            `json:"memory_limit"`
    SubmitCount   int            `json:"submit_count"`
    AcceptedCount int            `json:"accepted_count"`
    IsPublic      bool           `json:"is_public"`
    AuthorID      uint           `json:"author_id"`
    Author        User           `json:"author"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
}
```

### TestCase 模型
```go
type TestCase struct {
    ID        uint           `json:"id"`
    ProblemID uint           `json:"problem_id"`
    Input     string         `json:"input"`
    Output    string         `json:"output"`
    IsSample  bool           `json:"is_sample"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
}
```

## 使用示例

### 创建题目的完整流程

1. **用户登录获取Token**
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "teacher001", "password": "password123"}'
```

2. **创建题目**
```bash
curl -X POST http://localhost:8080/api/v1/problems \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Verilog加法器",
    "description": "实现一个简单的加法器",
    "difficulty": "Easy",
    "category": "数字逻辑",
    "tags": ["基础", "加法器"],
    "time_limit": 1000,
    "memory_limit": 128
  }'
```

3. **添加测试用例**
```bash
curl -X POST http://localhost:8080/api/v1/problems/1/testcases \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "input": "1 1",
    "output": "2",
    "is_sample": true
  }'
```

4. **发布题目**
```bash
curl -X PUT http://localhost:8080/api/v1/problems/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"is_public": true}'
```

## 注意事项

1. **标签格式**: 标签在数据库中存储为JSON字符串格式
2. **时间限制**: 以毫秒为单位
3. **内存限制**: 以MB为单位
4. **权限控制**: 所有修改操作都需要相应的权限
5. **数据验证**: 所有输入都会进行格式和长度验证 