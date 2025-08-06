# 1. 项目规则
- 本项目使用docker开发和部署， 各个项目的dockerfile都在docker文件夹中，env文件和compose文件在根目录，分为dev和prod两个环境，平常都使用dev环境，只有在部署到生产环境时才使用prod环境
- 因为开发环境处于中国境内，所以docker需要对系统进行镜像替换，比如软件源、go库、node库等
- 本地环境下的python脚本都是用uv来进行管理运行
- 目前后端还没有编写单元测试，前端也没有编写单元测试。 后端的api测试由tests目录的python脚本进行，可以通过运行test_all.py脚本进行测试。

# 2. 分层架构设计规范

## 2.1 三层架构模式
- **Handler层**：使用DTO（数据传输对象）用于API交互
- **Service层**：使用Domain Entity处理业务逻辑
- **Repository层**：使用Data Model进行数据库持久化

## 2.2 文件组织结构
```
internal/
├── handlers/          # API接口层
│   └── *.go           # 处理器实现
├── services/          # 业务逻辑层
│   └── *.go           # 服务实现
├── repository/        # 数据访问层
│   └── *.go           # 仓储实现
├── dto/               # 数据传输对象
│   └── *.go           # API请求/响应结构体
├── domain/            # 领域实体
│   └── *.go           # 业务领域对象
└── models/            # 数据模型
    └── *.go           # GORM数据库模型
```

## 2.3 实体命名规范
- **DTO层**：`CreateXxxRequest`、`UpdateXxxRequest`、`XxxResponse`
- **Domain层**：`XxxEntity`或简单名称（如`User`、`Problem`）
- **Model层**：保持当前GORM模型命名

## 2.4 依赖倒置原则（DIP）
依赖倒置原则是SOLID原则之一，要求：
- **高层模块不应该依赖低层模块，两者都应该依赖抽象**
- **抽象不应该依赖细节，细节应该依赖抽象**

在本项目中的体现：
- Service层定义Repository接口（抽象），不直接依赖具体的数据库实现
- Handler层依赖Service接口，不关心具体的业务实现
- 通过接口隔离，实现层次间的解耦

## 2.5 依赖注入（DI）
依赖注入是实现依赖倒置的具体手段：
- **构造函数注入**：通过构造函数传入依赖的接口实例
- **接口注入**：定义清晰的接口边界
- **容器管理**：使用DI容器统一管理依赖关系

示例结构：
```go
// Service层定义Repository接口
type UserRepository interface {
    GetByID(id uint) (*domain.User, error)
    Create(user *domain.User) error
}

// Service实现依赖接口
type UserService struct {
    userRepo UserRepository  // 依赖抽象而非具体实现
}

// 构造函数注入
func NewUserService(userRepo UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}
```

## 2.6 转换层设计
转换层负责不同层次间的数据转换，**不会影响依赖倒置原则**：
- 转换函数是纯函数，不包含业务逻辑
- 各层通过转换函数进行数据传递，保持层次独立
- DTO ↔ Domain：`ToDomain()`、`ToDTO()`、`FromDomain()`
- Domain ↔ Model：`ToModel()`、`FromModel()`

转换层的优势：
- **解耦数据结构**：各层可以独立演化数据结构
- **类型安全**：编译时检查数据转换的正确性
- **职责分离**：转换逻辑与业务逻辑分离

## 2.7 实施策略

### 2.7.1 阶段一：基础架构重构（当前阶段）
**目标**：建立清晰的分层架构和依赖注入基础

**步骤**：
1. **创建目录结构**
   - 在 `internal/` 下创建 `dto/`、`domain/`、`repository/` 目录
   - 保持现有 `models/`、`handlers/`、`services/` 结构

2. **定义Repository接口**
   - 在 `services/` 层定义各业务模块的Repository接口
   - 例如：`UserRepository`、`ProblemRepository`、`SubmissionRepository`

3. **实现Repository具体类**
   - 在 `repository/` 目录实现具体的Repository
   - 依赖现有的 `models/` 中的GORM模型

4. **引入依赖注入**
   - 在各层目录中创建 `wire.go` 文件
   - 使用Google Wire框架进行依赖注入
   - 在 `internal/wire.go` 中组装整体依赖关系
   - 在 `main.go` 中调用 `InitializeApp()` 初始化应用

5. **添加DTO层**
   - 在 `dto/` 目录创建API请求/响应结构体
   - 在 `handlers/` 层逐步替换直接使用model的地方

### 2.7.2 阶段二：转换层完善（中期目标）
**目标**：实现各层间的数据转换和进一步解耦

