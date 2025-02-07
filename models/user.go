package models

import (
	"github.com/google/uuid"
	"time"
)

const UserStatusActive = "active"
const UserStatusInactive = "inactive"

const (
	UserRoleOwner     = 1
	UserRoleAdmin     = 2
	UserRoleDeveloper = 3
	UserRoleExplorer  = 4
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Status    string    `db:"status"`
	RoleID    int       `db:"role_id"`
	Bio       string    `db:"bio"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AuthUser struct {
	ID     uuid.UUID
	RoleID int
}

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u User) GetRoles() []int {
	return []int{UserRoleOwner, UserRoleAdmin, UserRoleDeveloper, UserRoleExplorer}
}

func (u User) IsOwner() bool {
	return u.RoleID == UserRoleOwner
}

func (u User) IsAdminOrMore() bool {
	return u.RoleID <= UserRoleAdmin
}

func (u User) IsDeveloperOrMore() bool {
	return u.RoleID <= UserRoleDeveloper
}

func (u User) IsExplorerOrMore() bool {
	return u.RoleID <= UserRoleExplorer
}

func (au AuthUser) IsOwner() bool {
	return au.RoleID == UserRoleOwner
}

func (au AuthUser) IsAdminOrMore() bool {
	return au.RoleID <= UserRoleAdmin
}

func (au AuthUser) IsDeveloperOrMore() bool {
	return au.RoleID <= UserRoleDeveloper
}

func (au AuthUser) IsExplorerOrMore() bool {
	return au.RoleID <= UserRoleExplorer
}
