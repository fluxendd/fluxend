package requests

import (
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

type ColumnRenameRequest struct {
	Name           string `json:"name"`
	OrganizationID uint   `json:"-"`
}

func (r *ColumnRenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID, err := utils.ConvertStringToUint(c.Request().Header.Get("X-OrganizationID"))
	if err != nil {
		return []string{"Organization ID is required and must be a number"}
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
