package service

import (
	"context"
	"errors"
	"fmt"

	"maas-platform/model-registry/internal/model"
	"maas-platform/model-registry/internal/repository"
	"maas-platform/model-registry/pkg/logger"
)

// Service errors
var (
	ErrModelNotFound  = errors.New("model not found")
	ErrDuplicateModel = errors.New("model already exists")
	ErrInvalidInput   = errors.New("invalid input")
)

// ModelService defines the interface for model business logic
type ModelService interface {
	CreateModel(ctx context.Context, req CreateModelRequest) (*model.Model, error)
	GetModel(ctx context.Context, id string) (*model.Model, error)
	ListModels(ctx context.Context, filter ListModelsFilter) (*ListModelsResponse, error)
	UpdateModel(ctx context.Context, id string, req UpdateModelRequest) (*model.Model, error)
	DeleteModel(ctx context.Context, id string) error
	UpdateModelStatus(ctx context.Context, id string, status model.ModelStatus) error
	AddModelTags(ctx context.Context, id string, tags []string) error
	RemoveModelTags(ctx context.Context, id string, tags []string) error
	SetModelMetadata(ctx context.Context, id string, metadata map[string]string) error
	GetModelMetadata(ctx context.Context, id string) (map[string]string, error)
}

// CreateModelRequest represents a request to create a model
type CreateModelRequest struct {
	Name        string
	Description string
	Version     string
	Framework   model.ModelFramework
	Tags        []string
	Metadata    map[string]string
	OwnerID     string
	TenantID    string
	IsPublic    bool
}

// UpdateModelRequest represents a request to update a model
type UpdateModelRequest struct {
	Name        *string
	Description *string
	Tags        []string
	Metadata    map[string]string
	IsPublic    *bool
}

// ListModelsFilter represents filters for listing models
type ListModelsFilter struct {
	Name      string
	Framework model.ModelFramework
	Status    model.ModelStatus
	OwnerID   string
	TenantID  string
	Tags      []string
	IsPublic  *bool
	Page      int
	Limit     int
}

// ListModelsResponse represents the response for listing models
type ListModelsResponse struct {
	Models []*model.Model
	Total  int64
	Page   int
	Limit  int
}

// modelService implements ModelService
type modelService struct {
	repo   repository.ModelRepository
	logger *logger.Logger
}

// NewModelService creates a new model service
func NewModelService(repo repository.ModelRepository, logger *logger.Logger) ModelService {
	return &modelService{
		repo:   repo,
		logger: logger,
	}
}

// CreateModel creates a new model
func (s *modelService) CreateModel(ctx context.Context, req CreateModelRequest) (*model.Model, error) {
	// Validate framework
	if !isValidFramework(req.Framework) {
		return nil, fmt.Errorf("invalid framework: %s", req.Framework)
	}

	// Create model entity
	m := &model.Model{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Framework:   req.Framework,
		Status:      model.ModelStatusPending,
		OwnerID:     req.OwnerID,
		TenantID:    req.TenantID,
		IsPublic:    req.IsPublic,
	}

	// Save to database
	if err := s.repo.Create(ctx, m); err != nil {
		s.logger.Error("Failed to create model", "error", err)
		return nil, err
	}

	// Add tags if provided
	if len(req.Tags) > 0 {
		if err := s.repo.AddTags(ctx, m.ID, req.Tags); err != nil {
			s.logger.Error("Failed to add tags", "error", err)
		}
	}

	// Add metadata if provided
	if len(req.Metadata) > 0 {
		if err := s.repo.SetMetadata(ctx, m.ID, req.Metadata); err != nil {
			s.logger.Error("Failed to set metadata", "error", err)
		}
	}

	s.logger.Info("Model created",
		"model_id", m.ID,
		"name", m.Name,
		"version", m.Version,
	)

	return m, nil
}

// GetModel retrieves a model by ID
func (s *modelService) GetModel(ctx context.Context, id string) (*model.Model, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get model", "id", id, "error", err)
		return nil, err
	}
	return m, nil
}

// ListModels retrieves a paginated list of models
func (s *modelService) ListModels(ctx context.Context, filter ListModelsFilter) (*ListModelsResponse, error) {
	repoFilter := repository.ModelFilter{
		Name:      filter.Name,
		Framework: filter.Framework,
		Status:    filter.Status,
		OwnerID:   filter.OwnerID,
		TenantID:  filter.TenantID,
		Tags:      filter.Tags,
		IsPublic:  filter.IsPublic,
	}

	pagination := repository.Pagination{
		Page:  filter.Page,
		Limit: filter.Limit,
	}

	models, total, err := s.repo.List(ctx, repoFilter, pagination)
	if err != nil {
		s.logger.Error("Failed to list models", "error", err)
		return nil, err
	}

	return &ListModelsResponse{
		Models: models,
		Total:  total,
		Page:   filter.Page,
		Limit:  filter.Limit,
	}, nil
}

