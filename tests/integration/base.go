// tests/integration/base.go - Updated with shared methods
package integration

import (
	"bytes"
	"encoding/json"
	"fluxend/internal/app/commands"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"fluxend/internal/app"
	"fluxend/internal/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors,omitempty"`
}

type TestServer struct {
	Server       *httptest.Server
	EchoApp      *echo.Echo
	DB           *sqlx.DB
	Client       *http.Client
	BaseURL      string
	cleanupFuncs []func() error
}

func NewTestServer() *TestServer {
	setTestEnvVars()

	// Initialize database
	database.InitDB()
	db := database.GetDB()

	container := app.InitializeContainer()
	e := commands.SetupServer(container)

	server := httptest.NewServer(e)

	return &TestServer{
		Server:       server,
		EchoApp:      e,
		DB:           db,
		Client:       &http.Client{Timeout: 10 * time.Second},
		BaseURL:      server.URL,
		cleanupFuncs: make([]func() error, 0),
	}
}

func (ts *TestServer) Close() {
	// Run cleanup functions in reverse order
	for i := len(ts.cleanupFuncs) - 1; i >= 0; i-- {
		ts.cleanupFuncs[i]()
	}
	ts.Server.Close()
}

// AddCleanup adds a cleanup function to be run when the server closes
func (ts *TestServer) AddCleanup(fn func() error) {
	ts.cleanupFuncs = append(ts.cleanupFuncs, fn)
}

// PostJSON sends a POST request with JSON body
func (ts *TestServer) PostJSON(endpoint string, data interface{}) *http.Response {
	jsonData, _ := json.Marshal(data)
	resp, _ := ts.Client.Post(ts.BaseURL+endpoint, "application/json", bytes.NewBuffer(jsonData))
	return resp
}

func (ts *TestServer) PutJSON(endpoint string, data interface{}) *http.Response {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", ts.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := ts.Client.Do(req)
	return resp
}

// GetWithAuth sends a GET request with Authorization header
func (ts *TestServer) GetWithAuth(endpoint, token string) *http.Response {
	req, _ := http.NewRequest("GET", ts.BaseURL+endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := ts.Client.Do(req)
	return resp
}

// PutJSONWithAuth sends a PUT request with JSON body and Authorization header
func (ts *TestServer) PutJSONWithAuth(endpoint, token string, data interface{}) *http.Response {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", ts.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := ts.Client.Do(req)
	return resp
}

// PostWithAuth sends a POST request with Authorization header
func (ts *TestServer) PostWithAuth(endpoint, token string) *http.Response {
	req, _ := http.NewRequest("POST", ts.BaseURL+endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := ts.Client.Do(req)
	return resp
}

// CleanupUser removes a user and related data from the database
func (ts *TestServer) CleanupUser(userUUID uuid.UUID) error {
	// Clean up JWT versions
	_, err := ts.DB.Exec("DELETE FROM authentication.jwt_versions WHERE user_id = $1", userUUID)
	if err != nil {
		fmt.Printf("Warning: Failed to cleanup JWT versions: %v\n", err)
	}

	// Clean up organizations (if any)
	_, err = ts.DB.Exec("DELETE FROM organizations WHERE created_by = $1", userUUID)
	if err != nil {
		fmt.Printf("Warning: Failed to cleanup organizations: %v\n", err)
	}

	// Clean up user
	_, err = ts.DB.Exec("DELETE FROM authentication.users WHERE uuid = $1", userUUID)
	if err != nil {
		fmt.Printf("Warning: Failed to cleanup user: %v\n", err)
	}

	return nil
}

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
