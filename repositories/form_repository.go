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
)

type FormRepository struct {
	db *sqlx.DB
}

func NewFormRepository(injector *do.Injector) (*FormRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &FormRepository{db: db}, nil
}

func (r *FormRepository) ListForProject(paginationParams utils.PaginationParams, projectUUID uuid.UUID) ([]models.Form, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			fluxton.forms WHERE project_uuid = :project_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, utils.GetColumns[models.Form]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var forms []models.Form
	for rows.Next() {
		var organization models.Form
		if err := rows.StructScan(&organization); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		forms = append(forms, organization)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return forms, nil
}

func (r *FormRepository) GetByUUID(formUUID uuid.UUID) (models.Form, error) {
	query := "SELECT %s FROM fluxton.forms WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Form]())

	var form models.Form
	err := r.db.Get(&form, query, formUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Form{}, errs.NewNotFoundError("form.error.notFound")
		}

		return models.Form{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return form, nil
}

func (r *FormRepository) ExistsByUUID(formUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.forms WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, formUUID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *FormRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.forms WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, projectUUID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *FormRepository) Create(form *models.Form) (*models.Form, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	query := "INSERT INTO fluxton.projects (project_uuid, name, description, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING uuid"
	queryErr := tx.QueryRowx(query, form.ProjectUuid, form.Name, form.Description, form.CreatedBy, form.UpdatedBy).Scan(&form.Uuid)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create project: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return form, nil
}

func (r *FormRepository) Update(form *models.Form) (*models.Form, error) {
	query := `
		UPDATE fluxton.forms 
		SET name = :name, description = :description, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, form)
	if err != nil {
		return &models.Form{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Form{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return form, nil
}

func (r *FormRepository) Delete(projectId uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.forms WHERE uuid = $1"
	res, err := r.db.Exec(query, projectId)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
