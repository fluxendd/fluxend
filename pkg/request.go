package pkg

import (
	"strings"
)

const DefaultSchema = "public"

func ParseTableName(fullName string) (schema string, table string) {
	parts := strings.SplitN(fullName, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1] // schema, table
	}

	return DefaultSchema, fullName // Default to "public" schema if none provided
}
