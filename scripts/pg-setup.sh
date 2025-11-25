#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/..")"

# Load .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found. Please run generate-env.sh first."
    exit 1
fi

# Run the create_documents.sql script inside the Postgres container
docker exec -i postgres psql -U "$DB_USERNAME" -d "$DB_NAME" -f "$PROJECT_ROOT"/services/pg/schema/create_documents.sql

echo "Postgres documents table created."
