from api_gateway.internal.handler.health_handler import router as health_router
from api_gateway.internal.handler.user_handler import router as user_router
from api_gateway.internal.handler.model_handler import router as model_router

__all__ = ["health_router", "user_router", "model_router"]
