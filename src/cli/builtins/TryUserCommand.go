package builtins

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"cedpm.org/project_manifest"
	"cedpm.org/utils"
)

var (
	ErrNotFound       = errors.New("command or path not found")
	ErrUnknownCommand = errors.New("unknown command error")
)



// Todo: return command and arguments
func TryUsingCommands(commandList project_manifest.CommandList, input []string, startIndex int) error {
	if startIndex >= len(input) {
		return ErrNotFound
	}

	key := strings.Join(input[:len(input)-startIndex], " ")
	fmt.Printf("Trying, [%s]\n", key)
	command, exists := commandList.Commands[key]
	if exists {
		fmt.Printf("EXISTS\n")
		for _, action := range command.Actions {
			fmt.Printf("Action[]\n")
			err := utils.ExecuteCommand(strings.Replace(action, "{{@}}", strings.Join(input[len(input)-startIndex:], " "), -1))
			if err != nil {
				return err
			}
		}
		return nil
	}

	return TryUsingCommands(commandList, input, startIndex+1)
}

// Todo: return command and arguments
func TryUsingPaths(commandList project_manifest.CommandList, input []string) error {
	if len(input) == 0 {
		return ErrNotFound
	}

	path := strings.Join(input, "/")
	for _, p := range commandList.Paths {
		if _, err := os.Stat(p + "/" + path); !os.IsNotExist(err) {
			cmd := exec.Command(p+"/"+path) // TODO : add arguments
			return cmd.Run()
		}
	}

	return TryUsingPaths(commandList, input[:len(input)-1])
}

// Todo: use returned command and arguments to execute Before, After, and On
// Todo: determine if we keep events for a childs, exemple: we have 'test' as command key, and we run 'test foo'
//		would 'test' s events apply ?
func TryUserCommand(commandList project_manifest.CommandList, input []string) error {
	fmt.Println("trying with commands...")
	if err := TryUsingCommands(commandList, input, 0); err != nil {
		if err != ErrNotFound {
			return err
		}
	} else {
		return nil
	}
	fmt.Println("trying with paths...")
	return TryUsingPaths(commandList, input);
}
