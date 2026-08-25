# IMPLEMENTED - GET getMediaForPlaylist
query parameter playlistId
returns array of Media in results field

Implementation:
- Validates user token from cookie (jtoken12)
- Verifies playlist belongs to authenticated user
- Joins media_playlist to media table and returns ordered by seq_no
- Returns empty array if no media found
