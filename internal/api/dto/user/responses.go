package user

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid             uuid.UUID `json:"uuid"`
	OrganizationUuid uuid.UUID `json:"organizationUuid"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Status           string    `json:"status"`
	RoleID           int       `json:"roleId"`
	Bio              string    `json:"bio"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}
