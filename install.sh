#!/usr/bin/env bash
set -e

# ccconfig One-Click Installer
# Usage: curl -fsSL https://jiangtao.vercel.app/install.sh | bash
#        or: curl -fsSL https://jiangtao.vercel.app/install.sh | bash

REPO="jiangtao/cc-config"
GITHUB_BASE="https://github.com/${REPO}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}➜${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Darwin)
        OS="darwin"
        ;;
    Linux)
        OS="linux"
        ;;
    MINGW*|MSYS*|CYGWIN*)
        error "Windows not supported yet. Please download from ${GITHUB_BASE}/releases"
        exit 1
        ;;
    *)
        error "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    i686|i386)
        ARCH="386"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    *)
        error "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest version (with fallback for API rate limiting)
info "Detecting latest version..."

# Method 1: Try GitHub API (may hit rate limits)
LATEST_VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

# Method 2: Fallback - parse from GitHub releases page redirect
if [ -z "$LATEST_VERSION" ]; then
    REDIRECT_URL=$(curl -fsSL "https://github.com/${REPO}/releases/latest" -w "%{url_effective}" -o /dev/null 2>/dev/null)
    LATEST_VERSION=$(echo "$REDIRECT_URL" | sed -E 's|.*/tag/(v[0-9.]+)$|\1|')
fi

# Method 3: Last resort - use 'latest' redirect directly
if [ -z "$LATEST_VERSION" ] || [[ "$LATEST_VERSION" == *"http"* ]]; then
    info "Using latest release redirect (API rate limited)..."
    BINARY_URL="${GITHUB_BASE}/releases/latest/download/ccconfig-${OS}-${ARCH}"
    USE_LATEST_REDIRECT=true
else
    info "Latest version: ${LATEST_VERSION}"
    BINARY_URL="${GITHUB_BASE}/releases/download/${LATEST_VERSION}/ccconfig-${OS}-${ARCH}"
    USE_LATEST_REDIRECT=false
fi

BINARY_NAME="ccconfig"
TEMP_DIR=$(mktemp -d)

info "Downloading ccconfig for ${OS}-${ARCH} from:"
echo "  ${BINARY_URL}"
echo ""

if ! curl -fSL "$BINARY_URL" -o "${TEMP_DIR}/${BINARY_NAME}"; then
    error "Download failed!"
    echo ""
    echo "The binary URL was: ${BINARY_URL}"
    echo ""
    echo "Possible reasons:"
    echo "  1. Release hasn't been created yet (GitHub Actions still running)"
    echo "  2. Release doesn't have a binary for ${OS}-${ARCH}"
    echo "  3. Network connectivity issues"
    echo ""
    echo "To check the release status:"
    echo "  - Visit: ${GITHUB_BASE}/releases"
    echo "  - Or check GitHub Actions: ${GITHUB_BASE}/actions"
    echo ""
    echo "If this is a fresh release, wait a few minutes for GitHub Actions to complete."
    rm -rf "$TEMP_DIR"
    exit 1
fi

chmod +x "${TEMP_DIR}/${BINARY_NAME}"

# Determine install location
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    # Check if we can use sudo
    if command -v sudo >/dev/null 2>&1; then
        USE_SUDO=true
    else
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        info "Installing to ${INSTALL_DIR} (add to PATH: export PATH=\"\$HOME/.local/bin:\$PATH\")"
        USE_SUDO=false
    fi
else
    USE_SUDO=false
fi

# Install binary
info "Installing to ${INSTALL_DIR}..."
if [ "$USE_SUDO" = true ]; then
    sudo mv "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    mv "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

# Cleanup
rm -rf "$TEMP_DIR"

# Verify installation
if command -v ccconfig >/dev/null 2>&1; then
    INSTALLED_VERSION=$(ccconfig --version 2>/dev/null || echo "unknown")
    echo ""
    info "Successfully installed ccconfig!"
    echo "   Version: ${INSTALLED_VERSION}"
    echo "   Location: ${INSTALL_DIR}/${BINARY_NAME}"
    echo ""
    echo "Quick start:"
    echo "  ccconfig backup --repo ~/cc-config"
    echo ""
else
    warn "Installation completed, but 'ccconfig' is not in PATH."
    echo "   Add the following to your ~/.bashrc or ~/.zshrc:"
    echo "   export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
