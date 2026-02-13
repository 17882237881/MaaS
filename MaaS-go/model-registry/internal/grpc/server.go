package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"maas-platform/model-registry/internal/model"
	"maas-platform/model-registry/internal/service"
	modelpb "maas-platform/shared/proto"
)

// GRPCServer implements the gRPC ModelService
type GRPCServer struct {
	modelpb.UnimplementedModelServiceServer
	service service.ModelService
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(svc service.ModelService) *GRPCServer {
	return &GRPCServer{
		service: svc,
	}
}

// CreateModel creates a new model via gRPC
func (s *GRPCServer) CreateModel(ctx context.Context, req *modelpb.CreateModelRequest) (*modelpb.CreateModelResponse, error) {
	createReq := service.CreateModelRequest{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Framework:   model.ModelFramework(req.Framework),
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		OwnerID:     req.OwnerId,
		TenantID:    req.TenantId,
		IsPublic:    req.IsPublic,
	}

	m, err := s.service.CreateModel(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create model: %v", err)
	}

	return &modelpb.CreateModelResponse{
		Model: convertModelToProto(m),
	}, nil
}

// GetModel retrieves a model by ID via gRPC
func (s *GRPCServer) GetModel(ctx context.Context, req *modelpb.GetModelRequest) (*modelpb.GetModelResponse, error) {
	m, err := s.service.GetModel(ctx, req.Id)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get model: %v", err)
	}

	return &modelpb.GetModelResponse{
		Model: convertModelToProto(m),
	}, nil
}

// ListModels lists models with filtering via gRPC
func (s *GRPCServer) ListModels(ctx context.Context, req *modelpb.ListModelsRequest) (*modelpb.ListModelsResponse, error) {
	filter := service.ListModelsFilter{
		Name:     req.Name,
		Page:     int(req.Page),
		Limit:    int(req.Limit),
		OwnerID:  req.OwnerId,
		TenantID: req.TenantId,
		Tags:     req.Tags,
	}
	if req.IsPublic {
		isPublic := true
		filter.IsPublic = &isPublic
	}

	if req.Framework != "" {
		filter.Framework = model.ModelFramework(req.Framework)
	}
	if req.Status != "" {
		filter.Status = model.ModelStatus(req.Status)
	}

	resp, err := s.service.ListModels(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list models: %v", err)
	}

	models := make([]*modelpb.Model, len(resp.Models))
	for i, m := range resp.Models {
		models[i] = convertModelToProto(m)
	}

	return &modelpb.ListModelsResponse{
		Models: models,
		Total:  resp.Total,
		Page:   int32(resp.Page),
		Limit:  int32(resp.Limit),
	}, nil
}

// UpdateModel updates a model via gRPC
func (s *GRPCServer) UpdateModel(ctx context.Context, req *modelpb.UpdateModelRequest) (*modelpb.UpdateModelResponse, error) {
	updateReq := service.UpdateModelRequest{
		Tags:     req.Tags,
		Metadata: req.Metadata,
	}

	if req.Name != "" {
		updateReq.Name = &req.Name
	}
	if req.Description != "" {
		updateReq.Description = &req.Description
	}
	if req.IsPublic {
		isPublic := true
		updateReq.IsPublic = &isPublic
	}

	m, err := s.service.UpdateModel(ctx, req.Id, updateReq)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update model: %v", err)
	}

	return &modelpb.UpdateModelResponse{
		Model: convertModelToProto(m),
	}, nil
}

// AddModelTags adds tags to model via gRPC
func (s *GRPCServer) AddModelTags(ctx context.Context, req *modelpb.AddModelTagsRequest) (*emptypb.Empty, error) {
	err := s.service.AddModelTags(ctx, req.ModelId, req.Tags)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to add model tags: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// RemoveModelTags removes tags from model via gRPC
func (s *GRPCServer) RemoveModelTags(ctx context.Context, req *modelpb.RemoveModelTagsRequest) (*emptypb.Empty, error) {
	err := s.service.RemoveModelTags(ctx, req.ModelId, req.Tags)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to remove model tags: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// SetModelMetadata sets metadata for model via gRPC
func (s *GRPCServer) SetModelMetadata(ctx context.Context, req *modelpb.SetModelMetadataRequest) (*emptypb.Empty, error) {
	err := s.service.SetModelMetadata(ctx, req.ModelId, req.Metadata)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to set model metadata: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetModelMetadata gets metadata for model via gRPC
func (s *GRPCServer) GetModelMetadata(ctx context.Context, req *modelpb.GetModelMetadataRequest) (*modelpb.GetModelMetadataResponse, error) {
	metadata, err := s.service.GetModelMetadata(ctx, req.ModelId)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get model metadata: %v", err)
	}

	return &modelpb.GetModelMetadataResponse{Metadata: metadata}, nil
}

// DeleteModel deletes a model via gRPC
func (s *GRPCServer) DeleteModel(ctx context.Context, req *modelpb.DeleteModelRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteModel(ctx, req.Id)
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete model: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// UpdateModelStatus updates model status via gRPC
func (s *GRPCServer) UpdateModelStatus(ctx context.Context, req *modelpb.UpdateModelStatusRequest) (*modelpb.UpdateModelStatusResponse, error) {
	err := s.service.UpdateModelStatus(ctx, req.Id, model.ModelStatus(req.Status))
	if err != nil {
		if err == service.ErrModelNotFound {
			return nil, status.Errorf(codes.NotFound, "model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update model status: %v", err)
	}

	// Get updated model
	m, err := s.service.GetModel(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated model: %v", err)
	}

	return &modelpb.UpdateModelStatusResponse{
		Model: convertModelToProto(m),
	}, nil
}

// convertModelToProto converts internal model to protobuf model
func convertModelToProto(m *model.Model) *modelpb.Model {
	return &modelpb.Model{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Version:     m.Version,
		Framework:   string(m.Framework),
		Status:      string(m.Status),
		Size:        m.Size,
		Checksum:    m.Checksum,
		StoragePath: m.StoragePath,
		DockerImage: m.DockerImage,
		Tags:        getTagNames(m.Tags),
		OwnerId:     m.OwnerID,
		TenantId:    m.TenantID,
		IsPublic:    m.IsPublic,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

// getTagNames extracts tag names from tags
func getTagNames(tags []model.Tag) []string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return names
}
