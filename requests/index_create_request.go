package requests

import (
	"fluxton/configs"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
	"strings"
)

type IndexCreateRequest struct {
	Name     string   `json:"name"`
	Columns  []string `json:"columns"`
	IsUnique bool     `json:"is_unique"`
}

var reservedIndexNames = map[string]bool{
	"primary": true,
	"unique":  true,
	"foreign": true,
	"exclude": true,
}

func (r *IndexCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	// Validate base request fields
	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Index name is required"),
			validation.Length(
				configs.MinIndexNameLength, configs.MaxIndexNameLength,
			).Error(
				fmt.Sprintf(
					"Index name be between %d and %d characters",
					configs.MinIndexNameLength,
					configs.MaxIndexNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscorePattern()),
			).Error("Index name must be alphanumeric with underscores"),
		),
		validation.Field(&r.Columns, validation.Required.Error("At least one column is required")),
	)

	if err != nil {
		if ve, ok := err.(validation.Errors); ok {
			for _, validationErr := range ve {
				errors = append(errors, validationErr.Error())
			}
		}

		return errors
	}

	// Validate index name restrictions
	if reservedIndexNames[strings.ToLower(r.Name)] {
		errors = append(errors, fmt.Sprintf("Index name '%s' is reserved and cannot be used", r.Name))
	}

	// Ensure unique column names for the index
	seen := make(map[string]bool)
	for _, column := range r.Columns {
		if strings.TrimSpace(column) == "" {
			errors = append(errors, "Column name in index cannot be empty")
			continue
		}

		if seen[strings.ToLower(column)] {
			errors = append(errors, fmt.Sprintf("Duplicate column '%s' in index definition", column))
		}
		seen[strings.ToLower(column)] = true
	}

	return errors
}
