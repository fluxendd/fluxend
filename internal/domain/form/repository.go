package form

import (
	"fluxton/internal/api/dto"
	"github.com/google/uuid"
)

type Repository interface {
	ListForProject(paginationParams dto.PaginationParams, projectUUID uuid.UUID) ([]Form, error)
	GetProjectUUIDByFormUUID(formUUID uuid.UUID) (uuid.UUID, error)
	GetByUUID(formUUID uuid.UUID) (Form, error)
	ExistsByUUID(formUUID uuid.UUID) (bool, error)
	ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error)
	Create(form *Form) (*Form, error)
	Update(formInput *Form) (*Form, error)
	Delete(projectUUID uuid.UUID) (bool, error)
}
