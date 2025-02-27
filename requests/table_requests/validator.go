package table_requests

import (
	"fluxton/requests"
	"fmt"
	"strings"
)

func validateName(value interface{}) error {
	name := value.(string)

	if requests.IsReservedTableName(strings.ToLower(name)) {
		return fmt.Errorf("table name '%s' is reserved and cannot be used", name)
	}

	return nil
}
