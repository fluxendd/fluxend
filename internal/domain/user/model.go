package user

import (
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type User struct {
	shared.BaseModel
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

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	Bio string `json:"bio"`
}

func (u User) IsActive() bool {
	return u.Status == constants.UserStatusActive
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
