package service

import (
	"database/sql"
	"mplay2-go/internal/data"
)

type Servr struct {
	db     *sql.DB
	config *data.ConfigData
}

func NewServr(db *sql.DB) *Servr {
	return &Servr{}
}
