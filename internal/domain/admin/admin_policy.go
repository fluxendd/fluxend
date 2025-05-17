package admin

import (
	"fluxton/models"
)

type AdminPolicy struct {
}

func NewAdminPolicy() *AdminPolicy {
	return &AdminPolicy{}
}

func (s *AdminPolicy) CanCreate(authUser models.AuthUser) bool {
	return authUser.IsSuperman()
}

func (s *AdminPolicy) CanAccess(authUser models.AuthUser) bool {
	return authUser.IsSuperman()
}

func (s *AdminPolicy) CanUpdate(authUser models.AuthUser) bool {
	return authUser.IsSuperman()
}
