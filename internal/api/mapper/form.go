package mapper

import (
	formDto "fluxend/internal/api/dto/form"
	formDomain "fluxend/internal/domain/form"
)

func ToFormResource(form *formDomain.Form) formDto.Response {
	return formDto.Response{
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

func ToFormResourceCollection(forms []formDomain.Form) []formDto.Response {
	resourceForms := make([]formDto.Response, len(forms))
	for i, currentForm := range forms {
		resourceForms[i] = ToFormResource(&currentForm)
	}

	return resourceForms
}

func ToFieldResource(formField *formDomain.Field) formDto.FieldResponseApi {
	return formDto.FieldResponseApi{
		Uuid:         formField.Uuid,
		FormUuid:     formField.FormUuid,
		Label:        formField.Label,
		Description:  formField.Description,
		Type:         formField.Type,
		IsRequired:   formField.IsRequired,
		Options:      formField.Options,
		MinLength:    formField.MinLength,
		MaxLength:    formField.MaxLength,
		MinValue:     formField.MinValue,
		MaxValue:     formField.MaxValue,
		Pattern:      formField.Pattern,
		DefaultValue: formField.DefaultValue,
		StartDate:    formField.StartDate,
		EndDate:      formField.EndDate,
		DateFormat:   formField.DateFormat,
		CreatedAt:    formField.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    formField.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToFieldResourceCollection(formFields []formDomain.Field) []formDto.FieldResponseApi {
	resourceFormFields := make([]formDto.FieldResponseApi, len(formFields))
	for i, formField := range formFields {
		resourceFormFields[i] = ToFieldResource(&formField)
	}

	return resourceFormFields
}

func ToResponseResource(formResponse *formDomain.FormResponse) formDto.ResponseForAPI {
	response := formDto.ResponseForAPI{
		Uuid:     formResponse.Uuid,
		FormUuid: formResponse.FormUuid,
	}

	for _, formFieldResponse := range formResponse.Responses {
		response.Responses = append(response.Responses, formDto.FieldResponseForAPI{
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

func ToResponseResourceCollection(formFields []formDomain.FormResponse) []formDto.ResponseForAPI {
	resourceFormFields := make([]formDto.ResponseForAPI, len(formFields))
	for i, formField := range formFields {
		resourceFormFields[i] = ToResponseResource(&formField)
	}

	return resourceFormFields
}
