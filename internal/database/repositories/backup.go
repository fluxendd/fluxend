package repositories

import (
	"fluxend/internal/domain/backup"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type BackupRepository struct {
	db shared.DB
}

func NewBackupRepository(injector *do.Injector) (backup.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &BackupRepository{db: db}, nil
}

func (r *BackupRepository) ListForProject(projectUUID uuid.UUID) ([]backup.Backup, error) {
	query := `
       SELECT %s FROM storage.backups WHERE project_uuid = :project_uuid
       ORDER BY started_at DESC
    `

	query = fmt.Sprintf(query, pkg.GetColumns[backup.Backup]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
	}

	var backups []backup.Backup
	return backups, r.db.SelectNamedList(&backups, query, params)
}

func (r *BackupRepository) GetByUUID(backupUUID uuid.UUID) (backup.Backup, error) {
	query := "SELECT %s FROM storage.backups WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[backup.Backup]())

	var form backup.Backup
	return form, r.db.GetWithNotFound(&form, "backup.error.notFound", query, backupUUID)
}

func (r *BackupRepository) ExistsByUUID(backupUUID uuid.UUID) (bool, error) {
	return r.db.Exists("storage.backups", "uuid = $1", backupUUID)
}

func (r *BackupRepository) Create(backup *backup.Backup) (*backup.Backup, error) {
	return backup, r.db.WithTransaction(func(tx shared.Tx) error {
		query := `
        INSERT INTO storage.backups (
            project_uuid, status, error, started_at
        ) VALUES (
            $1, $2, $3, $4
        )
        RETURNING uuid
        `

		return tx.QueryRowx(
			query,
			backup.ProjectUuid, backup.Status, backup.Error, backup.StartedAt,
		).Scan(&backup.Uuid)
	})
}

func (r *BackupRepository) UpdateStatus(backupUUID uuid.UUID, status, error string, completedAt time.Time) error {
	_, err := r.db.ExecWithRowsAffected("UPDATE storage.backups SET status = $1, error = $2, completed_at = $3 WHERE uuid = $4", status, error, completedAt, backupUUID)
	return err
}

func (r *BackupRepository) Delete(backupUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM storage.backups WHERE uuid = $1", backupUUID)
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
