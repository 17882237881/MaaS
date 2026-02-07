# 节点1.4：日志与监控基础

> 📅 **学习时间**：4天  
> 🎯 **目标**：实现结构化日志系统和基础监控指标采集

## 本节你将学到

1. Zap日志库的进阶用法
2. 日志文件轮转（Rotation）
3. Prometheus指标采集
4. HTTP请求指标埋点
5. 健康检查和状态监控

---

## 技术详解

### 1. 为什么需要结构化日志？

**传统文本日志**：
```
2024-01-01 10:00:00 User john logged in from 192.168.1.1
2024-01-01 10:00:05 Request GET /api/users took 100ms
```

**问题**：
- 难以解析和查询
- 无法聚合统计
- 格式不统一

**结构化日志（JSON）**：
```json
{
  "timestamp": "2024-01-01T10:00:00Z",
  "level": "info",
  "msg": "User login",
  "user": "john",
  "ip": "192.168.1.1"
}
```

**优势**：
- 易于机器解析
- 支持日志聚合系统（ELK、Loki）
- 可索引和搜索
- 统一格式

### 2. 日志轮转（Log Rotation）

**为什么需要轮转？**
- 防止磁盘空间耗尽
- 便于归档和清理
- 提高日志查询效率

**轮转策略**：
- **按大小**：单个文件达到100MB时创建新文件
- **按时间**：每天创建一个新文件
- **保留策略**：保留最近7天或最近10个文件
- **压缩**：轮转后的文件进行压缩

**实现库**：lumberjack
```go
&lumberjack.Logger{
    Filename:   "/var/log/app.log",
    MaxSize:    100, // MB
    MaxBackups: 10,
    MaxAge:     7,   // days
    Compress:   true,
}
```

### 3. 可观测性三支柱

**1. Metrics（指标）**：
- 聚合数据（QPS、延迟、错误率）
- 时序数据
- 适合监控和告警

**2. Logging（日志）**：
- 离散事件
- 详细信息
- 适合故障排查

**3. Tracing（追踪）**：
- 请求链路
- 分布式追踪
- 适合性能分析

### 4. Prometheus指标类型

**Counter（计数器）**：
- 只增不减
- 适合：请求总数、错误总数

**Gauge（仪表盘）**：
- 可增可减
- 适合：当前连接数、队列长度

**Histogram（直方图）**：
- 采样分布
- 适合：请求延迟、响应大小
- 自动计算分位数（P50、P95、P99）

**Summary（摘要）**：
- 类似Histogram，但计算滑动窗口分位数

### 5. HTTP指标采集

**关键指标**：
- 请求总数（按方法、路径、状态码）
- 请求延迟（P50、P95、P99）
- 请求/响应大小
- 活跃连接数
- 错误率

**埋点位置**：
- 中间件中统一采集
- 请求开始前记录开始时间
- 请求结束后计算延迟并记录

---

## 实操任务

### 任务1：完善Zap日志系统

更新 `api-gateway/pkg/logger/logger.go`：
- 添加配置结构体（级别、格式、输出方式）
- 实现文件输出和轮转
- 添加全局日志函数
- 支持结构化字段

### 任务2：创建Prometheus指标

创建 `api-gateway/pkg/metrics/metrics.go`：
- 定义HTTP指标（Counter、Histogram、Gauge）
- 创建Gin中间件自动采集
- 提供辅助函数

### 任务3：更新main.go集成监控

在main.go中：
- 添加/metrics端点
- 集成Prometheus中间件
- 设置服务状态指标

### 任务4：验证监控数据

```bash
# 启动服务
go run api-gateway/cmd/main.go

# 查看指标
curl http://localhost:8080/metrics

# 产生一些请求
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/models

# 再次查看指标，应该看到增长的计数器
curl http://localhost:8080/metrics | grep http_requests_total
```

---

## 代码变更记录

### 提交信息
```
feat(phase1/node1.4): implement logging and monitoring foundation

- Enhance Zap logger with file rotation and structured logging
- Add Prometheus metrics collection
- Implement HTTP request metrics middleware
- Add /metrics endpoint for Prometheus scraping
- Update main.go to integrate monitoring
```

### 修改的文件

#### 1. api-gateway/pkg/logger/logger.go
**大幅更新**
从简单的Logger封装更新为完整的日志系统：

**新增功能**：
1. **配置结构体**：
   ```go
   type Config struct {
       Level      string // debug, info, warn, error
       Format     string // json, console
       Output     string // stdout, file, both
       FilePath   string // log file path
       MaxSize    int    // megabytes
       MaxAge     int    // days
       MaxBackups int    // number of backups
       Compress   bool   // compress rotated files
   }
   ```

2. **文件轮转支持**：
   ```go
   func createFileSyncer(config Config) zapcore.WriteSyncer {
       lumberjackLogger := &lumberjack.Logger{
           Filename:   config.FilePath,
           MaxSize:    config.MaxSize,
           MaxBackups: config.MaxBackups,
           MaxAge:     config.MaxAge,
           Compress:   config.Compress,
       }
       return zapcore.AddSync(lumberjackLogger)
   }
   ```

