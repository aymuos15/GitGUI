#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Linux)
    case "$ARCH" in
      x86_64)
        PLATFORM="linux-amd64"
        ;;
      aarch64)
        PLATFORM="linux-arm64"
        ;;
      *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
    esac
    ;;
  Darwin)
    case "$ARCH" in
      x86_64)
        PLATFORM="darwin-amd64"
        ;;
      arm64)
        PLATFORM="darwin-arm64"
        ;;
      *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
    esac
    ;;
  *)
    echo -e "${RED}Unsupported OS: $OS${NC}"
    exit 1
    ;;
esac

# Repository and version
REPO="aymuos15/GitGUI"
VERSION="${1:-latest}"
INSTALL_DIR="$HOME/bin"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Download URL
if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/v0.1.0/gg-$PLATFORM"
else
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/gg-$PLATFORM"
fi

echo "Installing gg ($PLATFORM)..."
echo "Downloading from: $DOWNLOAD_URL"

# Download binary
if ! curl -fsSL -L -o "$INSTALL_DIR/gg" "$DOWNLOAD_URL"; then
  echo -e "${RED}Failed to download gg${NC}"
  exit 1
fi

# Make executable
chmod +x "$INSTALL_DIR/gg"

# Check if ~/bin is in PATH
if [[ ":$PATH:" == *":$HOME/bin:"* ]]; then
  echo -e "${GREEN}✓ gg installed successfully to $INSTALL_DIR/gg${NC}"
else
  echo -e "${RED}Warning: $INSTALL_DIR is not in your PATH${NC}"
  echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc):"
  echo "  export PATH=\"\$HOME/bin:\$PATH\""
fi
