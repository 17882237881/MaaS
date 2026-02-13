import asyncio
import logging
import multiprocessing
import time
from uuid import uuid4

import grpc

from api_gateway.internal.client.grpc_client import ModelServiceClient
from model_registry.main import serve
from shared.proto import model_pb2


def start_server() -> None:
    """Starts the gRPC server in a separate process."""
    logging.basicConfig(level=logging.INFO)
    asyncio.run(serve())


async def run_client_test() -> None:
    """Runs client tests against the running server."""
    client = ModelServiceClient(target="localhost:9090")

    print("--- Testing Create Model ---")
    try:
        create_req = model_pb2.CreateModelRequest(
            name=f"grpc-test-model-{uuid4().hex[:8]}",
            description="Created via gRPC client",
            version="v1.0.0",
            framework="pytorch",
            tags=["grpc", "test"],
            metadata={"source": "test_grpc"},
            owner_id=str(uuid4()),
            tenant_id=str(uuid4()),
            is_public=True,
        )
        model = await client.create_model(create_req)
        print(f"Created Model: ID={model.id}, Name={model.name}, Status={model.status}")

        print("\n--- Testing Get Model ---")
        fetched = await client.get_model(model.id)
        print(f"Fetched Model: ID={fetched.id}, Name={fetched.name}")

        print("\n--- Testing List Models ---")
        list_req = model_pb2.ListModelsRequest(page=1, limit=10)
        models, total, _, _ = await client.list_models(list_req)
        print(f"List Models: total={total}, returned={len(models)}")

    except grpc.RpcError as e:
        print(f"RPC Failed: {e}")
    finally:
        await client.close()


if __name__ == "__main__":
    server_process = multiprocessing.Process(target=start_server)
    server_process.start()

    time.sleep(2)

    try:
        asyncio.run(run_client_test())
    finally:
        server_process.terminate()
        server_process.join()
