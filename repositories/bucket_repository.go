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

type BucketRepository struct {
	db *sqlx.DB
}

func NewBucketRepository(injector *do.Injector) (*BucketRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &BucketRepository{db: db}, nil
}

func (r *BucketRepository) ListForProject(paginationParams utils.PaginationParams, projectUUID uuid.UUID) ([]models.Bucket, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			storage.buckets WHERE project_uuid = :project_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, utils.GetColumns[models.Bucket]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var buckets []models.Bucket
	for rows.Next() {
		var bucket models.Bucket
		if err := rows.StructScan(&bucket); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		buckets = append(buckets, bucket)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return buckets, nil
}

func (r *BucketRepository) GetByUUID(bucketUUID uuid.UUID) (models.Bucket, error) {
	query := "SELECT %s FROM storage.buckets WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Bucket]())

	var bucket models.Bucket
	err := r.db.Get(&bucket, query, bucketUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Bucket{}, errs.NewNotFoundError("bucket.error.notFound")
		}

		return models.Bucket{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return bucket, nil
}

func (r *BucketRepository) ExistsByUUID(bucketUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.buckets WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, bucketUUID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *BucketRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.buckets WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, projectUUID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *BucketRepository) Create(bucket *models.Bucket) (*models.Bucket, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	query := `
    INSERT INTO storage.buckets (
        project_uuid, name, aws_name, description, is_public, url, max_file_size, created_by, updated_by
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		bucket.ProjectUuid,
		bucket.Name,
		bucket.AwsName,
		bucket.Description,
		bucket.IsPublic,
		bucket.Url,
		bucket.MaxFileSize,
		bucket.CreatedBy,
		bucket.UpdatedBy,
	).Scan(&bucket.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("could not create bucket: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return bucket, nil
}

func (r *BucketRepository) Update(bucket *models.Bucket) (*models.Bucket, error) {
	query := `
		UPDATE storage.buckets 
		SET name = :name, description = :description, is_public = :is_public, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, bucket)
	if err != nil {
		return &models.Bucket{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Bucket{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return bucket, nil
}

func (r *BucketRepository) IncrementTotalFiles(bucketUUID uuid.UUID) error {
	query := "UPDATE storage.buckets SET total_files = total_files + 1 WHERE uuid = $1"
	_, err := r.db.Exec(query, bucketUUID)
	if err != nil {
		return fmt.Errorf("could not update row: %v", err)
	}

	return nil
}

func (r *BucketRepository) DecrementTotalFiles(bucketUUID uuid.UUID) error {
	query := "UPDATE storage.buckets SET total_files = total_files - 1 WHERE uuid = $1"
	_, err := r.db.Exec(query, bucketUUID)
	if err != nil {
		return fmt.Errorf("could not update row: %v", err)
	}

	return nil
}

func (r *BucketRepository) Delete(bucketUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.buckets WHERE uuid = $1"
	res, err := r.db.Exec(query, bucketUUID)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
