#!/bin/bash

SERVICE_NAME="strawd.service"
SYSTEMD_DIR="$HOME/.config/systemd/user"
SERVICE_FILE="$SYSTEMD_DIR/$SERVICE_NAME"

echo "Stopping and disabling Straw daemon service..."

if systemctl --user is-active --quiet "$SERVICE_NAME"; then
    systemctl --user stop "$SERVICE_NAME"
fi

if systemctl --user is-enabled --quiet "$SERVICE_NAME"; then
    systemctl --user disable "$SERVICE_NAME"
fi

if [ -f "$SERVICE_FILE" ]; then
    echo "Removing service file: $SERVICE_FILE"
    rm "$SERVICE_FILE"
    systemctl --user daemon-reload
fi

echo "✅ Service uninstalled successfully."
