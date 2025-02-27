package project_requests

import (
	"fluxton/requests"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type ProjectUpdateRequest struct {
	requests.BaseRequest
	Name string `json:"name"`
}

func (r *ProjectUpdateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_ -]*$`)).Error("Name must start with a letter and contain only alphanumeric characters, underscores, spaces, and dashes"),
			validation.Length(3, 100).Error("Name must be between 3 and 100 characters")),
	)

	return r.ExtractValidationErrors(err)
}
