# IMPLEMENTED - POST addMediaToPlaylist

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
- Verifies playlist belongs to authenticated user
- Verifies all media ids exist and belong to authenticated user
- Checks which media ids are already in media_playlist for this playlist and ignores them
- For new media:
  - Creates media_playlist_id as UUID v7
  - Uses provided playlist_id
  - Uses each media_id from ids array
  - Gets seq_no from GetNextSequenceValue()
- Returns success:true and count of newly added media
