package table

import (
	"fluxton/internal/api/dto"
	column2 "fluxton/internal/api/dto/database/column"
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/database/column"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	dto.BaseRequest
	Name    string          `json:"name"`
	Columns []column.Column `json:"columns"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
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

	for _, column := range r.Columns {
		if err := column2.ValidateColumn(column); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
}

func (r *CreateRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscorePattern()),
			).Error("Table name must be alphanumeric with underscores"),
			validation.Length(
				constants.MinTableNameLength, constants.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"Name must be between %d and %d characters",
					constants.MinTableNameLength,
					constants.MaxTableNameLength,
				),
			),
			validation.By(validateName),
		),
		validation.Field(
			&r.Columns,
			validation.Required.Error("Columns are required"),
		),
	)
}
