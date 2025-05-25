package form

import (
	"fluxend/internal/domain/form"
)

func ToCreateFormInput(request *CreateRequest) *form.CreateFormInput {
	return &form.CreateFormInput{
		Name:        request.Name,
		Description: request.Description,
	}
}

func ToCreateFormFieldInput(request *CreateFormFieldsRequest) *form.CreateFormFieldsInput {
	fields := make([]form.FieldInput, len(request.Fields))
	for i, field := range request.Fields {
		fields[i] = form.FieldInput{
			Label:        field.Label,
			Type:         field.Type,
			IsRequired:   field.IsRequired,
			MinLength:    field.MinLength,
			MaxLength:    field.MaxLength,
			Pattern:      field.Pattern,
			Description:  field.Description,
			Options:      field.Options,
			DefaultValue: field.DefaultValue,
			MinValue:     field.MinValue,
			MaxValue:     field.MaxValue,
			StartDate:    field.StartDate,
			EndDate:      field.EndDate,
			DateFormat:   field.DateFormat,
		}
	}
	return &form.CreateFormFieldsInput{
		ProjectUUID: request.ProjectUUID,
		Fields:      fields,
	}
}

func ToUpdateFormFieldInput(request *UpdateFormFieldRequest) *form.UpdateFormFieldsInput {
	return &form.UpdateFormFieldsInput{
		ProjectUUID: request.ProjectUUID,
		FieldInput: form.FieldInput{
			Label:        request.Label,
			Type:         request.Type,
			IsRequired:   request.IsRequired,
			MinLength:    request.MinLength,
			MaxLength:    request.MaxLength,
			Pattern:      request.Pattern,
			Description:  request.Description,
			Options:      request.Options,
			DefaultValue: request.DefaultValue,
			MinValue:     request.MinValue,
			MaxValue:     request.MaxValue,
			StartDate:    request.StartDate,
			EndDate:      request.EndDate,
			DateFormat:   request.DateFormat,
		},
	}
}

func ToCreateFormResponseInput(request *CreateResponseRequest) *form.CreateResponseInput {
	return &form.CreateResponseInput{
		Response: request.Response,
	}
}
