# 节点1.2：API Gateway核心

> 📅 **学习时间**：5天  
> 🎯 **目标**：使用Gin框架实现API Gateway核心功能

## 本节你将学到

1. Gin框架核心概念和用法
2. RESTful API设计规范
3. 中间件开发（Logger、Recovery、CORS、RequestID）
4. 健康检查接口实现
5. 优雅关闭HTTP服务器

---

## 技术详解

### 1. Gin框架简介

Gin是Go语言中速度最快的Web框架之一，使用httprouter作为路由引擎。

**核心特性**：
- 高性能：使用radix tree路由算法
- 中间件支持：可以定义全局和分组中间件
- 错误管理：集中处理HTTP错误
- JSON验证：内置JSON数据验证
- 路由分组：支持路由分组和嵌套

**基本使用**：
```go
r := gin.Default()  // 默认带有Logger和Recovery中间件

r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
})

r.Run()  // 默认监听 :8080
```

### 2. 中间件（Middleware）

中间件是在请求处理前/后执行的函数，用于实现日志、认证、错误恢复等功能。

**中间件执行顺序**：
```
请求 → 中间件1 → 中间件2 → Handler → 中间件2 → 中间件1 → 响应
```

**常用中间件**：

**Recovery中间件**：捕获panic，防止服务器崩溃
```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录错误并返回500
                c.AbortWithStatusJSON(500, gin.H{
                    "error": "Internal server error",
                })
            }
        }()
        c.Next()
    }
}
```

**Logger中间件**：记录请求信息
```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()  // 执行后续处理器
        
        latency := time.Since(start)
        log.Printf("%s %s %d %v", 
            c.Request.Method,
            path,
            c.Writer.Status(),
            latency,
        )
    }
}
```

**CORS中间件**：处理跨域请求
```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

### 3. 优雅关闭

生产环境需要优雅关闭服务器，确保正在处理的请求完成后再退出。

```go
srv := &http.Server{
    Addr:    ":8080",
    Handler: r,
}

// 启动优雅关闭监听
go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    srv.Shutdown(ctx)  // 优雅关闭
}()

srv.ListenAndServe()
```

### 4. 配置管理（Viper）

Viper支持多种配置来源：配置文件、环境变量、命令行参数。

**配置优先级**（高→低）：
1. 显式调用Set
2. 命令行参数
3. 环境变量
4. 配置文件
5. 默认值

```go
viper.SetDefault("port", 8080)
viper.SetEnvPrefix("MAAS")  // 环境变量前缀
viper.AutomaticEnv()        // 自动读取环境变量

viper.SetConfigName("config")
viper.AddConfigPath(".")
viper.ReadInConfig()

var config Config
viper.Unmarshal(&config)
```

### 5. 结构化日志（Zap）

Zap是Uber开源的高性能日志库，比标准库快10倍以上。

**日志级别**：
- Debug：调试信息
- Info：一般信息
- Warn：警告信息
- Error：错误信息
- Fatal：致命错误（会退出程序）

```go
logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("request completed",
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Int("status", 200),
    zap.Duration("latency", 100*time.Millisecond),
)
```

---

## 实操任务

### 任务1：更新main.go实现HTTP服务器

参考上方代码，实现：
1. 使用gin.New()创建路由
2. 添加中间件链
3. 实现健康检查接口
4. 添加优雅关闭逻辑

### 任务2：实现中间件

创建以下中间件：
1. Recovery：捕获panic
2. Logger：记录请求日志
3. CORS：处理跨域
4. RequestID：生成请求ID

### 任务3：更新配置管理

使用Viper实现：
1. 从配置文件读取配置
2. 支持环境变量覆盖
3. 设置合理的默认值

### 任务4：实现Logger封装

封装Zap日志库：
1. 支持日志级别配置
2. JSON格式输出
3. 包含调用信息

### 任务5：实现Handler和Router

创建RESTful API：
1. 用户认证接口（/auth/login, /auth/register）
2. 模型管理接口（/models）
3. 推理接口（/inference）
4. 健康检查（/health）

---

## 代码变更记录

### 提交信息
```
feat(phase1/node1.2): implement API Gateway core with Gin

