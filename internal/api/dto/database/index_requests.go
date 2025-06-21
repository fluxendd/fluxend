package database

import (
	"fluxend/internal/api/dto"
	"fluxend/internal/config/constants"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
	"strings"
)

type CreateIndexRequest struct {
	dto.DefaultRequestWithProjectHeader
	Name     string   `json:"name"`
	Columns  []string `json:"columns"`
	IsUnique bool     `json:"is_unique"`
}

func (r *CreateIndexRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	var errors []string

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Index name is required"),
			validation.Length(
				constants.MinIndexNameLength, constants.MaxIndexNameLength,
			).Error(
				fmt.Sprintf(
					"Index name be between %d and %d characters",
					constants.MinIndexNameLength,
					constants.MaxIndexNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(constants.AlphanumericWithUnderscorePattern),
			).Error("Index name must be alphanumeric with underscores"),
		),
		validation.Field(&r.Columns, validation.Required.Error("At least one column is required")),
	)

	errors = append(errors, r.ExtractValidationErrors(err)...)
	if len(errors) > 0 {
		return errors
	}

	if dto.IsReservedIndexName(strings.ToLower(r.Name)) {
		return append(errors, fmt.Sprintf("Index name '%s' is reserved and cannot be used", r.Name))
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
