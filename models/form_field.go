package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"time"
)

type FormField struct {
	BaseModel
	Uuid         uuid.UUID   `db:"uuid" json:"uuid"`
	FormUuid     uuid.UUID   `db:"form_uuid" json:"formUuid"`
	Label        string      `db:"label" json:"label"`
	Type         string      `db:"type" json:"type"`
	IsRequired   bool        `db:"is_required" json:"isRequired"`
	Description  null.String `db:"description" json:"description"`
	Options      null.String `db:"options" json:"options"`
	MinLength    null.Int    `db:"min_length" json:"minLength"`
	MaxLength    null.Int    `db:"max_length" json:"maxLength"`
	MinValue     null.Int    `db:"min_value" json:"minValue"`
	MaxValue     null.Int    `db:"max_value" json:"maxValue"`
	Pattern      null.String `db:"pattern" json:"pattern"`
	DefaultValue null.String `db:"default_value" json:"defaultValue"`
	StartDate    null.String `db:"start_date" json:"startDate"`
	EndDate      null.String `db:"end_date" json:"endDate"`
	DateFormat   null.String `db:"date_format" json:"dateFormat"`
	CreatedAt    time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updatedAt"`
}

func (ff FormField) GetTableName() string {
	return "fluxton.form_fields"
}
