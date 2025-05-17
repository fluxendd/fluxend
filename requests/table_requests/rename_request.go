package table_requests

import (
	"fluxton/constants"
	"fluxton/pkg"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
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

	r.SetContext(c)

	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Name, validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscorePattern()),
			).Error("Table name must be alphanumeric with underscores"),
			validation.Length(
				constants.MinTableNameLength, constants.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"Name must be between %d and %d characters",
					constants.MinTableNameLength,
					constants.MaxTableNameLength,
				),
			),
			validation.By(validateName),
		),
	)

	return r.ExtractValidationErrors(err)
}
