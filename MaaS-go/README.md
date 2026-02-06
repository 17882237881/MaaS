# MaaS Platform

Model-as-a-Service Platform - 模型即服务平台

## 项目简介

MaaS（模型即服务）平台提供统一的模型接入、推理服务、版本管理、计量计费等能力。

## 技术栈

- **语言**：Go 1.21+
- **框架**：Gin, gRPC
- **数据库**：PostgreSQL, Redis
- **消息队列**：Kafka
- **对象存储**：MinIO
- **容器化**：Docker, Kubernetes
- **监控**：Prometheus, Grafana

## 项目结构

```
├── api-gateway/          # API网关服务
│   ├── cmd/             # 程序入口
│   ├── internal/        # 内部代码
│   │   ├── config/     # 配置管理
│   │   ├── handler/    # HTTP处理器
│   │   ├── middleware/ # 中间件
│   │   ├── model/      # 数据模型
│   │   ├── repository/ # 数据访问层
│   │   ├── router/     # 路由
│   │   └── service/    # 业务逻辑
│   └── pkg/            # 公共库
│       ├── logger/     # 日志
│       └── utils/      # 工具函数
│
├── model-registry/       # 模型注册服务
│   └── ...             # 相同结构
│
├── deploy/              # 部署配置
│   ├── docker/         # Docker配置
│   └── k8s/            # Kubernetes配置
│
├── docs/                # 文档
│   ├── 01-overview/    # 项目总览
│   ├── 02-phase1/      # 阶段1文档
│   ├── 03-phase2/      # 阶段2文档
│   ├── 04-phase3/      # 阶段3文档
│   └── 05-phase4/      # 阶段4文档
│
└── shared/              # 共享代码
    ├── proto/          # Protocol Buffers
    └── errors/         # 公共错误

```

## 快速开始

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### 本地开发

1. 克隆仓库
   ```bash
   git clone https://github.com/17882237881/MaaS.git
   cd MaaS
   ```

2. 运行API Gateway
   ```bash
   cd api-gateway
   go run cmd/main.go
   ```

3. 运行Model Registry
   ```bash
   cd model-registry
   go run cmd/main.go
   ```

## 文档

详细的学习文档请参考 `docs/` 目录：

- [项目总览](./docs/01-overview/README.md)
- [阶段1：基础架构](./docs/02-phase1/README.md)
- [阶段2：核心功能](./docs/03-phase2/README.md)
- [阶段3：企业级特性](./docs/04-phase3/README.md)
- [阶段4：高级优化](./docs/05-phase4/README.md)

## 许可证

MIT License

## 学习路径

这是一个系统性的MaaS平台学习项目，按照以下阶段逐步完成：

1. **阶段1（4周）**：基础架构搭建
2. **阶段2（5周）**：核心功能开发
3. **阶段3（4周）**：企业级特性
4. **阶段4（3周）**：高级特性与优化

每个阶段都有详细的文档和代码示例，请参考 `docs/` 目录。
