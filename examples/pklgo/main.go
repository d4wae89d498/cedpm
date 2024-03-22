package main
import (
	"context"
	"fmt"

	"github.com/apple/pkl-go/pkl"
)


// Project is the go representation of pkl.Project.
type Project struct {
	ProjectFileUri   string                    `pkl:"projectFileUri"`
	Package          *ProjectPackage           `pkl:"package"`
	EvaluatorSetings *ProjectEvaluatorSettings `pkl:"evaluatorSettings"`
	Tests            []string                  `pkl:"tests"`

	// internal field; use Project.Dependencies instead.
	// values are either *Project or *ProjectRemoteDependency
	RawDependencies map[string]any `pkl:"dependencies"`

	dependencies *pkl.ProjectDependencies `pkl:"-"`
}

// ProjectPackage is the go representation of pkl.Project#Package.
type ProjectPackage struct {
	Name                string   `pkl:"name"`
	BaseUri             string   `pkl:"baseUri"`
	Version             string   `pkl:"version"`
	PackageZipUrl       string   `pkl:"packageZipUrl"`
	Description         string   `pkl:"description"`
	Authors             []string `pkl:"authors"`
	Website             string   `pkl:"website"`
	Documentation       string   `pkl:"documentation"`
	SourceCode          string   `pkl:"sourceCode"`
	SourceCodeUrlScheme string   `pkl:"sourceCodeUrlScheme"`
	License             string   `pkl:"license"`
	LicenseText         string   `pkl:"licenseText"`
	IssueTracker        string   `pkl:"issueTracker"`
	ApiTests            []string `pkl:"apiTests"`
	Exclude             []string `pkl:"exclude"`
	Uri                 string   `pkl:"uri"`
}

// ProjectEvaluatorSettings is the Go representation of pkl.Project#EvaluatorSettings
type ProjectEvaluatorSettings struct {
	ExternalProperties map[string]string `pkl:"externalProperties"`
	Env                map[string]string `pkl:"env"`
	AllowedModules     []string          `pkl:"allowedModules"`
	AllowedResources   []string          `pkl:"allowedResources"`
	NoCache            *bool             `pkl:"noCache"`
	ModulePath         []string          `pkl:"modulePath"`
	Timeout            pkl.Duration          `pkl:"timeout"`
	ModuleCacheDir     string            `pkl:"moduleCacheDir"`
	RootDir            string            `pkl:"rootDir"`
}

func init() {
	pkl.RegisterMapping("pkl.AppleProject", Project{})
	pkl.RegisterMapping("pkl.AppleProject#RemoteDependency", pkl.ProjectRemoteDependency{})
}



func main() {

	evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
	if err != nil {
		panic(err)
	}
	defer evaluator.Close()

	var cfg Project
	if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource("./../PklProject"), &cfg); err != nil {
		fmt.Printf("error in eval\n")
		panic(err)
	}
	fmt.Printf("Got module: %+v", cfg)
}
