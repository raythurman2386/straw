#!/bin/bash

# Unified installation script for Straw

set -e # Exit on error

echo "---------------------------------------------------"
echo "  Straw — File Automation System Installation"
echo "---------------------------------------------------"

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed. Please install Go 1.25+ from https://go.dev/dl/"; exit 1; }

GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | head -1 | sed 's/go//')
GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 25 ]); then
    echo "Error: Go version $GO_VERSION is too old. Requires Go 1.25 or later."
    echo "Please install Go 1.25+ from https://go.dev/dl/"
    exit 1
fi

echo "✓ Go $GO_VERSION found"

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 1. Run environment setup
echo "Step 1: Setting up environment and configuration..."
"$SCRIPT_DIR/setup.sh"

# 2. Build binaries as current user (not root)
echo "Step 2: Building binaries..."
cd "$PROJECT_ROOT"
make build

# 3. Install binaries (requires sudo)
echo "Step 3: Installing binaries to /usr/local/bin (requires sudo)..."
sudo make install

# 4. Install the background service
echo "Step 4: Installing the background daemon service..."
"$SCRIPT_DIR/install_service.sh"

echo ""
echo "---------------------------------------------------"
echo "🚀 Installation Complete!"
echo ""
echo "You can now run 'straw' in your terminal to"
echo "monitor your file automation in real-time."
echo "---------------------------------------------------"
