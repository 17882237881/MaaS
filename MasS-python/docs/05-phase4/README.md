# 阶段4：高级特性与优化

> 🚀 **阶段目标**：为MaaS平台添加高级特性，进行性能优化，搭建CI/CD流水线，打造企业级AI平台

## 阶段总览

### 学习时间
**3周**（约15个工作日）

### 核心目标
1. 实现GPU资源调度和模型部署
2. 设计多租户隔离和配额管理
3. 进行全面性能优化
4. 搭建完整的CI/CD流水线

### 最终产出
- ✅ GPU调度器和模型部署控制器
- ✅ 多租户隔离和资源配额
- ✅ 性能优化后达标
- ✅ GitHub Actions自动化CI/CD流水线

## 技术栈详解

### 1. GPU调度与模型部署

**NVIDIA Docker**：在容器中使用GPU资源。

**Kubernetes GPU调度**：
```yaml
resources:
  limits:
    nvidia.com/gpu: 1  # 请求1个GPU
```

**调度策略**：
| 策略 | 说明 | 适用场景 |
|------|------|----------|
| 独占模式 | 一个模型独占一个GPU | 大模型推理 |
| 共享模式 | 多个模型共享GPU（MPS/MIG） | 小模型推理 |
| 弹性伸缩 | 按需分配GPU | 流量波动大 |

### 2. 多租户系统

**租户隔离级别**：
| 级别 | 隔离方式 | 优点 | 缺点 |
|------|----------|------|------|
| 数据库级 | 每租户一个数据库 | 隔离性最强 | 成本高 |
| Schema级 | 每租户一个Schema | 较好隔离 | 管理复杂 |
| 行级 | tenant_id字段区分 | 成本低 | 需代码保证 |

### 3. 性能优化

**Python性能分析工具**：
| 工具 | 用途 |
|------|------|
| py-spy | 采样式性能分析（无需修改代码） |
| cProfile | Python内置的确定性分析 |
| line_profiler | 逐行性能分析 |
| memory_profiler | 内存使用分析 |
| yappi | 支持多线程/协程的分析器 |

### 4. CI/CD

**GitHub Actions**：GitHub原生的CI/CD平台。

**典型流程**：
```
Push → Lint/Test → Build → Push Image → Deploy to Staging → Deploy to Production
```

---

## 节点详解

### 节点4.1：模型部署与调度（5天）

**学习目标**：
- 实现模型部署控制器
- 设计GPU资源调度策略
- 实现灰度发布
- 配置自动伸缩

**部署控制器设计**：
```python
from abc import ABC, abstractmethod
from enum import Enum

class DeploymentStrategy(str, Enum):
    ROLLING = "rolling"       # 滚动更新
    BLUE_GREEN = "blue_green" # 蓝绿部署
    CANARY = "canary"         # 金丝雀发布

class DeploymentController:
    def __init__(self, k8s_client, scheduler):
        self.k8s = k8s_client
        self.scheduler = scheduler

    async def deploy_model(
        self,
        model_id: str,
        version: str,
        strategy: DeploymentStrategy,
        gpu_request: int = 0,
    ):
        # 1. 检查资源可用性
        available = await self.scheduler.check_resources(gpu_request)
        if not available:
            raise ResourceError("Insufficient GPU resources")

        # 2. 创建推理容器
        deployment = self._build_deployment(model_id, version, gpu_request)

        # 3. 根据策略部署
        if strategy == DeploymentStrategy.CANARY:
            await self._canary_deploy(deployment)
        elif strategy == DeploymentStrategy.BLUE_GREEN:
            await self._blue_green_deploy(deployment)
        else:
            await self._rolling_deploy(deployment)

    async def _canary_deploy(self, deployment):
        """金丝雀发布：先部署10%流量，观察后逐步扩大"""
        # 1. 部署新版本（1个副本）
        await self.k8s.create_deployment(deployment, replicas=1)
        # 2. 配置流量权重（10%）
        await self.k8s.update_traffic_weight(deployment, weight=10)
        # 3. 监控指标...
```

**GPU调度器**：
```python
class GPUScheduler:
    async def check_resources(self, gpu_count: int) -> bool:
        """检查集群中是否有足够的GPU资源"""
        nodes = await self.k8s.list_nodes(label_selector="gpu=true")
        available = sum(
            node.status.allocatable.get("nvidia.com/gpu", 0)
            - node.status.allocated.get("nvidia.com/gpu", 0)
            for node in nodes
        )
        return available >= gpu_count

    async def schedule(self, model_id: str, gpu_count: int) -> str:
        """调度模型到合适的GPU节点"""
        nodes = await self._get_available_nodes(gpu_count)
        if not nodes:
            raise SchedulerError("No available GPU nodes")

        # 选择最优节点（最少负载优先）
        best_node = min(nodes, key=lambda n: n.gpu_utilization)
        return best_node.name
```

