package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clean environment variables for this test
	os.Unsetenv("BASTION_HOST")
	os.Unsetenv("BASTION_PORT")

	cfg, err := LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8287", cfg.URL)
	assert.Equal(t, "8287", cfg.Port)
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Setup environment
	os.Setenv("BASTION_HOST", "https://api.bastion.cloud")
	os.Setenv("BASTION_PORT", "9090")
	defer os.Unsetenv("BASTION_HOST")
	defer os.Unsetenv("BASTION_PORT")

	cfg, err := LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "https://api.bastion.cloud", cfg.URL)
	assert.Equal(t, "9090", cfg.Port)
}

func TestLoadConfig_MissingProtocol(t *testing.T) {
	os.Setenv("BASTION_HOST", "localhost:8287")
	defer os.Unsetenv("BASTION_HOST")

	cfg, err := LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8287", cfg.URL)
}

func TestConfigStringRepresentation(t *testing.T) {
	cfg := &Config{
		URL:  "http://bastion.test",
		Port: "1234",
	}

	assert.Equal(t, "Config{URL: http://bastion.test, Port: 1234}", cfg.String())
}
