package repositories

import (
	"database/sql"
	"errors"
	"fluxton/internal/api/dto"
	"fluxton/models"
	"fluxton/pkg"
	flxErrs "fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(injector *do.Injector) (*FileRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &FileRepository{db: db}, nil
}

func (r *FileRepository) ListForContainer(paginationParams dto.PaginationParams, containerUUID uuid.UUID) ([]models.File, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			storage.files WHERE container_uuid = :container_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, pkg.GetColumns[models.File]())

	params := map[string]interface{}{
		"container_uuid": containerUUID,
		"sort":           paginationParams.Sort,
		"limit":          paginationParams.Limit,
		"offset":         offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.StructScan(&file); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return files, nil
}

func (r *FileRepository) GetByUUID(fileUUID uuid.UUID) (models.File, error) {
	query := "SELECT %s FROM storage.files WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[models.File]())

	var file models.File
	err := r.db.Get(&file, query, fileUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.File{}, flxErrs.NewNotFoundError("file.error.notFound")
		}

		return models.File{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return file, nil
}

func (r *FileRepository) ExistsByUUID(containerUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.files WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, containerUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *FileRepository) ExistsByNameForContainer(name string, containerUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.files WHERE full_file_name = $1 AND container_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, containerUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *FileRepository) Create(file *models.File) (*models.File, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	query := `
    INSERT INTO storage.files (
        container_uuid, full_file_name, size, mime_type, created_by, updated_by, created_at, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		file.ContainerUuid,
		file.FullFileName,
		file.Size,
		file.MimeType,
		file.CreatedBy,
		file.UpdatedBy,
		file.CreatedAt,
		file.UpdatedAt,
	).Scan(&file.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, pkg.FormatError(queryErr, "insert", pkg.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return file, nil
}

func (r *FileRepository) Rename(container *models.File) (*models.File, error) {
	query := `
		UPDATE storage.files 
		SET full_file_name = :full_file_name, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, container)
	if err != nil {
		return &models.File{}, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.File{}, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return container, nil
}

func (r *FileRepository) Delete(fileUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.files WHERE uuid = $1"
	res, err := r.db.Exec(query, fileUUID)
	if err != nil {
		return false, pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}