**实操任务**：
1. 实现部署控制器框架
2. 实现GPU资源调度器
3. 实现滚动更新策略
4. 实现金丝雀发布流程
5. 配置K8s HPA自动伸缩

**检查点**：
- [ ] 模型能部署到K8s集群
- [ ] GPU资源调度正常
- [ ] 金丝雀发布流程完整
- [ ] 自动伸缩策略生效

---

### 节点4.2：多租户与配额（4天）

**学习目标**：
- 实现租户隔离（行级隔离）
- 设计资源配额系统
- 实现计费基础
- 添加租户管理接口

**行级租户隔离实现**：
```python
from sqlalchemy import event

class TenantMixin:
    """租户隔离Mixin，所有多租户表继承此类"""
    tenant_id: Mapped[str] = mapped_column(index=True)

class Model(Base, TenantMixin):
    __tablename__ = "models"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str]
    # ...

# SQLAlchemy事件：自动添加tenant_id过滤
class TenantAwareRepository:
    def __init__(self, session_factory, tenant_id: str):
        self._session_factory = session_factory
        self._tenant_id = tenant_id

    async def list_models(self) -> list[Model]:
        async with self._session_factory() as session:
            stmt = select(Model).where(Model.tenant_id == self._tenant_id)
            result = await session.execute(stmt)
            return result.scalars().all()
```

**资源配额管理**：
```python
from pydantic import BaseModel

class TenantQuota(BaseModel):
    max_models: int = 100
    max_storage_gb: float = 50.0
    max_gpu_hours: float = 100.0
    max_api_calls_per_day: int = 10000

class QuotaManager:
    async def check_quota(self, tenant_id: str, resource: str, amount: float) -> bool:
        quota = await self.get_quota(tenant_id)
        usage = await self.get_usage(tenant_id, resource)
        limit = getattr(quota, f"max_{resource}")
        return usage + amount <= limit

    async def record_usage(self, tenant_id: str, resource: str, amount: float):
        """记录资源使用量"""
        key = f"usage:{tenant_id}:{resource}"
        await self.redis.incrbyfloat(key, amount)
```

**FastAPI租户依赖注入**：
```python
async def get_current_tenant(user = Depends(get_current_user)) -> str:
    """从JWT中提取tenant_id"""
    return user.get("tenant_id")

@app.get("/api/v1/models")
async def list_models(tenant_id: str = Depends(get_current_tenant)):
    repo = TenantAwareRepository(session_factory, tenant_id)
    return await repo.list_models()
```

**实操任务**：
1. 实现TenantMixin和行级隔离
2. 实现TenantAwareRepository
3. 设计配额模型和检查逻辑
4. 实现资源用量统计（Redis）
5. 添加租户管理API

**检查点**：
- [ ] 不同租户数据完全隔离
- [ ] 超出配额请求被拒绝
- [ ] 资源用量统计准确
- [ ] 租户管理接口正常

---

### 节点4.3：性能优化（4天）

**学习目标**：
- 使用py-spy进行性能分析
- 优化数据库查询
- 优化异步并发
- 实现连接池调优

**性能分析流程**：
```bash
# 使用py-spy采样（无需修改代码，对标Go的pprof）
py-spy record -o profile.svg -- python -m uvicorn api_gateway.main:app

# 使用cProfile分析
python -m cProfile -o profile.pstat api_gateway/main.py

# 可视化（snakeviz）
pip install snakeviz
snakeviz profile.pstat
```

**数据库查询优化**：
```python
# 1. N+1问题优化（使用selectinload预加载）
stmt = select(Model).options(
    selectinload(Model.tags),
    selectinload(Model.versions)
).where(Model.status == "active")

# 2. 批量操作（对标GORM的CreateInBatches）
async with session.begin():
    session.add_all([model1, model2, model3])

# 3. 只查询需要的列
stmt = select(Model.id, Model.name, Model.version).where(...)

# 4. 索引优化
class Model(Base):
    __table_args__ = (
        Index("idx_models_status", "status"),
        Index("idx_models_tenant_owner", "tenant_id", "owner_id"),
    )
```

**异步并发优化**：
```python
import asyncio

# 并发执行多个独立操作（对标Go的errgroup）
async def get_model_detail(model_id: str):
    model, tags, versions = await asyncio.gather(
        repo.get_model(model_id),
        repo.get_tags(model_id),
        repo.get_versions(model_id),
    )
    return ModelDetail(model=model, tags=tags, versions=versions)

# 连接池配置优化
engine = create_async_engine(
    url,
    pool_size=20,          # 连接池大小
    max_overflow=10,       # 最大溢出连接数
    pool_timeout=30,       # 获取连接超时
    pool_recycle=1800,     # 连接回收时间
    echo=False,            # 关闭SQL日志
)
```

