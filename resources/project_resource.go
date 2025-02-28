package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type ProjectResponse struct {
	Uuid             uuid.UUID `json:"uuid"`
	OrganizationUuid uuid.UUID `json:"organizationUuid"`
	CreatedBy        uuid.UUID `json:"createdBy"`
	UpdatedBy        uuid.UUID `json:"updatedBy"`
	Name             string    `json:"name"`
	DBName           string    `json:"dbName"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}

func ProjectResource(project *models.Project) ProjectResponse {
	return ProjectResponse{
		Uuid:             project.Uuid,
		OrganizationUuid: project.OrganizationUuid,
		CreatedBy:        project.CreatedBy,
		UpdatedBy:        project.UpdatedBy,
		Name:             project.Name,
		DBName:           project.DBName,
		CreatedAt:        project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProjectResourceCollection(organizations []models.Project) []ProjectResponse {
	resourceNotes := make([]ProjectResponse, len(organizations))
	for i, organization := range organizations {
		resourceNotes[i] = ProjectResource(&organization)
	}

	return resourceNotes
}
