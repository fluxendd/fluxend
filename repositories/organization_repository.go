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

type OrganizationRepository struct {
	db *sqlx.DB
}

func NewOrganizationRepository(injector *do.Injector) (*OrganizationRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
	return &OrganizationRepository{db: db}, nil
}

func (r *OrganizationRepository) List(paginationParams utils.PaginationParams) ([]models.Organization, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := `
		SELECT 
			%s 
		FROM organizations 
			ORDER BY :sort DESC
		LIMIT :limit 
		OFFSET :offset
	`

	query = fmt.Sprintf(query, models.Organization{}.GetFields())

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

	var organizations []models.Organization
	for rows.Next() {
		var organization models.Organization
		if err := rows.StructScan(&organization); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		organizations = append(organizations, organization)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return organizations, nil
}

func (r *OrganizationRepository) GetByID(id uint) (models.Organization, error) {
	query := "SELECT %s FROM organizations WHERE id = $1"
	query = fmt.Sprintf(query, models.Organization{}.GetFields())

	var organization models.Organization
	err := r.db.Get(&organization, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Organization{}, errs.NewNotFoundError("organization.error.notFound")
		}

		return models.Organization{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return organization, nil
}

func (r *OrganizationRepository) ExistsByID(id uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM organizations WHERE id = $1)"
	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *OrganizationRepository) Create(organization *models.Organization) (*models.Organization, error) {
	query := "INSERT INTO organizations (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRowx(query, organization.Name).Scan(&organization.ID)
	if err != nil {
		return &models.Organization{}, fmt.Errorf("could not create row: %v", err)
	}

	return organization, nil
}

func (r *OrganizationRepository) Update(id uint, organization *models.Organization) (*models.Organization, error) {
	organization.UpdatedAt = time.Now()
	organization.ID = id

	query := `
		UPDATE organizations 
		SET name = :name, updated_at = :updated_at 
		WHERE id = :id`

	res, err := r.db.NamedExec(query, organization)
	if err != nil {
		return &models.Organization{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Organization{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return organization, nil
}

func (r *OrganizationRepository) Delete(organizationId uint) (bool, error) {
	query := "DELETE FROM organizations WHERE id = $1"
	res, err := r.db.Exec(query, organizationId)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
