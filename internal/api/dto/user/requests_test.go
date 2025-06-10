package user

import (
	"fluxend/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestCreateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("CreateRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"username": "validuser123",
			"email":    "test@example.com",
			"password": "password123",
			"bio":      "This is a valid bio",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "validuser123", r.Username)
		assert.Equal(t, "test@example.com", r.Email)
		assert.Equal(t, "password123", r.Password)
		assert.Equal(t, "This is a valid bio", r.Bio)
	})

	t.Run("CreateRequest: valid without bio", func(t *testing.T) {
		payload := map[string]interface{}{
			"username": "user_with_dash-and_underscore",
			"email":    "user@domain.co.uk",
			"password": "12345",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r CreateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "user_with_dash-and_underscore", r.Username)
		assert.Equal(t, "user@domain.co.uk", r.Email)
		assert.Equal(t, "12345", r.Password)
		assert.Empty(t, r.Bio)
	})

	t.Run("CreateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing username",
				payload: map[string]interface{}{
					"email":    "test@example.com",
					"password": "password123",
				},
				expected: []string{"Username is required"},
			},
			{
				name: "Empty username",
				payload: map[string]interface{}{
					"username": "",
					"email":    "test@example.com",
					"password": "password123",
				},
				expected: []string{"Username is required"},
			},
			{
				name: "Username too short",
				payload: map[string]interface{}{
					"username": "ab",
					"email":    "test@example.com",
					"password": "password123",
				},
				expected: []string{"Username must be between 3 and 100 characters"},
			},
			{
				name: "Username too long",
				payload: map[string]interface{}{
					"username": strings.Repeat("a", 101),
					"email":    "test@example.com",
					"password": "password123",
				},
				expected: []string{"Username must be between 3 and 100 characters"},
			},
			{
				name: "Username with spaces",
				payload: map[string]interface{}{
					"username": "user name",
					"email":    "test@example.com",
					"password": "password123",
				},
				expected: []string{"Username must not contain spaces or special characters"},
			},
			{
				name: "Username with special characters",
				payload: map[string]interface{}{
					"username": "user@name!",
					"email":    "test@example.com",
					"password": "password123",
				},
				expected: []string{"Username must not contain spaces or special characters"},
			},
			{
				name: "Missing email",
				payload: map[string]interface{}{
					"username": "testuser",
					"password": "password123",
				},
				expected: []string{"Email is required"},
			},
			{
				name: "Empty email",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "",
					"password": "password123",
				},
				expected: []string{"Email is required"},
			},
			{
				name: "Invalid email format",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "invalid-email",
					"password": "password123",
				},
				expected: []string{"Email must be a valid email address"},
			},
			{
				name: "Invalid email format - missing domain",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "test@",
					"password": "password123",
				},
				expected: []string{"Email must be a valid email address"},
			},
			{
				name: "Missing password",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
				},
				expected: []string{"Password is required"},
			},
			{
				name: "Empty password",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "",
				},
				expected: []string{"Password is required"},
			},
			{
				name: "Password too short",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "1234",
				},
				expected: []string{"Password must be at least 5 characters"},
			},
			{
				name: "Bio too long",
				payload: map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "password123",
					"bio":      strings.Repeat("a", 501),
				},
				expected: []string{"Bio must be between 0 and 500 characters"},
			},
			{
				name: "Multiple validation errors",
				payload: map[string]interface{}{
					"username": "a",
					"email":    "invalid-email",
					"password": "123",
					"bio":      strings.Repeat("a", 501),
				},
				expected: []string{
					"Username must be between 3 and 100 characters",
					"Email must be a valid email address",
					"Password must be at least 5 characters",
					"Bio must be between 0 and 500 characters",
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
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}

func TestLoginRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("LoginRequest: valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r LoginRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "test@example.com", r.Email)
		assert.Equal(t, "password123", r.Password)
	})

	t.Run("LoginRequest: valid with minimum password length", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "user@domain.org",
			"password": "12345",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, payload)

		var r LoginRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "user@domain.org", r.Email)
		assert.Equal(t, "12345", r.Password)
	})

	t.Run("LoginRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Missing email",
				payload: map[string]interface{}{
					"password": "password123",
				},
				expected: []string{"Email is required"},
			},
			{
				name: "Empty email",
				payload: map[string]interface{}{
					"email":    "",
					"password": "password123",
				},
				expected: []string{"Email is required"},
			},
			{
				name: "Invalid email format",
				payload: map[string]interface{}{
					"email":    "not-an-email",
					"password": "password123",
				},
				expected: []string{"Email must be a valid email address"},
			},
			{
				name: "Missing password",
				payload: map[string]interface{}{
					"email": "test@example.com",
				},
				expected: []string{"Password is required"},
			},
			{
				name: "Empty password",
				payload: map[string]interface{}{
					"email":    "test@example.com",
					"password": "",
				},
				expected: []string{"Password is required"},
			},
			{
				name: "Password too short",
				payload: map[string]interface{}{
					"email":    "test@example.com",
					"password": "1234",
				},
				expected: []string{"Password must be at least 5 characters"},
			},
			{
				name: "Multiple validation errors",
				payload: map[string]interface{}{
					"email":    "invalid-email",
					"password": "123",
				},
				expected: []string{
					"Email must be a valid email address",
					"Password must be at least 5 characters",
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPost, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r LoginRequest
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

	t.Run("UpdateRequest: valid with bio", func(t *testing.T) {
		payload := map[string]interface{}{
			"bio": "This is my updated bio",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, "This is my updated bio", r.Bio)
	})

	t.Run("UpdateRequest: valid without bio", func(t *testing.T) {
		payload := map[string]interface{}{}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Empty(t, r.Bio)
	})

	t.Run("UpdateRequest: valid with empty bio", func(t *testing.T) {
		payload := map[string]interface{}{
			"bio": "",
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Empty(t, r.Bio)
	})

	t.Run("UpdateRequest: valid with maximum bio length", func(t *testing.T) {
		payload := map[string]interface{}{
			"bio": strings.Repeat("a", 500),
		}

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Equal(t, strings.Repeat("a", 500), r.Bio)
	})

	t.Run("UpdateRequest: invalid", func(t *testing.T) {
		tests := []struct {
			name     string
			payload  map[string]interface{}
			headers  map[string]string
			expected []string
		}{
			{
				name: "Bio too long",
				payload: map[string]interface{}{
					"bio": strings.Repeat("a", 501),
				},
				expected: []string{"Bio must be between 0 and 500 characters"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, tc.payload)

				for key, value := range tc.headers {
					ctx.Request().Header.Set(key, value)
				}

				var r UpdateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}
	})
}
