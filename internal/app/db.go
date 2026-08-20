package app

import (
	"database/sql"
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)

func initDB(cfg data.DbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Machine, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Database connection successful")
	return db, nil
}
