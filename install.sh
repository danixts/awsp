#!/usr/bin/env bash
set -e

REPO="${AWS_PROFILE_REPO:-danixts/awsp}"
VERSION="${AWS_PROFILE_VERSION:-latest}"
INSTALL_DIR="${HOME}/.local/bin"
BIN_NAME="aws-profile"
CONFIG_AWSP="${HOME}/.config/awsp"

case "$(uname -s)" in
  Linux)   OS="linux";;
  Darwin)  OS="darwin";;
  MINGW*)  OS="windows";;
  *)       echo "Unsupported OS: $(uname -s)"; exit 1;;
esac

case "$(uname -m)" in
  x86_64|amd64) ARCH="amd64";;
  aarch64|arm64) ARCH="arm64";;
  *)            echo "Unsupported arch: $(uname -m)"; exit 1;;
esac

if [ "$OS" = "windows" ]; then
  ASSET="${BIN_NAME}-${OS}-${ARCH}.exe"
  OUT="${INSTALL_DIR}/${BIN_NAME}.exe"
else
  ASSET="${BIN_NAME}-${OS}-${ARCH}"
  OUT="${INSTALL_DIR}/${BIN_NAME}"
fi

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

echo "Installing aws-profile (${OS}-${ARCH}) from ${REPO}..."

mkdir -p "$INSTALL_DIR"
rm -f "$OUT" "${OUT}.exe" 2>/dev/null || true
rm -f "${CONFIG_AWSP}/completion.zsh" "${CONFIG_AWSP}/completion.bash" 2>/dev/null || true

TMP=$(mktemp)
if command -v curl &>/dev/null; then
  curl -sSL -o "$TMP" "$URL"
elif command -v wget &>/dev/null; then
  wget -q -O "$TMP" "$URL"
else
  echo "Need curl or wget"; exit 1
fi
mv "$TMP" "$OUT"
chmod +x "$OUT" 2>/dev/null || true

export PATH="${INSTALL_DIR}:${PATH}"
echo "Binary: $OUT"
echo "Configuring shell (awsp + completion)..."
"$OUT" install
