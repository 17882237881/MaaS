# 阶段1：基础架构搭建

> 🏗️ **阶段目标**：搭建MaaS平台的基础架构，实现API Gateway和Model Registry两个核心服务

## 阶段总览

### 学习时间
**4周**（约20个工作日）

### 核心目标
1. 掌握Go语言企业级项目结构组织
2. 学会使用Gin框架开发RESTful API
3. 掌握GORM进行数据库设计和操作
4. 学会Docker容器化部署
5. 理解微服务拆分的基本原则

### 最终产出
- ✅ 可独立运行的API Gateway服务
- ✅ 可独立运行的Model Registry服务
- ✅ 支持基础CRUD操作
- ✅ Docker Compose一键启动所有服务

## 技术栈详解

### 1. Go语言（1.21+）
**为什么要用这个版本？**
- Go 1.21引入了内置的`slog`日志库（虽然我们用Zap）
- 改进了泛型的使用体验
- 性能持续优化

**核心特性回顾**：
```go
// 结构体和方法
type User struct {
    ID   string
    Name string
}

func (u *User) GetName() string {
    return u.Name
}

// 接口
type Service interface {
    DoSomething() error
}

// Goroutine和Channel（并发）
go func() {
    // 异步执行任务
}()

ch := make(chan string)
ch <- "message"  // 发送
msg := <-ch      // 接收
```

### 2. Gin Web框架
**什么是Gin？**
Gin是Go语言中速度最快的Web框架之一，类似于Python的Flask或Node.js的Express。

**核心概念**：
- **路由（Router）**：URL路径和处理函数的映射
- **中间件（Middleware）**：请求处理链，可插入日志、认证等功能
- **上下文（Context）**：包含请求和响应的所有信息
- **参数绑定**：自动将JSON/表单数据绑定到结构体

**简单示例**：
```go
r := gin.Default()

// 定义路由
r.GET("/hello", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Hello World",
    })
})

// 带参数的路由
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})

// POST请求+JSON绑定
type UserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

r.POST("/users", func(c *gin.Context) {
    var req UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // 处理请求...
})
```

### 3. GORM（Go ORM库）
**什么是ORM？**
ORM（Object-Relational Mapping）对象关系映射，让你用操作对象的方式操作数据库，不用写SQL。

**GORM的核心功能**：
- 自动建表（Auto Migration）
- 链式查询
- 关联（一对一、一对多、多对多）
- 钩子（BeforeCreate, AfterUpdate等）
- 事务支持

**示例**：
```go
// 定义模型
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"size:255"`
    Email     string `gorm:"uniqueIndex"`
    CreatedAt time.Time
}

// 连接数据库
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 自动迁移（创建表）
db.AutoMigrate(&User{})

// 创建
user := User{Name: "John", Email: "john@example.com"}
db.Create(&user)

// 查询
var result User
db.First(&result, 1)           // 按主键查找
db.First(&result, "name = ?", "John")  // 按条件查找

// 更新
db.Model(&user).Update("name", "Jane")

// 删除
db.Delete(&user)
```

### 4. PostgreSQL
**为什么选择PostgreSQL？**
- 功能最强大的开源关系型数据库
- 支持JSON、数组、GIS等高级类型
- 事务完整性和并发控制优秀
- 云原生友好（各大云厂商都支持）

**关键概念**：
- **Schema**：数据库中的命名空间，用于组织表
- **索引（Index）**：加速查询的数据结构
- **事务（Transaction）**：保证数据一致性的操作单元
- **连接池**：复用数据库连接，提高性能

### 5. Zap日志库
**为什么要用Zap？**
- 性能极高（比标准库快10倍以上）
- 支持结构化日志（JSON格式）
- 支持日志级别（Debug/Info/Error）
- 可以对接ELK等日志收集系统

**日志级别说明**：
- **Debug**：开发调试信息（如变量值）
- **Info**：正常流程信息（如请求处理完成）
- **Warn**：警告信息（如性能下降）
- **Error**：错误信息（如数据库连接失败）
- **Fatal**：致命错误（程序无法继续）

### 6. Docker
**什么是Docker？**
Docker是一种容器化技术，把应用和依赖打包成一个"集装箱"，在任何地方都能一致运行。

**核心概念**：
- **镜像（Image）**：只读的模板，包含运行应用所需的一切
- **容器（Container）**：镜像的运行实例
- **Dockerfile**：定义如何构建镜像的脚本
- **Docker Compose**：定义和运行多容器应用的工具

**Dockerfile示例**：
```dockerfile
# 基础镜像
FROM golang:1.21-alpine

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./main"]
```

## 节点详解

### 节点1.1：项目初始化与架构设计（3天）

**学习目标**：
- 理解Go Modules工作区管理
- 掌握微服务目录结构组织
- 学会依赖版本管理
- 建立Git工作流规范

**技术介绍**：

**1. Go Modules**
Go Modules是Go 1.11+引入的依赖管理工具，类似于Node.js的npm、Python的pip。

**关键命令**：
```bash
# 初始化模块
go mod init maas-platform

