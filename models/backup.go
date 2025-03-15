package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	BackupStatusPending   = "pending"
	BackupStatusCompleted = "completed"
	BackupStatusFailed    = "failed"
)

type Backup struct {
	Uuid        uuid.UUID  `db:"uuid" json:"uuid"`
	ProjectUuid uuid.UUID  `db:"project_uuid" json:"projectUuid"`
	Status      string     `db:"status" json:"status"`
	Error       string     `db:"error" json:"error"`
	StartedAt   time.Time  `db:"started_at" json:"startedAt"`
	CompletedAt *time.Time `db:"completed_at" json:"completedAt"`
}
