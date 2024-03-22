package main

import (
	"bytes"
	"errors"
	"net/url"
	"os/exec"
//	"fmt"
)

// cliCommandReader executes shell commands based on URLs with the cli:// scheme.
type cliCommandReader struct{

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

	cmdStr := url.String()[6:]
	// Extracting the command from the URL path
	//commandWithArgs := url.Path // Removing the leading slash
	//if commandWithArgs == "" {
	//	return nil, errors.New("no command provided")
	//}

	// Splitting command and arguments
	//parts := bytes.SplitN([]byte(commandWithArgs), []byte(" "), 2)
	cmd := exec.Command(cmdStr)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
