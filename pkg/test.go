package pkg

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func CreateTestContext(t *testing.T, e *echo.Echo, method, target string, payload map[string]interface{}) echo.Context {
	body, err := json.Marshal(payload)
	assert.NoError(t, err, "Failed to marshal payload")

	fakeRequest := httptest.NewRequest(method, target, bytes.NewReader(body))
	fakeRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return e.NewContext(fakeRequest, httptest.NewRecorder())
}

func AssertErrorContains(t *testing.T, errs []string, expectedSubstring string) {
	t.Helper()

	for _, err := range errs {
		if strings.Contains(err, expectedSubstring) {
			return
		}
	}

	t.Errorf("Expected error containing '%s', but got errors: %v", expectedSubstring, errs)
}
