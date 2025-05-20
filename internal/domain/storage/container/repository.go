package container

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
)

type Repository interface {
	ListForProject(paginationParams shared.PaginationParams, projectUUID uuid.UUID) ([]Container, error)
	GetByUUID(containerUUID uuid.UUID) (Container, error)
	ExistsByUUID(containerUUID uuid.UUID) (bool, error)
	ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error)
	Create(container *Container) (*Container, error)
	Update(container *Container) (*Container, error)
	IncrementTotalFiles(containerUUID uuid.UUID) error
	DecrementTotalFiles(containerUUID uuid.UUID) error
	Delete(containerUUID uuid.UUID) (bool, error)
}
