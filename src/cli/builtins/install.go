package builtins

import (
//	"fmt"
	"os"
//	"evaluator"
	"path/filepath"
	"os/exec"
)

func Install(projectDir string) {

	projectFile := filepath.Join(projectDir, "Project")
	pklProjectPath := filepath.Join(projectDir, "PklProject")
	depsFile := filepath.Join(projectDir, ".dependencies.json")
	symlinkPath := filepath.Join(projectDir, "PklProject.deps.json")

	// Attempt to remove the generated files and ignore "not found" errors.
	filesToRemove := []string{pklProjectPath, depsFile, symlinkPath}
	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			panic(err)
		}
	}

	// Create a symlink from Project to PklProject
	if err := os.Symlink(projectFile, pklProjectPath); err != nil {
		panic(err)
	}

	// Ensure the symlink and the .deps.json file are cleaned up afterwards.
	defer func() {
		if err := os.Remove(pklProjectPath); err != nil && !os.IsNotExist(err) {
			panic(err)
		}
		if err := os.Rename(symlinkPath, depsFile); err != nil {
			panic(err)
		}
	}()

	pklPath := "pkl"
	pklExecEnv := os.Getenv("PKL_EXEC")
	if pklExecEnv != "" {
		pklPath = pklExecEnv
	}

	cmd := exec.Command(pklPath, "project", "resolve")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
