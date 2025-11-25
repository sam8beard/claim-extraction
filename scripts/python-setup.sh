#!/usr/bin/env bash
set -e

# Paths to your Python directories
PYTHON_DIRS=("internal/conversion/python" "internal/processing/python")

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
