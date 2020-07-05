package main

import (
	"Monkey/evaluator"
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"Monkey/repl"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

func main() {

	// Run File
	if len(os.Args) == 2 {

		filename := os.Args[1]
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("An Error occurred reading the file")
		}

		l := lexer.New(string(content), filename)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()

		evaluator.InitEvaluator()
		evaluator.Eval(program, env)

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
