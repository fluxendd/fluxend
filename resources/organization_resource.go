package resources

import "myapp/models"

type OrganizationResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func OrganizationResource(organization *models.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:        organization.ID,
		Name:      organization.Name,
		CreatedAt: organization.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: organization.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func OrganizationResourceCollection(organizations []models.Organization) []OrganizationResponse {
	resourceNotes := make([]OrganizationResponse, len(organizations))
	for i, organization := range organizations {
		resourceNotes[i] = OrganizationResource(&organization)
	}

	return resourceNotes
}
