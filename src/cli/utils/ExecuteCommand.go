package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func ExecuteCommand(commandStr string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", commandStr)
	} else {
		cmd = exec.Command("sh", "-c", commandStr)
	}

	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing command: %w, output: %s", err, cmdOutput)
	}

	fmt.Printf("%s", cmdOutput)
	return nil
}
