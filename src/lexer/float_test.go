package lexer

import (
	"testing"

	"github.com/duwa-lang/duwa/src/token"
)

func TestFloatLiterals(t *testing.T) {
	input := `1.0 0.5 3.14159 100.25`

	l := NewLexel()
	l.Handle([]byte(input))

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "1.0"},
		{token.INT, "0.5"},
		{token.INT, "3.14159"},
		{token.INT, "100.25"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestFloatVsMethodCall(t *testing.T) {
	// Make sure we distinguish between "1.5" (float) and "arr.length" (method call)
	input := `arr.length 1.5 x.method()`

	l := NewLexel()
	l.Handle([]byte(input))

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "arr"},
		{token.FULL_STOP, "."},
		{token.IDENT, "length"},
		{token.INT, "1.5"},
		{token.IDENT, "x"},
		{token.FULL_STOP, "."},
		{token.IDENT, "method"},
		{token.OPENING_PAREN, "("},
		{token.CLOSING_PAREN, ")"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestFloatArithmetic(t *testing.T) {
	// Test that floats work in arithmetic expressions
	input := `1.5 + 2.5`

	l := NewLexel()
	l.Handle([]byte(input))

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "1.5"},
		{token.PLUS, "+"},
		{token.INT, "2.5"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
