# 阶段3：企业级特性

> 🏢 **阶段目标**：为MaaS平台添加生产级特性，使其具备真正的高可用、可观测和弹性能力

## 阶段总览

### 学习时间
**4周**（约20个工作日）

### 核心目标
1. 掌握消息队列异步处理（Kafka）
2. 学会Kubernetes容器编排
3. 实现限流与熔断机制
4. 处理分布式事务
5. 搭建完整的监控告警体系

### 最终产出
- ✅ Kafka异步任务处理正常运行
- ✅ Kubernetes部署配置完成
- ✅ 限流和熔断机制生效
- ✅ 分布式事务保障一致性
- ✅ Prometheus + Grafana监控仪表盘

## 技术栈详解

### 1. Kafka消息队列

**什么是消息队列？**
消息队列在发送者和接收者之间传递消息，实现了系统解耦和异步处理。Python版使用 `confluent-kafka-python`。

**使用场景**：
- **异步处理**：用户上传模型后，异步进行验证、转换、部署
- **削峰填谷**：高并发时将请求放入队列缓冲
- **系统解耦**：服务间通过消息通信，互不依赖

**核心概念**：
| 概念 | 说明 |
|------|------|
| Topic | 消息的分类（如 model-events, inference-tasks） |
| Partition | Topic的分区，实现并行处理 |
| Producer | 消息生产者 |
| Consumer | 消息消费者 |
| Consumer Group | 消费者组，组内成员分担消息 |
| Offset | 消息在分区中的位置偏移 |

### 2. Kubernetes容器编排

**什么是Kubernetes（K8s）？**
Kubernetes是容器编排平台，自动化部署、扩缩容和管理容器化应用。

**核心概念**：
| 概念 | 说明 |
|------|------|
| Pod | 最小部署单元，包含1+个容器 |
| Deployment | 管理Pod的副本数和更新策略 |
| Service | 为Pod提供稳定的网络端口 |
| ConfigMap | 存储配置数据 |
| Secret | 存储敏感数据 |
| Ingress | HTTP路由规则 |
| HPA | 自动水平伸缩 |

### 3. 限流与熔断

**限流（Rate Limiting）**
限制单位时间内的请求数量，防止系统过载。Python版使用 `slowapi`（基于limits库）。

**常见算法**：
1. **固定窗口**：每分钟最多N个请求
2. **滑动窗口**：更精确的流量控制
3. **令牌桶**：以固定速率放入令牌，有桶容量限制
4. **漏桶**：以固定速率处理请求

**熔断（Circuit Breaker）**
当下游服务异常时自动"切断"请求，避免级联故障。Python版使用 `pybreaker`。

**熔断器三种状态**：
```
Closed（关闭）→ 正常通过请求
    ↓ 错误率超过阈值
Open（打开）→ 拒绝所有请求
    ↓ 等待超时
Half-Open（半开）→ 允许少量请求试探
    ↓ 成功则 → Closed
    ↓ 失败则 → Open
```

### 4. 分布式事务

**什么是分布式事务？**
跨多个服务/数据库的操作需要保持原子性。

**常见方案**：
| 方案 | 一致性 | 复杂度 | 性能 |
|------|--------|--------|------|
| 2PC（两阶段提交） | 强一致 | 高 | 低 |
| Saga模式 | 最终一致 | 中 | 高 |
| 本地消息表 | 最终一致 | 低 | 高 |
| TCC | 强一致 | 很高 | 中 |

### 5. 监控告警体系

**可观测性三支柱**：
1. **Metrics（指标）**：Prometheus + Grafana（数值时序数据）
2. **Logging（日志）**：ELK Stack / Loki（事件文本数据）
3. **Tracing（链路追踪）**：OpenTelemetry + Jaeger（请求链路）

---

## 节点详解

### 节点3.1：消息队列集成（5天）

**学习目标**：
- 理解消息队列的核心概念
- 掌握Kafka Producer/Consumer（confluent-kafka）
- 实现可靠消息投递
- 处理消息消费失败

**Kafka Producer示例**：
```python
from confluent_kafka import Producer

conf = {
    'bootstrap.servers': 'localhost:9092',
    'client.id': 'model-registry',
}

producer = Producer(conf)

def delivery_report(err, msg):
    if err:
        logger.error(f"Delivery failed: {err}")
    else:
        logger.info(f"Message delivered to {msg.topic()} [{msg.partition()}]")

# 发送消息
import json
event = {"model_id": "abc123", "action": "created", "timestamp": "2024-01-01T00:00:00"}
producer.produce(
    topic="model-events",
    key="abc123",                    # 分区键
    value=json.dumps(event).encode(),
    callback=delivery_report
)
producer.flush()  # 等待所有消息发送完成
```

