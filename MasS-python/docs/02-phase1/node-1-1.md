# 阶段1 节点1：项目初始化与架构设计

> 🏗️ **节点目标**：使用 `uv` 初始化项目，搭建符合微服务规范的目录结构，并运行第一个FastAPI应用。

## 1. 环境准备

### 1.1 安装 uv
我们使用 `uv` 替代 pip/conda/poetry 进行包管理。

**Windows (PowerShell)**:
```powershell
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
```

**验证安装**:
```bash
uv --version
# 输出示例: uv 0.1.0 (... 2024-02-01)
```

### 1.2 uv vs Conda 常用命令对照

如果你习惯使用 Conda，可以参考下表快速上手 `uv`：

| 操作 | Conda 命令 | uv 命令 | 说明 |
|------|-----------|---------|------|
| **创建环境** | `conda create -n myenv python=3.11` | `uv venv` | uv 默认在项目目录下创建 `.venv` |
| **激活环境** | `conda activate myenv` | `.venv\Scripts\activate` | Windows下激活方式 |
| **安装包** | `conda install numpy` | `uv add numpy` | uv 会自动更新 pyproject.toml |
| **安装开发包** | (无直接对应与项目关联) | `uv add --dev pytest` | 标记为开发依赖 |
| **查看环境** | `conda list` | `uv pip list` | 或 `uv tree` 查看依赖树 |
| **导出依赖** | `conda env export > env.yml` | (自动维护 `uv.lock`) | lock 文件锁定确切版本 |
| **同步环境** | `conda env create -f env.yml` | `uv sync` | 一键同步环境到 lock 文件状态 |

---

## 2. 项目初始化

### 2.1 初始化 uv 项目

```bash
# 确保在项目根目录 (D:\code\MaaS\MasS-python)
cd D:\code\MaaS\MasS-python

# 初始化项目
uv init
```

这会创建 `pyproject.toml` 文件。

### 2.2 安装核心依赖

```bash
# 添加 Web 框架和服务器
uv add fastapi uvicorn[standard]

# 添加配置管理和日志
uv add pydantic-settings loguru

# 添加开发工具 (测试、代码检查)
uv add --dev pytest ruff mypy
```

执行后，你会发现多了一个 `uv.lock` 文件，它记录了所有依赖的确切版本，确保团队协作时环境一致。

---

## 3. 搭建目录结构

我们将创建符合企业级微服务规范的目录结构。

### 3.1 创建核心目录

在 `MasS-python` 目录下执行：

```powershell
# API Gateway 服务目录
mkdir api_gateway/cmd -Force
mkdir api_gateway/config -Force
mkdir api_gateway/internal/handler -Force
mkdir api_gateway/internal/middleware -Force
mkdir api_gateway/internal/model -Force
mkdir api_gateway/internal/repository -Force
mkdir api_gateway/internal/router -Force
mkdir api_gateway/internal/service -Force
mkdir api_gateway/pkg/logger -Force

# Model Registry 服务目录 (暂时建立基础)
mkdir model_registry -Force

# 共享代码目录
mkdir shared/proto -Force

# 测试目录
mkdir tests -Force
```

### 3.2 目录结构说明

```
MasS-python/
├── api_gateway/              # API网关服务
│   ├── main.py              # 服务入口
│   ├── config/              # 配置文件
│   ├── internal/            # 内部业务逻辑 (不对外暴露)
│   │   ├── handler/         # HTTP 接口层 (Controller)
│   │   ├── service/         # 业务逻辑层 (Service)
│   │   ├── repository/      # 数据访问层 (DAO)
│   │   └── model/           # 数据模型
│   └── pkg/                 # 公共组件 (可被其他项目引用)
├── pyproject.toml            # 项目依赖配置
├── uv.lock                   # 依赖锁定文件
└── README.md
```

---

## 4. 编写第一个服务

### 4.1 创建入口文件 `api_gateway/main.py`

```python
from fastapi import FastAPI
from loguru import logger

# 初始化 FastAPI 应用
app = FastAPI(
    title="MaaS Platform API Gateway",
    description="Model-as-a-Service 平台 API 网关",
    version="0.1.0"
)

# 启动事件
@app.on_event("startup")
async def startup_event():
    logger.info("API Gateway 服务启动中...")

# 基础路由
@app.get("/")
async def root():
    return {"message": "Welcome to MaaS Platform", "status": "running"}

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "api-gateway"}

if __name__ == "__main__":
    import uvicorn
    # 开发模式启动
    uvicorn.run("api_gateway.main:app", host="0.0.0.0", port=8000, reload=True)
```

### 4.2 创建运行脚本 `Makefile` (可选)

如果你安装了 make 工具，可以创建 `Makefile` 简化命令。如果Windows下没有make，可以跳过。

```makefile
dev:
	uv run uvicorn api_gateway.main:app --reload --port 8000
```

---

## 5. 运行与验证

### 5.1 启动服务

使用 `uv run` 可以在虚拟环境中执行命令，无需手动激活环境。

```bash
uv run uvicorn api_gateway.main:app --reload --port 8000
```

### 5.2 验证接口

打开浏览器或使用 curl 访问：

1.  **首页**: http://localhost:8000/
    - 预期响应: `{"message": "Welcome to MaaS Platform", "status": "running"}`
2.  **健康检查**: http://localhost:8000/health
    - 预期响应: `{"status": "ok", "service": "api-gateway"}`
3.  **API 文档**: http://localhost:8000/docs
    - FastAPI 自动生成的 Swagger UI

---

## ✅ 完成检查清单

- [ ] 已安装 `uv` 并验证版本
- [ ] 项目初始化完成 (`uv init`)
- [ ] 核心依赖已安装 (`fastapi`, `uvicorn`, `loguru`)
- [ ] 项目目录结构已创建
- [ ] `api_gateway/main.py` 代码已编写
- [ ] 服务能正常启动并访问 API 文档
