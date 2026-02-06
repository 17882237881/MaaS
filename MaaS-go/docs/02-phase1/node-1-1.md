# 节点1.1：项目初始化与架构设计

> 📅 **学习时间**：3天  
> 🎯 **目标**：搭建项目骨架，理解微服务架构设计原则

## 本节你将学到

1. Go Modules依赖管理
2. 微服务目录结构设计
3. Git工作流规范
4. 项目初始化实操

---

## 技术详解

### 1. Go Modules 深度解析

#### 什么是Go Modules？
Go Modules是Go语言的依赖管理系统，从Go 1.11开始引入，在Go 1.16成为默认方式。

**解决的问题**：
- 依赖版本管理（不再依赖GOPATH）
- 可重现的构建（go.sum锁定版本）
- 支持语义化版本（Semantic Versioning）

#### 核心文件说明

**go.mod** - 模块定义文件
```go
module maas-platform  // 模块名称（通常是仓库路径）

go 1.21  // 要求的Go版本

require (
    // 直接依赖
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
    
    // indirect 间接依赖
    github.com/go-playground/validator/v10 v10.16.0 // indirect
)
```

**go.sum** - 依赖校验文件
```
github.com/gin-gonic/gin v1.9.1 h1:Q3nJ5xbvCcG...=
github.com/gin-gonic/gin v1.9.1/go.mod h1:YyFQF...=
```
每行包含：模块路径、版本、哈希值（确保下载的代码未被篡改）

#### 常用命令详解

```bash
# 1. 初始化模块（只需执行一次）
go mod init maas-platform
# 说明：在当前目录创建go.mod文件

# 2. 下载所有依赖（根据go.mod）
go mod download
# 说明：下载依赖到 $GOPATH/pkg/mod 缓存目录

# 3. 整理依赖（添加缺失、删除多余）
go mod tidy
# 说明：分析代码中import的依赖，自动更新go.mod

# 4. 查看依赖关系树
go mod graph
# 输出示例：
# maas-platform github.com/gin-gonic/gin@v1.9.1
# github.com/gin-gonic/gin@v1.9.1 github.com/go-playground/validator/v10@v10.16.0

# 5. 清理未使用的缓存
go clean -modcache

# 6. 更新依赖到最新版本
go get -u ./...
# -u = update，更新所有依赖

# 7. 查看可更新的依赖
go list -u -m all
```

#### 版本管理策略

**语义化版本（SemVer）**：
```
v1.2.3
│ │ │
│ │ └── Patch：Bug修复
│ └──── Minor：新功能（向后兼容）
└────── Major：重大变更（可能不兼容）
```

**go.mod中的版本标识**：
```go
require (
    github.com/foo/bar v1.2.3      // 精确版本
    github.com/foo/bar v1.2.3+incompatible  // 非模块版本
    github.com/foo/bar v0.0.0-20231201123456-abcdef123456  // 伪版本（commit hash）
)
```

#### 私有模块配置
如果使用了私有仓库（如公司内部的GitLab）：
```bash
# 配置Git使用SSH而不是HTTPS（针对私有仓库）
git config --global url."git@github.com:".insteadOf "https://github.com/"

# 或者配置GOPRIVATE环境变量（不通过代理）
export GOPRIVATE=gitlab.company.com
```

---

### 2. 微服务架构设计

#### 什么是微服务？

**单体应用（Monolithic）**：
```
┌─────────────────────────────────┐
│           单个程序               │
│  ┌─────┐ ┌─────┐ ┌─────┐       │
│  │用户 │ │订单 │ │支付 │       │
│  │模块 │ │模块 │ │模块 │       │
│  └─────┘ └─────┘ └─────┘       │
│                                 │
│  共用数据库、代码耦合            │
└─────────────────────────────────┘
优点：开发简单、部署方便
缺点：代码膨胀、技术栈单一、扩展困难
```

