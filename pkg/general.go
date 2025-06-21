package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

func FormatError(err error, errType, method string) error {
	// containerRepo.ListForProject: select.err => <error>
	return fmt.Errorf("%s: %s.err => %v", method, errType, err)
}

func MatchRegex(input, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return re.MatchString(input), nil
}

func DumpJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
