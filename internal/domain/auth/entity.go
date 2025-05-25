package auth

import (
	"fluxend/internal/config/constants"
	"github.com/google/uuid"
)

type User struct {
	Uuid   uuid.UUID
	RoleID int
}

func (au User) IsOwner() bool {
	return au.RoleID == constants.UserRoleOwner
}

func (au User) IsSuperman() bool {
	return au.RoleID == constants.UserRoleSuperman
}

func (au User) IsAdminOrMore() bool {
	return au.RoleID <= constants.UserRoleAdmin
}

func (au User) IsDeveloperOrMore() bool {
	return au.RoleID <= constants.UserRoleDeveloper
}

func (au User) IsExplorerOrMore() bool {
	return au.RoleID <= constants.UserRoleExplorer
}
