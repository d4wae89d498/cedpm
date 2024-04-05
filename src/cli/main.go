package main

import (
	"fmt"
	"os"
	"flag"
	"cedpm.org/project_manifest"
	"cedpm.org/evaluator"
	"cedpm.org/verbrose"
	"path"
)

const version = "1.0.0"
const debug = true

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

	verbrose.Enabled = debug

	// Init

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	verbrose.Println("cedpm in ", currentDir)

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
	verbrose.Println("cedpm project file :", projectDir)
	evaluator.EvaluateFile(path.Join(projectDir, project_manifest.ProjectFileName), projectDir)

	// Complex arguments that use the project file

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
