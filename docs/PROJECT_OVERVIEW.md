# MaaS 平台项目总览（功能、技术栈与节点规划）

> 目标：这份文档是项目的“全景说明书”。阅读后应能理解：平台要做什么、技术栈如何选型、当前代码结构与模块划分、以及完整开发节点/步骤及每一步实现的功能。

## 1. 项目定位与目标
MaaS（Model-as-a-Service）平台面向内部业务线，提供统一的模型接入、推理服务、异步生成、计量计费、审计合规与可观测能力。平台以 CloudWeGo 技术栈为核心，强调高并发、可扩展、多租户与 GPU 资源调度。

## 2. 平台核心功能清单
- **模型接入与版本管理**：模型注册、版本发布、运行时参数管理
- **推理服务（同步 + 异步）**：统一 REST API 与内部 RPC 编排
- **多租户与鉴权**：租户、项目、API Key 管理，逻辑隔离
- **配额与限流**：按租户/模型维度执行配额控制
- **计量计费**：token/时长/像素维度的用量统计与账单
- **审计合规**：审计日志与不可变留存
- **可观测与稳定性**：监控、日志、SLO、灰度发布

## 3. 技术栈（当前与规划）
**语言与框架**
- Go 语言
- CloudWeGo：Hertz（HTTP）、Kitex（RPC）、Protobuf（IDL）

**基础设施与中间件**
- Kubernetes（阿里云 ACK）
- OSS（对象存储）
- PostgreSQL（元数据）
- Redis（缓存/限流计数）
- Kafka（异步任务/计量事件流）

**可观测**
- Prometheus + Grafana（已规划 Helm 基线）
- Loki/Promtail（可选）

**部署与工程化**
- Terraform（基础设施 IaC）
- Helm（部署打包）
- GitHub Actions（CI）

## 4. 当前代码结构（含目录职责）

```
.
├─ .github/workflows/ci.yml        # CI：go mod tidy / test / vet
├─ api/
│  └─ proto/tenant.proto           # Tenant 领域 IDL（后续用于 Kitex 生成）
├─ cmd/
│  └─ tenant-service/main.go        # 租户服务入口（Hertz）
├─ deploy/
│  └─ helm/
│     ├─ maas-middleware/           # PostgreSQL/Redis/Kafka umbrella chart
│     └─ maas-observability/        # 监控栈 umbrella chart
├─ docs/
│  ├─ node-01/                      # 节点1文档
│  ├─ node-02/                      # 节点2文档
│  └─ node-03/                      # 节点3文档（已开始）
├─ infra/
│  └─ terraform/                    # ACK、VPC、OSS、GPU NodePool 基线
├─ internal/
│  ├─ config/                       # 统一配置读取
│  ├─ logging/                      # Hertz/Kitex 日志适配
│  └─ tenant/                       # Tenant 领域模型 + 服务逻辑
├─ scripts/                         # 本地脚本（待补充）
├─ Makefile
├─ README.md
├─ go.mod / go.sum
└─ .gitignore
```

### 4.1 `internal/` 目录细节
- `internal/config/config.go`
  - 负责读取环境变量并提供默认值（`APP_ENV`, `HTTP_ADDR`, `RPC_ADDR`, `LOG_LEVEL`）
- `internal/logging/logging.go`
  - 统一初始化 Hertz 与 Kitex 日志级别与输出
- `internal/tenant/*`
  - `model.go`：Tenant 领域模型
  - `store.go`：内存存储（InMemoryStore）
  - `service.go`：业务逻辑（Create/Get/List/Delete）
  - `handler.go`：Hertz 路由与 HTTP 处理函数

### 4.2 `cmd/` 入口
- `cmd/tenant-service/main.go`
  - 加载配置与日志
  - 初始化 Hertz 服务器
  - 注册 Tenant 路由

### 4.3 `infra/terraform/` 基础设施 IaC
- `main.tf`：VPC、VSwitch、ACK 集群、GPU 节点池
- `oss.tf`：OSS Bucket
- `variables.tf`：可配置变量
- `outputs.tf`：输出资源 ID
- `README.md`：使用说明

### 4.4 `deploy/helm/` 部署基线
- `maas-middleware/`：依赖 PostgreSQL/Redis/Kafka
- `maas-observability/`：依赖 kube-prometheus-stack（可选 Loki/Promtail）

