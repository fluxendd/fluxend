package repositories

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type FormResponseRepository struct {
	db *sqlx.DB
}

func NewFormResponseRepository(injector *do.Injector) (*FormResponseRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &FormResponseRepository{db: db}, nil
}

func (r *FormResponseRepository) ListForForm(formUUID uuid.UUID) ([]models.FormResponse, error) {
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
		FROM fluxton.form_responses fr
		LEFT JOIN fluxton.form_field_responses ffr 
			ON fr.uuid = ffr.form_response_uuid
		WHERE fr.form_uuid = $1
		ORDER BY fr.created_at, ffr.created_at;
	`

	rows, err := r.db.Queryx(query, formUUID)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	formResponseMap := make(map[uuid.UUID]*models.FormResponse)

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
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}

		formResponse, exists := formResponseMap[row.FormResponseUUID]
		if !exists {
			formResponse = &models.FormResponse{
				Uuid:      row.FormResponseUUID,
				FormUuid:  row.FormUUID,
				CreatedAt: row.FormResponseCreatedAt,
				UpdatedAt: row.FormResponseUpdatedAt,
			}
			formResponseMap[row.FormResponseUUID] = formResponse
		}

		if row.FormFieldResponseUUID != nil {
			formResponse.Responses = append(formResponse.Responses, models.FormFieldResponse{
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

	formResponses := make([]models.FormResponse, 0, len(formResponseMap))
	for _, formResponse := range formResponseMap {
		formResponses = append(formResponses, *formResponse)
	}

	return formResponses, nil
}

func (r *FormResponseRepository) GetByUUID(formResponseUUID uuid.UUID) (*models.FormResponse, error) {
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
		FROM fluxton.form_responses fr
		LEFT JOIN fluxton.form_field_responses ffr 
			ON fr.uuid = ffr.form_response_uuid
		WHERE fr.uuid = $1
		ORDER BY fr.created_at, ffr.created_at;
	`

	rows, err := r.db.Queryx(query, formResponseUUID)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	formResponseMap := make(map[uuid.UUID]*models.FormResponse)

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
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}

		formResponse, exists := formResponseMap[row.FormResponseUUID]
		if !exists {
			formResponse = &models.FormResponse{
				Uuid:      row.FormResponseUUID,
				FormUuid:  row.FormUUID,
				CreatedAt: row.FormResponseCreatedAt,
				UpdatedAt: row.FormResponseUpdatedAt,
			}
			formResponseMap[row.FormResponseUUID] = formResponse
		}

		if row.FormFieldResponseUUID != nil {
			formResponse.Responses = append(formResponse.Responses, models.FormFieldResponse{
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
		return nil, errs.NewNotFoundError("formResponse.error.notFound")
	}

	formResponse := formResponseMap[formResponseUUID]

	return formResponse, nil
}

func (r *FormResponseRepository) Create(
	formResponse *models.FormResponse,
	formFieldResponse *[]models.FormFieldResponse,
) (*models.FormResponse, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `INSERT INTO fluxton.form_responses (form_uuid) VALUES ($1) RETURNING uuid`

	queryErr := tx.QueryRowx(
		query,
		formResponse.FormUuid).Scan(&formResponse.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	for _, ffr := range *formFieldResponse {
		query = `INSERT INTO fluxton.form_field_responses (form_response_uuid, form_field_uuid, value) VALUES ($1, $2, $3) RETURNING uuid`
		queryErr = tx.QueryRowx(
			query,
			formResponse.Uuid,
			ffr.FormFieldUuid,
			ffr.Value).Scan(&ffr.Uuid)
		if queryErr != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	formResponse.Responses = *formFieldResponse

	return formResponse, nil
}

func (r *FormResponseRepository) Delete(formResponseUUID uuid.UUID) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `DELETE FROM fluxton.form_field_responses WHERE form_response_uuid = $1`
	_, err = tx.Exec(query, formResponseUUID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return utils.FormatError(err, "delete", utils.GetMethodName())
	}

	query = `DELETE FROM fluxton.form_responses WHERE uuid = $1`
	_, err = tx.Exec(query, formResponseUUID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return utils.FormatError(err, "delete", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return nil
}
