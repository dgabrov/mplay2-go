# IMPLEMENTED - POST removeMediaFromPlaylist
payload structure:
```json
{
  "playlistId": "playlist id",
  "ids": ["first", "second", "third"]
}
```

response is SuccessResponse (defined in @internal/data/business.go)

Implementation:
- Validates user token from cookie (jtoken12)
- Verifies playlistId belongs to authenticated user
- Deletes records from media_playlist table where playlist_id matches and media_id is in the provided ids list
- Returns success:true on completion
