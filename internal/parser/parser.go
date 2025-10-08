package parser

import (
	"bufio"
	"fmt"
	"math"
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
	inLoop       int
	isREPL       bool
}

type FlowAddr struct {
	addr  int
	token tokens.Token
}

func New(tokens []tokens.Token, errorHandler func(), lines []string, isREPL bool) *Parser {
	if errorHandler == nil {
		errorHandler = func() {}
	}

	return &Parser{
		Tokens:       tokens,
		pos:          0,
		errorHandler: errorHandler,
		lines:        lines,
		inLoop:       0,
		isREPL:       isREPL,
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
		case tokens.Le:
			return x <= y, tokens.Bool, nil
		case tokens.Ge:
			return x >= y, tokens.Bool, nil
		case tokens.Exp:
			return math.Pow(float64(x), float64(y)), tokens.Float, nil
		case tokens.Mod:
			if y == 0 {
				return nil, tokens.Nil, fmt.Errorf("modulo by zero")
			}

			r := x % y
			if r != 0 && ((y > 0 && r < 0) || (y < 0 && r > 0)) {
				r += y
			}
			return r, tokens.Int, nil
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
		case tokens.Le:
			return x <= y, tokens.Bool, nil
		case tokens.Ge:
			return x >= y, tokens.Bool, nil
		case tokens.Exp:
			return math.Pow(x, y), tokens.Float, nil
		case tokens.Mod:
			if y == 0 {
				return nil, tokens.Nil, fmt.Errorf("modulo by zero")
			}

			r := math.Mod(x, y)
			if r != 0 && ((y > 0 && r < 0) || (y < 0 && r > 0)) {
				r += y
			}
			return r, tokens.Float, nil
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

func (p Parser) handleControlFlow() map[string][]tokens.Token {
	addrInfo := []FlowAddr{}
	var top FlowAddr
	var e error

	defs := make(map[string][]tokens.Token)
	var keys []string

	var blockStack []BlockType

	var idx int = 0

	for {
		if idx >= len(p.Tokens) || p.Tokens[idx].Type == tokens.EOF {
			break
		}

		token := p.Tokens[idx]

		switch token.Type {
		case tokens.If:
			blockStack = append(blockStack, BlockIf)
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})
			idx++

		case tokens.Elif, tokens.Else:
			if len(blockStack) == 0 || blockStack[len(blockStack)-1] != BlockIf {
				err.SyntaxError(token,
					fmt.Sprintf("'%s' must follow an 'if ... do' or 'elif ... do' block.", token.Literal),
					p.lines,
				)
				p.errorHandler()
				break
			}

			addrInfo, top, e = Pop(addrInfo)
			if e != nil || top.token.Type != tokens.Do {
				err.SyntaxError(token,
					fmt.Sprintf("Invalid '%s' usage. Expected 'if ... do' or 'elif ... do'.", token.Literal),
					p.lines,
				)
				p.errorHandler()
				break
			}

			p.Tokens[top.addr].JmpTo = idx + 1
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})
			idx++

		case tokens.For:
			blockStack = append(blockStack, BlockFor)
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})
			idx++

		case tokens.Do:
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})
			idx++

		case tokens.Define:
			blockStack = append(blockStack, BlockDefine)
			addrInfo = append(addrInfo, FlowAddr{addr: idx, token: token})

			if idx+1 >= len(p.Tokens) || p.Tokens[idx+1].Type != tokens.Identifier {
				err.SyntaxError(token,
					fmt.Sprintf("Expected identifier after 'define' keyword, but got '%s'.",
						strings.ToLower(string(p.Tokens[idx+1].Type))),
					p.lines,
				)
				p.errorHandler()
				break
			}

			key := p.Tokens[idx+1].Literal.(string)
			defs[key] = []tokens.Token{}
			keys = append(keys, key)

			idx += 2
			continue

		case tokens.End:
			if len(blockStack) == 0 {
				err.SyntaxError(token,
					"Invalid 'end' usage. No matching block found.",
					p.lines,
				)
				p.errorHandler()
				break
			}

			current := blockStack[len(blockStack)-1]
			blockStack = blockStack[:len(blockStack)-1]

			switch current {
			case BlockFor:
				if len(addrInfo) < 2 {
					err.SyntaxError(token,
						"Invalid 'end' usage. No matching 'for .. do' block found.",
						p.lines,
					)
					p.errorHandler()
					break
				}
				forFlow := addrInfo[len(addrInfo)-2]
				doFlow := addrInfo[len(addrInfo)-1]

				p.Tokens[idx].JmpTo = forFlow.addr
				p.Tokens[doFlow.addr].JmpTo = idx + 1
				addrInfo = addrInfo[:len(addrInfo)-2]

			case BlockDefine:
				if len(addrInfo) == 0 {
					err.SyntaxError(token,
						"Invalid 'end' usage. No matching 'define' block found.",
						p.lines,
					)
					p.errorHandler()
					break
				}
				defineFlow := addrInfo[len(addrInfo)-1]
				if defineFlow.token.Type != tokens.Define {
					err.SyntaxError(token,
						fmt.Sprintf("Mismatched 'end' block. Expected to close 'define', but found '%s'.",
							defineFlow.token.Literal),
						p.lines,
					)
					p.errorHandler()
					break
				}
				p.Tokens[defineFlow.addr].JmpTo = idx + 1
				addrInfo = addrInfo[:len(addrInfo)-1]
				keys = keys[:len(keys)-1]

			case BlockIf:
				for {
					addrInfo, top, e = Pop(addrInfo)
					if e != nil {
						err.SyntaxError(token,
							"Unbalanced 'end'. No matching 'if' block found.",
							p.lines,
						)
						p.errorHandler()
						break
					}
					if top.token.Type != tokens.If {
						p.Tokens[top.addr].JmpTo = idx + 1
					}
					if top.token.Type == tokens.If {
						break
					}
				}
			}
			idx++
			continue

		default:
			if len(blockStack) > 0 && blockStack[len(blockStack)-1] == BlockDefine {
				if len(keys) == 0 {
					err.SyntaxError(token,
						fmt.Sprintf("Expected identifier after 'define' keyword, but got '%s'.",
							strings.ToLower(string(token.Type))),
						p.lines,
					)
					p.errorHandler()
					break
				}
				key := keys[len(keys)-1]
				defs[key] = append(defs[key], token)
			}
			idx++
		}
	}

	return defs
}

