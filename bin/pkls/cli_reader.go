package main

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os/exec"
	//"fmt"
	"runtime"
	"github.com/apple/pkl-go/pkl"
)

// executeCommand executes a shell command, automatically adjusting for the
// operating system. On Unix-like systems, it uses `bash -c`, and on Windows,
// it uses `cmd /C`.
func executeCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
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
func (c *cliCommandReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	return nil, errors.New("ListElements is not supported by cliCommandReader")
}

// Read executes the shell command encoded in base64 specified in the URL and returns its output.
func (c *cliCommandReader) Read(url url.URL) ([]byte, error) {

	//fmt.Println("Running cmd...");
	cmdStr := url.String()[6:]
	decodedBytes, err := base64.StdEncoding.DecodeString(cmdStr)
	if err != nil {
		panic(err)
		return nil, err
	}
	output, err := executeCommand(string(decodedBytes))
    if err != nil {
		panic(err)
        return nil, err
    }
	return []byte(output), nil
}
