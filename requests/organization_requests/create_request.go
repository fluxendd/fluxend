package organization_requests

import (
	"fluxton/constants"
	"fluxton/internal/api/dto"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	dto.BaseRequest
	Name string `json:"name"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	r.SetContext(c)

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(
				constants.MinOrganizationNameLength, constants.MaxOrganizationNameLength,
			).Error(
				fmt.Sprintf(
					"Organization name be between %d and %d characters",
					constants.MinOrganizationNameLength,
					constants.MaxOrganizationNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("Organization name must be alphanumeric with underscores, spaces and dashes"),
		),
	)

	return r.ExtractValidationErrors(err)
}