**微服务（Microservices）**：
```
┌─────────┐ ┌─────────┐ ┌─────────┐
│ 用户服务 │ │ 订单服务 │ │ 支付服务 │
│  Go     │ │  Java   │ │  Go     │
│  MySQL  │ │ Postgre │ │ MongoDB │
└────┬────┘ └────┬────┘ └────┬────┘
     └───────────┴───────────┘
               │
         ┌─────┴─────┐
         │ API网关   │
         └───────────┘
优点：独立部署、技术异构、团队自治
缺点：分布式复杂度、运维成本
```

#### MaaS平台的服务拆分

我们的平台拆分为以下服务：

| 服务 | 职责 | 数据库 |
|------|------|--------|
| **API Gateway** | 统一入口、路由转发、认证鉴权 | 无 |
| **Model Registry** | 模型元数据管理、版本控制 | PostgreSQL |
| **Inference Engine** | 推理请求处理、负载均衡 | Redis |
| **Deployment Controller** | 模型部署、扩缩容 | PostgreSQL |
| **User Center** | 用户管理、权限控制 | PostgreSQL |
| **Billing Service** | 计费、配额管理 | PostgreSQL |

**拆分原则**：
1. **单一职责**：一个服务只做一件事
2. **独立部署**：服务间通过API通信，可独立发布
3. **独立团队**：理想情况下一个团队维护一个服务
4. **独立数据库**：每个服务有自己的数据存储

---

### 3. 目录结构设计

#### 分层架构（Layered Architecture）

```
service/
├── cmd/                    # 程序入口（main函数）
│   └── main.go            # 可以构建成二进制文件
│
├── internal/              # 私有代码（其他模块不能导入）
│   ├── config/           # 配置管理
│   │   └── config.go
│   ├── handler/          # HTTP处理器（Controller层）
│   │   ├── user_handler.go
│   │   └── model_handler.go
│   ├── middleware/       # HTTP中间件
│   │   ├── auth.go
│   │   └── logger.go
│   ├── model/            # 数据模型（Entity层）
│   │   ├── user.go
│   │   └── model.go
│   ├── repository/       # 数据访问层（DAO层）
│   │   ├── user_repository.go
│   │   └── model_repository.go
│   ├── router/           # 路由定义
│   │   └── router.go
│   └── service/          # 业务逻辑层
│       ├── user_service.go
│       └── model_service.go
│
├── pkg/                   # 公共库（可被其他模块导入）
│   ├── logger/           # 日志工具
│   │   └── logger.go
│   └── utils/            # 工具函数
│       └── utils.go
│
└── go.mod                # 模块定义
```

#### 各层职责说明

**1. Handler层（控制器层）**
```go
// 职责：处理HTTP请求，参数校验，调用Service
// 不包含业务逻辑

type UserHandler struct {
    service UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    // 1. 参数绑定和校验
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 2. 调用Service
    user, err := h.service.CreateUser(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 3. 返回响应
    c.JSON(201, user)
}
```

**2. Service层（业务逻辑层）**
```go
// 职责：实现业务逻辑，编排Repository操作
// 处理事务，协调多个Repository

type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id string) (*User, error)
}

type userService struct {
    repo UserRepository
    // 可以依赖其他Service
}

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // 业务逻辑：检查邮箱是否已存在
    existing, _ := s.repo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }
    
    // 业务逻辑：密码加密
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
    
    // 创建用户
    user := &User{
        ID:       uuid.New().String(),
        Email:    req.Email,
        Password: string(hashedPassword),
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

**3. Repository层（数据访问层）**
```go
// 职责：数据库操作，封装查询逻辑
// 不包含业务逻辑，只处理数据存取

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter UserFilter) ([]*User, error)
}

type gormUserRepository struct {
    db *gorm.DB
}

