package form_requests

import (
	"fluxton/requests"
	"github.com/labstack/echo/v4"
)

type UpdateFormFieldRequest struct {
	requests.BaseRequest
	FieldRequest
}

func (r *UpdateFormFieldRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload" + err.Error()}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	return r.ExtractValidationErrors(validateField(r.FieldRequest))
}
