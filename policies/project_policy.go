package policies

import (
	"github.com/samber/do"
	"myapp/models"
	"myapp/repositories"
)

type ProjectPolicy struct {
	organizationRepo *repositories.OrganizationRepository
}

func NewProjectPolicy(injector *do.Injector) (*ProjectPolicy, error) {
	repo := do.MustInvoke[*repositories.OrganizationRepository](injector)

	return &ProjectPolicy{
		organizationRepo: repo,
	}, nil
}

func (s *ProjectPolicy) CanCreate(organizationId uint, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsLordOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanView(organizationId uint, authenticatedUser models.AuthenticatedUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanUpdate(organizationId uint, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsLordOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
