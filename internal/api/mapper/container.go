package mapper

import (
	containerDto "fluxend/internal/api/dto/storage/container"
	containerDomain "fluxend/internal/domain/storage/container"
)

func ToContainerResource(container *containerDomain.Container) containerDto.Response {
	return containerDto.Response{
		Uuid:        container.Uuid,
		ProjectUuid: container.ProjectUuid,
		Name:        container.Name,
		Description: container.Description,
		IsPublic:    container.IsPublic,
		Url:         container.Url,
		TotalFiles:  container.TotalFiles,
		MaxFileSize: container.MaxFileSize,
		CreatedBy:   container.CreatedBy,
		UpdatedBy:   container.UpdatedBy,
		CreatedAt:   container.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   container.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToContainerResourceCollection(containers []containerDomain.Container) []containerDto.Response {
	resourceContainers := make([]containerDto.Response, len(containers))
	for i, currentContainer := range containers {
		resourceContainers[i] = ToContainerResource(&currentContainer)
	}

	return resourceContainers
}
