package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "foobar"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

	LT = "<"
	GT = ">"

	EQ            = "=="
	NOT_EQ        = "!="
	EQ_STRICT     = "==="
	NOT_EQ_STRICT = "!=="
	AND           = "&&"
	OR            = "||"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	DOT       = "."
	COLON     = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	EXPORT   = "EXPORT"
	CONST    = "CONST"
	VAR      = "VAR"
	AWAIT    = "AWAIT"
	DECLARE  = "DECLARE"
	IMPORT   = "IMPORT"
	FROM     = "FROM"
	AS       = "AS"
)

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"export":   EXPORT,
	"const":    CONST,
	"var":      VAR,
	"await":    AWAIT,
	"declare":  DECLARE,
	"import":   IMPORT,
	"from":     FROM,
	"as":       AS,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
