package project

import (
	"bytes"
	"encoding/json"
	"fluxton/pkg"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

		body, _ := json.Marshal(payload)

		fakeRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		fakeRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := e.NewContext(fakeRequest, httptest.NewRecorder())

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
				body, _ := json.Marshal(tc.payload)

				fakeRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
				fakeRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				ctx := e.NewContext(fakeRequest, httptest.NewRecorder())

				var r CreateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					found := false
					pkg.DumpJSON(errs)
					for _, err := range errs {
						if contains(err, expected) {
							found = true
							break
						}
					}

					assert.True(t, found, "expected error to contain: "+expected)
				}
			})
		}
	})
}

// contains is a loose string matcher to account for dynamic error text (like length ranges).
func contains(s, substr string) bool {
	return s != "" && substr != "" && bytes.Contains([]byte(s), []byte(substr))
}
