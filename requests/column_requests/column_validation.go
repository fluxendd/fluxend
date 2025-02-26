package column_requests

import (
	"fmt"
	"strings"
)

var (
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
)

func validateName(value interface{}) error {
	name := value.(string)

	if reservedColumnNames[strings.ToLower(name)] {
		return fmt.Errorf("Column name '%s' is reserved and cannot be used", name)
	}

	return nil
}

func validateType(value interface{}) error {
	columnType := value.(string)

	if !allowedColumnTypes[strings.ToLower(columnType)] {
		return fmt.Errorf("column type '%s' is not allowed", columnType)
	}

	return nil
}
