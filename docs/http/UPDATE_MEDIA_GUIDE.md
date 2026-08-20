# POST /updateMedia Endpoint Guide

This guide demonstrates how to use the `/updateMedia` endpoint to add and update media files.

## Overview

The endpoint accepts multipart form data with two parts:
1. **data** (JSON): Metadata with structure `{adding: bool, id: string, description: string}`
2. **file** (optional): The media file to upload

## Automatic Features

The endpoint automatically:
- **Detects media type** (audio or video) from file content
- **Calculates file size** in bytes
- **Extracts video dimensions** (width/height) for video files
- **Resets dimensions to 0** when updating from video to audio

## Quick Start

### 1. Login to get a session token

```bash
curl -X POST http://localhost:8080/mediaserv/login \
  -H "Content-Type: application/json" \
  -d '{"login": "test1", "password": "test1"}' \
  -c cookies.txt
```

### 2. Add new media

```bash
curl -X POST http://localhost:8080/mediaserv/updateMedia \
  -H "Cookie: token=YOUR_TOKEN" \
  -F 'data={"adding": true, "id": "video-001", "description": "My Video"};type=application/json' \
  -F "file=@docs/http/sample.mp4"
```

### 3. Update description only

```bash
curl -X POST http://localhost:8080/mediaserv/updateMedia \
  -H "Cookie: token=YOUR_TOKEN" \
  -F 'data={"adding": false, "id": "video-001", "description": "Updated Description"};type=application/json'
```

### 4. Update with new file

```bash
curl -X POST http://localhost:8080/mediaserv/updateMedia \
  -H "Cookie: token=YOUR_TOKEN" \
  -F 'data={"adding": false, "id": "video-001", "description": "Replaced Video"};type=application/json' \
  -F "file=@docs/http/sample.mp4"
```

## Using the Test Script

For convenience, use the provided test script:

```bash
# Make script executable
chmod +x docs/http/test_update_media.sh

# Login
./docs/http/test_update_media.sh login

# Add new media (replace TOKEN with actual token)
TOKEN=<your_token> ./docs/http/test_update_media.sh add video-001 "My Video"

# Update description
TOKEN=<your_token> ./docs/http/test_update_media.sh update video-001 "Updated"

# Update with new file
TOKEN=<your_token> ./docs/http/test_update_media.sh update-file video-001 "New Video"
```

## Using HTTP Client Files

If using VS Code REST Client or JetBrains HTTP client, use the `update_media.http` file:

1. Open `docs/http/update_media.http` in your editor
2. Replace `YOUR_TOKEN_HERE` with actual token from login response
3. Click "Send Request" on each request

## Request Details

### Adding Media (adding: true)

**Required:**
- `adding`: true
- `id`: unique media identifier
- `file`: the media file (required)

**Optional:**
- `description`: media description
  - If not provided, auto-filled from filename (without extension)

**Errors:**
- Returns `"cannot add media without the actual content"` if file is missing

### Updating Media (adding: false)

**Required:**
- `adding`: false
- `id`: media identifier to update

**Optional:**
- `description`: new description (if omitted, keeps existing)
- `file`: new media file (if omitted, keeps existing file)

**Validations:**
- Media must exist
- Media must belong to logged-in user

## Response Examples

### Success Response
```json
{
  "success": true
}
```

### Error Response
```json
{
  "message": "cannot add media without the actual content",
  "items": []
}
```

## Database Fields Populated

After upload, the following fields are stored:

| Field | Source | Notes |
|-------|--------|-------|
| `media_id` | Provided in request | UUID if not provided |
| `user_id` | From session token | |
| `description` | Provided or auto-filled from filename | |
| `content_type` | Auto-detected from file | "audio" or "video" |
| `size` | Calculated during upload | Bytes |
| `width` | Extracted from video metadata | 0 if audio |
| `height` | Extracted from video metadata | 0 if audio |

## Large File Support

The endpoint supports large files (tested with 2GB+) through:
- Streaming file transfer (doesn't load entire file into memory)
- Efficient metadata extraction (only reads necessary bytes)

## Supported Formats

### Video
- MP4/MOV (ftyp)
- WebM/Matroska (mkv, mka)
- AVI
- MPEG-TS
- FLV

### Audio
- MP3
- WAV
- OGG
- FLAC
- AAC/M4A

## File Storage

Files are stored as `{media_id}.dat` in the folder specified by `UploadsFolder` configuration.

Example: Media with id `abc-123` is stored as `abc-123.dat`
