#!/usr/bin/env bash
# ==============================================================================
# InfraPilot Enterprise Native Installer (No Docker Required)
# ==============================================================================
set -euo pipefail

INFRAPILOT_VERSION="1.0.0"
INFRAPILOT_HOME="${INFRAPILOT_HOME:-$HOME/.infrapilot}"

log_info() { echo -e "\033[0;32m[INFO]\033[0m $1"; }
log_error() { echo -e "\033[0;31m[ERROR]\033[0m $1"; }

log_info "Installing InfraPilot Enterprise (Native Bare-Metal Mode)..."
mkdir -p "$INFRAPILOT_HOME"

# Check Go
if ! command -v go &> /dev/null; then
    log_error "Go is required. Please install Go 1.21+."
    exit 1
fi

# Check Node
if ! command -v npm &> /dev/null; then
    log_error "Node.js/npm is required. Please install Node.js."
    exit 1
fi

# Build components
log_info "Building Go Backend API..."
cd backend
go build -o "$INFRAPILOT_HOME/infrapilot-server" ./cmd/server/main.go

log_info "Building React Frontend..."
cd ../frontend
npm install --silent
npm run build

log_info "Building Agent..."
cd ../agent
go build -o "$INFRAPILOT_HOME/infrapilot-agent" ./cmd/agent/main.go || true

log_info "Installation complete in $INFRAPILOT_HOME!"
echo "To start backend: $INFRAPILOT_HOME/infrapilot-server"
echo "To start frontend: cd frontend && npm run dev"