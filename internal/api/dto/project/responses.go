package project

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid             uuid.UUID `json:"uuid"`
	OrganizationUuid uuid.UUID `json:"organizationUuid"`
	CreatedBy        uuid.UUID `json:"createdBy"`
	UpdatedBy        uuid.UUID `json:"updatedBy"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Description      string    `json:"description"`
	DBName           string    `json:"dbName"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}
