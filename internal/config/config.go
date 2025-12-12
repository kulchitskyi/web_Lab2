package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerPort             string `json:"server_port"`
	DatabaseURL            string `json:"database_url"`
	LogLevel               string `json:"log_level"`
	EnableCache            bool   `json:"enable_cache"`
	RequestTimeoutSeconds  int    `json:"request_timeout_seconds"`
	MaxConnections         int    `json:"max_connections"`
	EnableRequestLogging   bool   `json:"enable_request_logging"`
	AllowPlaceDeletion     bool   `json:"allow_place_deletion"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")

	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		cfg.ServerPort = envPort
	}
	if envLog := os.Getenv("LOG_LEVEL"); envLog != "" {
		cfg.LogLevel = envLog
	}

	return &cfg, nil
}
