import sys

from loguru import logger

from api_gateway.internal.config.settings import settings


def setup_logger() -> None:
    """Configures the Loguru logger based on settings."""
    logger.remove()  # Remove default handler

    # Console handler
    logger.add(
        sys.stderr,
        level=settings.log.level,
        format=(
            "{time:YYYY-MM-DD HH:mm:ss} | {level} | {message}"
            if settings.log.format == "text"
            else "{message}"
        ),
        serialize=(settings.log.format == "json"),
    )

    # File handler
    logger.add(
        "logs/api_gateway.log",
        rotation="100 MB",
        retention="30 days",
        level=settings.log.level,
        serialize=True,  # Always JSON for files
        enqueue=True,  # Async writing
        encoding="utf-8",
    )