# 下载依赖
go mod download

# 整理并下载缺失的依赖
go mod tidy

# 查看依赖关系
go mod graph
```

**2. 微服务目录结构**
推荐的目录组织方式：
```
project-root/
├── api-gateway/              # 服务目录
│   ├── cmd/                  # 程序入口
│   │   └── main.go          # main函数
│   ├── internal/            # 私有代码
│   │   ├── config/         # 配置
│   │   ├── handler/        # HTTP处理器
│   │   ├── middleware/     # 中间件
│   │   ├── model/          # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── router/         # 路由定义
│   │   └── service/        # 业务逻辑层
│   └── pkg/                 # 公共库
│       ├── logger/         # 日志工具
│       └── utils/          # 工具函数
├── model-registry/          # 另一个服务
├── shared/                  # 共享代码
│   ├── proto/              # Protocol Buffers定义
│   └── errors/             # 公共错误定义
├── deploy/                  # 部署配置
│   ├── docker/             # Docker配置
│   └── k8s/                # Kubernetes配置
└── docs/                    # 文档
```

**分层架构说明**：
- **Handler层**：处理HTTP请求，参数校验，调用Service
- **Service层**：业务逻辑，事务管理
- **Repository层**：数据访问，数据库操作
- **Model层**：数据结构定义

**为什么这样分层？**
- **单一职责**：每层只负责一件事
- **可测试性**：每层可以独立测试
- **可替换性**：可以替换某层实现（如换数据库）

**实操任务**：
1. 创建项目根目录和go.mod
2. 创建api-gateway和model-registry目录结构
3. 创建基本的main.go文件（能打印Hello World）
4. 初始化Git仓库，提交第一次commit
5. 编写README.md说明项目结构

**检查点**：
- [ ] 项目目录结构完整
- [ ] `go mod init`成功
- [ ] 每个服务能独立编译运行
- [ ] Git仓库初始化完成

---

### 节点1.2：API Gateway核心（5天）

**学习目标**：
- 掌握Gin框架核心用法
- 实现全局异常处理
- 学会使用中间件
- 生成Swagger API文档

**技术介绍**：

**1. HTTP基础回顾**
- **请求方法**：GET（获取）、POST（创建）、PUT（更新）、DELETE（删除）
- **状态码**：200（成功）、400（请求错误）、401（未授权）、500（服务器错误）
- **Header**：Content-Type、Authorization等
- **Body**：请求/响应的数据体

**2. RESTful API设计规范**
```
GET    /api/v1/users          # 获取用户列表
GET    /api/v1/users/:id      # 获取单个用户
POST   /api/v1/users          # 创建用户
PUT    /api/v1/users/:id      # 更新用户
DELETE /api/v1/users/:id      # 删除用户
```

**3. 中间件（Middleware）**
中间件是处理HTTP请求的"钩子"，可以在请求处理前/后执行代码。

**中间件执行顺序**：
```
请求 → 中间件1 → 中间件2 → Handler → 中间件2 → 中间件1 → 响应
```

**常见中间件类型**：
- **日志中间件**：记录请求信息
- **恢复中间件**：捕获panic，防止程序崩溃
- **CORS中间件**：处理跨域请求
- **认证中间件**：验证用户身份

**4. Swagger/OpenAPI**
Swagger是API文档规范，可以从代码注释自动生成可交互的API文档。

**实操任务**：
1. 创建Gin路由器
2. 实现基础中间件（Logger、Recovery、CORS）
3. 创建健康检查接口（GET /health）
4. 集成Swagger，生成API文档
5. 实现简单的CRUD接口示例

**代码结构**：
```go
// main.go - 程序入口
func main() {
    r := gin.Default()
    
    // 全局中间件
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    r.Use(middleware.CORS())
    
    // 路由注册
    router.RegisterRoutes(r)
    
    r.Run(":8080")
}

// middleware/logger.go - 日志中间件
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()  // 执行后续处理器
        duration := time.Since(start)
        
        log.Printf("%s %s %d %v", 
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            duration,
        )
    }
}

