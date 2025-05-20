package admin

import (
	"fluxton/internal/domain/auth"
)

type Policy struct {
}

func NewAdminPolicy() *Policy {
	return &Policy{}
}

func (s *Policy) CanCreate(authUser auth.User) bool {
	return authUser.IsSuperman()
}

func (s *Policy) CanAccess(authUser auth.User) bool {
	return authUser.IsSuperman()
}

func (s *Policy) CanUpdate(authUser auth.User) bool {
	return authUser.IsSuperman()
}
