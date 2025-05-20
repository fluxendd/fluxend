package backup

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid        uuid.UUID `json:"uuid"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	Status      string    `json:"status"`
	Error       string    `json:"error"`
	StartedAt   string    `json:"startedAt"`
	CompletedAt string    `json:"completedAt"`
}
