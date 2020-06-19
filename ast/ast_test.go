package ast

import (
	"Monkey/token"
	"testing"
)

func TestToString(t *testing.T) {

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

	if program.ToString() != "let foo = bar;" {
		t.Errorf("program.ToString() incorrect. got=%q",
			program.ToString())
	}
}
