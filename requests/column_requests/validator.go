package column_requests

import (
	"errors"
	"fluxton/internal/api/dto"
	"fluxton/models"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

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

func validateForeignKeyConstraints(column models.Column) validation.RuleFunc {
	return func(value interface{}) error {
		if column.Foreign {
			if !column.ReferenceTable.Valid || !column.ReferenceColumn.Valid {
				return errors.New("reference table and column are required for foreign key constraints")
			}
		}

		return nil
	}
}
