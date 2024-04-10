package main

import (
	"fmt"
	"os"
	"flag"
	"cedpm.org/project_manifest"
	"cedpm.org/project_skeleton"
	"cedpm.org/evaluator"
	"cedpm.org/internal"
	"cedpm.org/builtins"
	"path"
)

const version = "1.0.0"
const debugEnabled = true

func showUsage() {
	fmt.Printf(`Usage: %s [command] [options]
Commands:
  install    Install the package
  search     Search for the package
  uninstall  Uninstall the package

  help       Show help message
  version    Show version information
`, os.Args[0])
}

func main() {

	internal.DebugEnabled = debugEnabled

	// Init

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	internal.Debug("cedpm in %s", currentDir)

	// Basic arguments

	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Printf("Error: No command provided. Type %s help for usage.", os.Args[0])
		println()
		return
	}

	if os.Args[1] == "help" {
		showUsage()
		println()
		return
	}

	if os.Args[1] == "version" {
		fmt.Println("cedpm version: %s", version)
		return
	}

	if os.Args[1] == "init" {
		project_manifest.Init()
		return
	}

	// Require project file

	projectDir, err := project_manifest.FindProjectDir(currentDir)
	if err != nil {
		fmt.Println("Unable to use project file:", err)
		return
	}

	// Init .cedpm directory
	project_skeleton.Init(projectDir);

	projectFile := path.Join(projectDir, project_manifest.ProjectFileName)
	internal.Debug("cedpm project file : %s", projectFile)
	jsonData, err := evaluator.EvaluateFile(projectFile, projectDir)
	if err != nil {
		fmt.Printf("Error reading Project Manifest\n", err)
	}
	_ = jsonData

	//ParseProjectManifest(jsonData)
	//ParseProjectDeps(projectDir)

	// Complex arguments that use the project file

	if os.Args[1] == "install" {
		if len(os.Args) == 2 {
			builtins.Install(projectDir)

			internal.Debug("Project dependancies installed successfully.")
			return
		}
		// TODO : add support for package name

		fmt.Printf("Usage: %s install\n", os.Args[0])
		return
	}
	return

	if os.Args[1] == "eval" {
		if len(os.Args) != 3 {
			fmt.Printf("Usage: %s eval <filename>", os.Args[0])
			return
		}

		builtins.EvaluateFile(os.Args[2], projectDir);
		return
	}

	// lire le .dependencies json
	// puis lire tout les Project Manifest,
	// fournir les commandes

	command := os.Args[1]
	switch command {
	case "install", "search", "uninstall", "eval":	// its a builtin
		fmt.Printf("Executing '%s' command...\n", command)


		os.Exit(0)
	default:
	}

	// try in all path
/*



	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("PklProject file found at: %s\n", pklProjectPath)
	*/
}
