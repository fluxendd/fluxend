package policies

import (
	"fluxton/models"
	"fluxton/repositories"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type OrganizationPolicy struct {
	organizationRepo *repositories.OrganizationRepository
}

func NewOrganizationPolicy(injector *do.Injector) (*OrganizationPolicy, error) {
	repo := do.MustInvoke[*repositories.OrganizationRepository](injector)

	return &OrganizationPolicy{
		organizationRepo: repo,
	}, nil
}

func (s *OrganizationPolicy) CanCreate(authenticatedUser models.AuthenticatedUser) bool {
	return authenticatedUser.IsBishopOrMore()
}

func (s *OrganizationPolicy) CanView(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *OrganizationPolicy) CanUpdate(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsBishopOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
