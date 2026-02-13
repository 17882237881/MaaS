from __future__ import annotations

from dataclasses import dataclass
from uuid import UUID, uuid4

from model_registry.internal.domain.model import Model
from model_registry.internal.repository.postgres_repo import (
    DuplicateModelError,
    ModelFilter,
    ModelNotFoundError,
    Pagination,
    SqlAlchemyModelRepository,
)


class ModelServiceError(Exception):
    """Base error for model service."""


class ModelNotFoundServiceError(ModelServiceError):
    """Model not found."""


class DuplicateModelServiceError(ModelServiceError):
    """Duplicate model."""


class InvalidInputError(ModelServiceError):
    """Invalid input."""


@dataclass
class CreateModelRequest:
    name: str
    description: str
    version: str
    framework: str
    tags: list[str]
    metadata: dict[str, str]
    owner_id: str
    tenant_id: str
    is_public: bool


@dataclass
class UpdateModelRequest:
    name: str | None = None
    description: str | None = None
    tags: list[str] | None = None
    metadata: dict[str, str] | None = None
    is_public: bool | None = None


@dataclass
class ListModelsFilter:
    name: str = ""
    framework: str = ""
    status: str = ""
    owner_id: str = ""
    tenant_id: str = ""
    tags: list[str] | None = None
    is_public: bool | None = None
    page: int = 1
    limit: int = 20


@dataclass
class ListModelsResponse:
    models: list[Model]
    total: int
    page: int
    limit: int


class ModelService:
    """Model registry business logic."""

    _valid_frameworks = {
        "pytorch",
        "tensorflow",
        "onnx",
        "sklearn",
        "xgboost",
        "custom",
    }

    def __init__(self, repo: SqlAlchemyModelRepository) -> None:
        self._repo = repo

    async def create_model(self, req: CreateModelRequest) -> Model:
        if req.framework and req.framework not in self._valid_frameworks:
            raise InvalidInputError(f"invalid framework: {req.framework}")

        owner_id = _safe_uuid(req.owner_id)
        tenant_id = _safe_uuid(req.tenant_id)

        model = Model(
            name=req.name,
            description=req.description,
            version=req.version,
            framework=req.framework,
            status="pending",
            owner_id=owner_id,
            tenant_id=tenant_id,
            is_public=req.is_public,
        )

        try:
            created = await self._repo.create(model)
        except DuplicateModelError as exc:
            raise DuplicateModelServiceError(str(exc)) from exc

        if req.tags:
            await self._repo.add_tags(created.id, req.tags)
        if req.metadata:
            await self._repo.set_metadata(created.id, req.metadata)

        return created

    async def get_model(self, model_id: str) -> Model:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        model = await self._repo.get_by_id(model_uuid)
        if model is None:
            raise ModelNotFoundServiceError("model not found")
        return model

    async def list_models(self, req: ListModelsFilter) -> ListModelsResponse:
        filter = ModelFilter(
            name=req.name,
            framework=req.framework,
            status=req.status,
            owner_id=_parse_uuid(req.owner_id),
            tenant_id=_parse_uuid(req.tenant_id),
            tags=req.tags or [],
            is_public=req.is_public,
        )
        page = req.page if req.page > 0 else 1
        limit = req.limit if 0 < req.limit <= 100 else 20
        pagination = Pagination(page=page, limit=limit)
        models, total = await self._repo.list(filter, pagination)
        return ListModelsResponse(models=models, total=total, page=pagination.page, limit=pagination.limit)

    async def update_model(self, model_id: str, req: UpdateModelRequest) -> Model:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        model = await self._repo.get_by_id(model_uuid)
        if model is None:
            raise ModelNotFoundServiceError("model not found")

        if req.name is not None:
            model.name = req.name
        if req.description is not None:
            model.description = req.description
        if req.is_public is not None:
            model.is_public = req.is_public

        updated = await self._repo.update(model)

        if req.tags is not None:
            existing = {tag.name for tag in model.tags}
            desired = set(req.tags)
            await self._repo.add_tags(model_uuid, list(desired - existing))
            await self._repo.remove_tags(model_uuid, list(existing - desired))

        if req.metadata is not None:
            await self._repo.set_metadata(model_uuid, req.metadata)

        return updated

    async def delete_model(self, model_id: str) -> None:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        try:
            await self._repo.delete(model_uuid)
        except ModelNotFoundError as exc:
            raise ModelNotFoundServiceError(str(exc)) from exc

    async def update_model_status(self, model_id: str, status: str) -> Model:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        try:
            await self._repo.update_status(model_uuid, status)
        except ModelNotFoundError as exc:
            raise ModelNotFoundServiceError(str(exc)) from exc

        model = await self._repo.get_by_id(model_uuid)
        if model is None:
            raise ModelNotFoundServiceError("model not found")

        return model

    async def add_model_tags(self, model_id: str, tags: list[str]) -> None:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        try:
            await self._repo.add_tags(model_uuid, tags)
        except ModelNotFoundError as exc:
            raise ModelNotFoundServiceError(str(exc)) from exc

    async def remove_model_tags(self, model_id: str, tags: list[str]) -> None:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        try:
            await self._repo.remove_tags(model_uuid, tags)
        except ModelNotFoundError as exc:
            raise ModelNotFoundServiceError(str(exc)) from exc

    async def set_model_metadata(self, model_id: str, metadata: dict[str, str]) -> None:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        try:
            await self._repo.set_metadata(model_uuid, metadata)
        except ModelNotFoundError as exc:
            raise ModelNotFoundServiceError(str(exc)) from exc

    async def get_model_metadata(self, model_id: str) -> dict[str, str]:
        model_uuid = _parse_uuid(model_id)
        if model_uuid is None:
            raise InvalidInputError("invalid model id")

        return await self._repo.get_metadata(model_uuid)


def _safe_uuid(value: str) -> UUID:
    if not value:
        return uuid4()
    try:
        return UUID(value)
    except ValueError:
        return uuid4()


def _parse_uuid(value: str) -> UUID | None:
    if not value:
        return None
    try:
        return UUID(value)
    except ValueError:
        return None
