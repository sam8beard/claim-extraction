#!/usr/bin/env bash
set -e

# Load .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found. Please run generate-env.sh first."
    exit 1
fi

# Run the create_documents.sql script inside the Postgres container
docker exec -i postgres psql -U "$DB_USERNAME" -d "$DB_NAME" -f /pg/schema/create_documents.sql

echo "Postgres documents table created."
