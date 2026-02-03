# Node 01 / Step 03 - 配置与日志规范

## 背景/目标
- 统一服务的配置来源（环境变量）
- 统一 Hertz/Kitex 的日志级别与输出

## 范围与非目标
- 范围：基础配置加载、日志初始化
- 非目标：配置中心接入、结构化日志与链路追踪

## 设计与接口
- 配置模型：`internal/config.AppConfig`
- 日志入口：`internal/logging.Init(level, out)`
- 运行参数（环境变量）：
  - `APP_ENV` 默认 `dev`
  - `HTTP_ADDR` 默认 `:8080`
  - `RPC_ADDR` 默认 `:9090`
  - `LOG_LEVEL` 默认 `info`

## 实现步骤
1. 新增 `internal/config/config.go` 定义配置结构与加载函数
2. 新增 `internal/logging/logging.go` 适配 Hertz/Kitex 日志
3. 引入 CloudWeGo 依赖：`github.com/cloudwego/hertz`、`github.com/cloudwego/kitex`

## 测试与验收
- `go test ./...` 通过（当前无单测）
- `go vet ./...` 通过
- 通过设置 `LOG_LEVEL=debug` 能调整日志级别

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `chore: add config and logging foundation`
