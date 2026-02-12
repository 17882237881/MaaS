from __future__ import annotations

from api_gateway.internal.model.user import UserCreate, UserResponse, UserUpdate
from api_gateway.internal.model.user_entity import User
from api_gateway.internal.repository.user_repository import UserRepository


class UserAlreadyExistsError(ValueError):
    pass


class UserNotFoundError(ValueError):
    pass


class UserService:
    def __init__(self, repository: UserRepository) -> None:
        self._repository = repository

    async def list_users(self) -> list[UserResponse]:
        users = await self._repository.list_users()
        return [self._to_response(user) for user in users]

    async def get_user(self, user_id: str) -> UserResponse:
        user = await self._repository.get_by_id(user_id)
        if user is None:
            raise UserNotFoundError("User not found")
        return self._to_response(user)

    async def create_user(self, payload: UserCreate) -> UserResponse:
        existing = await self._repository.get_by_email(payload.email)
        if existing is not None:
            raise UserAlreadyExistsError("Email already exists")

        user = User(name=payload.name, email=payload.email)
        created = await self._repository.create(user)
        return self._to_response(created)

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
        return self._to_response(updated)

    async def delete_user(self, user_id: str) -> None:
        deleted = await self._repository.delete(user_id)
        if not deleted:
            raise UserNotFoundError("User not found")

    @staticmethod
    def _to_response(user: User) -> UserResponse:
        return UserResponse(
            id=user.id,
            name=user.name,
            email=user.email,
            created_at=user.created_at,
            updated_at=user.updated_at,
        )
