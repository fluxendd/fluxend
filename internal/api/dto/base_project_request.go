package dto

import (
	"github.com/labstack/echo/v4"
)

type DefaultRequestWithProjectHeader struct {
	BaseRequest
}

func (r *DefaultRequestWithProjectHeader) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	return nil
}
