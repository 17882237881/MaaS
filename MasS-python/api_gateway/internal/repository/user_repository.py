from __future__ import annotations

from datetime import datetime, timezone

from api_gateway.internal.model.user import UserResponse, UserUpdate


class InMemoryUserRepository:
    def __init__(self) -> None:
        self._users: dict[str, UserResponse] = {}

    async def list_users(self) -> list[UserResponse]:
        return sorted(self._users.values(), key=lambda user: user.created_at)

    async def get_by_id(self, user_id: str) -> UserResponse | None:
        return self._users.get(user_id)

    async def get_by_email(self, email: str) -> UserResponse | None:
        for user in self._users.values():
            if user.email == email:
                return user
        return None

    async def create(self, user: UserResponse) -> None:
        self._users[user.id] = user

    async def update(self, user_id: str, update: UserUpdate) -> UserResponse | None:
        existing = self._users.get(user_id)
        if existing is None:
            return None

        data = existing.model_dump()
        if update.name is not None:
            data["name"] = update.name
        if update.email is not None:
            data["email"] = update.email
        data["updated_at"] = datetime.now(timezone.utc)

        updated = UserResponse(**data)
        self._users[user_id] = updated
        return updated

    async def delete(self, user_id: str) -> bool:
        return self._users.pop(user_id, None) is not None
