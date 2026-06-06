#!/usr/bin/env sh
set -eu

APP_NAME="vw"

YES=0
PURGE=0
PREFIX="${HOME}/.local"

usage() {
  cat <<EOF2
Usage: scripts/uninstall.sh [options]

Options:
  --prefix DIR   Installation prefix. Default: \$HOME/.local
  --purge        Also remove vw config, data, cache, and state directories
  --yes, -y      Do not prompt for confirmation
  --help, -h     Show this help

Examples:
  scripts/uninstall.sh
  scripts/uninstall.sh --purge
  scripts/uninstall.sh --purge --yes
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
  --purge)
    PURGE=1
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

remove_path() {
  path="$1"

  if [ -e "$path" ] || [ -L "$path" ]; then
    rm -rf "$path"
    echo "removed: $path"
  fi
}

echo "Uninstalling ${APP_NAME}"

remove_path "${PREFIX}/bin/vw"

if confirm "Remove bundled bw from ${PREFIX}/bin/bw if present?"; then
  remove_path "${PREFIX}/bin/bw"
fi

remove_path "${HOME}/.local/share/vw/bin/bw"

if [ "$PURGE" -eq 1 ]; then
  echo "Purging vw user data"

  remove_path "${HOME}/.config/vw"
  remove_path "${HOME}/.cache/vw"
  remove_path "${HOME}/.local/state/vw"
  remove_path "${HOME}/.local/share/vw"

  echo "Note: if you stored BW_SESSION in your OS keyring, run 'vw lock' before uninstalling,"
  echo "or remove the 'vw' entry manually from your system keychain/keyring."
else
  echo "Kept user config/data. Re-run with --purge to remove:"
  echo "  ${HOME}/.config/vw"
  echo "  ${HOME}/.cache/vw"
  echo "  ${HOME}/.local/state/vw"
  echo "  ${HOME}/.local/share/vw"
fi

echo "Done."
