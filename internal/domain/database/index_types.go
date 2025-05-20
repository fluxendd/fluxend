package database

import (
	"github.com/google/uuid"
)

type CreateIndexInput struct {
	ProjectUUID uuid.UUID `json:"projectUUID,omitempty"`
	Name        string    `json:"name"`
	Columns     []string  `json:"columns"`
	IsUnique    bool      `json:"is_unique"`
}
