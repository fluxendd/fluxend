package form

import (
	"fluxton/internal/domain/form"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type FormFieldResponse struct {
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

func FormFieldResource(formField *form.Field) FormFieldResponse {
	return FormFieldResponse{
		Uuid:         formField.Uuid,
		FormUuid:     formField.FormUuid,
		Label:        formField.Label,
		Description:  formField.Description,
		Type:         formField.Type,
		IsRequired:   formField.IsRequired,
		Options:      formField.Options,
		MinLength:    formField.MinLength,
		MaxLength:    formField.MaxLength,
		MinValue:     formField.MinValue,
		MaxValue:     formField.MaxValue,
		Pattern:      formField.Pattern,
		DefaultValue: formField.DefaultValue,
		StartDate:    formField.StartDate,
		EndDate:      formField.EndDate,
		DateFormat:   formField.DateFormat,
		CreatedAt:    formField.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    formField.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FormFieldResourceCollection(formFields []form.Field) []FormFieldResponse {
	resourceFormFields := make([]FormFieldResponse, len(formFields))
	for i, formField := range formFields {
		resourceFormFields[i] = FormFieldResource(&formField)
	}

	return resourceFormFields
}
