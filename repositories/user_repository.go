package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(injector *do.Injector) (*UserRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) List(paginationParams utils.PaginationParams) ([]models.User, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := fmt.Sprintf(
		"SELECT %s FROM authentication.users ORDER BY :sort DESC LIMIT :limit OFFSET :offset",
		utils.GetColumns[models.User](),
	)

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.StructScan(&user); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return users, nil
}

func (r *UserRepository) GetByID(userUUID uuid.UUID) (models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM authentication.users WHERE uuid = $1", utils.GetColumns[models.User]())
	var user models.User
	err := r.db.Get(&user, query, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errs.NewNotFoundError("user.error.notFound")
		}

		return models.User{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return user, nil
}

func (r *UserRepository) ExistsByID(userUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM authentication.users WHERE uuid = $1)"
	var exists bool
	err := r.db.Get(&exists, query, userUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM authentication.users WHERE email = $1)"
	var exists bool
	err := r.db.Get(&exists, query, email)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM authentication.users WHERE username = $1)"
	var exists bool
	err := r.db.Get(&exists, query, username)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM authentication.users WHERE email = $1", utils.GetColumns[models.User]())
	var user models.User
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errs.NewNotFoundError("user.error.notFound")
		}

		return models.User{}, fmt.Errorf("could not fetch rowx: %v", err)
	}

	return user, nil
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	query := "INSERT INTO authentication.users (username, email, status, role_id, password) VALUES ($1, $2, $3, $4, $5) RETURNING uuid"
	err := r.db.QueryRowx(query, user.Username, user.Email, models.UserStatusActive, user.RoleID, utils.HashPassword(user.Password)).Scan(&user.Uuid)
	if err != nil {
		return &models.User{}, fmt.Errorf("could not create row: %v", err)
	}

	return user, nil
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
			return 0, errs.NewUnauthorizedError("auth.error.tokenExpired")
		}

		return 0, fmt.Errorf("could not fetch JWT version: %v", err)
	}

	return version, nil
}

func (r *UserRepository) Update(userUUID uuid.UUID, user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()
	user.Uuid = userUUID

	query := `
		UPDATE authentication.users 
		SET bio = :bio, updated_at = :updated_at 
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, user)
	if err != nil {
		return &models.User{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.User{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return user, nil
}

func (r *UserRepository) Delete(userUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM authentication.users WHERE uuid = $1"
	res, err := r.db.Exec(query, userUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
