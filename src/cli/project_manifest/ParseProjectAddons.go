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

// TODO : fix that to incrementally add after/before/on even if commnad was not set before, set command actions when 'commands' is present, do error when 'commands' is set and actions is already defined
/*func ParseProjectAddons(jsonData string, commandList *CommandList) error {
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
*/

func ParseProjectAddons(jsonData string, commandList *CommandList) error {
	var addons projectFileAddons

	if err := json.Unmarshal([]byte(jsonData), &addons); err != nil {
		return err
	}

	// Append paths
	commandList.Paths = append(commandList.Paths, addons.Paths...)

	if commandList.Commands == nil {
		commandList.Commands = make(map[string]Command)
	}

	// Process before, after, and on for initialization or appending to existing commands
	for _, actionType := range []struct {
		Data map[string][]string
		Type string
	}{
		{addons.Before, "before"},
		{addons.After, "after"},
		{addons.On, "on"},
	} {
		for cmdKey, values := range actionType.Data {
			cmd, exists := commandList.Commands[cmdKey]
			if !exists {
				cmd = Command{}
			}

			switch actionType.Type {
			case "before":
				cmd.Before = append(cmd.Before, values...)
			case "after":
				cmd.After = append(cmd.After, values...)
			case "on":
				cmd.On = append(cmd.On, values...)
			}

			commandList.Commands[cmdKey] = cmd
		}
	}

	// Process commands, ensuring no action conflicts
	for cmdKey, actions := range addons.Commands {
		cmd, exists := commandList.Commands[cmdKey]
		if exists && len(cmd.Actions) > 0 {
			return errors.New("actions already defined for command " + cmdKey)
		}
		cmd.Actions = actions
		commandList.Commands[cmdKey] = cmd
	}

	return nil
}
