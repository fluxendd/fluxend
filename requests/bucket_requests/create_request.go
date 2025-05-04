package bucket_requests

import (
	"fluxton/constants"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CreateRequest struct {
	requests.BaseRequest
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	MaxFileSize int    `json:"max_file_size"`
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
				constants.MinBucketNameLength, constants.MaxBucketNameLength,
			).Error(
				fmt.Sprintf(
					"Bucket name must be between %d and %d characters",
					constants.MinBucketNameLength,
					constants.MaxBucketNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Bucket name must be alphanumeric with underscores and dashes")),
		validation.Field(&r.IsPublic, validation.Required.Error("IsPublic is required")),
		validation.Field(&r.MaxFileSize,
			validation.Required.Error("max_file_size is required"),
			validation.Min(1).Error("max_file_size must be a positive number"),
		),
		validation.Field(
			&r.Description,
			validation.Length(constants.MinBucketDescriptionLength, constants.MaxBucketDescriptionLength).Error(
				fmt.Sprintf(
					"Bucket description must be less than %d characters",
					constants.MaxBucketDescriptionLength,
				),
			),
		),
	)

	return r.ExtractValidationErrors(err)
}
