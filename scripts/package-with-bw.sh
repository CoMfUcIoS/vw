#!/usr/bin/env sh
set -eu

APP_NAME="vw"
VERSION="${VERSION:-dev}"
OS="${OS:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
ARCH="${ARCH:-$(uname -m)}"

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
DIST_DIR="${ROOT_DIR}/dist"
BIN_DIR="${ROOT_DIR}/bin"

case "$ARCH" in
x86_64 | amd64)
  ARCH="amd64"
  ;;
arm64 | aarch64)
  ARCH="arm64"
  ;;
*)
  echo "error: unsupported architecture: $ARCH" >&2
  exit 1
  ;;
esac

case "$OS" in
linux)
  BW_OS="linux"
  ARCHIVE_EXT="tar.gz"
  ;;
darwin)
  BW_OS="macos"
  ARCHIVE_EXT="tar.gz"
  ;;
*)
  echo "error: unsupported OS for bundled bw packaging: $OS" >&2
  exit 1
  ;;
esac

VW_BIN="${BIN_DIR}/vw"
BW_BIN="${BW_PATH:-${BIN_DIR}/bw}"

if [ ! -x "$VW_BIN" ]; then
  echo "error: vw binary not found at ${VW_BIN}" >&2
  echo "Run: make build" >&2
  exit 1
fi

if [ ! -x "$BW_BIN" ]; then
  echo "error: bw binary not found at ${BW_BIN}" >&2
  echo "Set BW_PATH=/path/to/bw or place bw at ${BIN_DIR}/bw" >&2
  echo "You can also run: vw bootstrap-bw" >&2
  exit 1
fi

PACKAGE_NAME="${APP_NAME}-${VERSION}-${OS}-${ARCH}"
STAGE_DIR="${DIST_DIR}/${PACKAGE_NAME}"

rm -rf "$STAGE_DIR"
mkdir -p "${STAGE_DIR}/bin"
mkdir -p "${STAGE_DIR}/scripts"

cp "$VW_BIN" "${STAGE_DIR}/bin/vw"
cp "$BW_BIN" "${STAGE_DIR}/bin/bw"

cp "${ROOT_DIR}/README.md" "${STAGE_DIR}/README.md"

if [ -f "${ROOT_DIR}/LICENSE" ]; then
  cp "${ROOT_DIR}/LICENSE" "${STAGE_DIR}/LICENSE"
fi

if [ -f "${ROOT_DIR}/NOTICE" ]; then
  cp "${ROOT_DIR}/NOTICE" "${STAGE_DIR}/NOTICE"
fi

cp "${ROOT_DIR}/scripts/install.sh" "${STAGE_DIR}/scripts/install.sh"
cp "${ROOT_DIR}/scripts/uninstall.sh" "${STAGE_DIR}/scripts/uninstall.sh"

chmod +x "${STAGE_DIR}/bin/vw"
chmod +x "${STAGE_DIR}/bin/bw"
chmod +x "${STAGE_DIR}/scripts/install.sh"
chmod +x "${STAGE_DIR}/scripts/uninstall.sh"

cat >"${STAGE_DIR}/VERSION" <<EOF2
${VERSION}
EOF2

cat >"${STAGE_DIR}/BUNDLE" <<EOF2
name=${APP_NAME}
version=${VERSION}
os=${OS}
arch=${ARCH}
bw_os=${BW_OS}
contains_bw=true
EOF2

mkdir -p "$DIST_DIR"

(
  cd "$DIST_DIR"

  case "$ARCHIVE_EXT" in
  tar.gz)
    tar -czf "${PACKAGE_NAME}.tar.gz" "$PACKAGE_NAME"
    ;;
  *)
    echo "error: unsupported archive extension: ${ARCHIVE_EXT}" >&2
    exit 1
    ;;
  esac
)

echo "Created ${DIST_DIR}/${PACKAGE_NAME}.tar.gz"