func (r *gormUserRepository) Create(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

**4. Model层（实体层）**
```go
// 职责：定义数据结构
// 可以包含简单的业务方法

type User struct {
    ID        string         `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex" json:"email"`
    Password  string         `json:"-"`  // 不序列化到JSON
    Name      string         `json:"name"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 表名
func (User) TableName() string {
    return "users"
}

// 验证方法
func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("email is required")
    }
    return nil
}
```

#### 为什么要分层？

**1. 关注点分离**
- Handler只关心HTTP协议
- Service只关心业务逻辑
- Repository只关心数据库

**2. 可测试性**
```go
// 可以Mock Repository测试Service
func TestCreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
        Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).
        Return(nil)
    
    user, err := service.CreateUser(context.Background(), req)
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

**3. 可替换性**
- 可以更换数据库（MySQL→PostgreSQL）而不影响Service
- 可以更换Web框架（Gin→Echo）而不影响业务逻辑

---

### 4. Git工作流规范

#### 分支模型（Git Flow简化版）

```
main (生产分支)
 │
 ├── feature/user-auth (功能分支)
 │      │
 │      └── commit "feat: add JWT auth"
 │      └── commit "feat: add login API"
 │
 ├── feature/model-upload (功能分支)
 │
 └── hotfix/fix-memory-leak (热修复分支)
```

**分支说明**：
- **main**：生产分支，永远可部署
- **feature/***：功能分支，从main创建，开发完合并回main
- **hotfix/***：热修复分支，从main创建，修复紧急Bug

#### Commit Message规范（Conventional Commits）

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型（type）**：
- **feat**: 新功能
- **fix**: Bug修复
- **docs**: 文档更新
- **style**: 代码格式（不影响功能）
- **refactor**: 重构
- **test**: 测试相关
- **chore**: 构建/工具相关

**示例**：
```bash
# 简单提交
git commit -m "feat(api): add user login endpoint"

# 详细提交
git commit -m "feat(api): add user login endpoint

- Implement JWT token generation
- Add password validation
- Update swagger documentation

Closes #123"
```

---

## 实操任务

### 任务1：创建项目根目录和初始化Git

```bash
# 1. 创建项目目录
mkdir MaaS-platform
cd MaaS-platform

# 2. 初始化Git仓库
git init

# 3. 创建.gitignore
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Environment variables
.env
.env.local

# Build output
bin/
dist/

# Log files
*.log
EOF

# 4. 提交
git add .gitignore
git commit -m "chore: add .gitignore"
```

### 任务2：初始化Go Modules

```bash
# 初始化模块（模块名通常是仓库路径）
go mod init github.com/17882237881/MaaS

# 查看生成的go.mod
cat go.mod

# 提交
git add go.mod
git commit -m "chore: initialize go modules"
```

### 任务3：创建服务目录结构

```bash
# 创建API Gateway目录结构
mkdir -p api-gateway/{cmd,internal/{config,handler,middleware,model,repository,router,service},pkg/{logger,utils}}

# 创建Model Registry目录结构
mkdir -p model-registry/{cmd,internal/{config,handler,middleware,model,repository,router,service},pkg/{logger,utils}}

# 创建其他目录
mkdir -p {deploy/{docker,k8s},docs,shared/{proto,errors}}

# 创建占位文件（Go文件需要package声明才能编译）
```

### 任务4：编写第一个main.go

**api-gateway/cmd/main.go**：
```go
package main

import "fmt"

func main() {
    fmt.Println("MaaS API Gateway Starting...")
    fmt.Println("Version: 0.1.0")
    fmt.Println("Listening on :8080")
    
    // 保持运行
    select {}
}
```

**model-registry/cmd/main.go**：
```go
package main

import "fmt"

func main() {
    fmt.Println("MaaS Model Registry Starting...")
    fmt.Println("Version: 0.1.0")
    fmt.Println("Listening on :8081")
    
    select {}
}
```

### 任务5：创建README.md

```markdown
# MaaS Platform

Model-as-a-Service Platform - 模型即服务平台

## 项目结构

```
├── api-gateway/          # API网关服务
├── model-registry/       # 模型注册服务
├── deploy/               # 部署配置
│   ├── docker/          # Docker配置
│   └── k8s/             # Kubernetes配置
├── docs/                 # 文档
└── shared/               # 共享代码
    ├── proto/           # Protocol Buffers
    └── errors/          # 公共错误定义
```

## 快速开始

### 环境要求
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### 本地开发

1. 克隆仓库
   ```bash
   git clone https://github.com/17882237881/MaaS.git
   cd MaaS
   ```

2. 启动依赖服务
   ```bash
   cd deploy/docker
   docker-compose up -d postgres redis
   ```

3. 运行API Gateway
   ```bash
   cd api-gateway
   go run cmd/main.go
   ```

4. 运行Model Registry
   ```bash
   cd model-registry
   go run cmd/main.go
   ```

## 文档

- [项目总览](./docs/01-overview/README.md)
- [阶段1：基础架构](./docs/02-phase1/README.md)

## 许可证

MIT
```

### 任务6：提交代码

```bash
# 添加所有文件
git add .

# 提交
git commit -m "feat(project): initialize project structure

- Add api-gateway and model-registry service structure
- Create basic main.go for each service
- Add README.md with project overview
- Setup Go modules"

# 推送到GitHub（先创建远程仓库）
git remote add origin https://github.com/17882237881/MaaS.git
git branch -M main
git push -u origin main
```

---

## 检查清单

完成本节点后，请确认：

- [ ] 项目目录结构完整
- [ ] `go mod init`成功执行，生成了go.mod
- [ ] 每个服务的main.go能独立编译运行
- [ ] Git仓库初始化，有至少2个commit
- [ ] README.md包含项目结构说明
- [ ] .gitignore配置正确

---

## 验证方法

### 1. 验证目录结构
```bash
# 应该看到类似输出
tree -L 3 -d
.
├── api-gateway
│   ├── cmd
│   ├── internal
│   │   ├── config
│   │   ├── handler
│   │   ├── middleware
│   │   ├── model
│   │   ├── repository
│   │   ├── router
│   │   └── service
│   └── pkg
│       ├── logger
│       └── utils
├── deploy
│   ├── docker
│   └── k8s
├── docs
│   ├── 01-overview
│   └── 02-phase1
├── model-registry
│   ├── cmd
│   ├── internal
│   └── pkg
└── shared
    ├── errors
    └── proto
```

### 2. 验证Go模块
```bash
# 应该看到go.mod和go.sum（如果有依赖）
ls -la *.mod

# 检查模块内容
cat go.mod
```

### 3. 验证Git
```bash
# 查看提交历史
git log --oneline

# 应该看到：
# abc1234 feat(project): initialize project structure
# def5678 chore: initialize go modules
# ghi9012 chore: add .gitignore
```

---

## 常见问题

### Q: go mod init报错"already exists"
**A**: 删除已有的go.mod重新初始化
```bash
rm go.mod
go mod init github.com/17882237881/MaaS
```

### Q: 应该选择什么作为模块名？
**A**: 推荐使用仓库路径：
- GitHub: `github.com/username/repo`
- GitLab: `gitlab.com/username/repo`
- 私有仓库: `company.com/project/module`

### Q: internal目录的作用是什么？
**A**: Go 1.4引入的特殊目录，其中的代码只能被该目录的父目录导入。用于封装实现细节，防止外部滥用。

### Q: pkg目录什么时候用？
**A**: 当代码需要被其他模块导入时使用。internal中的代码只能在当前模块使用。

---

## 下一步

完成本节点后，你已经搭建了项目的骨架。接下来进入：

**节点1.2：API Gateway核心** → [继续学习](./node-1-2.md)

在那里你将：
- 学习Gin框架核心概念
- 实现第一个HTTP接口
- 添加中间件支持
- 集成Swagger文档

---

## 代码变更记录

本节详细记录了节点1.1中创建的所有文件及其代码内容。

### 提交信息
```
feat(phase1/node1.1): initialize project structure

- Add api-gateway service with complete directory structure
- Add model-registry service with complete directory structure  
- Create go.mod with project dependencies
- Add README.md with project overview
- Add Makefile for build automation
- Add .gitignore for Go projects
- Create placeholder files for all layers
```

### 创建的文件清单

#### 根目录文件
1. `go.mod` - Go模块定义
2. `README.md` - 项目说明文档
3. `.gitignore` - Git忽略文件配置
4. `Makefile` - 构建自动化脚本

#### API Gateway服务 (api-gateway/)
1. `cmd/main.go` - 服务入口
2. `internal/config/config.go` - 配置管理
3. `internal/handler/handler.go` - HTTP处理器
4. `internal/middleware/middleware.go` - 中间件
5. `internal/model/model.go` - 数据模型
6. `internal/router/router.go` - 路由定义
7. `internal/service/service.go` - 业务逻辑
8. `internal/repository/repository.go` - 数据访问
9. `pkg/logger/logger.go` - 日志工具
10. `pkg/utils/utils.go` - 工具函数

#### Model Registry服务 (model-registry/)
1. `cmd/main.go` - 服务入口
2. `internal/config/config.go` - 配置管理
3. `internal/handler/handler.go` - HTTP处理器
4. `internal/middleware/middleware.go` - 中间件
5. `internal/model/model.go` - 数据模型
6. `internal/router/router.go` - 路由定义
7. `internal/service/service.go` - 业务逻辑
8. `internal/repository/repository.go` - 数据访问
9. `pkg/logger/logger.go` - 日志工具
10. `pkg/utils/utils.go` - 工具函数

### 文件内容详情

#### 1. go.mod
```go
module maas-platform

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.5.0
	go.uber.org/zap v1.26.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
	github.com/redis/go-redis/v9 v9.3.0
	github.com/spf13/viper v1.18.1
	github.com/segmentio/kafka-go v0.4.46
	google.golang.org/grpc v1.60.0
	google.golang.org/protobuf v1.32.0
	github.com/hibiken/asynq v0.24.1
	github.com/casbin/casbin/v2 v2.79.0
	github.com/prometheus/client_golang v1.17.0
	github.com/jaegertracing/jaeger-client-go v2.30.0+incompatible
	github.com/swaggo/swag v1.16.2
	github.com/swaggo/gin-swagger v1.6.0
)
```

**说明**：定义了项目依赖，包括Gin框架、GORM、Redis、Kafka、gRPC等后续阶段需要的库。

#### 2. README.md
包含项目简介、技术栈、项目结构、快速开始指南等。

#### 3. .gitignore
```gitignore
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Environment variables
.env
.env.local

