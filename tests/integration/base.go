package integration

import (
	"encoding/json"
	"fluxend/internal/app/commands"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"fluxend/internal/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer wraps the echo server for testing
type TestServer struct {
	Server  *httptest.Server
	EchoApp *echo.Echo
	BaseURL string
}

// NewTestServer creates a new test server instance
func NewTestServer() *TestServer {
	// Set up test environment variables
	setTestEnvVars()

	// Initialize container and setup server
	container := app.InitializeContainer()
	e := commands.SetupServer(container) // You'll need to export this function

	// Create test server
	server := httptest.NewServer(e)

	return &TestServer{
		Server:  server,
		EchoApp: e,
		BaseURL: server.URL,
	}
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	ts.Server.Close()
}

// setTestEnvVars sets up required environment variables for testing
func setTestEnvVars() {
	testEnvVars := map[string]string{
		"APP_ENV":                 "test",
		"BASE_URL":                "http://localhost:8080",
		"APP_URL":                 "http://localhost:3000",
		"API_URL":                 "http://localhost:8080",
		"APP_CONTAINER_NAME":      "test-app",
		"DATABASE_CONTAINER_NAME": "test-db",
		"FRONTEND_CONTAINER_NAME": "test-frontend",
		"VITE_FLX_INTERNAL_URL":   "http://localhost:8080",
		"VITE_FLX_API_URL":        "http://localhost:8080",
		"VITE_FLX_BASE_DOMAIN":    "localhost",
		"VITE_FLX_HTTP_SCHEME":    "http",

		"DATABASE_HOST":     "localhost",
		"DATABASE_USER":     "fluxend",
		"DATABASE_PASSWORD": "fluxend",
		"DATABASE_NAME":     "fluxend",
		"DATABASE_SSL_MODE": "disable",

		"JWT_SECRET":               "test_jwt_secret_key_that_is_long_enough_for_validation",
		"STORAGE_DRIVER":           "local",
		"POSTGREST_DB_HOST":        "localhost",
		"POSTGREST_DB_USER":        "test_user",
		"POSTGREST_DB_PASSWORD":    "test_password",
		"POSTGREST_DEFAULT_SCHEMA": "public",
		"POSTGREST_DEFAULT_ROLE":   "test_role",
		"CUSTOM_ORIGINS":           "http://localhost:3000,http://api.fluxend.localhost",
	}

	for key, value := range testEnvVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// TestHealthEndpoint demonstrates a simple integration test
func TestHealthEndpoint(t *testing.T) {
	// Start test server
	testServer := NewTestServer()
	defer testServer.Close()

	// Make HTTP request
	client := &http.Client{Timeout: 10 * time.Second}

	// Test health endpoint (assuming you have one)
	resp, err := client.Get(testServer.BaseURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse and print JSON response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Pretty print JSON
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	require.NoError(t, err)
	fmt.Printf("Health endpoint response:\n%s\n", string(jsonBytes))
}

// TestAPIWithCustomHost demonstrates testing with custom host header
func TestAPIWithCustomHost(t *testing.T) {
	testServer := NewTestServer()
	defer testServer.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	// Create request with custom host header
	req, err := http.NewRequest("GET", testServer.BaseURL+"/docs/", nil)
	require.NoError(t, err)

	// Set host header to simulate api.fluxend.localhost
	req.Host = "api.fluxend.localhost"
	req.Header.Set("Origin", "http://api.fluxend.localhost")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Headers: %+v\n", resp.Header)

	// If it's JSON, parse and print
	if resp.Header.Get("Content-Type") == "application/json" {
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
			jsonBytes, _ := json.MarshalIndent(response, "", "  ")
			fmt.Printf("JSON Response:\n%s\n", string(jsonBytes))
		}
	}
}

// TestUserEndpoint demonstrates testing a protected endpoint
func TestUserEndpoint(t *testing.T) {
	testServer := NewTestServer()
	defer testServer.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	// Test without authentication (should fail)
	resp, err := client.Get(testServer.BaseURL + "/users")
	require.NoError(t, err)
	defer resp.Body.Close()

	fmt.Printf("Unauthenticated request status: %d\n", resp.StatusCode)

	// Parse response if it's JSON
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
		jsonBytes, _ := json.MarshalIndent(response, "", "  ")
		fmt.Printf("Error response:\n%s\n", string(jsonBytes))
	}
}

// Helper function to create authenticated request
func createAuthenticatedRequest(method, url, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}

// TestIntegrationSuite runs multiple integration tests
func TestIntegrationSuite(t *testing.T) {
	testServer := NewTestServer()
	defer testServer.Close()

	tests := []struct {
		name     string
		endpoint string
		method   string
		headers  map[string]string
	}{
		{"Swagger Docs", "/docs/", "GET", nil},
		{"Health Check", "/health", "GET", nil},
		{"Users Endpoint", "/users", "GET", map[string]string{"Origin": "http://api.fluxend.localhost"}},
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, testServer.BaseURL+test.endpoint, nil)
			require.NoError(t, err)

			// Set custom headers
			for key, value := range test.headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			fmt.Printf("\n=== %s ===\n", test.name)
			fmt.Printf("URL: %s\n", testServer.BaseURL+test.endpoint)
			fmt.Printf("Status: %d\n", resp.StatusCode)

			// Try to parse as JSON
			var response interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
				jsonBytes, _ := json.MarshalIndent(response, "", "  ")
				fmt.Printf("Response:\n%s\n", string(jsonBytes))
			} else {
				fmt.Printf("Response is not JSON or empty\n")
			}
		})
	}
}
