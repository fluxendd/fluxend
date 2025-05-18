package file_import

import (
	"bytes"
	"encoding/csv"
	"fluxton/internal/config/constants"
	"mime/multipart"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestNewFileImportService(t *testing.T) {
	injector := do.New()
	service, err := NewFileImportService(injector)

	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.IsType(t, &ServiceImpl{}, service)
}

func TestImportCSV_EmptyFile(t *testing.T) {
	service := &ServiceImpl{}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Flush()

	file := createMultipartFile(t, buf.Bytes())

	columns, rows, err := service.ImportCSV(file)

	assert.Error(t, err)
	assert.Nil(t, columns)
	assert.Nil(t, rows)
	assert.Equal(t, "fileImport.error.emptyFile", err.Error())
}

func TestImportCSV_HeadersOnly(t *testing.T) {
	service := &ServiceImpl{}

	// Create a CSV with headers but no data
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"name", "age", "active"})
	assert.NoError(t, err)
	writer.Flush()

	file := createMultipartFile(t, buf.Bytes())

	columns, rows, err := service.ImportCSV(file)

	assert.NoError(t, err)
	assert.Len(t, columns, 3)
	assert.Empty(t, rows)

	// When no data is present, we default to text type, and they are all consider nullable
	assert.Equal(t, "name", columns[0].Name)
	assert.Equal(t, "text", columns[0].Type)
	assert.False(t, columns[0].NotNull)

	assert.Equal(t, "age", columns[1].Name)
	assert.Equal(t, "text", columns[1].Type)
	assert.False(t, columns[1].NotNull)

	assert.Equal(t, "active", columns[2].Name)
	assert.Equal(t, "text", columns[2].Type)
	assert.False(t, columns[2].NotNull)
}

func TestImportCSV_CompleteData(t *testing.T) {
	service := &ServiceImpl{}

	// Create a CSV with headers and data
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"name", "age", "active", "salary", "details"})
	assert.NoError(t, err)
	err = writer.Write([]string{"John", "30", "true", "50000.50", "{\"key\": \"value\"}"})
	assert.NoError(t, err)
	err = writer.Write([]string{"Jane", "25", "false", "60000.75", ""})
	assert.NoError(t, err)
	writer.Flush()

	file := createMultipartFile(t, buf.Bytes())

	columns, rows, err := service.ImportCSV(file)

	assert.NoError(t, err)
	assert.Len(t, columns, 5)
	assert.Len(t, rows, 2)

	assert.Equal(t, "name", columns[0].Name)
	assert.Equal(t, "varchar(4)", columns[0].Type)
	assert.True(t, columns[0].NotNull)

	assert.Equal(t, "age", columns[1].Name)
	assert.Equal(t, constants.ColumnTypeInteger, columns[1].Type)
	assert.True(t, columns[1].NotNull)

	assert.Equal(t, "active", columns[2].Name)
	assert.Equal(t, constants.ColumnTypeBoolean, columns[2].Type)
	assert.True(t, columns[2].NotNull)

	assert.Equal(t, "salary", columns[3].Name)
	assert.Equal(t, "numeric(7,2)", columns[3].Type) // TODO: handle column types with precision and scale etc in requests
	assert.True(t, columns[3].NotNull)

	assert.Equal(t, "details", columns[4].Name)
	assert.Equal(t, constants.ColumnTypeJSON, columns[4].Type)
	assert.False(t, columns[4].NotNull) // the second row contains empty value so not null is false
}

func TestDetermineColumns(t *testing.T) {
	service := &ServiceImpl{}

	headers := []string{"Name", "Age", "Email"}
	dataRows := [][]string{
		{"John Doe", "30", "john@example.com"},
		{"Jane Smith", "25", "jane@example.com"},
	}

	columns, err := service.determineColumns(headers, dataRows)

	assert.NoError(t, err)
	assert.Len(t, columns, 3)

	// Check sanitized column names
	assert.Equal(t, "name", columns[0].Name)
	assert.Equal(t, "age", columns[1].Name)
	assert.Equal(t, "email", columns[2].Name)
}

