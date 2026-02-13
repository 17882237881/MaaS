from __future__ import annotations

from datetime import datetime
from uuid import UUID, uuid4

from sqlalchemy import BigInteger, Boolean, ForeignKey, String, Text, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy.sql import func

from api_gateway.internal.repository.database import Base


class Model(Base):
    __tablename__ = "models"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    version: Mapped[str] = mapped_column(String(50), nullable=False)
    description: Mapped[str | None] = mapped_column(Text)
    framework: Mapped[str] = mapped_column(
        String(50), nullable=False
    )  # pytorch, tensorflow, etc.
    status: Mapped[str] = mapped_column(String(50), default="pending")
    storage_path: Mapped[str | None] = mapped_column(String(512))
    size_bytes: Mapped[int] = mapped_column(BigInteger, default=0)
    checksum: Mapped[str | None] = mapped_column(String(64))
    docker_image: Mapped[str | None] = mapped_column(String(255))

    is_public: Mapped[bool] = mapped_column(Boolean, default=False)
    owner_id: Mapped[UUID] = mapped_column(nullable=False)
    tenant_id: Mapped[UUID] = mapped_column(nullable=False)

    created_at: Mapped[datetime] = mapped_column(server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        server_default=func.now(), onupdate=func.now()
    )

    # Relationships
    tags: Mapped[list["Tag"]] = relationship(
        secondary="model_tags", back_populates="models"
    )
    metadata_entries: Mapped[list["ModelMetadata"]] = relationship(
        back_populates="model", cascade="all, delete-orphan"
    )

    __table_args__ = (
        UniqueConstraint("name", "version", name="uq_model_name_version"),
    )


class Tag(Base):
    __tablename__ = "tags"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    name: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)

    models: Mapped[list["Model"]] = relationship(
        secondary="model_tags", back_populates="tags"
    )


class ModelTag(Base):
    __tablename__ = "model_tags"

    model_id: Mapped[UUID] = mapped_column(ForeignKey("models.id"), primary_key=True)
    tag_id: Mapped[UUID] = mapped_column(ForeignKey("tags.id"), primary_key=True)


class ModelMetadata(Base):
    __tablename__ = "model_metadata"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    model_id: Mapped[UUID] = mapped_column(ForeignKey("models.id"), nullable=False)
    key: Mapped[str] = mapped_column(String(100), nullable=False)
    value: Mapped[str | None] = mapped_column(String(500))
    created_at: Mapped[datetime] = mapped_column(server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        server_default=func.now(), onupdate=func.now()
    )

    model: Mapped["Model"] = relationship(back_populates="metadata_entries")
