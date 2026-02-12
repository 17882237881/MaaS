# 阶段1：基础架构搭建

> 🏗️ **阶段目标**：搭建MaaS平台的基础架构，实现API Gateway和Model Registry两个核心服务

## 阶段总览

### 学习时间
**4周**（约20个工作日）

### 核心目标
1. 掌握Python企业级项目结构组织（uv + pyproject.toml）
2. 学会使用FastAPI框架开发异步RESTful API
3. 掌握SQLAlchemy 2.0进行数据库设计和操作
4. 学会Docker容器化部署
5. 理解微服务拆分的基本原则

### 最终产出
- ✅ 可独立运行的API Gateway服务
- ✅ 可独立运行的Model Registry服务
- ✅ 支持基础CRUD操作
- ✅ Docker Compose一键启动所有服务

## 技术栈详解

### 1. Python 3.11+

**为什么要用这个版本？**
- Python 3.11性能提升25%（CPython优化）
- 改进的错误消息和异常链
- `asyncio.TaskGroup` 简化并发任务管理
- 更好的类型注解支持

**核心特性回顾**：
```python
# 数据类（类似Go的struct）
from dataclasses import dataclass

@dataclass
class User:
    id: str
    name: str

    def get_name(self) -> str:
        return self.name

# 抽象基类（类似Go的interface）
from abc import ABC, abstractmethod

class Service(ABC):
    @abstractmethod
    async def do_something(self) -> None:
        pass

# 异步编程（类似goroutine）
import asyncio

async def task():
    await asyncio.sleep(1)
    return "done"

# 并发执行多个任务
async def main():
    async with asyncio.TaskGroup() as tg:
        t1 = tg.create_task(task())
        t2 = tg.create_task(task())
    print(t1.result(), t2.result())
```

**Go vs Python 对比**：
| 概念 | Go | Python |
|------|-----|--------|
| 结构体/类 | `type User struct {}` | `class User:` / `@dataclass` |
| 接口 | `type Service interface {}` | `class Service(ABC):` |
| 并发 | `go func() {}()` | `asyncio.create_task()` |
| 通道/队列 | `ch := make(chan T)` | `asyncio.Queue()` |
| 错误处理 | `if err != nil { return err }` | `try/except` + 自定义Exception |

### 1. Modern Python Package Management: uv

**什么是 uv？**
`uv` 是一个由 Rust 编写的极速 Python 包管理工具，旨在替代 `pip`、`pip-tools` 和 `poetry`。

**为什么选择 uv？**
- **极速**：比 pip 快 10-100 倍
- **统一**：集成了 Python 版本管理、包管理、虚拟环境管理
- **兼容**：使用标准的 `pyproject.toml`
- **磁盘空间优化**：全局缓存，支持硬链接

**uv vs Conda 对比**：
| 特性 | Conda | uv |
|------|-------|----|
| 语言 | Python/C | Rust |
| 包源 | Anaconda/conda-forge | PyPI |
| 环境隔离 | 强（包含非Python依赖） | 标准 venv |
| 速度 | 较慢 | **极速** |
| 依赖解析 | 较慢 | **极速** |
| 用法复杂度 | 中等 | **极简** |

**常用命令**：
```bash
# 初始化项目
uv init

# 添加依赖
uv add fastapi

# 运行命令（自动同步环境）
uv run python main.py

# 查看依赖树
uv tree

# 同步环境（根据 lock 文件）
uv sync
```

### 2. FastAPI Web框架

**什么是FastAPI？**
FastAPI是Python中速度最快的异步Web框架之一，基于Starlette和Pydantic。对标Go中的Gin框架。

**核心概念**：
- **路由（Router）**：URL路径和处理函数的映射
- **中间件（Middleware）**：请求处理链，可插入日志、认证等功能
- **依赖注入（Depends）**：FastAPI的核心特色，自动解析依赖
- **Pydantic模型**：自动将JSON数据绑定到类型安全的模型

**简单示例**：
```python
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

app = FastAPI()

# 定义路由（对标Gin的r.GET）
@app.get("/hello")
async def hello():
    return {"message": "Hello World"}

# 带参数的路由（对标Gin的c.Param）
@app.get("/users/{user_id}")
async def get_user(user_id: str):
    return {"id": user_id}

# POST请求 + Pydantic校验（对标Gin的ShouldBindJSON）
class UserRequest(BaseModel):
    name: str = Field(..., min_length=1)
    email: str = Field(..., pattern=r"^[\w.-]+@[\w.-]+\.\w+$")

@app.post("/users")
async def create_user(req: UserRequest):
    # Pydantic自动校验，无效数据返回422
    return {"name": req.name, "email": req.email}
```

