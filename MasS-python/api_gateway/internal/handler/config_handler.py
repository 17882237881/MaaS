from fastapi import APIRouter, HTTPException

from api_gateway.internal.config.settings import settings

router = APIRouter()


@router.get("/config")
async def get_config():
    """获取当前配置（仅开发环境）"""
    if settings.environment != "development":
        raise HTTPException(status_code=403, detail="Config endpoint disabled in production")

    return {
        "environment": settings.environment,
        "server": {
            "host": settings.server.host,
            "port": settings.server.port,
            "reload": settings.server.reload,
        },
        "database": {
            "host": settings.database.host,
            "port": settings.database.port,
            "user": settings.database.user,
            "database": settings.database.database,
            # 不暴露密码
        },
        "redis": {
            "host": settings.redis.host,
            "port": settings.redis.port,
            "db": settings.redis.db,
            # 不暴露密码
        },
        "log": {
            "level": settings.log.level,
            "format": settings.log.format,
        },
    }
