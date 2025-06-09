package setting

import (
	"fluxend/internal/config/constants"
	"fluxend/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const dummyProjectUUID = "550e8400-e29b-41d4-a716-446655440000"

// Common test cases for individual setting validation
var individualSettingValidationTests = []struct {
	name     string
	setting  IndividualSetting
	expected []string
}{
	{
		name:     "Missing name",
		setting:  IndividualSetting{Name: "", Value: "test_value"},
		expected: []string{"Name is required"},
	},
	{
		name:     "Missing value",
		setting:  IndividualSetting{Name: "test_name", Value: ""},
		expected: []string{"Value is required"},
	},
	{
		name:     "Missing both name and value",
		setting:  IndividualSetting{Name: "", Value: ""},
		expected: []string{"Name is required", "Value is required"},
	},
}

func createValidSettings() []IndividualSetting {
	return []IndividualSetting{
		{
			Name:  "theme",
			Value: "dark",
		},
		{
			Name:  "language",
			Value: "en",
		},
		{
			Name:  "notifications",
			Value: "enabled",
		},
	}
}

func createRequestPayload(settings []IndividualSetting) map[string]interface{} {
	return map[string]interface{}{
		"settings": settings,
	}
}

func TestUpdateRequest_BindAndValidate_Suite(t *testing.T) {
	e := echo.New()

	t.Run("UpdateRequest: valid", func(t *testing.T) {
		settings := createValidSettings()
		payload := createRequestPayload(settings)

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Len(t, r.Settings, 3)
		assert.Equal(t, "theme", r.Settings[0].Name)
		assert.Equal(t, "dark", r.Settings[0].Value)
		assert.Equal(t, "language", r.Settings[1].Name)
		assert.Equal(t, "en", r.Settings[1].Value)
	})

	t.Run("UpdateRequest: valid single setting", func(t *testing.T) {
		settings := []IndividualSetting{
			{Name: "single_setting", Value: "single_value"},
		}
		payload := createRequestPayload(settings)

		ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
		ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

		var r UpdateRequest
		errs := r.BindAndValidate(ctx)

		assert.Len(t, errs, 0)
		assert.Len(t, r.Settings, 1)
		assert.Equal(t, "single_setting", r.Settings[0].Name)
		assert.Equal(t, "single_value", r.Settings[0].Value)
	})

	t.Run("UpdateRequest: invalid", func(t *testing.T) {
		// Test missing settings
		t.Run("Missing settings", func(t *testing.T) {
			payload := map[string]interface{}{}

			ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
			ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

			var r UpdateRequest
			errs := r.BindAndValidate(ctx)

			pkg.AssertErrorContains(t, errs, "Settings required")
		})

		t.Run("Empty settings array", func(t *testing.T) {
			payload := createRequestPayload([]IndividualSetting{})

			ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
			ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

			var r UpdateRequest
			errs := r.BindAndValidate(ctx)

			pkg.AssertErrorContains(t, errs, "Settings required")
		})

		// Test individual setting validation using shared test cases
		for _, tc := range individualSettingValidationTests {
			t.Run(tc.name, func(t *testing.T) {
				settings := []IndividualSetting{tc.setting}
				payload := createRequestPayload(settings)

				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
				ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

				var r UpdateRequest
				errs := r.BindAndValidate(ctx)

				for _, expected := range tc.expected {
					pkg.AssertErrorContains(t, errs, expected)
				}
			})
		}

		// Test multiple invalid settings
		t.Run("Multiple invalid settings", func(t *testing.T) {
			settings := []IndividualSetting{
				{Name: "", Value: "value1"},                // Missing name
				{Name: "name2", Value: ""},                 // Missing value
				{Name: "", Value: ""},                      // Missing both
				{Name: "valid_name", Value: "valid_value"}, // Valid setting
			}
			payload := createRequestPayload(settings)

			ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
			ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

			var r UpdateRequest
			errs := r.BindAndValidate(ctx)

			// Should have errors for the first 3 settings
			pkg.AssertErrorContains(t, errs, "Setting[0]")
			pkg.AssertErrorContains(t, errs, "Setting[1]")
			pkg.AssertErrorContains(t, errs, "Setting[2]")
			pkg.AssertErrorContains(t, errs, "Name is required")
			pkg.AssertErrorContains(t, errs, "Value is required")
		})

		// Test invalid JSON payload
		t.Run("Invalid JSON payload", func(t *testing.T) {
			// This would typically be tested by sending malformed JSON,
			// but since we're using pkg.CreateFakeRequestContext, we'll simulate
			// the error by testing the behavior when bind fails

			// Create a context with invalid structure that would cause bind to fail
			invalidPayload := map[string]interface{}{
				"settings": "invalid_not_array",
			}

			ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, invalidPayload)
			ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

			var r UpdateRequest
			errs := r.BindAndValidate(ctx)

			// This should result in a bind error or empty settings
			assert.True(t, len(errs) > 0)
		})
	})
}

func TestIndividualSetting_Validation(t *testing.T) {
	// Test individual setting validation scenarios in isolation
	t.Run("IndividualSetting validation scenarios", func(t *testing.T) {
		e := echo.New()

		validationTests := []struct {
			name     string
			settings []IndividualSetting
			expected map[string][]string // map of setting index to expected errors
		}{
			{
				name: "All valid settings",
				settings: []IndividualSetting{
					{Name: "setting1", Value: "value1"},
					{Name: "setting2", Value: "value2"},
				},
				expected: map[string][]string{},
			},
			{
				name: "First setting invalid name",
				settings: []IndividualSetting{
					{Name: "", Value: "value1"},
					{Name: "setting2", Value: "value2"},
				},
				expected: map[string][]string{
					"0": {"Setting[0] - name: Name is required"},
				},
			},
			{
				name: "Second setting invalid value",
				settings: []IndividualSetting{
					{Name: "setting1", Value: "value1"},
					{Name: "setting2", Value: ""},
				},
				expected: map[string][]string{
					"1": {"Setting[1] - value: Value is required"},
				},
			},
			{
				name: "Multiple settings with different validation errors",
				settings: []IndividualSetting{
					{Name: "", Value: ""},           // Both invalid
					{Name: "valid", Value: "valid"}, // Valid
					{Name: "", Value: "value3"},     // Name invalid
					{Name: "name4", Value: ""},      // Value invalid
				},
				expected: map[string][]string{
					"0": {"Setting[0] - name: Name is required", "Setting[0] - value: Value is required"},
					"2": {"Setting[2] - name: Name is required"},
					"3": {"Setting[3] - value: Value is required"},
				},
			},
		}

		for _, tc := range validationTests {
			t.Run(tc.name, func(t *testing.T) {
				payload := createRequestPayload(tc.settings)

				ctx := pkg.CreateFakeRequestContext(t, e, http.MethodPut, payload)
				ctx.Request().Header.Set(constants.ProjectHeaderKey, dummyProjectUUID)

				var r UpdateRequest
				errs := r.BindAndValidate(ctx)

				if len(tc.expected) == 0 {
					assert.Len(t, errs, 0)
				} else {
					for _, expectedErrors := range tc.expected {
						for _, expectedError := range expectedErrors {
							pkg.AssertErrorContains(t, errs, expectedError)
						}
					}
				}
			})
		}
	})
}
