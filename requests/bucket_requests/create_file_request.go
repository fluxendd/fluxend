package bucket_requests

import (
	"fluxton/constants"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"mime/multipart"
)

type CreateFileRequest struct {
	requests.BaseRequest
	FullFileName string                `json:"-" form:"full_file_name"`
	File         *multipart.FileHeader `json:"-" form:"file"`
}

func (r *CreateFileRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.FullFileName,
			validation.Required.Error("full_file_name is required"),
			validation.Length(
				constants.MinBucketNameLength, constants.MaxBucketNameLength,
			).Error(
				fmt.Sprintf(
					"File name must be between %d and %d characters",
					constants.MinBucketNameLength,
					constants.MaxBucketNameLength,
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
