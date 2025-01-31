package requests

import (
	"github.com/labstack/echo/v4"
	"myapp/utils"
)

type ProjectDefaultRequest struct {
	OrganizationID uint `json:"-"`
}

func (r *ProjectDefaultRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID, err := utils.ConvertStringToUint(c.Request().Header.Get("X-OrganizationID"))
	if err != nil {
		return []string{"Organization ID is required and must be a number"}
	}

	r.OrganizationID = organizationID

	return nil
}
