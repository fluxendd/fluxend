package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type FormResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	CreatedBy   uuid.UUID `json:"createdB"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

func FormResource(form *models.Form) FormResponse {
	return FormResponse{
		Uuid:        form.Uuid,
		Name:        form.Name,
		Description: form.Description,
		ProjectUuid: form.ProjectUuid,
		CreatedBy:   form.CreatedBy,
		UpdatedBy:   form.UpdatedBy,
		CreatedAt:   form.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   form.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FormResourceCollection(forms []models.Form) []FormResponse {
	resourceForms := make([]FormResponse, len(forms))
	for i, organization := range forms {
		resourceForms[i] = FormResource(&organization)
	}

	return resourceForms
}
