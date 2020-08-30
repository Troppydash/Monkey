package runner

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/options"
	"Monkey/parser"
	"Monkey/tmp"
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	// TODO cache this

	content, err := r.ReadFile(filename)
	if err != nil {
		// graceful return
		if ce, ok := err.(CError); ok {
			if options.Debug {
				fmt.Println("Circular Dependency Warning")
				for _, file := range ce.files {
					fmt.Printf("  %s\n", file)
				}
			}
			return &ast.Program{
				Statements: []ast.Statement{},
			}, nil
		}

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

	end := len(r.files) - index - 1
	for i := len(r.files) - 1; i > end; i-- {
		r.files = remove(r.files, len(r.files)-1)
	}
	fmt.Printf("")
}
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func (r *Runner) ToAbsolute(location string) string {
	var filename string
	// if is a std include
	if !strings.HasSuffix(location, ".mky") {
		filename = filepath.Join(tmp.STDDirectory, filepath.Base(location), location+".mky")
		tmp.SetAbsoluteDirectory(filepath.Dir(filename))
	} else {
		re := regexp.MustCompile("[/\\\\]")
		folders := re.Split(location, -1)
		dirs := append([]string{tmp.CurrentProcessingFileDirectory}, folders...)
		filename = filepath.Join(dirs...)
		tmp.SetAbsoluteDirectory(filepath.Dir(filename))
	}

	return filename
}

func (r *Runner) ParseProgram(content string, filename string) *ast.Program {
	l := lexer.New(content, filename)
	p := parser.New(l)
	return p.ParseProgram()
}

type CError struct {
	files []string
}

func NewCError(files []string) CError {
	return CError{
		files,
	}
}

func (c CError) Error() string {
	return "Circular Dependency Error"
}

func (r *Runner) ReadFile(filename string) ([]byte, error) {
	for _, file := range r.files {
		if file == filename {
			return nil, NewCError(append(r.files, filename))
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