func TestImportCSV_MixedNullability(t *testing.T) {
	service := &ServiceImpl{}

	// Create a CSV with headers and data, including nulls
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"name", "age", "active", "salary"})
	assert.NoError(t, err)
	err = writer.Write([]string{"John", "30", "true", "50000.50"})
	assert.NoError(t, err)
	err = writer.Write([]string{"Jane", "", "false", ""})
	assert.NoError(t, err)
	writer.Flush()

	file := createMultipartFile(t, buf.Bytes())
	columns, rows, err := service.ImportCSV(file)

	assert.NoError(t, err)
	assert.Len(t, columns, 4)
	assert.Len(t, rows, 2)

	// Check nullability
	assert.Equal(t, "name", columns[0].Name)
	assert.True(t, columns[0].NotNull)

	assert.Equal(t, "age", columns[1].Name)
	assert.False(t, columns[1].NotNull)

	assert.Equal(t, "active", columns[2].Name)
	assert.True(t, columns[2].NotNull)

	assert.Equal(t, "salary", columns[3].Name)
	assert.False(t, columns[3].NotNull)
}

func TestImportCSV_TypedHeaders(t *testing.T) {
	service := &ServiceImpl{}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"name [text]", "age [integer]", "joined [timestamp]", "balance [numeric:12,2]"})
	assert.NoError(t, err)
	err = writer.Write([]string{"John", "30", "2023-01-01", "1234.56"})
	assert.NoError(t, err)
	writer.Flush()

	file := createMultipartFile(t, buf.Bytes())
	columns, _, err := service.ImportCSV(file)

	assert.NoError(t, err)
	assert.Len(t, columns, 4)

	// Check if type hints were properly applied
	assert.Equal(t, "name", columns[0].Name)
	assert.Equal(t, constants.ColumnTypeText, columns[0].Type)

	assert.Equal(t, "age", columns[1].Name)
	assert.Equal(t, constants.ColumnTypeInteger, columns[1].Type)

	assert.Equal(t, "joined", columns[2].Name)
	assert.Equal(t, constants.ColumnTypeTimestamp, columns[2].Type)

	assert.Equal(t, "balance", columns[3].Name)
	assert.Equal(t, "numeric(12,2)", columns[3].Type)
}

func TestParseHeaderWithType(t *testing.T) {
	service := &ServiceImpl{}

	testCases := []struct {
		header     string
		name       string
		forcedType string
	}{
		{
			header:     "Name",
			name:       "Name",
			forcedType: "",
		},
		{
			header:     "Name [text]",
			name:       "Name",
			forcedType: "text",
		},
		{
			header:     "Age [integer]",
			name:       "Age",
			forcedType: "integer",
		},
		{
			header:     "Price [numeric:10,2]",
			name:       "Price",
			forcedType: "numeric(10,2)",
		},
		{
			header:     "Email Address [varchar:100]",
			name:       "Email Address",
			forcedType: "varchar(100)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.header, func(t *testing.T) {
			name, forcedType := service.parseHeaderWithType(tc.header)

			assert.Equal(t, tc.name, name)
			assert.Equal(t, tc.forcedType, forcedType)
		})
	}
}

