package container

import (
	"fluxton/internal/domain/storage/container"
	"github.com/google/uuid"
)

type ContainerResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"isPublic"`
	Url         string    `json:"url"`
	TotalFiles  int       `json:"totalFiles"`
	MaxFileSize int       `json:"maxFileSize"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

func ContainerResource(container *container.Container) ContainerResponse {
	return ContainerResponse{
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

func ContainerResourceCollection(containers []container.Container) []ContainerResponse {
	resourceContainers := make([]ContainerResponse, len(containers))
	for i, container := range containers {
		resourceContainers[i] = ContainerResource(&container)
	}

	return resourceContainers
}
