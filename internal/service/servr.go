package service

import (
	"database/sql"
	"github.com/amanagement24/mplay2-go/internal/data"
)

type Servr struct {
	db     *sql.DB
	config *data.ConfigData
}

func NewServr(db *sql.DB) *Servr {
	return &Servr{}
}
