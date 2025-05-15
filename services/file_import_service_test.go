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
