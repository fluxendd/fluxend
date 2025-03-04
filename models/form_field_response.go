package models

import (
	"github.com/google/uuid"
	"time"
)

type FormFieldResponse struct {
	Uuid             uuid.UUID `db:"uuid" json:"uuid"`
	FormResponseUuid uuid.UUID `db:"form_response_uuid" json:"formResponseUuid"`
	FormFieldUuid    uuid.UUID `db:"form_field_uuid" json:"formFieldUuid"`
	Value            string    `db:"value" json:"value"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

func (u FormFieldResponse) GetTableName() string {
	return "fluxton.form_field_responses"
}
