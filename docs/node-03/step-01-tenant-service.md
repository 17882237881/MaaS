# Node 03 / Step 01 - Tenant Service（基础能力）

## 背景/目标
- 提供租户基本管理能力（Create/Get/List/Delete）
- 形成后续 Kitex RPC 与控制面统一的领域模型

## 范围与非目标
- 范围：Tenant 领域模型、Hertz HTTP API、内存存储
- 非目标：持久化数据库、鉴权、配额、审计

## 设计与接口
- Proto：`api/proto/tenant.proto`
- HTTP API：
  - `POST /tenants`
  - `GET /tenants/{id}`
  - `GET /tenants`
  - `DELETE /tenants/{id}`
- 服务入口：`cmd/tenant-service/main.go`

## 实现步骤
1. 新增 Tenant Proto 定义（后续用于 Kitex 生成）
2. 新增领域模型与内存存储实现
3. 新增 Hertz 路由与处理函数
4. 初始化服务入口，接入配置与日志

## 测试与验收
- `go test ./...` 通过（当前无单测）
- 启动服务后：
  - 创建租户成功返回 JSON
  - 查询/删除不存在租户返回 404

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `feat: add tenant service skeleton`
