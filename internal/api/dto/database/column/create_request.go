package column

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/database/column"
	"fluxton/pkg"
	"fluxton/requests/column_requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	dto.BaseRequest
	Columns []column.Column `json:"columns"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload: " + err.Error()}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	var errors []string

	for _, column := range r.Columns {
		if err := ValidateColumn(column); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
}

func ValidateColumn(column column.Column) error {
	return validation.ValidateStruct(&column,
		validation.Field(
			&column.Name,
			validation.Required.Error("Column name is required"),
			validation.Length(
				constants.MinColumnNameLength, constants.MaxColumnNameLength,
			).Error(
				fmt.Sprintf(
					"Column name be between %d and %d characters",
					constants.MinColumnNameLength,
					constants.MaxTableNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Column name must be alphanumeric and start with a letter"),
			validation.By(column_requests.validateName),
		),
		validation.Field(
			&column.Type,
			validation.Required.Error("Column type is required"),
			validation.By(column_requests.validateType),
		),
		validation.Field(
			&column.Foreign,
			validation.By(column_requests.validateForeignKeyConstraints(column)),
		),
	)
}
