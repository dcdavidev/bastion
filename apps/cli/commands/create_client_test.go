package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClientCommand(t *testing.T) {
	// 1. Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/clients", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)
		assert.Equal(t, "TestClient", payload["name"])

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	// 2. Setup test config and token
	dir := setupTestConfig(t)
	mockToken(t, dir, "test-token")

	// 3. Execute command with flags
	_, err := executeCommand(rootCmd, "create", "client", "--name", "TestClient", "--url", server.URL)

	// 4. Verify results
	assert.NoError(t, err)
}