// UpdateModel updates a model
func (s *modelService) UpdateModel(ctx context.Context, id string, req UpdateModelRequest) (*model.Model, error) {
	// Get existing model
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = *req.Description
	}
	if req.IsPublic != nil {
		m.IsPublic = *req.IsPublic
	}

	// Save changes
	if err := s.repo.Update(ctx, m); err != nil {
		s.logger.Error("Failed to update model", "id", id, "error", err)
		return nil, err
	}

	// Update tags if provided
	if req.Tags != nil {
		existing := make(map[string]struct{}, len(m.Tags))
		for _, tag := range m.Tags {
			existing[tag.Name] = struct{}{}
		}

		desired := make(map[string]struct{}, len(req.Tags))
		for _, tag := range req.Tags {
			desired[tag] = struct{}{}
		}

		toAdd := make([]string, 0)
		for tag := range desired {
			if _, ok := existing[tag]; !ok {
				toAdd = append(toAdd, tag)
			}
		}

		toRemove := make([]string, 0)
		for tag := range existing {
			if _, ok := desired[tag]; !ok {
				toRemove = append(toRemove, tag)
			}
		}

		if len(toAdd) > 0 {
			if err := s.repo.AddTags(ctx, id, toAdd); err != nil {
				s.logger.Error("Failed to add model tags", "error", err)
				return nil, err
			}
		}

		if len(toRemove) > 0 {
			if err := s.repo.RemoveTags(ctx, id, toRemove); err != nil {
				s.logger.Error("Failed to remove model tags", "error", err)
				return nil, err
			}
		}
	}

	// Update metadata if provided
	if req.Metadata != nil {
		if err := s.repo.SetMetadata(ctx, id, req.Metadata); err != nil {
			s.logger.Error("Failed to update metadata", "error", err)
		}
	}

	s.logger.Info("Model updated", "model_id", id)
	return m, nil
}

// AddModelTags adds tags to a model
func (s *modelService) AddModelTags(ctx context.Context, id string, tags []string) error {
	if err := s.repo.AddTags(ctx, id, tags); err != nil {
		s.logger.Error("Failed to add model tags", "id", id, "error", err)
		return err
	}
	return nil
}

// RemoveModelTags removes tags from a model
func (s *modelService) RemoveModelTags(ctx context.Context, id string, tags []string) error {
	if err := s.repo.RemoveTags(ctx, id, tags); err != nil {
		s.logger.Error("Failed to remove model tags", "id", id, "error", err)
		return err
	}
	return nil
}

// SetModelMetadata sets metadata for a model
func (s *modelService) SetModelMetadata(ctx context.Context, id string, metadata map[string]string) error {
	if err := s.repo.SetMetadata(ctx, id, metadata); err != nil {
		s.logger.Error("Failed to set model metadata", "id", id, "error", err)
		return err
	}
	return nil
}

// GetModelMetadata gets metadata for a model
func (s *modelService) GetModelMetadata(ctx context.Context, id string) (map[string]string, error) {
	metadata, err := s.repo.GetMetadata(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get model metadata", "id", id, "error", err)
		return nil, err
	}
	return metadata, nil
}

// DeleteModel deletes a model
func (s *modelService) DeleteModel(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete model", "id", id, "error", err)
		return err
	}

	s.logger.Info("Model deleted", "model_id", id)
	return nil
}

// UpdateModelStatus updates the status of a model
func (s *modelService) UpdateModelStatus(ctx context.Context, id string, status model.ModelStatus) error {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		s.logger.Error("Failed to update model status",
			"id", id,
			"status", status,
			"error", err,
		)
		return err
	}

	s.logger.Info("Model status updated",
		"model_id", id,
		"status", status,
	)
	return nil
}

// isValidFramework checks if a framework is valid
func isValidFramework(f model.ModelFramework) bool {
	switch f {
	case model.FrameworkPyTorch, model.FrameworkTensorFlow, model.FrameworkONNX,
		model.FrameworkSKLearn, model.FrameworkXGBoost, model.FrameworkCustom:
		return true
	}
	return false
}
