**step 1 - IMPLEMENTED**

I need a GET endpoint called playMedia
- It gets a query parameter called id
- It gets a range header

Actions:
- check the user is logged ✅
- check the id passed as query parameter is a good media entry in media table and that is belongs to the logged in user ✅
- you serve slices from the file associated with the media entry, in the folder UploadsFolder in ConfigData, I remind you the name of the file would be <<id>>.dat where id is the mediaId ✅
- when you serve the slices, you put as well the content type that you retrieve from the database and the Content-Range header that shows exactly how much data you passed ✅
- up to you if you go and take the file size from the file itself or you load it from the database alongside the content type - IMPLEMENTED: Takes file size from file itself ✅
- the slices are like that
  - you serve maximum of configuration MediaSlice from ConfigData. Of course, if there aren't as many bytes from the required start index, you serve whatever you have available ✅

- if error, of course, you return some http error code 400+. But if all goes well, you always return 206, of course, use the standard http constant ✅

## Implementation Details

### Endpoint: GET /playMedia

**Query Parameters:**
- `id` (required): Media ID to stream

**Headers:**
- `Range` (optional): HTTP range request header (RFC 7233)
  - Format: `bytes=START-END` or `bytes=START-` or `bytes=-SUFFIX`
  - Examples: `bytes=0-999`, `bytes=500000-`, `bytes=-1000`

**Authentication:**
- Session token via Cookie header

**Response Codes:**
- `206 Partial Content`: Successful range request
- `400 Bad Request`: Missing or invalid parameters
- `416 Range Not Satisfiable`: Range exceeds file size
- Other `4xx`: Validation or authorization errors

**Response Headers:**
- `Content-Type`: From database (audio or video)
- `Content-Length`: Bytes being sent
- `Content-Range`: Format `bytes START-END/TOTAL`
- `Accept-Ranges: bytes`: Indicates range support

**Features:**
- Full HTTP range request support (RFC 7233)
- Configurable chunk size via `mediaSlice` config (default: 150KB)
- Efficient streaming without loading entire file to memory
- File seek support for resume capability
- User ownership validation
- Token auto-refresh on each request

**Files Created:**
- `internal/endpoint/get_play_media.go` - Main handler with range parsing
- `docs/http/play_media.http` - HTTP client examples
- `docs/http/test_play_media.sh` - Bash test script
- `docs/http/PLAY_MEDIA_GUIDE.md` - Complete documentation

**Usage Examples:**

```bash
# Stream from beginning
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN"

# Resume from byte 500000
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN" \
  -H "Range: bytes=500000-"

# Request specific range
curl -X GET "http://localhost:8080/mediaserv/playMedia?id=video-001" \
  -H "Cookie: token=YOUR_TOKEN" \
  -H "Range: bytes=0-99999"
```
