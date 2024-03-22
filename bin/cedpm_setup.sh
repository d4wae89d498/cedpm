#!/bin/bash

SCRIPT_PATH="$(realpath "$0")"
SCRIPT_ROOT="$(dirname "$SCRIPT_PATH")"
PKL_MIN_VERSION=0.25

#if which pkl >/dev/null 2>&1; then
#	echo "Found Apple PKL in path."
#	PKL=$(which pkl)
#else
	PKL_LATEST=$(curl -s  https://api.github.com/repos/apple/pkl/releases | grep tag_name | grep "\"$PKL_MIN_VERSION" | head -n 1 | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
	printf "Downloading Apple PKL $PKL_LATEST ... "
	if [ "$HOST_OS_NAME" = "Linux" ]; then
		if [ "$HOST_ISA" = "x86_64" ]; then
			PKL_EXECUTABLE_NAME="pkl-linux-amd64"
		elif [ "$HOST_ISA" = "aarch64" ]; then
			PKL_EXECUTABLE_NAME="pkl-linux-aarch64"
		fi
	elif [ "$HOST_OS_NAME" = "MacOS" ]; then
		if [ "$HOST_ISA" = "x86_64" ]; then
			PKL_EXECUTABLE_NAME="pkl-macos-amd64"
		elif [ "$HOST_ISA" = "arm64" ]; then
			PKL_EXECUTABLE_NAME="pkl-macos-aarch64"
		fi
	elif [ "$HOST_OS_NAME" = "Linux" ] && grep -q 'alpine' /etc/os-release; then
		if [ "$HOST_ISA" = "x86_64" ]; then
			PKL_EXECUTABLE_NAME="pkl-alpine-linux-amd64"
		fi
	fi
	if [ -z "$PKL_EXECUTABLE_NAME" ]; then
		echo "Unsupported OS or architecture for Apple PKL. Please install it manually and set in your PATH."
		exit 1;
	fi
	PKL=$SCRIPT_ROOT/$PKL_EXECUTABLE_NAME
	wget -q https://github.com/apple/pkl/releases/download/$PKL_LATEST/$PKL_EXECUTABLE_NAME -O $PKL
	chmod +x $PKL
	echo "done"
#fi;

#-------------------------------
#		C R E A T I N G
# . 		E N V
#-------------------------------
ENV_FILE=$SCRIPT_ROOT/env.sh
rm -f $ENV_FILE
export_env() {
    local var_name="$1"
    local var_value="${!var_name}"
    local file_path="$ENV_FILE"
    if [ -z "$var_value" ]; then
        echo "Notice: Variable '$var_name' is unset or set to empty."
    fi
    echo "export $var_name=\"$var_value\"" >> "$file_path"
}

export_env PKL
