package column_requests

import (
	"fluxton/requests"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type ColumnRenameRequest struct {
	requests.BaseRequest
	Name string `json:"name"`
}

func (r *ColumnRenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("New name is required for column"),
			validation.By(validateName),
		),
	)

	return r.ExtractValidationErrors(err)
}
