package container

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type Container struct {
	shared.BaseEntity
	Uuid        uuid.UUID `db:"uuid" json:"uuid"`
	ProjectUuid uuid.UUID `db:"project_uuid" json:"projectUuid"`
	Name        string    `db:"name" json:"name"`
	NameKey     string    `db:"name_key" json:"nameKey"`
	Provider    string    `db:"provider" json:"provider"`
	Description string    `db:"description" json:"description"`
	IsPublic    bool      `db:"is_public" json:"isPublic"`
	Url         string    `db:"url" json:"url"`
	TotalFiles  int       `db:"total_files" json:"totalFiles"`
	MaxFileSize int       `db:"max_file_size" json:"maxFileSize"` // in KB
	CreatedBy   uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy   uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