3. **多输出支持**：
   - stdout：标准输出
   - file：文件输出（带轮转）
   - both：同时输出到stdout和file

4. **全局日志函数**：
   ```go
   func Debug(msg string, keysAndValues ...interface{})
   func Info(msg string, keysAndValues ...interface{})
   func Warn(msg string, keysAndValues ...interface{})
   func Error(msg string, keysAndValues ...interface{})
   func Fatal(msg string, keysAndValues ...interface{})
   ```

5. **辅助方法**：
   - `With()`：创建带字段的子logger
   - `WithContext()`：创建带request_id的logger
   - `Sync()`：刷新缓冲区

#### 2. api-gateway/pkg/metrics/metrics.go
**新增文件**
完整的Prometheus指标采集系统：

**定义的指标**：
1. **HTTPRequestDuration**：请求延迟直方图
   - Labels: method, path, status
   - Buckets: 1ms ~ 10s

2. **HTTPRequestTotal**：请求总数计数器
   - Labels: method, path, status

3. **HTTPRequestSize**：请求大小直方图
   - Labels: method, path

4. **HTTPResponseSize**：响应大小直方图
   - Labels: method, path

5. **ActiveConnections**：活跃连接数仪表盘

6. **ServiceUp**：服务状态（1=up, 0=down）

7. **ServiceInfo**：服务信息（版本、环境）

**中间件实现**：
```go
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        ActiveConnections.Inc()
        defer ActiveConnections.Dec()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        HTTPRequestDuration.WithLabelValues(...).Observe(duration)
        HTTPRequestTotal.WithLabelValues(...).Inc()
    }
}
```

#### 3. api-gateway/cmd/main.go
**更新内容**：
1. 添加Prometheus中间件：
   ```go
   r.Use(metrics.PrometheusMiddleware())
   ```

2. 添加/metrics端点：
   ```go
   r.GET("/metrics", metrics.Handler())
   ```

3. 设置服务指标：
   ```go
   metrics.SetServiceUp(true)
   metrics.SetServiceInfo("1.0.0", cfg.Environment)
   ```

---

## 验证步骤

### 1. 日志系统验证

```bash
# 1. 启动服务
cd api-gateway
go run cmd/main.go

# 2. 查看日志输出（JSON格式）
# 应该看到结构化的JSON日志

# 3. 产生一些日志
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/models

# 4. 检查日志文件（如果配置了文件输出）
ls -la logs/
cat logs/app.log
```

### 2. 指标采集验证

```bash
# 1. 查看基础指标
curl http://localhost:8080/metrics

# 2. 查看特定指标
curl -s http://localhost:8080/metrics | grep http_requests_total
curl -s http://localhost:8080/metrics | grep http_request_duration_seconds
curl -s http://localhost:8080/metrics | grep service_up

# 3. 产生请求并观察指标变化
for i in {1..10}; do
    curl -s http://localhost:8080/health > /dev/null
done

# 4. 再次查看计数器，应该增加了10
curl -s http://localhost:8080/metrics | grep 'http_requests_total{method="GET",path="/health"}'
```

### 3. 延迟分布验证

```bash
# 产生不同延迟的请求
# 快速请求
curl http://localhost:8080/health

# 慢速请求（模拟）
sleep 0.1 && curl http://localhost:8080/health
sleep 0.5 && curl http://localhost:8080/health

# 查看延迟分布
curl -s http://localhost:8080/metrics | grep http_request_duration_seconds_bucket
```

---

## 检查清单

完成本节点后，请确认：

- [ ] Zap日志输出JSON格式
- [ ] 支持多级别日志（Debug/Info/Warn/Error）
- [ ] 日志文件轮转正常工作
- [ ] /metrics端点可访问
- [ ] http_requests_total计数器正常增长
- [ ] http_request_duration_seconds记录延迟分布
- [ ] service_up指标显示为1
- [ ] 活跃连接数正确显示

---

## Prometheus查询示例

在Prometheus中可以使用以下PromQL查询：

```promql
# 请求速率（每秒请求数）
rate(http_requests_total[5m])

# 平均延迟ate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# P95延迟
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# 错误率
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# 活跃连接数
http_active_connections
```

---

## 下一步

完成本节点后，你已经实现了完善的日志和监控基础。接下来进入：

**节点1.5：数据库层设计** → [继续学习](./node-1-5.md)

在那里你将：
- 设计数据库表结构
- 使用GORM进行ORM操作
- 实现Repository模式
- 添加数据库迁移

---

## 参考资源

- [Zap官方文档](https://github.com/uber-go/zap)
- [Lumberjack日志轮转](https://github.com/natefinch/lumberjack)
- [Prometheus Go客户端](https://github.com/prometheus/client_golang)
- [PromQL查询语言](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Gin Prometheus中间件示例](https://github.com/zsais/go-gin-prometheus)
