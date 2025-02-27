package organization_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	requests.BaseRequest
	Name string `json:"name"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(
				configs.MinOrganizationNameLength, configs.MaxOrganizationNameLength,
			).Error(
				fmt.Sprintf(
					"Organization name be between %d and %d characters",
					configs.MinOrganizationNameLength,
					configs.MaxOrganizationNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("Organization name must be alphanumeric with underscores, spaces and dashes"),
		),
	)

	return r.ExtractValidationErrors(err)
}
