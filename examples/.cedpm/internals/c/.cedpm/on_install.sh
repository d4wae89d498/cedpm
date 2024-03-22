#!/bin/bash

SCRIPT_PATH="$(realpath "$0")"
SCRIPT_ROOT="$(dirname "$SCRIPT_PATH")"

#echo "Script root directory: $SCRIPT_ROOT"


#-------------------------------
#		      S Y S
#          A N A L Y S E
#-------------------------------

printf "Detecting c compiler... "

detect_executable_and_version() {
    for cmd in "$@"; do
        if which "$cmd" >/dev/null 2>&1; then
            echo $(which $cmd)
			return
        fi
    done
}
TEST_SRC=$SCRIPT_ROOT/.test_prog.c
TEST_BIN=$SCRIPT_ROOT/.test.out
echo "int main() { return 0; }" > $TEST_SRC
CC=$(detect_executable_and_version clang gcc tinycc)
CFLAGS="-O2 -Wall -Werror -Wextra"
CSTANDARDS=""
for standard in c89 c99 c11 c17; do
    if $CC -std=$standard  $TEST_SRC -o $TEST_BIN >/dev/null 2>&1; then
        CSTANDARDS="$CSTANDARDS $standard"
        rm -f $TEST_BIN
    fi
done
rm -f  $TEST_SRC

echo "done. "

echo ""
echo "---------"
echo "CC:			$CC"
echo "CFLAGS:			$CFLAGS"
echo "CSTANDARDS:		$CSTANDARDS"
echo "---------"

#-------------------------------
#		C R E A T I N G
# . 		E N V
#-------------------------------

ENV_FILE=$ENV_DIR/c-compiler.sh
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

export_env CC
export_env CFLAGS
export_env CSTANDARDS
