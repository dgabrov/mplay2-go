#!/bin/bash

# Test script for POST /updateMedia endpoint
# This script demonstrates all operations available through the endpoint

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_URL="${SERVER_URL:-http://localhost:8080/mediaserv}"
TOKEN="${TOKEN:-}"
SAMPLE_FILE="${SCRIPT_DIR}/sample.mp4"

usage() {
  cat << EOF
Usage: $0 {login|add|update|update-file} [options]

Commands:
  login                        Get session token (required before other operations)
  add <id> [description]       Add new media with sample.mp4
  update <id> <description>    Update media description only
  update-file <id> <desc>      Update media with new file

Environment Variables:
  TOKEN                        Session token (obtained from login)
  SERVER_URL                   Server URL (default: http://localhost:8080/mediaserv)

Examples:
  # 1. Login to get token
  $0 login

  # 2. Add new media
  TOKEN=<token> $0 add video-001 "My Video"

  # 3. Update description only
  TOKEN=<token> $0 update video-001 "Updated description"

  # 4. Update with new file
  TOKEN=<token> $0 update-file video-001 "New video"
EOF
  exit 1
}

login() {
  echo "Logging in..."
  curl -X POST "$SERVER_URL/login" \
    -H "Content-Type: application/json" \
    -d '{
      "login": "test1",
      "password": "test1"
    }' \
    -c /tmp/cookies.txt \
    -b /tmp/cookies.txt \
    -w "\n" \
    -s | jq '.' || true

  echo ""
  echo "Token has been saved to cookies. Use -b /tmp/cookies.txt with curl commands."
}

add_media() {
  local id=$1
  local description=${2:-$(basename "$SAMPLE_FILE" .mp4)}

  if [ -z "$TOKEN" ]; then
    echo "Error: TOKEN environment variable not set. Run 'login' first."
    exit 1
  fi

  if [ ! -f "$SAMPLE_FILE" ]; then
    echo "Error: Sample file not found at $SAMPLE_FILE"
    exit 1
  fi

  echo "Adding media '$id' with description '$description'..."
  curl -X POST "$SERVER_URL/updateMedia" \
    -H "Cookie: token=$TOKEN" \
    -F "data={\"adding\": true, \"id\": \"$id\", \"description\": \"$description\"};type=application/json" \
    -F "file=@$SAMPLE_FILE" \
    -w "\n" \
    -s | jq '.' || true
}

update_description() {
  local id=$1
  local description=$2

  if [ -z "$TOKEN" ]; then
    echo "Error: TOKEN environment variable not set. Run 'login' first."
    exit 1
  fi

  echo "Updating media '$id' description to '$description'..."
  curl -X POST "$SERVER_URL/updateMedia" \
    -H "Cookie: token=$TOKEN" \
    -F "data={\"adding\": false, \"id\": \"$id\", \"description\": \"$description\"};type=application/json" \
    -w "\n" \
    -s | jq '.' || true
}

update_with_file() {
  local id=$1
  local description=$2

  if [ -z "$TOKEN" ]; then
    echo "Error: TOKEN environment variable not set. Run 'login' first."
    exit 1
  fi

  if [ ! -f "$SAMPLE_FILE" ]; then
    echo "Error: Sample file not found at $SAMPLE_FILE"
    exit 1
  fi

  echo "Updating media '$id' with new file and description '$description'..."
  curl -X POST "$SERVER_URL/updateMedia" \
    -H "Cookie: token=$TOKEN" \
    -F "data={\"adding\": false, \"id\": \"$id\", \"description\": \"$description\"};type=application/json" \
    -F "file=@$SAMPLE_FILE" \
    -w "\n" \
    -s | jq '.' || true
}

if [ $# -eq 0 ]; then
  usage
fi

case "$1" in
  login)
    login
    ;;
  add)
    if [ $# -lt 2 ]; then
      echo "Error: missing media id"
      usage
    fi
    add_media "$2" "$3"
    ;;
  update)
    if [ $# -lt 3 ]; then
      echo "Error: missing media id or description"
      usage
    fi
    update_description "$2" "$3"
    ;;
  update-file)
    if [ $# -lt 3 ]; then
      echo "Error: missing media id or description"
      usage
    fi
    update_with_file "$2" "$3"
    ;;
  *)
    echo "Error: unknown command '$1'"
    usage
    ;;
esac
