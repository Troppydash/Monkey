package evaluator

import (
	"Monkey/object"
	"Monkey/options"
	"Monkey/runner"
	"Monkey/tmp"
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
