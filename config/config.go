package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port int    `toml:"port"`
	Env  string `toml:"env"`
}

func NewConfig(cfgFile string) (*Config, error) {
	var cfg Config

	metadata, err := toml.DecodeFile(cfgFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	if len(metadata.Undecoded()) > 0 {
		return nil, fmt.Errorf("unknown configuration keys: %v", metadata.Undecoded())
	}

	return &cfg, nil
}
