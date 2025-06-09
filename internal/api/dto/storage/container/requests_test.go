package container

import (
	"fluxend/internal/config/constants"
	"fluxend/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

var dummyProjectUUID = "123e4567-e89b-12d3-a456-426614174000"

func TestCreateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":          "valid_container_name",
			"description":   "This is a valid container description",
			"is_public":     true,
			"max_file_size": 1024,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "valid_container_name", r.Name)
		assert.Equal(t, "This is a valid container description", r.Description)
		assert.True(t, r.IsPublic)
		assert.Equal(t, 1024, r.MaxFileSize)
	})

	t.Run("CreateRequest: valid with minimum values", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":          "abc", // minimum length
			"description":   "",    // minimum length (0)
			"is_public":     true,
			"max_file_size": 1, // minimum value
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "abc", r.Name)
		assert.Empty(t, r.Description)
		assert.True(t, r.IsPublic)
		assert.Equal(t, 1, r.MaxFileSize)
	})

	t.Run("CreateRequest: valid with maximum values", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":          strings.Repeat("a", constants.MaxContainerNameLength),        // maximum length
			"description":   strings.Repeat("d", constants.MaxContainerDescriptionLength), // maximum length
			"is_public":     true,
			"max_file_size": 999999999,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, strings.Repeat("a", constants.MaxContainerNameLength), r.Name)
		assert.Equal(t, strings.Repeat("d", constants.MaxContainerDescriptionLength), r.Description)
		assert.True(t, r.IsPublic)
		assert.Equal(t, 999999999, r.MaxFileSize)
	})

	t.Run("CreateRequest: valid with underscores and dashes", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":          "container_name-with-dashes_and_underscores123",
			"description":   "Valid container with special characters in name",
			"is_public":     false,
			"max_file_size": 512,
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "container_name-with-dashes_and_underscores123", r.Name)
		assert.Equal(t, "Valid container with special characters in name", r.Description)
		assert.False(t, r.IsPublic)
		assert.Equal(t, 512, r.MaxFileSize)
	})

	t.Run("CreateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing project header",
				payload: map[string]interface{}{
					"name":          "valid_name",
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers:  map[string]string{},
				expected: []string{"project"},
			},
			{
				name: "Missing name",
				payload: map[string]interface{}{
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Name is required"},
			},
			{
				name: "Empty name",
				payload: map[string]interface{}{
					"name":          "",
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Name is required"},
			},
			{
				name: "Name too short",
				payload: map[string]interface{}{
					"name":          "ab", // less than MinContainerNameLength (3)
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Container name must be between 3 and 63 characters"},
			},
			{
				name: "Name too long",
				payload: map[string]interface{}{
					"name":          strings.Repeat("a", constants.MaxContainerNameLength+1), // longer than MaxContainerNameLength (63)
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Container name must be between 3 and 63 characters"},
			},
			{
				name: "Name with invalid characters - spaces",
				payload: map[string]interface{}{
					"name":          "invalid name with spaces",
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Container name must be alphanumeric with underscores and dashes"},
			},
			{
				name: "Name with invalid characters - special symbols",
				payload: map[string]interface{}{
					"name":          "invalid@name!",
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Container name must be alphanumeric with underscores and dashes"},
			},
			{
				name: "Name with invalid characters - dots",
				payload: map[string]interface{}{
					"name":          "invalid.name",
					"description":   "valid description",
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Container name must be alphanumeric with underscores and dashes"},
			},
			{
				name: "Missing max_file_size",
				payload: map[string]interface{}{
					"name":        "valid_name",
					"description": "valid description",
					"is_public":   true,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"max_file_size is required"},
			},
			{
				name: "Description too long",
				payload: map[string]interface{}{
					"name":          "valid_name",
					"description":   strings.Repeat("d", constants.MaxContainerDescriptionLength+1),
					"is_public":     true,
					"max_file_size": 1024,
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Container description must be less than 255 characters"},
			},
			{
				name: "Multiple validation errors",
				payload: map[string]interface{}{
					"name":          "a",                                                            // too short
					"description":   strings.Repeat("d", constants.MaxContainerDescriptionLength+1), // too long
					"max_file_size": -5,                                                             // negative
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{
					"Container name must be between 3 and 63 characters",
					"max_file_size must be a positive number",
					"Container description must be less than 255 characters",
				},
			},
			{
				name: "Invalid JSON payload",
				payload: map[string]interface{}{
					"name": map[string]interface{}{"invalid": "structure"}, // invalid type for name
				},
				headers: map[string]string{
					constants.ProjectHeaderKey: dummyProjectUUID,
				},
				expected: []string{"Invalid request payload"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r CreateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
