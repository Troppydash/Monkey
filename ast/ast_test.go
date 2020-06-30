package ast

import (
	"Monkey/options"
	"Monkey/token"
	"testing"
)

// Test the ast tostring method
func TestToString(t *testing.T) {

	if options.NicerToString {
		return
	}

	// let foo = bar;
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "bar"},
					Value: "bar",
				},
			},
		},
	}

	if program.ToString() != "[(let (foo) = (bar);)]" {
		t.Errorf("program.ToString() incorrect. got=%q",
			program.ToString())
	}
}