**步骤**：
1. **创建Domain实体**
   - 在 `domain/` 目录定义纯业务领域对象
   - 不包含GORM标签和数据库相关逻辑

2. **实现转换函数**
   - DTO ↔ Domain 转换：在 `dto/` 包中实现
   - Domain ↔ Model 转换：在 `repository/` 包中实现
   - 确保转换函数是纯函数，无副作用

3. **重构Service层**
   - Service层只处理Domain对象
   - 通过Repository接口操作数据
   - 移除对Model的直接依赖

4. **完善依赖注入容器**
   - 支持接口自动注入
   - 添加配置管理和环境区分
   - 支持单例和原型模式

### 2.7.3 阶段三：架构优化（长期目标）
**目标**：实现完全的层次解耦和可测试性

**步骤**：
1. **接口抽象化**
   - 为Service层定义接口
   - Handler层依赖Service接口而非具体实现
   - 支持多种实现方式（如缓存装饰器）

2. **单元测试覆盖**
   - 为每层编写单元测试
   - 使用Mock对象测试依赖关系
   - 确保测试覆盖率达到80%以上

3. **性能优化**
   - 添加缓存层（Redis）
   - 实现读写分离
   - 优化数据库查询

### 2.7.4 依赖注入实施细节

**使用Wire框架进行依赖注入**：
本项目采用Google Wire框架实现依赖注入，各层在自己的目录中实现依赖注入配置。

**分层依赖注入结构**：
```
internal/
├── handlers/
│   ├── wire.go           # Handler层依赖注入
│   └── *.go
├── services/
│   ├── wire.go           # Service层依赖注入
│   └── *.go
├── repository/
│   ├── wire.go           # Repository层依赖注入
│   └── *.go
└── wire.go               # 根级别Wire配置
```

**Repository层Wire配置**：
```go
// internal/repository/wire.go
//go:build wireinject
// +build wireinject

package repository

import (
    "github.com/google/wire"
    "gorm.io/gorm"
)

// RepositorySet 提供所有Repository的Wire集合
var RepositorySet = wire.NewSet(
    NewUserRepository,
    NewProblemRepository,
    NewSubmissionRepository,
    wire.Bind(new(services.UserRepository), new(*UserRepository)),
    wire.Bind(new(services.ProblemRepository), new(*ProblemRepository)),
    wire.Bind(new(services.SubmissionRepository), new(*SubmissionRepository)),
)
```

**Service层Wire配置**：
```go
// internal/services/wire.go
//go:build wireinject
// +build wireinject

package services

import "github.com/google/wire"

// ServiceSet 提供所有Service的Wire集合
var ServiceSet = wire.NewSet(
    NewUserService,
    NewProblemService,
    NewSubmissionService,
    wire.Bind(new(UserServiceInterface), new(*UserService)),
    wire.Bind(new(ProblemServiceInterface), new(*ProblemService)),
    wire.Bind(new(SubmissionServiceInterface), new(*SubmissionService)),
)
```

**Handler层Wire配置**：
```go
// internal/handlers/wire.go
//go:build wireinject
// +build wireinject

package handlers

import "github.com/google/wire"

// HandlerSet 提供所有Handler的Wire集合
var HandlerSet = wire.NewSet(
    NewUserHandler,
    NewProblemHandler,
    NewSubmissionHandler,
)
```

**根级别Wire配置**：
```go
// internal/wire.go
//go:build wireinject
// +build wireinject

package internal

import (
    "github.com/google/wire"
    "gorm.io/gorm"
    "github.com/redis/go-redis/v9"
    "your-project/internal/handlers"
    "your-project/internal/services"
    "your-project/internal/repository"
)

// InitializeApp 初始化整个应用
func InitializeApp(db *gorm.DB, redis *redis.Client) (*App, error) {
    wire.Build(
        repository.RepositorySet,
        services.ServiceSet,
        handlers.HandlerSet,
        NewApp,
    )
    return &App{}, nil
}

type App struct {
    UserHandler       *handlers.UserHandler
    ProblemHandler    *handlers.ProblemHandler
    SubmissionHandler *handlers.SubmissionHandler
}

func NewApp(
    userHandler *handlers.UserHandler,
    problemHandler *handlers.ProblemHandler,
    submissionHandler *handlers.SubmissionHandler,
) *App {
    return &App{
        UserHandler:       userHandler,
        ProblemHandler:    problemHandler,
        SubmissionHandler: submissionHandler,
    }
}
```

**渐进式迁移策略**：
- 优先重构核心业务模块（User、Problem）
- 保持向后兼容，新旧代码并存
- 通过Feature Flag控制新架构的启用
- 逐步迁移其他业务模块