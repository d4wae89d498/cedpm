package builtins

import (
//	"fmt"
	"os"
//	"evaluator"
	"path/filepath"
	"os/exec"
	"cedpm.org/project_manifest"
	"cedpm.org/internal"
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

	pklPath := "pkl"
	pklExecEnv := os.Getenv("PKL_EXEC")
	if pklExecEnv != "" {
		pklPath = pklExecEnv
	}

	cmd := exec.Command(pklPath, "project", "resolve")
	cmd.Dir = projectDir
	cmdErr := cmd.Run();

	if err := os.Remove(pklProjectPath); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if err := os.Rename(symlinkPath, depsFile); err != nil {
		panic(err)
	}

	if cmdErr != nil {
		panic(cmdErr)
	}

	internal.Debug("Downloading dependencies ...")

	err := project_manifest.DownloadDependencies(projectDir);
	if err != nil {
		panic(err)
	}
}
