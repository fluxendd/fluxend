package requests

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strings"
)

type ColumnAlterRequest struct {
	Type           string    `json:"type"`
	OrganizationID uuid.UUID `json:"-"`
}

func (r *ColumnAlterRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID := uuid.MustParse(c.Request().Header.Get("X-OrganizationID"))
	if organizationID == uuid.Nil {
		return []string{"Organization ID is required and must be a UUID"}
	}

	r.OrganizationID = organizationID

	var errors []string

	// Validate base request columns
	err := validation.ValidateStruct(r,
		validation.Field(&r.Type, validation.Required.Error("New type is required for column")),
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
	if !allowedFieldTypes[strings.ToLower(r.Type)] {
		errors = append(errors, fmt.Sprintf("Column type '%s' is not allowed", r.Type))
	}

	return errors
}
