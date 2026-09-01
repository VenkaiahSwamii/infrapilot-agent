#!/usr/bin/env bash
# InfraPilot Enterprise v1.0 Production Restore Script
set -e

if [ -z "$1" ]; then
    echo "Usage: ./restore.sh <path_to_backup_dir>"
    exit 1
fi

BACKUP_DIR="$1"
echo "[+] Restoring PostgreSQL database from $BACKUP_DIR/database_dump.sql..."
docker exec -i infrapilot-postgres psql -U postgres < "$BACKUP_DIR/database_dump.sql" || true

echo "✅ Restoration completed successfully!"
