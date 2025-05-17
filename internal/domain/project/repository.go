package project

import (
	"fluxton/internal/api/dto"
	"github.com/google/uuid"
)

type Repository interface {
	ListForUser(paginationParams dto.PaginationParams, authUserId uuid.UUID) ([]Project, error)
	List(paginationParams dto.PaginationParams) ([]Project, error)
	GetByUUID(projectUUID uuid.UUID) (Project, error)
	GetDatabaseNameByUUID(projectUUID uuid.UUID) (string, error)
	GetUUIDByDatabaseName(dbName string) (uuid.UUID, error)
	GetOrganizationUUIDByProjectUUID(id uuid.UUID) (uuid.UUID, error)
	ExistsByUUID(id uuid.UUID) (bool, error)
	ExistsByNameForOrganization(name string, organizationUUID uuid.UUID) (bool, error)
	Create(project *Project) (*Project, error)
	Update(project *Project) (*Project, error)
	UpdateStatusByDatabaseName(databaseName, status string) (bool, error)
	Delete(projectUUID uuid.UUID) (bool, error)
}
