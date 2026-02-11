# MaaS Platform Skill（Python版）

## 概述

这是一个完整的**MaaS（模型即服务）平台学习项目（Python版）**，通过从零构建企业级后端系统，系统学习微服务架构、云原生技术和分布式系统。

本项目**严格对齐Go版（MaaS-go）的架构分层和阶段节点**，使用Python生态的最佳实践重新实现。

## 适用人群

- 有一定Python基础（能写简单程序）
- 想系统学习后端技术栈
- 准备入职大厂后端岗位
- 想了解MaaS/AI平台开发
- 已学习过Go版，想对比学习Python实现

## 核心特点

### 📚 完整的文档体系
- 项目总览文档
- 4个阶段详细规划（16周）
- 每个节点都有详细的技术讲解和实操指南
- 与Go版一一对应，方便对比学习

### 🎯 渐进式学习路径
- **阶段1**：基础架构（4周）- FastAPI、SQLAlchemy、Docker
- **阶段2**：核心功能（5周）- gRPC、Redis、认证、推理
- **阶段3**：企业级特性（4周）- Kafka、K8s、限流、监控
- **阶段4**：高级优化（3周）- GPU调度、多租户、CI/CD

### 🛠️ 企业级技术栈
- **语言**：Python 3.11+
- **框架**：FastAPI, gRPC (grpcio)
- **数据库**：PostgreSQL (SQLAlchemy 2.0 + asyncpg), Redis
- **消息队列**：Kafka (confluent-kafka)
- **对象存储**：MinIO
- **容器化**：Docker, Kubernetes
- **监控**：Prometheus, Grafana, Jaeger (OpenTelemetry)
- **CI/CD**：GitHub Actions

## 使用方式

### 学习方式

1. **阅读文档**：从 `docs/01-overview/README.md` 开始
2. **跟随实操**：每个节点都有详细的代码实操任务
3. **提交代码**：每完成一个节点提交到GitHub
4. **检查清单**：每个节点都有完成检查清单

### 开发流程

```bash
# 1. 克隆仓库
git clone https://github.com/17882237881/MaaS.git
cd MaaS/MasS-python

# 2. 阅读当前阶段文档
code docs/02-phase1/node-1-1.md

# 3. 按照文档实操
# ...编写代码...

# 4. 提交代码
git add .
git commit -m "feat(phase1/node1.1): 完成项目初始化"
git push origin main

# 5. 进入下一节点
```

## 文档结构

```
docs/
├── 01-overview/README.md      # 项目总览
├── 02-phase1/
│   ├── README.md              # 阶段1总览
│   ├── node-1-1.md            # 节点1.1：项目初始化
│   ├── node-1-2.md            # 节点1.2：API Gateway
│   ├── node-1-3.md            # 节点1.3：配置管理
│   ├── node-1-4.md            # 节点1.4：日志监控
│   └── node-1-5.md            # 节点1.5：数据库设计
├── 03-phase2/                 # 阶段2文档
├── 04-phase3/                 # 阶段3文档
└── 05-phase4/                 # 阶段4文档
```

## 学习成果

完成本项目后，你将掌握：

✅ **微服务架构设计**
- 服务拆分原则
- 分层架构设计
- API Gateway模式

✅ **Python企业级开发**
- 项目结构组织（Poetry + pyproject.toml）
- 异步编程（asyncio + async/await）
- 类型安全（Pydantic + mypy）

✅ **分布式系统**
- gRPC服务通信（grpcio）
- 消息队列异步处理（Kafka）
- 分布式事务
- 一致性保障

✅ **云原生技术**
- Docker容器化
- Kubernetes编排
- Helm包管理

✅ **高可用架构**
- 限流熔断（slowapi + pybreaker）
- 降级策略
- 负载均衡

✅ **可观测性**
- 结构化日志（Loguru）
- 指标监控（Prometheus/Grafana）
- 链路追踪（OpenTelemetry/Jaeger）

✅ **DevOps实践**
- CI/CD流水线（GitHub Actions）
- 自动化测试（pytest）
- 代码质量（ruff + mypy）

## 项目里程碑

### 阶段1里程碑（4周后）
- API Gateway独立运行
- Model Registry支持CRUD
- Docker Compose一键启动
- **可演示**：基础API调用

### 阶段2里程碑（9周后）
- 完整的上传→注册→推理流程
- JWT认证和RBAC权限
- Redis缓存支持
- **可演示**：端到端模型推理

### 阶段3里程碑（13周后）
- Kafka异步任务
- K8s生产部署
- 限流熔断生效
- 监控仪表盘
- **可演示**：生产级部署

### 阶段4里程碑（16周后）
- GPU推理调度
- 多租户隔离
- 性能优化完成
- 自动化流水线
- **可演示**：企业级MaaS平台

## Go版 vs Python版 对比

| 方面 | Go版 (MaaS-go) | Python版 (MasS-python) |
|------|----------------|------------------------|
| HTTP框架 | Gin | FastAPI |
| ORM | GORM | SQLAlchemy 2.0 |
| 配置管理 | Viper | Pydantic-Settings |
| 日志 | Zap | Loguru |
| 并发模型 | goroutine | asyncio |
| 包管理 | go.mod | Poetry |
| 类型系统 | 静态类型 | 类型注解 + mypy |

## 技能标签

`Python` `微服务` `后端开发` `Kubernetes` `Docker` `gRPC` `Redis` `PostgreSQL` `Kafka` `Prometheus` `分布式系统` `云原生` `DevOps` `FastAPI` `SQLAlchemy` `asyncio`

## 开始你的学习之旅

👉 **立即开始**：[项目总览文档](./01-overview/README.md)

## 许可证

MIT License

---

**祝你学习愉快，早日成为后端大牛！** 🚀
