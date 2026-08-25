package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/google/uuid"
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

func validateToken(ctx context.Context, tx *sql.Tx, token string, tokenValidity int) (string, error) {
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

	newExpiryDt := time.Now().Add(time.Duration(tokenValidity) * time.Second)
	updateQuery := "UPDATE session SET expiry_dt = ? WHERE token = ?"
	_, err = tx.ExecContext(ctx, updateQuery, newExpiryDt, token)
	if err != nil {
		return "", fmt.Errorf("failed to update token expiry: %w", err)
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
	// Build WHERE clause with AND conditions for each search term
	whereConditions := make([]string, len(searchTerms))
	args := make([]any, 0, len(searchTerms)+1)

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

func verifyMediaOwnership(ctx context.Context, tx *sql.Tx, userID string, mediaIds []string) error {
	if len(mediaIds) == 0 {
		return nil
	}

	placeholders := make([]string, len(mediaIds))
	args := make([]interface{}, len(mediaIds)+1)
	args[0] = userID

	for i, id := range mediaIds {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := "SELECT COUNT(*) FROM media WHERE user_id = ? AND media_id IN (" + strings.Join(placeholders, ",") + ")"
	var count int
	err := tx.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to verify media ownership: %w", err)
	}

	if count != len(mediaIds) {
		return fmt.Errorf("some media ids do not belong to this user or do not exist")
	}

	return nil
}

func deleteMedia(ctx context.Context, tx *sql.Tx, mediaIds []string) (int, error) {
	if len(mediaIds) == 0 {
		return 0, nil
	}

	placeholders := make([]string, len(mediaIds))
	args := make([]interface{}, len(mediaIds))

	for i, id := range mediaIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "DELETE FROM media WHERE media_id IN (" + strings.Join(placeholders, ",") + ")"
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete media: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func verifyPlaylistOwnership(ctx context.Context, tx *sql.Tx, userID string, playlistIds []string) error {
	if len(playlistIds) == 0 {
		return nil
	}

	placeholders := make([]string, len(playlistIds))
	args := make([]interface{}, len(playlistIds)+1)
	args[0] = userID

	for i, id := range playlistIds {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := "SELECT COUNT(*) FROM playlist WHERE user_id = ? AND playlist_id IN (" + strings.Join(placeholders, ",") + ")"
	var count int
	err := tx.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to verify playlist ownership: %w", err)
	}

	if count != len(playlistIds) {
		return fmt.Errorf("some playlist ids do not belong to this user or do not exist")
	}

	return nil
}

func deletePlaylist(ctx context.Context, tx *sql.Tx, playlistIds []string) (int, error) {
	if len(playlistIds) == 0 {
		return 0, nil
	}

	placeholders := make([]string, len(playlistIds))
	args := make([]interface{}, len(playlistIds))

	for i, id := range playlistIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "DELETE FROM playlist WHERE playlist_id IN (" + strings.Join(placeholders, ",") + ")"
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func addPlaylist(ctx context.Context, tx *sql.Tx, playlist *data.PlayList) error {
	query := "INSERT INTO playlist (playlist_id, user_id, description) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, playlist.PlaylistId, playlist.UserId, playlist.Description)
	if err != nil {
		return fmt.Errorf("failed to add playlist: %w", err)
	}

	return nil
}

func updatePlaylist(ctx context.Context, tx *sql.Tx, userID, playlistID, description string) error {
	query := "UPDATE playlist SET description = ? WHERE playlist_id = ? AND user_id = ?"
	_, err := tx.ExecContext(ctx, query, description, playlistID, userID)
	if err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}

	return nil
}

func getMediaForPlaylist(ctx context.Context, tx *sql.Tx, userID, playlistID string) ([]*data.ExtendedMedia, error) {
	query := `SELECT m.media_id, m.user_id, m.description, m.content_type, m.size, m.width, m.height, mp.seq_no
	FROM media m
	INNER JOIN media_playlist mp ON m.media_id = mp.media_id
	WHERE mp.playlist_id = ? AND m.user_id = ?
	ORDER BY mp.seq_no`

	rows, err := tx.QueryContext(ctx, query, playlistID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get media for playlist: %w", err)
	}
	defer rows.Close()

	var results []*data.ExtendedMedia
	for rows.Next() {
		var m data.ExtendedMedia
		err := rows.Scan(&m.Id, &m.UserId, &m.Description, &m.ContentType, &m.Size, &m.Width, &m.Height, &m.SeqNo)
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

func removeMediaFromPlaylist(ctx context.Context, tx *sql.Tx, userID, playlistID string, mediaIds []string) (int, error) {
	if len(mediaIds) == 0 {
		return 0, nil
	}

	// Verify playlist belongs to user first
	var count int
	query := "SELECT COUNT(*) FROM playlist WHERE playlist_id = ? AND user_id = ?"
	err := tx.QueryRowContext(ctx, query, playlistID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to verify playlist ownership: %w", err)
	}

	if count == 0 {
		return 0, fmt.Errorf("playlist does not belong to this user or does not exist")
	}

	placeholders := make([]string, len(mediaIds))
	args := make([]interface{}, len(mediaIds)+1)
	args[0] = playlistID

	for i, id := range mediaIds {
		placeholders[i] = "?"
		args[i+1] = id
	}

	deleteQuery := "DELETE FROM media_playlist WHERE playlist_id = ? AND media_id IN (" + strings.Join(placeholders, ",") + ")"
	result, err := tx.ExecContext(ctx, deleteQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to remove media from playlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func addMediaToPlaylist(ctx context.Context, tx *sql.Tx, userID, playlistID string, mediaIds []string) (int, error) {
	if len(mediaIds) == 0 {
		return 0, nil
	}

	// Verify playlist belongs to user
	var count int
	query := "SELECT COUNT(*) FROM playlist WHERE playlist_id = ? AND user_id = ?"
	err := tx.QueryRowContext(ctx, query, playlistID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to verify playlist ownership: %w", err)
	}

	if count == 0 {
		return 0, fmt.Errorf("playlist does not belong to this user or does not exist")
	}

	// Verify all media ids belong to user
	placeholders := make([]string, len(mediaIds))
	args := make([]interface{}, len(mediaIds)+1)
	args[0] = userID

	for i, id := range mediaIds {
		placeholders[i] = "?"
		args[i+1] = id
	}

	mediaQuery := "SELECT COUNT(*) FROM media WHERE user_id = ? AND media_id IN (" + strings.Join(placeholders, ",") + ")"
	var mediaCount int
	err = tx.QueryRowContext(ctx, mediaQuery, args...).Scan(&mediaCount)
	if err != nil {
		return 0, fmt.Errorf("failed to verify media ownership: %w", err)
	}

	if mediaCount != len(mediaIds) {
		return 0, fmt.Errorf("some media ids do not belong to this user or do not exist")
	}

	// Get existing media_ids already in this playlist
	existingQuery := "SELECT media_id FROM media_playlist WHERE playlist_id = ? AND media_id IN (" + strings.Join(placeholders, ",") + ")"
	existingArgs := make([]interface{}, len(mediaIds)+1)
	existingArgs[0] = playlistID
	copy(existingArgs[1:], args[1:])

	rows, err := tx.QueryContext(ctx, existingQuery, existingArgs...)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing media: %w", err)
	}
	defer rows.Close()

	existingIds := make(map[string]bool)
	for rows.Next() {
		var existingId string
		if err := rows.Scan(&existingId); err != nil {
			return 0, fmt.Errorf("failed to scan existing id: %w", err)
		}
		existingIds[existingId] = true
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating existing results: %w", err)
	}

	// Insert only new media to playlist
	addedCount := 0
	for _, mediaId := range mediaIds {
		if existingIds[mediaId] {
			continue // Skip already added media
		}

		seqNo, err := getNextSequenceValue(ctx, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to get sequence value: %w", err)
		}

		mediaPlaylistId := uuid.Must(uuid.NewV7()).String()
		insertQuery := "INSERT INTO media_playlist (media_playlist_id, playlist_id, media_id, seq_no) VALUES (?, ?, ?, ?)"
		_, err = tx.ExecContext(ctx, insertQuery, mediaPlaylistId, playlistID, mediaId, seqNo)
		if err != nil {
			return 0, fmt.Errorf("failed to add media to playlist: %w", err)
		}

		addedCount++
	}

	return addedCount, nil
}

func addMedia(ctx context.Context, tx *sql.Tx, media *data.Media) error {
	query := "INSERT INTO media (media_id, user_id, description, content_type, size, width, height) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, media.Id, media.UserId, media.Description, media.ContentType, media.Size, media.Width, media.Height)
	if err != nil {
		return fmt.Errorf("failed to add media: %w", err)
	}

	return nil
}

func getMedia(ctx context.Context, tx *sql.Tx, userID, mediaID string) (*data.Media, error) {
	var media data.Media
	query := "SELECT media_id, user_id, description, content_type, size, width, height FROM media WHERE media_id = ? AND user_id = ?"
	err := tx.QueryRowContext(ctx, query, mediaID, userID).Scan(&media.Id, &media.UserId, &media.Description, &media.ContentType, &media.Size, &media.Width, &media.Height)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return &media, nil
}

func updateMedia(ctx context.Context, tx *sql.Tx, userID, mediaID, description string) error {
	query := "UPDATE media SET description = ? WHERE media_id = ? AND user_id = ?"
	_, err := tx.ExecContext(ctx, query, description, mediaID, userID)
	if err != nil {
		return fmt.Errorf("failed to update media: %w", err)
	}

	return nil
}

func updateMediaWithType(ctx context.Context, tx *sql.Tx, userID, mediaID, description, contentType string) error {
	query := "UPDATE media SET description = ?, content_type = ? WHERE media_id = ? AND user_id = ?"
	_, err := tx.ExecContext(ctx, query, description, contentType, mediaID, userID)
	if err != nil {
		return fmt.Errorf("failed to update media: %w", err)
	}

	return nil
}

func updateMediaSize(ctx context.Context, tx *sql.Tx, userID, mediaID string, size int64) error {
	query := "UPDATE media SET size = ? WHERE media_id = ? AND user_id = ?"
	_, err := tx.ExecContext(ctx, query, size, mediaID, userID)
	if err != nil {
		return fmt.Errorf("failed to update media size: %w", err)
	}

	return nil
}

func updateMediaDimensions(ctx context.Context, tx *sql.Tx, userID, mediaID string, width, height int) error {
	query := "UPDATE media SET width = ?, height = ? WHERE media_id = ? AND user_id = ?"
	_, err := tx.ExecContext(ctx, query, width, height, mediaID, userID)
	if err != nil {
		return fmt.Errorf("failed to update media dimensions: %w", err)
	}

	return nil
}

func switchMediaSequence(ctx context.Context, tx *sql.Tx, userID, playlistID, media1ID, media2ID string) error {
	// Verify playlist belongs to user
	var count int
	query := "SELECT COUNT(*) FROM playlist WHERE playlist_id = ? AND user_id = ?"
	err := tx.QueryRowContext(ctx, query, playlistID, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to verify playlist ownership: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("playlist does not belong to this user or does not exist")
	}

	// Get seq_no for both media items
	var seq1, seq2 int64
	query = "SELECT seq_no FROM media_playlist WHERE playlist_id = ? AND media_id = ?"
	err = tx.QueryRowContext(ctx, query, playlistID, media1ID).Scan(&seq1)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("media1 is not part of this playlist")
		}
		return fmt.Errorf("failed to get media1 sequence: %w", err)
	}

	err = tx.QueryRowContext(ctx, query, playlistID, media2ID).Scan(&seq2)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("media2 is not part of this playlist")
		}
		return fmt.Errorf("failed to get media2 sequence: %w", err)
	}

	// Use intermediary value to switch sequences
	intermediaryValue := int64(-9999)

	// Update media1 to intermediary value
	updateQuery := "UPDATE media_playlist SET seq_no = ? WHERE playlist_id = ? AND media_id = ?"
	_, err = tx.ExecContext(ctx, updateQuery, intermediaryValue, playlistID, media1ID)
	if err != nil {
		return fmt.Errorf("failed to update media1 sequence: %w", err)
	}

	// Update media2 to media1's old seq_no
	_, err = tx.ExecContext(ctx, updateQuery, seq1, playlistID, media2ID)
	if err != nil {
		return fmt.Errorf("failed to update media2 sequence: %w", err)
	}

	// Update media1 to media2's old seq_no
	_, err = tx.ExecContext(ctx, updateQuery, seq2, playlistID, media1ID)
	if err != nil {
		return fmt.Errorf("failed to finalize media1 sequence: %w", err)
	}

	return nil
}
