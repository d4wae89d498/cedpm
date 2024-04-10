package project_manifest

import (
	"encoding/json"
_	"fmt"
//	"io/ioutil"
//	"os"
_	"cedpm.org/internal"
	"errors"
)

type Command struct {
	Before  	[]string
	After   	[]string
	On      	[]string
	Actions 	[]string
}

type CommandList struct {
	Paths     	[]string
	Commands 	map[string]Command
}

/*\
 *		./Project 	F i l e
\*/

type projectFileAddons struct {
	Paths		[]string			`json:"paths"`
	Before  	map[string][]string	`json:"before"`
	After   	map[string][]string	`json:"after"`
	On      	map[string][]string	`json:"on"`
	Commands 	map[string][]string	`json:"commands"`
}

//////////////////////////////////////////////////////////////

// TODO : fix that to incrementally add even if commnad was not set before
func ParseProjectAddons(jsonData string, commandList *CommandList) error {
	var addons projectFileAddons

	if err := json.Unmarshal([]byte(jsonData), &addons); err != nil {
		return err
	}

	// Append paths
	commandList.Paths = append(commandList.Paths, addons.Paths...)

	// Process commands
	for cmdKey, actions := range addons.Commands {
		if _, exists := commandList.Commands[cmdKey]; exists {
			return errors.New("duplicate command key found")
		}

		// Construct new Command struct with the addon data
		newCommand := Command{
			Before:  addons.Before[cmdKey],
			After:   addons.After[cmdKey],
			On:      addons.On[cmdKey],
			Actions: actions,
		}

		// Add the new Command to the CommandList's Commands map
		if commandList.Commands == nil {
			commandList.Commands = make(map[string]Command)
		}
		commandList.Commands[cmdKey] = newCommand
	}


	return nil
}
