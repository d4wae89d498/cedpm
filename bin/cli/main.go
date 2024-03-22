package main

import (
	"context"
	"fmt"
	"os"
)


type Instruction struct {
	Inputs 	[]string 			`pkl:"inputs"`
	Output 	string 	 			`pkl:"output"`
	Command string 				`pkl:"command"`
}

type Target struct {
	Name string					`pkl:"name"`
	Type string					`pkl:"type"`
	Instructions []Instruction	`pkl:"instructions"`
}

type MyConfig struct {
	Targets []Target 			`pkl:"targets"`
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
