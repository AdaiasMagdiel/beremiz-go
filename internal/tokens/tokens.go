package tokens

type TokenType string

const (
	Int   TokenType = "INT"
	Float TokenType = "FLOAT"

	Identifier TokenType = "IDENTIFIER"
	Show       TokenType = "SHOW"

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
	"show": Show,
}

func IsKeyword(content string) bool {
	_, ok := Keywords[content]
	return ok
}
