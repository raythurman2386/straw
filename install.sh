#!/bin/sh
# Straw installer script - downloads pre-built binaries from GitHub releases
# Usage: curl -fsSL https://raw.githubusercontent.com/raythurman2386/straw/main/install.sh | sh

set -e

REPO="raythurman2386/straw"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture names
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="armv7"
        ;;
    i386|i686)
        ARCH="386"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Validate OS (Linux and macOS only for now)
case "$OS" in
    linux|darwin)
        ;;
    *)
        echo "Error: Unsupported operating system: $OS"
        echo "Currently supported: Linux, macOS"
        exit 1
        ;;
esac

echo "Detected: $OS $ARCH"

# Get latest release version
if [ -z "$VERSION" ]; then
    echo "Fetching latest release..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest version"
        exit 1
    fi
fi

echo "Installing straw $VERSION..."

# GoReleaser strips the "v" prefix in archive filenames
VERSION_NUM="${VERSION#v}"

# Construct download URL
FILENAME="straw_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download binary archive
echo "Downloading from GitHub releases..."
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$URL" -o "$TMP_DIR/$FILENAME" 2>&1
elif command -v wget >/dev/null 2>&1; then
    wget -q "$URL" -O "$TMP_DIR/$FILENAME" 2>&1
else
    echo "Error: curl or wget is required"
    exit 1
fi

# Extract
echo "Extracting..."
cd "$TMP_DIR"
tar -xzf "$FILENAME"

# Verify binaries exist
if [ ! -f "straw" ] || [ ! -f "strawd" ]; then
    echo "Error: Expected binaries not found in archive"
    exit 1
fi

# Install binaries
if [ -w "$INSTALL_DIR" ]; then
    echo "Installing to $INSTALL_DIR..."
    cp straw strawd "$INSTALL_DIR/"
else
    echo "Installing to $INSTALL_DIR (requires sudo)..."
    sudo cp straw strawd "$INSTALL_DIR/"
fi

# Setup config directory and example config
echo "Setting up configuration..."
CONFIG_DIR="${HOME}/.config/straw"
mkdir -p "$CONFIG_DIR"

if [ -f "config.example.toml" ] && [ ! -f "$CONFIG_DIR/config.toml" ]; then
    cp "config.example.toml" "$CONFIG_DIR/config.toml"
    echo "Created default config at ~/.config/straw/config.toml"
fi

# Install systemd service on Linux
if [ "$OS" = "linux" ] && command -v systemctl >/dev/null 2>&1; then
    echo "Setting up systemd service..."
    
    SERVICE_DIR="${HOME}/.config/systemd/user"
    mkdir -p "$SERVICE_DIR"
    
    cat > "$SERVICE_DIR/strawd.service" << EOF
[Unit]
Description=Straw File Automation Daemon
After=network.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/strawd
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
EOF
    
    systemctl --user daemon-reload
    systemctl --user enable strawd.service
    systemctl --user start strawd.service
    
    echo "Started strawd service"
fi

echo ""
echo "Installation complete!"
echo ""
echo "Binaries installed:"
echo "  - straw (TUI client)"
echo "  - strawd (daemon)"
echo ""
echo "Configuration: ~/.config/straw/config.toml"

if [ "$OS" = "linux" ] && command -v systemctl >/dev/null 2>&1; then
    echo "Service status: systemctl --user status strawd"
fi

echo ""
echo "Run 'straw' to start the TUI client"