func expandDefs(defs map[string][]tokens.Token) {
	changed := true

	for i := 0; i < 20 && changed; i++ {
		changed = false

		for key, body := range defs {
			var expanded []tokens.Token

			for _, t := range body {
				if t.Type == tokens.Identifier {
					name := fmt.Sprintf("%v", t.Literal)
					if inner, ok := defs[name]; ok {
						expanded = append(expanded, inner...)
						changed = true
						continue
					}
				}
				expanded = append(expanded, t)
			}

			defs[key] = expanded
		}
	}
}

func (p *Parser) expandBlocks(defs map[string][]tokens.Token) {
	var expanded []tokens.Token

	for idx := 0; idx < len(p.Tokens); idx++ {
		tok := p.Tokens[idx]

		if tok.Type == tokens.Define {
			end := tok.JmpTo
			if end == 0 || end > len(p.Tokens) {
				depth := 1
				j := idx + 1
				for j < len(p.Tokens) && depth > 0 {
					switch p.Tokens[j].Type {
					case tokens.Define:
						depth++
					case tokens.End:
						depth--
					default:
						// no-op
					}
					j++
				}
				end = j
			}
			for i := idx; i < end; i++ {
				expanded = append(expanded, p.Tokens[i])
			}
			idx = end - 1
			continue
		}

		if tok.Type == tokens.Identifier {
			key := fmt.Sprintf("%v", tok.Literal)
			if body, ok := defs[key]; ok {
				for _, t := range body {
					clone := t
					if clone.Loc.File == "" {
						clone.Loc = tok.Loc
					}
					expanded = append(expanded, clone)
				}
				continue
			}
		}

		expanded = append(expanded, tok)
	}

	p.Tokens = expanded
}

