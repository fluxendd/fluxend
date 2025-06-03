package repositories

import (
	"fluxend/internal/domain/form"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/do"
)

type FormFieldRepository struct {
	db shared.DB
}

func NewFormFieldRepository(injector *do.Injector) (form.FieldRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &FormFieldRepository{db: db}, nil
}

func (r *FormFieldRepository) ListForForm(formUUID uuid.UUID) ([]form.Field, error) {
	query := "SELECT * FROM fluxend.form_fields WHERE form_uuid = $1;"

	var forms []form.Field
	return forms, r.db.Select(&forms, query, formUUID)
}

func (r *FormFieldRepository) GetByUUID(formUUID uuid.UUID) (form.Field, error) {
	query := "SELECT %s FROM fluxend.form_fields WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[form.Field]())

	var fetchedField form.Field
	return fetchedField, r.db.GetWithNotFound(&fetchedField, "form.error.notFound", query, formUUID)
}

func (r *FormFieldRepository) ExistsByUUID(formFieldUUID uuid.UUID) (bool, error) {
	return r.db.Exists("fluxend.form_fields", "uuid = $1", formFieldUUID)
}

func (r *FormFieldRepository) ExistsByAnyLabelForForm(labels []string, formUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxend.form_fields WHERE label = ANY($1) AND form_uuid = $2)"

	var exists bool
	return exists, r.db.Get(&exists, query, pq.Array(labels), formUUID)
}

func (r *FormFieldRepository) ExistsByLabelForForm(label string, formUUID uuid.UUID) (bool, error) {
	return r.db.Exists("fluxend.form_fields", "label = $1 AND form_uuid = $2", label, formUUID)
}

func (r *FormFieldRepository) Create(formField *form.Field) (*form.Field, error) {
	return formField, r.db.WithTransaction(func(tx shared.Tx) error {
		query := `
        INSERT INTO fluxend.form_fields (
            form_uuid,
            label,
            type,
            description,
            is_required,
            options,
            min_length,
            max_length,
            min_value,
            max_value,
            pattern,
            default_value,
            start_date,
            end_date,
            date_format
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
        )
        RETURNING uuid
        `

		return tx.QueryRowx(
			query,
			formField.FormUuid,
			formField.Label,
			formField.Type,
			formField.Description,
			formField.IsRequired,
			formField.Options,
			formField.MinLength,
			formField.MaxLength,
			formField.MinValue,
			formField.MaxValue,
			formField.Pattern,
			formField.DefaultValue,
			formField.StartDate,
			formField.EndDate,
			formField.DateFormat,
		).Scan(&formField.Uuid)
	})
}

func (r *FormFieldRepository) CreateMany(formFields []form.Field, formUUID uuid.UUID) ([]form.Field, error) {
	createdFields := make([]form.Field, 0, len(formFields))
	for i, formField := range formFields {
		formField.FormUuid = formUUID

		createdField, err := r.Create(&formField)
		if err != nil {
			return nil, fmt.Errorf("could not create form field at index %d: %v", i, err)
		}

		createdFields = append(createdFields, *createdField)
	}

	return createdFields, nil
}

func (r *FormFieldRepository) Update(formField *form.Field) (*form.Field, error) {
	query := `
		UPDATE fluxend.form_fields 
		SET 
		    label = :label, 
		    description = :description, 
		    type = :type, 
		    is_required = :is_required, 
		    options = :options, 
		    updated_at = :updated_at
		WHERE uuid = :uuid`

	err := r.db.ExecWithErr(query, formField)

	return formField, err
}

func (r *FormFieldRepository) Delete(formFieldUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM fluxend.form_fields WHERE uuid = $1", formFieldUUID)
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
