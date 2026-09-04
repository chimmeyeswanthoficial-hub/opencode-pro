#!/usr/bin/env bash
set -e

echo "Building OpenCode Pro binary..."
go build -ldflags="-s -w" -o opencode .
echo "Build complete! Run './opencode' or './opencode optimize \"prompt\"' to get started."
