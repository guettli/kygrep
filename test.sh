#!/usr/bin/env bash
# Test script for kygrep
# Uses bash strict mode for better error handling

set -euo pipefail
IFS=$'\n\t'

echo "Running tests..."
go test -v ./...
echo "All tests passed!"
