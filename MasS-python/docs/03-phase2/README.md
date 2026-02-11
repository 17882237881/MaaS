# 阶段2：核心功能开发

> 🚀 **阶段目标**：实现MaaS平台的核心业务流程，包括服务通信、缓存、认证、模型管理和推理服务

## 阶段总览

### 学习时间
**5周**（约25个工作日）

### 核心目标
1. 掌握gRPC服务间通信（grpcio）
2. 学会Redis缓存策略和设计
3. 实现JWT认证和RBAC权限控制
4. 处理大文件分片上传
5. 构建模型推理服务
6. 理解版本控制和依赖管理

### 最终产出
- ✅ API Gateway通过gRPC调用后端服务
- ✅ 完整的用户认证体系（JWT + RBAC）
- ✅ Redis多级缓存支持
- ✅ 模型文件上传和版本管理
- ✅ 可调用模型的推理接口

## 技术栈详解

### 1. gRPC与Protocol Buffers

**什么是gRPC？**
gRPC是Google开发的高性能RPC框架，基于HTTP/2和Protocol Buffers。Python版使用 `grpcio` 库。

**为什么使用gRPC？**
- **性能高**：基于HTTP/2，支持多路复用、头部压缩
- **跨语言**：支持Go/Java/Python/C++等10+语言（**Python版可直接与Go版互通**）
- **强类型**：使用Protocol Buffers定义接口，自动生成代码
- **流式通信**：支持客户端流、服务端流、双向流
- **异步支持**：grpcio提供 `grpc.aio` 异步API

**Protocol Buffers（protobuf）**
是一种语言中立、平台中立的数据序列化格式，比JSON更小更快。

> 💡 **关键优势**：Python版和Go版**共用同一份 `.proto` 文件**，只是生成的代码不同！

**示例对比**：
```protobuf
// protobuf定义（shared/proto/model.proto，与Go版完全相同）
syntax = "proto3";
package model;

service ModelService {
  rpc CreateModel(CreateModelRequest) returns (CreateModelResponse);
  rpc GetModel(GetModelRequest) returns (GetModelResponse);
  rpc ListModels(ListModelsRequest) returns (ListModelsResponse);
  rpc DeleteModel(DeleteModelRequest) returns (google.protobuf.Empty);
}
```

**Python生成代码**：
```bash
# 生成Python代码（对标 protoc --go_out --go-grpc_out）
python -m grpc_tools.protoc \
  -I./shared/proto \
  --python_out=./shared/proto \
  --grpc_python_out=./shared/proto \
  ./shared/proto/model.proto

# 生成文件：
# model_pb2.py       (对标 model.pb.go)
# model_pb2_grpc.py  (对标 model_grpc.pb.go)
```

### 2. Redis缓存

**什么是Redis？**
Redis是内存中的数据结构存储系统。Python版使用 `redis-py` 的异步模式。

**为什么要用缓存？**
- **减少数据库压力**：热点数据从内存读取
- **提高响应速度**：内存访问比磁盘快1000倍
- **降低系统成本**：减少数据库服务器配置

**缓存模式**：
1. **Cache-Aside（旁路缓存）**：应用先查缓存，没有则查DB并写入缓存
2. **Read-Through**：缓存服务自己处理DB查询
3. **Write-Through**：写数据时同时更新缓存和DB
4. **Write-Behind**：异步写DB，先写缓存

**常见问题**：
- **缓存穿透**：查询不存在的数据 → 缓存空值、布隆过滤器
- **缓存击穿**：热点key过期瞬间 → asyncio.Lock、逻辑过期
- **缓存雪崩**：大量key同时过期 → 随机过期时间、多级缓存

### 3. JWT认证

**什么是JWT？**
JSON Web Token是一种开放标准（RFC 7519），用于在网络应用间安全传输信息。Python版使用 `python-jose` 库。

