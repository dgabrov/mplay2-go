package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mplay2-go/internal/data"
	"mplay2-go/internal/endpoint"
	"net/http"
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

	return startServer(cfg)
}

func logConfigFields(cfg data.ConfigData) {
	slog.Info("Version", "value", cfg.Version)
	slog.Info("Server Address", "value", cfg.Server.Address)
	slog.Info("Server Port", "value", cfg.Server.Port)
}

func startServer(cfg data.ConfigData) error {
	mux := http.NewServeMux()
	endpoint.RegisterRoutes(mux)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	slog.Info("Starting server", "address", addr)

	return http.ListenAndServe(addr, mux)
}
