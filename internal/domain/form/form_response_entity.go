package form

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type FormResponse struct {
	shared.BaseEntity
	Uuid      uuid.UUID       `db:"uuid" json:"uuid"`
	FormUuid  uuid.UUID       `db:"form_uuid" json:"formUuid"`
	CreatedAt time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time       `db:"updated_at" json:"updatedAt"`
	Responses []FieldResponse `json:"responses"`
}

func (u FormResponse) GetTableName() string {
	return "fluxton.form_field_responses"
}
