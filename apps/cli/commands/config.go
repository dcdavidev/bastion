package commands

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabaseURL       string `yaml:"database_url"`
	AdminPasswordHash string `yaml:"admin_password_hash"`
	AdminPasswordSalt string `yaml:"admin_password_salt"`
	JWTSecret         string `yaml:"jwt_secret"`
	ServerURL         string `yaml:"server_url"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		_ = os.MkdirAll(configDir, 0755)
	}
	return filepath.Join(configDir, "bastion.yml"), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{ServerURL: "http://localhost:8081"}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.ServerURL == "" {
		config.ServerURL = "http://localhost:8081"
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
