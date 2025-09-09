package enums

type TokenType string

const (
	Int TokenType = "INT"

	Show TokenType = "SHOW"

	Plus TokenType = "+"
)

type Token struct {
	Type TokenType
	Literal string
}
