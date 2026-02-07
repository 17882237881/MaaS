# 节点1.3：配置管理体系

> 📅 **学习时间**：3天  
> 🎯 **目标**：实现完善的配置管理体系，支持多环境、配置验证、热更新

## 本节你将学到

1. Viper配置库的高级用法
2. 多环境配置管理（dev/staging/prod）
3. 配置验证和默认值
4. 配置文件热更新（Hot Reload）
5. 敏感信息处理（环境变量）

---

## 技术详解

### 1. 为什么需要配置管理？

不同环境需要不同的配置：
- **开发环境**：本地数据库、调试模式、详细日志
- **测试环境**：测试数据库、独立服务实例
- **生产环境**：集群数据库、生产级安全、监控告警

**配置管理目标**：
- 代码与配置分离（12-Factor App原则）
- 支持多环境部署
- 配置变更无需重新编译
- 敏感信息不泄露（密码、密钥）

### 2. Viper配置库详解

**Viper**是Go语言最流行的配置库，支持：
- JSON、YAML、TOML、HCL、envfile、Java properties
- 实时重新加载
- 从环境变量、命令行、远程配置系统读取
- 默认值设置

**配置优先级**（高→低）：
1. 显式调用Set
2. 命令行参数
3. 环境变量
4. 配置文件
5. 默认值

### 3. 配置结构设计

**分层配置**：
```yaml
# 基础配置
environment: development
port: 8080
log_level: info

# 数据库配置（嵌套结构）
database:
  host: localhost
  port: 5432
  name: maas_platform

# 服务配置（下游服务地址）
services:
  model_registry: http://localhost:8081
  inference: http://localhost:8082
```

**Go结构体映射**：
```go
type Config struct {
    Environment string          `mapstructure:"environment"`
    Database    DatabaseConfig  `mapstructure:"database"`
    Services    ServiceConfig   `mapstructure:"services"`
}

type DatabaseConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
```

### 4. 配置验证

**为什么需要验证？**
- 防止启动后因配置错误而崩溃
- 提前发现配置问题（如缺失必填项）
- 确保生产环境安全（如强密码）

**验证内容**：
- 必填项检查
- 数值范围检查
- 格式验证
- 生产环境特殊规则

### 5. 热更新（Hot Reload）

**什么是热更新？**
配置文件修改后，应用程序自动重新加载配置，无需重启。

**实现原理**：
```go
// 使用fsnotify监听文件变化
watcher, _ := fsnotify.NewWatcher()
watcher.Add("config.yaml")

// 文件修改时重新加载
for event := range watcher.Events {
    if event.Op&fsnotify.Write == fsnotify.Write {
        viper.ReadInConfig()  // 重新读取
        viper.Unmarshal(&config)  // 重新解析
    }
}
```

**适用场景**：
- 日志级别调整
- 限流阈值调整
- 开关功能切换

**不适用场景**：
- 数据库连接信息（需要重启连接池）
- 端口变更（需要重启监听）
- 安全配置变更（建议重启）

---

## 实操任务

### 任务1：创建配置文件

创建 `api-gateway/config/config.yaml`：

```yaml
environment: development
port: 8080
log_level: info

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: maas_platform
  ssl_mode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key
  expires_in: 86400

services:
  model_registry: http://localhost:8081
  inference: http://localhost:8082
  user_center: http://localhost:8083
  billing: http://localhost:8084

rate_limit:
  enabled: true
  rpm: 1000
  burst: 100
```

### 任务2：更新配置加载逻辑

增强 `config.go`：
1. 定义完整的配置结构体
2. 添加配置验证方法
3. 实现配置热更新
4. 添加辅助方法（DSN生成等）

### 任务3：环境变量覆盖

测试环境变量覆盖：
```bash
export MAAS_ENVIRONMENT=production
export MAAS_PORT=9090
export MAAS_DATABASE_PASSWORD=secret123
```

### 任务4：更新main.go

使用新的配置加载方式：
```go
cfg, err := config.LoadWithWatch(func(newCfg *config.Config) {
    log.Info("Config reloaded")
})
```

---

## 代码变更记录

### 提交信息
```
feat(phase1/node1.3): implement configuration management system

- Add comprehensive config.yaml with all settings
- Update config.go with full configuration management
  - Add nested configuration structures
  - Implement configuration validation
  - Add hot reload support with fsnotify
  - Add helper methods (DatabaseDSN, RedisAddr)
  - Add environment detection helpers
- Update main.go to use new config system
  - Add essential config validation
  - Add config summary logging
  - Add /config endpoint (dev only)
- Add detailed logging for config loading
```

### 创建的文件

#### 1. api-gateway/config/config.yaml
**新增文件**
完整的YAML配置文件，包含：
- 基础配置（环境、端口、日志级别）
- 数据库配置（PostgreSQL连接信息）
- Redis配置
- JWT配置
- 下游服务地址
- 限流配置