**压力测试（locust）**：
```python
from locust import HttpUser, task, between

class MaaSUser(HttpUser):
    wait_time = between(1, 3)

    @task(3)
    def list_models(self):
        self.client.get("/api/v1/models", headers=self.headers)

    @task(1)
    def create_model(self):
        self.client.post("/api/v1/models", json={...}, headers=self.headers)
```

**实操任务**：
1. 使用py-spy分析性能瓶颈
2. 优化Top-3慢查询（添加索引、预加载）
3. 优化异步并发（asyncio.gather）
4. 调优连接池参数
5. 使用locust进行压力测试
6. 生成性能报告

**检查点**：
- [ ] 识别Top-3性能瓶颈
- [ ] P99延迟降低30%+
- [ ] 数据库连接池使用率正常
- [ ] 压力测试QPS达标

---

### 节点4.4：CI/CD流水线（3天）

**学习目标**：
- 搭建GitHub Actions CI/CD
- 实现自动化测试流程
- 实现自动化构建和部署
- 配置代码质量检查

**GitHub Actions工作流**：
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - name: Install dependencies
        run: |
          pip install poetry
          poetry install
      - name: Lint
        run: |
          poetry run ruff check .
          poetry run mypy .

  test:
    runs-on: ubuntu-latest
    needs: lint
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
      redis:
        image: redis:7
        ports:
          - 6379:6379
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - name: Install dependencies
        run: |
          pip install poetry
          poetry install
      - name: Run tests
        run: poetry run pytest --cov=. --cov-report=xml -v
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      - name: Build Docker images
        run: |
          docker build -t maas/api-gateway:${{ github.sha }} ./api_gateway
          docker build -t maas/model-registry:${{ github.sha }} ./model_registry
      - name: Push to registry
        run: |
          docker push maas/api-gateway:${{ github.sha }}
          docker push maas/model-registry:${{ github.sha }}

  deploy:
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Deploy to K8s
        run: |
          kubectl set image deployment/api-gateway \
            api-gateway=maas/api-gateway:${{ github.sha }}
          kubectl set image deployment/model-registry \
            model-registry=maas/model-registry:${{ github.sha }}
```

**代码质量工具**：
```toml
# pyproject.toml
[tool.ruff]
line-length = 120
target-version = "py311"
select = ["E", "F", "W", "I", "N", "UP"]

[tool.mypy]
python_version = "3.11"
strict = true
warn_return_any = true
disallow_untyped_defs = true

[tool.pytest.ini_options]
asyncio_mode = "auto"
testpaths = ["tests"]
```

**实操任务**：
1. 创建GitHub Actions工作流（lint + test + build）
2. 配置ruff和mypy代码检查
3. 配置pytest自动化测试（含覆盖率）
4. 实现Docker镜像自动构建和推送
5. 实现K8s自动部署

**检查点**：
- [ ] Push代码自动触发CI
- [ ] Lint和Test通过才能合并
- [ ] Docker镜像自动构建
- [ ] 主分支合并自动部署

---

## 阶段4里程碑

### 完成检查清单

- [ ] GPU调度器正常工作
- [ ] 模型部署控制器支持多种策略
- [ ] 多租户数据完全隔离
- [ ] 配额管理系统生效
- [ ] 性能优化后P99延迟达标
- [ ] CI/CD流水线完整运行
- [ ] 代码质量检查通过

### 可演示功能
1. 部署模型到GPU节点
2. 金丝雀发布新版本
3. 不同租户看到不同数据
4. 超出配额时请求被拒绝
5. 性能测试报告
6. Push代码自动触发CI/CD

---

## 🎉 项目完成

恭喜你完成了整个MaaS平台的学习！

### 你现在掌握了：

| 领域 | Go版技术栈 | Python版技术栈 |
|------|-----------|---------------|
| HTTP框架 | Gin | FastAPI |
| ORM | GORM | SQLAlchemy 2.0 |
| RPC通信 | gRPC (Go) | gRPC (grpcio) |
| 配置管理 | Viper | Pydantic-Settings |
| 日志 | Zap | Loguru |
| 消息队列 | confluent-kafka-go | confluent-kafka-python |
| 缓存 | go-redis | redis-py |
| 限流 | 自定义 | slowapi |
| 熔断 | 自定义 | pybreaker |
| 追踪 | OpenTelemetry | OpenTelemetry |
| 测试 | go test | pytest |
| CI/CD | GitHub Actions | GitHub Actions |

### 下一步建议

1. **完善项目**：添加更多现实功能（如计费、审计日志）
2. **开源贡献**：参与FastAPI/SQLAlchemy等项目
3. **面试准备**：基于项目经验准备面试
4. **持续学习**：关注Python生态新发展

---

**感谢学习，祝你前程似锦！** 🚀