**JWT结构**：
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
│                        │                             │
└── Header（头部）        └── Payload（负载）           └── Signature（签名）
```

**JWT vs Session**：
| 特性 | JWT | Session |
|------|-----|---------| 
| 存储位置 | 客户端 | 服务端 |
| 可扩展性 | 好（无状态） | 差（需要共享session） |
| 安全性 | 中（不能撤销） | 高（可随时删除） |
| 适用场景 | 微服务、移动端 | 传统Web应用 |

### 4. MinIO对象存储

**什么是对象存储？**
对象存储是一种数据存储架构，将数据作为对象管理。Python版使用 `minio-py` 库。

**为什么用MinIO？**
- 开源的S3兼容对象存储
- Python客户端API简洁
- 易于部署
- 支持分布式部署

### 5. 模型版本管理

**Semantic Versioning（语义化版本）**：
```
版本格式：主版本号.次版本号.修订号（MAJOR.MINOR.PATCH）

1.0.0
│ │ │
│ │ └── Patch：Bug修复，向后兼容
│ └──── Minor：新功能，向后兼容
└────── Major：重大变更，可能不兼容
```

---

## 节点详解

### 节点2.1：服务间通信（5天）

**学习目标**：
- 理解RPC vs HTTP REST
- 掌握Protocol Buffers语法
- 实现gRPC服务端（grpc.aio）和客户端
- 配置服务发现和负载均衡

**技术详解**：

**RPC vs REST对比**：
| 特性 | REST | RPC（gRPC） |
|------|------|-------------|
| 协议 | HTTP/1.1 | HTTP/2 |
| 格式 | JSON/XML | Protocol Buffers |
| 性能 | 一般 | 高（5-10倍） |
| 流式 | 不支持 | 支持 |
| Python库 | httpx / requests | grpcio |

**gRPC异步服务端示例**：
```python
import grpc
from grpc import aio
from shared.proto import model_pb2, model_pb2_grpc

class ModelServicer(model_pb2_grpc.ModelServiceServicer):
    async def CreateModel(self, request, context):
        # 业务逻辑...
        return model_pb2.CreateModelResponse(model=proto_model)

    async def GetModel(self, request, context):
        model = await self.service.get_model(request.id)
        if model is None:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("model not found")
            return model_pb2.GetModelResponse()
        return model_pb2.GetModelResponse(model=convert_to_proto(model))

# 启动gRPC服务器
async def serve():
    server = aio.server()
    model_pb2_grpc.add_ModelServiceServicer_to_server(ModelServicer(), server)
    server.add_insecure_port("[::]:9090")
    await server.start()
    await server.wait_for_termination()
```

**gRPC异步客户端示例**：
```python
async def create_grpc_client(address: str):
    channel = aio.insecure_channel(address)
    client = model_pb2_grpc.ModelServiceStub(channel)
    return client, channel

# 调用
client, channel = await create_grpc_client("localhost:9090")
response = await client.CreateModel(model_pb2.CreateModelRequest(
    name="bert-base", version="1.0.0", framework="pytorch"
))
```

**实操任务**：
1. 安装grpcio和grpcio-tools
2. 复用Go版的.proto文件，生成Python代码
3. 实现gRPC异步服务端（Model Registry）
4. 实现gRPC异步客户端（API Gateway调用）
5. 测试服务间调用
6. 添加错误处理和超时控制

**检查点**：
- [ ] .proto文件生成Python代码成功
- [ ] gRPC服务端能启动
- [ ] API Gateway能成功调用Model Registry
- [ ] 支持错误处理和超时控制

---

### 节点2.2：缓存层设计（4天）

**学习目标**：
- 掌握Redis基础数据类型
- 实现多级缓存架构
- 处理缓存常见问题
- 实现缓存预热和刷新

**Cache-Aside模式代码示例**：
```python
import redis.asyncio as redis
import json

class ModelService:
    def __init__(self, repo, redis_client: redis.Redis):
        self.repo = repo
        self.redis = redis_client

    async def get_model(self, model_id: str) -> Model | None:
        # 1. 查Redis缓存
        key = f"model:{model_id}"
        cached = await self.redis.get(key)
        if cached:
            return Model(**json.loads(cached))

        # 2. 查数据库
        model = await self.repo.get_by_id(model_id)
        if model is None:
            return None

        # 3. 回填缓存（1小时过期）
        await self.redis.set(key, model.model_dump_json(), ex=3600)
        return model
