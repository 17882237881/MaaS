from __future__ import annotations

from dataclasses import dataclass, field
from typing import Protocol
from uuid import UUID

from sqlalchemy import delete, func, select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from model_registry.internal.domain.model import Model, ModelMetadata, Tag


class ModelRepositoryError(Exception):
    """Base error for repository operations."""


class ModelNotFoundError(ModelRepositoryError):
    """Raised when a model cannot be found."""


class DuplicateModelError(ModelRepositoryError):
    """Raised when a model already exists."""


@dataclass
class ModelFilter:
    name: str = ""
    framework: str = ""
    status: str = ""
    owner_id: UUID | None = None
    tenant_id: UUID | None = None
    tags: list[str] = field(default_factory=list)
    is_public: bool | None = None


@dataclass
class Pagination:
    page: int = 1
    limit: int = 20


class ModelRepository(Protocol):
    """Interface for Model repository."""

    async def create(self, model: Model) -> Model: ...

    async def get_by_id(self, id: UUID) -> Model | None: ...

    async def get_by_name_version(self, name: str, version: str) -> Model | None: ...

    async def list(
        self, filter: ModelFilter, pagination: Pagination
    ) -> tuple[list[Model], int]: ...

    async def update(self, model: Model) -> Model: ...

    async def delete(self, id: UUID) -> None: ...

    async def update_status(self, id: UUID, status: str) -> None: ...

    async def add_tags(self, model_id: UUID, tags: list[str]) -> None: ...

    async def remove_tags(self, model_id: UUID, tags: list[str]) -> None: ...

    async def set_metadata(self, model_id: UUID, metadata: dict[str, str]) -> None: ...

    async def get_metadata(self, model_id: UUID) -> dict[str, str]: ...


class SqlAlchemyModelRepository:
    """SQLAlchemy implementation of ModelRepository."""

    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def create(self, model: Model) -> Model:
        existing = await self.get_by_name_version(model.name, model.version)
        if existing:
            raise DuplicateModelError("model with this name and version already exists")

        self._session.add(model)
        try:
            await self._session.commit()
        except IntegrityError as exc:
            await self._session.rollback()
            raise DuplicateModelError("model with this name and version already exists") from exc
        await self._session.refresh(model)
        return model

    async def get_by_id(self, id: UUID) -> Model | None:
        stmt = (
            select(Model)
            .options(selectinload(Model.tags), selectinload(Model.metadata_entries))
            .where(Model.id == id)
        )
        result = await self._session.execute(stmt)
        return result.scalar_one_or_none()

    async def get_by_name_version(self, name: str, version: str) -> Model | None:
        stmt = (
            select(Model)
            .options(selectinload(Model.tags), selectinload(Model.metadata_entries))
            .where(Model.name == name, Model.version == version)
        )
        result = await self._session.execute(stmt)
        return result.scalar_one_or_none()

    async def list(
        self, filter: ModelFilter, pagination: Pagination
    ) -> tuple[list[Model], int]:
        stmt = select(Model).options(
            selectinload(Model.tags), selectinload(Model.metadata_entries)
        )

        if filter.name:
            stmt = stmt.where(Model.name.ilike(f"%{filter.name}%"))
        if filter.framework:
            stmt = stmt.where(Model.framework == filter.framework)
        if filter.status:
            stmt = stmt.where(Model.status == filter.status)
        if filter.owner_id:
            stmt = stmt.where(Model.owner_id == filter.owner_id)
        if filter.tenant_id:
            stmt = stmt.where(Model.tenant_id == filter.tenant_id)
        if filter.is_public is not None:
            stmt = stmt.where(Model.is_public == filter.is_public)
        if filter.tags:
            stmt = stmt.join(Model.tags).where(Tag.name.in_(filter.tags)).distinct()

        count_stmt = select(func.count(func.distinct(Model.id)))
        if filter.name:
            count_stmt = count_stmt.where(Model.name.ilike(f"%{filter.name}%"))
        if filter.framework:
            count_stmt = count_stmt.where(Model.framework == filter.framework)
        if filter.status:
            count_stmt = count_stmt.where(Model.status == filter.status)
        if filter.owner_id:
            count_stmt = count_stmt.where(Model.owner_id == filter.owner_id)
        if filter.tenant_id:
            count_stmt = count_stmt.where(Model.tenant_id == filter.tenant_id)
        if filter.is_public is not None:
            count_stmt = count_stmt.where(Model.is_public == filter.is_public)
        if filter.tags:
            count_stmt = count_stmt.join(Model.tags).where(Tag.name.in_(filter.tags))

        total = await self._session.scalar(count_stmt)
        total_count = int(total or 0)

        page = pagination.page if pagination.page > 0 else 1
        limit = pagination.limit if 0 < pagination.limit <= 100 else 20
        offset = (page - 1) * limit

        stmt = stmt.order_by(Model.created_at.desc()).offset(offset).limit(limit)

        result = await self._session.execute(stmt)
        return list(result.scalars().all()), total_count

    async def update(self, model: Model) -> Model:
        await self._session.commit()
        await self._session.refresh(model)
        return model

    async def delete(self, id: UUID) -> None:
        model = await self.get_by_id(id)
        if model is None:
            raise ModelNotFoundError("model not found")

        await self._session.delete(model)
        await self._session.commit()

    async def update_status(self, id: UUID, status: str) -> None:
        model = await self.get_by_id(id)
        if model is None:
            raise ModelNotFoundError("model not found")

        model.status = status
        await self._session.commit()

    async def add_tags(self, model_id: UUID, tags: list[str]) -> None:
        if not tags:
            return

        model = await self.get_by_id(model_id)
        if model is None:
            raise ModelNotFoundError("model not found")

        for name in tags:
            tag = await self._session.scalar(select(Tag).where(Tag.name == name))
            if tag is None:
                tag = Tag(name=name)
                self._session.add(tag)
                await self._session.flush()
            if tag not in model.tags:
                model.tags.append(tag)

        await self._session.commit()

    async def remove_tags(self, model_id: UUID, tags: list[str]) -> None:
        if not tags:
            return

        model = await self.get_by_id(model_id)
        if model is None:
            raise ModelNotFoundError("model not found")

        model.tags = [tag for tag in model.tags if tag.name not in tags]
        await self._session.commit()

    async def set_metadata(self, model_id: UUID, metadata: dict[str, str]) -> None:
        model = await self.get_by_id(model_id)
        if model is None:
            raise ModelNotFoundError("model not found")

        await self._session.execute(
            delete(ModelMetadata).where(ModelMetadata.model_id == model_id)
        )

        for key, value in metadata.items():
            self._session.add(
                ModelMetadata(model_id=model_id, key=key, value=value)
            )

        await self._session.commit()

    async def get_metadata(self, model_id: UUID) -> dict[str, str]:
        stmt = select(ModelMetadata).where(ModelMetadata.model_id == model_id)
        result = await self._session.execute(stmt)
        metadata = result.scalars().all()
        return {item.key: item.value or "" for item in metadata}
