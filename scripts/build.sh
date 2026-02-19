#!/usr/bin/env bash

set -e
set -oo pipefail
set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Fetch the current API version
API_URL="${NUON_API_URL:-https://api.nuon.co}"
API_VERSION=$(curl -sS "$API_URL/version" | python3 -c "import sys,json; print(json.load(sys.stdin)['version'])")

if [[ -z "$API_VERSION" ]]; then
  echo "error: failed to fetch API version from $API_URL/version" >&2
  exit 1
fi

echo "API version: $API_VERSION"

# Download the latest spec
echo "Downloading spec from $API_URL/docs/doc.json..."
curl -sS "$API_URL/docs/doc.json" -o "$ROOT_DIR/spec/doc.json"

# Build
echo "Building nuon-ext-api..."
cd "$ROOT_DIR"
GOWORK=off go build -ldflags "-s -w" -o "$ROOT_DIR/nuon-ext-api" .

echo "Built: $ROOT_DIR/nuon-ext-api (API $API_VERSION)"
