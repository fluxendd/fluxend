package mapper

import (
	organizationDto "fluxton/internal/api/dto/organization"
	organizationDomain "fluxton/internal/domain/organization"
)

func ToOrganizationResource(organization *organizationDomain.Organization) organizationDto.Response {
	return organizationDto.Response{
		Uuid:      organization.Uuid,
		Name:      organization.Name,
		CreatedBy: organization.CreatedBy,
		UpdatedBy: organization.UpdatedBy,
		CreatedAt: organization.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: organization.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToOrganizationResourceCollection(organizations []organizationDomain.Organization) []organizationDto.Response {
	resourceNotes := make([]organizationDto.Response, len(organizations))
	for i, organization := range organizations {
		resourceNotes[i] = ToOrganizationResource(&organization)
	}

	return resourceNotes
}
