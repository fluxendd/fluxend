package form_requests

import (
	"github.com/labstack/echo/v4"
)

type CreateResponseRequest struct {
	Response map[string]interface{} `json:"response"`
}

func (r *CreateResponseRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	return nil
}
