package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type FormResponseForAPI struct {
	Uuid      uuid.UUID                  `json:"uuid"`
	FormUuid  uuid.UUID                  `json:"formUuid"`
	Responses []models.FormFieldResponse `json:"responses"`
}

func FormResponseResource(formResponse *models.FormResponse) FormResponseForAPI {
	return FormResponseForAPI{
		Uuid:      formResponse.Uuid,
		FormUuid:  formResponse.FormUuid,
		Responses: formResponse.Responses,
	}
}

func FormResponseResourceCollection(formFields []models.FormResponse) []FormResponseForAPI {
	resourceFormFields := make([]FormResponseForAPI, len(formFields))
	for i, formField := range formFields {
		resourceFormFields[i] = FormResponseResource(&formField)
	}

	return resourceFormFields
}
