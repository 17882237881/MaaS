from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel, Field


_EMAIL_PATTERN = r"^[\w\.-]+@[\w\.-]+\.\w+$"


class UserBase(BaseModel):
    name: str = Field(..., min_length=1, max_length=100)
    email: str = Field(..., pattern=_EMAIL_PATTERN)


class UserCreate(UserBase):
    pass


class UserUpdate(BaseModel):
    name: str | None = Field(default=None, min_length=1, max_length=100)
    email: str | None = Field(default=None, pattern=_EMAIL_PATTERN)


class UserResponse(UserBase):
    id: str
    created_at: datetime
    updated_at: datetime


class UserListResponse(BaseModel):
    total: int
    items: list[UserResponse]
