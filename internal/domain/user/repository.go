package user

import (
	"fluxend/internal/domain/shared"
	"github.com/google/uuid"
)

type Repository interface {
	List(paginationParams shared.PaginationParams) ([]User, error)
	GetByID(userUUID uuid.UUID) (User, error)
	ExistsByID(userUUID uuid.UUID) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	GetByEmail(email string) (User, error)
	Create(user *User) (*User, error)
	CreateJWTVersion(userId uuid.UUID) (int, error)
	GetJWTVersion(userId uuid.UUID) (int, error)
	Update(userUUID uuid.UUID, user *User) (*User, error)
	Delete(userUUID uuid.UUID) (bool, error)
}
