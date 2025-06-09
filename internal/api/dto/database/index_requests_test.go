package database

import (
	"fluxend/internal/config/constants"
	"fluxend/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateIndexRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateIndexRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":      "test_index",
			"columns":   []string{"column1", "column2"},
			"is_unique": false,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateIndexRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Equal(t, []string{"column1", "column2"}, r.Columns)
		assert.Equal(t, false, r.IsUnique)
	})

	t.Run("CreateIndexRequest: valid unique index", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":      "unique_test_index",
			"columns":   []string{"email"},
			"is_unique": true,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateIndexRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Equal(t, []string{"email"}, r.Columns)
		assert.Equal(t, true, r.IsUnique)
	})

	t.Run("CreateIndexRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing project header",
				payload: map[string]interface{}{
					"name":    "test_index",
					"columns": []string{"column1"},
				},
				headers:  map[string]string{},
				expected: []string{"project"},
			},
			{
				name: "Missing name",
				payload: map[string]interface{}{
					"columns": []string{"column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Index name is required"},
			},
			{
				name: "Empty name",
				payload: map[string]interface{}{
					"name":    "",
					"columns": []string{"column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Index name is required"},
			},
			{
				name: "Invalid name - special characters",
				payload: map[string]interface{}{
					"name":    "invalid@index",
					"columns": []string{"column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Index name must be alphanumeric with underscores"},
			},
			{
				name: "Name too short",
				payload: map[string]interface{}{
					"name":    "a",
					"columns": []string{"column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Index name be between"},
			},
			{
				name: "Missing columns",
				payload: map[string]interface{}{
					"name": "test_index",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"At least one column is required"},
			},
			{
				name: "Empty columns array",
				payload: map[string]interface{}{
					"name":    "test_index",
					"columns": []string{},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"At least one column is required"},
			},
			{
				name: "Reserved index name",
				payload: map[string]interface{}{
					"name":    "primary",
					"columns": []string{"column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Index name 'primary' is reserved"},
			},
			{
				name: "Empty column name in array",
				payload: map[string]interface{}{
					"name":    "test_index",
					"columns": []string{"column1", "", "column3"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column name in index cannot be empty"},
			},
			{
				name: "Whitespace only column name",
				payload: map[string]interface{}{
					"name":    "test_index",
					"columns": []string{"column1", "   ", "column3"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column name in index cannot be empty"},
			},
			{
				name: "Duplicate column names",
				payload: map[string]interface{}{
					"name":    "test_index",
					"columns": []string{"column1", "column2", "column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Duplicate column 'column1' in index definition"},
			},
			{
				name: "Duplicate column names - case insensitive",
				payload: map[string]interface{}{
					"name":    "test_index",
					"columns": []string{"Column1", "column1"},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Duplicate column 'column1' in index definition"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r CreateIndexRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
