package commands

import (
	"github.com/dcdavidev/bastion/packages/config"
)

type Config = config.Config

func LoadConfig() (*Config, error) {
	return config.LoadConfig()
}
