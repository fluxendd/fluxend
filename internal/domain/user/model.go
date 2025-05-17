package user

import (
	"fluxton/internal/config/constants"
	"github.com/google/uuid"
	"time"
)

const UserStatusActive = "active"
const UserStatusInactive = "inactive"

type User struct {
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

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u User) GetRoles() []int {
	return []int{constants.UserRoleOwner, constants.UserRoleAdmin, constants.UserRoleDeveloper, constants.UserRoleExplorer}
}

func (u User) IsSuperman() bool {
	return u.RoleID == constants.UserRoleSuperman
}

func (u User) IsOwner() bool {
	return u.RoleID == constants.UserRoleOwner
}

func (u User) IsAdminOrMore() bool {
	return u.RoleID <= constants.UserRoleAdmin
}

func (u User) IsDeveloperOrMore() bool {
	return u.RoleID <= constants.UserRoleDeveloper
}

func (u User) IsExplorerOrMore() bool {
	return u.RoleID <= constants.UserRoleExplorer
}
