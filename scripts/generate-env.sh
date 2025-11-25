#!/usr/bin/env bash
set -e

echo "Generating .env file..."

ENV_FILE=".env"

cat > $ENV_FILE <<EOL

# DB config
DATABASE_URL=postgres://username:password@postgres:5432/claimex-db?sslmode=disable
DB_HOST=postgres
DB_NAME=claimex-db
DB_USERNAME=username
DB_PASSWORD=password
DB_PORT=5432

# MinIO config
MINIO_ENDPOINT=http://minio:9000
MINIO_USER=admin
MINIO_PASSWORD=password
USE_SSL=False
EOL

echo ".env template generated at $ENV_FILE"

