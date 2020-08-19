package tmp

import (
	"os"
	"path"
)

// TODO Fix
var Filename string

var CurrentProcessingFileDirectory string

// Command Directory
var CurrentDirectory string
var ExeDirectory string

var STDDirectory string

func SetAbsoluteDirectory(absolute string) {
	CurrentProcessingFileDirectory = absolute
}

func init() {
	CD, err := os.Getwd()
	if err != nil {
		panic("Cannot get current directory")
	}
	CurrentDirectory = CD

	EXE := os.Getenv("MKYROOT")
	if len(EXE) == 0 {
		panic("Cannot find executable location")
	}
	ExeDirectory = EXE

	STDDirectory = path.Join(ExeDirectory, "lib", "std")

	CurrentProcessingFileDirectory = CD
}
