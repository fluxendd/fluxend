package models

import (
	"github.com/google/uuid"
	"time"
)

const UserStatusActive = "active"
const UserStatusInactive = "inactive"

const (
	UserRoleKing    = 1
	UserRoleBishop  = 2
	UserRoleLord    = 3
	UserRolePeasant = 4
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

type AuthenticatedUser struct {
	ID     uuid.UUID
	RoleID int
}

func (u User) GetColumns() string {
	return "id, username, email, status, role_id, bio, password, created_at, updated_at"
}

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u User) GetRoles() []int {
	return []int{UserRoleKing, UserRoleBishop, UserRoleLord, UserRolePeasant}
}

func (u User) IsKing() bool {
	return u.RoleID == UserRoleKing
}

func (u User) IsBishopOrMore() bool {
	return u.RoleID <= UserRoleBishop
}

func (u User) IsLordOrMore() bool {
	return u.RoleID <= UserRoleLord
}

func (u User) IsPeasantOrMore() bool {
	return u.RoleID <= UserRolePeasant
}

func (au AuthenticatedUser) IsKing() bool {
	return au.RoleID == UserRoleKing
}

func (au AuthenticatedUser) IsBishopOrMore() bool {
	return au.RoleID <= UserRoleBishop
}

func (au AuthenticatedUser) IsLordOrMore() bool {
	return au.RoleID <= UserRoleLord
}

func (au AuthenticatedUser) IsPeasantOrMore() bool {
	return au.RoleID <= UserRolePeasant
}
