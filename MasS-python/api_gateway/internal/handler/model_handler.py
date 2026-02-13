from __future__ import annotations

from datetime import datetime
from uuid import uuid4

import grpc
from fastapi import APIRouter, Depends, HTTPException, Request, status
from fastapi import Query

from api_gateway.internal.client.grpc_client import ModelServiceClient
from api_gateway.internal.model.model import (
    ModelCreate,
    ModelListResponse,
    ModelMetadataUpdate,
    ModelResponse,
    ModelStatusUpdate,
    ModelTagsUpdate,
    ModelUpdate,
)
from shared.proto import model_pb2

router = APIRouter()


def _get_client(request: Request) -> ModelServiceClient:
    client = getattr(request.app.state, "model_registry_client", None)
    if client is None:
        raise HTTPException(status_code=503, detail="gRPC client not initialized")
    return client


def _to_response(model: model_pb2.Model) -> ModelResponse:
    created_at = model.created_at.ToDatetime() if model.created_at else datetime.utcnow()
    updated_at = model.updated_at.ToDatetime() if model.updated_at else datetime.utcnow()
    return ModelResponse(
        id=model.id,
        name=model.name,
        description=model.description or None,
        version=model.version,
        framework=model.framework,
        status=model.status,
        size=model.size,
        checksum=model.checksum or None,
        storage_path=model.storage_path or None,
        docker_image=model.docker_image or None,
        tags=list(model.tags),
        metadata={},
        owner_id=model.owner_id,
        tenant_id=model.tenant_id,
        is_public=model.is_public,
        created_at=created_at,
        updated_at=updated_at,
    )


def _handle_grpc_error(exc: grpc.RpcError) -> HTTPException:
    code = exc.code()
    if code == grpc.StatusCode.NOT_FOUND:
        return HTTPException(status_code=404, detail=str(exc))
    if code == grpc.StatusCode.ALREADY_EXISTS:
        return HTTPException(status_code=409, detail=str(exc))
    if code == grpc.StatusCode.INVALID_ARGUMENT:
        return HTTPException(status_code=400, detail=str(exc))
    return HTTPException(status_code=500, detail="internal server error")


def _resolve_identity(request: Request) -> tuple[str, str]:
    owner_id = getattr(request.state, "user_id", "")
    tenant_id = getattr(request.state, "tenant_id", "")
    if not owner_id:
        owner_id = str(uuid4())
    if not tenant_id:
        tenant_id = str(uuid4())
    return owner_id, tenant_id


@router.post("", response_model=ModelResponse, status_code=status.HTTP_201_CREATED)
async def create_model(
    payload: ModelCreate,
    request: Request,
    client: ModelServiceClient = Depends(_get_client),
) -> ModelResponse:
    owner_id, tenant_id = _resolve_identity(request)
    grpc_req = model_pb2.CreateModelRequest(
        name=payload.name,
        description=payload.description or "",
        version=payload.version,
        framework=payload.framework,
        tags=payload.tags,
        metadata=payload.metadata,
        owner_id=owner_id,
        tenant_id=tenant_id,
        is_public=payload.is_public,
    )
    try:
        model = await client.create_model(grpc_req)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc
    return _to_response(model)


@router.get("", response_model=ModelListResponse)
async def list_models(
    page: int = Query(1, ge=1),
    limit: int = Query(20, ge=1, le=100),
    name: str | None = None,
    framework: str | None = None,
    status: str | None = None,
    tags: list[str] | None = Query(default=None),
    client: ModelServiceClient = Depends(_get_client),
) -> ModelListResponse:
    grpc_req = model_pb2.ListModelsRequest(
        page=page,
        limit=limit,
    )
    if name:
        grpc_req.name = name
    if framework:
        grpc_req.framework = framework
    if status:
        grpc_req.status = status
    if tags:
        grpc_req.tags.extend(tags)
    try:
        models, total, page, limit = await client.list_models(grpc_req)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc

    return ModelListResponse(
        models=[_to_response(model) for model in models],
        total=total,
        page=page,
        limit=limit,
    )


@router.get("/{model_id}", response_model=ModelResponse)
async def get_model(
    model_id: str,
    client: ModelServiceClient = Depends(_get_client),
) -> ModelResponse:
    try:
        model = await client.get_model(model_id)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc
    return _to_response(model)


@router.put("/{model_id}", response_model=ModelResponse)
async def update_model(
    model_id: str,
    payload: ModelUpdate,
    client: ModelServiceClient = Depends(_get_client),
) -> ModelResponse:
    grpc_req = model_pb2.UpdateModelRequest(id=model_id)
    if payload.name is not None:
        grpc_req.name = payload.name
    if payload.description is not None:
        grpc_req.description = payload.description
    if payload.tags is not None:
        grpc_req.tags.extend(payload.tags)
    if payload.metadata is not None:
        grpc_req.metadata.update(payload.metadata)
    if payload.is_public is not None:
        grpc_req.is_public = payload.is_public

    try:
        model = await client.update_model(grpc_req)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc
    return _to_response(model)


@router.delete("/{model_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_model(
    model_id: str,
    client: ModelServiceClient = Depends(_get_client),
) -> None:
    try:
        await client.delete_model(model_id)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc


@router.patch("/{model_id}/status", response_model=ModelResponse)
async def update_model_status(
    model_id: str,
    payload: ModelStatusUpdate,
    client: ModelServiceClient = Depends(_get_client),
) -> ModelResponse:
    try:
        model = await client.update_model_status(model_id, payload.status)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc
    return _to_response(model)


@router.post("/{model_id}/tags", status_code=status.HTTP_200_OK)
async def add_model_tags(
    model_id: str,
    payload: ModelTagsUpdate,
    client: ModelServiceClient = Depends(_get_client),
) -> None:
    try:
        await client.add_model_tags(model_id, payload.tags)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc


@router.delete("/{model_id}/tags", status_code=status.HTTP_200_OK)
async def remove_model_tags(
    model_id: str,
    payload: ModelTagsUpdate,
    client: ModelServiceClient = Depends(_get_client),
) -> None:
    try:
        await client.remove_model_tags(model_id, payload.tags)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc


@router.get("/{model_id}/metadata")
async def get_model_metadata(
    model_id: str,
    client: ModelServiceClient = Depends(_get_client),
) -> dict[str, dict[str, str]]:
    try:
        metadata = await client.get_model_metadata(model_id)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc
    return {"metadata": metadata}


@router.put("/{model_id}/metadata", status_code=status.HTTP_200_OK)
async def set_model_metadata(
    model_id: str,
    payload: ModelMetadataUpdate,
    client: ModelServiceClient = Depends(_get_client),
) -> None:
    try:
        await client.set_model_metadata(model_id, payload.metadata)
    except grpc.RpcError as exc:
        raise _handle_grpc_error(exc) from exc
