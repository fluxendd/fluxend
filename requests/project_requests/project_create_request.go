package project_requests

import (
	"fluxton/requests"
	"fluxton/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"regexp"
)

type ProjectCreateRequest struct {
	requests.BaseRequest
	Name             string    `json:"name"`
	OrganizationUUID uuid.UUID `json:"organization_uuid"`
}

func (r *ProjectCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(3, 100).Error("Name must be between 3 and 100 characters"),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("Project name must be alphanumeric with underscores, spaces and dashes")),
		validation.Field(&r.OrganizationUUID, validation.Required.Error("Organization UUID is required")),
	)

	return r.ExtractValidationErrors(err)
}
