from datetime import datetime, timezone

from fastapi import APIRouter, Depends
from sqlalchemy import text
from sqlalchemy.ext.asyncio import AsyncSession

from api_gateway.internal.repository.database import get_session

router = APIRouter()


@router.get("/")
async def root() -> dict[str, str]:
    return {"message": "Welcome to MaaS Platform", "status": "running"}


@router.get("/health")
async def health_check(
    session: AsyncSession = Depends(get_session),
) -> dict[str, str | int]:
    db_status = "up"
    try:
        await session.execute(text("SELECT 1"))
    except Exception:
        db_status = "down"

    return {
        "status": "ok",
        "service": "api-gateway",
        "database": db_status,
        "timestamp": int(datetime.now(timezone.utc).timestamp()),
    }
