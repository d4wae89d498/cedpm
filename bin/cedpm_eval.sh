#!/bin/bash

cedpm_root=$(dirname "$(realpath "$0")")

source ${cedpm_root}/common/get_project_root.sh


#rm -f $project_root/PklProject $project_root/PklProject.deps.json

#ln -s $project_root/Project.pkl $project_root/PklProject
#ln -s $project_root/Project.deps.json $project_root/PklProject.deps.json

#if [ ! -f $project_root/.cedpm/PklProject.deps.json ]; then
#	echo "Error. Run npm install first."
#	exit 1;
#fi;

#source $project_root/.cedpm/env.sh

#cat $1 | preprocess > .tmp
#shift;
#cd pkl eval .tmp $@
if test $project_root != "null"; then
	source $project_root/.cedpm/env/*.sh

	content=$(${cedpm_root}/pkls/pkls $1 $project_root)

	echo $content;
	exit ;

	if ! echo "$content" | ${cedpm_root}/jmk/jmk; then
		echo "JSON parse error in: $content"
	fi;
else
	${cedpm_root}/pkls/pkls $1 #| ${cedpm_root}/jmk/jmk

#	${cedpm_root}/pkls/pkls $1
fi

#rm -f $project_root/PklProject $project_root/PklProject.deps.json