# Build output
bin/
dist/

# Log files
*.log

# Temporary files
tmp/
temp/
```

**说明**：标准的Go项目.gitignore，排除编译产物、IDE配置、临时文件等。

#### 4. Makefile
```makefile
.PHONY: help build run test clean

help:
	@echo "Available targets:"
	@echo "  build      - Build all services"
	@echo "  run-api    - Run API Gateway"
	@echo "  run-model  - Run Model Registry"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"

build:
	@echo "Building API Gateway..."
	@cd api-gateway && go build -o ../bin/api-gateway ./cmd/main.go
	@echo "Building Model Registry..."
	@cd model-registry && go build -o ../bin/model-registry ./cmd/main.go
	@echo "Build complete!"

run-api:
	@cd api-gateway && go run ./cmd/main.go

run-model:
	@cd model-registry && go run ./cmd/main.go

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@rm -rf bin/
	@echo "Clean complete!"
```

**说明**：提供了构建、运行、测试、清理等常用命令。

#### 5. api-gateway/cmd/main.go
```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("MaaS API Gateway Starting...")
	fmt.Println("Version: 0.1.0")
	fmt.Println("Listening on :8080")
	
	// 保持运行
	select {}
}
```

**说明**：API Gateway服务的入口，目前只是打印启动信息，后续会添加Gin框架。

#### 6. api-gateway/internal/config/config.go
```go
package config

