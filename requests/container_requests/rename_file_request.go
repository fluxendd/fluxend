package container_requests

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type RenameFileRequest struct {
	dto.BaseRequest
	FullFileName string `json:"full_file_name"`
}

func (r *RenameFileRequest) BindAndValidate(c echo.Context) []string {
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
			&r.FullFileName,
			validation.Required.Error("Name is required"),
			validation.Length(
				constants.MinContainerNameLength, constants.MaxContainerNameLength,
			).Error(
				fmt.Sprintf(
					"File name must be between %d and %d characters",
					constants.MinContainerNameLength,
					constants.MaxContainerNameLength,
				),
			),
			/*validation.Match(
				regexp.MustCompile(utils.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("File name must be alphanumeric with underscores, spaces and dashes")*/),
	)

	return r.ExtractValidationErrors(err)
}
