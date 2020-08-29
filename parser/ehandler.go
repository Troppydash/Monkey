package parser

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
)

const LinesAround = 3

func PrintParserError(err *ParseError) {
	fmt.Printf("Parser Error: %s, at %d:%d, in file %s\n",
		err.Message, err.RowNumber, err.ColumnNumber, err.Filename)
	rowNumber := int(err.RowNumber)

	fmt.Printf("[%s]\n", filepath.Base(err.Filename))
	rows, e := readFileRows(rowNumber, err.Filename)
	if e != nil {
		return
	}
	for index, row := range rows {
		// i have no idea how it came to this
		number := index + int(math.Max(0, float64(rowNumber-LinesAround+1))) + 1
		fmt.Printf("| %-3d %s", number, row)
		if number == rowNumber {
			fmt.Print("    <-- over here")
		}
		fmt.Println()
	}
	fmt.Println()
}

func readFileRows(rows int, filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Cannot read file %q\n", filename)
		return []string{}, err
	}

	cont := strings.ReplaceAll(string(content), "\r", "")
	lines := strings.Split(cont, "\n")

	return getRowsAround(lines, rows), nil
}

func getRowsAround(lines []string, rows int) []string {
	// IDK how this works, it just does
	top := math.Max(0, float64(rows-LinesAround+1))
	bottom := math.Min(float64(len(lines)), float64(rows+LinesAround))

	return lines[int(top):int(bottom)]
}
