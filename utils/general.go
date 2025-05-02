package utils

import (
	"fmt"
	"regexp"
)

func FormatError(err error, errType, method string) error {
	// bucketRepo.ListForProject: select.err => <error>
	return fmt.Errorf("%s: %s.err => %v", method, errType, err)
}

func MatchRegex(input, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return re.MatchString(input), nil
}
