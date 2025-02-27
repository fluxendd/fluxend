package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	reservedTableNames = map[string]bool{
		"pg_catalog":         true,
		"information_schema": true,
	}

	reservedColumnNames = map[string]bool{
		"oid":      true,
		"xmin":     true,
		"cmin":     true,
		"xmax":     true,
		"cmax":     true,
		"tableoid": true,
	}

	allowedColumnTypes = map[string]bool{
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

	reservedIndexNames = map[string]bool{
		"primary": true,
		"unique":  true,
		"foreign": true,
		"exclude": true,
	}
)

type BaseRequest struct{}

func (r *BaseRequest) ExtractValidationErrors(err error) []string {
	if err == nil {
		return nil
	}

	var errors []string
	if ve, ok := err.(validation.Errors); ok {
		for _, validationErr := range ve {
			errors = append(errors, validationErr.Error())
		}
	}

	return errors
}

func GetReservedTableNames() map[string]bool {
	return reservedTableNames
}

func IsReservedTableName(name string) bool {
	if _, ok := reservedTableNames[name]; ok {
		return true
	}

	return false
}

func GetReservedColumnNames() map[string]bool {
	return reservedColumnNames
}

func IsReservedColumnName(name string) bool {
	if _, ok := reservedColumnNames[name]; ok {
		return true
	}

	return false
}

func GetAllowedColumnTypes() map[string]bool {
	return allowedColumnTypes
}

func IsAllowedColumnType(columnType string) bool {
	if _, ok := allowedColumnTypes[columnType]; ok {
		return true
	}

	return false
}

func GetReservedIndexNames() map[string]bool {
	return reservedIndexNames
}

func IsReservedIndexName(name string) bool {
	if _, ok := reservedIndexNames[name]; ok {
		return true
	}

	return false
}
