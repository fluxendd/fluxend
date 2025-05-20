package repositories

import (
	"database/sql"
	"errors"
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/shared"
	"fluxton/internal/domain/user"
	"fluxton/pkg"
	"fluxton/pkg/auth"
	flxErrs "fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(injector *do.Injector) (user.Repository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
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

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var currentUser user.User
		if err := rows.StructScan(&currentUser); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		users = append(users, currentUser)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return users, nil
}

func (r *UserRepository) GetByID(userUUID uuid.UUID) (user.User, error) {
	query := fmt.Sprintf("SELECT %s FROM authentication.users WHERE uuid = $1", pkg.GetColumns[user.User]())
	var fetchedUser user.User
	err := r.db.Get(&fetchedUser, query, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, flxErrs.NewNotFoundError("user.error.notFound")
		}

		return user.User{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return fetchedUser, nil
}

func (r *UserRepository) ExistsByID(userUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM authentication.users WHERE uuid = $1)"
	var exists bool
	err := r.db.Get(&exists, query, userUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM authentication.users WHERE email = $1)"
	var exists bool
	err := r.db.Get(&exists, query, email)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM authentication.users WHERE username = $1)"
	var exists bool
	err := r.db.Get(&exists, query, username)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *UserRepository) GetByEmail(email string) (user.User, error) {
	query := fmt.Sprintf("SELECT %s FROM authentication.users WHERE email = $1", pkg.GetColumns[user.User]())
	var fetchedUser user.User
	err := r.db.Get(&fetchedUser, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, flxErrs.NewNotFoundError("user.error.notFound")
		}

		return user.User{}, fmt.Errorf("could not fetch rowx: %v", err)
	}

	return fetchedUser, nil
}

func (r *UserRepository) Create(input *user.User) (*user.User, error) {
	query := "INSERT INTO authentication.users (username, email, status, role_id, password) VALUES ($1, $2, $3, $4, $5) RETURNING uuid"
	err := r.db.QueryRowx(query, input.Username, input.Email, constants.UserStatusActive, input.RoleID, auth.HashPassword(input.Password)).Scan(&input.Uuid)
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
		if errors.Is(err, sql.ErrNoRows) {
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

	res, err := r.db.NamedExec(query, inputUser)
	if err != nil {
		return &user.User{}, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &user.User{}, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return inputUser, nil
}

func (r *UserRepository) Delete(userUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM authentication.users WHERE uuid = $1"
	res, err := r.db.Exec(query, userUUID)
	if err != nil {
		return false, pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}
