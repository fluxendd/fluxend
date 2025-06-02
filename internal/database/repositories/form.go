package repositories

import (
	"fluxend/internal/domain/form"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type FormRepository struct {
	db shared.DB
}

func NewFormRepository(injector *do.Injector) (form.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &FormRepository{db: db}, nil
}

func (r *FormRepository) ListForProject(paginationParams shared.PaginationParams, projectUUID uuid.UUID) ([]form.Form, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			fluxend.forms WHERE project_uuid = :project_uuid
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

	var forms []form.Form
	return forms, r.db.SelectNamedList(&forms, query, params)
}

func (r *FormRepository) GetProjectUUIDByFormUUID(formUUID uuid.UUID) (uuid.UUID, error) {
	query := "SELECT project_uuid FROM fluxend.forms WHERE uuid = $1"

	var projectUUID uuid.UUID
	return projectUUID, r.db.GetWithNotFound(&projectUUID, "form.error.notFound", query, formUUID)
}

func (r *FormRepository) GetByUUID(formUUID uuid.UUID) (form.Form, error) {
	query := "SELECT %s FROM fluxend.forms WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[form.Form]())

	var fetchedForm form.Form
	return fetchedForm, r.db.GetWithNotFound(&fetchedForm, "form.error.notFound", query, formUUID)
}

func (r *FormRepository) ExistsByUUID(formUUID uuid.UUID) (bool, error) {
	return r.db.Exists("fluxend.forms", "uuid = $1", formUUID)
}

func (r *FormRepository) ExistsByNameForProject(name string, projectUUID uuid.UUID) (bool, error) {
	return r.db.Exists("fluxend.forms", "name = $1 AND project_uuid = $2", name, projectUUID)
}

func (r *FormRepository) Create(form *form.Form) (*form.Form, error) {
	return form, r.db.WithTransaction(func(tx shared.Tx) error {
		query := `
        INSERT INTO fluxend.forms (
            project_uuid, name, description, created_by, updated_by
        ) VALUES (
            $1, $2, $3, $4, $5
        )
        RETURNING uuid
        `

		return tx.QueryRowx(
			query,
			form.ProjectUuid, form.Name, form.Description, form.CreatedBy, form.UpdatedBy,
		).Scan(&form.Uuid)
	})
}

func (r *FormRepository) Update(formInput *form.Form) (*form.Form, error) {
	query := `
		UPDATE fluxend.forms 
		SET name = :name, description = :description, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	return formInput, r.db.ExecWithErr(query, formInput)
}

func (r *FormRepository) Delete(projectUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM fluxend.forms WHERE uuid = $1", projectUUID)
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
