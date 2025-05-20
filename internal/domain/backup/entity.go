package backup

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type Backup struct {
	shared.BaseEntity
	Uuid        uuid.UUID  `db:"uuid" json:"uuid"`
	ProjectUuid uuid.UUID  `db:"project_uuid" json:"projectUuid"`
	Status      string     `db:"status" json:"status"`
	Error       string     `db:"error" json:"error"`
	StartedAt   time.Time  `db:"started_at" json:"startedAt"`
	CompletedAt *time.Time `db:"completed_at" json:"completedAt"`
}
