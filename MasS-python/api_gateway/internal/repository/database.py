from __future__ import annotations

import os
from typing import AsyncGenerator

from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.orm import DeclarativeBase

from api_gateway.internal.config.settings import settings


class Base(DeclarativeBase):
    pass


# Use configuration for database URL
DATABASE_URL = settings.database.url

ENGINE: AsyncEngine = create_async_engine(
    DATABASE_URL,
    pool_pre_ping=True,
)

SESSION_FACTORY = async_sessionmaker(ENGINE, expire_on_commit=False)


async def init_db() -> None:
    # Ensure models are imported before metadata create_all.
    from api_gateway.internal.model import user_entity  # noqa: F401

    async with ENGINE.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)


async def get_session() -> AsyncGenerator[AsyncSession, None]:
    async with SESSION_FACTORY() as session:
        yield session
