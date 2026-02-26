package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	URL  string
	Port string
}

func LoadConfig() (*Config, error) {
	url := os.Getenv("BASTION_HOST")
	if url == "" {
		url = "http://localhost:8287"
	}

	port := os.Getenv("BASTION_PORT")
	if port == "" {
		port = "8287"
	}

	// Ensure protocol
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	return &Config{
		URL:  url,
		Port: port,
	}, nil
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{URL: %s, Port: %s}", c.URL, c.Port)
}
