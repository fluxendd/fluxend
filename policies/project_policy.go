package policies

import (
	"fluxton/models"
	"fluxton/repositories"
	"github.com/google/uuid"
	"github.com/samber/do"
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

func (s *ProjectPolicy) CanCreate(organizationID uuid.UUID, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsDeveloperOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanList(organizationID uuid.UUID, authenticatedUserId uuid.UUID) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUserId)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanView(organizationID uuid.UUID, authenticatedUser models.AuthenticatedUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanUpdate(organizationID uuid.UUID, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsDeveloperOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
