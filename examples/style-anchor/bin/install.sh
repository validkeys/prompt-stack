#!/usr/bin/env sh
# simple installer for the single-binary example
set -e
BIN=bin/mytool
DEST=${DEST:-${HOME}/.local/bin}

if [ ! -f "$BIN" ]; then
	printf "Binary not found: %s\n" "$BIN" >&2
	exit 1
fi

mkdir -p "$DEST"
cp "$BIN" "$DEST/"
printf "Installed %s to %s\n" "$BIN" "$DEST"
