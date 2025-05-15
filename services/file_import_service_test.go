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
