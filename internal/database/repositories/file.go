package repositories

import (
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/storage/file"
	"fluxend/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type FileRepository struct {
	db shared.DB
}

func NewFileRepository(injector *do.Injector) (file.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &FileRepository{db: db}, nil
}

func (r *FileRepository) ListForContainer(paginationParams shared.PaginationParams, containerUUID uuid.UUID) ([]file.File, error) {
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

	query = fmt.Sprintf(query, pkg.GetColumns[file.File]())

	params := map[string]interface{}{
		"container_uuid": containerUUID,
		"sort":           paginationParams.Sort,
		"limit":          paginationParams.Limit,
		"offset":         offset,
	}

	var files []file.File
	return files, r.db.SelectNamedList(&files, query, params)
}

func (r *FileRepository) GetByUUID(fileUUID uuid.UUID) (file.File, error) {
	query := "SELECT %s FROM storage.files WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[file.File]())

	var fetchedFile file.File
	return fetchedFile, r.db.GetWithNotFound(&fetchedFile, "file.error.notFound", query, fileUUID)
}

func (r *FileRepository) ExistsByUUID(containerUUID uuid.UUID) (bool, error) {
	return r.db.Exists("storage.files", "uuid = $1", containerUUID)
}

func (r *FileRepository) ExistsByNameForContainer(name string, containerUUID uuid.UUID) (bool, error) {
	return r.db.Exists("storage.files", "full_file_name = $1 AND container_uuid = $2", name, containerUUID)
}

func (r *FileRepository) Create(file *file.File) (*file.File, error) {
	return file, r.db.WithTransaction(func(tx shared.Tx) error {
		query := `
        INSERT INTO storage.files (
            container_uuid, full_file_name, size, mime_type, created_by, updated_by, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
        RETURNING uuid
        `

		return tx.QueryRowx(
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
	})
}

func (r *FileRepository) Rename(container *file.File) (*file.File, error) {
	query := `
		UPDATE storage.files 
		SET full_file_name = :full_file_name, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	_, err := r.db.NamedExecWithRowsAffected(query, container)
	return container, err
}

func (r *FileRepository) Delete(fileUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM storage.files WHERE uuid = $1", fileUUID)
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