// router/router.go - 路由定义
func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    {
        api.GET("/health", handler.HealthCheck)
        api.GET("/users", handler.ListUsers)
        api.POST("/users", handler.CreateUser)
    }
}
```

**检查点**：
- [ ] 服务能启动并监听端口
- [ ] /health接口返回200
- [ ] 日志正确输出请求信息
- [ ] Swagger文档可访问
- [ ] 异常处理正常工作（panic不会导致服务崩溃）

---

### 节点1.3：配置管理体系（3天）

**学习目标**：
- 学会多环境配置管理
- 掌握Viper库使用
- 理解配置热更新
- 学会敏感信息处理

**技术介绍**：

**1. 为什么需要配置管理？**
不同环境（开发/测试/生产）需要不同的配置：
- 数据库连接信息
- API密钥
- 日志级别
- 服务端口号

**2. 配置来源优先级**（高→低）：
1. 命令行参数
2. 环境变量
3. 配置文件
4. 默认值

**3. Viper配置库**
Viper支持多种配置格式：JSON、YAML、TOML、HCL、envfile、Java properties

**示例配置**（config.yaml）：
```yaml
environment: development
port: 8080
log_level: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: secret
  name: maas_platform

redis:
  host: localhost
  port: 6379
```

**Go结构体映射**：
```go
type Config struct {
    Environment string `mapstructure:"environment"`
    Port        int    `mapstructure:"port"`
    LogLevel    string `mapstructure:"log_level"`
    
    Database struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        User     string `mapstructure:"user"`
        Password string `mapstructure:"password"`
        Name     string `mapstructure:"name"`
    } `mapstructure:"database"`
}
```

**4. 环境变量使用**
```bash
# .env文件
DB_HOST=localhost
DB_PORT=5432

# 代码中读取
host := os.Getenv("DB_HOST")
```

**实操任务**：
1. 创建config.yaml配置文件
2. 使用Viper加载配置
3. 支持环境变量覆盖
4. 添加配置验证
5. 实现配置热更新（可选）

**检查点**：
- [ ] 配置能从YAML文件加载
- [ ] 环境变量能覆盖配置文件
- [ ] 配置验证正常工作
- [ ] 不同环境使用不同配置

---

### 节点1.4：日志与监控基础（4天）

**学习目标**：
- 掌握Zap日志库
- 实现结构化日志
- 添加基础Metrics
- 理解可观测性三支柱

**技术介绍**：

**1. 可观测性三支柱**
- **日志（Logging）**：记录离散事件（如"用户登录"）
- **指标（Metrics）**：记录聚合数据（如"QPS=100"）
- **追踪（Tracing）**：记录请求链路（如"API→Service→DB"）

**2. 结构化日志 vs 文本日志**
```go
// 文本日志（难解析）
log.Printf("User %s logged in from %s at %s", username, ip, time)

// 结构化日志（JSON，易解析）
logger.Info("user logged in",
    zap.String("username", username),
    zap.String("ip", ip),
    zap.Time("timestamp", time),
)
```

**3. Zap使用示例**
```go
// 创建logger
config := zap.NewProductionConfig()
config.EncoderConfig.TimeKey = "timestamp"
config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

logger, _ := config.Build()
defer logger.Sync()

// 使用
logger.Info("request processed",
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Int("status", 200),
    zap.Duration("latency", 100*time.Millisecond),
)

logger.Error("database connection failed",
    zap.Error(err),
    zap.String("host", "localhost"),
)
```

**4. 基础Metrics**
使用Prometheus客户端库记录指标：
```go
// 计数器（只增不减）
requestCounter := prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total HTTP requests",
    },
    []string{"method", "path", "status"},
)

// 直方图（记录分布）
requestDuration := prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "HTTP request duration",
        Buckets: []float64{0.1, 0.5, 1, 2, 5},
    },
    []string{"method", "path"},
)
```

**实操任务**：
1. 集成Zap日志库
2. 替换所有fmt.Print为结构化日志
3. 添加请求日志中间件
4. 集成Prometheus metrics
5. 创建/metrics接口

**检查点**：
- [ ] 日志输出为JSON格式
- [ ] 包含request_id追踪
- [ ] /metrics接口可访问
- [ ] 能记录请求延迟分布

---

### 节点1.5：数据库层设计（5天）

**学习目标**：
- 掌握数据库设计原则
- 学会GORM高级用法
- 实现Repository模式
- 掌握数据库迁移

**技术介绍**：

**1. 数据库设计原则**

**第一范式（1NF）**：每个字段都是原子性的
```sql
-- 错误：hobbies字段包含多个值
users (id, name, hobbies)
-- "reading,swimming,gaming"

-- 正确：拆分为单独表
users (id, name)
user_hobbies (user_id, hobby)
```

**第二范式（2NF）**：非主键字段完全依赖于主键
```sql
-- 错误：category_name只依赖于category_id
products (id, name, category_id, category_name)

-- 正确：拆分为两个表
products (id, name, category_id)
categories (id, name)
```

**索引设计原则**：
- 主键自动创建索引
- 外键通常需要索引
- 频繁查询的字段加索引
- 不要给所有字段都加索引（影响写入性能）

**2. Repository模式**
Repository模式是数据访问层的设计模式，将数据访问逻辑封装起来。

**优势**：
- 业务逻辑与数据访问解耦
- 易于测试（可Mock Repository）
- 易于切换数据库实现

**结构**：
```go
// 接口定义
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter UserFilter) ([]*User, error)
}

