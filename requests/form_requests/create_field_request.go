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

type FieldRequest struct {
	Label       string `json:"label"`
	Type        string `json:"type"`
	IsRequired  bool   `json:"is_required"`
	Description string `json:"description,omitempty"`
	Options     string `json:"options,omitempty"` // Optional for select/radio types
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

	err := validation.ValidateStruct(r,
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
			validation.In("text", "textarea", "number", "email", "date", "checkbox", "radio", "select").Error("Invalid field type"),
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
