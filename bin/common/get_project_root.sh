#!/bin/bash

project_root=$(pwd)
found=0

while [ "$project_root" != "/" ]; do
    if test -f "$project_root/Project" || test -f "$project_root/PklProject"; then
        echo "Found project in '$project_root'"
        found=1
        break
    fi
    project_root=$(dirname "$project_root")
done

if [ $found -ne 1 ]; then
	project_root="null"
fi
