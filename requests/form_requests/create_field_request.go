package form_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

var (
	fieldTypeText     = "text"
	fieldTypeTextarea = "textarea"
	fieldTypeNumber   = "number"
	fieldTypeEmail    = "email"
	fieldTypeDate     = "date"
	fieldTypeCheckbox = "checkbox"
	fieldTypeRadio    = "radio"
	fieldTypeSelect   = "select"
)

var allowedFieldTypes = []string{
	fieldTypeText,
	fieldTypeTextarea,
	fieldTypeNumber,
	fieldTypeEmail,
	fieldTypeDate,
	fieldTypeCheckbox,
	fieldTypeRadio,
	fieldTypeSelect,
}

type FieldRequest struct {
	// required fields
	Label      string `json:"label"`
	Type       string `json:"type"`
	IsRequired bool   `json:"is_required"`

	// all fields from this point are optional
	MinLength    int    `json:"min_length,omitempty"`
	MaxLength    int    `json:"max_length,omitempty"`
	Pattern      string `json:"pattern,omitempty"`
	Description  string `json:"description,omitempty"`
	Options      string `json:"options,omitempty"` // Options for select/radio types
	DefaultValue string `json:"default_value,omitempty"`

	// only applicable for number types
	MinValue int `json:"min_value,omitempty"`
	MaxValue int `json:"max_value,omitempty"`

	// only applicable for date types
	StartDate  string `json:"start_date,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
	DateFormat string `json:"date_format,omitempty"` // fails if provided and field value doesn't match
}

// CreateFormFieldsRequest represents multiple fields in a request
type CreateFormFieldsRequest struct {
	requests.BaseRequest
	Fields []FieldRequest `json:"fields"`
}

// BindAndValidate binds and validates the request
func (r *CreateFormFieldsRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
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
			validation.In(allowedFieldTypes).Error("Invalid field type"),
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
