package project_manifest

import (
	"cedpm.org/internal"
	"errors"
	"fmt"
)

func ParseDependenciesAddons(resolvedDependencies []Dependency, commandList *CommandList) error {
	for key, value := range resolvedDependencies {

		internal.Debug("[%d] %s", key, value)
	}
	fmt.Println("deps:", resolvedDependencies)

	panic(errors.New("Not implemented"))
	return nil
}
