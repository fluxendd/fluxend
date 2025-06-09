package organization

import (
	"bytes"
	"encoding/json"
	"fluxend/pkg"
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
		payload := map[string]interface{}{
			"name": "Valid_Organization-123",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, payload["name"], r.Name)
	})

	t.Run("CreateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			expected []string
		}{
			{
				name:    "Missing: name",
				payload: map[string]interface{}{},
				expected: []string{
					"Name is required",
				},
			},
			{
				name: "Empty name",
				payload: map[string]interface{}{
					"name": "",
				},
				expected: []string{"Name is required"},
			},
			{
				name: "Invalid characters in name",
				payload: map[string]interface{}{
					"name": "!!!BAD$$$",
				},
				expected: []string{
					"Organization name must be alphanumeric with underscores, spaces and dashes",
				},
			},
			{
				name: "Too short name",
				payload: map[string]interface{}{
					"name": "A",
				},
				expected: []string{
					"Organization name be between",
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

func TestMemberCreateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("MemberCreateRequest: valid", func(t *testing.T) {
		validUUID := uuid.New()
		payload := map[string]interface{}{
			"user_id": validUUID.String(),
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r MemberCreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, validUUID, r.UserID)
	})

	t.Run("MemberCreateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			expected []string
		}{
			{
				name:    "Missing: user_id",
				payload: map[string]interface{}{},
				expected: []string{
					"UserID is required",
				},
			},
			{
				name: "Invalid UUID format",
				payload: map[string]interface{}{
					"user_id": "not-a-uuid",
				},
				expected: []string{"Invalid request payload"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				body, err := json.Marshal(tc.payload)
				assert.NoError(t, err, "Failed to marshal payload")

				fakeRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
				fakeRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				ctx := e.NewContext(fakeRequest, httptest.NewRecorder())

				var r MemberCreateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
