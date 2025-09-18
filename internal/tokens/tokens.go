package tokens

type TokenType string

const (
	Int    TokenType = "INT"
	Float  TokenType = "FLOAT"
	String TokenType = "STRING"

	True  TokenType = "TRUE"
	False TokenType = "FALSE"
	Nil   TokenType = "NIL"

	Identifier TokenType = "IDENTIFIER"

	If    TokenType = "IF"
	Else  TokenType = "ELSE"
	Endif TokenType = "ENDIF"

	Write   TokenType = "WRITE"
	Writeln TokenType = "WRITELINE"
	Type    TokenType = "TYPE"

	Plus  TokenType = "PLUS"
	Minus TokenType = "MINUS"
	Times TokenType = "TIMES"
	Div   TokenType = "DIV"

	EOF TokenType = "EOF"
)

type Loc struct {
	File string
	Col  int
	Line int
}

type Token struct {
	Type    TokenType
	Literal any
	Loc     Loc
	JmpTo   int
}

var Operators map[string]TokenType = map[string]TokenType{
	"+": Plus,
	"-": Minus,
	"*": Times,
	"/": Div,
}

func IsOperator(args ...byte) bool {
	_, ok := Operators[string(args)]
	return ok
}

var Keywords map[string]TokenType = map[string]TokenType{
	"write":   Write,
	"writeln": Writeln,
	"type":    Type,
	"if":      If,
	"else":    Else,
	"endif":   Endif,
}

func IsKeyword(content string) bool {
	_, ok := Keywords[content]
	return ok
}
