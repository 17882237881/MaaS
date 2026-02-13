# 节点1.5：Docker化与部署（Docker & Deployment）

## 学习目标
- 掌握 Python 项目的 Dockerfile 编写最佳实践
- 理解多阶段构建（Multi-stage Build）以减小镜像体积
- 使用 Docker Compose 编排微服务环境（API + DB + Cache）
- 解决容器间网络通信和服务依赖问题

## 1. Dockerfile 最佳实践

### 为什么需要 Docker？
- **环境一致性**："在我的机器上能跑" -> "在任何地方都能跑"。
- **依赖隔离**：不再担心系统 Python 版本冲突。
- **快速部署**：一键启动完整环境。

### 多阶段构建 (Multi-stage Build)
我们使用 `uv` 进行依赖安装。为了保持生产镜像精简，我们将构建过程（下载、编译）和运行过程分离。

**Builder Stage**:
1.  下载并安装构建工具（uv, git, gcc 等）。
2.  复制依赖文件 (`pyproject.toml`, `uv.lock`)。
3.  编译并安装依赖到虚拟环境。

**Runtime Stage**:
1.  从 Builder 阶段复制编译好的虚拟环境。
2.  复制源代码。
3.  设置环境变量。
4.  启动应用。

这样，最终镜像不包含编译器和 uv 工具本身，体积更小，安全性更高。

## 2. Docker Compose 编排

在微服务架构中，我们的 API Gateway 不是孤岛，它依赖：
- **PostgreSQL**：持久化存储用户和模型数据。
- **Redis**：缓存热点数据和限流。

使用 `docker-compose.yml` 可以一次性定义和启动这些服务。

### 关键配置解析

```yaml
services:
  api-gateway:
    build: .
    environment:
      # 告诉服务：数据库在名为 "postgres" 的主机上，而不是 localhost
      - MAAS_DATABASE__HOST=postgres
    depends_on:
      postgres:
        condition: service_healthy  # 等数据库完全启动后再启动 API
```

### 网络服务发现
Docker Compose 会自动创建一个内部网络。服务之间可以通过**服务名**（service name）互相访问。
- API Gateway 访问数据库：`postgres:5432`
- API Gateway 访问缓存：`redis:6379`

**注意**：在容器内部，`localhost` 指的是容器自己，而不是宿主机。所以必须使用服务名。

## 3. 实操步骤

### 步骤1：创建 .dockerignore
防止不必要的文件（如本地虚拟环境、日志、git记录）被复制到镜像中。

### 步骤2：编写 Dockerfile
使用官方 Python Slim 镜像作为基础，结合 uv 进行依赖管理。

### 步骤3：编写 docker-compose.yml
定义三个服务：
1.  **postgres**: 使用官方 `postgres:15-alpine` 镜像。配置健康检查（pg_isready）。
2.  **redis**: 使用官方 `redis:7-alpine` 镜像。
3.  **api-gateway**: 构建当前目录的代码。配置环境变量覆盖默认的 `localhost` 地址。

### 步骤4：启动环境
```bash
docker-compose up -d
```

### 步骤5：验证
- 查看日志：`docker-compose logs -f`
- 检查状态：`docker-compose ps`
- 访问接口：`http://localhost:8000/health`

## 常见问题

### Q: 数据库连接失败？
**A**:
1.  检查 `depends_on` 是否配置了 `condition: service_healthy`。如果数据库还没准备好就启动 API，连接会报错。
2.  检查 `MAAS_DATABASE__HOST` 是否设置为了 `postgres`（compose 中的服务名）。

### Q: 镜像构建很慢？
**A**:
1.  确保 `.dockerignore` 过滤了 `.venv` 和 `logs`。
2.  Docker 利用层缓存（Layer Caching）。只要 `uv.lock` 没变，`RUN uv sync` 这一层就会被缓存，瞬间完成。

## 检查点
- [ ] 镜像构建成功且体积合理 (<200MB 理想，<500MB 可接受)
- [ ] `docker-compose up` 一键启动无报错
- [ ] API Gateway 能成功读写数据库
- [ ] `/metrics` 接口能访问
- [ ] 数据持久化（重启容器数据不丢失）

## 下一步
恭喜！你已经完成了**阶段1：基础框架搭建**的所有内容。
现在你拥有了一个：
- 结构清晰的 Python (FastAPI + uv) 项目
- 集成了 SQLAlchemy 2.0 数据库层
- 具备结构化日志和 Prometheus 监控
- 可以一键 Docker 部署的完整环境

接下来，我们将进入 **阶段2：核心功能开发**，开始构建真正的 Model Registry 业务逻辑。
