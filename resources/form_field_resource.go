package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type FormFieldResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	FormUuid    uuid.UUID `json:"formUuid"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	IsRequired  bool      `json:"isRequired"`
	Options     string    `json:"options"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

func FormFieldResource(formField *models.FormFiled) FormFieldResponse {
	return FormFieldResponse{
		Uuid:        formField.Uuid,
		FormUuid:    formField.FormUuid,
		Label:       formField.Label,
		Description: formField.Description,
		Type:        formField.Type,
		IsRequired:  formField.IsRequired,
		Options:     formField.Options,
		CreatedBy:   formField.CreatedBy,
		UpdatedBy:   formField.UpdatedBy,
		CreatedAt:   formField.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   formField.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FormFieldResourceCollection(formFields []models.FormFiled) []FormFieldResponse {
	resourceFormFields := make([]FormFieldResponse, len(formFields))
	for i, formField := range formFields {
		resourceFormFields[i] = FormFieldResource(&formField)
	}

	return resourceFormFields
}
