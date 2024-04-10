package project_manifest

import (
	"encoding/json"
	"fmt"
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


func ParseProjectAddons(jsonData string, commandList *CommandList) error {
	var addons projectFileAddons

	if err := json.Unmarshal([]byte(jsonData), &addons); err != nil {
		//("Error parsing JSON:", err)
		return err
	}

	fmt.Println("Project File addons:", addons)


	commandList.Paths = append(
		commandList.Paths,
		addons.Paths...
	)

	for key, value := range addons.Commands {
		_, exists := commandList.Commands[key]
		if exists {
			return errors.New("Err")
		}
		_ = value
		_ = key
    }


	return nil
}
