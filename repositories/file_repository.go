package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/requests"
	"fluxton/utils"
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

func (r *FileRepository) ListForBucket(paginationParams requests.PaginationParams, bucketUUID uuid.UUID) ([]models.File, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			storage.files WHERE bucket_uuid = :bucket_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, utils.GetColumns[models.File]())

	params := map[string]interface{}{
		"bucket_uuid": bucketUUID,
		"sort":        paginationParams.Sort,
		"limit":       paginationParams.Limit,
		"offset":      offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.StructScan(&file); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return files, nil
}

func (r *FileRepository) GetByUUID(fileUUID uuid.UUID) (models.File, error) {
	query := "SELECT %s FROM storage.files WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.File]())

	var file models.File
	err := r.db.Get(&file, query, fileUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.File{}, errs.NewNotFoundError("file.error.notFound")
		}

		return models.File{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return file, nil
}

func (r *FileRepository) ExistsByUUID(bucketUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.files WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, bucketUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FileRepository) ExistsByNameForBucket(name string, bucketUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.files WHERE full_file_name = $1 AND bucket_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, bucketUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FileRepository) Create(file *models.File) (*models.File, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `
    INSERT INTO storage.files (
        bucket_uuid, full_file_name, size, mime_type, created_by, updated_by, created_at, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		file.BucketUuid,
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
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return file, nil
}

func (r *FileRepository) Rename(bucket *models.File) (*models.File, error) {
	query := `
		UPDATE storage.files 
		SET full_file_name = :full_file_name, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, bucket)
	if err != nil {
		return &models.File{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.File{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return bucket, nil
}

func (r *FileRepository) Delete(fileUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.files WHERE uuid = $1"
	res, err := r.db.Exec(query, fileUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
