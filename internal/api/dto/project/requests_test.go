package project

import (
	"fluxend/pkg"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateRequest: valid", func(t *testing.T) {
		validUUID := uuid.New()
		payload := map[string]interface{}{
			"name":              "Valid_Project-123",
			"description":       "Some description",
			"organization_uuid": validUUID.String(),
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Equal(t, payload["description"], r.Description)
		assert.Equal(t, validUUID, r.OrganizationUUID)
	})

	t.Run("CreateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			expected []string
		}{
			{
				name: "Missing: name",
				payload: map[string]interface{}{
					"organization_uuid": uuid.New().String(),
				},
				expected: []string{
					"Name is required",
				},
			},
			{
				name: "Missing: Organization UUID",
				payload: map[string]interface{}{
					"name": "SomeName",
				},
				expected: []string{
					"Organization UUID is required",
				},
			},
			{
				name: "Empty name",
				payload: map[string]interface{}{
					"name":              "",
					"organization_uuid": uuid.New().String(),
				},
				expected: []string{"Name is required"},
			},
			{
				name: "Invalid UUID format",
				payload: map[string]interface{}{
					"name":              "ValidName",
					"organization_uuid": "not-a-uuid",
				},
				expected: []string{"Invalid request payload"},
			},
			{
				name: "Invalid characters in name",
				payload: map[string]interface{}{
					"name":              "!!!BAD$$$",
					"organization_uuid": uuid.New().String(),
				},
				expected: []string{
					"Project name must be alphanumeric with underscores, spaces and dashes",
				},
			},
			{
				name: "Too short name",
				payload: map[string]interface{}{
					"name":              "A",
					"organization_uuid": uuid.New().String(),
				},
				expected: []string{
					"Project name be between",
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				var r CreateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestUpdateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("UpdateRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Updated_Project-123",
			"description": "Updated desc",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
		assert.Equal(t, payload["description"], r.Description)
	})

	t.Run("UpdateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			expected []string
		}{
			{
				name: "Missing name",
				payload: map[string]interface{}{
					"description": "some desc",
				},
				expected: []string{
					"Name is required",
				},
			},
			{
				name: "Name too short",
				payload: map[string]interface{}{
					"name": "a",
				},
				expected: []string{
					"Name must be between",
				},
			},
			{
				name: "Invalid name characters",
				payload: map[string]interface{}{
					"name": "@@@@!!!!",
				},
				expected: []string{
					"Project name must be alphanumeric with underscores, spaces and dashes",
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, tc.payload)

				var r UpdateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
