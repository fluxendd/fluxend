package table

import (
	"fluxton/internal/api/dto"
	columnDto "fluxton/internal/api/dto/database/column"
	"fluxton/internal/config/constants"
	columnDomain "fluxton/internal/domain/database"
	"strings"

	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"regexp"
)

type CreateRequest struct {
	dto.BaseRequest
	Name    string                `json:"name"`
	Columns []columnDomain.Column `json:"columns"`
}

type RenameRequest struct {
	dto.BaseRequest
	Name string `json:"name"`
}

type UploadRequest struct {
	dto.BaseRequest
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

func (r *RenameRequest) BindAndValidate(c echo.Context) []string {
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
			&r.Name, validation.Required.Error("Name is required"),
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
	)

	return r.ExtractValidationErrors(err)
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

	var errors []string

	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)

		return errors
	}

	for _, currentColumn := range r.Columns {
		if err := columnDto.Validate(currentColumn); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
}

func (r *CreateRequest) validate() error {
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
			&r.Columns,
			validation.Required.Error("Columns are required"),
		),
	)
}

func validateName(value interface{}) error {
	name := value.(string)

	if dto.IsReservedTableName(strings.ToLower(name)) {
		return fmt.Errorf("table name '%s' is reserved and cannot be used", name)
	}

	return nil
}
