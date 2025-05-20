package database

import (
	"github.com/google/uuid"
)

type CreateColumnInput struct {
	ProjectUUID uuid.UUID `json:"projectUUID,omitempty"`
	Columns     []Column  `json:"columns"`
}

type RenameColumnInput struct {
	ProjectUUID uuid.UUID `json:"projectUUID,omitempty"`
	Name        string    `json:"name"`
}
