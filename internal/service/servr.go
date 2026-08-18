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
		return nil, err
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
		return err
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
		return nil, err
	}

	return session, nil
}

func (s *Servr) ValidateToken(ctx context.Context, token string) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	userID, err := validateToken(ctx, tx, token)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return userID, nil
}

func (s *Servr) GetNextSequenceValue(ctx context.Context) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	seqval, err := getNextSequenceValue(ctx, tx)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return seqval, nil
}

func (s *Servr) SearchMedia(ctx context.Context, userID string, searchTerms []string) ([]*data.Media, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	results, err := searchMedia(ctx, tx, userID, searchTerms)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Servr) SearchPlaylist(ctx context.Context, userID string, searchTerms []string) ([]*data.PlayList, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	results, err := searchPlaylist(ctx, tx, userID, searchTerms)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Servr) VerifyMediaOwnership(ctx context.Context, userID string, mediaIds []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := verifyMediaOwnership(ctx, tx, userID, mediaIds); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Servr) DeleteMedia(ctx context.Context, mediaIds []string) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	count, err := deleteMedia(ctx, tx, mediaIds)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Servr) VerifyPlaylistOwnership(ctx context.Context, userID string, playlistIds []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := verifyPlaylistOwnership(ctx, tx, userID, playlistIds); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Servr) DeletePlaylist(ctx context.Context, playlistIds []string) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	count, err := deletePlaylist(ctx, tx, playlistIds)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Servr) AddPlaylist(ctx context.Context, playlist *data.PlayList) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := addPlaylist(ctx, tx, playlist); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Servr) UpdatePlaylist(ctx context.Context, userID, playlistID, description string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := updatePlaylist(ctx, tx, userID, playlistID, description); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Servr) GetMediaForPlaylist(ctx context.Context, userID, playlistID string) ([]*data.Media, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	media, err := getMediaForPlaylist(ctx, tx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return media, nil
}
