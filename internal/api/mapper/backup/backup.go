package backup

import (
	backupDto "fluxton/internal/api/dto/backup"
	"fluxton/internal/domain/backup"
)

func ToResource(backup *backup.Backup) backupDto.Response {
	completedAt := ""
	if backup.CompletedAt != nil {
		completedAt = backup.CompletedAt.Format("2006-01-02 15:04:05")
	}

	return backupDto.Response{
		Uuid:        backup.Uuid,
		ProjectUuid: backup.ProjectUuid,
		Status:      backup.Status,
		Error:       backup.Error,
		StartedAt:   backup.StartedAt.Format("2006-01-02 15:04:05"),
		CompletedAt: completedAt,
	}
}

func ToResourceCollection(backups []backup.Backup) []backupDto.Response {
	resourceBackups := make([]backupDto.Response, len(backups))
	for i, currentBackup := range backups {
		resourceBackups[i] = ToResource(&currentBackup)
	}

	return resourceBackups
}
