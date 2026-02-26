package commands

import (
	"github.com/dcdavidev/bastion/packages/core/config"
)

type Config = config.Config

func LoadConfig() (*Config, error) {
	return config.LoadConfig()
}
