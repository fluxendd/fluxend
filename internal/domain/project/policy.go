package project

import (
	"fluxton/internal/database/repositories"
	"fluxton/internal/domain/auth"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Policy struct {
	organizationRepo *repositories.OrganizationRepository
}

func NewProjectPolicy(injector *do.Injector) (*Policy, error) {
	repo := do.MustInvoke[*repositories.OrganizationRepository](injector)

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
