package project_requests

import (
	"fluxton/requests"
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
			validation.Match(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_ -]*$`)).Error("Name must start with a letter and contain only alphanumeric characters, underscores, spaces, and dashes"),
			validation.Length(3, 100).Error("Name must be between 3 and 100 characters")),
		validation.Field(&r.OrganizationUUID, validation.Required.Error("Organization UUID is required")),
	)

	return r.ExtractValidationErrors(err)
}