**Go Gin vs Python FastAPI 对比**：
| 功能 | Gin (Go) | FastAPI (Python) |
|------|----------|------------------|
| 路由定义 | `r.GET("/path", handler)` | `@app.get("/path")` |
| JSON绑定 | `c.ShouldBindJSON(&req)` | 参数类型注解自动绑定 |
| 路径参数 | `c.Param("id")` | 函数参数 `id: str` |
| 查询参数 | `c.Query("page")` | 函数参数 `page: int = 1` |
| 中间件 | `r.Use(middleware)` | `app.add_middleware(cls)` |
| Swagger | 需要swag注解 | **自动生成** |
| 响应 | `c.JSON(200, data)` | 直接return dict |

### 3. SQLAlchemy 2.0（ORM）

**什么是ORM？**
ORM（Object-Relational Mapping）对象关系映射，让你用操作对象的方式操作数据库，不用写SQL。SQLAlchemy 2.0对标Go中的GORM。

**SQLAlchemy 2.0的核心功能**：
- 原生异步支持（AsyncSession）
- 类型安全的 Mapped 注解
- 关联（一对一、一对多、多对多）
- 事务支持
- 性能极佳

**示例**：
```python
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship
from sqlalchemy import String, DateTime, func
from datetime import datetime

# 基类
class Base(DeclarativeBase):
    pass

# 定义模型（对标GORM的Model struct）
class User(Base):
    __tablename__ = "users"

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(255))
    email: Mapped[str] = mapped_column(String(255), unique=True)
    created_at: Mapped[datetime] = mapped_column(
        DateTime, server_default=func.now()
    )

# 异步连接数据库
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker

engine = create_async_engine("postgresql+asyncpg://user:pass@localhost/db")
async_session = async_sessionmaker(engine)

# CRUD操作
async with async_session() as session:
    # 创建（对标 db.Create）
    user = User(name="John", email="john@example.com")
    session.add(user)
    await session.commit()

    # 查询（对标 db.First）
    from sqlalchemy import select
    result = await session.execute(select(User).where(User.id == 1))
    user = result.scalar_one()

    # 更新（对标 db.Model.Update）
    user.name = "Jane"
    await session.commit()

    # 删除（对标 db.Delete）
    await session.delete(user)
    await session.commit()
```

**GORM vs SQLAlchemy 2.0 对比**：
| 功能 | GORM (Go) | SQLAlchemy 2.0 (Python) |
|------|-----------|-------------------------|
| 模型定义 | struct tag | Mapped 类型注解 |
| 连接 | `gorm.Open(postgres.Open(dsn))` | `create_async_engine(url)` |
| 创建 | `db.Create(&user)` | `session.add(user)` |
| 查询 | `db.First(&user, id)` | `session.execute(select(User))` |
| 预加载 | `db.Preload("Tags")` | `selectinload(Model.tags)` |
| 迁移 | `db.AutoMigrate(&User{})` | Alembic |
| 事务 | `db.Transaction(func(tx) error)` | `async with session.begin():` |

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
- **连接池**：复用数据库连接，提高性能（asyncpg天然支持）

### 5. Loguru日志库

**为什么要用Loguru？**
- 零配置即可使用，API极简
- 支持结构化日志（JSON格式）
- 自动包含调用者信息
- 内置日志文件轮转
- 对标Go中的Zap日志库

**日志级别说明**：
- **DEBUG**：开发调试信息（如变量值）
- **INFO**：正常流程信息（如请求处理完成）
- **WARNING**：警告信息（如性能下降）
- **ERROR**：错误信息（如数据库连接失败）
- **CRITICAL**：致命错误（程序无法继续）

**Loguru示例**：
```python
from loguru import logger

# 零配置即用
logger.info("Server starting", port=8080)

# JSON格式输出（对标Zap的JSON encoder）
logger.add("logs/app.log", serialize=True, rotation="100 MB")

# 带上下文字段（对标Zap的With fields）
req_logger = logger.bind(request_id="abc-123")
req_logger.info("Request processed", method="GET", path="/api/users")

# 异常自动捕获堆栈
@logger.catch
async def risky_function():
    raise ValueError("something went wrong")
```

