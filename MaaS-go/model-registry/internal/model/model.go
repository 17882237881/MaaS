package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelStatus represents model status
type ModelStatus string

const (
	ModelStatusPending   ModelStatus = "pending"
	ModelStatusBuilding  ModelStatus = "building"
	ModelStatusReady     ModelStatus = "ready"
	ModelStatusDeploying ModelStatus = "deploying"
	ModelStatusRunning   ModelStatus = "running"
	ModelStatusFailed    ModelStatus = "failed"
	ModelStatusArchived  ModelStatus = "archived"
)

// ModelFramework represents supported ML frameworks
type ModelFramework string

const (
	FrameworkPyTorch    ModelFramework = "pytorch"
	FrameworkTensorFlow ModelFramework = "tensorflow"
	FrameworkONNX       ModelFramework = "onnx"
	FrameworkSKLearn    ModelFramework = "sklearn"
	FrameworkXGBoost    ModelFramework = "xgboost"
	FrameworkCustom     ModelFramework = "custom"
)

// Model represents a machine learning model
type Model struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Version     string         `gorm:"type:varchar(50);not null" json:"version"`
	Framework   ModelFramework `gorm:"type:varchar(50);not null" json:"framework"`
	Status      ModelStatus    `gorm:"type:varchar(50);default:'pending'" json:"status"`
	Size        int64          `gorm:"default:0" json:"size"`
	Checksum    string         `gorm:"type:varchar(64)" json:"checksum"`
	StoragePath string         `gorm:"type:varchar(512)" json:"storage_path"`
	DockerImage string         `gorm:"type:varchar(255)" json:"docker_image"`

	// Relationships
	Tags     []Tag          `gorm:"many2many:model_tags;" json:"tags"`
	Metadata []Metadata     `gorm:"foreignKey:ModelID" json:"metadata"`
	Versions []ModelVersion `gorm:"foreignKey:ModelID" json:"versions"`

	// Ownership
	OwnerID  string `gorm:"type:uuid;not null;index" json:"owner_id"`
	TenantID string `gorm:"type:uuid;not null;index" json:"tenant_id"`

	// Visibility
	IsPublic bool `gorm:"default:false" json:"is_public"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (Model) TableName() string {
	return "models"
}

// BeforeCreate hook to generate UUID
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// Tag represents a model tag
type Tag struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook for Tag
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Tag) TableName() string {
	return "tags"
}

// Metadata represents model metadata key-value pairs
type Metadata struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ModelID   string    `gorm:"type:uuid;not null;index" json:"model_id"`
	Key       string    `gorm:"type:varchar(100);not null" json:"key"`
	Value     string    `gorm:"type:varchar(500)" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook for Metadata
func (m *Metadata) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Metadata) TableName() string {
	return "model_metadata"
}

// ModelVersion represents a specific version of a model
type ModelVersion struct {
	ID          string      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ModelID     string      `gorm:"type:uuid;not null;index" json:"model_id"`
	Version     string      `gorm:"type:varchar(50);not null" json:"version"`
	Status      ModelStatus `gorm:"type:varchar(50);not null" json:"status"`
	Size        int64       `json:"size"`
	Checksum    string      `gorm:"type:varchar(64)" json:"checksum"`
	StoragePath string      `gorm:"type:varchar(512)" json:"storage_path"`
	DockerImage string      `gorm:"type:varchar(255)" json:"docker_image"`
	ChangeLog   string      `gorm:"type:text" json:"change_log"`
	CreatedBy   string      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
}

// BeforeCreate hook for ModelVersion
func (v *ModelVersion) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (ModelVersion) TableName() string {
	return "model_versions"
}