func TestSanitizeColumnName(t *testing.T) {
	service := &ServiceImpl{}

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Name",
			expected: "name",
		},
		{
			input:    "First Name",
			expected: "first_name",
		},
		{
			input:    "EmailAddress",
			expected: "email_address",
		},
		{
			input:    "123Name",
			expected: "_123name",
		},
		{
			input:    "user-id",
			expected: "user_id",
		},
		{
			input:    "  Spaced   Name  ",
			expected: "spaced_name",
		},
		{
			input:    "multiple__underscores___here",
			expected: "multiple_underscores_here",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := service.sanitizeColumnName(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertValueToType(t *testing.T) {
	service := &ServiceImpl{}

	testCases := []struct {
		name      string
		value     string
		colType   string
		expected  interface{}
		expectErr bool
	}{
		{
			name:      "String to text",
			value:     "Hello",
			colType:   "text",
			expected:  "Hello",
			expectErr: false,
		},
		{
			name:      "String to varchar",
			value:     "World",
			colType:   "varchar(10)",
			expected:  "World",
			expectErr: false,
		},
		{
			name:      "String to boolean - true",
			value:     "true",
			colType:   "boolean",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "String to boolean - 1",
			value:     "1",
			colType:   "boolean",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "String to boolean - false",
			value:     "false",
			colType:   "boolean",
			expected:  false,
			expectErr: false,
		},
		{
			name:      "String to integer - valid",
			value:     "123",
			colType:   "integer",
			expected:  123,
			expectErr: false,
		},
		{
			name:      "String to integer - invalid",
			value:     "abc",
			colType:   "integer",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "String to float - valid",
			value:     "123.45",
			colType:   "float",
			expected:  123.45,
			expectErr: false,
		},
		{
			name:      "String to float - invalid",
			value:     "abc",
			colType:   "float",
			expected:  0.0,
			expectErr: true,
		},
		{
			name:      "String to numeric",
			value:     "123.45",
			colType:   "numeric(10,2)",
			expected:  "123.45",
			expectErr: false,
		},
		{
			name:      "String to json - valid",
			value:     `{"name":"John"}`,
			colType:   "json",
			expected:  map[string]interface{}{"name": "John"},
			expectErr: false,
		},
		{
			name:      "String to json - invalid",
			value:     `{invalid json}`,
			colType:   "json",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "String to timestamp - ISO format",
			value:     "2023-01-01T12:00:00",
			colType:   "timestamp",
			expectErr: false,
		},
		{
			name:      "String to timestamp - date only",
			value:     "2023-01-01",
			colType:   "timestamp",
			expectErr: false,
		},
		{
			name:      "String to timestamp - invalid",
			value:     "not a date",
			colType:   "timestamp",
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.convertValueToType(tc.value, tc.colType)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Special case for timestamp which we can't directly compare
				if tc.colType == "timestamp" && !tc.expectErr {
					assert.NotNil(t, result)
				} else {
					assert.Equal(t, tc.expected, result)
				}
			}
		})
	}
}

func TestDetectColumnType(t *testing.T) {
	service := &ServiceImpl{}

	testCases := []struct {
		name         string
		values       []string
		expectedType string
		notNull      bool
	}{
		{
			name:         "Boolean Column",
			values:       []string{"true", "false", "true"},
			expectedType: "boolean",
			notNull:      true,
		},
		{
			name:         "Integer Column",
			values:       []string{"1", "2", "3"},
			expectedType: "integer",
			notNull:      true,
		},
		{
			name:         "Float Column",
			values:       []string{"1.1", "2.2", "3.3"},
			expectedType: "numeric(2,1)",
			notNull:      true,
		},
		{
			name:         "Text Column",
			values:       []string{"abc", "def", "ghi"},
			expectedType: "varchar(3)",
			notNull:      true,
		},
		{
			name:         "Long Text Column",
			values:       []string{"a very long text that exceeds 255 characters" + string(make([]byte, 300))},
			expectedType: "text",
			notNull:      true,
		},
		{
			name:         "JSON Column",
			values:       []string{`{"name":"John"}`, `{"name":"Jane"}`},
			expectedType: "json",
			notNull:      true,
		},
		{
			name:         "Timestamp Column",
			values:       []string{"2023-01-01", "2023-02-01"},
			expectedType: "timestamp",
			notNull:      true,
		},
		{
			name:         "Nullable Column",
			values:       []string{"1", "", "3"},
			expectedType: "integer",
			notNull:      false,
		},
		{
			name:         "Empty Column",
			values:       []string{"", "", ""},
			expectedType: "text",
			notNull:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rows := make([][]string, len(tc.values))
			for i, val := range tc.values {
				rows[i] = []string{val}
			}

			colType, notNull := service.detectColumnType(rows, 0)

			assert.Equal(t, tc.expectedType, colType, "Type detection failed for "+tc.name)
			assert.Equal(t, tc.notNull, notNull, "Nullability detection failed for "+tc.name)
		})
	}
}

// Helper function to create a multipart file from bytes
func createMultipartFile(t *testing.T, data []byte) multipart.File {
	return &bytesFile{
		Reader: bytes.NewReader(data),
		size:   int64(len(data)),
	}
}

// Mock implementation of multipart.File interface
type bytesFile struct {
	*bytes.Reader
	size int64
}

func (f *bytesFile) Close() error {
	return nil
}
