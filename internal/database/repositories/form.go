package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/internal/domain/form"
	"fluxton/requests"
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

func (r *FormRepository) ListForProject(paginationParams requests.PaginationParams, projectUUID uuid.UUID) ([]form.Form, error) {
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

	query = fmt.Sprintf(query, utils.GetColumns[form.Form]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var forms []form.Form
	for rows.Next() {
		var form form.Form
		if err := rows.StructScan(&form); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return forms, nil
}

func (r *FormRepository) GetProjectUUIDByFormUUID(formUUID uuid.UUID) (uuid.UUID, error) {
	query := "SELECT project_uuid FROM fluxton.forms WHERE uuid = $1"

	var projectUUID uuid.UUID
	err := r.db.Get(&projectUUID, query, formUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, errs.NewNotFoundError("form.error.notFound")
		}

		return uuid.UUID{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return projectUUID, nil
}

func (r *FormRepository) GetByUUID(formUUID uuid.UUID) (form.Form, error) {
	query := "SELECT %s FROM fluxton.forms WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[form.Form]())

	var form form.Form
	err := r.db.Get(&form, query, formUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return form.Form{}, errs.NewNotFoundError("form.error.notFound")
		}

		return form.Form{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return form, nil
}

func (r *FormRepository) ExistsByUUID(formUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.forms WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, formUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FormRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.forms WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, projectUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FormRepository) Create(form *form.Form) (*form.Form, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `
    INSERT INTO fluxton.forms (
        project_uuid, name, description, created_by, updated_by
    ) VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		form.ProjectUuid, form.Name, form.Description, form.CreatedBy, form.UpdatedBy,
	).Scan(&form.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return form, nil
}

func (r *FormRepository) Update(form *form.Form) (*form.Form, error) {
	query := `
		UPDATE fluxton.forms 
		SET name = :name, description = :description, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, form)
	if err != nil {
		return &form.Form{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &form.Form{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return form, nil
}

func (r *FormRepository) Delete(projectUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.forms WHERE uuid = $1"
	res, err := r.db.Exec(query, projectUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
