package repositories

import (
	"database/sql"
	"errors"
	"fluxton/internal/api/dto"
	"fluxton/internal/domain/storage/container"
	"fluxton/pkg"
	flxErrs "fluxton/pkg/errors"
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

func (r *ContainerRepository) ListForProject(paginationParams dto.PaginationParams, projectUUID uuid.UUID) ([]container.Container, error) {
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

	query = fmt.Sprintf(query, pkg.GetColumns[container.Container]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var containers []container.Container
	for rows.Next() {
		var container container.Container
		if err := rows.StructScan(&container); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		containers = append(containers, container)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return containers, nil
}

func (r *ContainerRepository) GetByUUID(containerUUID uuid.UUID) (container.Container, error) {
	query := "SELECT %s FROM storage.containers WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[container.Container]())

	var container container.Container
	err := r.db.Get(&container, query, containerUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return container.Container{}, flxErrs.NewNotFoundError("container.error.notFound")
		}

		return container.Container{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return container, nil
}

func (r *ContainerRepository) ExistsByUUID(containerUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.containers WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, containerUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *ContainerRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM storage.containers WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, projectUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *ContainerRepository) Create(container *container.Container) (*container.Container, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	query := `
    INSERT INTO storage.containers (
        project_uuid, name, name_key, provider, description, is_public, url, max_file_size, created_by, updated_by
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		container.ProjectUuid,
		container.Name,
		container.NameKey,
		container.Provider,
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
		return nil, pkg.FormatError(err, "create", pkg.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return container, nil
}

func (r *ContainerRepository) Update(container *container.Container) (*container.Container, error) {
	query := `
		UPDATE storage.containers 
		SET 
		    name = :name, 
		    description = :description, 
		    is_public = :is_public, 
		    max_file_size = :max_file_size,
		    updated_at = :updated_at, 
		    updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, container)
	if err != nil {
		return &container.Container{}, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &container.Container{}, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return container, nil
}

func (r *ContainerRepository) IncrementTotalFiles(containerUUID uuid.UUID) error {
	query := "UPDATE storage.containers SET total_files = total_files + 1 WHERE uuid = $1"
	_, err := r.db.Exec(query, containerUUID)
	if err != nil {
		return pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	return nil
}

func (r *ContainerRepository) DecrementTotalFiles(containerUUID uuid.UUID) error {
	query := "UPDATE storage.containers SET total_files = total_files - 1 WHERE uuid = $1"
	_, err := r.db.Exec(query, containerUUID)
	if err != nil {
		return pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	return nil
}

func (r *ContainerRepository) Delete(containerUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM storage.containers WHERE uuid = $1"
	res, err := r.db.Exec(query, containerUUID)
	if err != nil {
		return false, pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}
