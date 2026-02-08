# Protocol Buffers 安装指南（Windows）

## 步骤1：下载 protoc

1. 访问 https://github.com/protocolbuffers/protobuf/releases
2. 下载最新版本的 `protoc-<version>-win64.zip`（例如：`protoc-24.4-win64.zip`）
3. 解压到任意目录（例如：`C:\protoc`）

## 步骤2：添加到环境变量

1. 右键点击"此电脑" → 属性 → 高级系统设置
2. 点击"环境变量"
3. 在"系统变量"中找到 `Path`，点击"编辑"
4. 点击"新建"，添加 protoc 的 bin 目录路径（例如：`C:\protoc\bin`）
5. 点击"确定"保存

## 步骤3：验证安装

打开新的命令提示符（CMD）或 PowerShell：

```cmd
protoc --version
```

应该显示类似：`libprotoc 24.4`

## 步骤4：安装 Go 插件

在命令提示符中执行：

```cmd
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 步骤5：确保插件在 PATH 中

将 Go 的 bin 目录添加到环境变量：

1. 找到你的 Go bin 目录（通常是 `%USERPROFILE%\go\bin` 或 `%GOPATH%\bin`）
2. 将其添加到系统的 PATH 环境变量中

验证插件安装：

```cmd
protoc-gen-go --version
protoc-gen-go-grpc --version
```

## 步骤6：生成代码

在项目根目录执行：

```cmd
cd D:\code\MaaS\MaaS-go
scripts\generate-proto.bat
```

或者手动执行：

```cmd
cd D:\code\MaaS\MaaS-go
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative shared/proto/model.proto
```

## 预期输出

成功后应该看到：

```
Generating Protocol Buffers Go code...
Generating Go code from model.proto...

✅ Successfully generated protobuf code!
Generated files:
  - shared/proto/model.pb.go
  - shared/proto/model_grpc.pb.go
```

## 常见问题

### 问题1：'protoc' 不是内部或外部命令

**解决**：检查 PATH 环境变量是否正确设置，确保添加了 `C:\protoc\bin`

### 问题2：protoc-gen-go: program not found

**解决**：确保 `protoc-gen-go.exe` 在你的 PATH 中，可以尝试：

```cmd
where protoc-gen-go
```

如果没找到，将 `%USERPROFILE%\go\bin` 添加到 PATH

### 问题3：Import 错误

如果生成后导入报错，确保在代码中正确导入：

```go
import modelpb "maas-platform/shared/proto/model"
```

## 验证编译

生成代码后，运行：

```cmd
cd D:\code\MaaS\MaaS-go
go build ./...
```

如果没有错误，说明 protobuf 代码生成成功！
