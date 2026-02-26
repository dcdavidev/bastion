package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectCommands(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
		} else {
			json.NewEncoder(w).Encode([]models.Project{{ID: uuid.New(), Name: "P1"}})
		}
	})
	mux.HandleFunc("/api/v1/clients", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]models.Client{{ID: uuid.New(), Name: "C1"}})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	dir := setupTestConfig(t)
	mockToken(t, dir, "test-token")

	t.Run("CreateProject", func(t *testing.T) {
		_, err := executeCommand(rootCmd, "create", "project", "--client", uuid.New().String(), "--name", "TestProject", "--url", server.URL)
		assert.NoError(t, err)
	})

	t.Run("ListClients", func(t *testing.T) {
		_, err := executeCommand(rootCmd, "list", "clients", "--url", server.URL)
		assert.NoError(t, err)
	})

	t.Run("ListProjects", func(t *testing.T) {
		_, err := executeCommand(rootCmd, "list", "projects", "--client", uuid.New().String(), "--url", server.URL)
		assert.NoError(t, err)
	})
}

func TestLoginCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"token": "new-token"})
	}))
	defer server.Close()

	setupTestConfig(t)

	_, err := executeCommand(rootCmd, "login", "--email", "test@example.com", "--password", "password123", "--url", server.URL)
	assert.NoError(t, err)

	token, _ := loadToken()
	assert.Equal(t, "new-token", token)
}

func TestRunCommand(t *testing.T) {
	projectID := uuid.New()
	adminPassword := "vault-pass"
	salt := []byte("salt-16-bytes--")
	mk, _ := crypto.GenerateRandomKey()
	dk, _ := crypto.GenerateRandomKey()
	
	adminKEK := crypto.DeriveKey([]byte(adminPassword), salt)
	wrappedMK, _ := crypto.WrapKey(adminKEK, mk)
	wrappedDK, _ := crypto.WrapKey(mk, dk)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault/config", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"wrapped_master_key": hex.EncodeToString(wrappedMK),
			"master_key_salt":    hex.EncodeToString(salt),
		})
	})
	mux.HandleFunc(fmt.Sprintf("/api/v1/projects/%s/key", projectID), func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"wrapped_data_key": hex.EncodeToString(wrappedDK),
		})
	})
	mux.HandleFunc("/api/v1/secrets", func(w http.ResponseWriter, r *http.Request) {
		val, _ := crypto.Encrypt(dk, []byte("injected-secret"))
		secrets := []models.Secret{
			{Key: "TEST_SECRET", Value: hex.EncodeToString(val)},
		}
		json.NewEncoder(w).Encode(secrets)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	dir := setupTestConfig(t)
	mockToken(t, dir, "test-token")

	// Run "env" command and check if TEST_SECRET is present in output
	output, err := executeCommand(rootCmd, "run", "--project", projectID.String(), "--password", adminPassword, "--url", server.URL, "--", "env")
	assert.NoError(t, err)
	assert.Contains(t, output, "TEST_SECRET=injected-secret")
}

func TestSetCommand(t *testing.T) {
	projectID := uuid.New()
	adminPassword := "vault-pass"
	salt := []byte("salt-16-bytes--")
	mk, _ := crypto.GenerateRandomKey()
	dk, _ := crypto.GenerateRandomKey()
	
	adminKEK := crypto.DeriveKey([]byte(adminPassword), salt)
	wrappedMK, _ := crypto.WrapKey(adminKEK, mk)
	wrappedDK, _ := crypto.WrapKey(mk, dk)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault/config", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"wrapped_master_key": hex.EncodeToString(wrappedMK),
			"master_key_salt":    hex.EncodeToString(salt),
		})
	})
	mux.HandleFunc(fmt.Sprintf("/api/v1/projects/%s", projectID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":"%s","name":"Test","wrapped_data_key":"%s"}`, projectID, hex.EncodeToString(wrappedDK))
	})
	mux.HandleFunc("/api/v1/secrets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	dir := setupTestConfig(t)
	mockToken(t, dir, "test-token")

	_, err := executeCommand(rootCmd, "set", "MY_SECRET", "secret-value", "--project", projectID.String(), "--password", adminPassword, "--url", server.URL)
	assert.NoError(t, err)
}
