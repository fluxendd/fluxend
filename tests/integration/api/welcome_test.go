package api

import (
	"encoding/json"
	"fluxend/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestWelcomeEndpoint(t *testing.T) {
	testServer := integration.NewTestServer()
	defer testServer.Close()

	resp, err := testServer.Client.Get(testServer.BaseURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.True(t, strings.Contains(response["message"].(string), "Fluxend"))
}
