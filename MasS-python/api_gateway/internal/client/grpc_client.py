import logging

import grpc

from shared.proto import model_pb2, model_pb2_grpc


class ModelServiceClient:
    def __init__(self, target: str = "localhost:9090") -> None:
        self.target = target
        self.channel: grpc.aio.Channel | None = None
        self.stub: model_pb2_grpc.ModelServiceStub | None = None

    async def connect(self) -> None:
        if not self.channel:
            self.channel = grpc.aio.insecure_channel(self.target)
            self.stub = model_pb2_grpc.ModelServiceStub(self.channel)
            logging.info("Connected to Model Registry at %s", self.target)

    async def close(self) -> None:
        if self.channel:
            await self.channel.close()
            self.channel = None
            self.stub = None

    async def create_model(self, request: model_pb2.CreateModelRequest) -> model_pb2.Model:
        await self.connect()
        try:
            response = await self.stub.CreateModel(request)
            return response.model
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def get_model(self, model_id: str) -> model_pb2.Model:
        await self.connect()
        request = model_pb2.GetModelRequest(id=model_id)
        try:
            response = await self.stub.GetModel(request)
            return response.model
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def list_models(
        self, request: model_pb2.ListModelsRequest
    ) -> tuple[list[model_pb2.Model], int, int, int]:
        await self.connect()
        try:
            response = await self.stub.ListModels(request)
            return list(response.models), int(response.total), response.page, response.limit
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def update_model(self, request: model_pb2.UpdateModelRequest) -> model_pb2.Model:
        await self.connect()
        try:
            response = await self.stub.UpdateModel(request)
            return response.model
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def delete_model(self, model_id: str) -> None:
        await self.connect()
        request = model_pb2.DeleteModelRequest(id=model_id)
        try:
            await self.stub.DeleteModel(request)
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def update_model_status(self, model_id: str, status: str) -> model_pb2.Model:
        await self.connect()
        request = model_pb2.UpdateModelStatusRequest(id=model_id, status=status)
        try:
            response = await self.stub.UpdateModelStatus(request)
            return response.model
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def add_model_tags(self, model_id: str, tags: list[str]) -> None:
        await self.connect()
        request = model_pb2.AddModelTagsRequest(model_id=model_id, tags=tags)
        try:
            await self.stub.AddModelTags(request)
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def remove_model_tags(self, model_id: str, tags: list[str]) -> None:
        await self.connect()
        request = model_pb2.RemoveModelTagsRequest(model_id=model_id, tags=tags)
        try:
            await self.stub.RemoveModelTags(request)
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def set_model_metadata(self, model_id: str, metadata: dict[str, str]) -> None:
        await self.connect()
        request = model_pb2.SetModelMetadataRequest(model_id=model_id, metadata=metadata)
        try:
            await self.stub.SetModelMetadata(request)
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise

    async def get_model_metadata(self, model_id: str) -> dict[str, str]:
        await self.connect()
        request = model_pb2.GetModelMetadataRequest(model_id=model_id)
        try:
            response = await self.stub.GetModelMetadata(request)
            return dict(response.metadata)
        except grpc.RpcError as exc:
            logging.error("gRPC call failed: %s", exc)
            raise


ModelRegistryClient = ModelServiceClient
