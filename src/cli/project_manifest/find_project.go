package project_manifest

import (
	"fmt"
	"os"
	"path/filepath"
)

const ProjectFileName = "PklProject"

func FindProjectDir(startDir string) (string, error) {
	pklProjectDir := filepath.Join(startDir, ProjectFileName)
	if _, err := os.Stat(pklProjectDir); err == nil {
		return startDir, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if startDir == "/" {
		return "", fmt.Errorf("project file not found")
	}

	parentDir := filepath.Dir(startDir)
	return FindProjectDir(parentDir)
}