**Kafka Consumer示例**：
```python
from confluent_kafka import Consumer

conf = {
    'bootstrap.servers': 'localhost:9092',
    'group.id': 'model-processor',
    'auto.offset.reset': 'earliest',
    'enable.auto.commit': False,       # 手动提交保证可靠性
}

consumer = Consumer(conf)
consumer.subscribe(['model-events'])

try:
    while True:
        msg = consumer.poll(timeout=1.0)
        if msg is None:
            continue
        if msg.error():
            logger.error(f"Consumer error: {msg.error()}")
            continue

        # 处理消息
        event = json.loads(msg.value().decode())
        await process_model_event(event)

        # 手动提交offset
        consumer.commit(msg)
except KeyboardInterrupt:
    pass
finally:
    consumer.close()
```

**异步消费（asyncio集成）**：
```python
import asyncio
from confluent_kafka import Consumer

async def consume_loop(consumer: Consumer, topics: list[str]):
    """在asyncio事件循环中运行Kafka消费者"""
    consumer.subscribe(topics)
    loop = asyncio.get_event_loop()

    while True:
        msg = await loop.run_in_executor(None, consumer.poll, 1.0)
        if msg is None:
            continue
        # 处理...
```

**实操任务**：
1. Docker Compose部署Kafka和Zookeeper
2. 实现Producer封装（支持JSON序列化）
3. 实现Consumer封装（支持手动commit）
4. Model Registry创建模型时发送事件
5. 实现消费者处理模型事件
6. 处理消息消费失败和重试

**检查点**：
- [ ] Kafka集群正常启动
- [ ] Producer能发送消息
- [ ] Consumer能接收并处理消息
- [ ] 消费失败时有重试机制
- [ ] 消息不丢失（手动commit）

---

### 节点3.2：容器化与编排（5天）

**学习目标**：
- 编写生产级Dockerfile
- 创建Kubernetes Deployment/Service/ConfigMap
- 使用Helm管理部署
- 配置健康检查和资源限制

**生产级Python Dockerfile**：
```dockerfile
# 多阶段构建
FROM python:3.11-slim as builder

WORKDIR /app
# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

# Install dependencies
COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-cache
RUN pip install --no-cache-dir -r requirements.txt

FROM python:3.11-slim as runtime
WORKDIR /app

# 安全：非root用户
RUN groupadd -r appuser && useradd -r -g appuser appuser

COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY . .

USER appuser
EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=10s CMD curl -f http://localhost:8000/health || exit 1
CMD ["python", "-m", "uvicorn", "api_gateway.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

**Kubernetes Deployment示例**：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  labels:
    app: api-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: maas/api-gateway:latest
        ports:
        - containerPort: 8000
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8000
          periodSeconds: 10
        envFrom:
        - configMapRef:
            name: api-gateway-config
```

**实操任务**：
1. 编写各服务的生产级Dockerfile（多阶段构建）
2. 创建K8s Deployment/Service/ConfigMap
3. 配置健康检查和资源限制
4. 使用Helm打包部署
5. 配置HPA自动伸缩

**检查点**：
- [ ] Docker镜像构建成功
- [ ] K8s Deployment部署成功
- [ ] 健康检查通过
- [ ] Service能正常路由
- [ ] HPA自动伸缩配置完成

---

### 节点3.3：限流与熔断（4天）

**学习目标**：
- 实现API限流（slowapi）
- 实现熔断器模式（pybreaker）
- 配置降级策略
- 添加限流/熔断指标

**slowapi限流示例**：
```python
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address
from slowapi.errors import RateLimitExceeded

limiter = Limiter(key_func=get_remote_address)
app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)

@app.get("/api/v1/models")
@limiter.limit("100/minute")
async def list_models(request: Request):
    return await service.list_models()

# 基于用户的限流
@app.post("/api/v1/inference")
@limiter.limit("10/minute", key_func=lambda request: get_user_id(request))
async def inference(request: Request):
    return await service.inference()
```

**pybreaker熔断器示例**：
```python
import pybreaker

# 创建熔断器（5次失败后打开，30秒后半开）
model_breaker = pybreaker.CircuitBreaker(
    fail_max=5,
    reset_timeout=30,
    listeners=[],
)

@model_breaker
async def call_model_registry(request):
    """调用Model Registry（带熔断保护）"""
    response = await grpc_client.GetModel(request)
    return response

# 降级处理
async def get_model_with_fallback(model_id: str):
    try:
        return await call_model_registry(model_id)
    except pybreaker.CircuitBreakerError:
        # 熔断器打开时的降级策略
        logger.warning("Circuit breaker is open, using cache")
        return await redis.get(f"model:{model_id}")
```

**实操任务**：
1. 集成slowapi实现接口限流
2. 集成pybreaker实现熔断器
3. 实现降级策略（缓存回退）
4. 添加限流/熔断Prometheus指标
5. 测试各种故障场景

