package form

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"time"
)

type FieldInput struct {
	// required fields
	Label      string
	Type       string
	IsRequired bool

	// all fields from this point are optional
	MinLength    null.Int
	MaxLength    null.Int
	Pattern      null.String
	Description  null.String
	Options      null.String
	DefaultValue null.String

	// only applicable for number types
	MinValue null.Int
	MaxValue null.Int

	// only applicable for date types
	StartDate  null.String
	EndDate    null.String
	DateFormat null.String
}

type CreateFormFieldsInput struct {
	ProjectUUID uuid.UUID
	Fields      []FieldInput
}

type UpdateFormFieldsInput struct {
	ProjectUUID uuid.UUID
	FieldInput
}

func (ff Field) GetTableName() string {
	return "fluxton.form_fields"
}
