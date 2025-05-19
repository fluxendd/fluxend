package database

import (
	"github.com/google/uuid"
)

type Function struct {
	Name          string `db:"routine_name" json:"name"`
	Type          string `db:"routine_type" json:"type"`
	DataType      string `db:"data_type" json:"dataType"`
	TypeUdtName   string `db:"type_udt_name" json:"typeUdtName"`
	Definition    string `db:"routine_definition" json:"definition"`
	Language      string `db:"external_language" json:"language"`
	SqlDataAccess string `db:"sql_data_access" json:"sqlDataAccess"`
}

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
