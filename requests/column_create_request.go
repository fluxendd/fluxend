package requests

import (
	"fluxton/types"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strings"
)

type ColumnCreateRequest struct {
	Column         types.TableColumn `json:"column"`
	OrganizationID uuid.UUID         `json:"-"`
}

func (r *ColumnCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID := uuid.MustParse(c.Request().Header.Get("X-OrganizationID"))
	if organizationID == uuid.Nil {
		return []string{"Organization ID is required and must be a UUID"}
	}

	r.OrganizationID = organizationID

	var errors []string

	// Validate base request columns
	err := validation.ValidateStruct(r,
		validation.Field(&r.Column, validation.Required.Error("Column is required")),
	)

	if err != nil {
		if ve, ok := err.(validation.Errors); ok {
			for _, validationErr := range ve {
				errors = append(errors, validationErr.Error())
			}
		}

		return errors
	}

	// Validate column
	if r.Column.Name == "" {
		errors = append(errors, "Field name is required")
	}

	if r.Column.Type == "" {
		errors = append(errors, fmt.Sprintf("Field type is required for column %s", r.Column.Name))
	}

	// Check for reserved column names
	if reservedFieldNames[strings.ToLower(r.Column.Name)] {
		errors = append(errors, fmt.Sprintf("Field name '%s' is reserved and cannot be used", r.Column.Name))
	}

	// Check for valid column types
	if !allowedFieldTypes[strings.ToLower(r.Column.Type)] {
		errors = append(errors, fmt.Sprintf("Field type '%s' is not allowed", r.Column.Type))
	}

	return errors
}
