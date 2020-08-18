package main

import (
	"Monkey/repl"
	"Monkey/runner"
	"fmt"
	"os"
	"os/user"
)

func main() {

	// Run File
	if len(os.Args) == 2 {

		filename := os.Args[1]
		// TODO: Make everything runner
		r := runner.New()
		r.Execute(filename)

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
