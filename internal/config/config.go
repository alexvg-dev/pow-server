package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"strings"
)

type Config struct {
	Server         ServerConfig `json:"server"`
	SessionTtlSec  int          `json:"sessionTtlSec" env:"SESSION_TTL_SEC" default:"10"`
	QuotesFilePath string       `json:"quotesFilePath" env:"QUOTES_FILE_PATH"`
}

type ServerConfig struct {
	Port int `json:"port" env:"SERVER_PORT"`
}

func GetConfig() (*Config, error) {

	configPath := os.Getenv("PATH_CONFIG")
	if strings.TrimSpace(configPath) == "" {
		return nil, errors.New("PATH_CONFIG is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist - `%s`", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("config reading: %w", err)
	}

	return &cfg, nil
}
