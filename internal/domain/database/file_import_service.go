package database

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/samber/do"
)

type Row map[string]interface{}

type FileImportService interface {
	ImportCSV(file multipart.File) ([]Column, [][]string, error)
}

type FileImportServiceImpl struct {
}

func NewFileImportService(injector *do.Injector) (FileImportService, error) {
	return &FileImportServiceImpl{}, nil
}

func (s *FileImportServiceImpl) ImportCSV(file multipart.File) ([]Column, [][]string, error) {
	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, nil, errors.New("fileImport.error.emptyFile")
	}

	// Process headers
	headers := records[0]
	if len(headers) == 0 {
		return nil, nil, errors.New("fileImport.error.emptyHeaders")
	}

	// Process data rows
	var dataRows [][]string
	if len(records) > 1 {
		dataRows = records[1:]
	}

	// Determine columns with their types
	columns, err := s.determineColumns(headers, dataRows)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to determine columns: %w", err)
	}

	return columns, dataRows, nil
}

func (s *FileImportServiceImpl) determineColumns(headers []string, dataRows [][]string) ([]Column, error) {
	columns := make([]Column, 0, len(headers))
	columnNames := make(map[string]int) // Track sanitized column names for uniqueness

	for i, header := range headers {
		name, forcedType := s.parseHeaderWithType(header)
		name = s.sanitizeColumnName(name)

		// Handle duplicate column names
		if count, exists := columnNames[name]; exists {
			columnNames[name] = count + 1
			name = fmt.Sprintf("%s_%d", name, count)
		} else {
			columnNames[name] = 1
		}

		var colType string
		var notNull bool

		if forcedType != "" {
			// Use forced type if provided
			colType = forcedType
			notNull = true // Default to not null for forced types
		} else {
			// Auto-detect type and nullability from data
			colType, notNull = s.detectColumnType(dataRows, i)
		}

		columns = append(columns, Column{
			Name:     name,
			Position: i,
			NotNull:  notNull,
			Type:     colType,
			Primary:  false,
			Unique:   false,
			Foreign:  false,
			Default:  "",
		})
	}

	return columns, nil
}

func (s *FileImportServiceImpl) detectColumnType(rows [][]string, colIndex int) (string, bool) {
	if len(rows) == 0 {
		return "text", false // Default to text if no data
	}

	// Check if all cells in the column are empty
	allEmpty := true
	for _, row := range rows {
		if colIndex < len(row) && row[colIndex] != "" {
			allEmpty = false
			break
		}
	}

	if allEmpty {
		return "text", false // Default to nullable text if all values are empty
	}

	// Try to determine type from data
	isBoolean := true
	isInteger := true
	isFloat := true
	isJSON := true
	isTimestamp := true
	maxLength := 0
	maxPrecision := 0
	maxScale := 0

	for _, row := range rows {
		// Skip empty cells when detecting type
		if colIndex >= len(row) || row[colIndex] == "" {
			continue
		}

		value := row[colIndex]

		// Update max length for varchar/text determination
		if len(value) > maxLength {
			maxLength = len(value)
		}

		// Check if value is boolean
		if isBoolean {
			lowerValue := strings.ToLower(value)
			if lowerValue != "true" && lowerValue != "false" && lowerValue != "1" && lowerValue != "0" {
				isBoolean = false
			}
		}

		// Check if value is integer
		if isInteger {
			_, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				isInteger = false
			}
		}

		// Check if value is float and track precision/scale
		if isFloat {
			floatVal, err := strconv.ParseFloat(value, 64)
			if err != nil {
				isFloat = false
			} else {
				// Calculate precision and scale
				parts := strings.Split(value, ".")
				intDigits := len(parts[0])

				var decimalDigits int
				if len(parts) > 1 {
					decimalDigits = len(parts[1])
				}

				totalDigits := intDigits + decimalDigits

				if totalDigits > maxPrecision {
					maxPrecision = totalDigits
				}

				if decimalDigits > maxScale {
					maxScale = decimalDigits
				}

				// Check if it's actually an integer
				if floatVal == math.Trunc(floatVal) && len(parts) == 1 {
					// It's an integer value without decimal point
				} else {
					// It has decimal point or fraction
					isInteger = false
				}
			}
		}

		// Check if value is valid JSON
		if isJSON {
			var js interface{}
			if err := json.Unmarshal([]byte(value), &js); err != nil {
				isJSON = false
			} else {
				// Only consider it JSON if it's an object or array
				switch js.(type) {
				case map[string]interface{}, []interface{}:
					// Valid JSON object or array
				default:
					isJSON = false
				}
			}
		}

		// Check if value is timestamp
		if isTimestamp {
			isValidDate := false

			// Try common date formats
			formats := []string{
				time.RFC3339,
				"2006-01-02T15:04:05",
				"2006-01-02 15:04:05",
				"2006-01-02",
				"01/02/2006",
				"02/01/2006",
				"2006/01/02",
			}

			for _, format := range formats {
				_, err := time.Parse(format, value)
				if err == nil {
					isValidDate = true
					break
				}
			}

			if !isValidDate {
				isTimestamp = false
			}
		}
	}

	// Determine nullability based on presence of empty cells
	nullable := false
	for _, row := range rows {
		if colIndex >= len(row) || row[colIndex] == "" {
			nullable = true
			break
		}
	}

	// Determine the type based on checks
	if isBoolean {
		return "boolean", !nullable
	}

	if isInteger {
		return "integer", !nullable
	}

	if isFloat {
		// Use numeric with appropriate precision and scale for financial data
		if maxPrecision > 0 {
			// Ensure reasonable defaults and limits
			if maxPrecision > 38 {
				maxPrecision = 38 // PostgreSQL limit
			}

			if maxScale > maxPrecision-1 {
				maxScale = maxPrecision - 1
			}

			if maxScale < 0 {
				maxScale = 0
			}

			return fmt.Sprintf("numeric(%d,%d)", maxPrecision, maxScale), !nullable
		}
		return "float", !nullable
	}

	if isJSON {
		return "json", !nullable
	}

	if isTimestamp {
		return "timestamp", !nullable
	}

	// Default to varchar for short strings, text for longer ones
	if maxLength <= 255 {
		return fmt.Sprintf("varchar(%d)", maxLength), !nullable
	}

	return "text", !nullable
}