// Config holds all configuration for the API Gateway
type Config struct {
	Environment string
	Port        int
	LogLevel    string
}

// Load returns the application configuration
func Load() *Config {
	return &Config{
		Environment: "development",
		Port:        8080,
		LogLevel:    "info",
	}
}
```

**说明**：配置结构体，后续会使用Viper从配置文件或环境变量加载。

#### 7. api-gateway/internal/handler/handler.go
```go
package handler

// Handler handles HTTP requests
type Handler struct{}

// NewHandler creates a new handler
func NewHandler() *Handler {
	return &Handler{}
}
```

**说明**：HTTP处理器占位，后续会添加具体的请求处理方法。

#### 8. api-gateway/internal/middleware/middleware.go
```go
package middleware

// Middleware placeholder for future implementation
```

**说明**：中间件占位，后续会添加Logger、Recovery、CORS等中间件。

#### 9. api-gateway/internal/model/model.go
```go
package model

// Model placeholder for future data models
```

**说明**：数据模型占位，后续会定义用户、模型等实体。

#### 10. api-gateway/internal/router/router.go
```go
package router

// Router placeholder for future route definitions
```

**说明**：路由定义占位，后续会添加具体的路由规则。

#### 11. api-gateway/internal/service/service.go
```go
package service

// Service placeholder for future business logic
```

**说明**：业务逻辑层占位，后续会实现业务功能。

#### 12. api-gateway/internal/repository/repository.go
```go
package repository

