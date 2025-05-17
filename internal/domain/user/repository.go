package user

import (
	"fluxton/internal/api/dto"
	"github.com/google/uuid"
	"os/user"
)

type Repository interface {
	List(paginationParams dto.PaginationParams) ([]user.User, error)
	GetByID(userUUID uuid.UUID) (user.User, error)
	ExistsByID(userUUID uuid.UUID) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	GetByEmail(email string) (user.User, error)
	Create(user *user.User) (*user.User, error)
	CreateJWTVersion(userId uuid.UUID) (int, error)
	GetJWTVersion(userId uuid.UUID) (int, error)
	Update(userUUID uuid.UUID, user *user.User) (*user.User, error)
	Delete(userUUID uuid.UUID) (bool, error)
}
