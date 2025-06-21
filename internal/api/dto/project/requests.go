package project

import (
	"fluxend/internal/api/dto"
	"fluxend/internal/config/constants"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	dto.DefaultRequestWithProjectHeader
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	OrganizationUUID uuid.UUID `json:"organization_uuid"`
}

type UpdateRequest struct {
	dto.DefaultRequestWithProjectHeader
	Name        string `json:"name"`
	Description string `json:"description"`
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
				constants.MinProjectNameLength, constants.MaxProjectNameLength,
			).Error(
				fmt.Sprintf(
					"Project name be between %d and %d characters",
					constants.MinProjectNameLength,
					constants.MaxProjectNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(constants.AlphanumericWithSpaceUnderScoreAndDashPattern),
			).Error("Project name must be alphanumeric with underscores, spaces and dashes")),
		validation.Field(
			&r.OrganizationUUID,
			validation.By(func(value interface{}) error {
				if uuidValue, ok := value.(uuid.UUID); ok {
					if uuidValue == uuid.Nil {
						return fmt.Errorf("organization UUID is required")
					}
				}
				return nil
			}),
		),
	)

	return r.ExtractValidationErrors(err)
}

func (r *UpdateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(3, 100).Error("Name must be between 3 and 100 characters"),
			validation.Match(
				regexp.MustCompile(constants.AlphanumericWithSpaceUnderScoreAndDashPattern),
			).Error("Project name must be alphanumeric with underscores, spaces and dashes")),
	)

	return r.ExtractValidationErrors(err)
}
