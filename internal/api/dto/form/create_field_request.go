package form

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/guregu/null/v6"
	"github.com/labstack/echo/v4"
	"regexp"
)

const (
	FieldTypeText     = "text"
	FieldTypeTextarea = "textarea"
	FieldTypeNumber   = "number"
	FieldTypeEmail    = "email"
	FieldTypeDate     = "date"
	FieldTypeCheckbox = "checkbox"
	FieldTypeRadio    = "radio"
	FieldTypeSelect   = "select"
)

var allowedFieldTypes = []interface{}{
	FieldTypeText,
	FieldTypeTextarea,
	FieldTypeNumber,
	FieldTypeEmail,
	FieldTypeDate,
	FieldTypeCheckbox,
	FieldTypeRadio,
	FieldTypeSelect,
}

type FieldRequest struct {
	// required fields
	Label      string `json:"label"`
	Type       string `json:"type"`
	IsRequired bool   `json:"is_required"`

	// all fields from this point are optional
	MinLength    null.Int    `json:"min_length" swaggertype:"integer"`
	MaxLength    null.Int    `json:"max_length" swaggertype:"integer"`
	Pattern      null.String `json:"pattern" swaggertype:"string"`
	Description  null.String `json:"description" swaggertype:"string"`
	Options      null.String `json:"options" swaggertype:"string"` // Options for select/radio types
	DefaultValue null.String `json:"default_value" swaggertype:"string"`

	// only applicable for number types
	MinValue null.Int `json:"min_value" swaggertype:"integer"`
	MaxValue null.Int `json:"max_value" swaggertype:"integer"`

	// only applicable for date types
	StartDate  null.String `json:"start_date" swaggertype:"string"`
	EndDate    null.String `json:"end_date" swaggertype:"string"`
	DateFormat null.String `json:"date_format" swaggertype:"string"` // fails if provided and field value doesn't match
}

// CreateFormFieldsRequest represents multiple fields in a request
type CreateFormFieldsRequest struct {
	dto.BaseRequest
	Fields []FieldRequest `json:"fields"`
}

// BindAndValidate binds and validates the request
func (r *CreateFormFieldsRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload: " + err.Error()}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	err = validation.ValidateStruct(r,
		validation.Field(&r.Fields, validation.Required.Error("Fields array is required"), validation.Each(
			validation.By(validateField),
		)),
	)

	return r.ExtractValidationErrors(err)
}

func validateField(value interface{}) error {
	field, ok := value.(FieldRequest)
	if !ok {
		return fmt.Errorf("invalid field format")
	}

	return validation.ValidateStruct(&field,
		validation.Field(
			&field.Label,
			validation.Required.Error("Label is required"),
			validation.Length(
				constants.MinFormFieldLabelLength, constants.MaxFormFieldLabelLength,
			).Error(
				fmt.Sprintf(
					"Label name must be between %d and %d characters",
					constants.MinFormFieldLabelLength,
					constants.MaxFormFieldLabelLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Label must be alphanumeric with underscores and dashes"),
		),
		validation.Field(
			&field.Type,
			validation.Required.Error("Type is required"),
			validation.In(allowedFieldTypes...).Error("Invalid field type"),
		),
		validation.Field(&field.IsRequired, validation.Required.Error("IsRequired is required")),
		validation.Field(
			&field.Description,
			validation.Length(
				constants.MinFormFieldDescriptionLength, constants.MaxFormFieldDescriptionLength,
			).Error("Description must be between 0 and 255 characters"),
		),
	)
}
