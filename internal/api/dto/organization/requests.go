package organization

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
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
				regexp.MustCompile(constants.AlphanumericWithSpaceUnderScoreAndDashPattern),
			).Error("Organization name must be alphanumeric with underscores, spaces and dashes"),
		),
	)

	return r.ExtractValidationErrors(err)
}

type MemberCreateRequest struct {
	dto.BaseRequest
	UserID uuid.UUID `json:"user_id"`
}

func (r *MemberCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required.Error("UserID is required")),
	)

	return r.ExtractValidationErrors(err)
}
