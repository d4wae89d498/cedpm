package builtins

import (
	"cedpm.org/evaluator"
	"cedpm.org/internal"
	"fmt"
	"path/filepath"
	"os"
)

func EvaluateFile(filePath string, projectDir string) {
	internal.Debug("evaluating '%s'\n", filePath)

	// Define the paths
	depsFile := filepath.Join(projectDir, ".dependencies.json")
	symlinkPath := filepath.Join(projectDir, "PklProject.deps.json")

	// Remove symlink if already exists
	if err := os.Remove(symlinkPath); err != nil {
		if os.IsNotExist(err) {
			internal.Debug("File %s does not exists, continuing...", symlinkPath)
		} else {
			panic(err)
		}
	}

	// Create a symbolic link to .dependencies.json named PklProject.deps.json
	if err := os.Symlink(depsFile, symlinkPath); err != nil {
		panic(err)
	}

	// Defer the removal of the symbolic link until the end of the function's execution
	defer func() {
		if err := os.Remove(symlinkPath); err != nil {
			if os.IsNotExist(err) {
				internal.Debug("File %s does not exists, continuing...", symlinkPath)
			} else {
				panic(err)
			}
		}
	}()

	json, err := evaluator.EvaluateFile(filePath, projectDir)
	if err != nil {
		panic(err)
	}
	fmt.Println(json)
}
