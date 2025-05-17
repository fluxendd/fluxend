package project

import (
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/organization"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Policy struct {
	organizationRepo *organization.Repository
}

func NewProjectPolicy(injector *do.Injector) (*Policy, error) {
	repo := do.MustInvoke[*organization.Repository](injector)

	return &Policy{
		organizationRepo: repo,
	}, nil
}

func (s *Policy) CanCreate(organizationUUID uuid.UUID, authUser auth.AuthUser) bool {
	if !authUser.IsDeveloperOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *Policy) CanAccess(organizationUUID uuid.UUID, authUser auth.AuthUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *Policy) CanUpdate(organizationUUID uuid.UUID, authUser auth.AuthUser) bool {
	if !authUser.IsDeveloperOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationMember(organizationUUID, authUser.Uuid)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
