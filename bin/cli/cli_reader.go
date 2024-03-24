package main

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os/exec"
	"fmt"
	"runtime"
)

// executeCommand executes a shell command, automatically adjusting for the
// operating system. On Unix-like systems, it uses `bash -c`, and on Windows,
// it uses `cmd /C`.
func executeCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: Use cmd /C
		cmd = exec.Command("cmd", "/C", command)
	} else {
		// Unix-like: Use bash -c
		cmd = exec.Command("bash", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}



///--------

// cliCommandReader executes shell commands based on URLs with the cli:// scheme.
type cliCommandReader struct {
}

// Scheme returns the scheme this reader handles.
func (c *cliCommandReader) Scheme() string {
	return "cli"
}

// IsGlobbable returns false as globbing is not applicable to shell command execution.
func (c *cliCommandReader) IsGlobbable() bool {
	return false
}

// HasHierarchicalUris returns false as CLI command execution doesn't involve hierarchical URIs.
func (c *cliCommandReader) HasHierarchicalUris() bool {
	return false
}

// ListElements is not applicable for cliCommandReader as it does not involve browsing a file system structure.
func (c *cliCommandReader) ListElements(url url.URL) ([]PathElement, error) {
	return nil, errors.New("ListElements is not supported by cliCommandReader")
}

// Read executes the shell command specified in the URL and returns its output.
func (c *cliCommandReader) Read(url url.URL) ([]byte, error) {

	fmt.Println("Running cmd...");
	cmdStr := url.String()[6:]
	decodedBytes, err := base64.StdEncoding.DecodeString(cmdStr)
	if err != nil {
		return nil, err
	}
	output, err := executeCommand(string(decodedBytes))
    if err != nil {
        return nil, err
    }
	return []byte(output), nil
}
