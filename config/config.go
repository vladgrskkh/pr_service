package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Port    int    `toml:"port"`
	Env     string `toml:"env"`
	Version string `toml:"version"`
	DB      struct {
		DSN          string
		MaxOpenConns int32  `toml:"max_open_conns"`
		MaxIdleTime  string `toml:"max_idle_time"`
	} `toml:"db"`
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

	err = godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	if cfg.Env == "development" {
		cfg.DB.DSN = os.Getenv("DB_DSN_LOCAL")
	} else {
		cfg.DB.DSN = os.Getenv("DB_DSN")
	}

	return &cfg, nil
}
