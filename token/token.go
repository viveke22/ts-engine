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
	IDENT         = "IDENT"  // add, foobar, x, y, ...
	INT           = "INT"    // 1343456
	STRING        = "STRING" // "foobar"
	TEMPLATE_HEAD = "TEMPLATE_HEAD"
	TEMPLATE_MID  = "TEMPLATE_MID"
	TEMPLATE_TAIL = "TEMPLATE_TAIL"

	// Operators
	ASSIGN      = "="
	PLUS        = "+"
	MINUS       = "-"
	BANG        = "!"
	ASTERISK    = "*"
	SLASH       = "/"
	MOD         = "%"
	ARROW       = "=>"
	PLUS_PLUS   = "++"
	MINUS_MINUS = "--"

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

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	EXPORT    = "EXPORT"
	CONST     = "CONST"
	VAR       = "VAR"
	AWAIT     = "AWAIT"
	DECLARE   = "DECLARE"
	IMPORT    = "IMPORT"
	FROM      = "FROM"
	AS        = "AS"
	FOR       = "FOR"
	WHILE     = "WHILE"
	BREAK     = "BREAK"
	CONTINUE  = "CONTINUE"
	NULL      = "NULL"
	UNDEFINED = "UNDEFINED"
	CLASS     = "CLASS"
	NEW       = "NEW"
	THIS      = "THIS"
)

var keywords = map[string]TokenType{
	"function":  FUNCTION,
	"let":       LET,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"export":    EXPORT,
	"const":     CONST,
	"var":       VAR,
	"await":     AWAIT,
	"declare":   DECLARE,
	"import":    IMPORT,
	"from":      FROM,
	"as":        AS,
	"for":       FOR,
	"while":     WHILE,
	"break":     BREAK,
	"continue":  CONTINUE,
	"null":      NULL,
	"undefined": UNDEFINED,
	"class":     CLASS,
	"new":       NEW,
	"this":      THIS,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
