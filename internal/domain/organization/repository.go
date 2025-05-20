package organization

import (
	"fluxton/internal/domain/shared"
	"fluxton/internal/domain/user"
	"github.com/google/uuid"
)

type Repository interface {
	ListForUser(paginationParams shared.PaginationParams, authUserID uuid.UUID) ([]Organization, error)
	ListUsers(organizationUUID uuid.UUID) ([]user.User, error)
	GetUser(organizationUUID, userUUID uuid.UUID) (user.User, error)
	CreateUser(organizationUUID, userUUID uuid.UUID) error
	DeleteUser(organizationUUID, userUUID uuid.UUID) error
	GetByUUID(organizationUUID uuid.UUID) (Organization, error)
	ExistsByID(organizationUUID uuid.UUID) (bool, error)
	Create(organization *Organization, authUserID uuid.UUID) (*Organization, error)
	Update(organization *Organization) (*Organization, error)
	Delete(organizationUUID uuid.UUID) (bool, error)
	IsOrganizationMember(organizationUUID, authUserID uuid.UUID) (bool, error)
}
