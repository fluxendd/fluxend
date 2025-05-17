package organization

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
	UpdatedBy uuid.UUID `json:"updatedBy"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}
