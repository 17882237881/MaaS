from __future__ import annotations

from api_gateway.internal.model.user import UserUpdate
from api_gateway.internal.model.user_entity import User
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession


class UserRepository:
    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def list_users(self) -> list[User]:
        result = await self._session.execute(select(User).order_by(User.created_at))
        return list(result.scalars().all())

    async def get_by_id(self, user_id: str) -> User | None:
        return await self._session.get(User, user_id)

    async def get_by_email(self, email: str) -> User | None:
        result = await self._session.execute(select(User).where(User.email == email))
        return result.scalar_one_or_none()

    async def create(self, user: User) -> User:
        self._session.add(user)
        await self._session.commit()
        await self._session.refresh(user)
        return user

    async def update(self, user_id: str, update: UserUpdate) -> User | None:
        user = await self.get_by_id(user_id)
        if user is None:
            return None

        if update.name is not None:
            user.name = update.name
        if update.email is not None:
            user.email = update.email

        await self._session.commit()
        await self._session.refresh(user)
        return user

    async def delete(self, user_id: str) -> bool:
        user = await self.get_by_id(user_id)
        if user is None:
            return False

        await self._session.delete(user)
        await self._session.commit()
        return True
