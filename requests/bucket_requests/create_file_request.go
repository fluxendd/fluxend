package bucket_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"mime/multipart"
)

type CreateFileRequest struct {
	requests.BaseRequest
	Name string                `json:"-" form:"name"`
	File *multipart.FileHeader `json:"-" form:"file"`
}

func (r *CreateFileRequest) BindAndValidate(c echo.Context) []string {
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
		validation.Field(&r.File, validation.By(fileRequired)),
	)

	return r.ExtractValidationErrors(err)
}

func fileRequired(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return fmt.Errorf("file is required")
	}
	return nil
}
