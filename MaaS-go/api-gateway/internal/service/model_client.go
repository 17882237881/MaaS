package service

import (
	"context"

	"maas-platform/api-gateway/pkg/grpc"
	"maas-platform/api-gateway/pkg/logger"
	modelpb "maas-platform/shared/proto/model"
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
