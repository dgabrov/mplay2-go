#!/bin/bash

# Test script for GET /playMedia endpoint
# Demonstrates range request support for streaming media

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_URL="${SERVER_URL:-http://localhost:8080/mediaserv}"
TOKEN="${TOKEN:-}"

usage() {
  cat << EOF
Usage: $0 {login|play|play-range} [options]

Commands:
  login                        Get session token (required before other operations)
  play <id>                    Play media from beginning
  play-range <id> <start> [end]   Play media with range request

Environment Variables:
  TOKEN                        Session token (obtained from login)
  SERVER_URL                   Server URL (default: http://localhost:8080/mediaserv)

Examples:
  # 1. Login to get token
  $0 login

  # 2. Play media from beginning
  TOKEN=<token> $0 play video-001

  # 3. Play media from byte 500000 onwards
  TOKEN=<token> $0 play-range video-001 500000

  # 4. Play media from byte 500000 to 600000
  TOKEN=<token> $0 play-range video-001 500000 600000

  # 5. Download to file with resume support
  TOKEN=<token> $0 play-range video-001 0 > video.mp4
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
    -w "\n" \
    -s | jq '.' || true

  echo ""
  echo "Token logged in. Use -b /tmp/cookies.txt with curl commands."
}

play_media() {
  local id=$1

  if [ -z "$TOKEN" ]; then
    echo "Error: TOKEN environment variable not set. Run 'login' first."
    exit 1
  fi

  echo "Playing media '$id' from beginning..."
  curl -X GET "$SERVER_URL/playMedia?id=$id" \
    -H "Cookie: token=$TOKEN" \
    -i
}

play_media_range() {
  local id=$1
  local start=$2
  local end=$3

  if [ -z "$TOKEN" ]; then
    echo "Error: TOKEN environment variable not set. Run 'login' first."
    exit 1
  fi

  if [ -z "$start" ]; then
    echo "Error: missing start byte position"
    usage
  fi

  if [ -n "$end" ]; then
    range="bytes=$start-$end"
  else
    range="bytes=$start-"
  fi

  echo "Playing media '$id' with range: $range"
  curl -X GET "$SERVER_URL/playMedia?id=$id" \
    -H "Cookie: token=$TOKEN" \
    -H "Range: $range" \
    -i
}

if [ $# -eq 0 ]; then
  usage
fi

case "$1" in
  login)
    login
    ;;
  play)
    if [ $# -lt 2 ]; then
      echo "Error: missing media id"
      usage
    fi
    play_media "$2"
    ;;
  play-range)
    if [ $# -lt 3 ]; then
      echo "Error: missing media id or start position"
      usage
    fi
    play_media_range "$2" "$3" "$4"
    ;;
  *)
    echo "Error: unknown command '$1'"
    usage
    ;;
esac
