#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/..")"

# Paths to Python directories
PYTHON_DIRS=(
	"$PROJECT_ROOT/internal/conversion/python" 
	"$PROJECT_ROOT/internal/processing/python"
)

for dir in "${PYTHON_DIRS[@]}"; do
    echo "Setting up virtual environment in $dir..."
    python3 -m venv "$dir/.venv"
    source "$dir/.venv/bin/activate"
    if [ -f "$dir/requirements.txt" ]; then
        pip install --upgrade pip
        pip install -r "$dir/requirements.txt"
        echo "Installed dependencies from $dir/requirements.txt"
    else
        echo "No requirements.txt found in $dir"
    fi
    deactivate
done

echo "All Python virtual environments created and dependencies installed."
