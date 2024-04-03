package main

import (
	"context"
	"fmt"
	"os"
	"github.com/apple/pkl-go/pkl"
	"errors"
	"encoding/json"
	//"strings"
)

type Target struct {
	Dependencies 	[]string 			`pkl:"dependencies"`
	Name 			string 	 			`pkl:"name"`
	Command 		string 				`pkl:"command"`
}

type Artifact struct {
	Targets			[]Target				`pkl:"targets"`
}

type MyConfig struct {
	Artifacts 		[]Artifact 			`pkl:"artifacts"`
}

func main() {

	// Getting cli args
	filePath, projectDir := "", ""
	if len(os.Args) == 2 {
		filePath = os.Args[1]
	} else if len(os.Args) == 3 {
		filePath, projectDir = os.Args[1], os.Args[2]
	} else {
		fmt.Printf("Usage: %s <FilePath> [ProjectDir]\n", os.Args[0])
		os.Exit(1)
	}

	// Default PKL options
	opts := func(opts *pkl.EvaluatorOptions) {
		pkl.PreconfiguredOptions(opts)
		reader := &cliCommandReader{}
		pkl.WithResourceReader(reader)(opts)
	}

	var evaluator pkl.Evaluator
	var err error
	if projectDir != "" {

		// Setting up the project configuration
		//fmt.Println("using projectDir: %s", projectDir)
		manager := pkl.NewEvaluatorManager()
		projectEvaluator, err := manager.NewEvaluator(context.Background(), opts)
		if err != nil {
			panic(err)
		}

		var projectFileName string = projectDir + "/Project.pkl";
		_, err = os.Stat(projectFileName)
		if err != nil && os.IsNotExist(err) {
			projectFileName = projectDir + "/PklProject"
		}
		_, err = os.Stat(projectFileName)
		if err != nil && os.IsNotExist(err) {
			panic(errors.New("No 'Project.pkl' or 'PklProject' found."))
		}
		//fmt.Println("Using projectFile: %s", projectFileName)

		pkl.RegisterMapping("pkl.AppleProject#RemoteDependency", pkl.ProjectRemoteDependency{})
		project, err := pkl.LoadProjectFromEvaluator(context.Background(), projectEvaluator, projectFileName)
		if err != nil {
			panic(err)
		}
		var withProject = func(project *pkl.Project) func(opts *pkl.EvaluatorOptions) {
			return func(opts *pkl.EvaluatorOptions) {
				pkl.WithProjectEvaluatorSettings(project)(opts)
				opts.ProjectDir = projectDir //strings.TrimPrefix(strings.TrimSuffix(project.ProjectFileUri, "/PklProject"), "file://")
				opts.DeclaredProjectDepenedencies = project.Dependencies()

			}
		}
		newOpts := []func(opts *pkl.EvaluatorOptions){
			withProject(project),
		}
		newOpts = append(newOpts, opts)
		evaluator, err = manager.NewEvaluator(context.Background(), newOpts...)
		if err != nil {
			panic(err)
		}

	} else {
		evaluator, err = pkl.NewEvaluator(context.Background(), opts)
	}
	if err != nil {
		panic(err)
	}
	defer evaluator.Close()

	// Evalute provided PKL file
	var cfg MyConfig
	if err := evaluator.EvaluateModule(context.Background(), pkl.FileSource(filePath), &cfg); err != nil {
		//fmt.Println("error in eval")
		panic(err)
	}
	//fmt.Printf("%+v\n", cfg)

	// Conversion de l'instance de Person en JSON
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	// Conversion des données JSON en chaîne (pour l'affichage)
	jsonString := string(jsonData)

	fmt.Println(jsonString)
}
