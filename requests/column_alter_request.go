package requests

import (
	"fluxton/types"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

type ColumnAlterRequest struct {
	Columns []types.TableColumn `json:"columns"`
}

func (r *ColumnAlterRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	// Validate each column in the request
	for _, column := range r.Columns {
		// Validate column name and type
		err := validation.ValidateStruct(&column,
			validation.Field(&column.Name, validation.Required.Error("Column name is required")),
			validation.Field(&column.Type, validation.Required.Error("Column type is required")),
		)

		if err != nil {
			if ve, ok := err.(validation.Errors); ok {
				for _, validationErr := range ve {
					errors = append(errors, fmt.Sprintf("Column '%s': %s", column.Name, validationErr.Error()))
				}
			}
		}

		// Check for valid column types
		if !allowedFieldTypes[strings.ToLower(column.Type)] {
			errors = append(errors, fmt.Sprintf("Column '%s': type '%s' is not allowed", column.Name, column.Type))
		}
	}

	return errors
}
