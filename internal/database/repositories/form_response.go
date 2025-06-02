package repositories

import (
	"fluxend/internal/domain/form"
	"fluxend/internal/domain/shared"
	flxErrs "fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type FormResponseRepository struct {
	db shared.DB
}

func NewFormResponseRepository(injector *do.Injector) (form.FieldResponseRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &FormResponseRepository{db: db}, nil
}

func (r *FormResponseRepository) ListForForm(formUUID uuid.UUID) ([]form.FormResponse, error) {
	query := `
		SELECT 
			fr.uuid AS form_response_uuid, 
			fr.form_uuid, 
			fr.created_at AS form_response_created_at, 
			fr.updated_at AS form_response_updated_at, 
			ffr.uuid AS form_field_response_uuid, 
			ffr.form_response_uuid, 
			ffr.form_field_uuid, 
			ffr.value, 
			ffr.created_at AS form_field_response_created_at, 
			ffr.updated_at AS form_field_response_updated_at
		FROM fluxend.form_responses fr
		LEFT JOIN fluxend.form_field_responses ffr 
			ON fr.uuid = ffr.form_response_uuid
		WHERE fr.form_uuid = $1
		ORDER BY fr.created_at, ffr.created_at;
	`

	rows, err := r.db.Query(query, formUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanFormResponseRows(rows)
}

func (r *FormResponseRepository) GetByUUID(formResponseUUID uuid.UUID) (*form.FormResponse, error) {
	query := `
		SELECT 
			fr.uuid AS form_response_uuid, 
			fr.form_uuid, 
			fr.created_at AS form_response_created_at, 
			fr.updated_at AS form_response_updated_at, 
			ffr.uuid AS form_field_response_uuid, 
			ffr.form_response_uuid, 
			ffr.form_field_uuid, 
			ffr.value, 
			ffr.created_at AS form_field_response_created_at, 
			ffr.updated_at AS form_field_response_updated_at
		FROM fluxend.form_responses fr
		LEFT JOIN fluxend.form_field_responses ffr 
			ON fr.uuid = ffr.form_response_uuid
		WHERE fr.uuid = $1
		ORDER BY fr.created_at, ffr.created_at;
	`

	rows, err := r.db.Query(query, formResponseUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formResponses, err := r.scanFormResponseRows(rows)
	if err != nil {
		return nil, err
	}

	if len(formResponses) == 0 {
		return nil, flxErrs.NewNotFoundError("formResponse.error.notFound")
	}

	return &formResponses[0], nil
}

func (r *FormResponseRepository) Create(
	formResponse *form.FormResponse,
	formFieldResponse *[]form.FieldResponse,
) (*form.FormResponse, error) {
	return formResponse, r.db.WithTransaction(func(tx shared.Tx) error {
		// Insert form response
		query := `INSERT INTO fluxend.form_responses (form_uuid) VALUES ($1) RETURNING uuid`
		if err := tx.QueryRowx(query, formResponse.FormUuid).Scan(&formResponse.Uuid); err != nil {
			return err
		}

		// Insert form field responses
		for i, ffr := range *formFieldResponse {
			query = `INSERT INTO fluxend.form_field_responses (form_response_uuid, form_field_uuid, value) VALUES ($1, $2, $3) RETURNING uuid`
			if err := tx.QueryRowx(query, formResponse.Uuid, ffr.FormFieldUuid, ffr.Value).Scan(&(*formFieldResponse)[i].Uuid); err != nil {
				return err
			}
		}

		formResponse.Responses = *formFieldResponse
		return nil
	})
}

func (r *FormResponseRepository) Delete(formResponseUUID uuid.UUID) error {
	return r.db.WithTransaction(func(tx shared.Tx) error {
		// Delete form field responses first (foreign key constraint)
		if _, err := tx.Exec(`DELETE FROM fluxend.form_field_responses WHERE form_response_uuid = $1`, formResponseUUID); err != nil {
			return err
		}

		// Delete form response
		if _, err := tx.Exec(`DELETE FROM fluxend.form_responses WHERE uuid = $1`, formResponseUUID); err != nil {
			return err
		}

		return nil
	})
}

// Helper method to scan form response rows with complex join logic
func (r *FormResponseRepository) scanFormResponseRows(rows interface{}) ([]form.FormResponse, error) {
	formResponseMap := make(map[uuid.UUID]*form.FormResponse)

	formResponses := make([]form.FormResponse, 0, len(formResponseMap))
	for _, formResponse := range formResponseMap {
		formResponses = append(formResponses, *formResponse)
	}

	return formResponses, nil
}
