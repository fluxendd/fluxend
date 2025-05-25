package repositories

import (
	"fluxend/internal/domain/form"
	"fluxend/pkg"
	flxErrs "fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type FormResponseRepository struct {
	db *sqlx.DB
}

func NewFormResponseRepository(injector *do.Injector) (form.FieldResponseRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

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

	rows, err := r.db.Queryx(query, formUUID)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	formResponseMap := make(map[uuid.UUID]*form.FormResponse)

	for rows.Next() {
		var row struct {
			FormResponseUUID           uuid.UUID  `db:"form_response_uuid"`
			FormUUID                   uuid.UUID  `db:"form_uuid"`
			FormResponseCreatedAt      time.Time  `db:"form_response_created_at"`
			FormResponseUpdatedAt      time.Time  `db:"form_response_updated_at"`
			FormFieldResponseUUID      *uuid.UUID `db:"form_field_response_uuid"`
			FormFieldUUID              *uuid.UUID `db:"form_field_uuid"`
			Value                      *string    `db:"value"`
			FormFieldResponseCreatedAt *time.Time `db:"form_field_response_created_at"`
			FormFieldResponseUpdatedAt *time.Time `db:"form_field_response_updated_at"`
		}

		if err := rows.StructScan(&row); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}

		formResponse, exists := formResponseMap[row.FormResponseUUID]
		if !exists {
			formResponse = &form.FormResponse{
				Uuid:      row.FormResponseUUID,
				FormUuid:  row.FormUUID,
				CreatedAt: row.FormResponseCreatedAt,
				UpdatedAt: row.FormResponseUpdatedAt,
			}
			formResponseMap[row.FormResponseUUID] = formResponse
		}

		if row.FormFieldResponseUUID != nil {
			formResponse.Responses = append(formResponse.Responses, form.FieldResponse{
				Uuid:             *row.FormFieldResponseUUID,
				FormResponseUuid: row.FormResponseUUID,
				FormFieldUuid:    *row.FormFieldUUID,
				Value:            *row.Value,
				CreatedAt:        *row.FormFieldResponseCreatedAt,
				UpdatedAt:        *row.FormFieldResponseUpdatedAt,
			})
		}

		formResponseMap[row.FormResponseUUID] = formResponse
	}

	formResponses := make([]form.FormResponse, 0, len(formResponseMap))
	for _, formResponse := range formResponseMap {
		formResponses = append(formResponses, *formResponse)
	}

	return formResponses, nil
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

	rows, err := r.db.Queryx(query, formResponseUUID)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	formResponseMap := make(map[uuid.UUID]*form.FormResponse)

	for rows.Next() {
		var row struct {
			FormResponseUUID           uuid.UUID  `db:"form_response_uuid"`
			FormUUID                   uuid.UUID  `db:"form_uuid"`
			FormResponseCreatedAt      time.Time  `db:"form_response_created_at"`
			FormResponseUpdatedAt      time.Time  `db:"form_response_updated_at"`
			FormFieldResponseUUID      *uuid.UUID `db:"form_field_response_uuid"`
			FormFieldUUID              *uuid.UUID `db:"form_field_uuid"`
			Value                      *string    `db:"value"`
			FormFieldResponseCreatedAt *time.Time `db:"form_field_response_created_at"`
			FormFieldResponseUpdatedAt *time.Time `db:"form_field_response_updated_at"`
		}

		if err := rows.StructScan(&row); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}

		formResponse, exists := formResponseMap[row.FormResponseUUID]
		if !exists {
			formResponse = &form.FormResponse{
				Uuid:      row.FormResponseUUID,
				FormUuid:  row.FormUUID,
				CreatedAt: row.FormResponseCreatedAt,
				UpdatedAt: row.FormResponseUpdatedAt,
			}
			formResponseMap[row.FormResponseUUID] = formResponse
		}

		if row.FormFieldResponseUUID != nil {
			formResponse.Responses = append(formResponse.Responses, form.FieldResponse{
				Uuid:             *row.FormFieldResponseUUID,
				FormResponseUuid: row.FormResponseUUID,
				FormFieldUuid:    *row.FormFieldUUID,
				Value:            *row.Value,
				CreatedAt:        *row.FormFieldResponseCreatedAt,
				UpdatedAt:        *row.FormFieldResponseUpdatedAt,
			})
		}

		break
	}

	if len(formResponseMap) == 0 {
		return nil, flxErrs.NewNotFoundError("formResponse.error.notFound")
	}

	formResponse := formResponseMap[formResponseUUID]

	return formResponse, nil
}

func (r *FormResponseRepository) Create(
	formResponse *form.FormResponse,
	formFieldResponse *[]form.FieldResponse,
) (*form.FormResponse, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	query := `INSERT INTO fluxend.form_responses (form_uuid) VALUES ($1) RETURNING uuid`

	queryErr := tx.QueryRowx(
		query,
		formResponse.FormUuid).Scan(&formResponse.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, pkg.FormatError(queryErr, "insert", pkg.GetMethodName())
	}

	for _, ffr := range *formFieldResponse {
		query = `INSERT INTO fluxend.form_field_responses (form_response_uuid, form_field_uuid, value) VALUES ($1, $2, $3) RETURNING uuid`
		queryErr = tx.QueryRowx(
			query,
			formResponse.Uuid,
			ffr.FormFieldUuid,
			ffr.Value).Scan(&ffr.Uuid)
		if queryErr != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, pkg.FormatError(queryErr, "insert", pkg.GetMethodName())
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	formResponse.Responses = *formFieldResponse

	return formResponse, nil
}

func (r *FormResponseRepository) Delete(formResponseUUID uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	query := `DELETE FROM fluxend.form_field_responses WHERE form_response_uuid = $1`
	_, err = tx.Exec(query, formResponseUUID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	query = `DELETE FROM fluxend.form_responses WHERE uuid = $1`
	_, err = tx.Exec(query, formResponseUUID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return nil
}
