package container

import (
	"fluxend/internal/domain/storage/container"
)

func ToCreateContainerInput(request *CreateRequest) *container.CreateContainerInput {
	return &container.CreateContainerInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
		Description: request.Description,
		IsPublic:    request.IsPublic,
		MaxFileSize: request.MaxFileSize,
	}
}
