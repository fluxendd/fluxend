package organization

import (
	"fluxton/internal/domain/auth"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Policy struct {
	organizationRepo *Repository
}

func NewOrganizationPolicy(injector *do.Injector) (*Policy, error) {
	repo := do.MustInvoke[*Repository](injector)

	return &Policy{
		organizationRepo: repo,
	}, nil
}

func (s *Policy) CanCreate(authUser auth.AuthUser) bool {
	return authUser.IsAdminOrMore()
}

func (s *Policy) CanAccess(organizationUUID uuid.UUID, authUser auth.AuthUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *Policy) CanUpdate(organizationUUID uuid.UUID, authUser auth.AuthUser) bool {
	if !authUser.IsAdminOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
