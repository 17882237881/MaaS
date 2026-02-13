import asyncio
import logging
import grpc
from shared.proto import model_pb2_grpc
from model_registry.internal.server.grpc_server import ModelServiceServicer
from api_gateway.internal.config.settings import settings

async def serve():
    server = grpc.aio.server()
    model_pb2_grpc.add_ModelServiceServicer_to_server(
        ModelServiceServicer(), server
    )
    listen_addr = f"[::]:{settings.grpc.model_registry_port}"
    server.add_insecure_port(listen_addr)
    logging.info(f"Starting Model Registry gRPC server on {listen_addr}")
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    try:
        asyncio.run(serve())
    except KeyboardInterrupt:
        pass
