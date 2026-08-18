package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mplay2-go/internal/data"
	"os"
)

func Start() error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return fmt.Errorf("CONFIG_FILE environment variable not set")
	}

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg data.ConfigData
	if err := json.Unmarshal(configBytes, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	logConfigFields(cfg)

	return nil
}

func logConfigFields(cfg data.ConfigData) {
	slog.Info("Version", "value", cfg.Version)
}
