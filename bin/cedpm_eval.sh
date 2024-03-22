#!/bin/bash

filename="Project.pkl"
project_root=$(pwd)
found=0

while [ "$project_root" != "/" ]; do
    if [ -f "$project_root/$filename" ]; then
        echo "Found '$filename' in '$project_root'"
        found=1
        break
    fi
    project_root=$(dirname "$project_root")
done

if [ $found -ne 1 ]; then
    echo "No '$filename' was not found in any parent directory. Run a cedpm init first."
	exit 1;
fi

if test -f "$project_root/PklProject" && ! test -L "$project_root/PklProject"; then
	echo "Error: a non-symlink PklProject is present before Project.pkl."
	exit 1;
fi;

rm -f $project_root/PklProject $project_root/PklProject.deps.json

ln -s $project_root/Project.pkl $project_root/PklProject
ln -s $project_root/Project.deps.json $project_root/PklProject.deps.json

#if [ ! -f $project_root/.cedpm/PklProject.deps.json ]; then
#	echo "Error. Run npm install first."
#	exit 1;
#fi;

#source $project_root/.cedpm/env.sh

cat $1 | preprocess > .tmp
shift;
pkl eval .tmp $@

rm -f $project_root/PklProject $project_root/PklProject.deps.json
