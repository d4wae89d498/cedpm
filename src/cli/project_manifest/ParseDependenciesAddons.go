package project_manifest

import (
_	"cedpm.org/internal"
_	"errors"
	"fmt"
	"path"
	"path/filepath"
	"os"
	"cedpm.org/evaluator"
)

func ParseDependenciesAddons(projectDir string, resolvedDependencies []Dependency, commandList *CommandList) error {
	var baseDir = path.Join(projectDir, ".cedpm", "packages")

	for _, dep := range resolvedDependencies {
		depPath := path.Join(baseDir, dep.Name)
		projectFile := filepath.Join(depPath, "Project")

		if _, err := os.Stat(projectFile); os.IsNotExist(err) {
			// Project file does not exist, skip this dependency
			continue
		}

		// TODO: what project path to give here ? none ?
		jsonData, err := evaluator.EvaluateFile(projectFile, depPath)
		if err != nil {
			return fmt.Errorf("failed to evaluate project file for dependency '%s': %v", dep.Name, err)
		}

		if err := ParseProjectAddons(jsonData, commandList); err != nil {
			return fmt.Errorf("failed to parse project addons for dependency '%s': %v", dep.Name, err)
		}
	}

	return nil
}
