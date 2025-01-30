package models

import (
	"fmt"
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
	ID        uint      `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Status    string    `db:"status"`
	RoleID    int       `db:"role_id"`
	Bio       string    `db:"bio"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u User) GetFields() string {
	return fmt.Sprintf("id, username, email, status, role, bio, password, created_at, updated_at")
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
