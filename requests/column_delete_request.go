package requests

import (
	"fluxton/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type ColumnDeleteRequest struct {
	Name           string `json:"name"`
	OrganizationID uint   `json:"-"`
}

func (r *ColumnDeleteRequest) BindAndValidate(c echo.Context) []string {
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
		validation.Field(&r.Name, validation.Required.Error("name: name of column is required")),
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
