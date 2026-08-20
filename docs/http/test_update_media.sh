#!/bin/bash

# Test script for POST /updateMedia endpoint
# Usage: ./test_update_media.sh [add|update] [token] [media-id] [description] [file]

# Example token (should be obtained from login endpoint)
TOKEN="${2:- }"
MEDIA_ID="${3:-test-media-1}"
DESCRIPTION="${4:-Test media}"
FILE="${5:-./sample.mp4}"

if [ "$1" = "add" ]; then
  echo "Adding new media..."
  curl -X POST http://localhost:8080/mediaserv/updateMedia \
    -H "Cookie: token=$TOKEN" \
    -F "data={\"adding\": true, \"id\": \"$MEDIA_ID\", \"description\": \"$DESCRIPTION\"};type=application/json" \
    -F "file=@$FILE" \
    -v
elif [ "$1" = "update" ]; then
  echo "Updating media (description only)..."
  curl -X POST http://localhost:8080/mediaserv/updateMedia \
    -H "Cookie: token=$TOKEN" \
    -F "data={\"adding\": false, \"id\": \"$MEDIA_ID\", \"description\": \"$DESCRIPTION\"};type=application/json" \
    -v
elif [ "$1" = "update-with-file" ]; then
  echo "Updating media (with new file)..."
  curl -X POST http://localhost:8080/mediaserv/updateMedia \
    -H "Cookie: token=$TOKEN" \
    -F "data={\"adding\": false, \"id\": \"$MEDIA_ID\", \"description\": \"$DESCRIPTION\"};type=application/json" \
    -F "file=@$FILE" \
    -v
else
  echo "Usage: $0 {add|update|update-with-file} [token] [media-id] [description] [file]"
  echo ""
  echo "Examples:"
  echo "  # Add new media (requires login first to get token)"
  echo "  $0 add <token> my-media-1 'My Video' ./sample.mp4"
  echo ""
  echo "  # Update description only"
  echo "  $0 update <token> my-media-1 'Updated description'"
  echo ""
  echo "  # Update with new file"
  echo "  $0 update-with-file <token> my-media-1 'Updated video' ./sample.mp4"
fi
