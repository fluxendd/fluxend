package database

import (
	"fluxton/internal/domain/database"
)

func ToCreateIndexInput(request CreateIndexRequest) database.CreateIndexInput {
	return database.CreateIndexInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
		Columns:     request.Columns,
		IsUnique:    request.IsUnique,
	}
}

func ToCreateColumnInput(request CreateColumnRequest) database.CreateColumnInput {
	return database.CreateColumnInput{
		ProjectUUID: request.ProjectUUID,
		Columns:     request.Columns,
	}
}

func ToRenameColumnInput(request RenameColumnRequest) database.RenameColumnInput {
	return database.RenameColumnInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
	}
}
