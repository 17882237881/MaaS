# 节点1.3：配置管理体系（Configuration Management System）

## 学习目标
- 掌握 Pydantic-Settings 进行配置管理
- 学会 YAML 配置文件组织
- 理解环境变量覆盖机制
- 实现多环境配置切换（dev/prod）
- 配置验证和类型安全

## 为什么需要配置管理？

在实际项目中，不同环境（开发/测试/生产）需要不同的配置：

- **数据库连接**：开发环境用本地 PostgreSQL，生产环境用 AWS RDS
- **日志级别**：开发环境 DEBUG，生产环境 WARNING
- **服务端口**：开发环境 8000，生产环境 80/443
- **API 密钥**：不同环境使用不同密钥

**核心问题**：如何优雅地管理这些配置？

## Python vs Go 配置管理对比

| 特性 | Go (Viper) | Python (Pydantic-Settings) |
|------|------------|----------------------------|
| 配置文件格式 | YAML/JSON/TOML | YAML/JSON/.env |
| 环境变量支持 | `SetEnvPrefix("MAAS")` | `env_prefix="MAAS_"` |
| 嵌套配置 | `Get("database.host")` | `settings.database.host` |
| 类型安全 | 需手动类型转换 | 自动类型验证（Pydantic） |
| 配置验证 | 需手动编写 | 内置 validator |
| 热更新 | `WatchConfig()` | 需额外集成 Watchdog |

**Python 优势**：Pydantic 提供了强大的类型检查和自动验证，在编译时就能发现配置错误。

## 配置优先级（从高到低）

1.  **命令行参数** (暂不实现)
2.  **环境变量** (`.env` 文件或系统环境变量)
3.  **配置文件** (`config.yaml` 或 `config.{environment}.yaml`)
4.  **代码默认值** (Settings 类中的 `Field(default=...)`)

这和 Go 版的 Viper 优先级一致。

## 核心技术：Pydantic-Settings

### 1. 基础配置类

```python
from pydantic_settings import BaseSettings
from pydantic import Field

class DatabaseConfig(BaseSettings):
    host: str = "localhost"
    port: int = 5432
    user: str = "maas"
    password: str = "maas"    # 生产环境用环境变量覆盖
    database: str = "maas"
    
    @property
    def url(self) -> str:
        return f"postgresql+asyncpg://{self.user}:{self.password}@{self.host}:{self.port}/{self.database}"

class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_prefix="MAAS_",              # 环境变量前缀
        env_nested_delimiter="__",        # 嵌套配置分隔符
        case_sensitive=False,             # 不区分大小写
    )
    
    environment: str = "development"
    database: DatabaseConfig = Field(default_factory=DatabaseConfig)
```

### 2. 环境变量覆盖

假设你有如下配置类：

```python
class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="MAAS_")
    
    server_port: int = 8000
    database_host: str = "localhost"
```

可以通过环境变量覆盖：

```bash
export MAAS_SERVER_PORT=9000
export MAAS_DATABASE_HOST=prod-db.example.com
```

Pydantic-Settings 会自动读取这些环境变量（支持 `.env` 文件）。

**嵌套配置**：

```python
class Settings(BaseSettings):
    database: DatabaseConfig
```

通过环境变量覆盖：

```bash
export MAAS_DATABASE__HOST=prod-db.example.com  # 注意双下划线
export MAAS_DATABASE__PORT=5433
```

### 3. 多环境配置文件

**目录结构**：

```
configs/
├── config.yaml                # 基础默认配置
├── config.development.yaml    # 开发环境覆盖
└── config.production.yaml     # 生产环境覆盖
```

**config.yaml**（默认值）：

```yaml
environment: development

server:
  host: 0.0.0.0
  port: 8000
  reload: true

database:
  host: localhost
  port: 5432
  user: maas
  password: maas
  database: maas

log:
  level: INFO
  format: json
```

**config.production.yaml**（生产环境）：

```yaml
server:
  reload: false  # 生产环境不启用热重载

log:
  level: WARNING  # 生产环境减少日志
```

**加载逻辑**（基于环境变量 `MAAS_ENVIRONMENT`）：

```python
import os
import yaml
from pathlib import Path

def load_settings() -> Settings:
    env = os.getenv("MAAS_ENVIRONMENT", "development")
    
    # 1. 加载基础配置
    base_config = yaml.safe_load(Path("configs/config.yaml").read_text())
    
    # 2. 加载环境特定配置（如果存在）
    env_config_path = Path(f"configs/config.{env}.yaml")
    if env_config_path.exists():
        env_config = yaml.safe_load(env_config_path.read_text())
        base_config.update(env_config)  # 合并配置
    
    # 3. 创建 Pydantic Settings（会自动读取环境变量）
    return Settings(**base_config)
```

### 4. 配置验证

Pydantic 提供强大的验证能力：

