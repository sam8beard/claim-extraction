#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/..")"

# Ensure mc binary exists
if ! command -v mc &> /dev/null; then
    echo "Downloading MinIO client..."
    wget https://dl.min.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin/
fi

# Set alias for MinIO admin
echo "Setting up admin alias..."
mc alias set minio-admin http://localhost:9000 admin password

# Verify connection
echo "Verifying MinIO server..."
mc admin info minio-admin

# Add user for database operations
DB_USER="username"
DB_PASS="password"
echo "Adding MinIO user for DB operations..."
mc admin user add minio-admin $DB_USER $DB_PASS

# Verify user
mc admin user ls minio-admin

# Attach policy
mc admin policy attach minio-admin readwrite --user $DB_USER

# Connect with the new user
mc alias set s3 http://localhost:9000 $DB_USER $DB_PASS

# Create bucket
mc mb s3/claim-pipeline-docstore || echo "Bucket already exists"

echo "MinIO client setup complete."
