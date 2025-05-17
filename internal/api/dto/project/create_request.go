package project

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	dto.BaseRequest
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	OrganizationUUID uuid.UUID `json:"organization_uuid"`
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
				constants.MinProjectNameLength, constants.MaxProjectNameLength,
			).Error(
				fmt.Sprintf(
					"Project name be between %d and %d characters",
					constants.MinProjectNameLength,
					constants.MaxProjectNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("Project name must be alphanumeric with underscores, spaces and dashes")),
		validation.Field(&r.OrganizationUUID, validation.Required.Error("Organization UUID is required")),
	)

	return r.ExtractValidationErrors(err)
}
