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
