package database

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"mime/multipart"
)

type Table struct {
	shared.BaseModel
	Id            int    `db:"id"`
	Name          string `db:"name"`
	Schema        string `db:"schema"`
	EstimatedRows int    `db:"estimated_rows"`
	TotalSize     string `db:"total_size"`
}

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
