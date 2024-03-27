package main

import (
	"context"
	"fmt"
	"os"
)


type Target struct {
	Dependencies 	[]string 			`pkl:"dependencies"`
	Name 			string 	 			`pkl:"name"`
	Command 		string 				`pkl:"command"`
}

type Artifact struct {
//	Name string					`pkl:"name"`
//	Type string					`pkl:"type"`
Targets			[]Target				`pkl:"targets"`
}

type MyConfig struct {
	Artifacts 		[]Artifact 			`pkl:"artifacts"`
}


func main() {

	err := os.Chdir("/Users/mfaussur/Desktop/cedpm/examples")
	if err != nil {
		// Handle the error
		fmt.Println("Error changing directory:", err)
		return
	}


	opts := func(opts *EvaluatorOptions){
		PreconfiguredOptions(opts)
		reader := &cliCommandReader{}
		WithResourceReader(reader)(opts)
	}


	evaluator, err := NewProjectEvaluator(context.Background(), "/Users/mfaussur/Desktop/cedpm/examples/", opts)
	if err != nil {
		panic(err)
	}
	defer evaluator.Close()


	var cfg MyConfig
	if err = evaluator.EvaluateModule(context.Background(), FileSource("/Users/mfaussur/Desktop/cedpm/examples/MathOperations.pkl"), &cfg); err != nil {
		fmt.Printf("error in eval\n")
		panic(err)
	}
	fmt.Printf("Got module: %+v", cfg)

}
