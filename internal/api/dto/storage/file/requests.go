package file

import (
	"fluxend/internal/api/dto"
	"fluxend/internal/config/constants"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"mime/multipart"
)

type CreateRequest struct {
	dto.DefaultRequestWithProjectHeader
	FullFileName string                `json:"-" form:"full_file_name"`
	File         *multipart.FileHeader `json:"-" form:"file"`
}

type RenameRequest struct {
	dto.DefaultRequestWithProjectHeader
	FullFileName string `json:"full_file_name"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.FullFileName,
			validation.Required.Error("full_file_name is required"),
			validation.Length(
				constants.MinContainerNameLength, constants.MaxContainerNameLength,
			).Error(
				fmt.Sprintf(
					"File name must be between %d and %d characters",
					constants.MinContainerNameLength,
					constants.MaxContainerNameLength,
				),
			)),
		validation.Field(&r.File, validation.By(fileRequired)),
	)

	return r.ExtractValidationErrors(err)
}

func (r *RenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	err := validation.ValidateStruct(r,
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
			)),
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
