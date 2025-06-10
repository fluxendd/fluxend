package database

import (
	"fluxend/internal/config/constants"
	"fluxend/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateFunctionRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateFunctionRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "test_function",
			"parameters":  []functionParameter{},
			"definition":  "BEGIN RETURN 1; END",
			"language":    "plpgsql",
			"return_type": "integer",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateFunctionRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Equal(t, payload["definition"], r.Definition)
		assert.Equal(t, payload["language"], r.Language)
		assert.Equal(t, payload["return_type"], r.ReturnType)
	})

	t.Run("CreateFunctionRequest: valid with parameters", func(t *testing.T) {
		parameters := []functionParameter{
			{Name: "param1", Type: constants.ColumnTypeInteger},
			{Name: "param2", Type: constants.ColumnTypeVarchar},
		}

		payload := map[string]interface{}{
			"name":        "test_function_with_params",
			"parameters":  parameters,
			"definition":  "BEGIN RETURN param1 + 1; END",
			"language":    "sql",
			"return_type": "bigint",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateFunctionRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Len(t, r.Parameters, 2)
		assert.Equal(t, "param1", r.Parameters[0].Name)
		assert.Equal(t, "integer", r.Parameters[0].Type)
	})

	t.Run("CreateFunctionRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing project header",
				payload: map[string]interface{}{
					"name":        "test_function",
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers:  map[string]string{},
				expected: []string{"project"},
			},
			{
				name: "Missing name",
				payload: map[string]interface{}{
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"name is required"},
			},
			{
				name: "Empty name",
				payload: map[string]interface{}{
					"name":        "",
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"name is required"},
			},
			{
				name: "Invalid name - special characters",
				payload: map[string]interface{}{
					"name":        "invalid@function",
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"name must be alphanumeric with underscores"},
			},
			{
				name: "Name too short",
				payload: map[string]interface{}{
					"name":        "a",
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"name must be between"},
			},
			{
				name: "Missing return type",
				payload: map[string]interface{}{
					"name":       "test_function",
					"definition": "BEGIN RETURN 1; END",
					"language":   "plpgsql",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"return_type is required"},
			},
			{
				name: "Invalid return type",
				payload: map[string]interface{}{
					"name":        "test_function",
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "invalid_type",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"invalid return type: invalid_type"},
			},
			{
				name: "Missing language",
				payload: map[string]interface{}{
					"name":        "test_function",
					"definition":  "BEGIN RETURN 1; END",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"language is required"},
			},
			{
				name: "Invalid language",
				payload: map[string]interface{}{
					"name":        "test_function",
					"definition":  "BEGIN RETURN 1; END",
					"language":    "invalid_lang",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"invalid language: invalid_lang"},
			},
			{
				name: "Missing definition",
				payload: map[string]interface{}{
					"name":        "test_function",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"definition is required"},
			},
			{
				name: "Invalid definition - missing BEGIN",
				payload: map[string]interface{}{
					"name":        "test_function",
					"definition":  "RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"invalid definition"},
			},
			{
				name: "Invalid definition - missing END",
				payload: map[string]interface{}{
					"name":        "test_function",
					"definition":  "BEGIN RETURN 1;",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"invalid definition"},
			},
			{
				name: "Invalid parameter type",
				payload: map[string]interface{}{
					"name": "test_function",
					"parameters": []functionParameter{
						{Name: "param1", Type: "invalid_type"},
					},
					"definition":  "BEGIN RETURN 1; END",
					"language":    "plpgsql",
					"return_type": "integer",
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"invalid parameter type: invalid_type"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r CreateFunctionRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
