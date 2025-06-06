package mapper

import (
	backupDto "fluxend/internal/api/dto/backup"
	"fluxend/internal/domain/backup"
)

func ToBackupResource(backup *backup.Backup) backupDto.Response {
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

func ToBackupResourceCollection(backups []backup.Backup) []backupDto.Response {
	resourceBackups := make([]backupDto.Response, len(backups))
	for i, currentBackup := range backups {
		resourceBackups[i] = ToBackupResource(&currentBackup)
	}

	return resourceBackups
}
