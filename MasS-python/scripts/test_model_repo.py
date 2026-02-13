import asyncio
import uuid
from api_gateway.internal.repository.database import SESSION_FACTORY, init_db
from api_gateway.internal.repository.model_repository import SqlAlchemyModelRepository
from api_gateway.internal.model.model_registry import Model

async def main():
    # Ensure tables exist (redundant if alembic ran, but harmless)
    # await init_db() 
    
    async with SESSION_FACTORY() as session:
        repo = SqlAlchemyModelRepository(session)
        
        # 1. Create
        print("--- Testing Create ---")
        model_id = uuid.uuid4()
        new_model = Model(
            id=model_id,
            name="test-model",
            version="v1.0.0",
            framework="pytorch",
            owner_id=uuid.uuid4()
        )
        created_model = await repo.create(new_model)
        print(f"Created model: {created_model.name} ({created_model.id})")
        
        # 2. Get by ID
        print("\n--- Testing Get by ID ---")
        fetched_model = await repo.get_by_id(model_id)
        if fetched_model:
            print(f"Fetched model: {fetched_model.name} - {fetched_model.version}")
        else:
            print("Model not found!")

        # 3. Get by Name/Version
        print("\n--- Testing Get by Name/Version ---")
        nv_model = await repo.get_by_name_version("test-model", "v1.0.0")
        if nv_model:
             print(f"Fetched by N/V: {nv_model.name}")

        # 4. List
        print("\n--- Testing List ---")
        models = await repo.list()
        print(f"List count: {len(models)}")
        
        # 5. Delete
        print("\n--- Testing Delete ---")
        deleted = await repo.delete(model_id)
        print(f"Deleted: {deleted}")
        
        # Verify deletion
        gone_model = await repo.get_by_id(model_id)
        print(f"Exists after delete? {gone_model is not None}")

if __name__ == "__main__":
    asyncio.run(main())
