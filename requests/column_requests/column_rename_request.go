package column_requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type ColumnRenameRequest struct {
	Name string `json:"name"`
}

func (r *ColumnRenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	// Validate base request columns
	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("New name is required for column"),
			validation.By(validateName),
		),
	)

	if err != nil {
		if ve, ok := err.(validation.Errors); ok {
			for _, validationErr := range ve {
				errors = append(errors, validationErr.Error())
			}
		}

		return errors
	}

	return errors
}
