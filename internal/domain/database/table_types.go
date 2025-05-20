package database

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type CreateTableInput struct {
	ProjectUUID uuid.UUID `json:"projectUUID,omitempty"`
	Name        string    `json:"name"`
	Columns     []Column  `json:"columns"`
}

type RenameTableInput struct {
	ProjectUUID uuid.UUID `json:"projectUUID,omitempty"`
	Name        string    `json:"name"`
}

type UploadTableInput struct {
	ProjectUUID uuid.UUID             `json:"projectUUID,omitempty"`
	Name        string                `json:"name"`
	File        *multipart.FileHeader `form:"file"`
}
