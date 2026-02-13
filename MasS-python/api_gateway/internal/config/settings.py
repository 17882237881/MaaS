from __future__ import annotations

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class DatabaseConfig(BaseSettings):
    """数据库配置"""

    host: str = "localhost"
    port: int = 5432
    user: str = "maas"
    password: str = "maas"
    database: str = "maas"

    @property
    def url(self) -> str:
        """构建数据库连接 URL"""
        return f"postgresql+asyncpg://{self.user}:{self.password}@{self.host}:{self.port}/{self.database}"


class RedisConfig(BaseSettings):
    """Redis 配置"""

    host: str = "localhost"
    port: int = 6379
    password: str = ""
    db: int = 0


class ServerConfig(BaseSettings):
    """服务器配置"""

    host: str = "0.0.0.0"
    port: int = 8000
    reload: bool = True


class LogConfig(BaseSettings):
    """日志配置"""

    level: str = "INFO"
    format: str = "json"


class GrpcConfig(BaseSettings):
    """gRPC 配置"""

    model_registry_host: str = "localhost"
    model_registry_port: int = 9090

    @property
    def model_registry_target(self) -> str:
        return f"{self.model_registry_host}:{self.model_registry_port}"


class Settings(BaseSettings):
    """应用程序主配置"""

    model_config = SettingsConfigDict(
        env_prefix="MAAS_",
        env_nested_delimiter="__",
        case_sensitive=False,
    )

    environment: str = "development"
    server: ServerConfig = Field(default_factory=ServerConfig)
    database: DatabaseConfig = Field(default_factory=DatabaseConfig)
    redis: RedisConfig = Field(default_factory=RedisConfig)
    grpc: GrpcConfig = Field(default_factory=GrpcConfig)
    log: LogConfig = Field(default_factory=LogConfig)


# 全局配置实例
settings = Settings()
