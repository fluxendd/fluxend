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

func (s *ProjectPolicy) CanCreate(organizationID uuid.UUID, authUser models.AuthUser) bool {
	if !authUser.IsDeveloperOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanAccess(organizationID uuid.UUID, authUser models.AuthUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanUpdate(organizationID uuid.UUID, authUser models.AuthUser) bool {
	if !authUser.IsDeveloperOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
