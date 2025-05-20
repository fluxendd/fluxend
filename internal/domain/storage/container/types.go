package container

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateContainerInput struct {
	Context     echo.Context
	ProjectUUID uuid.UUID `db:"project_uuid" json:"projectUUID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	MaxFileSize int       `json:"max_file_size"`
}
