#!/bin/bash

CONFIG_DIR="$HOME/.config/straw"
STATE_DIR="$HOME/.local/state/straw"
CONFIG_FILE="$CONFIG_DIR/config.toml"

echo "Setting up Straw environment..."

mkdir -p "$CONFIG_DIR"
mkdir -p "$STATE_DIR"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Creating default configuration at $CONFIG_FILE"
    cat > "$CONFIG_FILE" <<EOF
# Straw Configuration

[tui]
theme = "default"

[[watch]]
path = "$HOME/Downloads"
recursive = true

[[rules]]
name = "Organize PDFs"
enabled = true
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "$HOME/Documents/PDFs"

[[rules.actions]]
type = "shell"
target = "notify-send"
args = ["Straw", "Moved \$FILE to PDFs"]
EOF
    echo "Default config created. Please edit it to match your folders."
else
    echo "Config file already exists at $CONFIG_FILE"
fi

echo "Setup complete. Run 'make build' to compile the binaries."
