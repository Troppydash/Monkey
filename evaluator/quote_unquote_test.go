package evaluator

import (
	"Monkey/object"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(5)`,
			`5`,
		},
		{
			`quote(5 + 8)`,
			`(5 + 8)`,
		},
		{
			`quote(foobar)`,
			`foobar`,
		},
		{
			`quote(foobar + barfoo)`,
			`(foobar + barfoo)`,
		},
		{
			`quote(unquote(true))`,
			`true`,
		},
		{
			`quote(unquote(true == false))`,
			`false`,
		},
		{
			`quote(unquote(quote(4 + 4)))`,
			`(4 + 4)`,
		},
		{
			`let quotedInfi = quote(4 + 4)
			quote(unquote(4 + 4) + unquote(quotedInfi))`,
			`(8 + (4 + 4))`,
		},
	}

	for _, tt := range tests {
		evaluated := CheckEvalNice(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote. got=%T (%+v)",
				evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.ToString() != tt.expected {
			t.Errorf("not equal. got=%q, want=%q",
				quote.Node.ToString(), tt.expected)
		}
	}
}

func TestQuoteUnQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(unquote(4))`,
			`4`,
		},
		{
			`quote(unquote(4 + 4))`,
			`8`,
		},
		{
			`quote(8 + unquote(4 + 4))`,
			`(8 + 8)`,
		},
		{
			`quote(unquote(4 + 4) + 8)`,
			`(8 + 8)`,
		},
		{
			`let foobar = 8
					quote(foobar)`,
			`foobar`,
		},
		{
			`let foobar = 8
					quote(unquote(foobar))`,
			`8`,
		},
	}

	for _, tt := range tests {
		evaluated := CheckEvalNice(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote. got=%T (%+v)",
				evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.ToString() != tt.expected {
			t.Errorf("not equal. got=%q, want=%q",
				quote.Node.ToString(), tt.expected)
		}
	}
}
