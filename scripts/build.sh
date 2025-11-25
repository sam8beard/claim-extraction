#!/usr/bin/env bash
set -e 

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/..")"

# echo "$PROJECT_ROOT"
ENTRY_POINT="$PROJECT_ROOT/cmd"

echo "Building claimex..."


sudo go build -o /usr/local/bin/claimex "$ENTRY_POINT"

echo "claimex built successfully" 

echo "Placed in /usr/local/bin/"

echo "Usage: claimex [query] [filecount]"