**Zap vs Loguru 对比**：
| 功能 | Zap (Go) | Loguru (Python) |
|------|----------|-----------------|
| 初始化 | `zap.NewProduction()` | 零配置，直接用 |
| 结构化日志 | `zap.String("key", val)` | `logger.info("msg", key=val)` |
| JSON输出 | `NewJSONEncoder` | `serialize=True` |
| 文件轮转 | 需要lumberjack | 内置 `rotation` |
| 级别过滤 | `zap.NewAtomicLevel()` | `logger.add(level="INFO")` |

### 6. Docker

**什么是Docker？**
Docker是一种容器化技术，把应用和依赖打包成一个"集装箱"，在任何地方都能一致运行。

**核心概念**：
- **镜像（Image）**：只读的模板，包含运行应用所需的一切
- **容器（Container）**：镜像的运行实例
- **Dockerfile**：定义如何构建镜像的脚本
- **Docker Compose**：定义和运行多容器应用的工具

**Python Dockerfile示例**：
```dockerfile
# 基础镜像
FROM python:3.11-slim

# 设置工作目录
WORKDIR /app

# 安装uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

# 复制依赖定义
COPY pyproject.toml uv.lock ./

# 安装依赖（无缓存模式减小镜像体积）
RUN uv sync --frozen --no-cache

# 复制源代码
COPY . .

# 暴露端口
EXPOSE 8000

# 运行（使用uv run）
CMD ["uv", "run", "uvicorn", "api_gateway.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

## 节点详解

### 节点1.1：项目初始化与架构设计（3天）

**学习目标**：
- 理解uv包管理工具
- 掌握Python微服务目录结构组织
- 学会依赖版本管理
- 建立Git工作流规范

**技术介绍**：

**1. uv包管理**
uv是Python的现代包管理和构建工具，比Poetry更快。

**关键命令**：
```bash
# 初始化项目（对标 go mod init）
uv init

# 安装依赖（对标 go mod download）
uv sync

# 添加依赖（对标 go get）
uv add fastapi uvicorn

# 添加开发依赖
uv add --dev pytest ruff mypy

# 查看依赖树（对标 go mod graph）
uv tree
```

**2. 微服务目录结构**
推荐的目录组织方式（严格对齐Go版）：
```
MasS-python/
├── api_gateway/              # 对齐 api-gateway/
│   ├── main.py              # 对齐 cmd/main.go
│   ├── config/
│   │   ├── config.py        # 对齐 internal/config/config.go
│   │   └── config.yaml
│   ├── internal/
│   │   ├── handler/         # HTTP处理器
│   │   ├── middleware/      # 中间件
│   │   ├── model/           # 数据模型
│   │   ├── repository/      # 数据访问层
│   │   ├── router/          # 路由定义
│   │   └── service/         # 业务逻辑层
│   └── pkg/
│       ├── logger/          # 日志工具
│       ├── metrics/         # Prometheus指标
│       └── grpc_client/     # gRPC客户端
├── model_registry/           # 对齐 model-registry/
│   └── ...                  # 相同结构
├── shared/                   # 共享代码
│   └── proto/               # Protocol Buffers定义
├── tests/                    # 测试目录
├── pyproject.toml            # 对齐 go.mod
├── uv.lock                   # 依赖锁文件
├── Makefile
└── README.md
```

**分层架构说明**（与Go版完全一致）：
- **Handler层**：处理HTTP请求，参数校验，调用Service
- **Service层**：业务逻辑，事务管理
- **Repository层**：数据访问，数据库操作
- **Model层**：数据结构定义

**为什么这样分层？**
- **单一职责**：每层只负责一件事
- **可测试性**：每层可以独立测试（用mock替换依赖）
- **可替换性**：可以替换某层实现（如换数据库）

**实操任务**：
1. 安装uv，创建pyproject.toml
2. 创建api_gateway和model_registry目录结构
3. 创建基本的main.py文件（能启动FastAPI）
4. 初始化Git仓库，提交第一次commit
5. 编写README.md说明项目结构

**检查点**：
- [ ] 项目目录结构完整
- [ ] `uv sync` 成功
- [ ] 每个服务能独立启动运行
- [ ] Git仓库初始化完成

---

### 节点1.2：API Gateway核心（5天）

**学习目标**：
- 掌握FastAPI核心用法
- 实现全局异常处理
- 学会使用中间件
- 自动生成OpenAPI文档

**技术介绍**：

**1. HTTP基础回顾**
- **请求方法**：GET（获取）、POST（创建）、PUT（更新）、DELETE（删除）
- **状态码**：200（成功）、400（请求错误）、401（未授权）、500（服务器错误）
- **Header**：Content-Type、Authorization等
- **Body**：请求/响应的数据体

**2. RESTful API设计规范**
```
GET    /api/v1/users          # 获取用户列表
GET    /api/v1/users/{id}     # 获取单个用户
POST   /api/v1/users          # 创建用户
PUT    /api/v1/users/{id}     # 更新用户
DELETE /api/v1/users/{id}     # 删除用户
```

**3. 中间件（Middleware）**
中间件是处理HTTP请求的"钩子"，可以在请求处理前/后执行代码。

**中间件执行顺序**：
```
请求 → 中间件1 → 中间件2 → Handler → 中间件2 → 中间件1 → 响应
```

**常见中间件类型**：
- **日志中间件**：记录请求信息
- **异常中间件**：捕获异常，防止程序崩溃（对标Gin的Recovery）
- **CORS中间件**：处理跨域请求
- **认证中间件**：验证用户身份

**4. OpenAPI（Swagger）**
FastAPI自动生成OpenAPI文档，无需额外注解！这是FastAPI相比Gin的一大优势。

**实操任务**：
1. 创建FastAPI应用实例
2. 实现基础中间件（Logger、Recovery、CORS、RequestID）
3. 创建健康检查接口（GET /health）
4. 访问自动生成的API文档（/docs）
5. 实现简单的CRUD接口示例

**代码结构**：
```python
# main.py - 程序入口（对标cmd/main.go）
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from api_gateway.internal.middleware import LoggerMiddleware, RequestIDMiddleware
from api_gateway.internal.router import register_routes

