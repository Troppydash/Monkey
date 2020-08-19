package runner

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/parser"
	"Monkey/tmp"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// Singleton
var runner = &Runner{}

type Runner struct {
	files []string
}

func GetInstance() *Runner {
	return runner
}

func (r *Runner) CompileAbs(filename string) (*ast.Program, error) {
	content, err := r.ReadFile(path.Join(tmp.CurrentProcessingFileDirectory, filename))
	if err != nil {
		return nil, err
	}
	return r.ParseProgram(string(content), filename), nil
}

func (r *Runner) Pop(filename string) {
	var index int
	for i, file := range r.files {
		if filename == file {
			index = i
		}
	}

	for i := len(r.files); i > len(r.files)-index; i++ {
		r.files = r.files[:1]
	}
}

func (r *Runner) Compile(location string) (*ast.Program, error) {
	var filename string
	// if is a std include
	if !strings.HasSuffix(location, ".mky") {
		filename = path.Join(tmp.STDDirectory, location+".mky")
		tmp.SetAbsoluteDirectory(path.Dir(filename))
	} else {
		re := regexp.MustCompile("[/\\\\]")
		folders := re.Split(location, -1)
		dirs := append([]string{tmp.CurrentProcessingFileDirectory}, folders...)
		filename = path.Join(dirs...)
		tmp.SetAbsoluteDirectory(path.Dir(filename))
	}

	return r.CompileAbs(path.Base(filename))
}

func (r *Runner) ParseProgram(content string, filename string) *ast.Program {
	l := lexer.New(content, filename)
	p := parser.New(l)
	return p.ParseProgram()
}

func (r *Runner) ReadFile(filename string) ([]byte, error) {
	for _, file := range r.files {
		if file == filename {
			fmt.Println("Circular Dependency Error")
			for _, file := range r.files {
				fmt.Printf("  %s\n", file)
			}
			fmt.Printf("  %s\n", filename)
			os.Exit(1)
		}
	}
	r.files = append(r.files, filename)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Cannot read file %q\n", filename)
		return []byte{}, err
	}

	return content, nil
}