## 5. API 入口与现有接口
**当前已实现（HTTP/Hertz）**
- `POST /tenants`
- `GET /tenants/{id}`
- `GET /tenants`
- `DELETE /tenants/{id}`
- `GET /healthz`

**规划中的 RPC**
- `TenantService`（Kitex）将根据 `api/proto/tenant.proto` 生成代码并对内提供 RPC

## 6. 开发节点与步骤（全量规划）
> 总计 8 个节点，细化到每一步可独立提交与验收。每步均有对应文档：`docs/node-XX/step-YY-*.md`。

### Node 01：项目初始化与工程基线（已完成）
- Step 01：仓库结构 + Go Module + Makefile
  - 实现：统一目录结构、初始化 Go module、添加 README 与 .gitignore
  - 产物：`go.mod`, `Makefile`, `README.md`, 目录树
  - 文档：`docs/node-01/step-01-bootstrap.md`
- Step 02：GitHub Actions CI
  - 实现：push/PR 触发 tidy/test/vet
  - 产物：`.github/workflows/ci.yml`
  - 文档：`docs/node-01/step-02-ci.md`
- Step 03：配置与日志基线
  - 实现：统一配置加载、Hertz/Kitex 日志适配
  - 产物：`internal/config`, `internal/logging`
  - 文档：`docs/node-01/step-03-config-logging.md`

### Node 02：基础设施与部署基线（已完成）
- Step 01：Terraform 基线（ACK + GPU 节点池）
  - 实现：VPC/VSwitch/ACK/GPU NodePool
  - 产物：`infra/terraform/*.tf`
  - 文档：`docs/node-02/step-01-terraform-base.md`
- Step 02：Helm 中间件基线
  - 实现：PostgreSQL/Redis/Kafka umbrella chart
  - 产物：`deploy/helm/maas-middleware`
  - 文档：`docs/node-02/step-02-helm-middleware.md`
- Step 03：OSS + 可观测基线
  - 实现：OSS bucket + Prometheus/Grafana chart
  - 产物：`infra/terraform/oss.tf`, `deploy/helm/maas-observability`
  - 文档：`docs/node-02/step-03-observability-base.md`

### Node 03：租户与鉴权（进行中）
- Step 01：Tenant Service 骨架（已完成）
  - 实现：Tenant 领域模型、内存存储、Hertz API
  - 产物：`internal/tenant/*`, `cmd/tenant-service/main.go`, `api/proto/tenant.proto`
  - 文档：`docs/node-03/step-01-tenant-service.md`
- Step 02：鉴权中间件（规划）
  - 实现：API Key 校验、请求鉴权
  - 预期产物：`internal/auth/*`, Hertz middleware
  - 文档：`docs/node-03/step-02-auth.md`
- Step 03：配额模型（规划）
  - 实现：租户配额、限流策略
  - 预期产物：`internal/quota/*`
  - 文档：`docs/node-03/step-03-quota.md`

### Node 04：模型注册与版本管理（规划）
- Step 01：模型注册 API
- Step 02：版本与运行时配置
- Step 03：模型生命周期状态机

### Node 05：同步推理链路（规划）
- Step 01：Inference Gateway + 鉴权
- Step 02：vLLM/Triton 运行时接入
- Step 03：统一错误码与返回协议

### Node 06：异步任务与回调（规划）
- Step 01：任务模型 + Kafka 队列
- Step 02：GPU Worker + 结果存储
- Step 03：回调与轮询 API

### Node 07：计量计费与审计（规划）
- Step 01：计量事件上报
- Step 02：账单生成
- Step 03：审计日志与留存

### Node 08：稳定性与 SLO（规划）
- Step 01：限流/熔断/重试
- Step 02：SLO 告警与仪表盘
- Step 03：灰度发布与回滚流程

## 7. 当前实现与下一步建议
- 已具备最小可运行的 Tenant Service（HTTP API）与完整基础设施/部署基线。
- 下一步可开始 Node03/Step02（鉴权）或先解决 cwgo 代码生成问题再进入 Kitex。

## 8. 文档索引
- `docs/node-01/`：工程基线
- `docs/node-02/`：基础设施与部署
- `docs/node-03/`：租户与鉴权（进行中）

---
若希望本文档继续扩充（例如加入模块间调用时序、RPC 接口细节、数据库表结构等），可在对应节点完成后迭代补充。
