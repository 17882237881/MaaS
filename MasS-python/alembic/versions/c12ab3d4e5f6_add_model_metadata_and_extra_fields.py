"""add model metadata and extra fields

Revision ID: c12ab3d4e5f6
Revises: fb879fde9ca5
Create Date: 2026-02-13 13:55:00.000000

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "c12ab3d4e5f6"
down_revision: Union[str, Sequence[str], None] = "fb879fde9ca5"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    op.add_column("models", sa.Column("checksum", sa.String(length=64), nullable=True))
    op.add_column("models", sa.Column("docker_image", sa.String(length=255), nullable=True))
    op.add_column("models", sa.Column("tenant_id", sa.Uuid(), nullable=True))

    op.create_table(
        "model_metadata",
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("model_id", sa.Uuid(), nullable=False),
        sa.Column("key", sa.String(length=100), nullable=False),
        sa.Column("value", sa.String(length=500), nullable=True),
        sa.Column(
            "created_at",
            sa.DateTime(),
            server_default=sa.text("now()"),
            nullable=False,
        ),
        sa.Column(
            "updated_at",
            sa.DateTime(),
            server_default=sa.text("now()"),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(["model_id"], ["models.id"]),
        sa.PrimaryKeyConstraint("id"),
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.drop_table("model_metadata")
    op.drop_column("models", "tenant_id")
    op.drop_column("models", "docker_image")
    op.drop_column("models", "checksum")
