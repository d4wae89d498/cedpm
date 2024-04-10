#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PARENT_DIR=$(dirname "$SCRIPT_DIR")
TARGET="$PARENT_DIR/cedpm_skeleton"

for dir in "$SCRIPT_DIR"/*; do
  if [ -d "$dir" ]; then
    LINK_PATH="$dir/.cedpm"

	ln -s "$TARGET" "$LINK_PATH"
	echo "Created symlink in $LINK_PATH"
  fi
done
