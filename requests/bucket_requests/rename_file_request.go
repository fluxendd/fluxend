package bucket_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type RenameFileRequest struct {
	requests.BaseRequest
	Name string `json:"name"`
}

func (r *RenameFileRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(
				configs.MinBucketNameLength, configs.MaxBucketNameLength,
			).Error(
				fmt.Sprintf(
					"File name must be between %d and %d characters",
					configs.MinBucketNameLength,
					configs.MaxBucketNameLength,
				),
			),
			/*validation.Match(
				regexp.MustCompile(utils.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("File name must be alphanumeric with underscores, spaces and dashes")*/),
	)

	return r.ExtractValidationErrors(err)
}
