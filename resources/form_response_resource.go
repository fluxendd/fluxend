package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type FormResponseForAPI struct {
	Uuid      uuid.UUID                 `json:"uuid"`
	FormUuid  uuid.UUID                 `json:"formUuid"`
	Responses []FormFieldResponseForAPI `json:"responses"`
}

type FormFieldResponseForAPI struct {
	Uuid             uuid.UUID `json:"uuid"`
	FormResponseUuid uuid.UUID `json:"formResponseUuid"`
	FormFieldUuid    uuid.UUID `json:"formFieldUuid"`
	Value            string    `json:"value"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}

func FormResponseResource(formResponse *models.FormResponse) FormResponseForAPI {
	response := FormResponseForAPI{
		Uuid:     formResponse.Uuid,
		FormUuid: formResponse.FormUuid,
	}

	for _, formFieldResponse := range formResponse.Responses {
		response.Responses = append(response.Responses, FormFieldResponseForAPI{
			Uuid:             formFieldResponse.Uuid,
			FormResponseUuid: formFieldResponse.FormResponseUuid,
			FormFieldUuid:    formFieldResponse.FormFieldUuid,
			Value:            formFieldResponse.Value,
			CreatedAt:        formFieldResponse.CreatedAt.String(),
			UpdatedAt:        formFieldResponse.UpdatedAt.String(),
		})
	}

	return response
}

func FormResponseResourceCollection(formFields []models.FormResponse) []FormResponseForAPI {
	resourceFormFields := make([]FormResponseForAPI, len(formFields))
	for i, formField := range formFields {
		resourceFormFields[i] = FormResponseResource(&formField)
	}

	return resourceFormFields
}
