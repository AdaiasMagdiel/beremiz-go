package parser

import (
	"fmt"
	"os"

	"github.com/adaiasmagdiel/beremiz-go/internal/err"
	"github.com/adaiasmagdiel/beremiz-go/internal/tokens"
)

type Parser struct {
	Tokens       []tokens.Token
	pos          int
	errorHandler func()
}

func New(tokens []tokens.Token, errorHandler func()) *Parser {
	if errorHandler == nil {
		errorHandler = func() {}
	}

	return &Parser{
		Tokens:       tokens,
		pos:          0,
		errorHandler: errorHandler,
	}
}

// evalNumBin aplica um operador numérico (+, -, *, /) em dois tokens numéricos.
// Regras:
// - int op int => int (exceto /, que vira float)
// - qualquer caso com float => float
// - checa divisão por zero
func evalNumBin(op, a, b tokens.Token) (tokens.Token, error) {
	intOp := func(x, y int64) (any, tokens.TokenType, error) {
		switch op.Type {
		case tokens.Plus:
			return x + y, tokens.Int, nil
		case tokens.Minus:
			return x - y, tokens.Int, nil
		case tokens.Times:
			return x * y, tokens.Int, nil
		case tokens.Div:
			if y == 0 {
				return nil, tokens.NIL, fmt.Errorf("division by zero")
			}
			return float64(x) / float64(y), tokens.Float, nil
		default:
			return nil, tokens.NIL, fmt.Errorf("unsupported op: %s", op.Type)
		}
	}
	floatOp := func(x, y float64) (any, tokens.TokenType, error) {
		switch op.Type {
		case tokens.Plus:
			return x + y, tokens.Float, nil
		case tokens.Minus:
			return x - y, tokens.Float, nil
		case tokens.Times:
			return x * y, tokens.Float, nil
		case tokens.Div:
			if y == 0 {
				return nil, tokens.NIL, fmt.Errorf("division by zero")
			}
			return x / y, tokens.Float, nil
		default:
			return nil, tokens.NIL, fmt.Errorf("unsupported op: %s", op.Type)
		}
	}

	switch av := a.Literal.(type) {
	case int64:
		switch bv := b.Literal.(type) {
		case int64:
			val, ty, err := intOp(av, bv)
			if err != nil {
				return tokens.Token{}, err
			}
			return tokens.Token{Type: ty, Literal: val, Loc: op.Loc}, nil
		case float64:
			val, ty, err := floatOp(float64(av), bv)
			if err != nil {
				return tokens.Token{}, err
			}
			return tokens.Token{Type: ty, Literal: val, Loc: op.Loc}, nil
		default:
			return tokens.Token{}, fmt.Errorf("unsupported rhs type: %T", b.Literal)
		}

	case float64:
		switch bv := b.Literal.(type) {
		case int64:
			val, ty, err := floatOp(av, float64(bv))
			if err != nil {
				return tokens.Token{}, err
			}
			return tokens.Token{Type: ty, Literal: val, Loc: op.Loc}, nil
		case float64:
			val, ty, err := floatOp(av, bv)
			if err != nil {
				return tokens.Token{}, err
			}
			return tokens.Token{Type: ty, Literal: val, Loc: op.Loc}, nil
		default:
			return tokens.Token{}, fmt.Errorf("unsupported rhs type: %T", b.Literal)
		}

	default:
		return tokens.Token{}, fmt.Errorf("unsupported lhs type: %T", a.Literal)
	}
}

func (p *Parser) peek() tokens.Token {
	return p.Tokens[p.pos]
}

func (p *Parser) consume() tokens.Token {
	token := p.Tokens[p.pos]
	p.pos++

	return token
}

func (p *Parser) isAtEnd() bool {
	return p.pos >= len(p.Tokens) || p.Tokens[p.pos].Type == tokens.EOF
}

func (p *Parser) Eval() {
	var stack = []tokens.Token{}

	for {
		if p.isAtEnd() {
			break
		}

		token := p.peek()

		switch token.Type {
		case tokens.Int,
			tokens.Float:

			stack = append(stack, p.consume())
		case tokens.Plus,
			tokens.Minus,
			tokens.Times,
			tokens.Div:

			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)))
				p.errorHandler()
				return
			}

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			// garanta que são numéricos (opcional, já validamos nos type assertions)
			if (a.Type != tokens.Int && a.Type != tokens.Float) ||
				(b.Type != tokens.Int && b.Type != tokens.Float) {
				err.SyntaxError(token, fmt.Sprintf(
					"Operator '%s' expects int or float.", token.Literal))
				p.errorHandler()
				return
			}

			res, e := evalNumBin(token, a, b)
			if e != nil {
				err.SyntaxError(token, e.Error())
				p.errorHandler()
				return
			}
			stack = append(stack, res)
			p.consume()

		case tokens.Show,
			tokens.ShowLN:
			if len(stack) == 0 {
				err.SyntaxError(token, "This keyword requires value in stack. Stack is empty.")
				p.errorHandler()
				return
			}

			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			p.consume()

			if token.Type == tokens.Show {
				fmt.Print(a.Literal)
			} else {
				fmt.Println(a.Literal)

			}

		default:
			err.Error(fmt.Sprintf("Not implemented case for TokenType '%s'.", token.Type))
			fmt.Fprintf(os.Stderr, "\x1b[31m%s:%d:%d:\x1b[0m '%s'", token.Loc.File, token.Loc.Line, token.Loc.Col, token.Literal)
			p.errorHandler()
			return
		}
	}
}
