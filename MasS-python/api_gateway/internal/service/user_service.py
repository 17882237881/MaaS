from __future__ import annotations

from datetime import datetime, timezone
from uuid import uuid4

from api_gateway.internal.model.user import UserCreate, UserResponse, UserUpdate
from api_gateway.internal.repository.user_repository import InMemoryUserRepository


class UserAlreadyExistsError(ValueError):
    pass


class UserNotFoundError(ValueError):
    pass


class UserService:
    def __init__(self, repository: InMemoryUserRepository | None = None) -> None:
        self._repository = repository or InMemoryUserRepository()

    async def list_users(self) -> list[UserResponse]:
        return await self._repository.list_users()

    async def get_user(self, user_id: str) -> UserResponse:
        user = await self._repository.get_by_id(user_id)
        if user is None:
            raise UserNotFoundError("User not found")
        return user

    async def create_user(self, payload: UserCreate) -> UserResponse:
        existing = await self._repository.get_by_email(payload.email)
        if existing is not None:
            raise UserAlreadyExistsError("Email already exists")

        now = datetime.now(timezone.utc)
        user = UserResponse(
            id=str(uuid4()),
            name=payload.name,
            email=payload.email,
            created_at=now,
            updated_at=now,
        )
        await self._repository.create(user)
        return user

    async def update_user(self, user_id: str, payload: UserUpdate) -> UserResponse:
        user = await self._repository.get_by_id(user_id)
        if user is None:
            raise UserNotFoundError("User not found")

        if payload.email:
            existing = await self._repository.get_by_email(payload.email)
            if existing is not None and existing.id != user_id:
                raise UserAlreadyExistsError("Email already exists")

        updated = await self._repository.update(user_id, payload)
        if updated is None:
            raise UserNotFoundError("User not found")
        return updated

    async def delete_user(self, user_id: str) -> None:
        deleted = await self._repository.delete(user_id)
        if not deleted:
            raise UserNotFoundError("User not found")
