# 节点2.1：gRPC服务间通信

> 📅 **学习时间**：5天  
> 🎯 **目标**：实现 API Gateway 到 Model Registry 的 gRPC 通信，并完成多语言协议对齐

## 本节你将学到

1. Protocol Buffers 定义接口
2. gRPC 服务端实现
3. gRPC 客户端调用
4. 连接池与超时配置
5. 错误处理与状态码映射

---

## 技术详解

### 1. 为什么使用 gRPC？

| 特性 | REST (HTTP/JSON) | gRPC (HTTP/2 + Protobuf) |
|------|------------------|---------------------------|
| 协议 | HTTP/1.1 | HTTP/2 |
| 格式 | JSON (文本) | Protobuf (二进制) |
| 性能 | 一般 | 高（5-10倍） |
| 类型 | 弱类型 | 强类型 |
| 流式 | 不支持 | 支持 |
| 浏览器 | 原生支持 | 需要 gRPC-Web |

**适用场景**：
- 微服务内部通信
- 高性能要求
- 多语言环境（Go / Python）

---

## 实操任务

### 任务1：定义 Proto 接口

创建 `shared/proto/model.proto`：
- `ModelService` 服务定义（10个 RPC）
- `Model` 与各类 Request/Response
- 包名统一为 `model`

> ✅ 已与 `MaaS-go/shared/proto/model.proto` 对齐

### 任务2：生成 Python 代码

```bash
# 在 MasS-python 目录下执行
.venv\Scripts\python -m grpc_tools.protoc \
  -I shared/proto \
  --python_out=shared/proto \
  --grpc_python_out=shared/proto \
  shared/proto/model.proto
```

生成文件：
- `shared/proto/model_pb2.py`
- `shared/proto/model_pb2_grpc.py`

### 任务3：实现 gRPC 服务端

在 Model Registry 中：
- 文件：`model_registry/internal/server/grpc_server.py`
- 实现 `ModelServiceServicer`
- 连接 Repository & Service
- 错误码映射：
  - NotFound → `codes.NOT_FOUND`
  - InvalidArgument → `codes.INVALID_ARGUMENT`
  - Duplicate → `codes.ALREADY_EXISTS`

### 任务4：实现 gRPC 客户端

在 API Gateway 中：
- 文件：`api_gateway/internal/client/grpc_client.py`
- 封装 `ModelServiceClient`
- 提供 `create/get/list/update/delete` 等方法

### 任务5：更新服务启动

- `model_registry/main.py`：启动 gRPC Server（默认端口 9090）
- `api_gateway/main.py`：初始化 gRPC Client

---

## 验证步骤

### 1. 启动 Model Registry gRPC 服务端

```bash
uv run python model_registry/main.py
```

### 2. 启动 API Gateway

```bash
uv run python api_gateway/main.py
```

### 3. 调用 API Gateway

```bash
curl -X POST http://localhost:8000/api/v1/models \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-model",
    "version": "1.0.0",
    "framework": "pytorch",
    "tags": ["cv"],
    "metadata": {"source": "demo"}
  }'
```

---

## 检查清单

- [ ] Proto 文件定义完整
- [ ] Python gRPC 代码生成成功
- [ ] gRPC 服务端可启动（端口 9090）
- [ ] gRPC 客户端可连接
- [ ] API Gateway 通过 gRPC 调用 Model Registry
- [ ] 错误码映射正确

---

## 下一步

完成本节点后，进入：

**节点2.2：Redis 缓存层设计** → [继续学习](./node-2-2.md)