app = FastAPI(title="MaaS Platform API", version="1.0.0")

# 全局中间件（对标r.Use）
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"])
app.add_middleware(LoggerMiddleware)
app.add_middleware(RequestIDMiddleware)

# 路由注册
register_routes(app)

# middleware/logger.py - 日志中间件（对标middleware.Logger）
from starlette.middleware.base import BaseHTTPMiddleware
import time

class LoggerMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        start = time.time()
        response = await call_next(request)
        duration = time.time() - start
        logger.info("Request completed",
            method=request.method,
            path=request.url.path,
            status=response.status_code,
            duration=f"{duration:.3f}s"
        )
        return response

# router/router.py - 路由定义（对标router.RegisterRoutes）
from fastapi import APIRouter

def register_routes(app: FastAPI):
    api = APIRouter(prefix="/api/v1")
    api.include_router(health_router)
    api.include_router(user_router)
    app.include_router(api)
```

**检查点**：
- [ ] 服务能启动并监听端口
- [ ] /health接口返回200
- [ ] 日志正确输出请求信息
- [ ] /docs 自动生成的API文档可访问
- [ ] 异常处理正常工作（未捕获异常不会导致服务崩溃）

---

### 节点1.3：配置管理体系（3天）

**学习目标**：
- 学会多环境配置管理
- 掌握Pydantic-Settings使用
- 理解配置热更新（Watchdog）
- 学会敏感信息处理

**技术介绍**：

**1. 为什么需要配置管理？**
不同环境（开发/测试/生产）需要不同的配置：
- 数据库连接信息
- API密钥
- 日志级别
- 服务端口号

**2. 配置来源优先级**（高→低，与Go版Viper一致）：
1. 命令行参数
2. 环境变量
3. 配置文件
4. 默认值

**3. Pydantic-Settings配置管理**
Pydantic-Settings对标Go的Viper库，支持YAML配置文件 + 环境变量 + 类型校验。

**示例配置**（config.yaml）：
```yaml
environment: development
port: 8000
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

**Python配置类映射**（对标Go的Config struct）：
```python
from pydantic_settings import BaseSettings
from pydantic import Field

class DatabaseConfig(BaseSettings):
    host: str = "localhost"
    port: int = 5432
    user: str = "postgres"
    password: str = "postgres"
    name: str = "maas_platform"
    ssl_mode: str = "disable"

class RedisConfig(BaseSettings):
    host: str = "localhost"
    port: int = 6379
    password: str = ""
    db: int = 0

class Config(BaseSettings):
    environment: str = "development"
    port: int = 8000
    log_level: str = "info"
    database: DatabaseConfig = DatabaseConfig()
    redis: RedisConfig = RedisConfig()

    model_config = {"env_prefix": "MAAS_"}  # 对标Viper的SetEnvPrefix
```

