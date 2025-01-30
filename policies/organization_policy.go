package policies

import (
	"github.com/samber/do"
	"myapp/models"
	"myapp/repositories"
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
	return authenticatedUser.RoleID == models.UserRoleKing || authenticatedUser.RoleID == models.UserRoleBishop
}

func (s *OrganizationPolicy) CanView(organizationUserId, authenticatedUserId uint) bool {
	return organizationUserId == authenticatedUserId
}

func (s *OrganizationPolicy) CanUpdate(organizationId, authenticatedUserId uint) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationId, authenticatedUserId)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
