package repl

import (
	"Monkey/lexer"
	"Monkey/token"
	"bufio"
	"fmt"
	"io"
)

// Console Prompt header
const PROMPT = ">> "

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

		// Loop through all tokens, print them out
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
