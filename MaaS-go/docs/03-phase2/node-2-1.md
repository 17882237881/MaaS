# 节点2.1：gRPC服务间通信

> 📅 **学习时间**：5天  
> 🎯 **目标**：实现API Gateway到Model Registry的gRPC通信

## 本节你将学到

1. Protocol Buffers定义接口
2. gRPC服务端实现
3. gRPC客户端调用
4. 连接池和超时配置
5. 错误处理

---

## 技术详解

### 1. 为什么使用gRPC？

**gRPC vs REST**：

| 特性 | REST (HTTP/JSON) | gRPC (HTTP/2 + Protobuf) |
|------|------------------|---------------------------|
| 协议 | HTTP/1.1 | HTTP/2 |
| 格式 | JSON (文本) | Protobuf (二进制) |
| 性能 | 一般 | 高（5-10倍） |
| 类型 | 弱类型 | 强类型 |
| 流式 | 不支持 | 支持 |
| 浏览器 | 原生支持 | 需要gRPC-Web |

**适用场景**：
- 微服务内部通信
- 高性能要求
- 多语言环境

### 2. Protocol Buffers

**什么是Protobuf？**
语言中立、平台中立的数据序列化格式，定义在.proto文件中。

**示例**：
```protobuf
syntax = "proto3";

message Model {
  string id = 1;
  string name = 2;
  string version = 3;
}

service ModelService {
  rpc GetModel(GetModelRequest) returns (GetModelResponse);
}
```

**优势**：
- 体积小（二进制编码）
- 速度快（解析快）
- 自动生成代码

### 3. gRPC四种服务类型

**1. Unary RPC（一元RPC）**：
```
客户端 → 单个请求 → 服务端 → 单个响应
```

**2. Server Streaming RPC（服务端流）**：
```
客户端 → 单个请求 → 服务端 → 多个响应
```

**3. Client Streaming RPC（客户端流）**：
```
客户端 → 多个请求 → 服务端 → 单个响应
```

**4. Bidirectional Streaming RPC（双向流）**：
```
客户端 ←→ 多个消息 ←→ 服务端
```

### 4. gRPC连接配置

**连接选项**：
```go
opts := []grpc.DialOption{
    // 不安全的连接（开发环境）
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    
    // 保活配置
    grpc.WithKeepaliveParams(keepalive.ClientParameters{
        Time:                10 * time.Second,
        Timeout:             20 * time.Second,
        PermitWithoutStream: true,
    }),
    
    // 超时配置
    grpc.WithTimeout(30 * time.Second),
}
```

---

## 实操任务

### 任务1：定义Proto接口

创建 `shared/proto/model.proto`：
- 定义Model消息
- 定义CRUD服务接口
- 包含标签和元数据操作

### 任务2：生成Go代码

使用protoc生成：
```bash
protoc --go_out=. --go-grpc_out=. model.proto
```

生成两个文件：
- `model.pb.go` - 消息类型
- `model_grpc.pb.go` - gRPC接口

### 任务3：实现gRPC服务端

在Model Registry中：
- 实现ModelServiceServer接口
- 转换内部模型到Protobuf模型
- 处理gRPC错误码

### 任务4：实现gRPC客户端

在API Gateway中：
- 创建gRPC客户端
- 封装调用方法
- 处理连接管理

### 任务5：更新服务启动

更新main.go：
- Model Registry启动gRPC服务器
- API Gateway初始化gRPC客户端

---

## 代码变更记录

### 提交信息
```
feat(phase2/node2.1): implement gRPC service communication

- Add Protocol Buffers definition (model.proto)
- Generate gRPC Go code from proto
- Implement gRPC server in Model Registry
- Implement gRPC client in API Gateway
- Add connection pool and timeout configuration
```

### 新增文件

#### 1. shared/proto/model.proto
**新增文件**
Protocol Buffers定义文件，包含：
- Model消息定义（所有字段）
- CRUD请求/响应消息
- ModelService服务定义（10个RPC方法）

#### 2. shared/proto/model_grpc.pb.go
**生成的文件**
gRPC接口代码，包含：
- ModelServiceClient接口
- ModelServiceServer接口
- 客户端实现
- 服务端注册

#### 3. model-registry/internal/grpc/server.go
**新增文件**
gRPC服务端实现：
- GRPCServer结构体
- 实现所有10个RPC方法
- 错误码映射（NotFound→codes.NotFound）
- 模型转换函数

