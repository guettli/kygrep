#!/usr/bin/env bash
# Build script for kygrep
# Uses bash strict mode for better error handling

set -euo pipefail
IFS=$'\n\t'

echo "Building kygrep..."
go build -o kygrep main.go
echo "Build complete. Binary: ./kygrep"