func (s *FileImportServiceImpl) parseHeaderWithType(header string) (string, string) {
	r := regexp.MustCompile(`(.*?)\s*\[(.*?)\]$`)
	matches := r.FindStringSubmatch(header)

	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), s.parseTypeHint(matches[2])
	}

	return header, ""
}

func (s *FileImportServiceImpl) parseTypeHint(typeHint string) string {
	typeHint = strings.TrimSpace(typeHint)
	typeHint = strings.ToLower(typeHint)

	// Handle specific type formats
	if strings.HasPrefix(typeHint, "varchar") {
		r := regexp.MustCompile(`varchar:(\d+)`)
		matches := r.FindStringSubmatch(typeHint)
		if len(matches) == 2 {
			return fmt.Sprintf("varchar(%s)", matches[1])
		}
		return "varchar(255)" // Default
	}

	if strings.HasPrefix(typeHint, "numeric") {
		r := regexp.MustCompile(`numeric:(\d+),(\d+)`)
		matches := r.FindStringSubmatch(typeHint)
		if len(matches) == 3 {
			return fmt.Sprintf("numeric(%s,%s)", matches[1], matches[2])
		}
		return "numeric(10,2)" // Default
	}

	// Standard types
	switch typeHint {
	case "bool", "boolean":
		return "boolean"
	case "int", "integer":
		return "integer"
	case "float", "real":
		return "float"
	case "text":
		return "text"
	case "json":
		return "json"
	case "timestamp", "datetime":
		return "timestamp"
	default:
		return typeHint
	}
}

func (s *FileImportServiceImpl) sanitizeColumnName(name string) string {
	// Convert to snake_case
	var result strings.Builder
	var prevChar rune

	name = strings.TrimSpace(name)

	for i, char := range name {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
			if i > 0 && (unicode.IsUpper(char) || unicode.IsDigit(char)) && unicode.IsLower(prevChar) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else if unicode.IsSpace(char) || char == '-' {
			result.WriteRune('_')
		}
		prevChar = char
	}

	sanitizedName := result.String()

	// Replace multiple underscores with a single one
	reg := regexp.MustCompile(`_+`)
	sanitizedName = reg.ReplaceAllString(sanitizedName, "_")

	// Ensure it starts with a letter or underscore
	if len(sanitizedName) > 0 && !unicode.IsLetter(rune(sanitizedName[0])) && sanitizedName[0] != '_' {
		sanitizedName = "_" + sanitizedName
	}

	return sanitizedName
}

func (s *FileImportServiceImpl) convertValueToType(value string, colType string) (interface{}, error) {
	// Handle different types
	if strings.HasPrefix(colType, "varchar") || colType == "text" {
		return value, nil
	}

	if colType == "boolean" {
		lowerValue := strings.ToLower(value)
		return lowerValue == "true" || lowerValue == "1", nil
	}

	if colType == "integer" {
		return strconv.Atoi(value)
	}

	if colType == "float" {
		return strconv.ParseFloat(value, 64)
	}

	if strings.HasPrefix(colType, "numeric") {
		// For numeric types, return the string value
		// The database will handle conversion to numeric type
		return value, nil
	}

	if colType == "json" {
		var parsedJSON interface{}
		err := json.Unmarshal([]byte(value), &parsedJSON)
		return parsedJSON, err
	}

	if colType == "timestamp" {
		// Try common date formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"01/02/2006",
			"02/01/2006",
			"2006/01/02",
		}

		for _, format := range formats {
			parsedTime, err := time.Parse(format, value)
			if err == nil {
				return parsedTime, nil
			}
		}

		return nil, fmt.Errorf("could not parse timestamp: %s", value)
	}

	// Default: return as string for unknown types
	return value, nil
}
