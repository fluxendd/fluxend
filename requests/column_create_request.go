package requests

import (
	"fluxton/types"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

type ColumnCreateRequest struct {
	Columns []types.TableColumn `json:"columns"`
}

func (r *ColumnCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	// Validate base request columns
	err := validation.ValidateStruct(r,
		validation.Field(&r.Columns, validation.Required.Error("At least one column is required")),
	)

	if err != nil {
		if ve, ok := err.(validation.Errors); ok {
			for _, validationErr := range ve {
				errors = append(errors, validationErr.Error())
			}
		}

		return errors
	}

	// Validate each column
	for _, column := range r.Columns {
		if column.Name == "" {
			errors = append(errors, "Field name is required")
		}

		if column.Type == "" {
			errors = append(errors, fmt.Sprintf("Field type is required for column %s", column.Name))
		}

		// Check for reserved column names
		if reservedFieldNames[strings.ToLower(column.Name)] {
			errors = append(errors, fmt.Sprintf("Field name '%s' is reserved and cannot be used", column.Name))
		}

		// Check for valid column types
		if !allowedFieldTypes[strings.ToLower(column.Type)] {
			errors = append(errors, fmt.Sprintf("Field type '%s' is not allowed", column.Type))
		}
	}

	return errors
}
