package project

import (
	"fluxton/internal/domain/project"
	"github.com/google/uuid"
)

type ProjectResponse struct {
	Uuid             uuid.UUID `json:"uuid"`
	OrganizationUuid uuid.UUID `json:"organizationUuid"`
	CreatedBy        uuid.UUID `json:"createdBy"`
	UpdatedBy        uuid.UUID `json:"updatedBy"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Description      string    `json:"description"`
	DBName           string    `json:"dbName"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}

func ProjectResource(project *project.Project) ProjectResponse {
	return ProjectResponse{
		Uuid:             project.Uuid,
		OrganizationUuid: project.OrganizationUuid,
		CreatedBy:        project.CreatedBy,
		UpdatedBy:        project.UpdatedBy,
		Name:             project.Name,
		Status:           project.Status,
		Description:      project.Description,
		DBName:           project.DBName,
		CreatedAt:        project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProjectResourceCollection(organizations []project.Project) []ProjectResponse {
	resourceNotes := make([]ProjectResponse, len(organizations))
	for i, organization := range organizations {
		resourceNotes[i] = ProjectResource(&organization)
	}

	return resourceNotes
}
