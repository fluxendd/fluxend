package requests

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DefaultRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}

func (r *DefaultRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID := uuid.MustParse(c.Request().Header.Get("X-OrganizationID"))
	if organizationID == uuid.Nil {
		return []string{"Organization ID is required and must be a UUID"}
	}

	r.OrganizationID = organizationID

	return nil
}
