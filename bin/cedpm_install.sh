# sip first arg

if [ $3 ]; then
	echo "Too many args."
	exit ;
fi;

if ! test -f "$(pwd)/Project.pkl";
then
	echo "Error: no Project.pkl file found. Please run cedpm init before."
	exit 1;
fi;

if test -f "$(pwd)/PklProject" && ! test -L "$(pwd)/PklProject"; then
	echo "Error: a non-symlink PklProject is present."
	exit 1;
fi;

mkdir -p $(pwd)/.cedpm
rm -f $(pwd)/PklProject
ln -s  $(pwd)/Project.pkl $(pwd)/PklProject

export ENV_DIR=$(pwd)/.cedpm/env/

if [ -z $2 ];
	then
	printf "Resolving dependancies of $(pwd) ... "
	if ! $PKL project resolve; then
		rm -f $(pwd)/PklProject
		echo "Error resolving dependancies, please review your Project file."

		exit 1;
	fi;
	echo "done."
echo "- - - - - -"
rm -f $(pwd)/PklProject

if ! test -f "PklProject.deps.json"; then
	echo "Nothing to install.";
	exit 1;
fi;

mv "PklProject.deps.json" "Project.deps.json"
json_file="Project.deps.json"

jq -c '.resolvedDependencies[] | {type, uri, path}' "$json_file" | while read -r dep; do
  type=$(echo "$dep"  | tr '\n' ' ' | jq -r '.type')
  uri=$(echo "$dep"  | tr '\n' ' ' | jq -r '.uri' | sed 's/projectpackage:\/\//https:\/\//g' )
  if [ "$type" = "remote" ]; then
    pkg=$(curl -s -L "$uri")
	name=$(echo "$pkg" | tr '\n' ' ' | jq -r '.name')
  	target_dir="$(pwd)/.cedpm/packages/$name"
	printf "fetching $name... "
    mkdir -p "$target_dir"
#	echo "---"
	zip=$(echo "$pkg" | tr '\n' ' ' | jq -r '.packageZipUrl')
#	echo "waas trying \n\n\n$zip"
    wget -q -O "$target_dir/.remote.zip" "$zip"  # Adjust command if necessary
	unzip -oq  "$target_dir/.remote.zip" -d "$target_dir" 2>/dev/null
	echo " done."
    rm "$target_dir/.remote.zip"  # Cleanup the zip file
	mkdir -p $target_dir/.cedpm
	if test -f $target_dir/.cedpm/on_install.sh; then
		chmod +x $target_dir/.cedpm/on_install.sh
		$target_dir/.cedpm/on_install.sh
	fi;
	rm -rf $target_dir/.cedpm/1.0
	ln -s $target_dir/../../1.0 $target_dir/.cedpm/1.0
  elif [ "$type" = "local" ]; then
  	path=$(echo "$dep"  | tr '\n' ' ' | jq -r '.path')
#	echo "[[$path]]] [[$dep]]"
	name=$(pkl eval "$path/Project.pkl" --format json | tr '\n' ' ' | jq -r '.package.name')
	echo "Found local dep $name"
#  	target_dir=".cedpm/packages/$name"
    if [ -n "$path" ]; then
	  rm -rf $target_dir
#	  echo  "-rP" "$path" "$target_dir"
#      cp -rP "$path" "$target_dir"
	  mkdir -p $path/.cedpm
	  if test -f $path/.cedpm/on_install.sh; then
		chmod +x $path/.cedpm/on_install.sh
		$path/.cedpm/on_install.sh
	  fi;
#	  rm -rf $path/.cedpm/1.0
#	  ln -s $(dirname $path/../../1.0) $path/.cedpm/1.0
#	  echo "Running ln -s $(dirname $path/../../1.0) $path/.cedpm/1.0"
    else
      echo "No path specified for local dependency '$name'. Skipping."
    fi
  fi
done

	exit;
fi;

package_name=$2
echo "Installing package $2 ..."
CONFIG_FILE="Project.pkl"
cp "$CONFIG_FILE" "${CONFIG_FILE}.bak"
NEW_DEPENDENCY='  ["new_dependency"] {\n    uri = "package://pkg.pkl-lang.org/new_dependency@1.0.0"\n  }'
awk -v new_dep="$NEW_DEPENDENCY" '
BEGIN { brace_counter = 0; append_done = 0; }
/dependencies \{/ { inside_dependencies_block = 1; }
{
  if (inside_dependencies_block == 1) {
    if (/}/ && brace_counter == 1) {
      print new_dep; # Append new dependency before the closing brace of the dependencies block
      append_done = 1;
    }
    if (/}/) { brace_counter--; }
    if (/{/) { brace_counter++; }
  }
  print; # Print the current line
  if (/}/ && append_done == 1) { inside_dependencies_block = 0; }
}' "$CONFIG_FILE" > temp_file && mv temp_file "$CONFIG_FILE"

echo "New dependency added successfully."