- Add Gin framework integration
- Implement middleware chain (Recovery, Logger, CORS, RequestID)
- Add health check endpoint
- Implement graceful shutdown
- Add Viper configuration management
- Integrate Zap structured logging
- Create RESTful API handlers (auth, models, inference)
- Setup router with route groups
```

### 修改的文件清单

#### 1. api-gateway/cmd/main.go
从占位文件更新为完整的HTTP服务器：
- 集成Gin框架
- 添加中间件链
- 实现健康检查
- 添加优雅关闭
- 集成Swagger文档注解

**主要变化**：
- 添加Gin框架初始化
- 配置全局中间件
- 实现HTTP服务器和优雅关闭
- 添加Swagger注解（为后续文档生成准备）

#### 2. api-gateway/internal/config/config.go
从简单配置结构更新为使用Viper：
- 添加Viper配置加载
- 支持配置文件和环境变量
- 设置默认值
- 添加完整的数据库、Redis、JWT等配置项

**主要变化**：
- 引入viper库
- 实现Load()函数
- 添加配置默认值
- 支持多来源配置（文件、环境变量）

#### 3. api-gateway/pkg/logger/logger.go
从占位文件更新为完整的Zap日志实现：
- 封装zap.Logger
- 支持日志级别配置
- JSON格式输出
- 添加快捷方法（Info, Error, Fatal等）

**主要变化**：
- 引入zap库
- 实现Logger结构体
- 配置JSON编码器
- 添加日志级别支持

#### 4. api-gateway/internal/middleware/middleware.go
从占位文件更新为完整的中间件实现：
- Recovery：捕获panic并记录堆栈
- Logger：记录请求信息和延迟
- CORS：处理跨域请求
- RequestID：生成和传递请求ID

**主要变化**：
- 实现4个核心中间件
- 集成logger进行日志记录
- 使用uuid生成请求ID

#### 5. api-gateway/internal/handler/handler.go
从占位文件更新为基础Handler结构：
- 添加Config和Logger依赖
- 实现统一的响应方法（Success, Error等）
- 添加HTTP状态码快捷方法

**主要变化**：
- 添加Handler结构体
- 实现响应封装
- 添加快捷错误处理方法

#### 6. api-gateway/internal/handler/user_handler.go
**新增文件**
实现用户相关接口：
- Login：用户登录
- Register：用户注册
- GetCurrentUser：获取当前用户

**代码结构**：
- 定义请求/响应结构体
- 实现JWT Token返回（占位）
- 参数校验（binding tags）

#### 7. api-gateway/internal/handler/model_handler.go
**新增文件**
实现模型管理接口：
- CreateModel：创建模型
- ListModels：列出模型
- GetModel：获取模型详情
- DeleteModel：删除模型

**代码结构**：
- 定义ModelRequest/ModelResponse
- 实现CRUD接口（占位）
- 查询参数处理

#### 8. api-gateway/internal/handler/inference_handler.go
**新增文件**
实现推理接口：
- RunInference：执行模型推理

**代码结构**：
- 定义InferenceRequest/InferenceResponse
- 返回推理结果（占位）

#### 9. api-gateway/internal/router/router.go
从占位文件更新为完整的路由配置：
- 注册认证路由（/auth）
- 注册用户路由（/users）
- 注册模型路由（/models）
- 注册推理路由（/inference）

**路由结构**：
```
/api/v1/
├── /auth
│   ├── POST /login
│   └── POST /register
├── /users
│   └── GET /me
├── /models
│   ├── POST /
│   ├── GET /
│   ├── GET /:id
│   └── DELETE /:id
└── /inference
    └── POST /
```

---

## 验证步骤

### 1. 编译验证
```bash
cd MaaS-go
make build
# 应该成功编译出 bin/api-gateway
```

### 2. 运行验证
```bash
# 终端1
make run-api
# 输出：
# {"level":"info","msg":"Starting API Gateway",...}
# {"level":"info","msg":"Server starting","addr":":8080"}

# 终端2 - 测试健康检查
curl http://localhost:8080/health
# 输出：{"status":"ok","service":"api-gateway","timestamp":...}

# 测试登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
# 输出：{"code":0,"message":"success","data":{...}}
```

### 3. 日志验证
```bash
# 查看日志输出，应该包含：
# - 启动日志（JSON格式）
# - 请求日志（方法、路径、状态码、延迟）
# - 请求ID（X-Request-ID header）
```

---

## 检查清单

完成本节点后，请确认：

- [ ] Gin框架正确集成
- [ ] 服务能启动并监听端口
- [ ] /health接口返回200和JSON
- [ ] 中间件正常工作（日志、CORS、RequestID）
- [ ] 配置管理支持环境变量
- [ ] 日志输出为JSON格式
- [ ] API路由正确注册
- [ ] 优雅关闭正常工作（Ctrl+C能正常退出）

---

## 下一步

完成本节点后，API Gateway的核心功能已经完成。接下来进入：

**节点1.3：配置管理体系** → [继续学习](./node-1-3.md)

在那里你将：
- 完善配置文件管理
- 实现配置热更新
- 添加配置验证
- 处理敏感信息

---

## 参考资源

- [Gin官方文档](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Viper配置库](https://github.com/spf13/viper)
- [Zap日志库](https://github.com/uber-go/zap)
