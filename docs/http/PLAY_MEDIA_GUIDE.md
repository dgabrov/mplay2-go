# GET /playMedia Endpoint Guide

This guide demonstrates how to use the `/playMedia` endpoint for streaming and downloading media files with HTTP range request support.

## Overview

The `/playMedia` endpoint streams media files in chunks, supporting standard HTTP range requests for resume capability and seeking.

## Features

- ✅ HTTP 206 Partial Content responses
- ✅ Range request support (resume, seeking)
- ✅ Configurable chunk size (default: 150KB)
- ✅ Automatic content-type detection from database
- ✅ User authorization and ownership verification
- ✅ Efficient streaming (doesn't load entire file in memory)

## Quick Start

### 1. Login to get a session token

```bash
curl -X POST http://localhost:8080/mediaserv/login \
  -H "Content-Type: application/json" \
  -d '{"login": "test1", "password": "test1"}'
```

### 2. Stream media from beginning

```bash
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN"
```

### 3. Stream with range request (bytes 0-99999)

```bash
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Range: bytes=0-99999"
```

### 4. Resume download from byte 500000

```bash
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN" \
  -H "Range: bytes=500000-" \
  -o video.mp4
```

## Range Request Formats

The endpoint supports all standard HTTP range request formats:

| Format | Example | Meaning |
|--------|---------|---------|
| `bytes=0-999` | Request first 1000 bytes | Specific range |
| `bytes=500000-` | Request from byte 500000 to end | Open-ended range |
| `bytes=-1000` | Request last 1000 bytes | Suffix range |
| (no Range header) | Request entire file in chunks | Full download |

## Response Headers

### Success Response (206 Partial Content)

```
HTTP/1.1 206 Partial Content
Content-Type: video/mp4
Content-Length: 150000
Content-Range: bytes 0-149999/3114374
Accept-Ranges: bytes
```

### Error Response (416 Requested Range Not Satisfiable)

```
HTTP/1.1 416 Range Not Satisfiable
Content-Range: bytes */3114374
```

## Query Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `id` | Yes | Media ID to play |

## Request Headers

| Header | Optional | Description |
|--------|----------|-------------|
| `Range` | Yes | HTTP range specification (e.g., `bytes=0-99999`) |
| `Cookie` | Yes | Session token (e.g., `token=...`) |

## Using the Test Script

```bash
# Make script executable
chmod +x docs/http/test_play_media.sh

# Login
./docs/http/test_play_media.sh login

# Stream from beginning
TOKEN=<your_token> ./docs/http/test_play_media.sh play video-001

# Stream with range
TOKEN=<your_token> ./docs/http/test_play_media.sh play-range video-001 500000

# Stream specific range
TOKEN=<your_token> ./docs/http/test_play_media.sh play-range video-001 500000 600000
```

## Example: Download with Resume Support

```bash
# Download with resume support
curl -C - -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN" \
  -o video.mp4

# curl's -C - flag automatically handles resume using Range header
```

## Example: Video Player Integration

```bash
# Stream to video player (VLC, ffplay, etc.)
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN" | ffplay -
```

## Chunk Size Configuration

The chunk size is configured via `mediaSlice` in config.json. Default is 150,000 bytes per request.

For a 3GB file:
- First request returns bytes 0-149,999 (150KB)
- Second request with `Range: bytes=150000-` returns next 150KB
- And so on...

## Database Lookups

The endpoint performs these database operations:

1. **Token validation**: Verify session token is valid
2. **Media lookup**: Verify media exists and belongs to user
3. **Size retrieval**: Get file size from media table
4. **Type retrieval**: Get content-type from media table

## Error Handling

| Status | Meaning | Example |
|--------|---------|---------|
| 400 | Bad Request | Missing `id` parameter or invalid Range header |
| 401 | Unauthorized | Invalid/missing session token |
| 404 | Not Found | Media doesn't exist or doesn't belong to user |
| 416 | Range Not Satisfiable | Range exceeds file size |
| 206 | Partial Content | Success with range request |

## Performance Characteristics

- **Memory usage**: O(chunk size) - typically 150KB
- **File I/O**: Single file seek + sequential read per request
- **Database**: Single query per request (caches token)
- **Network**: Efficient streaming without buffering entire file

## Browser Compatibility

This endpoint works with all modern browsers that support:
- HTTP Range requests
- HTML5 `<video>` element
- HTML5 `<audio>` element

Example HTML5 Video:
```html
<video width="320" height="240" controls>
  <source src="http://localhost:8080/mediaserv/playMedia?id=video-001" type="video/mp4">
  Your browser does not support the video tag.
</video>
```

## Supported Content Types

The endpoint supports any content-type stored in the database:
- Audio: `audio/mpeg`, `audio/wav`, `audio/ogg`, `audio/flac`, etc.
- Video: `video/mp4`, `video/webm`, `video/x-msvideo`, etc.

The actual type is determined during media upload via format detection.
