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

type CreateFieldRequest struct {
	requests.BaseRequest
	Label       string `json:"label"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsRequired  bool   `json:"isRequired"`
	Options     string `json:"options"`
}

func (r *CreateFieldRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Label,
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
			).Error("Form name must be alphanumeric with underscores and dashes"),
		),
		validation.Field(
			&r.Description,
			validation.Length(
				configs.MinFormFieldDescriptionLength, configs.MaxFormFieldDescriptionLength,
			).Error(
				fmt.Sprintf(
					"Description name must be between %d and %d characters",
					configs.MinFormFieldDescriptionLength,
					configs.MaxFormFieldDescriptionLength,
				),
			),
		),
		validation.Field(
			&r.Type,
			validation.Required.Error("Type is required"),
			validation.In("text", "textarea", "number", "email", "date", "checkbox", "radio", "select"),
		),
	)

	return r.ExtractValidationErrors(err)
}
