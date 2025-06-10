package database

import (
	"bytes"
	"fluxend/internal/config/constants"
	"fluxend/internal/domain/database"
	"fluxend/pkg"
	"github.com/guregu/null/v6"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Common test cases for table name validation
var tableNameValidationTests = []struct {
	name     string
	payload  map[string]interface{}
	headers  map[string]string
	expected []string
}{
	{
		name:     "Missing project header",
		payload:  map[string]interface{}{"name": "test_name"},
		headers:  map[string]string{},
		expected: []string{"project"},
	},
	{
		name:    "Missing name",
		payload: map[string]interface{}{},
		headers: map[string]string{
			constants.ProjectHeaderKey: dummyProjectUUID,
		},
		expected: []string{"Name is required"},
	},
	{
		name: "Empty name",
		payload: map[string]interface{}{
			"name": "",
		},
		headers: map[string]string{
			constants.ProjectHeaderKey: dummyProjectUUID,
		},
		expected: []string{"Name is required"},
	},
	{
		name: "Invalid name - special characters",
		payload: map[string]interface{}{
			"name": "invalid@table",
		},
		headers: map[string]string{
			constants.ProjectHeaderKey: dummyProjectUUID,
		},
		expected: []string{"Table name must be alphanumeric with underscores"},
	},
	{
		name: "Name too short",
		payload: map[string]interface{}{
			"name": "a",
		},
		headers: map[string]string{
			constants.ProjectHeaderKey: dummyProjectUUID,
		},
		expected: []string{"Name must be between"},
	},
	{
		name: "Reserved table name",
		payload: map[string]interface{}{
			"name": "information_schema",
		},
		headers: map[string]string{
			constants.ProjectHeaderKey: dummyProjectUUID,
		},
		expected: []string{"table name 'information_schema' is reserved"},
	},
}

func createValidColumns() []database.Column {
	return []database.Column{
		{
			Name:    "valid_column",
			Type:    constants.ColumnTypeVarchar,
			NotNull: true,
		},
		{
			Name:   "another_column",
			Type:   constants.ColumnTypeInteger,
			Unique: true,
		},
	}
}

func createMultipartRequest(e *echo.Echo, name string, hasFile bool) echo.Context {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", name)

	if hasFile {
		part, _ := writer.CreateFormFile("file", "test.csv")
		_, _ = part.Write([]byte("col1,col2\nvalue1,value2"))
	}
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestCreateTableRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateTableRequest: valid", func(t *testing.T) {
		columns := createValidColumns()

		payload := map[string]interface{}{
			"name":    "valid_table",
			"columns": columns,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateTableRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Len(t, r.Columns, 2)
		assert.Equal(t, "valid_column", r.Columns[0].Name)
		assert.Equal(t, constants.ColumnTypeVarchar, r.Columns[0].Type)
	})

	t.Run("CreateTableRequest: valid with foreign key", func(t *testing.T) {
		columns := []database.Column{
			{
				Name:            "foreign_key_column",
				Type:            constants.ColumnTypeInteger,
				NotNull:         true,
				Foreign:         true,
				ReferenceTable:  null.StringFrom("users"),
				ReferenceColumn: null.StringFrom("id"),
			},
		}

		payload := map[string]interface{}{
			"name":    "table_with_fk",
			"columns": columns,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateTableRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Len(t, r.Columns, 1)
		assert.True(t, r.Columns[0].Foreign)
	})

	t.Run("CreateTableRequest: invalid", func(t *testing.T) {
		// Test table name validation using shared test cases
		for _, tc := range tableNameValidationTests {
			t.Run(tc.name, func(t *testing.T) {
				// Add columns to payload to make it a valid CreateTableRequest structure
				payload := make(map[string]interface{})
				for k, v := range tc.payload {
					payload[k] = v
				}
				payload["columns"] = createValidColumns()

				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r CreateTableRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}

		// Additional CreateTableRequest-specific validation tests
		additionalTests := []struct {
			name     string
			payload  map[string]interface{}
			expected []string
		}{
			{
				name: "Missing columns",
				payload: map[string]interface{}{
					"name": "test_table",
				},
				expected: []string{"Columns are required"},
			},
			{
				name: "Empty columns array",
				payload: map[string]interface{}{
					"name":    "test_table",
					"columns": []database.Column{},
				},
				expected: []string{"Columns are required"},
			},
			{
				name: "Invalid column name",
				payload: map[string]interface{}{
					"name": "test_table",
					"columns": []database.Column{
						{Name: "", Type: constants.ColumnTypeVarchar},
					},
				},
				expected: []string{"Column name is required"},
			},
			{
				name: "Invalid column type",
				payload: map[string]interface{}{
					"name": "test_table",
					"columns": []database.Column{
						{Name: "test_column", Type: "invalid_type"},
					},
				},
				expected: []string{"column type 'invalid_type' is not allowed"},
			},
		}

		for _, tc := range additionalTests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)
				ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

				var r CreateTableRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestRenameTableRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("RenameTableRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "new_table_name",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r RenameTableRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
	})

	t.Run("RenameTableRequest: invalid", func(t *testing.T) {
		// Use shared table name validation tests
		for _, tc := range tableNameValidationTests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r RenameTableRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestUploadTableRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("UploadTableRequest: valid", func(t *testing.T) {
		ctx := createMultipartRequest(e, "uploaded_table", true)

		var r UploadTableRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "uploaded_table", r.Name)
		assert.NotNil(t, r.File)
		assert.Equal(t, "test.csv", r.File.Filename)
	})

	t.Run("UploadTableRequest: invalid", func(t *testing.T) {
		t.Run("Missing file", func(t *testing.T) {
			ctx := createMultipartRequest(e, "test_table", false)

			var r UploadTableRequest
			errs := r.BindAndValidate(ctx)

			pkg.AssertErrorContains(t, errs, "File is required")
		})
	})
}
