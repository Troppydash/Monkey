package runner

import (
	"Monkey/evaluator"
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"Monkey/tmp"
	"fmt"
	"io/ioutil"
)

type Runner struct {
	content []byte
}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Execute(filename string) {
	// TODO FIx
	tmp.Filename = filename
	r.Include(filename)

	l := lexer.New(string(r.content), filename)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	evaluator.Eval(program, env)
}

func (r *Runner) Include(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("An Error occurred reading the file")
	}

	r.content = append(r.content, content...)
}
