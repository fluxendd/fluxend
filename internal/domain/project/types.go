package project

import (
	"github.com/google/uuid"
)

type CreateProjectInput struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	OrganizationUUID uuid.UUID `json:"organization_uuid"`
}

type UpdateProjectInput struct {
	ProjectUUID uuid.UUID `json:"project_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
