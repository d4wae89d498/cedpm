#!/bin/bash

SCRIPT_PATH="$(realpath "$0")"
SCRIPT_ROOT="$(dirname "$SCRIPT_PATH")"

#echo "Script root directory: $SCRIPT_ROOT"


#-------------------------------
#		      S Y S
#          A N A L Y S E
#-------------------------------

printf "Detecting environment... "
HOST_KERNEL_NAME="$(uname -s)"
HOST_KERNEL_VERSION="$(uname -r)"
HOST_ISA=$(uname -m)
HOST_OS_NAME="$(uname -s)"
case "${HOST_OS_NAME}" in
    Darwin)
        HOST_OS_NAME="MacOS"
        HOST_OS_VERSION="$(sw_vers -productVersion)"
        ;;
    Linux)
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            HOST_OS_NAME=$NAME
            HOST_OS_VERSION=$VERSION_ID
        elif [ -f /etc/lsb-release ]; then
            . /etc/lsb-release
            HOST_OS_NAME=$DISTRIB_ID
            HOST_OS_VERSION=$DISTRIB_RELEASE
        else
            HOST_OS_NAME="linux"
            HOST_OS_VERSION="unknown"
        fi
        ;;
    *)
        HOST_OS_NAME="unknown"
        HOST_OS_VERSION="unknown"
        ;;
esac

if test $(printf "\x01\x00" | od -An -t u2 | tr -d ' ') = "1"; then
    HOST_ENDIANNESS="little"
else
    HOST_ENDIANNESS="big"
fi

echo "done."

#-------------------------------
#		S Y S   I N F O
# . 	    D U M P
#-------------------------------

echo "---------"
echo "HOST_ISA:		$HOST_ISA"
echo "HOST_ENDIANNESS:	$HOST_ENDIANNESS"
echo "HOST_KERNEL_NAME:	$HOST_KERNEL_NAME"
echo "HOST_KERNEL_VERSION:	$HOST_KERNEL_VERSION"
echo "HOST_OS_NAME:		$HOST_OS_NAME"
echo "HOST_OS_VERSION:	$HOST_OS_VERSION"
echo "---------"

#-------------------------------
#		C R E A T I N G
# . 		E N V
#-------------------------------
ENV_FILE=$ENV_DIR/system.sh
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
export_env HOST_ISA
export_env HOST_ENDIANNESS
export_env HOST_KERNEL_NAME
export_env HOST_KERNEL_VERSION
export_env HOST_OS_NAME
export_env HOST_OS_VERSION

TARGET_ISA=$HOST_ISA
TARGET_ENDIANNESS=$HOST_ENDIANNESS
TARGET_KERNEL_NAME=$HOST_KERNEL_NAME
TARGET_KERNEL_VERSION=$HOST_KERNEL_VERSION
TARGET_OS_NAME=$HOST_OS_NAME
TARGET_OS_VERSION=$HOST_OS_VERSION

export_env TARGET_ISA
export_env TARGET_ENDIANNESS
export_env TARGET_KERNEL_NAME
export_env TARGET_KERNEL_VERSION
export_env TARGET_OS_NAME
export_env TARGET_OS_VERSION

