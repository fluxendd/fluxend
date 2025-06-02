package repositories

import (
	"fluxend/internal/config/constants"
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/user"
	"fluxend/pkg"
	"fluxend/pkg/auth"
	flxErrs "fluxend/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type UserRepository struct {
	db shared.DB
}

func NewUserRepository(injector *do.Injector) (user.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) List(paginationParams shared.PaginationParams) ([]user.User, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := fmt.Sprintf(
		"SELECT %s FROM authentication.users ORDER BY :sort DESC LIMIT :limit OFFSET :offset",
		pkg.GetColumns[user.User](),
	)

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	var users []user.User
	return users, r.db.SelectNamedList(&users, query, params)
}

func (r *UserRepository) GetByID(userUUID uuid.UUID) (user.User, error) {
	query := fmt.Sprintf("SELECT %s FROM authentication.users WHERE uuid = $1", pkg.GetColumns[user.User]())

	var fetchedUser user.User
	return fetchedUser, r.db.GetWithNotFound(&fetchedUser, "user.error.notFound", query, userUUID)
}

func (r *UserRepository) ExistsByID(userUUID uuid.UUID) (bool, error) {
	return r.db.Exists("authentication.users", "uuid = $1", userUUID)
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	return r.db.Exists("authentication.users", "email = $1", email)
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	return r.db.Exists("authentication.users", "username = $1", username)
}

func (r *UserRepository) GetByEmail(email string) (user.User, error) {
	query := fmt.Sprintf("SELECT %s FROM authentication.users WHERE email = $1", pkg.GetColumns[user.User]())

	var fetchedUser user.User
	return fetchedUser, r.db.GetWithNotFound(&fetchedUser, "user.error.notFound", query, email)
}

func (r *UserRepository) Create(input *user.User) (*user.User, error) {
	query := "INSERT INTO authentication.users (username, email, status, role_id, password) VALUES ($1, $2, $3, $4, $5) RETURNING uuid"

	err := r.db.QueryRow(query, input.Username, input.Email, constants.UserStatusActive, input.RoleID, auth.HashPassword(input.Password)).Scan(&input.Uuid)
	if err != nil {
		return &user.User{}, fmt.Errorf("could not create row: %v", err)
	}

	return input, nil
}

func (r *UserRepository) CreateJWTVersion(userId uuid.UUID) (int, error) {
	var version int
	query := `
		INSERT INTO authentication.jwt_versions (user_id, version, updated_at)
		VALUES ($1, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) 
		DO UPDATE 
		SET version = authentication.jwt_versions.version + 1, 
		    updated_at = CURRENT_TIMESTAMP
		RETURNING version;
	`

	err := r.db.QueryRow(query, userId).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("could not create or update JWT version: %v", err)
	}

	return version, nil
}

func (r *UserRepository) GetJWTVersion(userId uuid.UUID) (int, error) {
	query := "SELECT version FROM authentication.jwt_versions WHERE user_id = $1"
	var version int

	err := r.db.Get(&version, query, userId)
	if err != nil {
		// Check if it's a not found error and return custom error
		if err.Error() == "sql: no rows in result set" {
			return 0, flxErrs.NewUnauthorizedError("auth.error.tokenExpired")
		}
		return 0, fmt.Errorf("could not fetch JWT version: %v", err)
	}

	return version, nil
}

func (r *UserRepository) Update(userUUID uuid.UUID, inputUser *user.User) (*user.User, error) {
	inputUser.UpdatedAt = time.Now()
	inputUser.Uuid = userUUID

	query := `
		UPDATE authentication.users 
		SET bio = :bio, updated_at = :updated_at 
		WHERE uuid = :uuid`

	_, err := r.db.NamedExecWithRowsAffected(query, inputUser)
	return inputUser, err
}

func (r *UserRepository) Delete(userUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM authentication.users WHERE uuid = $1", userUUID)
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
