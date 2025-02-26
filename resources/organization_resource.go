package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type OrganizationResponse struct {
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func OrganizationResource(organization *models.Organization) OrganizationResponse {
	return OrganizationResponse{
		Uuid:      organization.Uuid,
		Name:      organization.Name,
		CreatedBy: organization.CreatedBy,
		UpdatedBy: organization.UpdatedBy,
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
