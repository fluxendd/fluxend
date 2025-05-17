package form

import (
	formDto "fluxton/internal/api/dto/form"
	formDomain "fluxton/internal/domain/form"
)

func ToResource(form *formDomain.Form) formDto.Response {
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

func ToResourceCollection(forms []formDomain.Form) []formDto.Response {
	resourceForms := make([]formDto.Response, len(forms))
	for i, currentForm := range forms {
		resourceForms[i] = ToResource(&currentForm)
	}

	return resourceForms
}
