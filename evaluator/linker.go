package evaluator

import (
	"Monkey/object"
	"Monkey/runner"
	"Monkey/tmp"
)

func LinkAndEval(filename string, env *object.Environment) error {
	old := tmp.CurrentProcessingFileDirectory
	abs := runner.GetInstance().ToAbsolute(filename)
	p, e := runner.GetInstance().CompileAbs(abs)
	if e != nil {
		//fmt.Printf("Failed to compile file %q\n", filename)
		return e
	}
	Eval(p, env)
	runner.GetInstance().Pop(abs)
	tmp.CurrentProcessingFileDirectory = old
	return nil
}
