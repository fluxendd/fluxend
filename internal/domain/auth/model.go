package auth

import (
	"fluxton/internal/config/constants"
	"github.com/google/uuid"
)

type AuthUser struct {
	Uuid   uuid.UUID
	RoleID int
}

func (au AuthUser) IsOwner() bool {
	return au.RoleID == constants.UserRoleOwner
}

func (au AuthUser) IsSuperman() bool {
	return au.RoleID == constants.UserRoleSuperman
}

func (au AuthUser) IsAdminOrMore() bool {
	return au.RoleID <= constants.UserRoleAdmin
}

func (au AuthUser) IsDeveloperOrMore() bool {
	return au.RoleID <= constants.UserRoleDeveloper
}

func (au AuthUser) IsExplorerOrMore() bool {
	return au.RoleID <= constants.UserRoleExplorer
}
