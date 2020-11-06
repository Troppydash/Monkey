package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
	"Monkey/options"
	"Monkey/runner"
	"Monkey/tmp"
	"Monkey/token"
	"errors"
)

var std = []string{
	"std",
}

func LinkSTD(env *object.Environment) error {
	for _, stdLocation := range std {
		err := LinkAndEval(stdLocation, env)
		if err != nil {
			return err
		}

		if options.FatalErrors {
			return errors.New("FatalError Encountered")
		}
	}

	return nil
}

func LinkAndEvalModule(filename string, module *object.Module, token token.Token) error {
	old := tmp.CurrentProcessingFileDirectory
	abs := runner.GetInstance().ToAbsolute(filename)
	program, e := runner.GetInstance().CompileAbs(abs)
	if e != nil {
		return e
	}
	module.Body = &ast.BlockStatement{
		Token:      token,
		Statements: program.Statements,
	}
	DefineMacros(program, module.Env)
	expanded := ExpandMacros(program, module.Env)

	Eval(expanded, module.Env)
	if options.FatalErrors {
		return errors.New("FatalError Encountered")
	}
	runner.GetInstance().Pop(abs)
	tmp.CurrentProcessingFileDirectory = old
	return nil
}

func LinkAndEval(filename string, env *object.Environment) error {
	old := tmp.CurrentProcessingFileDirectory
	abs := runner.GetInstance().ToAbsolute(filename)
	program, e := runner.GetInstance().CompileAbs(abs)
	if e != nil {
		return e
	}
	DefineMacros(program, env)
	expanded := ExpandMacros(program, env)

	Eval(expanded, env)
	if options.FatalErrors {
		return errors.New("FatalError Encountered")
	}
	runner.GetInstance().Pop(abs)
	tmp.CurrentProcessingFileDirectory = old
	return nil
}