**4. 环境变量使用**
```bash
# .env文件
MAAS_DATABASE__HOST=localhost
MAAS_DATABASE__PORT=5432

# 环境变量自动覆盖配置文件（Pydantic-Settings内置支持）
```

**实操任务**：
1. 创建config.yaml配置文件
2. 使用Pydantic-Settings加载配置
3. 支持环境变量覆盖（MAAS_ 前缀）
4. 添加配置验证（Pydantic validator）
5. 实现配置热更新（Watchdog监听文件变化）

**检查点**：
- [ ] 配置能从YAML文件加载
- [ ] 环境变量能覆盖配置文件
- [ ] 配置验证正常工作
- [ ] 不同环境使用不同配置

---

### 节点1.4：日志与监控基础（4天）

**学习目标**：
- 掌握Loguru日志库
- 实现结构化日志
- 添加基础Prometheus Metrics
- 理解可观测性三支柱

**技术介绍**：

**1. 可观测性三支柱**
- **日志（Logging）**：记录离散事件（如"用户登录"）
- **指标（Metrics）**：记录聚合数据（如"QPS=100"）
- **追踪（Tracing）**：记录请求链路（如"API→Service→DB"）

**2. 结构化日志 vs 文本日志**
```python
# 文本日志（难解析）
print(f"User {username} logged in from {ip} at {time}")

# 结构化日志（JSON，易解析）- Loguru
logger.bind(username=username, ip=ip).info("user logged in")
# 输出: {"text": "user logged in", "username": "john", "ip": "1.2.3.4", ...}
```

**3. Loguru使用示例**
```python
from loguru import logger
import sys

# 配置日志（对标Zap的Config）
logger.remove()  # 移除默认handler
logger.add(
    sys.stdout,
    format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {message}",
    level="INFO",
)
logger.add(
    "logs/app.log",
    serialize=True,       # JSON格式（对标Zap的JSONEncoder）
    rotation="100 MB",    # 文件轮转
    retention="7 days",   # 保留天数
    compression="gz",     # 压缩
)

# 使用
logger.info("request processed",
    method="GET", path="/api/users", status=200, latency=0.1)
logger.error("database connection failed", error=str(err), host="localhost")
```

**4. 基础Prometheus Metrics**
使用prometheus-client库记录指标（对标Go的prometheus/client_golang）：
```python
from prometheus_client import Counter, Histogram, Gauge

# 计数器（只增不减）
REQUEST_TOTAL = Counter(
    "http_requests_total",
    "Total HTTP requests",
    ["method", "path", "status"]
)

# 直方图（记录分布）
REQUEST_DURATION = Histogram(
    "http_request_duration_seconds",
    "HTTP request duration",
    ["method", "path"],
    buckets=[0.1, 0.5, 1, 2, 5]
)

# 仪表盘（可增可减）
ACTIVE_CONNECTIONS = Gauge(
    "http_active_connections",
    "Active HTTP connections"
)
```

**实操任务**：
1. 集成Loguru日志库
2. 配置JSON格式输出和文件轮转
3. 添加请求日志中间件
4. 集成prometheus-client
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
- 学会SQLAlchemy 2.0高级用法
- 实现Repository模式
- 掌握Alembic数据库迁移

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
Repository模式是数据访问层的设计模式，将数据访问逻辑封装起来。对标Go版的interface + struct实现。

**优势**：
- 业务逻辑与数据访问解耦
- 易于测试（可Mock Repository）
- 易于切换数据库实现

**结构**（对标Go版的 `ModelRepository` interface）：
```python
from abc import ABC, abstractmethod

# 接口定义（对标Go的interface）
class ModelRepository(ABC):
    @abstractmethod
    async def create(self, model: Model) -> None: ...

    @abstractmethod
    async def get_by_id(self, id: str) -> Model | None: ...

    @abstractmethod
    async def update(self, model: Model) -> None: ...

    @abstractmethod
    async def delete(self, id: str) -> None: ...

    @abstractmethod
    async def list(self, filter: ModelFilter, pagination: Pagination) -> tuple[list[Model], int]: ...

# 实现（对标Go的GormModelRepository struct）
class SQLAlchemyModelRepository(ModelRepository):
    def __init__(self, session_factory):
        self._session_factory = session_factory

    async def create(self, model: Model) -> None:
        async with self._session_factory() as session:
            session.add(model)
            await session.commit()
```

