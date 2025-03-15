package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type BackupRepository struct {
	db *sqlx.DB
}

func NewBackupRepository(injector *do.Injector) (*BackupRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &BackupRepository{db: db}, nil
}

func (r *BackupRepository) ListForProject(projectUUID uuid.UUID) ([]models.Backup, error) {
	query := `
		SELECT %s FROM storage.backups WHERE project_uuid = :project_uuid
		ORDER BY completed_at DESC
	`

	query = fmt.Sprintf(query, utils.GetColumns[models.Backup]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var backups []models.Backup
	for rows.Next() {
		var form models.Backup
		if err := rows.StructScan(&form); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		backups = append(backups, form)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return backups, nil
}

func (r *BackupRepository) GetByUUID(backupUUID uuid.UUID) (models.Backup, error) {
	query := "SELECT %s FROM storage.backups WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Backup]())

	var form models.Backup
	err := r.db.Get(&form, query, backupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Backup{}, errs.NewNotFoundError("backup.error.notFound")
		}

		return models.Backup{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return form, nil
}

func (r *BackupRepository) ExistsByUUID(backupUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.backups WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, backupUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *BackupRepository) Create(backup *models.Backup) (*models.Backup, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `
    INSERT INTO storage.backups (
        project_uuid, status, error, started_at
    ) VALUES (
        $1, $2, $3, $4
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		backup.ProjectUuid, backup.Status, backup.Error, backup.Status,
	).Scan(&backup.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return backup, nil
}

func (r *BackupRepository) Update(backup *models.Backup) (*models.Backup, error) {
	query := `
		UPDATE storage.backups 
		SET status = :status, error = :error, completed_at = :completed_at
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, backup)
	if err != nil {
		return &models.Backup{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Backup{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return backup, nil
}

func (r *BackupRepository) Delete(backupUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.backups WHERE uuid = $1"
	res, err := r.db.Exec(query, backupUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
