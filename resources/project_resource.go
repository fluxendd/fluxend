package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type ProjectResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedBy      uuid.UUID `json:"created_by"`
	UpdatedBy      uuid.UUID `json:"updated_by"`
	Name           string    `json:"name"`
	DBName         string    `json:"db_name"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func ProjectResource(project *models.Project) ProjectResponse {
	return ProjectResponse{
		ID:             project.ID,
		OrganizationID: project.OrganizationID,
		CreatedBy:      project.CreatedBy,
		UpdatedBy:      project.UpdatedBy,
		Name:           project.Name,
		DBName:         project.DBName,
		CreatedAt:      project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProjectResourceCollection(organizations []models.Project) []ProjectResponse {
	resourceNotes := make([]ProjectResponse, len(organizations))
	for i, organization := range organizations {
		resourceNotes[i] = ProjectResource(&organization)
	}

	return resourceNotes
}
