package file

import (
	"fluxend/internal/config/constants"
	"fluxend/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

var dummyProjectUUID = "123e4567-e89b-12d3-a456-426614174000"

func TestCreateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateRequest: valid", func(t *testing.T) {
		// Create a mock file header
		fileHeader := &multipart.FileHeader{
			Filename: "test.txt",
			Size:     1024,
		}

		payload := map[string]interface{}{
			"full_file_name": "valid_file_name.txt",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		r.FullFileName = "valid_file_name.txt"
		r.File = fileHeader

		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "valid_file_name.txt", r.FullFileName)
		assert.NotNil(t, r.File)
		assert.Equal(t, "test.txt", r.File.Filename)
	})

	t.Run("CreateRequest: valid with minimum filename length", func(t *testing.T) {
		fileHeader := &multipart.FileHeader{
			Filename: "a.txt",
			Size:     512,
		}

		payload := map[string]interface{}{
			"full_file_name": "abc",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		r.FullFileName = "abc"
		r.File = fileHeader

		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "abc", r.FullFileName)
		assert.NotNil(t, r.File)
	})

	t.Run("CreateRequest: valid with maximum filename length", func(t *testing.T) {
		fileHeader := &multipart.FileHeader{
			Filename: "document.pdf",
			Size:     2048,
		}

		maxLengthName := strings.Repeat("a", constants.MaxContainerNameLength)
		payload := map[string]interface{}{
			"full_file_name": maxLengthName,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		r.FullFileName = maxLengthName
		r.File = fileHeader

		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, maxLengthName, r.FullFileName)
		assert.NotNil(t, r.File)
	})

	t.Run("CreateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			file     *multipart.FileHeader
			expected []string
		}{
			{
				name: "Missing project header",
				payload: map[string]interface{}{
					"full_file_name": "valid_file.txt",
				},
				headers: map[string]string{},
				file: &multipart.FileHeader{
					Filename: "test.txt",
					Size:     1024,
				},
				expected: []string{"project"},
			},
			{
				name:    "Missing full_file_name",
				payload: map[string]interface{}{},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				file: &multipart.FileHeader{
					Filename: "test.txt",
					Size:     1024,
				},
				expected: []string{"full_file_name is required"},
			},
			{
				name: "Empty full_file_name",
				payload: map[string]interface{}{
					"full_file_name": "",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				file: &multipart.FileHeader{
					Filename: "test.txt",
					Size:     1024,
				},
				expected: []string{"full_file_name is required"},
			},
			{
				name: "Full file name too short",
				payload: map[string]interface{}{
					"full_file_name": "ab",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				file: &multipart.FileHeader{
					Filename: "test.txt",
					Size:     1024,
				},
				expected: []string{"File name must be between 3 and 63 characters"},
			},
			{
				name: "Full file name too long",
				payload: map[string]interface{}{
					"full_file_name": strings.Repeat("a", constants.MaxContainerNameLength+1),
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				file: &multipart.FileHeader{
					Filename: "test.txt",
					Size:     1024,
				},
				expected: []string{"File name must be between 3 and 63 characters"},
			},
			{
				name: "Missing file",
				payload: map[string]interface{}{
					"full_file_name": "valid_file.txt",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				file:     nil,
				expected: []string{"file is required"},
			},
			{
				name: "Multiple validation errors",
				payload: map[string]interface{}{
					"full_file_name": "a", // too short
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				file: nil, // missing file
				expected: []string{
					"File name must be between 3 and 63 characters",
					"file is required",
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r CreateRequest
				r.File = tc.file
				if tc.payload["full_file_name"] != nil {
					if filename, ok := tc.payload["full_file_name"].(string); ok {
						r.FullFileName = filename
					}
				}

				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestRenameRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("RenameRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"full_file_name": "renamed_file.txt",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r RenameRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "renamed_file.txt", r.FullFileName)
	})

	t.Run("RenameRequest: valid with minimum filename length", func(t *testing.T) {
		payload := map[string]interface{}{
			"full_file_name": "abc", // minimum length (3 characters)
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r RenameRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "abc", r.FullFileName)
	})

	t.Run("RenameRequest: valid with maximum filename length", func(t *testing.T) {
		payload := map[string]interface{}{
			"full_file_name": strings.Repeat("a", constants.MaxContainerNameLength),
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r RenameRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, strings.Repeat("a", constants.MaxContainerNameLength), r.FullFileName)
	})

	t.Run("RenameRequest: valid with special characters", func(t *testing.T) {
		payload := map[string]interface{}{
			"full_file_name": "my_file-v2.backup.tar.gz",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r RenameRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "my_file-v2.backup.tar.gz", r.FullFileName)
	})

	t.Run("RenameRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing project header",
				payload: map[string]interface{}{
					"full_file_name": "valid_file.txt",
				},
				headers:  map[string]string{},
				expected: []string{"project"},
			},
			{
				name:    "Missing full_file_name",
				payload: map[string]interface{}{},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Name is required"},
			},
			{
				name: "Empty full_file_name",
				payload: map[string]interface{}{
					"full_file_name": "",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Name is required"},
			},
			{
				name: "Full file name too short",
				payload: map[string]interface{}{
					"full_file_name": "ab",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"File name must be between 3 and 63 characters"},
			},
			{
				name: "Full file name too long",
				payload: map[string]interface{}{
					"full_file_name": strings.Repeat("a", constants.MaxContainerNameLength+1),
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"File name must be between 3 and 63 characters"},
			},
			{
				name: "Invalid JSON payload",
				payload: map[string]interface{}{
					"full_file_name": map[string]interface{}{"invalid": "structure"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Invalid request payload"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r RenameRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestFileRequired(t *testing.T) {
	t.Run("fileRequired: valid file", func(t *testing.T) {
		fileHeader := &multipart.FileHeader{
			Filename: "test.txt",
			Size:     1024,
		}

		err := fileRequired(fileHeader)
		assert.NoError(t, err)
	})

	t.Run("fileRequired: nil file", func(t *testing.T) {
		err := fileRequired(nil)
		assert.Error(t, err)
		assert.Equal(t, "file is required", err.Error())
	})

	t.Run("fileRequired: nil pointer to FileHeader", func(t *testing.T) {
		var fileHeader *multipart.FileHeader = nil
		err := fileRequired(fileHeader)
		assert.Error(t, err)
		assert.Equal(t, "file is required", err.Error())
	})

	t.Run("fileRequired: wrong type", func(t *testing.T) {
		err := fileRequired("not a file header")
		assert.Error(t, err)
		assert.Equal(t, "file is required", err.Error())
	})

	t.Run("fileRequired: wrong type - int", func(t *testing.T) {
		err := fileRequired(123)
		assert.Error(t, err)
		assert.Equal(t, "file is required", err.Error())
	})
}
