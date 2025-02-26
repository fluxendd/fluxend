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

func (s *OrganizationPolicy) CanCreate(authUser models.AuthUser) bool {
	return authUser.IsAdminOrMore()
}

func (s *OrganizationPolicy) CanAccess(organizationUUID uuid.UUID, authUser models.AuthUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *OrganizationPolicy) CanUpdate(organizationUUID uuid.UUID, authUser models.AuthUser) bool {
	if !authUser.IsAdminOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
