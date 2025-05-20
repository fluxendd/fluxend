package file

import (
	"fluxton/internal/domain/storage/file"
)

func ToCreateFileInput(request *CreateRequest) *file.CreateFileInput {
	return &file.CreateFileInput{
		ProjectUUID:  request.ProjectUUID,
		FullFileName: request.FullFileName,
		File:         request.File,
	}
}

func ToRenameFileInput(request *RenameRequest) *file.RenameFileInput {
	return &file.RenameFileInput{
		ProjectUUID:  request.ProjectUUID,
		FullFileName: request.FullFileName,
	}
}
