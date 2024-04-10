package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
//	"path/filepath"
	"time"
)

type Target struct {
	Dependencies []string `json:"Dependencies"`
	Name         string   `json:"Name"`
	Command      string   `json:"Command"`
}

type Artifact struct {
	Targets []Target `json:"Targets"`
}

type MyConfig struct {
	Artifacts []Artifact `json:"Artifacts"`
}

func main() {
	targetName := flag.String("target", "", "Name of the target to build")
	changeDir := flag.String("C", ".", "Change to directory before doing anything")
	flag.Parse()

	if err := os.Chdir(*changeDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change directory: %s\n", err)
		os.Exit(1)
	}

	// Read JSON from stdin
	configData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %s\n", err)
		os.Exit(1)
	}

	var config MyConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %s\n", err)
		os.Exit(1)
	}

	for _, artifact := range config.Artifacts {
		for _, target := range artifact.Targets {
			if *targetName == "" || *targetName == target.Name {
				if targetUpToDate(target) {
					fmt.Printf("\033[1;33mTarget '%s' is up to date.\033[0m\n", target.Name)
					continue
				}
				executeTarget(target)
			}
		}
	}
}

func targetUpToDate(target Target) bool {
	targetModTime, err := modTime(target.Name)
	if err != nil {
		return false
	}

	for _, dep := range target.Dependencies {
		depModTime, err := modTime(dep)
		if err != nil || depModTime.After(targetModTime) {
			return false
		}
	}
	return true
}

func modTime(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

func executeTarget(target Target) {
	fmt.Printf("\033[1;34mBuilding target: %s\033[0m\n", target.Name)
	cmd := exec.Command("sh", "-c", target.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("\033[1;31mFailed to build %s: %s\033[0m\n", target.Name, err)
	} else {
		fmt.Printf("\033[1;32mSuccessfully built %s\033[0m\n", target.Name)
	}
}
