package requests

import (
	"fluxton/constants"
	"fluxton/internal/api/dto"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
	"strings"
)

type functionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateFunctionRequest struct {
	dto.BaseRequest
	Name       string              `json:"name"`
	Parameters []functionParameter `json:"parameters"`
	Definition string              `json:"definition"`
	Language   string              `json:"language"`
	ReturnType string              `json:"return_type"`
}

func (r *CreateFunctionRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	var errors []string

	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)

		return errors
	}

	errors = append(errors, r.ExtractValidationErrors(r.validate())...)

	return errors
}

func (r *CreateFunctionRequest) validate() error {
	validTypes := map[string]bool{
		"integer": true, "bigint": true, "smallint": true, "serial": true, "bigserial": true,
		"text": true, "varchar": true, "char": true, "boolean": true,
		"real": true, "double precision": true, "numeric": true,
		"json": true, "jsonb": true, "uuid": true,
		"timestamp": true, "timestamptz": true, "date": true, "time": true,
		"bytea": true, "void": true, "record": true, "table": true,
	}

	validLanguages := map[string]bool{
		"plpgsql": true,
		"sql":     true,
	}

	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("name is required"),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscorePattern()),
			).Error("name must be alphanumeric with underscores"),
			validation.Length(
				constants.MinTableNameLength, constants.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"name must be between %d and %d characters",
					constants.MinTableNameLength,
					constants.MaxTableNameLength,
				),
			),
		),

		// Validate return type
		validation.Field(
			&r.ReturnType,
			validation.Required.Error("return_type is required"),
			validation.By(func(value interface{}) error {
				if _, exists := validTypes[value.(string)]; !exists {
					return fmt.Errorf("invalid return type: %s", value.(string))
				}
				return nil
			}),
		),

		// Validate language
		validation.Field(
			&r.Language,
			validation.Required.Error("language is required"),
			validation.By(func(value interface{}) error {
				if _, exists := validLanguages[value.(string)]; !exists {
					return fmt.Errorf("invalid language: %s", value.(string))
				}
				return nil
			}),
		),

		// Validate definition
		validation.Field(
			&r.Definition,
			validation.Required.Error("definition is required"),
			validation.By(func(value interface{}) error {
				if !strings.Contains(value.(string), "BEGIN") || !strings.Contains(value.(string), "END") {
					return fmt.Errorf("invalid definition: %s", value.(string))
				}

				return nil
			}),
		),

		// Validate function parameters
		validation.Field(
			&r.Parameters,
			validation.By(func(value interface{}) error {
				params, ok := value.([]functionParameter)
				if !ok {
					return fmt.Errorf("invalid parameters format")
				}
				for _, param := range params {
					if _, exists := validTypes[param.Type]; !exists {
						return fmt.Errorf("invalid parameter type: %s", param.Type)
					}
				}
				return nil
			}),
		),
	)
}
