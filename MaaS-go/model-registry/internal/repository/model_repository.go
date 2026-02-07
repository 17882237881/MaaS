package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"maas-platform/model-registry/internal/model"
)

var (
	ErrModelNotFound  = errors.New("model not found")
	ErrDuplicateModel = errors.New("model with this name and version already exists")
	ErrInvalidFilter  = errors.New("invalid filter parameters")
)

// ModelRepository defines the interface for model data access
type ModelRepository interface {
	Create(ctx context.Context, m *model.Model) error
	GetByID(ctx context.Context, id string) (*model.Model, error)
	GetByNameAndVersion(ctx context.Context, name, version string) (*model.Model, error)
	List(ctx context.Context, filter ModelFilter, pagination Pagination) ([]*model.Model, int64, error)
	Update(ctx context.Context, m *model.Model) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status model.ModelStatus) error

	// Tag operations
	AddTags(ctx context.Context, modelID string, tags []string) error
	RemoveTags(ctx context.Context, modelID string, tags []string) error

	// Metadata operations
	SetMetadata(ctx context.Context, modelID string, metadata map[string]string) error
	GetMetadata(ctx context.Context, modelID string) (map[string]string, error)
}

// ModelFilter defines filter criteria for listing models
type ModelFilter struct {
	Name      string
	Framework model.ModelFramework
	Status    model.ModelStatus
	OwnerID   string
	TenantID  string
	Tags      []string
	IsPublic  *bool
}

// Pagination defines pagination parameters
type Pagination struct {
	Page  int
	Limit int
}

// GormModelRepository implements ModelRepository using GORM
type GormModelRepository struct {
	db *gorm.DB
}

// NewGormModelRepository creates a new GORM model repository
func NewGormModelRepository(db *gorm.DB) ModelRepository {
	return &GormModelRepository{db: db}
}

// Create creates a new model
func (r *GormModelRepository) Create(ctx context.Context, m *model.Model) error {
	// Check for duplicate
	var existing model.Model
	result := r.db.WithContext(ctx).
		Where("name = ? AND version = ?", m.Name, m.Version).
		First(&existing)

	if result.Error == nil {
		return ErrDuplicateModel
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return r.db.WithContext(ctx).Create(m).Error
}

// GetByID retrieves a model by ID
func (r *GormModelRepository) GetByID(ctx context.Context, id string) (*model.Model, error) {
	var m model.Model
	result := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Metadata").
		First(&m, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrModelNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &m, nil
}

// GetByNameAndVersion retrieves a model by name and version
func (r *GormModelRepository) GetByNameAndVersion(ctx context.Context, name, version string) (*model.Model, error) {
	var m model.Model
	result := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Metadata").
		Where("name = ? AND version = ?", name, version).
		First(&m)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrModelNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &m, nil
}

// List retrieves a paginated list of models with optional filtering
func (r *GormModelRepository) List(ctx context.Context, filter ModelFilter, pagination Pagination) ([]*model.Model, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Model{})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", filter.Name))
	}
	if filter.Framework != "" {
		query = query.Where("framework = ?", filter.Framework)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OwnerID != "" {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}
	if len(filter.Tags) > 0 {
		query = query.Joins("JOIN model_tags ON models.id = model_tags.model_id").
			Joins("JOIN tags ON model_tags.tag_id = tags.id").
			Where("tags.name IN ?", filter.Tags)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 || pagination.Limit > 100 {
		pagination.Limit = 20
	}
	offset := (pagination.Page - 1) * pagination.Limit

	var models []*model.Model
	result := query.Preload("Tags").
		Offset(offset).
		Limit(pagination.Limit).
		Order("created_at DESC").
		Find(&models)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return models, total, nil
}

// Update updates a model
func (r *GormModelRepository) Update(ctx context.Context, m *model.Model) error {
	result := r.db.WithContext(ctx).Save(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrModelNotFound
	}
	return nil
}

// Delete soft-deletes a model
func (r *GormModelRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Model{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrModelNotFound
	}
	return nil
}

// UpdateStatus updates the status of a model
func (r *GormModelRepository) UpdateStatus(ctx context.Context, id string, status model.ModelStatus) error {
	result := r.db.WithContext(ctx).
		Model(&model.Model{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrModelNotFound
	}
	return nil
}

// AddTags adds tags to a model
func (r *GormModelRepository) AddTags(ctx context.Context, modelID string, tagNames []string) error {
	// Get or create tags
	var tags []model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		result := r.db.WithContext(ctx).Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name})
		if result.Error != nil {
			return result.Error
		}
		tags = append(tags, tag)
	}

	// Associate tags with model
	return r.db.WithContext(ctx).Model(&model.Model{ID: modelID}).Association("Tags").Append(tags)
}

// RemoveTags removes tags from a model
func (r *GormModelRepository) RemoveTags(ctx context.Context, modelID string, tagNames []string) error {
	var tags []model.Tag
	if err := r.db.WithContext(ctx).Where("name IN ?", tagNames).Find(&tags).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&model.Model{ID: modelID}).Association("Tags").Delete(tags)
}

// SetMetadata sets metadata for a model
func (r *GormModelRepository) SetMetadata(ctx context.Context, modelID string, metadata map[string]string) error {
	// Delete existing metadata
	if err := r.db.WithContext(ctx).Where("model_id = ?", modelID).Delete(&model.Metadata{}).Error; err != nil {
		return err
	}

	// Create new metadata
	for key, value := range metadata {
		m := &model.Metadata{
			ModelID: modelID,
			Key:     key,
			Value:   value,
		}
		if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetMetadata retrieves metadata for a model
func (r *GormModelRepository) GetMetadata(ctx context.Context, modelID string) (map[string]string, error) {
	var metadata []model.Metadata
	if err := r.db.WithContext(ctx).Where("model_id = ?", modelID).Find(&metadata).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, m := range metadata {
		result[m.Key] = m.Value
	}

	return result, nil
}
