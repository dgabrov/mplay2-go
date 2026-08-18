package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/endpoint"
	"log/slog"
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

	db, err := initDB(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	return startServer(cfg, db)
}

func logConfigFields(cfg data.ConfigData) {
	slog.Info("Version", "value", cfg.Version)
	slog.Info("Server Address", "value", cfg.Server.Address)
	slog.Info("Server Port", "value", cfg.Server.Port)
	slog.Info("DB Machine", "value", cfg.DB.Machine)
	slog.Info("DB Port", "value", cfg.DB.Port)
	slog.Info("DB Database", "value", cfg.DB.Database)
	slog.Info("DB User", "value", cfg.DB.User)
	slog.Info("Auth URL", "value", cfg.Auth.URL)
	slog.Info("Auth Right", "value", cfg.Auth.Right)
	slog.Info("Context", "value", cfg.Context)
}

func startServer(cfg data.ConfigData, db *sql.DB) error {
	mux := http.NewServeMux()
	endpoint.RegisterRoutes(mux, cfg.Context)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	slog.Info("Starting server", "address", addr)

	return http.ListenAndServe(addr, mux)
}