// Repository placeholder for future data access layer
```

**说明**：数据访问层占位，后续会实现数据库操作。

#### 13. api-gateway/pkg/logger/logger.go
```go
package logger

// Logger placeholder for future logging implementation
```

**说明**：日志工具占位，后续会集成Zap日志库。

#### 14. api-gateway/pkg/utils/utils.go
```go
package utils

// Utils placeholder for future utility functions
```

**说明**：工具函数占位，后续会添加通用工具。

#### 15-24. model-registry/ 下文件
与api-gateway结构相同，只是端口改为8081。

### 目录结构总览

执行 `tree -L 3` 后的输出：
```
.
├── Makefile
├── README.md
├── api-gateway
│   ├── cmd
│   │   └── main.go
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── handler
│   │   │   └── handler.go
│   │   ├── middleware
│   │   │   └── middleware.go
│   │   ├── model
│   │   │   └── model.go
│   │   ├── repository
│   │   │   └── repository.go
│   │   ├── router
│   │   │   └── router.go
│   │   └── service
│   │       └── service.go
│   └── pkg
│       ├── logger
│       │   └── logger.go
│       └── utils
│           └── utils.go
├── deploy
├── docs
├── go.mod
├── model-registry
│   ├── cmd
│   │   └── main.go
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── handler
│   │   │   └── handler.go
│   │   ├── middleware
│   │   │   └── middleware.go
│   │   ├── model
│   │   │   └── model.go
│   │   ├── repository
│   │   │   └── repository.go
│   │   ├── router
│   │   │   └── router.go
│   │   └── service
│   │       └── service.go
│   └── pkg
│       ├── logger
│       │   └── logger.go
│       └── utils
│           └── utils.go
└── shared
```

### 验证步骤

1. **验证Go模块**
   ```bash
   cat go.mod
   # 应该看到 module maas-platform 和 go 1.21
   ```

2. **验证编译**
   ```bash
   make build
   # 应该生成 bin/api-gateway 和 bin/model-registry
   ```

3. **验证运行**
   ```bash
   # 终端1
   make run-api
   # 输出: MaaS API Gateway Starting...

   # 终端2
   make run-model
   # 输出: MaaS Model Registry Starting...
   ```

4. **验证Git**
   ```bash
   git log --oneline
   # 应该看到提交记录
   ```

---

## 参考资源

- [Go Modules官方文档](https://golang.org/ref/mod)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Conventional Commits规范](https://www.conventionalcommits.org/)
- [微服务设计模式](https://microservices.io/patterns/index.html)