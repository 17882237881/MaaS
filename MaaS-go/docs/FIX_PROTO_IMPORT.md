# 修复 protobuf 导入问题

## 问题原因
Go 编译器无法找到 `maas-platform/shared/proto/model` 包，即使文件存在。

## 解决步骤

### 1. 清理 Go 缓存
```bash
cd D:\code\MaaS\MaaS-go
go clean -cache
go clean -modcache
```

### 2. 重新生成 protobuf 代码
```bash
# 删除旧的生成文件
del shared\proto\model.pb.go
del shared\proto\model_grpc.pb.go

# 重新生成
scripts\generate-proto.bat
```

### 3. 确保生成成功
检查文件是否存在：
```bash
dir shared\proto\*.go
```

应该看到：
- model.pb.go
- model_grpc.pb.go

### 4. 验证编译
```bash
go build ./...
```

## 如果还是报错

可能是 protobuf 生成不完整，请手动检查：

1. 检查 `shared/proto/model.pb.go` 第7行应该是：
   ```go
   package modelpb
   ```

2. 检查 `shared/proto/model_grpc.pb.go` 第7行应该是：
   ```go
   package modelpb
   ```

3. 如果 package 名不对，需要重新运行 protoc 生成。

## 临时绕过方案

如果 protobuf 生成有问题，可以暂时不使用 gRPC，直接通过 HTTP 调用 Model Registry：

```go
// 在 api-gateway 中，修改 handler，直接调用 Model Registry 的 HTTP API
// 而不是通过 gRPC
```

建议先尝试重新生成 protobuf 代码！
