package services

import (
	"bytes"
	"encoding/csv"
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
	assert.IsType(t, &FileImportServiceImpl{}, service)
}

func TestImportCSV_EmptyFile(t *testing.T) {
	service := &FileImportServiceImpl{}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Flush()

	file := createMultipartFile(t, buf.Bytes())

	columns, rows, err := service.ImportCSV(file)

	assert.Error(t, err)
	assert.Nil(t, columns)
	assert.Nil(t, rows)
	assert.Equal(t, "File is empty", err.Error())
}

func TestImportCSV_HeadersOnly(t *testing.T) {
	service := &FileImportServiceImpl{}

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
	service := &FileImportServiceImpl{}

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
	assert.Equal(t, "integer", columns[1].Type)
	assert.True(t, columns[1].NotNull)

	assert.Equal(t, "active", columns[2].Name)
	assert.Equal(t, "boolean", columns[2].Type)
	assert.True(t, columns[2].NotNull)

	assert.Equal(t, "salary", columns[3].Name)
	assert.Equal(t, "numeric(7,2)", columns[3].Type)
	assert.True(t, columns[3].NotNull)

	assert.Equal(t, "details", columns[4].Name)
	assert.Equal(t, "json", columns[4].Type)
	assert.False(t, columns[4].NotNull) // the second row contains empty value so not null is false
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
