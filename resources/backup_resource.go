package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type BackupResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	Status      string    `json:"status"`
	Error       string    `json:"error"`
	StartedAt   string    `json:"startedAt"`
	CompletedAt string    `json:"completedAt"`
}

func BackupResource(backup *models.Backup) BackupResponse {
	completedAt := ""
	if backup.CompletedAt != nil {
		completedAt = backup.CompletedAt.Format("2006-01-02 15:04:05")
	}

	return BackupResponse{
		Uuid:        backup.Uuid,
		ProjectUuid: backup.ProjectUuid,
		Status:      backup.Status,
		Error:       backup.Error,
		StartedAt:   backup.StartedAt.Format("2006-01-02 15:04:05"),
		CompletedAt: completedAt,
	}
}

func BackupResourceCollection(backups []models.Backup) []BackupResponse {
	resourceBackups := make([]BackupResponse, len(backups))
	for i, backup := range backups {
		resourceBackups[i] = BackupResource(&backup)
	}

	return resourceBackups
}
