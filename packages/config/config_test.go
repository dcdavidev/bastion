package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Ensure config file doesn't exist to trigger defaults
	home, _ := os.UserHomeDir()
	configPath := home + "/.bastion/config.yaml"
	os.Remove(configPath)

	cfg, err := LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "default", cfg.ActiveProfile)
	assert.Contains(t, cfg.Profiles, "default")
	assert.Equal(t, "http://localhost:8287", cfg.Profiles["default"].URL)
}

func TestConfig_GetActiveProfile(t *testing.T) {
	cfg := &Config{
		ActiveProfile: "test",
		Profiles: map[string]Profile{
			"test": {
				Name: "test",
				URL:  "http://test.bastion",
			},
		},
	}

	profile := cfg.GetActiveProfile()
	require.NotNil(t, profile)
	assert.Equal(t, "test", profile.Name)
	assert.Equal(t, "http://test.bastion", profile.URL)
}

func TestConfig_GetActiveProfile_Fallback(t *testing.T) {
	cfg := &Config{
		ActiveProfile: "non-existent",
		Profiles: map[string]Profile{
			"only": {
				Name: "only",
				URL:  "http://only.bastion",
			},
		},
	}

	profile := cfg.GetActiveProfile()
	require.NotNil(t, profile)
	assert.Equal(t, "only", profile.Name)
}
