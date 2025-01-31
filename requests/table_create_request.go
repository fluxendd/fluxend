package requests

import (
	"fluxton/types"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

type TableCreateRequest struct {
	Name           string             `json:"name"`
	OrganizationID uint               `json:"-"`
	Fields         []types.TableField `json:"fields"`
}

var (
	reservedTableNames = map[string]bool{
		"pg_catalog":         true,
		"information_schema": true,
	}

	reservedFieldNames = map[string]bool{
		"oid":      true,
		"xmin":     true,
		"cmin":     true,
		"xmax":     true,
		"cmax":     true,
		"tableoid": true,
	}

	allowedFieldTypes = map[string]bool{
		"int":       true,
		"serial":    true,
		"varchar":   true,
		"text":      true,
		"boolean":   true,
		"date":      true,
		"timestamp": true,
		"float":     true,
		"uuid":      true,
	}
)

func (r *TableCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	organizationID, err := utils.ConvertStringToUint(c.Request().Header.Get("X-OrganizationID"))
	if err != nil {
		return []string{"Organization ID is required and must be a number"}
	}

	r.OrganizationID = organizationID

	var errors []string

	// Validate base request fields
	err = validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Name is required"), validation.Length(3, 100).Error("Name must be between 3 and 100 characters")),
		validation.Field(&r.Fields, validation.Required.Error("Fields are required")),
	)

	if err != nil {
		if ve, ok := err.(validation.Errors); ok {
			for _, validationErr := range ve {
				errors = append(errors, validationErr.Error())
			}
		}

		return errors
	}

	// Validate table name restrictions
	if reservedTableNames[strings.ToLower(r.Name)] {
		errors = append(errors, fmt.Sprintf("Table name '%s' is not allowed", r.Name))
	}

	// Validate fields
	for _, field := range r.Fields {
		if field.Name == "" {
			errors = append(errors, "Field name is required")
		}

		if field.Type == "" {
			errors = append(errors, fmt.Sprintf("Field type is required for field %s", field.Name))
		}

		// Check for reserved field names
		if reservedFieldNames[strings.ToLower(field.Name)] {
			errors = append(errors, fmt.Sprintf("Field name '%s' is reserved and cannot be used", field.Name))
		}

		// Check for valid field types
		if !allowedFieldTypes[strings.ToLower(field.Type)] {
			errors = append(errors, fmt.Sprintf("Field type '%s' is not allowed", field.Type))
		}
	}

	return errors
}
