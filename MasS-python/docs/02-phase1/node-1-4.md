# 节点1.4：日志与监控基础（Logging & Monitoring）

## 学习目标
- 掌握 Loguru 日志库的高级用法
- 理解结构化日志（Structured Logging）的重要性
- 实现 Prometheus 监控指标（Metrics）
- 理解 RED 方法论（Rate, Errors, Duration）
- 搭建可观测性基础架构

## 1. 结构化日志 (Structured Logging)

### 为什么 `print()` 是不够的？
在简单的脚本中，`print("User logged in")` 也许够用。但在微服务和生产环境中，我们需要回答复杂的问题：
- "过去一小时有多少用户登录失败？"
- "这个 API 的平均延迟是多少？"
- "某个请求 ID 在所有服务中的完整路径是什么？"

文本日志（Text Logs）像这样：
```text
2024-03-20 10:00:00 [INFO] User 123 login failed: bad password
```
很难被机器解析和聚合。

结构化日志（JSON Logs）像这样：
```json
{
  "time": "2024-03-20 10:00:00",
  "level": "INFO",
  "event": "user_login_failed",
  "user_id": 123,
  "reason": "bad_password",
  "request_id": "req-abc-123"
}
```
可以直接被 ELK (Elasticsearch), Splunk, Datadog 等工具索引和查询。

### Loguru 最佳实践
[Loguru](https://github.com/Delgan/loguru) 是 Python 中最优雅的日志库。

**配置策略**：
- **开发环境**：输出带颜色的文本日志到控制台，方便人类阅读。
- **生产环境**：输出 JSON 格式日志到文件或标准输出（ stdout ），方便机器采集。
- **文件轮转 (Rotation)**：每 100MB 或每天轮转一次，防止单个日志文件无限膨胀。
- **保留策略 (Retention)**：保留最近 30 天的日志。

## 2. 监控指标 (Metrics)

日志记录了"由于什么原因发生了什么事"（离散事件），而监控记录了"系统的健康状况如何"（聚合趋势）。

### Prometheus 三大指标类型

1.  **Counter (计数器)**
    - 特点：只增不减（除非重启）。
    - 场景：请求总数 (`http_requests_total`)、错误总数、完成任务数。
    - QL示例：`rate(http_requests_total[5m])` 计算 QPS。

2.  **Gauge (仪表盘)**
    - 特点：可增可减，反映瞬时状态。
    - 场景：当前并发请求数 (`http_requests_in_flight`)、内存使用量、Goroutine/线程数。

3.  **Histogram (直方图)**
    - 特点：将数据落入不同的桶 (Buckets) 中，统计分布。
    - 场景：请求延迟 (`http_request_duration_seconds`)、响应大小。
    - QL示例：`histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` 计算 P99 延迟。

### RED 方法论
在微服务监控中，黄金指标（Golden Signals）通常被概括为 RED：
- **Rate (速率)**：每秒请求数 (QPS)。
- **Errors (错误)**：每秒失败请求数 (Error Rate)。
- **Duration (耗时)**：请求处理耗时 (Latency)。

## 3. 实现细节

### 安装依赖
```bash
uv add prometheus-client
```

### Loguru 配置 (`pkg/logger/setup.py`)
```python
logger.add(
    "logs/app.log",
    rotation="100 MB",
    retention="10 days",
    serialize=True  # 输出 JSON
)
```

### Prometheus 中间件 (`internal/middleware/metrics.py`)
我们将实现一个中间件，自动拦截所有请求，并记录：
- `http_requests_total`: Counter, labels=[method, path, status]
- `http_request_duration_seconds`: Histogram, labels=[method, path]

### 暴露端点
Prometheus 需要一个 HTTP 接口来"抓取"（Scrape）数据。通常是 `/metrics`。
FastAPI 可以通过集成 `make_asgi_app` 轻松实现：

```python
from prometheus_client import make_asgi_app
app.mount("/metrics", make_asgi_app())
```

## 检查点
- [ ] 开发环境下控制台看到彩色日志
- [ ] 生产环境下 logs 目录看到 JSON 日志文件
- [ ] 访问 `/metrics` 能看到 `http_requests_total` 等指标
- [ ] 发送请求后，指标数值正确增加

## 下一步
下一节我们将进入 Phase 1 的尾声：[节点1.5：Docker化与部署](./node-1-5.md)。
