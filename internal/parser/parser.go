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
	lines        []string
}

type FlowAddr struct {
	addr  int
	token tokens.Token
}

func New(tokens []tokens.Token, errorHandler func(), lines []string) *Parser {
	if errorHandler == nil {
		errorHandler = func() {}
	}

	return &Parser{
		Tokens:       tokens,
		pos:          0,
		errorHandler: errorHandler,
		lines:        lines,
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
		case tokens.Lt:
			return x < y, tokens.Bool, nil
		case tokens.Gt:
			return x > y, tokens.Bool, nil
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
		case tokens.Lt:
			return x < y, tokens.Bool, nil
		case tokens.Gt:
			return x > y, tokens.Bool, nil
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

func (p Parser) peek() tokens.Token {
	return p.Tokens[p.pos]
}

func (p *Parser) consume() tokens.Token {
	token := p.Tokens[p.pos]
	p.pos++

	return token
}

func (p Parser) isAtEnd() bool {
	return p.pos >= len(p.Tokens) || p.Tokens[p.pos].Type == tokens.EOF
}

func (p Parser) handleControlFlow() {
	addrInfo := []FlowAddr{}
	var top FlowAddr
	var e error

	// Tracks the active block type: "if" (1), "for" (2), or none (0).
	// We need to distinguish them, since "elif"/"else" are only valid inside "if",
	// and "for" needs a loop-back connection not present in "if".
	var blockType uint8 = 0

	for idx, token := range p.Tokens {
		switch token.Type {
		case tokens.If:
			// When encountering "if", we mark the start of a new conditional block.
			// Push it onto the stack so we can later match it with a "do" and an "end".
			blockType = 1
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})

		case tokens.For:
			// When encountering "for", we mark the start of a new conditional block.
			// Push it onto the stack so we can later match it with a "do" and an "end".
			blockType = 2
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})

		case tokens.Elif, tokens.Else:
			// "elif" and "else" can only exist inside a valid "if" block.
			// If no "if" was opened, the control flow is invalid.
			if blockType != 1 {
				err.SyntaxError(
					token,
					fmt.Sprintf("'%s' must follow an 'if ... do' or 'elif ... do' block.", token.Literal),
					p.lines,
				)
				p.errorHandler()
				break
			}

			// Each new "elif" or "else" must close the previous "do" block.
			// This ensures that execution can jump to the correct branch.
			addrInfo, top, e = Pop(addrInfo)

			if e != nil || top.token.Type != tokens.Do {
				// If we didn't find a "do", then the user wrote an invalid structure.
				// This protects against "elif"/"else" without a proper "if ... do" or "elif ... do".
				err.SyntaxError(
					token,
					fmt.Sprintf("Invalid '%s' usage. Expected 'if ... do' or 'elif ... do'.", token.Literal),
					p.lines,
				)
				p.errorHandler()
				break
			}

			// The closed "do" now knows where execution continues:
			// immediately after this "elif" or "else".
			p.Tokens[top.addr].JmpTo = idx + 1

			// Push the new "elif"/"else" as part of the current block chain.
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})

		case tokens.Do:
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})

		case tokens.End:
			// "end" must close a previously opened "if".
			// If the stack is empty, the code has an extra "end".
			if len(addrInfo) == 0 {
				err.SyntaxError(
					token,
					"Invalid 'end' usage. No matching 'if' block found.",
					p.lines,
				)
				p.errorHandler()
				break
			}

			// Special case: handle "for .. do .. end".
			if blockType == 2 {
				if len(addrInfo) < 2 {
					// Why: a valid for-loop requires at least a "for" + "do".
					err.SyntaxError(
						token,
						"Invalid 'end' usage. No matching 'for .. do' block found.",
						p.lines,
					)
					p.errorHandler()
					break
				}

				forFlow := addrInfo[len(addrInfo)-2]
				doFlow := addrInfo[len(addrInfo)-1]

				if forFlow.token.Type != tokens.For || doFlow.token.Type != tokens.Do {
					// Why: structure must be exactly "for .. do .. end".
					err.SyntaxError(
						token,
						"Invalid 'end' usage. No matching 'for .. do' block found.",
						p.lines,
					)
					p.errorHandler()
					break
				}

				// Why: loop exit jumps back to "for",
				// and "do" body knows where to continue if loop ends.
				p.Tokens[idx].JmpTo = forFlow.addr
				p.Tokens[doFlow.addr].JmpTo = idx + 1
				addrInfo = addrInfo[:len(addrInfo)-2]
				continue
			}

			// We unwind the stack until we find the matching "if".
			// Along the way, we patch all "do"/"elif"/"else" so they
			// jump right after this "end", ensuring correct flow.
			for {
				addrInfo, top, e = Pop(addrInfo)
				if e != nil {
					// If we exhaust the stack before finding "if",
					// the structure is unbalanced.
					err.SyntaxError(
						token,
						"Unbalanced 'end'. No matching 'if' block found.",
						p.lines,
					)
					p.errorHandler()
					break
				}

				if top.token.Type != tokens.If {
					p.Tokens[top.addr].JmpTo = idx + 1
				}

				// Once we find the opening "if", the block is complete.
				if top.token.Type == tokens.If {
					break
				}
			}
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
			tokens.Bool,
			tokens.Nil:
			stack = append(stack, p.consume())

		case tokens.Plus,
			tokens.Minus,
			tokens.Times,
			tokens.Div,
			tokens.Lt,
			tokens.Gt:

			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			if (a.Type != tokens.Int && a.Type != tokens.Float) ||
				(b.Type != tokens.Int && b.Type != tokens.Float) {
				err.SyntaxError(token, fmt.Sprintf(
					"Operator '%s' expects int or float.", token.Literal), p.lines)
				p.errorHandler()
				return
			}

			res, e := evalNumBin(token, a, b)
			if e != nil {
				err.SyntaxError(token, e.Error(), p.lines)
				p.errorHandler()
				return
			}
			stack = append(stack, res)
			p.consume()

		case tokens.Write,
			tokens.Writeln:

			if len(stack) == 0 {
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal), p.lines)
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
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal), p.lines)
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

		case tokens.Dup:
			if len(stack) == 0 {
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal), p.lines)
				p.errorHandler()
				return
			}

			a := stack[len(stack)-1]
			p.consume()

			stack = append(stack, a)

		case tokens.Eq:
			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			p.consume()

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			stack = append(stack, tokens.Token{
				Type:    tokens.Bool,
				Literal: a.Literal == b.Literal,
				Loc:     token.Loc,
				JmpTo:   0,
			})

		case tokens.Neq:
			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			p.consume()

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			stack = append(stack, tokens.Token{
				Type:    tokens.Bool,
				Literal: a.Literal != b.Literal,
				Loc:     token.Loc,
				JmpTo:   0,
			})

		case tokens.If,
			tokens.For:
			p.consume()

		case tokens.Elif,
			tokens.Else:
			p.pos = p.consume().JmpTo

		case tokens.Do:
			var top tokens.Token
			var e error
			stack, top, e = Pop(stack)

			if e != nil {
				err.SyntaxError(token, "The 'do' keyword requires value in stack. Stack is empty.", p.lines)
				p.errorHandler()
				return
			}

			var cond bool
			switch top.Type {
			case tokens.Bool:
				cond = top.Literal.(bool)
			case tokens.Nil:
				cond = false
			case tokens.Int:
				cond = top.Literal != 0
			case tokens.Float:
				cond = top.Literal != 0.0
			case tokens.String:
				cond = top.Literal != ""
			default:
				err.SyntaxError(token, fmt.Sprintf(
					"Invalid condition type '%s' cannot be used in a boolean context", top.Type,
				), p.lines)
				p.errorHandler()
				return
			}

			p.consume()

			if !cond {
				p.pos = token.JmpTo
			}

		case tokens.End:
			if token.JmpTo > 0 {
				p.pos = token.JmpTo
			}

			p.consume()

		default:
			err.Error(fmt.Sprintf("Not implemented case for TokenType '%s'.", token.Type))
			fmt.Fprintf(os.Stderr, "\x1b[31m%s:%d:%d:\x1b[0m '%s'", token.Loc.File, token.Loc.Line, token.Loc.Col, token.Literal)
			p.errorHandler()
			return
		}
	}
}
