package requests

import (
	"fluxton/configs"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type functionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateFunctionRequest struct {
	BaseRequest
	Name       string              `json:"name"`
	Schema     string              `json:"schema"`
	Parameters []functionParameter `json:"parameters"`
	Definition string              `json:"definition"`
	Language   string              `json:"language"`
	ReturnType string              `json:"return_type"`
}

func (r *CreateFunctionRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)

		return errors
	}

	errors = append(errors, r.validate().Error())

	return errors
}

func (r *CreateFunctionRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscorePattern()),
			).Error("Table name must be alphanumeric with underscores"),
			validation.Length(
				configs.MinTableNameLength, configs.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"Name must be between %d and %d characters",
					configs.MinTableNameLength,
					configs.MaxTableNameLength,
				),
			),
		),
	)
}
