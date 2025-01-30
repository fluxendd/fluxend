package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"myapp/errs"
	"myapp/models"
	"myapp/utils"
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
		"SELECT %s FROM users ORDER BY :sort DESC LIMIT :limit OFFSET :offset",
		models.User{}.GetFields(),
	)

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return users, nil
}

func (r *UserRepository) GetByID(id uint) (models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM users WHERE id = $1", models.User{}.GetFields())
	var user models.User
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errs.NewNotFoundError("user.error.notFound")
		}

		return models.User{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return user, nil
}

func (r *UserRepository) ExistsByID(id uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"
	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM users WHERE email = $1", models.User{}.GetFields())
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
	query := "INSERT INTO users (username, email, status, role_id, password) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRowx(query, user.Username, user.Email, models.UserStatusActive, user.RoleID, utils.HashPassword(user.Password)).Scan(&user.ID)
	if err != nil {
		return &models.User{}, fmt.Errorf("could not create row: %v", err)
	}

	return user, nil
}

func (r *UserRepository) Update(id uint, user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()
	user.ID = id

	query := `
		UPDATE users 
		SET bio = :bio, updated_at = :updated_at 
		WHERE id = :id`

	res, err := r.db.NamedExec(query, user)
	if err != nil {
		return &models.User{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.User{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return user, nil
}

func (r *UserRepository) Delete(userId uint) (bool, error) {
	query := "DELETE FROM users WHERE id = $1"
	res, err := r.db.Exec(query, userId)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
