package utils

import (
	"fmt"
	"strconv"
)

func ConvertStringToInt(param string) (int, error) {
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("provided value is not a valid integer: %w", err)
	}

	return value, nil
}

func ConvertPointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
