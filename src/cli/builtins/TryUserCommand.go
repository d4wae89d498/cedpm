package builtins

import (
	"errors"
_	"cedpm.org/internal"
	"cedpm.org/project_manifest"
)

func TryUserCommand(commandList project_manifest.CommandList, input []string) error {
	_ = commandList
	_ = input
	return errors.New("Not implemented")
}
