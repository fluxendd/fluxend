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

func ToCreateTableInput(request CreateTableRequest) database.CreateTableInput {
	return database.CreateTableInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
		Columns:     request.Columns,
	}
}

func ToRenameTableInput(request RenameTableRequest) database.RenameTableInput {
	return database.RenameTableInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
	}
}

func ToUploadTableInput(request UploadTableRequest) database.UploadTableInput {
	return database.UploadTableInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
		File:        request.File,
	}
}
