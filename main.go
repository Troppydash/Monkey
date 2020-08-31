package main

import (
	"Monkey/evaluator"
	"Monkey/object"
	"Monkey/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {

	// Run File
	if len(os.Args) == 2 {

		// Get filename
		filename := os.Args[1]

		// Create env
		env := object.NewEnvironment()

		// Link std
		err := evaluator.LinkAndEval("std", env)
		if err != nil {
			fmt.Printf("Failed to compile file %q\n", filename)
			return
		}

		// Compile
		err = evaluator.LinkAndEval(filename, env)

		if err != nil {
			fmt.Printf("Failed to compile file %q\n", filename)
			return
		}
		return
	}

	// Retrieve os user
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Display Welcoming message
	fmt.Printf("Hello %s! Welcome to the Monkey Programming Language!\n",
		usr.Username)
	fmt.Printf("REPL Started!\n")

	// Start the repl
	repl.Start(os.Stdin, os.Stdout)
}

//func LinkFile(libraryName string, env *object.Environment) {
//	p, e := runner.GetInstance().Compile(libraryName)
//	if e != nil {
//		panic("Failed to link library " + libraryName)
//	}
//	evaluator.Eval(p, env)
//}