**特点**：
- 清晰的注释说明
- 合理的默认值
- 支持环境变量覆盖（MAAS_前缀）

#### 2. api-gateway/internal/config/config.go
**大幅更新**
从简单配置更新为完整的配置管理系统：

**主要变化**：
1. **结构体优化**：
   ```go
   type Config struct {
       Environment string          `mapstructure:"environment"`
       Port        int             `mapstructure:"port"`
       Database    DatabaseConfig  `mapstructure:"database"`
       // ... 更多嵌套结构
   }
   ```

2. **增强的Load函数**：
   - 支持配置文件多路径搜索
   - 环境变量自动映射
   - 合理的默认值设置
   - 配置验证

3. **热更新支持**：
   ```go
   func LoadWithWatch(onChange func(*Config)) (*Config, error)
   ```
   - 使用fsnotify监听文件变化
   - 自动重新加载配置
   - 回调通知配置变更

4. **配置验证**：
   ```go
   func (c *Config) Validate() error
   ```
   - 环境值验证（development/staging/production）
   - 端口范围验证（1-65535）
   - 日志级别验证
   - 数据库必填项验证
   - 生产环境安全检查（JWT密钥长度）

5. **辅助方法**：
   ```go
   func (c *Config) IsDevelopment() bool
   func (c *Config) IsProduction() bool
   func (c *Config) DatabaseDSN() string
   func (c *Config) RedisAddr() string
   ```

### 修改的文件

#### api-gateway/cmd/main.go
**更新内容**：
1. 使用新的配置加载方式：
   ```go
   cfg, err := config.LoadWithWatch(func(newCfg *config.Config) {
       fmt.Println("Configuration reloaded")
   })
   ```

2. 添加配置验证：
   ```go
   if err := validateEssentialConfig(cfg); err != nil {
       fmt.Printf("Configuration validation failed: %v\n", err)
       os.Exit(1)
   }
   ```

3. 添加配置信息端点（仅开发环境）：
   ```go
   if cfg.IsDevelopment() {
       r.GET("/config", func(c *gin.Context) {
           c.JSON(200, gin.H{
               "environment": cfg.Environment,
               "port": cfg.Port,
               // ...
           })
       })
   }
   ```

4. 添加配置摘要日志：
   ```go
   func printConfigSummary(cfg *config.Config, log *logger.Logger)
   ```
   - 输出关键配置信息
   - 便于调试和验证

---

## 验证步骤

### 1. 配置文件加载验证

```bash
# 1. 启动服务（自动加载config.yaml）
cd api-gateway
go run cmd/main.go

# 应该看到：
# {"level":"info","msg":"Configuration loaded",...}
# {"level":"info","msg":"Database configuration",...}
# {"level":"info","msg":"Redis configuration",...}
```

### 2. 环境变量覆盖验证

```bash
# 设置环境变量
export MAAS_ENVIRONMENT=staging
export MAAS_PORT=9090
export MAAS_LOG_LEVEL=debug

# 启动服务
go run cmd/main.go

# 验证端口变为9090，日志级别变为debug
```

### 3. 配置端点验证（开发环境）

```bash
# 访问配置信息端点
curl http://localhost:8080/config

# 应该返回：
# {
#   "environment": "development",
#   "port": 8080,
#   "log_level": "info",
#   "database": {...},
#   "redis": {...}
# }
```

### 4. 热更新验证

```bash
# 1. 启动服务
go run cmd/main.go

# 2. 修改 config/config.yaml
# 修改 log_level 从 info 改为 debug

# 3. 应该看到控制台输出：
# Config file modified, reloading...
# Config reloaded successfully
```

### 5. 配置验证失败测试

```bash
# 测试必填项缺失
export MAAS_ENVIRONMENT=production
export MAAS_JWT_SECRET=short  # 太短

# 启动应该失败并提示：
# Configuration validation failed: JWT secret must be at least 32 characters
```

---

## 检查清单

完成本节点后，请确认：

- [ ] config.yaml文件结构清晰
- [ ] 服务能加载YAML配置
- [ ] 环境变量能覆盖配置文件
- [ ] 配置验证正常工作（必填项、范围、格式）
- [ ] 热更新功能正常（修改文件自动加载）
- [ ] 辅助方法工作正常（DatabaseDSN、RedisAddr）
- [ ] 生产环境安全检查生效
- [ ] /config端点只在开发环境可用

---

## 下一步

完成本节点后，你已经实现了完善的配置管理体系。接下来进入：

**节点1.4：日志与监控基础** → [继续学习](./node-1-4.md)

在那里你将：
- 完善Zap日志系统
- 添加日志轮转
- 集成Prometheus指标
- 实现健康检查指标

---

## 参考资源

- [Viper官方文档](https://github.com/spf13/viper)
- [fsnotify文件监控](https://github.com/fsnotify/fsnotify)
- [12-Factor App: Config](https://12factor.net/config)
- [Go配置最佳实践](https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64)
