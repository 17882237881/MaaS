from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession

from api_gateway.internal.model.user import (
    UserCreate,
    UserListResponse,
    UserResponse,
    UserUpdate,
)
from api_gateway.internal.repository.database import get_session
from api_gateway.internal.repository.user_repository import UserRepository
from api_gateway.internal.service.user_service import (
    UserAlreadyExistsError,
    UserNotFoundError,
    UserService,
)

router = APIRouter()


def _get_service(session: AsyncSession) -> UserService:
    return UserService(UserRepository(session))


@router.post("", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(
    payload: UserCreate,
    session: AsyncSession = Depends(get_session),
) -> UserResponse:
    service = _get_service(session)
    try:
        return await service.create_user(payload)
    except UserAlreadyExistsError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc


@router.get("", response_model=UserListResponse)
async def list_users(
    session: AsyncSession = Depends(get_session),
) -> UserListResponse:
    service = _get_service(session)
    items = await service.list_users()
    return UserListResponse(total=len(items), items=items)


@router.get("/{user_id}", response_model=UserResponse)
async def get_user(
    user_id: str,
    session: AsyncSession = Depends(get_session),
) -> UserResponse:
    service = _get_service(session)
    try:
        return await service.get_user(user_id)
    except UserNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc


@router.put("/{user_id}", response_model=UserResponse)
async def update_user(
    user_id: str,
    payload: UserUpdate,
    session: AsyncSession = Depends(get_session),
) -> UserResponse:
    service = _get_service(session)
    try:
        return await service.update_user(user_id, payload)
    except UserNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
    except UserAlreadyExistsError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc


@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(
    user_id: str,
    session: AsyncSession = Depends(get_session),
) -> None:
    service = _get_service(session)
    try:
        await service.delete_user(user_id)
    except UserNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
