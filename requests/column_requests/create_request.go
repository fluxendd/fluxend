package column_requests

import (
	"fluxton/configs"
	"fluxton/models"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	requests.BaseRequest
	Columns []models.Column `json:"columns"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload: " + err.Error()}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	var errors []string

	for _, column := range r.Columns {
		if err := ValidateColumn(column); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
}

func ValidateColumn(column models.Column) error {
	return validation.ValidateStruct(&column,
		validation.Field(
			&column.Name,
			validation.Required.Error("Column name is required"),
			validation.Length(
				configs.MinColumnNameLength, configs.MaxColumnNameLength,
			).Error(
				fmt.Sprintf(
					"Column name be between %d and %d characters",
					configs.MinColumnNameLength,
					configs.MaxTableNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Column name must be alphanumeric and start with a letter"),
			validation.By(validateName),
		),
		validation.Field(
			&column.Type,
			validation.Required.Error("Column type is required"),
			validation.By(validateType),
		),
		validation.Field(
			&column.Foreign,
			validation.By(validateForeignKeyConstraints(column)),
		),
	)
}
