package column

import (
	"errors"
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/database/column"
	"fluxton/pkg"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
	"strings"
)

type CreateRequest struct {
	dto.BaseRequest
	Columns []column.Column `json:"columns"`
}

type RenameRequest struct {
	dto.BaseRequest
	Name string `json:"name"`
}

func (r *CreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload: " + err.Error()}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	r.SetContext(c)

	var errors []string

	for _, currentColumn := range r.Columns {
		if err := validateColumn(currentColumn); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
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
			&r.Name,
			validation.Required.Error("New name is required for column"),
			validation.Length(
				constants.MinColumnNameLength, constants.MaxColumnNameLength,
			).Error(
				fmt.Sprintf(
					"Column name be between %d and %d characters",
					constants.MinColumnNameLength,
					constants.MaxColumnNameLength,
				),
			),
			validation.By(validateName),
		),
	)

	return r.ExtractValidationErrors(err)
}

func validateColumn(column column.Column) error {
	return validation.ValidateStruct(&column,
		validation.Field(
			&column.Name,
			validation.Required.Error("Column name is required"),
			validation.Length(
				constants.MinColumnNameLength, constants.MaxColumnNameLength,
			).Error(
				fmt.Sprintf(
					"Column name be between %d and %d characters",
					constants.MinColumnNameLength,
					constants.MaxTableNameLength,
				),
			),
			validation.Match(
				regexp.MustCompile(pkg.AlphanumericWithUnderscoreAndDashPattern()),
			).Error("Column name must be alphanumeric and start with a letter"),
			validation.By(validateName),
		),
		validation.Field(
			&column.Type,
			validation.Required.Error("Column type is required"),
			validation.By(validateType),
		),
		validation.Field(
			&column.Foreign,
			validation.By(validateForeignKeyConstraints(column)),
		),
	)
}

func validateName(value interface{}) error {
	name := value.(string)

	if dto.IsReservedColumnName(strings.ToLower(name)) {
		return fmt.Errorf("column name '%s' is reserved and cannot be used", name)
	}

	return nil
}

func validateType(value interface{}) error {
	columnType := value.(string)

	if !dto.IsAllowedColumnType(strings.ToLower(columnType)) {
		return fmt.Errorf("column type '%s' is not allowed", columnType)
	}

	return nil
}

func validateForeignKeyConstraints(column column.Column) validation.RuleFunc {
	return func(value interface{}) error {
		if column.Foreign {
			if !column.ReferenceTable.Valid || !column.ReferenceColumn.Valid {
				return errors.New("reference table and column are required for foreign key constraints")
			}
		}

		return nil
	}
}
