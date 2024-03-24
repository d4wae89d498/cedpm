#!/bin/bash

SCRIPT_PATH="$(realpath "$0")"
SCRIPT_ROOT="$(dirname "$SCRIPT_PATH")"

#echo "Script root directory: $SCRIPT_ROOT"


#-------------------------------
#		      S Y S
#          A N A L Y S E
#-------------------------------

printf "Detecting c++ compiler... "

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
CXX=$(detect_executable_and_version clang++ g++)
CXXFLAGS="-O2 -Wall -Werror -Wextra"
CXXSTANDARDS=""
for standard in c++98 c++03 c++11 c++14 c++17 c++20 c++23; do
    if $CXX -xc++ -std=$standard $TEST_SRC -o $TEST_BIN >/dev/null 2>&1; then
        CXXSTANDARDS="$CXXSTANDARDS $standard"
        rm -f $TEST_BIN
	else
		echo "ERRPR $CXX -std=$standard test_prog.c $TEST_SRC -o $TEST_BIN"
		exit;
	fi
done
rm -f  $TEST_SRC

echo "done"

echo ""
echo "---------"
echo "CXX:			$CXX"
echo "CXXFLAGS:			$CXXFLAGS"
echo "CXXSTANDARDS:		$CXXSTANDARDS"
echo "---------"

#-------------------------------
#		C R E A T I N G
# . 		E N V
#-------------------------------
ENV_FILE=$ENV_DIR/c++-compiler.sh
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

export_env CXX
export_env CXXFLAGS
export_env CXXSTANDARDS

