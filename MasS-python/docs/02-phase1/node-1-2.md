# 阶段1 节点2：API Gateway核心

> 🏗️ **节点目标**：掌握 FastAPI 核心用法，完成 API Gateway 的中间件、健康检查与基础 CRUD 示例。

## 1. 目录结构回顾

本节点主要使用以下目录（与 Go 版对齐）：

```
api_gateway/
├── main.py
└── internal/
    ├── handler/
    ├── middleware/
    ├── model/
    ├── repository/
    ├── router/
    └── service/
```

---

## 2. 中间件实现（Logger / Recovery / RequestID / CORS）

### 2.1 RequestID 中间件
**作用**：为每个请求生成 `X-Request-ID`，方便链路追踪。

文件：`api_gateway/internal/middleware/request_id.py`

```python
from uuid import uuid4
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import Response

class RequestIDMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next) -> Response:
        request_id = request.headers.get("X-Request-ID") or str(uuid4())
        request.state.request_id = request_id
        response = await call_next(request)
        response.headers["X-Request-ID"] = request_id
        return response
```

### 2.2 Logger 中间件
**作用**：记录请求方法、路径、耗时、状态码等。

文件：`api_gateway/internal/middleware/logger.py`

```python
class LoggerMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next) -> Response:
        start = time.perf_counter()
        response = await call_next(request)
        duration_ms = (time.perf_counter() - start) * 1000
        request_id = getattr(request.state, "request_id", "")
        logger.bind(
            request_id=request_id,
            method=request.method,
            path=request.url.path,
            status=response.status_code,
            duration_ms=round(duration_ms, 2),
        ).info("Request completed")
        return response
```

### 2.3 Recovery 中间件
**作用**：捕获未处理异常，返回统一错误格式。

文件：`api_gateway/internal/middleware/recovery.py`

```python
class RecoveryMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next) -> Response:
        try:
            return await call_next(request)
        except Exception:
            logger.exception("Unhandled exception")
            return JSONResponse(
                status_code=500,
                content={"error": "Internal server error", "code": "INTERNAL_ERROR"},
            )
```

### 2.4 CORS
FastAPI 内置 `CORSMiddleware`：

```python
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)
```

---

## 3. CRUD 示例（用户）

### 3.1 数据模型
文件：`api_gateway/internal/model/user.py`

```python
class UserCreate(BaseModel):
    name: str
    email: str

class UserUpdate(BaseModel):
    name: str | None
    email: str | None

class UserResponse(BaseModel):
    id: str
    name: str
    email: str
    created_at: datetime
    updated_at: datetime
```

### 3.2 Repository + Service
文件：
- `api_gateway/internal/repository/user_repository.py`
- `api_gateway/internal/service/user_service.py`

职责：
- Repository：提供数据访问接口（此处用内存字典模拟）
- Service：封装业务逻辑（校验唯一性、生成ID）

### 3.3 Handler + Router
文件：
- `api_gateway/internal/handler/user_handler.py`
- `api_gateway/internal/router/router.py`

路由示例：
```
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

---

## 4. 入口文件组装

文件：`api_gateway/main.py`

```python
app = FastAPI(...)
app.add_middleware(RecoveryMiddleware)
app.add_middleware(LoggerMiddleware)
app.add_middleware(RequestIDMiddleware)
app.add_middleware(CORSMiddleware, allow_origins=["*"])
register_routes(app)
```

---

## 5. 运行与验证

### 5.1 启动服务
```bash
uv run uvicorn api_gateway.main:app --reload --port 8000
```

### 5.2 验证接口
1. 访问健康检查：
```
GET http://localhost:8000/health
```

2. 创建用户：
```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

3. 获取用户列表：
```bash
curl http://localhost:8000/api/v1/users
```

4. API 文档：
```
http://localhost:8000/docs
```

---

## ✅ 完成检查清单

- [ ] FastAPI 应用可启动
- [ ] 中间件生效（日志、请求ID、异常捕获）
- [ ] /health 接口返回正常
- [ ] /docs 自动生成 OpenAPI 文档
- [ ] CRUD 示例接口可正常调用

---

🎉 **完成节点1.2后，你已经具备搭建 API Gateway 核心能力！**
