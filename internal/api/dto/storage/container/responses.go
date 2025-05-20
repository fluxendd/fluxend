package container

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid        uuid.UUID `json:"uuid"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"isPublic"`
	Url         string    `json:"url"`
	TotalFiles  int       `json:"totalFiles"`
	MaxFileSize int       `json:"maxFileSize"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}
