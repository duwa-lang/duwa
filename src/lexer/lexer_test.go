package lexer

import (
	"testing"

	"github.com/duwa-lang/duwa/src/token"
)

func TestNextToken(t *testing.T) {
	input := `
	nambala phatikiza(yambi, chiwiri) {
		bweza yamba + chiwiri;
	}
 	"foobar"
	"foo bar"
	[1, 2];
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INTEGER, "nambala"},
		{token.IDENT, "phatikiza"},
		{token.OPENING_PAREN, "("},
		{token.IDENT, "yambi"},
		{token.COMMA, ","},
		{token.IDENT, "chiwiri"},
		{token.CLOSING_PAREN, ")"},
		{token.OPENING_BRACE, "{"},
		{token.RETURN, "bweza"},
		{token.IDENT, "yamba"},
		{token.PLUS, "+"},
		{token.IDENT, "chiwiri"},
		{token.SEMICOLON, ";"},
		{token.CLOSING_BRACE, "}"},
		{token.STR, "foobar"},
		{token.STR, "foo bar"},
		{token.OPENING_BRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.CLOSING_BRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := NewLexel()
	l.Handle([]byte(input))
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
