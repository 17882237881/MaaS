package service

import (
	"context"

	"maas-platform/api-gateway/pkg/grpc"
	"maas-platform/api-gateway/pkg/logger"
	modelpb "maas-platform/shared/proto"
)

// ModelServiceClient wraps the gRPC client for model operations
type ModelServiceClient struct {
	client *grpc.Client
	logger *logger.Logger
}

// NewModelServiceClient creates a new model service client
func NewModelServiceClient(client *grpc.Client, logger *logger.Logger) *ModelServiceClient {
	return &ModelServiceClient{
		client: client,
		logger: logger,
	}
}

// CreateModel creates a new model via gRPC
func (s *ModelServiceClient) CreateModel(ctx context.Context, req *modelpb.CreateModelRequest) (*modelpb.Model, error) {
	resp, err := s.client.CreateModel(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create model via gRPC", "error", err)
		return nil, err
	}
	return resp.Model, nil
}

// GetModel gets a model by ID via gRPC
func (s *ModelServiceClient) GetModel(ctx context.Context, id string) (*modelpb.Model, error) {
	resp, err := s.client.GetModel(ctx, &modelpb.GetModelRequest{Id: id})
	if err != nil {
		s.logger.Error("Failed to get model via gRPC", "error", err, "id", id)
		return nil, err
	}
	return resp.Model, nil
}

// ListModels lists models via gRPC
func (s *ModelServiceClient) ListModels(ctx context.Context, req *modelpb.ListModelsRequest) ([]*modelpb.Model, int64, error) {
	resp, err := s.client.ListModels(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list models via gRPC", "error", err)
		return nil, 0, err
	}
	return resp.Models, resp.Total, nil
}

// UpdateModel updates a model via gRPC
func (s *ModelServiceClient) UpdateModel(ctx context.Context, req *modelpb.UpdateModelRequest) (*modelpb.Model, error) {
	resp, err := s.client.UpdateModel(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update model via gRPC", "error", err, "id", req.Id)
		return nil, err
	}
	return resp.Model, nil
}

// DeleteModel deletes a model via gRPC
func (s *ModelServiceClient) DeleteModel(ctx context.Context, id string) error {
	err := s.client.DeleteModel(ctx, &modelpb.DeleteModelRequest{Id: id})
	if err != nil {
		s.logger.Error("Failed to delete model via gRPC", "error", err, "id", id)
		return err
	}
	return nil
}

// UpdateModelStatus updates model status via gRPC
func (s *ModelServiceClient) UpdateModelStatus(ctx context.Context, id string, status string) (*modelpb.Model, error) {
	resp, err := s.client.UpdateModelStatus(ctx, &modelpb.UpdateModelStatusRequest{
		Id:     id,
		Status: status,
	})
	if err != nil {
		s.logger.Error("Failed to update model status via gRPC", "error", err, "id", id)
		return nil, err
	}
	return resp.Model, nil
}

// AddModelTags adds tags to a model via gRPC
func (s *ModelServiceClient) AddModelTags(ctx context.Context, modelID string, tags []string) error {
	err := s.client.AddModelTags(ctx, &modelpb.AddModelTagsRequest{
		ModelId: modelID,
		Tags:    tags,
	})
	if err != nil {
		s.logger.Error("Failed to add model tags via gRPC", "error", err, "model_id", modelID)
		return err
	}
	return nil
}

// RemoveModelTags removes tags from a model via gRPC
func (s *ModelServiceClient) RemoveModelTags(ctx context.Context, modelID string, tags []string) error {
	err := s.client.RemoveModelTags(ctx, &modelpb.RemoveModelTagsRequest{
		ModelId: modelID,
		Tags:    tags,
	})
	if err != nil {
		s.logger.Error("Failed to remove model tags via gRPC", "error", err, "model_id", modelID)
		return err
	}
	return nil
}

// SetModelMetadata sets model metadata via gRPC
func (s *ModelServiceClient) SetModelMetadata(ctx context.Context, modelID string, metadata map[string]string) error {
	err := s.client.SetModelMetadata(ctx, &modelpb.SetModelMetadataRequest{
		ModelId:  modelID,
		Metadata: metadata,
	})
	if err != nil {
		s.logger.Error("Failed to set model metadata via gRPC", "error", err, "model_id", modelID)
		return err
	}
	return nil
}

// GetModelMetadata gets model metadata via gRPC
func (s *ModelServiceClient) GetModelMetadata(ctx context.Context, modelID string) (map[string]string, error) {
	resp, err := s.client.GetModelMetadata(ctx, &modelpb.GetModelMetadataRequest{
		ModelId: modelID,
	})
	if err != nil {
		s.logger.Error("Failed to get model metadata via gRPC", "error", err, "model_id", modelID)
		return nil, err
	}
	return resp.Metadata, nil
}