```python
from pydantic import Field, field_validator

class DatabaseConfig(BaseSettings):
    port: int = Field(ge=1, le=65535)  # 端口范围验证
    
    @field_validator('user')
    @classmethod
    def validate_user(cls, v: str) -> str:
        if len(v) < 3:
            raise ValueError('Username must be at least 3 characters')
        return v
```

如果传入非法值，Pydantic 会在**启动时**直接报错，避免运行时崩溃。

## 实操步骤

### 步骤1：创建配置类

创建 `api_gateway/internal/config/settings.py`：

```python
from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import Field

class DatabaseConfig(BaseSettings):
    host: str = "localhost"
    port: int = 5432
    user: str = "maas"
    password: str = "maas"
    database: str = "maas"
    
    @property
    def url(self) -> str:
        return f"postgresql+asyncpg://{self.user}:{self.password}@{self.host}:{self.port}/{self.database}"

class RedisConfig(BaseSettings):
    host: str = "localhost"
    port: int = 6379
    password: str = ""
    db: int = 0

class ServerConfig(BaseSettings):
    host: str = "0.0.0.0"
    port: int = 8000
    reload: bool = True

class LogConfig(BaseSettings):
    level: str = "INFO"
    format: str = "json"

class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_prefix="MAAS_",
        env_nested_delimiter="__",
        case_sensitive=False,
    )
    
    environment: str = "development"
    server: ServerConfig = Field(default_factory=ServerConfig)
    database: DatabaseConfig = Field(default_factory=DatabaseConfig)
    redis: RedisConfig = Field(default_factory=RedisConfig)
    log: LogConfig = Field(default_factory=LogConfig)
```

### 步骤2：创建配置文件

创建 `configs/config.yaml`：

```yaml
environment: development

server:
  host: 0.0.0.0
  port: 8000
  reload: true

database:
  host: localhost
  port: 5432
  user: maas
  password: maas
  database: maas

redis:
  host: localhost
  port: 6379
  db: 0

log:
  level: INFO
  format: json
```

### 步骤3：创建环境变量模板

创建 `.env.example`：

```bash
# Server Configuration
MAAS_SERVER__HOST=0.0.0.0
MAAS_SERVER__PORT=8000

# Database Configuration
MAAS_DATABASE__HOST=localhost
MAAS_DATABASE__PORT=5432
MAAS_DATABASE__USER=maas
MAAS_DATABASE__PASSWORD=maas
MAAS_DATABASE__DATABASE=maas

# Redis Configuration
MAAS_REDIS__HOST=localhost
MAAS_REDIS__PORT=6379
```

### 步骤4：集成到应用

修改 `api_gateway/internal/repository/database.py`：

```python
from api_gateway.internal.config.settings import Settings

settings = Settings()  # 自动读取环境变量

DATABASE_URL = settings.database.url

ENGINE: AsyncEngine = create_async_engine(
    DATABASE_URL,
    pool_pre_ping=True,
)
```

修改 `main.py`：

```python
from api_gateway.internal.config.settings import Settings

settings = Settings()

if __name__ == "__main__":
    uvicorn.run(
        "api_gateway.main:app",
        host=settings.server.host,
        port=settings.server.port,
        reload=settings.server.reload,
    )
```

### 步骤5：验证配置加载

添加配置查看接口（仅开发环境）：

```python
@app.get("/config")
async def get_config():
    if settings.environment != "development":
        raise HTTPException(403, "Config endpoint disabled in production")
    
    return {
        "environment": settings.environment,
        "server": settings.server.model_dump(),
        "database": {
            "host": settings.database.host,
            "port": settings.database.port,
            # 不暴露密码
        }
    }
```

## 最佳实践

### 1. 敏感信息管理

**永远不要**将密码写入配置文件！使用环境变量：

```bash
export MAAS_DATABASE__PASSWORD=super_secret_password
```

或者使用 `.env` 文件（记得加入 `.gitignore`）。

### 2. 配置分层

```
默认配置 (config.yaml)
    ↓ 覆盖
环境配置 (config.development.yaml)
    ↓ 覆盖
环境变量 (.env 或系统环境变量)
```

### 3. 配置验证

在 `startup` 事件中验证关键配置：

```python
@app.on_event("startup")
async def validate_config():
    if settings.environment == "production":
        if settings.database.password == "maas":
            raise ValueError("生产环境不能使用默认密码！")
    
    logger.info(f"Configuration loaded for environment: {settings.environment}")
```

## 检查点

- [ ] 配置能从 YAML 文件加载
- [ ] 环境变量能正确覆盖 YAML 配置
- [ ] 多环境配置（dev/prod）能正确切换
- [ ] 配置验证能捕获非法值
- [ ] 敏感信息（如密码）不在代码仓库中

## 下一步

完成配置管理后，你将掌握：
- ✅ Pydantic-Settings 配置管理
- ✅ 多环境配置组织
- ✅ 环境变量最佳实践

**下一节**：[节点1.4：日志与监控基础](./node-1-4.md)
