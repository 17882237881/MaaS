from fastapi import APIRouter, FastAPI

from api_gateway.internal.handler import health_router, model_router, user_router
from api_gateway.internal.handler.config_handler import router as config_router


def register_routes(app: FastAPI) -> None:
    app.include_router(health_router)
    app.include_router(config_router)

    api_v1 = APIRouter(prefix="/api/v1")
    api_v1.include_router(user_router, prefix="/users", tags=["users"])
    api_v1.include_router(model_router, prefix="/models", tags=["models"])
    app.include_router(api_v1)
