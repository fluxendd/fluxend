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

type OrganizationRepository struct {
	db *sqlx.DB
}

func NewOrganizationRepository(injector *do.Injector) (*OrganizationRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
	return &OrganizationRepository{db: db}, nil
}

func (r *OrganizationRepository) ListForUser(paginationParams utils.PaginationParams, authUserID uuid.UUID) ([]models.Organization, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	modelSkeleton := models.Organization{}

	query := `
		SELECT 
			%s 
		FROM 
			fluxton.organizations organizations
		JOIN 
			fluxton.organization_users organization_users ON organizations.id = organization_users.organization_id
		WHERE 
			organization_users.user_id = :user_id
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, modelSkeleton.GetColumnsWithAlias(modelSkeleton.GetTableName()))

	params := map[string]interface{}{
		"user_id": authUserID,
		"sort":    paginationParams.Sort,
		"limit":   paginationParams.Limit,
		"offset":  offset,
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

func (r *OrganizationRepository) GetByIDForUser(id, authUserID uuid.UUID) (models.Organization, error) {
	query := "SELECT %s FROM fluxton.organizations WHERE id = $1"
	query = fmt.Sprintf(query, models.Organization{}.GetColumns())

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
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.organizations WHERE id = $1)"
	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *OrganizationRepository) Create(organization *models.Organization, authUserID uuid.UUID) (*models.Organization, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	// Insert into organizations table
	query := "INSERT INTO fluxton.organizations (name) VALUES ($1) RETURNING id"
	queryErr := tx.QueryRowx(query, organization.Name).Scan(&organization.ID)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create organization: %v", queryErr)
	}

	// Insert into organization_users pivot table
	queryErr = r.createOrganizationUser(tx, organization.ID, authUserID)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not insert into pivot table: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	organization.ID = organization.ID
	return organization, nil
}

func (r *OrganizationRepository) Update(id uuid.UUID, organization *models.Organization) (*models.Organization, error) {
	organization.UpdatedAt = time.Now()
	organization.ID = id

	query := `
		UPDATE fluxton.organizations 
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

func (r *OrganizationRepository) Delete(organizationID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.organizations WHERE id = $1"
	res, err := r.db.Exec(query, organizationID)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}

func (r *OrganizationRepository) IsOrganizationUser(organizationID, authUserID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.organization_users WHERE organization_id = $1 AND user_id = $2)"

	var exists bool
	err := r.db.Get(&exists, query, organizationID, authUserID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *OrganizationRepository) createOrganizationUser(tx *sqlx.Tx, organizationID, userId uuid.UUID) error {
	query := "INSERT INTO fluxton.organization_users (organization_id, user_id) VALUES ($1, $2)"
	_, err := tx.Exec(query, organizationID, userId)
	if err != nil {
		return fmt.Errorf("could not insert into pivot table: %v", err)
	}

	return nil
}
