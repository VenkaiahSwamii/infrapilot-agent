#!/usr/bin/env bash
# InfraPilot Enterprise v1.0 Production Linux Installer
set -e

echo "========================================================"
echo "   🚀 Installing InfraPilot Enterprise v1.0 Platform    "
echo "========================================================"

# Check docker & docker compose
if ! command -v docker &> /dev/null; then
    echo "[!] Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com | sh
fi

echo "[+] Pulling InfraPilot Enterprise v1.0 production containers..."
docker compose pull || true

echo "[+] Starting InfraPilot Enterprise services..."
docker compose up -d

echo "========================================================"
echo "   ✅ InfraPilot Enterprise v1.0 Successfully Deployed! "
echo "   Web Dashboard:  http://localhost:3000                 "
echo "   API Endpoint:   http://localhost:8080/api/v1/health   "
echo "========================================================"