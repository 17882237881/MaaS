# MaaS Python vs Go RPC 逻辑对比

## Python 实现步骤
1. **定义 RPC 协议**：`shared/proto/model.proto` 定义 `ModelService` 的 10 个 RPC 方法（CRUD + 状态 + tags + metadata）。
2. **启动 gRPC 服务端**：`model_registry/main.py` 创建 `grpc.aio.server()`，注册 `ModelServiceServicer` 并监听配置端口。
3. **服务端请求处理**：`model_registry/internal/server/grpc_server.py` 为每个 RPC 方法创建 session 与 `ModelService`，完成请求到领域模型转换，调用服务层，异常映射为 gRPC 状态码。
4. **服务层业务逻辑**：`model_registry/internal/service/model_service.py` 执行输入校验（框架、UUID）、CRUD、tags diff 更新、metadata 更新等。
5. **网关作为 gRPC 客户端**：`api_gateway/internal/client/grpc_client.py` 通过 `grpc.aio.insecure_channel` 建立连接并封装各 RPC。
6. **HTTP -> gRPC 转发**：`api_gateway/internal/handler/model_handler.py` 将 HTTP 请求转 protobuf，调用 gRPC client，并把 gRPC 错误映射为 HTTP 状态码。

## Go 实现步骤
1. **定义 RPC 协议**：`shared/proto/model.proto` 与 Python 使用同一协议。
2. **启动 model-registry gRPC 服务端**：`model-registry/cmd/main.go` 创建 `grpc.NewServer()`，注册 `GRPCServer`，监听 `:9090`。
3. **服务端请求处理**：`model-registry/internal/grpc/server.go` 处理 Create/Get/List/Update/Delete/UpdateStatus 六个方法，调用 service 层并做基础错误映射。
4. **服务层业务逻辑**：`model-registry/internal/service/model_service.go` 实现核心 CRUD 和状态更新，调用 repository。
5. **网关创建 gRPC 连接**：`api-gateway/pkg/grpc/client.go` `grpc.Dial` 建立连接并封装所有 proto 声明的方法。
6. **HTTP -> gRPC 转发**：`api-gateway/internal/handler/handler.go` 处理模型相关 HTTP 路由并调用 `ModelServiceClient`。

## 一致性结论
- **协议层一致**：两端都基于相同 `model.proto`，RPC 方法定义一致。
- **实现层不完全一致**：
  - Python 服务端实现了 proto 中全部 10 个 RPC；Go 服务端目前仅实现 6 个，缺少 `AddModelTags/RemoveModelTags/SetModelMetadata/GetModelMetadata`，会走未实现返回。
  - Go 的 `ListModels` gRPC 层未将 `tags`、`is_public` 从请求传入 service 过滤条件，Python 有完整映射。
  - Go `UpdateModel` 中“更新 tags”逻辑疑似错误地调用了 `SetMetadata`（清空 metadata），Python 使用 tags 差集更新。
  - Python 对 gRPC 错误到 HTTP 错误码映射更细（404/409/400）；Go API Gateway 当前主要按内部错误处理。

总体：**两者在 RPC 契约上保持一致，但在服务端方法覆盖、字段映射和错误处理细节上存在明显不一致，Go 版本实现完整度低于 Python。**
