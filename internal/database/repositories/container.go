package repositories

import (
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/storage/container"
	"fluxend/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type ContainerRepository struct {
	db shared.DB
}

func NewContainerRepository(injector *do.Injector) (container.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &ContainerRepository{db: db}, nil
}

func (r *ContainerRepository) ListForProject(paginationParams shared.PaginationParams, projectUUID uuid.UUID) ([]container.Container, error) {
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

	query = fmt.Sprintf(query, pkg.GetColumns[container.Container]()) // dynamically pulls columns from entity

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	var containers []container.Container
	return containers, r.db.SelectNamedList(&containers, query, params)
}

func (r *ContainerRepository) GetByUUID(containerUUID uuid.UUID) (container.Container, error) {
	query := "SELECT %s FROM storage.containers WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[container.Container]())

	var fetchedContainer container.Container
	return fetchedContainer, r.db.GetWithNotFound(&fetchedContainer, "container.error.notFound", query, containerUUID)
}

func (r *ContainerRepository) ExistsByUUID(containerUUID uuid.UUID) (bool, error) {
	return r.db.Exists("storage.containers", "uuid = $1", containerUUID)
}

func (r *ContainerRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	return r.db.Exists("storage.containers", "name = $1 AND project_uuid = $2", name, projectUUID)
}

func (r *ContainerRepository) Create(container *container.Container) (*container.Container, error) {
	return container, r.db.WithTransaction(func(tx shared.Tx) error {
		query := `
        INSERT INTO storage.containers (
            project_uuid, name, name_key, provider, description, is_public, url, max_file_size, created_by, updated_by
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
        RETURNING uuid
        `

		return tx.QueryRowx(
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
	})
}

func (r *ContainerRepository) Update(containerInput *container.Container) (*container.Container, error) {
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

	_, err := r.db.NamedExecWithRowsAffected(query, containerInput)
	return containerInput, err
}

func (r *ContainerRepository) IncrementTotalFiles(containerUUID uuid.UUID) error {
	_, err := r.db.ExecWithRowsAffected("UPDATE storage.containers SET total_files = total_files + 1 WHERE uuid = $1", containerUUID)
	return err
}

func (r *ContainerRepository) DecrementTotalFiles(containerUUID uuid.UUID) error {
	_, err := r.db.ExecWithRowsAffected("UPDATE storage.containers SET total_files = total_files - 1 WHERE uuid = $1", containerUUID)
	return err
}

func (r *ContainerRepository) Delete(containerUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM storage.containers WHERE uuid = $1", containerUUID)
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