#### 4. api-gateway/pkg/grpc/client.go
**新增文件**
gRPC客户端封装：
- Client结构体
- 连接管理（Dial、Close）
- 所有RPC方法的便捷调用
- 保活配置

### 修改的文件

#### model-registry/cmd/main.go
**更新**
- 添加 gRPC 服务器启动（端口 9090）
- 在 `startGRPCServer` 函数中创建并注册 ModelService
- 使用 goroutine 并行启动 gRPC 和 HTTP 服务器
- 导入 `google.golang.org/grpc` 和 protobuf 生成的代码

#### api-gateway/cmd/main.go
**更新**
- 添加 gRPC 客户端初始化
- 使用 `grpc.NewClient` 连接到 Model Registry 服务
- 创建 `ModelServiceClient` 包装 gRPC 调用
- 将 client 注入到 handler 中
- 添加连接日志和错误处理
- 程序退出时关闭 gRPC 连接

#### api-gateway/internal/service/model_client.go
**新增文件**
- 创建 `ModelServiceClient` 封装 gRPC 调用
- 提供 `CreateModel`, `GetModel`, `ListModels`, `DeleteModel` 方法
- 处理错误日志记录

#### api-gateway/internal/handler/handler.go
**更新**
- 添加 `modelClient` 字段到 Handler 结构体
- 更新 `New()` 函数接收 modelClient 参数
- 修改所有 model 相关 handler 方法
- 通过 gRPC 调用 Model Registry 服务（替换硬编码 mock 数据）
- 实现 `convertProtoModelToResponse` 转换函数

---

## 验证步骤

### 1. 生成Proto代码

```bash
# 安装protoc-gen-go和protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
cd shared/proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       model.proto
```

### 2. 启动gRPC服务端

```bash
cd model-registry
go run cmd/main.go

# 应该看到：
# {"msg":"gRPC server starting","port":9090}
```

### 3. 测试gRPC调用

```bash
# 使用grpcurl测试
grpcurl -plaintext localhost:9090 list model.ModelService
grpcurl -plaintext localhost:9090 model.ModelService/GetModel
```

### 4. API Gateway调用

```bash
cd api-gateway
go run cmd/main.go

# 应该看到：
# {"msg":"Connecting to Model Registry gRPC service...","address":"localhost:9090"}
# {"msg":"Connected to Model Registry gRPC service"}

# 调用API Gateway的接口，它会通过gRPC调用Model Registry
curl -X POST http://localhost:8080/api/v1/models \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-model",
    "version": "1.0.0",
    "framework": "pytorch"
  }'

# 查询列表（通过gRPC调用Model Registry）
curl http://localhost:8080/api/v1/models

# 查询单个模型
curl http://localhost:8080/api/v1/models/{id}
```

### 5. 验证调用链

确保调用链路完整：
1. API Gateway 接收到 HTTP 请求
2. Handler 调用 `modelClient.CreateModel()` 等方法
3. `ModelServiceClient` 通过 gRPC 发送请求到 Model Registry
4. Model Registry 操作数据库
5. 返回结果通过 gRPC → HTTP → 用户

---

## 检查清单

- [ ] Proto文件定义完整
- [ ] Go代码生成成功
- [ ] gRPC服务端可启动（端口9090）
- [ ] gRPC客户端可连接
- [ ] API Gateway handler 使用 gRPC 调用（非 mock 数据）
- [ ] ModelServiceClient 封装完整
- [ ] 方法调用正常
- [ ] 错误处理正确
- [ ] 数据流转：HTTP → gRPC → Database → gRPC → HTTP

---

## 下一步

完成本节点后，服务间通信已完成。接下来进入：

**节点2.2：Redis缓存层设计** → [继续学习](./node-2-2.md)

在那里你将：
- 集成Redis缓存
- 实现多级缓存
- 处理缓存穿透/击穿/雪崩

---

## 参考资源

- [gRPC官方文档](https://grpc.io/docs/)
- [Protocol Buffers指南](https://developers.google.com/protocol-buffers)
- [Go gRPC教程](https://grpc.io/docs/languages/go/)
- [gRPC vs REST](https://medium.com/@EmperorRXF/evaluating-performance-of-rest-vs-grpc-1b8bdf0b228d)
