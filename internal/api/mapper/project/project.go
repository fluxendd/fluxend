package project

import (
	projectDto "fluxton/internal/api/dto/project"
	projectDomain "fluxton/internal/domain/project"
)

func ToResource(project *projectDomain.Project) projectDto.Response {
	return projectDto.Response{
		Uuid:             project.Uuid,
		OrganizationUuid: project.OrganizationUuid,
		CreatedBy:        project.CreatedBy,
		UpdatedBy:        project.UpdatedBy,
		Name:             project.Name,
		Status:           project.Status,
		Description:      project.Description,
		DBName:           project.DBName,
		CreatedAt:        project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        project.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToResourceCollection(projects []projectDomain.Project) []projectDto.Response {
	resourceNotes := make([]projectDto.Response, len(projects))
	for i, currentProject := range projects {
		resourceNotes[i] = ToResource(&currentProject)
	}

	return resourceNotes
}
