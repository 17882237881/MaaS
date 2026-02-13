import time

from fastapi import Request
from prometheus_client import Counter, Gauge, Histogram
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint
from starlette.responses import Response

# Prometheus Metrics
REQUEST_COUNT = Counter(
    "http_requests_total",
    "Total HTTP requests",
    ["method", "path", "status"],
)
REQUEST_LATENCY = Histogram(
    "http_request_duration_seconds",
    "HTTP request latency",
    ["method", "path"],
)
IN_FLIGHT_REQUESTS = Gauge(
    "http_requests_in_flight",
    "Number of requests currently being processed",
    ["method"],
)


class MetricsMiddleware(BaseHTTPMiddleware):
    async def dispatch(
        self, request: Request, call_next: RequestResponseEndpoint
    ) -> Response:
        method = request.method
        path = request.url.path

        IN_FLIGHT_REQUESTS.labels(method=method).inc()
        start_time = time.time()

        try:
            response = await call_next(request)
            status_code = response.status_code
        except Exception:
            status_code = 500
            raise
        finally:
            duration = time.time() - start_time
            IN_FLIGHT_REQUESTS.labels(method=method).dec()
            REQUEST_COUNT.labels(method=method, path=path, status=status_code).inc()
            REQUEST_LATENCY.labels(method=method, path=path).observe(duration)

        return response
