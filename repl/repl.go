package repl

import (
	"Monkey/evaluator"
	"Monkey/lexer"
	"Monkey/parser"
	"bufio"
	"fmt"
	"io"
)

// Console Prompt header
const PROMPT = ">> "

// Print parsing errors
func PrintParserErrors(out io.Writer, errors []*parser.ParseError) {
	for _, err := range errors {
		message := fmt.Sprintf("On %d:%d, %s, in %q",
			err.RowNumber, err.ColumnNumber, err.Message, err.Filename)
		io.WriteString(out, message+"\n")
	}
}

// Start the REPL by repeating asking for input
func Start(in io.Reader, out io.Writer) {

	// Create a new Scanner
	scanner := bufio.NewScanner(in)

	for {
		// Display Prompt Header
		fmt.Printf(PROMPT)

		// Advance Scanner buffer
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Retrieve Buffer Text
		line := scanner.Text()
		// Create new lexer
		l := lexer.New(line, "REPL")
		// Create new parser
		p := parser.New(l)

		// Parse Program
		program := p.ParseProgram()
		if p.HasError() {
			PrintParserErrors(out, p.Errors())
			continue
		}

		// Eval it
		evaluated := evaluator.Eval(program)

		// Print the parsed program out
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
