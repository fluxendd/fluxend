package column_requests

import (
	"fluxton/types"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type ColumnCreateRequest struct {
	Columns []types.TableColumn `json:"columns"`
}

func (r *ColumnCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	for _, column := range r.Columns {
		// Validate column name and type
		err := validation.ValidateStruct(&column,
			validation.Field(
				&column.Name,
				validation.Required.Error("Column name is required"),
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
		)

		if err != nil {
			if ve, ok := err.(validation.Errors); ok {
				for _, validationErr := range ve {
					errors = append(errors, fmt.Sprintf("Column '%s': %s", column.Name, validationErr.Error()))
				}
			}
		}
	}

	return errors
}
