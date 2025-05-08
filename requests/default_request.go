package requests

import (
	"github.com/labstack/echo/v4"
)

type DefaultRequest struct {
	BaseRequest
}

func (r *DefaultRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	r.SetContext(c)

	return nil
}
