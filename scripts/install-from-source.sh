#!/bin/bash

# Unified installation script for Straw

set -e # Exit on error

echo "---------------------------------------------------"
echo "  Straw — File Automation System Installation"
echo "---------------------------------------------------"

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 1. Run environment setup
echo "Step 1: Setting up environment and configuration..."
"$SCRIPT_DIR/setup.sh"

# 2. Build and Install binaries
echo "Step 2: Building and installing binaries (may require sudo)..."
cd "$PROJECT_ROOT"
sudo make install

# 3. Install the background service
echo "Step 3: Installing the background daemon service..."
"$SCRIPT_DIR/install_service.sh"

echo ""
echo "---------------------------------------------------"
echo "🚀 Installation Complete!"
echo ""
echo "You can now run 'straw' in your terminal to"
echo "monitor your file automation in real-time."
echo "---------------------------------------------------"
