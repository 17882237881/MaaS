import logging
from datetime import datetime, timezone

import grpc
from google.protobuf import empty_pb2, timestamp_pb2

from shared.proto import model_pb2, model_pb2_grpc
from model_registry.internal.repository.postgres_repo import SqlAlchemyModelRepository
from model_registry.internal.service.model_service import (
    CreateModelRequest,
    DuplicateModelServiceError,
    InvalidInputError,
    ListModelsFilter,
    ModelNotFoundServiceError,
    ModelService,
    UpdateModelRequest,
)
from api_gateway.internal.repository.database import SESSION_FACTORY


class ModelServiceServicer(model_pb2_grpc.ModelServiceServicer):
    """gRPC servicer for Model Registry."""

    def __init__(self) -> None:
        self._logger = logging.getLogger(__name__)

    async def CreateModel(
        self,
        request: model_pb2.CreateModelRequest,
        context: grpc.aio.ServicerContext,
    ) -> model_pb2.CreateModelResponse:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            req = CreateModelRequest(
                name=request.name,
                description=request.description,
                version=request.version,
                framework=request.framework,
                tags=list(request.tags),
                metadata=dict(request.metadata),
                owner_id=request.owner_id,
                tenant_id=request.tenant_id,
                is_public=request.is_public,
            )
            try:
                model = await service.create_model(req)
            except DuplicateModelServiceError as exc:
                await context.abort(grpc.StatusCode.ALREADY_EXISTS, str(exc))
                return model_pb2.CreateModelResponse()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return model_pb2.CreateModelResponse()
            except Exception as exc:  # pragma: no cover - safety net
                self._logger.exception("CreateModel failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return model_pb2.CreateModelResponse()

            return model_pb2.CreateModelResponse(model=_model_to_proto(model))

    async def GetModel(
        self,
        request: model_pb2.GetModelRequest,
        context: grpc.aio.ServicerContext,
    ) -> model_pb2.GetModelResponse:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                model = await service.get_model(request.id)
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return model_pb2.GetModelResponse()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return model_pb2.GetModelResponse()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("GetModel failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return model_pb2.GetModelResponse()

            return model_pb2.GetModelResponse(model=_model_to_proto(model))

    async def ListModels(
        self,
        request: model_pb2.ListModelsRequest,
        context: grpc.aio.ServicerContext,
    ) -> model_pb2.ListModelsResponse:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            req = ListModelsFilter(
                name=request.name,
                framework=request.framework,
                status=request.status,
                owner_id=request.owner_id,
                tenant_id=request.tenant_id,
                tags=list(request.tags),
                is_public=True if request.is_public else None,
                page=request.page or 1,
                limit=request.limit or 20,
            )
            try:
                result = await service.list_models(req)
            except Exception as exc:  # pragma: no cover
                self._logger.exception("ListModels failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return model_pb2.ListModelsResponse()

            return model_pb2.ListModelsResponse(
                models=[_model_to_proto(model) for model in result.models],
                total=result.total,
                page=result.page,
                limit=result.limit,
            )

    async def UpdateModel(
        self,
        request: model_pb2.UpdateModelRequest,
        context: grpc.aio.ServicerContext,
    ) -> model_pb2.UpdateModelResponse:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            req = UpdateModelRequest(
                name=request.name or None,
                description=request.description or None,
                tags=list(request.tags) if request.tags else None,
                metadata=dict(request.metadata) if request.metadata else None,
            )
            try:
                model = await service.update_model(request.id, req)
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return model_pb2.UpdateModelResponse()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return model_pb2.UpdateModelResponse()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("UpdateModel failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return model_pb2.UpdateModelResponse()

            return model_pb2.UpdateModelResponse(model=_model_to_proto(model))

    async def DeleteModel(
        self,
        request: model_pb2.DeleteModelRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                await service.delete_model(request.id)
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return empty_pb2.Empty()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return empty_pb2.Empty()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("DeleteModel failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return empty_pb2.Empty()

            return empty_pb2.Empty()

    async def UpdateModelStatus(
        self,
        request: model_pb2.UpdateModelStatusRequest,
        context: grpc.aio.ServicerContext,
    ) -> model_pb2.UpdateModelStatusResponse:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                model = await service.update_model_status(request.id, request.status)
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return model_pb2.UpdateModelStatusResponse()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return model_pb2.UpdateModelStatusResponse()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("UpdateModelStatus failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return model_pb2.UpdateModelStatusResponse()

            return model_pb2.UpdateModelStatusResponse(model=_model_to_proto(model))

    async def AddModelTags(
        self,
        request: model_pb2.AddModelTagsRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                await service.add_model_tags(request.model_id, list(request.tags))
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return empty_pb2.Empty()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return empty_pb2.Empty()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("AddModelTags failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return empty_pb2.Empty()

            return empty_pb2.Empty()

    async def RemoveModelTags(
        self,
        request: model_pb2.RemoveModelTagsRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                await service.remove_model_tags(request.model_id, list(request.tags))
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return empty_pb2.Empty()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return empty_pb2.Empty()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("RemoveModelTags failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return empty_pb2.Empty()

            return empty_pb2.Empty()

    async def SetModelMetadata(
        self,
        request: model_pb2.SetModelMetadataRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                await service.set_model_metadata(request.model_id, dict(request.metadata))
            except ModelNotFoundServiceError as exc:
                await context.abort(grpc.StatusCode.NOT_FOUND, str(exc))
                return empty_pb2.Empty()
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return empty_pb2.Empty()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("SetModelMetadata failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return empty_pb2.Empty()

            return empty_pb2.Empty()

    async def GetModelMetadata(
        self,
        request: model_pb2.GetModelMetadataRequest,
        context: grpc.aio.ServicerContext,
    ) -> model_pb2.GetModelMetadataResponse:
        async with SESSION_FACTORY() as session:
            service = ModelService(SqlAlchemyModelRepository(session))
            try:
                metadata = await service.get_model_metadata(request.model_id)
            except InvalidInputError as exc:
                await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(exc))
                return model_pb2.GetModelMetadataResponse()
            except Exception as exc:  # pragma: no cover
                self._logger.exception("GetModelMetadata failed")
                await context.abort(grpc.StatusCode.INTERNAL, str(exc))
                return model_pb2.GetModelMetadataResponse()

            return model_pb2.GetModelMetadataResponse(metadata=metadata)


def _model_to_proto(model: object) -> model_pb2.Model:
    tags = []
    if hasattr(model, "tags"):
        tags = [tag.name for tag in model.tags]

    return model_pb2.Model(
        id=str(model.id),
        name=model.name,
        description=model.description or "",
        version=model.version,
        framework=model.framework,
        status=model.status,
        size=getattr(model, "size_bytes", 0),
        checksum=model.checksum or "",
        storage_path=model.storage_path or "",
        docker_image=model.docker_image or "",
        tags=tags,
        owner_id=str(model.owner_id) if model.owner_id else "",
        tenant_id=str(model.tenant_id) if model.tenant_id else "",
        is_public=model.is_public,
        created_at=_to_timestamp(model.created_at),
        updated_at=_to_timestamp(model.updated_at),
    )


def _to_timestamp(value: datetime) -> timestamp_pb2.Timestamp:
    ts = timestamp_pb2.Timestamp()
    if value.tzinfo is None:
        value = value.replace(tzinfo=timezone.utc)
    ts.FromDatetime(value)
    return ts
