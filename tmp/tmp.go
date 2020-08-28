package tmp

import (
	"fmt"
	"os"
	"path"
)

// the dir of te file that it is processing
var CurrentProcessingFileDirectory string

// Command Directory
var CurrentDirectory string

// root path
var ExeDirectory string

// lib directory
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
		fmt.Fprintln(os.Stderr, "Cannot find executable location, Proceeding with the current directory")
		//panic("Cannot find executable location")
		EXE = CD
	}
	ExeDirectory = EXE

	STDDirectory = path.Join(ExeDirectory, "lib")

	CurrentProcessingFileDirectory = CD
}
