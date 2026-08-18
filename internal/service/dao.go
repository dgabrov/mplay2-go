package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/google/uuid"
	"time"
)

func getUserByProvidedId(ctx context.Context, tx *sql.Tx, providedUserID string) (*data.User, error) {
	var user data.User
	query := "SELECT user_id, provided_user_id, login, name FROM user WHERE provided_user_id = ?"
	err := tx.QueryRowContext(ctx, query, providedUserID).Scan(&user.UserID, &user.ProvidedUserID, &user.Login, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func createUser(ctx context.Context, tx *sql.Tx, user *data.User) error {
	if user.UserID == "" {
		user.UserID = uuid.Must(uuid.NewV7()).String()
	}
	query := "INSERT INTO user (user_id, provided_user_id, login, name) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, user.UserID, user.ProvidedUserID, user.Login, user.Name)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func createSession(ctx context.Context, tx *sql.Tx, userID, token string, tokenValidity int) (*data.Session, error) {
	sessionID := uuid.Must(uuid.NewV7()).String()
	expiryDt := time.Now().Add(time.Duration(tokenValidity) * time.Second)

	query := "INSERT INTO session (session_id, user_id, token, expired_ind, expiry_dt) VALUES (?, ?, ?, 'N', ?)"
	_, err := tx.ExecContext(ctx, query, sessionID, userID, token, expiryDt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &data.Session{
		SessionID:  sessionID,
		UserID:     userID,
		Token:      token,
		ExpiredInd: "N",
		ExpiryDt:   expiryDt.Format(time.RFC3339),
	}, nil
}
