package models

import (
	"github.com/google/uuid"
	"time"
)

const UserStatusActive = "active"
const UserStatusInactive = "inactive"

const (
	UserRoleSuperman  = 1
	UserRoleOwner     = 2
	UserRoleAdmin     = 3
	UserRoleDeveloper = 4
	UserRoleExplorer  = 5
)

type User struct {
	BaseModel
	Uuid      uuid.UUID `db:"uuid"`
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
	Uuid   uuid.UUID
	RoleID int
}

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u User) GetRoles() []int {
	return []int{UserRoleOwner, UserRoleAdmin, UserRoleDeveloper, UserRoleExplorer}
}

func (u User) IsSuperman() bool {
	return u.RoleID == UserRoleSuperman
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

func (au AuthUser) IsSuperman() bool {
	return au.RoleID == UserRoleSuperman
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
