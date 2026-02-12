from fastapi import APIRouter, FastAPI

from api_gateway.internal.handler import health_router, user_router


def register_routes(app: FastAPI) -> None:
    app.include_router(health_router)

    api_v1 = APIRouter(prefix="/api/v1")
    api_v1.include_router(user_router, prefix="/users", tags=["users"])
    app.include_router(api_v1)
