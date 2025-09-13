package tokens

type TokenType string

const (
	Int    TokenType = "INT"
	Float  TokenType = "FLOAT"
	String TokenType = "STRING"

	Identifier TokenType = "IDENTIFIER"
	Write      TokenType = "WRITE"
	Writeln    TokenType = "WRITELINE"

	Plus  TokenType = "PLUS"
	Minus TokenType = "MINUS"
	Times TokenType = "TIMES"
	Div   TokenType = "DIV"

	EOF TokenType = "EOF"
	NIL TokenType = "NIL"
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
}

func IsKeyword(content string) bool {
	_, ok := Keywords[content]
	return ok
}
