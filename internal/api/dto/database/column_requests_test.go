package database

import (
	"fluxend/internal/config/constants"
	"fluxend/internal/domain/database"
	"fluxend/pkg"
	"github.com/guregu/null/v6"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var dummyProjectUUID = "123e4567-e89b-12d3-a456-426614174000"

func TestCreateColumnRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateColumnRequest: valid", func(t *testing.T) {
		columns := []database.Column{
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

		payload := map[string]interface{}{
			"columns": columns,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateColumnRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Len(t, r.Columns, 2)
		assert.Equal(t, "valid_column", r.Columns[0].Name)
		assert.Equal(t, constants.ColumnTypeVarchar, r.Columns[0].Type)
	})

	t.Run("CreateColumnRequest: valid with foreign key", func(t *testing.T) {
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
			"columns": columns,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateColumnRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Len(t, r.Columns, 1)
		assert.Equal(t, "foreign_key_column", r.Columns[0].Name)
		assert.True(t, r.Columns[0].Foreign)
	})

	t.Run("CreateColumnRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing project header",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "test", Type: constants.ColumnTypeVarchar},
					},
				},
				headers:  map[string]string{},
				expected: []string{"project"},
			},
			{
				name: "Missing column name",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "", Type: constants.ColumnTypeVarchar},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column name is required"},
			},
			{
				name: "Missing column type",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "test_column", Type: ""},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column type is required"},
			},
			{
				name: "Invalid column name - too short",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "a", Type: constants.ColumnTypeVarchar},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column name be between"},
			},
			{
				name: "Invalid column name - special characters",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "invalid@name", Type: constants.ColumnTypeVarchar},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column name must be alphanumeric and start with a letter"},
			},
			{
				name: "Reserved column name",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "oid", Type: constants.ColumnTypeVarchar},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"column name 'oid' is reserved"},
			},
			{
				name: "Invalid column type",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{Name: "test_column", Type: "invalid_type"},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"column type 'invalid_type' is not allowed"},
			},
			{
				name: "Foreign key without reference table",
				payload: map[string]interface{}{
					"columns": []database.Column{
						{
							Name:    "foreign_column",
							Type:    constants.ColumnTypeInteger,
							Foreign: true,
						},
					},
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"reference table and column are required for foreign key constraints"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r CreateColumnRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestRenameColumnRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("RenameColumnRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "new_column_name",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r RenameColumnRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
	})

	t.Run("RenameColumnRequest: invalid", func(t *testing.T) {
		tests := []struct {
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
				expected: []string{"New name is required for column"},
			},
			{
				name: "Empty name",
				payload: map[string]interface{}{
					"name": "",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"New name is required for column"},
			},
			{
				name: "Name too short",
				payload: map[string]interface{}{
					"name": "a",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Column name be between"},
			},
			{
				name: "Reserved column name",
				payload: map[string]interface{}{
					"name": "created_at",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"column name 'created_at' is reserved"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r RenameColumnRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
