package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
)

type Servr struct {
	db     *sql.DB
	config *data.ConfigData
}

func NewServr(db *sql.DB, cfg *data.ConfigData) *Servr {
	return &Servr{
		db:     db,
		config: cfg,
	}
}

func (s *Servr) GetUserByProvidedId(ctx context.Context, providedUserID string) (*data.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	user, err := getUserByProvidedId(ctx, tx, providedUserID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

func (s *Servr) CreateUser(ctx context.Context, user *data.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Servr) CreateSession(ctx context.Context, userID, token string) (*data.Session, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	session, err := createSession(ctx, tx, userID, token, s.config.TokenValidity)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return session, nil
}
