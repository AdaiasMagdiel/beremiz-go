package parser

import (
	"fmt"
	"os"
	"strings"

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
				return nil, tokens.Nil, fmt.Errorf("division by zero")
			}
			return float64(x) / float64(y), tokens.Float, nil
		default:
			return nil, tokens.Nil, fmt.Errorf("unsupported op: %s", op.Type)
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
				return nil, tokens.Nil, fmt.Errorf("division by zero")
			}
			return x / y, tokens.Float, nil
		default:
			return nil, tokens.Nil, fmt.Errorf("unsupported op: %s", op.Type)
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

func (p *Parser) handleControlFlow() {
	var ctrl []int

	for idx, token := range p.Tokens {
		switch token.Type {

		case tokens.If:
			ctrl = append(ctrl, idx)

		case tokens.Else:

			var top int
			var e error
			ctrl, top, e = Pop(ctrl)

			if e != nil || top < 0 {
				msg := "Unexpected 'else' — no active 'if' block"
				if top < 0 {
					msg = "Duplicate 'else' — 'if' block already has an 'else'"
				}
				err.SyntaxError(token, msg)
				p.errorHandler()
				return
			}

			p.Tokens[top].JmpTo = idx + 1

			ctrl = append(ctrl, -idx)

		case tokens.Endif:
			var top int
			var e error
			ctrl, top, e = Pop(ctrl)

			if e != nil {
				err.SyntaxError(token, "Unexpected 'endif' — no open 'if' block")
				p.errorHandler()
				return
			}

			if top < 0 {
				elseIdx := -top
				p.Tokens[elseIdx].JmpTo = idx + 1
			} else {
				ifIdx := top
				p.Tokens[ifIdx].JmpTo = idx + 1
			}
		}
	}

	for len(ctrl) > 0 {
		var top int
		ctrl, top, _ = Pop(ctrl)
		if top < 0 {
			elseIdx := -top
			err.SyntaxError(p.Tokens[elseIdx], "Syntax error: 'else' without matching 'endif'")
		} else {
			ifIdx := top
			err.SyntaxError(p.Tokens[ifIdx], "Syntax error: 'if' without matching 'endif'")
		}
	}
}

func (p *Parser) Eval() {
	var stack = []tokens.Token{}

	p.handleControlFlow()

	for {
		if p.isAtEnd() {
			break
		}

		token := p.peek()

		switch token.Type {
		case tokens.Int,
			tokens.Float,
			tokens.String,
			tokens.True,
			tokens.False,
			tokens.Nil:
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

		case tokens.Write,
			tokens.Writeln:

			if len(stack) == 0 {
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal))
				p.errorHandler()
				return
			}

			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			p.consume()

			if token.Type == tokens.Write {
				fmt.Print(a.Literal)
			} else {
				fmt.Println(a.Literal)
			}

		case tokens.Type:
			if len(stack) == 0 {
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal))
				p.errorHandler()
				return
			}

			a := stack[len(stack)-1]
			p.consume()

			if tokens.IsKeyword(strings.ToLower(string(a.Type))) {
				stack = append(stack, tokens.Token{
					Type:    tokens.String,
					Literal: "KEYWORD",
					Loc:     token.Loc,
				})
			} else {
				stack = append(stack, tokens.Token{
					Type:    tokens.String,
					Literal: string(a.Type),
					Loc:     token.Loc,
				})
			}

		case tokens.If:
			var top tokens.Token
			var e error
			stack, top, e = Pop(stack)

			if e != nil {
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal))
				p.errorHandler()
				return
			}

			var cond bool
			switch {
			case top.Type == tokens.True:
				cond = true
			case top.Type == tokens.False,
				top.Type == tokens.Nil:
				cond = false
			case top.Type == tokens.Int:
				cond = top.Literal != 0
			case top.Type == tokens.Float:
				cond = top.Literal != 0.0
			case top.Type == tokens.String:
				cond = top.Literal != ""
			default:
				err.SyntaxError(token, fmt.Sprintf(
					"Invalid condition type '%s' cannot be used in a boolean context", top.Type,
				))
				p.errorHandler()
				return
			}

			p.consume()

			if !cond {
				p.pos = token.JmpTo
			}

		case tokens.Else:
			p.pos = p.consume().JmpTo

		case tokens.Endif:
			p.consume()

		default:
			err.Error(fmt.Sprintf("Not implemented case for TokenType '%s'.", token.Type))
			fmt.Fprintf(os.Stderr, "\x1b[31m%s:%d:%d:\x1b[0m '%s'", token.Loc.File, token.Loc.Line, token.Loc.Col, token.Literal)
			p.errorHandler()
			return
		}
	}
}
