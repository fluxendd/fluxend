package requests

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strings"
)

type ColumnRenameRequest struct {
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"-"`
}

func (r *ColumnRenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID, err := uuid.Parse(c.Request().Header.Get("X-OrganizationID"))
	if err != nil {
		return []string{"Invalid organization ID"}
	}

	r.OrganizationID = organizationID

	var errors []string

	// Validate base request columns
	err = validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("New name is required for column")),
	)

	if err != nil {
		if ve, ok := err.(validation.Errors); ok {
			for _, validationErr := range ve {
				errors = append(errors, validationErr.Error())
			}
		}

		return errors
	}

	// Check for valid column types
	if reservedFieldNames[strings.ToLower(r.Name)] {
		errors = append(errors, fmt.Sprintf("Field name '%s' is reserved and cannot be used", r.Name))
	}

	return errors
}
