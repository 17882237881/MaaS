from __future__ import annotations

from fastapi import HTTPException
from loguru import logger
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import JSONResponse, Response


class RecoveryMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next) -> Response:
        try:
            return await call_next(request)
        except HTTPException:
            raise
        except Exception:
            request_id = getattr(request.state, "request_id", "")
            logger.bind(
                request_id=request_id,
                method=request.method,
                path=request.url.path,
            ).exception("Unhandled exception")

            return JSONResponse(
                status_code=500,
                content={
                    "error": "Internal server error",
                    "code": "INTERNAL_ERROR",
                    "request_id": request_id,
                },
            )
