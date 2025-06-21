package database

import (
	"fluxend/internal/api/dto"
	"fluxend/internal/config/constants"
	columnDomain "fluxend/internal/domain/database"
	"strings"

	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"regexp"
)

type CreateTableRequest struct {
	dto.DefaultRequestWithProjectHeader
	Name    string                `json:"name"`
	Columns []columnDomain.Column `json:"columns"`
}

type RenameTableRequest struct {
	dto.DefaultRequestWithProjectHeader
	Name string `json:"name"`
}

type UploadTableRequest struct {
	dto.DefaultRequestWithProjectHeader
	Name string                `form:"name"`
	File *multipart.FileHeader `form:"file"`
}

func (r *UploadTableRequest) BindAndValidate(c echo.Context) []string {
	file, err := c.FormFile("file")
	if err != nil {
		return []string{"File is required"}
	}

	r.Name = c.FormValue("name")
	r.File = file

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	var errors []string
	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)
		return errors
	}

	return errors
}

func (r *UploadTableRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(constants.AlphanumericWithUnderscorePattern),
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
			validation.By(validateTableName),
		),
		validation.Field(
			&r.File,
			validation.Required.Error("File is required"),
		),
	)
}

func (r *RenameTableRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name, validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(constants.AlphanumericWithUnderscorePattern),
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
			validation.By(validateTableName),
		),
	)

	return r.ExtractValidationErrors(err)
}

func (r *CreateTableRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	var errors []string

	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)

		return errors
	}

	for _, currentColumn := range r.Columns {
		if err := Validate(currentColumn); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
}

func (r *CreateTableRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(constants.AlphanumericWithUnderscorePattern),
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
			validation.By(validateTableName),
		),
		validation.Field(
			&r.Columns,
			validation.Required.Error("Columns are required"),
		),
	)
}

func validateTableName(value interface{}) error {
	name := value.(string)

	if dto.IsReservedTableName(strings.ToLower(name)) {
		return fmt.Errorf("table name '%s' is reserved and cannot be used", name)
	}

	return nil
}
