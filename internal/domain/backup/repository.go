package backup

import (
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	ListForProject(projectUUID uuid.UUID) ([]Backup, error)
	GetByUUID(backupUUID uuid.UUID) (Backup, error)
	ExistsByUUID(backupUUID uuid.UUID) (bool, error)
	Create(backup *Backup) (*Backup, error)
	UpdateStatus(backupUUID uuid.UUID, status, error string, completedAt time.Time) error
	Delete(backupUUID uuid.UUID) (bool, error)
}
