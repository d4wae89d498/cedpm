#!/bin/bash
############
script_root=$(dirname "$(realpath "$0")")
############

source ${script_root}/get_project_root.sh

if [ $found -ne 1 ]; then
    echo "No '$filename' was not found in any parent directory. Run a cedpm init first."
	exit 1;
fi
