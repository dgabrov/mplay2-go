package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/google/uuid"
	"strings"
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

func validateToken(ctx context.Context, tx *sql.Tx, token string) (string, error) {
	var userID string
	var expiredInd string
	var expiryDt time.Time

	query := "SELECT user_id, expired_ind, expiry_dt FROM session WHERE token = ?"
	err := tx.QueryRowContext(ctx, query, token).Scan(&userID, &expiredInd, &expiryDt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid token or token expired")
		}
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	if expiredInd == "Y" || time.Now().After(expiryDt) {
		return "", fmt.Errorf("invalid token or token expired")
	}

	return userID, nil
}

func getNextSequenceValue(ctx context.Context, tx *sql.Tx) (int64, error) {
	var seqval int64
	query := "SELECT seqval FROM seqvalues WHERE id = '1' FOR UPDATE"
	err := tx.QueryRowContext(ctx, query).Scan(&seqval)
	if err != nil {
		return 0, fmt.Errorf("failed to get sequence value: %w", err)
	}

	newSeqval := seqval + 1
	updateQuery := "UPDATE seqvalues SET seqval = ? WHERE id = '1'"
	_, err = tx.ExecContext(ctx, updateQuery, newSeqval)
	if err != nil {
		return 0, fmt.Errorf("failed to update sequence value: %w", err)
	}

	return newSeqval, nil
}

func searchMedia(ctx context.Context, tx *sql.Tx, userID string, searchTerms []string) ([]*data.Media, error) {
	if len(searchTerms) == 0 {
		return []*data.Media{}, nil
	}

	// Build WHERE clause with AND conditions for each search term
	whereConditions := make([]string, len(searchTerms))
	args := make([]interface{}, 0, len(searchTerms)+1)

	args = append(args, userID)
	for i := range searchTerms {
		whereConditions[i] = "description LIKE ?"
		args = append(args, searchTerms[i])
	}

	whereClause := "user_id = ? AND (" + strings.Join(whereConditions, " AND ") + ")"

	query := "SELECT media_id, user_id, description, content_type, size, width, height FROM media WHERE " + whereClause
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search media: %w", err)
	}
	defer rows.Close()

	var results []*data.Media
	for rows.Next() {
		var m data.Media
		err := rows.Scan(&m.Id, &m.UserId, &m.Description, &m.ContentType, &m.Size, &m.Width, &m.Height)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media: %w", err)
		}
		results = append(results, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return results, nil
}

func searchPlaylist(ctx context.Context, tx *sql.Tx, userID string, searchTerms []string) ([]*data.PlayList, error) {
	if len(searchTerms) == 0 {
		return []*data.PlayList{}, nil
	}

	// Build WHERE clause with AND conditions for each search term
	whereConditions := make([]string, len(searchTerms))
	args := make([]interface{}, 0, len(searchTerms)+1)

	args = append(args, userID)
	for i := range searchTerms {
		whereConditions[i] = "description LIKE ?"
		args = append(args, searchTerms[i])
	}

	whereClause := "user_id = ? AND (" + strings.Join(whereConditions, " AND ") + ")"

	query := "SELECT playlist_id, user_id, description FROM playlist WHERE " + whereClause
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlist: %w", err)
	}
	defer rows.Close()

	var results []*data.PlayList
	for rows.Next() {
		var p data.PlayList
		err := rows.Scan(&p.PlaylistId, &p.UserId, &p.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		results = append(results, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return results, nil
}
