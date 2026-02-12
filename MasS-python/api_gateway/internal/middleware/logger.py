from __future__ import annotations

import time

from loguru import logger
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import Response


class LoggerMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next) -> Response:
        start = time.perf_counter()
        response = await call_next(request)
        duration_ms = (time.perf_counter() - start) * 1000

        request_id = getattr(request.state, "request_id", "")
        path = request.url.path
        if request.url.query:
            path = f"{path}?{request.url.query}"

        log = logger.bind(
            request_id=request_id,
            method=request.method,
            path=path,
            status=response.status_code,
            duration_ms=round(duration_ms, 2),
            client_ip=request.client.host if request.client else "",
        )

        if response.status_code >= 400:
            log.warning("Request completed with error")
        else:
            log.info("Request completed")

        return response
