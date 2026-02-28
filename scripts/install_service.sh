#!/bin/bash

SERVICE_NAME="strawd.service"
SYSTEMD_DIR="$HOME/.config/systemd/user"
SERVICE_FILE="$SYSTEMD_DIR/$SERVICE_NAME"

# Find strawd binary
if command -v strawd &> /dev/null; then
    BIN_PATH=$(command -v strawd)
else
    echo "Error: 'strawd' not found in PATH."
    echo "Please ensure you have built and installed the project (sudo make install)."
    exit 1
fi

echo "Found strawd at: $BIN_PATH"

# Check if systemctl is available
if ! command -v systemctl >/dev/null 2>&1; then
    echo "Warning: 'systemctl' not found. Skipping automatic service installation."
    echo "You can still run 'strawd' manually to start the daemon."
    exit 0
fi

# Create systemd user directory if it doesn't exist
mkdir -p "$SYSTEMD_DIR"

# Generate the service file
echo "Generating $SERVICE_FILE..."

cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Straw File Automation Daemon
Documentation=https://github.com/raythurman2386/straw
After=network.target

[Service]
Type=simple
ExecStart=$BIN_PATH
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

# Environment variables can be added here if needed
# Environment=STRAW_LOG_LEVEL=debug

[Install]
WantedBy=default.target
EOF

echo "Installing service..."

# Stop existing service before updating
if systemctl --user is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
    echo "Stopping existing strawd service..."
    systemctl --user stop "$SERVICE_NAME"
    sleep 1
fi

systemctl --user daemon-reload
systemctl --user enable "$SERVICE_NAME"

# Start or restart the service
if systemctl --user is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
    systemctl --user restart "$SERVICE_NAME"
    echo "Restarted strawd service"
else
    systemctl --user start "$SERVICE_NAME"
    echo "Started strawd service"
fi

echo "---------------------------------------------------"
echo "✅ Straw daemon is now running as a user service!"
echo ""
echo "Manage it with:"
echo "  systemctl --user status strawd"
echo "  systemctl --user stop strawd"
echo "  systemctl --user restart strawd"
echo ""
echo "View logs with:"
echo "  journalctl --user -u strawd -f"
echo "---------------------------------------------------"
