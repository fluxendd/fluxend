package database

import (
	"errors"
	"fluxton/internal/api/dto"
	"fluxton/internal/config/constants"
	columnDomain "fluxton/internal/domain/database"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
	"strings"
)

type CreateColumnRequest struct {
	dto.BaseRequest
	Columns []columnDomain.Column `json:"columns"`
}

type RenameColumnRequest struct {
	dto.BaseRequest
	Name string `json:"name"`
}

func (r *CreateColumnRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload: " + err.Error()}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

	var requestErrors []string

	for _, currentColumn := range r.Columns {
		if err := Validate(currentColumn); err != nil {
			requestErrors = append(requestErrors, err.Error())

			return requestErrors
		}
	}

	return requestErrors
}

func (r *RenameColumnRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := r.WithProjectHeader(c)
	if err != nil {
		return []string{err.Error()}
	}

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
			validation.By(validateColumnName),
		),
	)

	return r.ExtractValidationErrors(err)
}

func Validate(column columnDomain.Column) error {
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
				regexp.MustCompile(constants.AlphanumericWithUnderscoreAndDashPattern),
			).Error("Column name must be alphanumeric and start with a letter"),
			validation.By(validateColumnName),
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

func validateColumnName(value interface{}) error {
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

func validateForeignKeyConstraints(column columnDomain.Column) validation.RuleFunc {
	return func(value interface{}) error {
		if column.Foreign {
			if !column.ReferenceTable.Valid || !column.ReferenceColumn.Valid {
				return errors.New("reference table and column are required for foreign key constraints")
			}
		}

		return nil
	}
}
