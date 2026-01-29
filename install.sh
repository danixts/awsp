#!/usr/bin/env bash
set -e
REPO="${AWS_PROFILE_REPO:-danixts/awsp}"
VERSION="${AWS_PROFILE_VERSION:-latest}"
INSTALL_DIR="${HOME}/.local/bin"
BIN_NAME="aws-profile"

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
else
  ASSET="${BIN_NAME}-${OS}-${ARCH}"
fi

URL="https://github.com/${REPO}/releases/${VERSION}/download/${ASSET}"
if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
fi

if [ "$OS" = "windows" ]; then
  OUT="${INSTALL_DIR}/${BIN_NAME}.exe"
else
  OUT="${INSTALL_DIR}/${BIN_NAME}"
fi

echo "Installing aws-profile (${OS}-${ARCH}) from ${REPO}..."
mkdir -p "$INSTALL_DIR"
if command -v curl &>/dev/null; then
  curl -sSL -o "$OUT" "$URL"
elif command -v wget &>/dev/null; then
  wget -q -O "$OUT" "$URL"
else
  echo "Need curl or wget"; exit 1
fi
chmod +x "$OUT" 2>/dev/null || true

export PATH="${INSTALL_DIR}:${PATH}"
echo "Binary: $OUT"
echo "Configuring shell (awsp + completion)..."
"$OUT" install