```

**实操任务**：
1. 安装并配置Redis（redis-py async）
2. 实现基础Redis操作封装
3. 实现模型查询缓存（Cache-Aside）
4. 处理缓存穿透保护
5. 编写缓存性能测试

**检查点**：
- [ ] Redis连接正常
- [ ] 缓存命中和未命中逻辑正确
- [ ] 缓存过期策略合理
- [ ] 缓存穿透/击穿问题已处理

---

### 节点2.3：认证与授权（5天）

**学习目标**：
- 实现JWT认证流程
- 设计RBAC权限模型
- 实现API密钥管理
- 添加请求鉴权中间件

**JWT认证流程**：
```
1. 用户登录 → POST /auth/login → 生成JWT → 返回Token
2. 访问受保护资源 → Header: Authorization: Bearer <token> → 验证JWT
3. Token刷新 → POST /auth/refresh → 生成新Token
```

**Python JWT示例**：
```python
from jose import jwt, JWTError
from datetime import datetime, timedelta

SECRET_KEY = "your-secret-key"

def create_access_token(data: dict, expires_delta: timedelta = timedelta(hours=24)):
    to_encode = data.copy()
    to_encode["exp"] = datetime.utcnow() + expires_delta
    return jwt.encode(to_encode, SECRET_KEY, algorithm="HS256")

def verify_token(token: str) -> dict:
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        return payload
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")
```

**FastAPI依赖注入实现认证**：
```python
from fastapi import Depends
from fastapi.security import HTTPBearer

security = HTTPBearer()

async def get_current_user(credentials = Depends(security)):
    payload = verify_token(credentials.credentials)
    return payload

@app.get("/api/v1/users/me")
async def get_me(user = Depends(get_current_user)):
    return user
```

**Casbin RBAC示例（pycasbin）**：
```python
import casbin

enforcer = casbin.Enforcer("model.conf", "policy.csv")

# 检查权限
if enforcer.enforce("alice", "model", "read"):
    # 允许访问
    pass
```

**实操任务**：
1. 实现用户注册和登录接口
2. 集成python-jose JWT生成和验证
3. 实现FastAPI认证依赖注入
4. 集成pycasbin权限控制
5. 实现Token刷新机制

**检查点**：
- [ ] 用户能注册和登录
- [ ] JWT Token正确生成和验证
- [ ] 未认证请求返回401
- [ ] 无权限请求返回403
- [ ] RBAC权限控制生效

---

### 节点2.4：模型上传与存储（5天）

**学习目标**：
- 实现大文件分片上传
- 集成MinIO对象存储（minio-py）
- 实现断点续传
- 添加文件校验

**MinIO集成示例**：
```python
from minio import Minio

client = Minio(
    "localhost:9000",
    access_key="access-key",
    secret_key="secret-key",
    secure=False,
)

# 上传文件
client.put_object(
    bucket_name="models",
    object_name="user1/model-v1.pkl",
    data=file_data,
    length=file_size,
)

# 生成下载URL（1小时有效）
url = client.presigned_get_object("models", "user1/model-v1.pkl", expires=timedelta(hours=1))
```

**FastAPI文件上传**：
```python
from fastapi import UploadFile, File

@app.post("/api/v1/models/{model_id}/upload")
async def upload_model(model_id: str, file: UploadFile = File(...)):
    content = await file.read()
    # 上传到MinIO...
