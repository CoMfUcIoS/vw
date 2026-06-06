#!/usr/bin/env sh
set -eu

APP_NAME="vw"
PREFIX="${PREFIX:-${HOME}/.local}"
INSTALL_BW=1
YES=0

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
BUNDLE_DIR="$(CDPATH= cd -- "${SCRIPT_DIR}/.." && pwd)"

usage() {
  cat <<EOF2
Usage: scripts/install.sh [options]

Options:
  --prefix DIR     Installation prefix. Default: \$HOME/.local
  --no-bw          Do not install bundled bw even if present
  --yes, -y        Do not prompt for confirmation
  --help, -h       Show this help

Examples:
  scripts/install.sh
  scripts/install.sh --prefix /usr/local
  scripts/install.sh --no-bw
EOF2
}

while [ "$#" -gt 0 ]; do
  case "$1" in
  --prefix)
    if [ "$#" -lt 2 ]; then
      echo "error: --prefix requires a value" >&2
      exit 2
    fi
    PREFIX="$2"
    shift 2
    ;;
  --no-bw)
    INSTALL_BW=0
    shift
    ;;
  --yes | -y)
    YES=1
    shift
    ;;
  --help | -h)
    usage
    exit 0
    ;;
  *)
    echo "error: unknown argument: $1" >&2
    usage >&2
    exit 2
    ;;
  esac
done

confirm() {
  message="$1"

  if [ "$YES" -eq 1 ]; then
    return 0
  fi

  printf "%s [y/N] " "$message"
  read -r answer || answer=""
  case "$answer" in
  y | Y | yes | YES)
    return 0
    ;;
  *)
    return 1
    ;;
  esac
}

install_file() {
  src="$1"
  dst="$2"

  if [ ! -f "$src" ]; then
    echo "error: missing file: $src" >&2
    exit 1
  fi

  mkdir -p "$(dirname -- "$dst")"
  cp "$src" "$dst"
  chmod +x "$dst"
  echo "installed: $dst"
}

VW_SRC="${BUNDLE_DIR}/bin/vw"
BW_SRC="${BUNDLE_DIR}/bin/bw"

if [ ! -f "$VW_SRC" ]; then
  echo "error: bundled vw binary not found at ${VW_SRC}" >&2
  echo "Run this script from an extracted vw release bundle." >&2
  exit 1
fi

echo "Installing ${APP_NAME} to ${PREFIX}/bin"

if [ "$YES" -eq 0 ]; then
  confirm "Continue installation?" || {
    echo "aborted"
    exit 1
  }
fi

install_file "$VW_SRC" "${PREFIX}/bin/vw"

if [ "$INSTALL_BW" -eq 1 ]; then
  if [ -f "$BW_SRC" ]; then
    install_file "$BW_SRC" "${PREFIX}/bin/bw"
  else
    echo "bundled bw not found; skipping bw install"
  fi
fi

cat <<EOF2

Installation complete.

Make sure this is in your PATH:

  ${PREFIX}/bin

Try:

  vw doctor
  vw setup --server https://vaultwarden.example.com

EOF2
