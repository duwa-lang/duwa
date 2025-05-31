package token

type TokenType string

const (
	// Single character tokens
	EOF             = "EOF"
	MINUS           = "-"
	ASTERISK        = "*"
	SLASH           = "/"
	PLUS            = "+"
	SEMICOLON       = ";"
	GREATER_THAN    = ">"
	LESS_THAN       = "<"
	ASSIGN          = "="
	COLON           = ":"
	COMMA           = ","
	OPENING_BRACE   = "{"
	CLOSING_BRACE   = "}"
	OPENING_PAREN   = "("
	CLOSING_PAREN   = ")"
	OPENING_BRACKET = "["
	CLOSING_BRACKET = "]"
	FULL_STOP       = "."
	BANG            = "!"
	AMPERSAND       = "&"
	OR              = "|"

	// One or two character token
	GREATER_THAN_OR_EQUAL_TO = ">="
	LESS_THAN_OR_EQUAL_TO    = "<="
	EQUAL_TO                 = "=="
	NOT_EQUAL_TO             = "!="
	MINUS_EQUAL              = "-="
	SLASH_EQUAL              = "/="
	STAR_EQUAL               = "*="
	AND_AND                  = "&&"
	OR_OR                    = "||"
	PLUS_EQUAL               = "+="
	PLUS_PLUS                = "++"
	MINUS_MINUS              = "--"

	ILLEGAL

	COMMENT           = "//"
	MULTILINE_COMMENT = "/*"

	IDENT = "identifier"

	// Literals
	INT = "INT"
	STR = "STR"

	// Keywords
	INTEGER  = "INTEGER"
	STRING   = "STRING"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	WHILE    = "WHILE"
	MAP      = "MAP"
	CLASS    = "CLASS"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	NULL     = "NULL"
	IMPORT   = "IMPORT"
	FROM     = "FROM"
)

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
	File    string
}

var keywords = map[string]TokenType{
	"nambala":    INTEGER,
	"mawu":       STRING,
	"zoona":      TRUE,
	"bodza":      FALSE,
	"ngati":      IF,
	"kapena":     ELSE,
	"bweza":      RETURN,
	"ndondomeko": FUNCTION,
	"za":         FOR,
	"pamene":     WHILE,
	"mgwirizano": MAP,
	"kalasi":     CLASS,
	"siya":       BREAK,
	"pitirizani": CONTINUE,
	"palibe":     NULL,
	"tenga":      IMPORT,
	"mu":  FROM,
}

var variableTypes = map[TokenType]TokenType{
	INTEGER: INTEGER,
	STRING:  STRING,
	MAP:     MAP,
}

func LookupVariableType(ident TokenType) TokenType {
	if tok, ok := variableTypes[ident]; ok {
		return tok
	}
	return ""
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func BooleanToString(b bool) string {
	if b {
		return "zoona"
	}
	return "bodza"
}
