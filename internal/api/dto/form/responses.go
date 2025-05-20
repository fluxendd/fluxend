package form

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Response struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	CreatedBy   uuid.UUID `json:"createdB"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

type FieldResponse struct {
	Uuid         uuid.UUID   `json:"uuid"`
	FormUuid     uuid.UUID   `json:"formUuid"`
	Label        string      `json:"label"`
	Type         string      `json:"type"`
	IsRequired   bool        `json:"isRequired"`
	Description  null.String `json:"description" swaggertype:"string"`
	Options      null.String `json:"options" swaggertype:"string"`
	MinLength    null.Int    `db:"min_length" json:"minLength" swaggertype:"integer"`
	MaxLength    null.Int    `db:"max_length" json:"maxLength" swaggertype:"integer"`
	MinValue     null.Int    `db:"min_value" json:"minValue" swaggertype:"integer"`
	MaxValue     null.Int    `db:"max_value" json:"maxValue" swaggertype:"integer"`
	Pattern      null.String `db:"pattern" json:"pattern" swaggertype:"string"`
	DefaultValue null.String `db:"default_value" json:"defaultValue" swaggertype:"string"`
	StartDate    null.String `db:"start_date" json:"startDate" swaggertype:"string"`
	EndDate      null.String `db:"end_date" json:"endDate" swaggertype:"string"`
	DateFormat   null.String `db:"date_format" json:"dateFormat" swaggertype:"string"`
	CreatedAt    string      `json:"createdAt"`
	UpdatedAt    string      `json:"updatedAt"`
}

type ResponseForAPI struct {
	Uuid      uuid.UUID             `json:"uuid"`
	FormUuid  uuid.UUID             `json:"formUuid"`
	Responses []FieldResponseForAPI `json:"responses"`
}

type FieldResponseForAPI struct {
	Uuid             uuid.UUID `json:"uuid"`
	FormResponseUuid uuid.UUID `json:"formResponseUuid"`
	FormFieldUuid    uuid.UUID `json:"formFieldUuid"`
	Value            string    `json:"value"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}
