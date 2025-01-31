package resources

import "myapp/models"

type ProjectResponse struct {
	ID             uint   `json:"id"`
	OrganizationID uint   `json:"organization_id"`
	Name           string `json:"name"`
	DBName         string `json:"db_name"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func ProjectResource(organization *models.Project) ProjectResponse {
	return ProjectResponse{
		ID:             organization.ID,
		OrganizationID: organization.OrganizationID,
		Name:           organization.Name,
		DBName:         organization.DBName,
		CreatedAt:      organization.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      organization.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProjectResourceCollection(organizations []models.Project) []ProjectResponse {
	resourceNotes := make([]ProjectResponse, len(organizations))
	for i, organization := range organizations {
		resourceNotes[i] = ProjectResource(&organization)
	}

	return resourceNotes
}
