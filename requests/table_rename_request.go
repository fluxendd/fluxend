package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TableRenameRequest struct {
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"-"`
}

func (r *TableRenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID := uuid.MustParse(c.Request().Header.Get("X-OrganizationID"))
	if organizationID == uuid.Nil {
		return []string{"Organization ID is required and must be a UUID"}
	}

	r.OrganizationID = organizationID

	err := validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Name is required"), validation.Length(3, 100).Error("Name must be between 3 and 100 characters")),
	)

	if err == nil {
		return nil
	}

	var errors []string
	if ve, ok := err.(validation.Errors); ok {
		for _, validationErr := range ve {
			errors = append(errors, validationErr.Error())
		}
	}

	return errors
}
