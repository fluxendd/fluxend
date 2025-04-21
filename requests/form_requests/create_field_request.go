package form_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fluxton/utils"
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
	MinLength    null.Int    `json:"min_length"`
	MaxLength    null.Int    `json:"max_length"`
	Pattern      null.String `json:"pattern"`
	Description  null.String `json:"description"`
	Options      null.String `json:"options"` // Options for select/radio types
	DefaultValue null.String `json:"default_value"`

	// only applicable for number types
	MinValue null.Int `json:"min_value"`
	MaxValue null.Int `json:"max_value"`

	// only applicable for date types
	StartDate  null.String `json:"start_date"`
	EndDate    null.String `json:"end_date"`
	DateFormat null.String `json:"date_format"` // fails if provided and field value doesn't match
}

// CreateFormFieldsRequest represents multiple fields in a request
type CreateFormFieldsRequest struct {
	requests.BaseRequest
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
				configs.MinFormFieldLabelLength, configs.MaxFormFieldLabelLength,
			).Error(
				fmt.Sprintf(
					"Label name must be between %d and %d characters",
					configs.MinFormFieldLabelLength,
					configs.MaxFormFieldLabelLength,
				),
			),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscoreAndDashPattern()),
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
				configs.MinFormFieldDescriptionLength, configs.MaxFormFieldDescriptionLength,
			).Error("Description must be between 0 and 255 characters"),
		),
	)
}
