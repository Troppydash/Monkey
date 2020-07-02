package repl

import (
	"Monkey/evaluator"
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/options"
	"Monkey/parser"
	"bufio"
	"fmt"
	"io"
	"strings"
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

	// REPL Environment
	env := object.NewEnvironment()

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

		if strings.Contains(line, "--") {
			ParseOptions(out, line)
			continue
		}

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

		io.WriteString(out, "  >> Parsed ")
		io.WriteString(out, program.ToString())
		io.WriteString(out, "\n")

		// Eval it
		evaluated := evaluator.Eval(program, env)

		// Print the parsed program out
		if evaluated != nil {
			io.WriteString(out, "  >> Evaluated ")
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func ParseOptions(out io.Writer, line string) {
	switch line {
	case "--on nicer":
		options.NicerToString = true
		io.WriteString(out, "Enabled Nicer ToString")
	case "--off nicer":
		options.NicerToString = false
		io.WriteString(out, "Disabled Nicer ToString")

	case "--on fatalErrors":
		options.FatalErrors = true
		io.WriteString(out, "Enabled FatalErrors")
	case "--off fatalErrors":
		options.FatalErrors = false
		io.WriteString(out, "Disabled FatalErrors")

	default:
		io.WriteString(out, "No options matching your request has been found")
	}
	io.WriteString(out, "\n")
}
