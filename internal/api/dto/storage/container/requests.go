package container

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	dto.BaseRequest
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	MaxFileSize int    `json:"max_file_size"`
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
				constants.MinContainerNameLength, constants.MaxContainerNameLength,
			).Error(
				fmt.Sprintf(
					"Container name must be between %d and %d characters",
					constants.MinContainerNameLength,
					constants.MaxContainerNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Container name must be alphanumeric with underscores and dashes")),
		validation.Field(&r.IsPublic, validation.Required.Error("IsPublic is required")),
		validation.Field(&r.MaxFileSize,
			validation.Required.Error("max_file_size is required"),
			validation.Min(1).Error("max_file_size must be a positive number"),
		),
		validation.Field(
			&r.Description,
			validation.Length(constants.MinContainerDescriptionLength, constants.MaxContainerDescriptionLength).Error(
				fmt.Sprintf(
					"Container description must be less than %d characters",
					constants.MaxContainerDescriptionLength,
				),
			),
		),
	)

	return r.ExtractValidationErrors(err)
}
