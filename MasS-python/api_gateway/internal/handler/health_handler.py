from fastapi import APIRouter

router = APIRouter()


@router.get("/")
async def root() -> dict[str, str]:
    return {"message": "Welcome to MaaS Platform", "status": "running"}


@router.get("/health")
async def health_check() -> dict[str, str]:
    return {"status": "ok", "service": "api-gateway"}
