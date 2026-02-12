#!/bin/bash

# Unified installation script for Straw

set -e # Exit on error

echo "---------------------------------------------------"
echo "  Straw — File Automation System Installation"
echo "---------------------------------------------------"

# 1. Run environment setup
echo "Step 1: Setting up environment and configuration..."
./scripts/setup.sh

# 2. Build and Install binaries
echo "Step 2: Building and installing binaries (may require sudo)..."
make build
sudo make install

# 3. Install the background service
echo "Step 3: Installing the background daemon service..."
./scripts/install_service.sh

echo ""
echo "---------------------------------------------------"
echo "🚀 Installation Complete!"
echo ""
echo "You can now run 'straw' in your terminal to"
echo "monitor your file automation in real-time."
echo "---------------------------------------------------"