// 实现
type GormUserRepository struct {
    db *gorm.DB
}

func (r *GormUserRepository) Create(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

**3. GORM高级特性**

**关联查询**：
```go
// 定义模型
type User struct {
    ID      string
    Name    string
    Orders  []Order  // 一对多
    Profile Profile  // 一对一
    Roles   []Role   `gorm:"many2many:user_roles;"` // 多对多
}

// 预加载（解决N+1问题）
var users []User
db.Preload("Orders").Preload("Profile").Find(&users)
```

**事务**：
```go
// 方式1：闭包（自动提交/回滚）
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // 自动回滚
    }
    if err := tx.Create(&profile).Error; err != nil {
        return err  // 自动回滚
    }
    return nil  // 自动提交
})

// 方式2：手动控制
tx := db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

**4. 数据库迁移**
迁移是管理数据库Schema变更的方式，比手动执行SQL更可靠。

**使用golang-migrate**：
```bash
# 安装
migrate -version

# 创建迁移文件
migrate create -ext sql -dir migrations -seq create_users_table

# 生成文件：
# 001_create_users_table.up.sql
# 001_create_users_table.down.sql

# 执行迁移
migrate -path migrations -database "postgres://user:pass@localhost/db" up
```

**实操任务**：
1. 设计Model Registry的数据库表结构
2. 创建GORM模型定义
3. 实现Repository接口
4. 编写数据库迁移脚本
5. 实现Service层调用Repository
6. 编写单元测试

**数据库设计（Model Registry）**：
```sql
-- models表：模型基本信息
CREATE TABLE models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) NOT NULL,
    framework VARCHAR(50) NOT NULL,  -- pytorch/tensorflow/onnx
    status VARCHAR(50) DEFAULT 'pending',  -- pending/ready/failed
    size_bytes BIGINT DEFAULT 0,
    checksum VARCHAR(64),
    storage_path VARCHAR(512),
    owner_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(name, version)
);

-- tags表：模型标签
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- model_tags表：模型和标签的多对多关系
CREATE TABLE model_tags (
    model_id UUID REFERENCES models(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (model_id, tag_id)
);

-- model_metadata表：模型元数据（键值对）
CREATE TABLE model_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_id UUID REFERENCES models(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**检查点**：
- [ ] 数据库表结构符合范式
- [ ] GORM模型定义正确
- [ ] Repository接口完整实现
- [ ] 迁移脚本可正常执行
- [ ] 基础CRUD接口可调用
- [ ] 有基本的单元测试

---

## 阶段1里程碑

### 完成检查清单

**API Gateway服务**：
- [ ] 可独立启动和运行
- [ ] 支持基础路由和中间件
- [ ] 有健康检查接口
- [ ] 日志输出正常
- [ ] 配置管理正常工作

**Model Registry服务**：
- [ ] 可独立启动和运行
- [ ] 数据库连接正常
- [ ] 支持模型的CRUD操作
- [ ] 有完整的Swagger文档

**基础设施**：
- [ ] Docker Compose可启动所有服务
- [ ] PostgreSQL容器正常运行
- [ ] 服务间可通过网络通信

### 可演示功能
1. 启动服务：`docker-compose up`
2. 访问API文档：http://localhost:8080/swagger/index.html
3. 创建模型：POST /api/v1/models
4. 查询模型列表：GET /api/v1/models
5. 查看日志输出

### 下一步
完成阶段1后，你将掌握：
- ✅ Go语言项目结构组织
- ✅ Gin框架开发RESTful API
- ✅ GORM数据库操作
- ✅ Docker容器化

准备好进入**阶段2：核心功能开发**了吗？

---

## 常见问题

### Q: 为什么用PostgreSQL而不是MySQL？
**A**: PostgreSQL功能更强大，支持JSON、数组、GIS等类型，事务和并发控制更好。两者都是优秀的数据库，实际工作中根据团队熟悉度和云厂商支持选择。

### Q: 为什么用GORM而不是原生SQL？
**A**: GORM提高开发效率，自动处理很多重复工作。但对于复杂查询，GORM支持原生SQL。实际项目中两者可以结合使用。

### Q: 微服务拆分的粒度怎么把握？
**A**: 没有绝对标准，一般原则：
- 独立部署（一个服务出问题不影响其他）
- 独立团队维护
- 业务边界清晰
- 初期可以粗粒度，逐步拆分

### Q: 如何调试Go程序？
**A**: 
1. 使用VSCode + Go插件，支持断点调试
2. 使用日志打印（zap或fmt）
3. 使用Delve命令行调试器：`dlv debug`

---

**下一步**：开始节点1.1的学习 → [节点1.1：项目初始化与架构设计](./node-1-1.md)