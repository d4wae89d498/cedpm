package project_skeleton

import (
	"path/filepath"
	"os"
	"cedpm.org/internal"
	"cedpm.org/utils"
)

func Init(projectDir string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	sourceDir := filepath.Join(exeDir, "..", "cedpm_skeleton")
	destDir := filepath.Join(projectDir, ".cedpm")

	err = utils.CopyDir(sourceDir, destDir)
	if err != nil {
		return err
	}

	internal.Debug(".cedpm directory copied successfully.")
	return nil
}
