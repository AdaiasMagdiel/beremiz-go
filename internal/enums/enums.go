package enums

type TokenType string

const (
	Int   TokenType = "INT"
	Float TokenType = "FLOAT"

	Show TokenType = "SHOW"

	Plus TokenType = "+"
)

type Loc struct {
	Col  int
	Line int
}

type Token struct {
	Type    TokenType
	Literal string
	Loc     Loc
}
