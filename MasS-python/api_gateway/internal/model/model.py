from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel, Field


class ModelCreate(BaseModel):
    name: str = Field(..., min_length=1, max_length=255)
    description: str | None = None
    version: str = Field(..., min_length=1, max_length=50)
    framework: str = Field(..., min_length=1, max_length=50)
    tags: list[str] = Field(default_factory=list)
    metadata: dict[str, str] = Field(default_factory=dict)
    is_public: bool = False


class ModelUpdate(BaseModel):
    name: str | None = Field(default=None, min_length=1, max_length=255)
    description: str | None = None
    tags: list[str] | None = None
    metadata: dict[str, str] | None = None
    is_public: bool | None = None


class ModelStatusUpdate(BaseModel):
    status: str = Field(..., min_length=1, max_length=50)


class ModelTagsUpdate(BaseModel):
    tags: list[str] = Field(...)


class ModelMetadataUpdate(BaseModel):
    metadata: dict[str, str] = Field(...)


class ModelResponse(BaseModel):
    id: str
    name: str
    description: str | None
    version: str
    framework: str
    status: str
    size: int
    checksum: str | None
    storage_path: str | None
    docker_image: str | None
    tags: list[str]
    metadata: dict[str, str]
    owner_id: str
    tenant_id: str
    is_public: bool
    created_at: datetime
    updated_at: datetime


class ModelListResponse(BaseModel):
    models: list[ModelResponse]
    total: int
    page: int
    limit: int
