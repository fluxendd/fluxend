package form_requests

import (
	"fluxton/constants"
	"fluxton/pkg"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	requests.BaseRequest
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(
				constants.MinProjectNameLength, constants.MaxProjectNameLength,
			).Error(
				fmt.Sprintf(
					"Form name must be between %d and %d characters",
					constants.MinProjectNameLength,
					constants.MaxProjectNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("Form name must be alphanumeric with underscores, spaces and dashes")),
	)

	return r.ExtractValidationErrors(err)
}
