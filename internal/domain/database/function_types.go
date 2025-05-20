package database

import (
	"github.com/google/uuid"
)

type FunctionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateFunctionInput struct {
	ProjectUUID uuid.UUID           `json:"projectUUID,omitempty"`
	Name        string              `json:"name"`
	Parameters  []FunctionParameter `json:"parameters"`
	Definition  string              `json:"definition"`
	Language    string              `json:"language"`
	ReturnType  string              `json:"return_type"`
}
