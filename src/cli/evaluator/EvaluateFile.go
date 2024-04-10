package evaluator

import (
	"context"
//	"fmt"
	"os"
	"cedpm.org/pkl-go/pkl"
	"errors"
)

// cliCommandReader is assumed to be defined elsewhere

func EvaluateFile(filePath string, projectDir string) (string, error) {
	opts := func(opts *pkl.EvaluatorOptions) {
		pkl.PreconfiguredOptions(opts)
		reader := &cliCommandReader{}
		pkl.WithResourceReader(reader)(opts)
		opts.OutputFormat = "json"
	}

	manager := pkl.NewEvaluatorManager()
	var evaluator pkl.Evaluator
	var err error
	if projectDir != "" {
		projectEvaluator, err := manager.NewEvaluator(context.Background(), opts)
		if err != nil {
			return "", err
		}

		projectFileName := projectDir + "/Project"
		if _, err := os.Stat(projectFileName); err != nil {
			return "", errors.New("No 'Project' found")
		}

		pkl.RegisterMapping("pkl.AppleProject#RemoteDependency", pkl.ProjectRemoteDependency{})
		project, err := pkl.LoadProjectFromEvaluator(context.Background(), projectEvaluator, projectFileName)
		if err != nil {
			return "", err
		}

		withProject := func(project *pkl.Project) func(opts *pkl.EvaluatorOptions) {
			return func(opts *pkl.EvaluatorOptions) {
				pkl.WithProjectEvaluatorSettings(project)(opts)
				opts.ProjectDir = projectDir
				opts.DeclaredProjectDepenedencies = project.Dependencies()
			}
		}
		newOpts := []func(opts *pkl.EvaluatorOptions){withProject(project), opts}
		evaluator, err = manager.NewEvaluator(context.Background(), newOpts...)
		if err != nil {
			return "", err
		}
	} else {
		evaluator, err = pkl.NewEvaluator(context.Background(), opts)
		if err != nil {
			return "", err
		}
	}
	defer evaluator.Close()

	return evaluator.EvaluateOutputText(context.Background(), pkl.FileSource(filePath))
//	if err != nil {
//		return "", err
//	}

/*
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil*/
}


/*
func EvaluateFile(filePath string, projectDir string) {
	// Default PKL options
	opts := func(opts *pkl.EvaluatorOptions) {
		pkl.PreconfiguredOptions(opts)
		reader := &cliCommandReader{}
		pkl.WithResourceReader(reader)(opts)
		opts.OutputFormat = "json"
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

		var projectFileName string = projectDir + "/Project";
		_, err = os.Stat(projectFileName)
		if err != nil && os.IsNotExist(err) {
			projectFileName = projectDir + "/PklProject"
		}
		_, err = os.Stat(projectFileName)
		if err != nil && os.IsNotExist(err) {
			panic(errors.New("No 'Project' or 'PklProject' found."))
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
				opts.ProjectDir = projectDir //+ "/.cedpm" //strings.TrimPrefix(strings.TrimSuffix(project.ProjectFileUri, "/PklProject"), "file://")
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
	var cfg string;
	cfg, err = evaluator.EvaluateOutputText(context.Background(), pkl.FileSource(filePath));

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cfg)
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

*/
