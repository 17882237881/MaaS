from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from loguru import logger
import uvicorn

from api_gateway.internal.middleware import (
    LoggerMiddleware,
    RecoveryMiddleware,
    RequestIDMiddleware,
)
from api_gateway.internal.router.router import register_routes

# 初始化 FastAPI 应用
app = FastAPI(
    title="MaaS Platform API Gateway",
    description="Model-as-a-Service 平台 API 网关",
    version="1.0.0",
)

# 中间件（顺序：Recovery -> Logger -> RequestID -> CORS）
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
    logger.info("API Gateway 服务启动中...")


if __name__ == "__main__":
    # 开发模式启动
    uvicorn.run("api_gateway.main:app", host="0.0.0.0", port=8000, reload=True)
