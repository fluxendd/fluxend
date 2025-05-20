package project

import (
	"fluxton/internal/domain/project"
)

func ToCreateProjectInput(request *CreateRequest) *project.CreateProjectInput {
	return &project.CreateProjectInput{
		Name:             request.Name,
		Description:      request.Description,
		OrganizationUUID: request.OrganizationUUID,
	}
}

func ToUpdateProjectInput(request *UpdateRequest) *project.UpdateProjectInput {
	return &project.UpdateProjectInput{
		ProjectUUID: request.ProjectUUID,
		Name:        request.Name,
		Description: request.Description,
	}
}
