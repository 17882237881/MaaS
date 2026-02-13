from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from loguru import logger
from prometheus_client import make_asgi_app
import uvicorn

from api_gateway.internal.config.settings import settings
from api_gateway.internal.middleware import (
    LoggerMiddleware,
    RecoveryMiddleware,
    RequestIDMiddleware,
)
from api_gateway.internal.middleware.metrics import MetricsMiddleware
from api_gateway.internal.repository.database import init_db
from api_gateway.internal.router.router import register_routes
from api_gateway.pkg.logger import setup_logger
from api_gateway.internal.client.grpc_client import ModelServiceClient

# 初始化 FastAPI 应用
app = FastAPI(
    title="MaaS Platform API Gateway",
    description="Model-as-a-Service 平台 API 网关",
    version="1.0.0",
)

# Prometheus Metrics Endpoint
metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)

# 中间件（顺序：metrics -> Recovery -> Logger -> RequestID -> CORS）
app.add_middleware(MetricsMiddleware)
app.add_middleware(RecoveryMiddleware)
app.add_middleware(LoggerMiddleware)
app.add_middleware(RequestIDMiddleware)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
    allow_credentials=True,
)

# 路由注册
register_routes(app)


@app.on_event("startup")
async def startup_event() -> None:
    setup_logger()  # 配置日志
    await init_db()
    app.state.model_registry_client = ModelServiceClient(
        settings.grpc.model_registry_target
    )
    await app.state.model_registry_client.connect()
    logger.info(f"API Gateway starting in {settings.environment} mode...")
    logger.info(f"Server: {settings.server.host}:{settings.server.port}")


@app.on_event("shutdown")
async def shutdown_event() -> None:
    client = getattr(app.state, "model_registry_client", None)
    if client is not None:
        await client.close()


if __name__ == "__main__":
    # 开发模式启动 - 使用配置
    uvicorn.run(
        "api_gateway.main:app",
        host=settings.server.host,
        port=settings.server.port,
        reload=settings.server.reload,
    )
