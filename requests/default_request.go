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

	organizationID, err := uuid.Parse(c.Request().Header.Get("X-OrganizationID"))
	if err != nil {
		return []string{"Invalid organization ID"}
	}

	r.OrganizationID = organizationID

	return nil
}