**检查点**：
- [ ] 超过频率限制返回429
- [ ] 下游服务故障时触发熔断
- [ ] 熔断时走降级逻辑
- [ ] 限流/熔断指标可观测

---

### 节点3.4：分布式事务（4天）

**学习目标**：
- 理解分布式事务问题
- 实现Saga事务编排
- 实现本地消息表模式
- 处理补偿操作

**Saga模式实现**：
```python
from dataclasses import dataclass
from typing import Callable, Awaitable

@dataclass
class SagaStep:
    name: str
    execute: Callable[..., Awaitable]   # 正向操作
    compensate: Callable[..., Awaitable]  # 补偿操作

class SagaOrchestrator:
    def __init__(self):
        self.steps: list[SagaStep] = []
        self.completed: list[SagaStep] = []

    def add_step(self, step: SagaStep):
        self.steps.append(step)

    async def execute(self, context: dict) -> bool:
        for step in self.steps:
            try:
                await step.execute(context)
                self.completed.append(step)
            except Exception as e:
                logger.error(f"Saga step '{step.name}' failed: {e}")
                await self._compensate(context)
                return False
        return True

    async def _compensate(self, context: dict):
        """反向补偿已完成的步骤"""
        for step in reversed(self.completed):
            try:
                await step.compensate(context)
                logger.info(f"Compensated: {step.name}")
            except Exception as e:
                logger.error(f"Compensation failed for '{step.name}': {e}")
```

**实操任务**：
1. 设计一个跨服务事务流程（如：创建模型+分配存储+注册版本）
2. 实现Saga编排器
3. 实现每步的补偿操作
4. 实现本地消息表作为备选方案
5. 测试各种失败场景

**检查点**：
- [ ] 正常流程全部步骤执行成功
- [ ] 某步骤失败时补偿操作执行
- [ ] 最终数据一致
- [ ] 日志记录完整

---

### 节点3.5：监控告警体系（4天）

**学习目标**：
- 集成OpenTelemetry链路追踪
- 搭建Prometheus + Grafana监控
- 配置Alertmanager告警规则
- 构建监控仪表盘

**OpenTelemetry链路追踪集成**：
```python
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor

# 配置TracerProvider
tracer_provider = TracerProvider()
jaeger_exporter = JaegerExporter(
    agent_host_name="localhost",
    agent_port=6831,
)
tracer_provider.add_span_processor(BatchSpanProcessor(jaeger_exporter))
trace.set_tracer_provider(tracer_provider)

# 自动注入FastAPI
FastAPIInstrumentor.instrument_app(app)

# 手动创建Span
tracer = trace.get_tracer(__name__)

async def create_model(data):
    with tracer.start_as_current_span("create_model") as span:
        span.set_attribute("model.name", data.name)

        # 嵌套Span
        with tracer.start_as_current_span("save_to_db"):
            await repo.create(model)

        with tracer.start_as_current_span("publish_event"):
            await kafka_producer.send(event)
```

**Prometheus指标（prometheus-client）**：
```python
from prometheus_client import make_asgi_app, Counter, Histogram

# 自定义指标
REQUEST_COUNT = Counter(
    "app_request_count_total",
    "Total request count",
    ["method", "endpoint", "status"]
)
REQUEST_LATENCY = Histogram(
    "app_request_latency_seconds",
    "Request latency",
    ["method", "endpoint"]
)

# 挂载/metrics端点
metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)
```

**Grafana Dashboard配置要点**：
- HTTP请求速率（QPS）
- 请求延迟分位数（P50/P95/P99）
- 错误率
- 活跃连接数
- 数据库连接池使用率
- 缓存命中率

**实操任务**：
1. 集成OpenTelemetry + Jaeger
2. 配置Prometheus指标收集
3. 搭建Grafana，导入Dashboard
4. 配置Alertmanager告警规则
5. 验证端到端的链路追踪

**检查点**：
- [ ] Jaeger能查看完整调用链
- [ ] Prometheus收集到所有指标
- [ ] Grafana仪表盘展示正常
- [ ] 告警规则能正确触发

---

## 阶段3里程碑

### 完成检查清单

- [ ] Kafka异步任务处理正常运行
- [ ] K8s部署配置完成，服务可运行
- [ ] 限流超过阈值时返回429
- [ ] 熔断器在故障时正确触发
- [ ] 分布式事务数据一致
- [ ] Prometheus + Grafana监控可用
- [ ] Jaeger链路追踪可用
- [ ] 告警规则配置完成

### 可演示功能
1. 创建模型 → Kafka事件 → 异步处理
2. 高并发请求 → 限流返回429
3. 关闭下游服务 → 熔断 → 降级
4. Grafana仪表盘展示QPS/延迟/错误率
5. Jaeger展示请求全链路

### 下一步
进入**阶段4：高级特性与优化**，学习GPU调度、多租户、性能优化和CI/CD。

---

**继续学习**：[阶段4文档](../05-phase4/README.md)
