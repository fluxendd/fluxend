package form

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	CreatedBy   uuid.UUID `json:"createdB"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}
