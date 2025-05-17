package table_requests

import (
	"fluxton/constants"
	"fluxton/pkg"
	"fluxton/requests"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"regexp"
)

type UploadRequest struct {
	requests.BaseRequest
	Name string                `form:"name"`
	File *multipart.FileHeader `form:"file"` // use FileHeader, not string
}

func (r *UploadRequest) BindAndValidate(c echo.Context) []string {
	file, err := c.FormFile("file")
	if err != nil {
		return []string{"File is required"}
	}

	r.Name = c.FormValue("name")
	r.File = file

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	var errors []string
	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)
		return errors
	}

	return errors
}

func (r *UploadRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscorePattern()),
			).Error("Table name must be alphanumeric with underscores"),
			validation.Length(
				constants.MinTableNameLength, constants.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"Name must be between %d and %d characters",
					constants.MinTableNameLength,
					constants.MaxTableNameLength,
				),
			),
			validation.By(validateName),
		),
		validation.Field(
			&r.File,
			validation.Required.Error("File is required"),
		),
	)
}