```

**实操任务**：
1. 部署MinIO服务
2. 实现分片上传接口
3. 实现断点续传功能
4. 添加文件MD5校验
5. 限制文件类型和大小

**检查点**：
- [ ] 小文件上传正常
- [ ] 大文件分片上传正常
- [ ] 断点续传功能工作
- [ ] 文件存储在MinIO中

---

### 节点2.5：模型版本管理（4天）

**学习目标**：
- 实现语义化版本控制
- 设计版本元数据结构
- 实现版本回滚
- 添加版本间差异比较

**版本数据模型**：
```python
class ModelVersion(Base):
    __tablename__ = "model_versions"

    id: Mapped[str] = mapped_column(primary_key=True, default=uuid4)
    model_id: Mapped[str] = mapped_column(ForeignKey("models.id"))
    version: Mapped[str]           # 版本号（如 1.0.0）
    status: Mapped[str]            # draft/published/deprecated
    change_log: Mapped[str | None]
    size: Mapped[int] = mapped_column(default=0)
    checksum: Mapped[str | None]
    storage_path: Mapped[str | None]
    created_by: Mapped[str]
    is_latest: Mapped[bool] = mapped_column(default=False)
    created_at: Mapped[datetime] = mapped_column(server_default=func.now())
    published_at: Mapped[datetime | None]
```

**实操任务**：
1. 设计版本数据库表结构
2. 实现版本创建和查询接口
3. 实现latest标签自动更新
4. 实现版本回滚功能
5. 添加版本变更日志

**检查点**：
- [ ] 版本号符合SemVer规范
- [ ] latest标签指向正确
- [ ] 版本历史可查询
- [ ] 版本回滚功能正常

---

### 节点2.6：推理服务基础（5天）

**学习目标**：
- 实现同步推理接口
- 设计请求队列和并发控制
- 实现超时和取消机制
- 添加推理结果缓存

**并发控制模式**：

**1. asyncio.Semaphore（对标Go的channel信号量）**：
```python
import asyncio

class InferenceLimiter:
    def __init__(self, max_concurrent: int):
        self._semaphore = asyncio.Semaphore(max_concurrent)

    async def run_inference(self, model_id: str, input_data: dict):
        async with self._semaphore:
            # 执行推理（最多max_concurrent个并发）
            return await self._do_inference(model_id, input_data)

limiter = InferenceLimiter(100)  # 最多100个并发
```

**2. 超时控制（对标Go的context.WithTimeout）**：
```python
import asyncio

async def inference_with_timeout(model_id: str, input_data: dict):
    try:
        result = await asyncio.wait_for(
            do_inference(model_id, input_data),
            timeout=30.0  # 30秒超时
        )
        return result
    except asyncio.TimeoutError:
        raise HTTPException(408, "inference timeout")
```

**实操任务**：
1. 创建推理服务框架（模拟实现）
2. 实现同步推理接口
3. 添加并发控制（Semaphore）
4. 实现超时和取消机制
5. 添加推理结果缓存
6. 记录推理指标

**检查点**：
- [ ] 推理接口可调用
- [ ] 并发控制生效
- [ ] 超时机制正常工作
- [ ] 推理结果正确返回

---

## 阶段2里程碑

### 完成检查清单

**服务间通信**：
- [ ] API Gateway通过gRPC调用Model Registry
- [ ] Protocol Buffers定义完整
- [ ] 支持错误处理和重试

**缓存**：
- [ ] Redis集成完成
- [ ] 模型查询有缓存
- [ ] 缓存命中率>50%

**认证授权**：
- [ ] JWT认证工作正常
- [ ] RBAC权限控制生效
- [ ] 用户能正常注册登录

**模型管理**：
- [ ] 模型文件可上传
- [ ] 版本管理完整
- [ ] 文件存储在MinIO

**推理服务**：
- [ ] 推理接口可调用
- [ ] 并发控制有效
- [ ] 超时机制正常

### 可演示功能
1. 用户注册登录获取JWT
2. 创建模型并上传文件
3. 发布模型版本
4. 调用推理接口
5. 查看缓存命中率

### 下一步
进入**阶段3：企业级特性**，学习：
- 消息队列（Kafka）
- Kubernetes部署
- 限流熔断
- 分布式事务
- 监控告警

---

**继续学习**：[阶段3文档](../04-phase3/README.md)
