from fastapi import APIRouter, HTTPException, status

from api_gateway.internal.model.user import (
    UserCreate,
    UserListResponse,
    UserResponse,
    UserUpdate,
)
from api_gateway.internal.service.user_service import (
    UserAlreadyExistsError,
    UserNotFoundError,
    UserService,
)

router = APIRouter()
_service = UserService()


@router.post("", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(payload: UserCreate) -> UserResponse:
    try:
        return await _service.create_user(payload)
    except UserAlreadyExistsError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc


@router.get("", response_model=UserListResponse)
async def list_users() -> UserListResponse:
    items = await _service.list_users()
    return UserListResponse(total=len(items), items=items)


@router.get("/{user_id}", response_model=UserResponse)
async def get_user(user_id: str) -> UserResponse:
    try:
        return await _service.get_user(user_id)
    except UserNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc


@router.put("/{user_id}", response_model=UserResponse)
async def update_user(user_id: str, payload: UserUpdate) -> UserResponse:
    try:
        return await _service.update_user(user_id, payload)
    except UserNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
    except UserAlreadyExistsError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc


@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(user_id: str) -> None:
    try:
        await _service.delete_user(user_id)
    except UserNotFoundError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
