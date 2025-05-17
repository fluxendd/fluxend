package repositories

import (
	"database/sql"
	"errors"
	"fluxton/internal/api/dto"
	"fluxton/internal/domain/form"
	"fluxton/pkg"
	flxErrs "fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type FormRepository struct {
	db *sqlx.DB
}

func NewFormRepository(injector *do.Injector) (form.Repository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &FormRepository{db: db}, nil
}

func (r *FormRepository) ListForProject(paginationParams dto.PaginationParams, projectUUID uuid.UUID) ([]form.Form, error) {
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

	query = fmt.Sprintf(query, pkg.GetColumns[form.Form]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var forms []form.Form
	for rows.Next() {
		var currentForm form.Form
		if err := rows.StructScan(&currentForm); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		forms = append(forms, currentForm)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return forms, nil
}

func (r *FormRepository) GetProjectUUIDByFormUUID(formUUID uuid.UUID) (uuid.UUID, error) {
	query := "SELECT project_uuid FROM fluxton.forms WHERE uuid = $1"

	var projectUUID uuid.UUID
	err := r.db.Get(&projectUUID, query, formUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, flxErrs.NewNotFoundError("form.error.notFound")
		}

		return uuid.UUID{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return projectUUID, nil
}

func (r *FormRepository) GetByUUID(formUUID uuid.UUID) (form.Form, error) {
	query := "SELECT %s FROM fluxton.forms WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[form.Form]())

	var fetchedForm form.Form
	err := r.db.Get(&fetchedForm, query, formUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return form.Form{}, flxErrs.NewNotFoundError("form.error.notFound")
		}

		return form.Form{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return fetchedForm, nil
}

func (r *FormRepository) ExistsByUUID(formUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.forms WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, formUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *FormRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.forms WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, projectUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *FormRepository) Create(form *form.Form) (*form.Form, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
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
		return nil, pkg.FormatError(queryErr, "insert", pkg.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return form, nil
}

func (r *FormRepository) Update(formInput *form.Form) (*form.Form, error) {
	query := `
		UPDATE fluxton.forms 
		SET name = :name, description = :description, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, formInput)
	if err != nil {
		return &form.Form{}, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &form.Form{}, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return formInput, nil
}

func (r *FormRepository) Delete(projectUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.forms WHERE uuid = $1"
	res, err := r.db.Exec(query, projectUUID)
	if err != nil {
		return false, pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}