func (p *Parser) Eval() {
	var stack = []tokens.Token{}

	defs := p.handleControlFlow()
	expandDefs(defs)
	p.expandBlocks(defs)
	p.handleControlFlow()

	var outputBuffer = bufio.NewWriter(os.Stdout)

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
			tokens.Gt,
			tokens.Le,
			tokens.Ge,
			tokens.Exp,
			tokens.Mod:

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

		case tokens.Concat:
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
			p.consume()

			result := toString(a) + toString(b)

			stack = append(stack, tokens.Token{
				Type:    tokens.String,
				Literal: result,
				Loc:     token.Loc,
			})

		case tokens.And, tokens.Or:
			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			left := stack[len(stack)-2]
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			condLeft := toBool(left)
			condRight := toBool(right)

			result := false

			switch token.Type {
			case tokens.And:
				if !condLeft {
					result = false
				} else {
					result = condRight
				}
			case tokens.Or:
				if condLeft {
					result = true
				} else {
					result = condRight
				}
			default:
				err.SyntaxError(token, fmt.Sprintf("unsupported logical operator: %s", token.Literal), p.lines)
				p.errorHandler()
				return
			}

			stack = append(stack, tokens.Token{
				Type:    tokens.Bool,
				Literal: result,
				Loc:     token.Loc,
			})

			p.consume()

		case tokens.Write, tokens.Writeln:
			if len(stack) == 0 {
				err.SyntaxError(token, fmt.Sprintf(
					"The keyword '%s' requires value in stack. Stack is empty.",
					token.Literal),
					p.lines)
				p.errorHandler()
				return
			}

			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			p.consume()

			str := toString(a)
			outputBuffer.WriteString(str)
			if token.Type == tokens.Writeln {
				outputBuffer.WriteByte('\n')
			}

			if p.inLoop > 0 || p.isREPL {
				outputBuffer.Flush()
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

		case tokens.Swap:
			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			p.consume()

			stack[len(stack)-1], stack[len(stack)-2] = stack[len(stack)-2], stack[len(stack)-1]

		case tokens.Pop:
			var e error

			stack, _, e = Pop(stack)
			if e != nil {
				err.SyntaxError(token, fmt.Sprintf("The keyword '%s' requires value in stack. Stack is empty.", token.Literal), p.lines)
				p.errorHandler()
				return
			}

			p.consume()

		case tokens.Over:
			if len(stack) < 2 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires two operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			p.consume()

			v := stack[len(stack)-2]

			stack = append(stack, v)

		case tokens.Rot:
			if len(stack) < 3 {
				err.SyntaxError(token, fmt.Sprintf(
					"The '%s' operator requires three operands in stack. Found %d.",
					token.Literal, len(stack)),
					p.lines)
				p.errorHandler()
				return
			}

			p.consume()

			n := len(stack)
			a := stack[n-3]
			b := stack[n-2]
			c := stack[n-1]
			stack[n-3], stack[n-2], stack[n-1] = b, c, a

		case tokens.Depth:
			p.consume()
			stack = append(stack, tokens.Token{
				Type:    tokens.Int,
				Literal: len(stack),
				Loc:     token.Loc,
			})

		case tokens.Dump:
			p.consume()
			fmt.Printf("Stack[%d]:\n", len(stack))
			for i, v := range stack {
				fmt.Printf("  %d: (%s) %v", i, strings.ToLower(string(v.Type)), v.Literal)
				if i == len(stack)-1 {
					fmt.Printf("  <- top")
				}
				fmt.Println()
			}

		case tokens.Clear:
			p.consume()
			stack = stack[:0]

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
			if token.Type == tokens.For {
				p.inLoop++
			}
			p.consume()

		case tokens.Elif,
			tokens.Else,
			tokens.Define:
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

			if p.inLoop > 0 {
				p.inLoop--
			}

			p.consume()

		case tokens.Identifier:
			// TODO: Add suport for identifiers
			err.SyntaxError(token, fmt.Sprintf("Name '%s' is not defined.", token.Literal), p.lines)
			p.errorHandler()
			p.consume()

		default:
			err.Error(fmt.Sprintf("Not implemented case for TokenType '%s'.", token.Type))
			fmt.Fprintf(os.Stderr, "\x1b[31m%s:%d:%d:\x1b[0m '%s'", token.Loc.File, token.Loc.Line, token.Loc.Col, token.Literal)
			p.errorHandler()
			return
		}
	}

	outputBuffer.Flush()
}
