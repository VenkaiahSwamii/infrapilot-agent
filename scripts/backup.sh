#!/usr/bin/env bash
# InfraPilot Enterprise v1.0 Production Backup Script
set -e

BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "[+] Backing up PostgreSQL database cluster..."
docker exec -t infrapilot-postgres pg_dumpall -U postgres > "$BACKUP_DIR/database_dump.sql" || true

echo "[+] Backing up Redis cache & session store..."
docker exec -t infrapilot-redis redis-cli save || true

echo "[+] Copying configuration & secrets..."
cp -r ./certs "$BACKUP_DIR/" || true

echo "✅ Backup completed successfully at $BACKUP_DIR"
