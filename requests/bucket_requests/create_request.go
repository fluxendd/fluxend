package bucket_requests

import (
	"fluxton/configs"
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
	IsPublic    bool   `json:"isPublic"`
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
				configs.MinBucketNameLength, configs.MaxBucketNameLength,
			).Error(
				fmt.Sprintf(
					"Bucket name must be between %d and %d characters",
					configs.MinBucketNameLength,
					configs.MaxBucketNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Bucket name must be alphanumeric with underscores and dashes")),
		validation.Field(&r.IsPublic, validation.Required.Error("IsPublic is required")),
		validation.Field(
			&r.Description,
			validation.Length(configs.MinBucketDescriptionLength, configs.MaxBucketDescriptionLength).Error(
				fmt.Sprintf(
					"Bucket description must be less than %d characters",
					configs.MaxBucketDescriptionLength,
				),
			),
		),
	)

	return r.ExtractValidationErrors(err)
}