**3. Alembic数据库迁移**（对标Go的GORM AutoMigrate + golang-migrate）
```bash
# 初始化Alembic
alembic init alembic

# 生成迁移脚本（自动检测模型变更）
alembic revision --autogenerate -m "create models table"

# 执行迁移
alembic upgrade head

# 回滚
alembic downgrade -1
```

**4. SQLAlchemy 2.0高级特性**

**关联查询**：
```python
class User(Base):
    __tablename__ = "users"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str]
    orders: Mapped[list["Order"]] = relationship(back_populates="user")  # 一对多
    roles: Mapped[list["Role"]] = relationship(secondary=user_roles)     # 多对多

# 预加载（解决N+1问题，对标GORM的Preload）
from sqlalchemy.orm import selectinload
stmt = select(User).options(selectinload(User.orders))
```

**事务**：
```python
# 方式1：context manager（对标GORM的Transaction闭包）
async with session.begin():
    session.add(user)
    session.add(profile)
    # 自动commit或rollback

# 方式2：手动控制
try:
    session.add(user)
    await session.commit()
except Exception:
    await session.rollback()
    raise
```

**实操任务**：
1. 设计Model Registry的数据库表结构
2. 创建SQLAlchemy 2.0模型定义
3. 实现Repository接口（ABC + SQLAlchemy实现）
4. 配置Alembic迁移
5. 实现Service层调用Repository
6. 编写单元测试（pytest-asyncio）

**数据库设计（Model Registry）**（与Go版完全对齐）：
```sql
-- models表：模型基本信息
CREATE TABLE models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) NOT NULL,
    framework VARCHAR(50) NOT NULL,  -- pytorch/tensorflow/onnx
    status VARCHAR(50) DEFAULT 'pending',
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

-- tags表、model_tags表、model_metadata表 与Go版完全相同
```

**检查点**：
- [ ] 数据库表结构符合范式
- [ ] SQLAlchemy模型定义正确
- [ ] Repository接口完整实现
- [ ] Alembic迁移脚本可正常执行
- [ ] 基础CRUD接口可调用
- [ ] 有基本的单元测试

---

## 阶段1里程碑

### 完成检查清单

**API Gateway服务**：
- [ ] 可独立启动和运行（uvicorn）
- [ ] 支持基础路由和中间件
- [ ] 有健康检查接口
- [ ] 日志输出正常（Loguru JSON）
- [ ] 配置管理正常工作

**Model Registry服务**：
- [ ] 可独立启动和运行
- [ ] 数据库连接正常（asyncpg）
- [ ] 支持模型的CRUD操作
- [ ] 有自动生成的OpenAPI文档

**基础设施**：
- [ ] Docker Compose可启动所有服务
- [ ] PostgreSQL容器正常运行
- [ ] 服务间可通过网络通信

### 可演示功能
1. 启动服务：`docker-compose up`
2. 访问API文档：http://localhost:8000/docs
3. 创建模型：POST /api/v1/models
4. 查询模型列表：GET /api/v1/models
5. 查看日志输出

### 下一步
完成阶段1后，你将掌握：
- ✅ Python项目结构组织（uv）
- ✅ FastAPI开发异步RESTful API
- ✅ SQLAlchemy 2.0数据库操作
- ✅ Docker容器化

准备好进入**阶段2：核心功能开发**了吗？

---

## 常见问题

### Q: uv sync 失败？

**A**:
1. 检查网络连接（PyPI源）
2. 尝试 `uv cache clean` 清理缓存
3. 使用 `uv sync -v` 查看详细错误

### Q: asyncpg连接PostgreSQL失败？
1. 确认PostgreSQL服务已启动
2. 检查端口号和密码
3. 确认数据库已创建

### Q: Alembic迁移报错？
1. 确认 `alembic.ini` 中的数据库URL正确
2. 运行 `alembic heads` 查看当前版本
3. 使用 `alembic stamp head` 重置版本

### Q: FastAPI启动端口被占用？
```bash
# 使用其他端口启动
uvicorn api_gateway.main:app --port 8001
```
