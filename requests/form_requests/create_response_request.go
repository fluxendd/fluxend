package form_requests

import (
	"fluxton/internal/api/dto"
	"github.com/labstack/echo/v4"
)

type CreateResponseRequest struct {
	dto.BaseRequest
	Response map[string]interface{} `json:"response"`
}

func (r *CreateResponseRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	return nil
}
