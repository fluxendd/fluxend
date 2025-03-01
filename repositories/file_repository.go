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

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(injector *do.Injector) (*FileRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &FileRepository{db: db}, nil
}

func (r *FileRepository) ListForBucket(paginationParams utils.PaginationParams, bucketUUID uuid.UUID) ([]models.File, error) {
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
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.StructScan(&file); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
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

		return models.File{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return file, nil
}

func (r *FileRepository) ExistsByUUID(bucketUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.files WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, bucketUUID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *FileRepository) ExistsByNameForBucket(name string, bucketUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.files WHERE name = $1 AND bucket_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, bucketUUID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *FileRepository) Create(file *models.File) (*models.File, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	query := `
    INSERT INTO storage.files (
        name, path, size, mime_type, created_by
    ) VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		file.Name,
		file.Path,
		file.Size,
		file.MimeType,
		file.CreatedBy,
	).Scan(&file.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("could not create file: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return file, nil
}

func (r *FileRepository) Rename(bucket *models.File) (*models.File, error) {
	query := `
		UPDATE storage.files 
		SET name = :name, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, bucket)
	if err != nil {
		return &models.File{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.File{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return bucket, nil
}

func (r *FileRepository) Delete(fileUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.files WHERE uuid = $1"
	res, err := r.db.Exec(query, fileUUID)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
