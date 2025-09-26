package tokens

type TokenType string

const (
	Int    TokenType = "INT"
	Float  TokenType = "FLOAT"
	String TokenType = "STRING"

	Bool TokenType = "BOOL"
	Nil  TokenType = "NIL"

	Identifier TokenType = "IDENTIFIER"

	For  TokenType = "FOR"
	If   TokenType = "IF"
	Else TokenType = "ELSE"
	Elif TokenType = "ELIF"
	Do   TokenType = "DO"
	End  TokenType = "END"

	Write   TokenType = "WRITE"
	Writeln TokenType = "WRITELINE"
	Type    TokenType = "TYPE"

	Plus  TokenType = "PLUS"
	Minus TokenType = "MINUS"
	Times TokenType = "TIMES"
	Div   TokenType = "DIV"

	Eq  TokenType = "EQUALS"
	Neq TokenType = "NOT_EQUALS"
	Lt  TokenType = "LOWER_THAN"
	Gt  TokenType = "GREATER_THAN"

	Dup  TokenType = "DUP"
	Pop  TokenType = "POP"
	Swap TokenType = "SWAP"

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
	"<": Lt,
	">": Gt,
}

func IsOperator(args ...byte) bool {
	_, ok := Operators[string(args)]
	return ok
}

var Keywords map[string]TokenType = map[string]TokenType{
	"write":   Write,
	"writeln": Writeln,
	"type":    Type,

	"nil": Nil,

	"for":  For,
	"if":   If,
	"else": Else,
	"elif": Elif,
	"do":   Do,
	"end":  End,

	"eq":   Eq,
	"neq":  Neq,
	"dup":  Dup,
	"pop":  Pop,
	"swap": Swap,
}

func IsKeyword(content string) bool {
	_, ok := Keywords[content]
	return ok
}
