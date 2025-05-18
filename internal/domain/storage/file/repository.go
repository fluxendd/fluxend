package file

import (
	"fluxton/internal/api/dto"
	"github.com/google/uuid"
)

type Repository interface {
	ListForContainer(paginationParams dto.PaginationParams, containerUUID uuid.UUID) ([]File, error)
	GetByUUID(fileUUID uuid.UUID) (File, error)
	ExistsByUUID(containerUUID uuid.UUID) (bool, error)
	ExistsByNameForContainer(name string, containerUUID uuid.UUID) (bool, error)
	Create(file *File) (*File, error)
	Rename(container *File) (*File, error)
	Delete(fileUUID uuid.UUID) (bool, error)
}
