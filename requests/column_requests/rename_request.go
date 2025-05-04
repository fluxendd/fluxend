package column_requests

import (
	"fluxton/constants"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type RenameRequest struct {
	requests.BaseRequest
	Name string `json:"name"`
}

func (r *RenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("New name is required for column"),
			validation.Length(
				constants.MinColumnNameLength, constants.MaxColumnNameLength,
			).Error(
				fmt.Sprintf(
					"Column name be between %d and %d characters",
					constants.MinColumnNameLength,
					constants.MaxColumnNameLength,
				),
			),
			validation.By(validateName),
		),
	)

	return r.ExtractValidationErrors(err)
}
