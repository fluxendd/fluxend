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

type ContainerRepository struct {
	db *sqlx.DB
}

func NewContainerRepository(injector *do.Injector) (*ContainerRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &ContainerRepository{db: db}, nil
}

func (r *ContainerRepository) ListForProject(paginationParams requests.PaginationParams, projectUUID uuid.UUID) ([]models.Container, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			storage.containers WHERE project_uuid = :project_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, utils.GetColumns[models.Container]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var containers []models.Container
	for rows.Next() {
		var container models.Container
		if err := rows.StructScan(&container); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		containers = append(containers, container)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return containers, nil
}

func (r *ContainerRepository) GetByUUID(containerUUID uuid.UUID) (models.Container, error) {
	query := "SELECT %s FROM storage.containers WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Container]())

	var container models.Container
	err := r.db.Get(&container, query, containerUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Container{}, errs.NewNotFoundError("container.error.notFound")
		}

		return models.Container{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return container, nil
}

func (r *ContainerRepository) ExistsByUUID(containerUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.containers WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, containerUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *ContainerRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.containers WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, projectUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *ContainerRepository) Create(container *models.Container) (*models.Container, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `
    INSERT INTO storage.containers (
        project_uuid, name, aws_name, description, is_public, url, max_file_size, created_by, updated_by
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		container.ProjectUuid,
		container.Name,
		container.NameKey,
		container.Description,
		container.IsPublic,
		container.Url,
		container.MaxFileSize,
		container.CreatedBy,
		container.UpdatedBy,
	).Scan(&container.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, utils.FormatError(err, "create", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return container, nil
}

func (r *ContainerRepository) Update(container *models.Container) (*models.Container, error) {
	query := `
		UPDATE storage.containers 
		SET 
		    name = :name, 
		    description = :description, 
		    is_public = :is_public, 
		    updated_at = :updated_at, 
		    updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, container)
	if err != nil {
		return &models.Container{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Container{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return container, nil
}

func (r *ContainerRepository) IncrementTotalFiles(containerUUID uuid.UUID) error {
	query := "UPDATE storage.containers SET total_files = total_files + 1 WHERE uuid = $1"
	_, err := r.db.Exec(query, containerUUID)
	if err != nil {
		return utils.FormatError(err, "update", utils.GetMethodName())
	}

	return nil
}

func (r *ContainerRepository) DecrementTotalFiles(containerUUID uuid.UUID) error {
	query := "UPDATE storage.containers SET total_files = total_files - 1 WHERE uuid = $1"
	_, err := r.db.Exec(query, containerUUID)
	if err != nil {
		return utils.FormatError(err, "update", utils.GetMethodName())
	}

	return nil
}

func (r *ContainerRepository) Delete(containerUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.containers WHERE uuid = $1"
	res, err := r.db.Exec(query, containerUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
